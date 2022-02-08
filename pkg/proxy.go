package pkg

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/pseidemann/finish"
	"github.com/redwebcreation/nest/global"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

type Proxy struct {
	Http               string
	Https              string
	Logger             *Logger
	Services           ServiceMap
	CertificateManager *autocert.Manager
	HostToIp           map[string]string
}

func NewProxy() *Proxy {
	proxy := &Proxy{
		Http:     "80",
		Https:    "443",
		Logger:   ProxyLogger,
		HostToIp: make(map[string]string),
	}

	proxy.UpdateServicesHostAndIps()

	proxy.CertificateManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			config, err := Config.Resolve()
			if err != nil {
				return err
			}

			contents, err := os.ReadFile(global.ContainerManifestFile)
			if err != nil {
				return err
			}

			var manifest Manifest
			err = json.Unmarshal(contents, &manifest)
			if err != nil {
				return err
			}

			for _, service := range config.Services {
				for _, comparison := range service.Hosts {
					if comparison == host {
						return nil
					}
				}
			}

			return nil
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
		Log:     p.Logger,
	}

	httpToHttps := p.newRedirector()

	finisher.Add(httpToHttps)
	finisher.Add(proxy)

	go func() {
		err := httpToHttps.ListenAndServe()
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

func (p *Proxy) UpdateServicesHostAndIps() error {
	config, err := Config.Resolve()
	if err != nil {
		return err
	}

	contents, err := os.ReadFile(global.ContainerManifestFile)
	if err != nil {
		return err
	}

	var manifest *Manifest
	err = json.Unmarshal(contents, &manifest)
	if err != nil {
		return err
	}

	p.Services = config.Services

	for _, service := range p.Services {
		for _, host := range service.Hosts {
			p.HostToIp[host] = manifest.Containers[service.Name][0].IP
		}
	}

	return nil
}

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	ip := p.HostToIp[r.Host]

	if ip == "" {
		err := p.UpdateServicesHostAndIps()
		if err != nil {
			p.Logger.Error(err)
		}

		ip = p.HostToIp[r.Host]

		if ip == "" {
			p.Log(r, LevelInfo, "host not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   ip,
	}).ServeHTTP(w, r)

	p.Log(r, LevelInfo, "request proxied")
}

func (p *Proxy) Log(r *http.Request, level Level, message string) {
	var ip string

	if ip = r.Header.Get("X-Forwarded-For"); ip == "" {
		ip = r.RemoteAddr
	}

	p.Logger.Log(
		level,
		message,
		NewField("host", r.Host),
		NewField("method", r.Method),
		NewField("path", r.URL.Path),
		NewField("ip", ip),
	)
}

func (p *Proxy) newRedirector() *http.Server {
	return &http.Server{
		Addr: ":" + p.Http,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p.CertificateManager.HTTPHandler(nil).ServeHTTP(w, r)

			p.Log(r, LevelInfo, "redirecting to https")
		}),
	}
}
