package proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var httpPort string
var httpsPort string

func runRunCommand(cmd *cobra.Command, args []string) error {
	config, err := pkg.Config.Resolve()
	if err != nil {
		return err
	}

	contents, err := os.ReadFile(global.ConfigHome + "/manifest.json")
	if err != nil {
		return err
	}

	var manifest *pkg.Manifest
	err = json.Unmarshal(contents, manifest)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		comparison, err := os.ReadFile(global.ConfigHome + "/manifest.json")
		if err == nil {
			if bytes.Compare(contents, comparison) != 0 {
				var newManifest *pkg.Manifest
				err = json.Unmarshal(comparison, newManifest)
				if err != nil {
					panic(err)
				}

				manifest = newManifest
			}
		}

		service := config.Services[r.Host]

		if service == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Load balancing would happen here
		container := manifest.Containers[service.Name][0]

		httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   container.Ip,
		}).ServeHTTP(w, r)
	})

	certificateManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if _, ok := config.Services[host]; ok {
				return nil
			}

			return fmt.Errorf("acme/autocert: host %s not configured", host)
		},
		Cache: autocert.DirCache(global.ConfigHome + "/certs"),
	}

	go func() {
		err = http.ListenAndServe(":"+httpPort, certificateManager.HTTPHandler(nil))

		if err != nil {
			fmt.Println(err)
		}
	}()

	server := &http.Server{
		Addr: ":" + httpsPort,
		TLSConfig: &tls.Config{
			MinVersion:     tls.VersionTLS13,
			GetCertificate: certificateManager.GetCertificate,
		},
		Handler: handler,
	}

	return server.ListenAndServeTLS("", "")
}

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Starts the proxy",
		RunE:  runRunCommand,
	}

	cmd.Flags().StringVar(&httpPort, "http", "80", "HTTP port")
	cmd.Flags().StringVar(&httpsPort, "https", "443", "HTTPS port")

	return cmd
}
