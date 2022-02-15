package proxy

import (
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

var config *pkg.Configuration

var http string
var https string
var selfSigned bool

func runRunCommand(cmd *cobra.Command, args []string) error {
	// update configuration according to flags
	config.Proxy.Http = http
	config.Proxy.Https = https
	config.Proxy.SelfSigned = selfSigned

	var manifest *pkg.Manifest
	var err error

	if len(args) > 0 {
		manifest, err = pkg.LoadManifest(args[0])
		if err != nil {
			return err
		}
	} else {
		manifest, err = pkg.GetLatestManifest()
		if err != nil {
			return err
		}
	}

	pkg.NewProxy(config, manifest).Run()

	return nil
}

// NewRunCommand starts the reverse proxy
func NewRunCommand() *cobra.Command {
	resolvedConfig, err := pkg.Locator.Resolve()

	cmd := &cobra.Command{
		Use:   "run [deployment]",
		Short: "Starts the proxy",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err != nil {
				return err
			}

			return runRunCommand(cmd, args)
		},
	}

	if err == nil {
		http = resolvedConfig.Proxy.Http
		https = resolvedConfig.Proxy.Https
		selfSigned = resolvedConfig.Proxy.SelfSigned
	}

	cmd.Flags().StringVar(&http, "http", http, "HTTP port")
	cmd.Flags().StringVar(&https, "https", https, "HTTPS port")
	cmd.Flags().BoolVarP(&selfSigned, "self-signed", "s", selfSigned, "Use a self-signed certificate")

	config = resolvedConfig

	return cmd
}
