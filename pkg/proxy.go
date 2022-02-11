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
	"time"
)

type Proxy struct {
	Http               string
	Https              string
	Logger             global.Logger
	Services           ServiceMap
	CertificateManager *autocert.Manager
	HostToIp           map[string]string
}

func NewProxy(http string, https string, services ServiceMap, manifest *Manifest) *Proxy {
	proxy := &Proxy{
		Http:     http,
		Https:    https,
		Logger:   global.ProxyLogger,
		HostToIp: make(map[string]string),
	}

	for _, service := range services {
		for _, host := range service.Hosts {
			proxy.HostToIp[host] = manifest.Containers[service.Name][0].IP
		}
	}

	proxy.CertificateManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if _, ok := proxy.HostToIp[host]; ok {
				return nil
			}

			return fmt.Errorf("acme/autocert: host %s not configured", host)
		},
		Cache: autocert.DirCache(global.CertsDir),
	}

	return proxy
}

func (p *Proxy) Run() {
	server := &http.Server{
		Addr: ":" + p.Https,
		TLSConfig: &tls.Config{
			MinVersion:     tls.VersionTLS13,
			GetCertificate: p.CertificateManager.GetCertificate,
		},
		Handler: http.HandlerFunc(p.handler),
	}

	p.start(server)
}

func (p *Proxy) start(proxy *http.Server) {
	finisher := &finish.Finisher{
		Timeout: 10 * time.Second,
		Log: &global.LogrusCompat{
			Logger: p.Logger,
		},
	}

	certsHandler := p.certificateCreationHandler()

	finisher.Add(certsHandler)
	finisher.Add(proxy)

	go func() {
		err := certsHandler.ListenAndServe()
		if err != nil {
			p.Logger.Error(err)
			os.Exit(1)
		}
	}()

	go func() {
		err := proxy.ListenAndServeTLS("", "")
		if err != nil {
			p.Logger.Error(err)
			os.Exit(1)
		}
	}()

	finisher.Wait()
}

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	ip := p.HostToIp[r.Host]

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

	p.Logger.Log(
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

func (p *Proxy) certificateCreationHandler() *http.Server {
	return &http.Server{
		Addr: ":" + p.Http,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p.CertificateManager.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" && r.Method != "HEAD" {
					http.Error(w, "Use HTTPS", http.StatusBadRequest)
					return
				}

				target := "https://" + replacePort(r.Host, p.Https) + r.URL.RequestURI()
				http.Redirect(w, r, target, http.StatusFound)
			})).ServeHTTP(w, r)

			p.Log(r, global.LevelInfo, "redirecting to https")
		}),
	}
}

func replacePort(url string, newPort string) string {
	host, _, err := net.SplitHostPort(url)
	if err != nil {
		return url
	}
	return net.JoinHostPort(host, newPort)
}
