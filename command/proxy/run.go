package proxy

import (
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
		err = http.ListenAndServe(":80", certificateManager.HTTPHandler(nil))

		if err != nil {
			fmt.Println(err)
		}
	}()

	server := &http.Server{
		Addr: ":443",
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

	return cmd
}
