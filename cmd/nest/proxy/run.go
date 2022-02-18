package proxy

import (
	config2 "github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/context"
	"github.com/redwebcreation/nest/deploy"
	"github.com/redwebcreation/nest/proxy"
	"github.com/spf13/cobra"
)

var config *config2.ServicesConfig

var http string
var https string
var selfSigned bool

type runOptions struct {
	deployment string
}

func runRunCommand(ctx *context.Context, opts *runOptions) error {
	// update config according to flags
	config.Proxy.HTTP = http
	config.Proxy.HTTPS = https
	config.Proxy.SelfSigned = selfSigned

	var manifest *deploy.Manifest
	var err error

	if opts.deployment != "" {
		manifest, err = ctx.ManifestManager().LoadWithID(opts.deployment)
		if err != nil {
			return err
		}
	} else {
		manifest, err = ctx.ManifestManager().Latest()
		if err != nil {
			return err
		}
	}

	proxy.NewProxy(ctx, config, manifest).Run()
	return err
}

// NewRunCommand creates a new `run` command
func NewRunCommand(ctx *context.Context) *cobra.Command {
	resolvedConfig, err := ctx.ServicesConfig()

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
