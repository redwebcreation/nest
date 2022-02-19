package proxy

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/pseidemann/finish"
	"github.com/redwebcreation/nest/config"
	context2 "github.com/redwebcreation/nest/context"
	"github.com/redwebcreation/nest/deploy"
	"github.com/redwebcreation/nest/loggy"
	"github.com/redwebcreation/nest/proxy/plane"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

type Proxy struct {
	Ctx                *context2.Context
	Config             *config.ServicesConfig
	CertificateManager *autocert.Manager
	hostToIP           map[string]string
}

func NewProxy(ctx *context2.Context, servicesConfig *config.ServicesConfig, manifest *deploy.Manifest) *Proxy {
	proxy := &Proxy{
		Ctx:      ctx,
		Config:   servicesConfig,
		hostToIP: make(map[string]string),
	}

	// todo: get container ip
	//for _, service := range servicesConfig.Services {
	//	for _, host := range service.Hosts {
	//		proxy.hostToIP[host] = manifest.Containers[service.Name].IP
	//	}
	//}

	proxy.CertificateManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if _, ok := proxy.hostToIP[host]; ok {
				return nil
			}

			return fmt.Errorf("acme/autocert: host %s not configured", host)
		},
		Cache: ctx.CertificateStore(),
	}

	return proxy
}

func (p *Proxy) Run() {
	server := p.newServer(p.Config.Proxy.HTTPS, func(w http.ResponseWriter, r *http.Request) {
		if (r.Host != "" && r.Host == p.Config.ControlPlane.Host) || r.Host == "control-plane" {
			p.Log(r, loggy.InfoLevel, "proxied request to plane")

			plane.New(p.Ctx).ServeHTTP(w, r)
			return
		}

		p.handler(w, r)
	})
	server.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if !p.Config.Proxy.SelfSigned {
				return p.CertificateManager.GetCertificate(info)
			}

			keyFile := "dev_key.pem"
			certFile := "dev_cert.pem"

			keyPEM, err := p.Ctx.CertificateStore().Get(context.Background(), keyFile)
			if err != nil {
				return handleSelfSignedCertificates(keyFile, certFile)
			}

			certPEM, certErr := p.Ctx.CertificateStore().Get(context.Background(), certFile)
			if certErr != nil {
				return handleSelfSignedCertificates(keyFile, certFile)
			}

			cert, err := tls.X509KeyPair(certPEM, keyPEM)
			if err != nil {
				return nil, err
			}

			return &cert, nil
		},
	}

	p.start(server)
}

func handleSelfSignedCertificates(keyFile, certFile string) (*tls.Certificate, error) {
	err := generateSelfSignedCertificates(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (p *Proxy) start(proxy *http.Server) {
	finisher := &finish.Finisher{
		Timeout: 10 * time.Second,
		Log: &FinisherLogger{
			Logger: p.Ctx.ProxyLogger(),
		},
	}

	certsHandler := p.certsCreationHandler()

	finisher.Add(certsHandler, finish.WithName("http"))
	finisher.Add(proxy, finish.WithName("https"))

	go func() {
		err := certsHandler.ListenAndServe()
		if err != nil {
			p.Ctx.ProxyLogger().Print(loggy.NewEvent(loggy.FatalLevel, err.Error(), nil))
			os.Exit(1)
		}
	}()

	go func() {
		err := proxy.ListenAndServeTLS("", "")
		if err != nil {
			p.Ctx.ProxyLogger().Print(loggy.NewEvent(loggy.FatalLevel, err.Error(), nil))
			os.Exit(1)
		}
	}()

	finisher.Wait()
}

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	ip := p.hostToIP[r.Host]

	if ip == "" {
		p.Log(r, loggy.InfoLevel, "host not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   ip,
	}).ServeHTTP(w, r)

	p.Log(r, loggy.InfoLevel, "request proxied")
}

func (p *Proxy) Log(r *http.Request, level loggy.Level, message string) {
	var ip string

	if ip = r.Header.Get("X-Forwarded-For"); ip == "" {
		ip = r.RemoteAddr
	}

	p.Ctx.ProxyLogger().Print(loggy.NewEvent(
		level,
		message,
		loggy.Fields{
			"tag":    "proxy",
			"method": r.Method,
			"host":   r.Host,
			"path":   r.URL.Path,
			"ip":     ip,
		},
	))
}

func (p *Proxy) certsCreationHandler() *http.Server {
	return p.newServer(p.Config.Proxy.HTTP, func(w http.ResponseWriter, r *http.Request) {
		p.CertificateManager.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" && r.Method != "HEAD" {
				http.Error(w, "Use HTTPS", http.StatusBadRequest)
				return
			}

			target := "https://" + replacePort(r.Host, p.Config.Proxy.HTTPS) + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusFound)
		})).ServeHTTP(w, r)

		p.Log(r, loggy.InfoLevel, "redirecting to https")
	})
}

func (p *Proxy) newServer(port string, handler http.HandlerFunc) *http.Server {
	return &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
		ErrorLog:       p.Ctx.ProxyLogger(),
	}
}

func replacePort(url string, newPort string) string {
	host, _, err := net.SplitHostPort(url)
	if err != nil {
		return url
	}
	return net.JoinHostPort(host, newPort)
}

func generateSelfSignedCertificates(certFile string, keyFile string) error {
	// create self-signed certificate using x509 package
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// create a new template for certificate
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1000, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// generate certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// write key to file
	keyFileHandle, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer keyFileHandle.Close()

	err = pem.Encode(keyFileHandle, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return err
	}

	// write certificate to file
	certFileHandle, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer certFileHandle.Close()

	err = pem.Encode(certFileHandle, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return err
	}

	return nil

}

type FinisherLogger struct {
	Logger *log.Logger
}

func (l FinisherLogger) Infof(message string, args ...any) {
	l.Logger.Print(loggy.NewEvent(loggy.InfoLevel, fmt.Sprintf(message, args...), loggy.Fields{
		"tag": "proxy.finisher",
	}))
}

func (l FinisherLogger) Errorf(message string, args ...any) {
	l.Logger.Print(loggy.NewEvent(loggy.ErrorLevel, fmt.Sprintf(message, args...), loggy.Fields{
		"tag": "proxy.finisher",
	}))
}
