package proxy

import (
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

var config *pkg.ServerConfiguration

var http string
var https string
var selfSigned bool

type runOptions struct {
	deployment string
}

func runRunCommand(ctx *pkg.Context, opts *runOptions) error {
	// update configuration according to flags
	config.Proxy.HTTP = http
	config.Proxy.HTTPS = https
	config.Proxy.SelfSigned = selfSigned

	var manifest *pkg.Manifest
	var err error

	if opts.deployment != "" {
		manifest, err = pkg.LoadManifest(opts.deployment)
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

// NewRunCommand creates a new `run` command
func NewRunCommand(ctx *pkg.Context) *cobra.Command {
	resolvedConfig, err := ctx.ServerConfiguration()

	cmd := &cobra.Command{
		Use:   "run [deployment]",
		Short: "Starts the proxy",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err != nil {
				return err
			}

			opts := &runOptions{}

			if len(args) > 0 {
				opts.deployment = args[0]
			}

			return runRunCommand(ctx, opts)
		},
	}

	if err == nil {
		http = resolvedConfig.Proxy.HTTP
		https = resolvedConfig.Proxy.HTTPS
		selfSigned = resolvedConfig.Proxy.SelfSigned
	}

	cmd.Flags().StringVar(&http, "http", http, "HTTP port")
	cmd.Flags().StringVar(&https, "https", https, "HTTPS port")
	cmd.Flags().BoolVarP(&selfSigned, "self-signed", "s", selfSigned, "Use a self-signed certificate")

	config = resolvedConfig

	return cmd
}