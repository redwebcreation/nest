package pkg

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/pseidemann/finish"
	"github.com/redwebcreation/nest/global"
	"golang.org/x/crypto/acme/autocert"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

type Proxy struct {
	Config             *Configuration
	CertificateManager *autocert.Manager
	hostToIp           map[string]string
}

func NewProxy(config *Configuration, manifest *Manifest) *Proxy {
	proxy := &Proxy{
		Config:   config,
		hostToIp: make(map[string]string),
	}

	// todo: get container ip
	//for _, service := range config.Services {
	//	for _, host := range service.Hosts {
	//		//proxy.hostToIp[host] = manifest.Containers[service.Name].IP
	//	}
	//}

	proxy.CertificateManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if _, ok := proxy.hostToIp[host]; ok {
				return nil
			}

			return fmt.Errorf("acme/autocert: host %s not configured", host)
		},
		Cache: autocert.DirCache(global.GetCertsDir()),
	}

	return proxy
}

func (p *Proxy) Run() {
	server := p.newServer(p.Config.Proxy.Https, func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "" && r.Host == p.Config.ControlPlane.Host {
			p.Log(r, global.LevelInfo, "proxied request to plane")

			NewRouter(p.Config).ServeHTTP(w, r)
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

			certFile := global.GetSelfSignedCertFile()
			keyFile := global.GetSelfSignedCertKeyFile()

			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if errors.Is(err, os.ErrNotExist) {
				err = createSelfSignedCertificates(certFile, keyFile)
				if err != nil {
					global.LogP(global.LevelError, "failed to create self signed certificates", global.Fields{
						"error": err.Error(),
					})
					return nil, err
				}

				cert, err = tls.LoadX509KeyPair(certFile, keyFile)
				if err != nil {
					global.LogP(global.LevelError, "failed to load self signed certificates", global.Fields{
						"error": err.Error(),
					})
					return nil, err
				}
			} else if err != nil {
				global.LogP(global.LevelError, "certificates exist but failed to load", global.Fields{
					"error": err.Error(),
				})
				return nil, err
			}

			return &cert, nil
		},
	}

	p.start(server)
}

func (p *Proxy) start(proxy *http.Server) {
	finisher := &finish.Finisher{
		Timeout: 10 * time.Second,
		Log: &global.FinisherLogger{
			Logger: global.ProxyLogger,
		},
	}

	certsHandler := p.certsCreationHandler()

	finisher.Add(certsHandler, finish.WithName("http"))
	finisher.Add(proxy, finish.WithName("https"))

	go func() {
		err := certsHandler.ListenAndServe()
		if err != nil {
			global.LogP(global.LevelFatal, err.Error(), nil)
			os.Exit(1)
		}
	}()

	go func() {
		err := proxy.ListenAndServeTLS("", "")
		if err != nil {
			global.LogP(global.LevelFatal, err.Error(), nil)
			os.Exit(1)
		}
	}()

	finisher.Wait()
}

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	ip := p.hostToIp[r.Host]

	if ip == "" {
		p.Log(r, global.LevelInfo, "host not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   ip,
	}).ServeHTTP(w, r)

	p.Log(r, global.LevelInfo, "request proxied")
}

func (p *Proxy) Log(r *http.Request, level global.Level, message string) {
	var ip string

	if ip = r.Header.Get("X-Forwarded-For"); ip == "" {
		ip = r.RemoteAddr
	}

	global.LogP(
		level,
		message,
		global.Fields{
			"tag":    "proxy",
			"method": r.Method,
			"host":   r.Host,
			"path":   r.URL.Path,
			"ip":     ip,
		},
	)
}

func (p *Proxy) certsCreationHandler() *http.Server {
	return p.newServer(p.Config.Proxy.Http, func(w http.ResponseWriter, r *http.Request) {
		p.CertificateManager.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" && r.Method != "HEAD" {
				http.Error(w, "Use HTTPS", http.StatusBadRequest)
				return
			}

			target := "https://" + replacePort(r.Host, p.Config.Proxy.Https) + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusFound)
		})).ServeHTTP(w, r)

		p.Log(r, global.LevelInfo, "redirecting to https")
	})
}

func (p *Proxy) newServer(port string, handler http.HandlerFunc) *http.Server {
	return &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
		ErrorLog:       global.ProxyLogger,
	}
}

func replacePort(url string, newPort string) string {
	host, _, err := net.SplitHostPort(url)
	if err != nil {
		return url
	}
	return net.JoinHostPort(host, newPort)
}

func createSelfSignedCertificates(certFile string, keyFile string) error {
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
