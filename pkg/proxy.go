package pkg

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/pseidemann/finish"
	"github.com/redwebcreation/nest/global"
	"golang.org/x/crypto/acme/autocert"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
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

			// todo: generate self-signed certificate using golang
			certFile := global.GetCertsDir() + "/dev_cert.pem"
			keyFile := global.GetCertsDir() + "/dev_key.pem"

			if _, err := os.Stat(keyFile); os.IsNotExist(err) {
				cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048", "-keyout", keyFile, "-out", certFile, "-days", "365", "-nodes", "-subj", "/CN=localhost")
				err = cmd.Run()
				if err != nil {
					return nil, err
				}
			} else if err != nil {
				return nil, err
			}

			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
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
