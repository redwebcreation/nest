package proxy

import (
	"github.com/redwebcreation/nest/context"
	"github.com/redwebcreation/nest/deploy"
	"github.com/redwebcreation/nest/proxy"
	"github.com/spf13/cobra"
)

type runOptions struct {
	deployment string
	HTTP       string
	HTTPS      string
	selfSigned bool
}

func runRunCommand(ctx *context.Context, opts *runOptions) error {

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

	config, err := ctx.ServicesConfig()
	if err != nil {
		return err
	}

	config.Proxy.HTTP = opts.HTTP
	config.Proxy.HTTPS = opts.HTTPS
	config.Proxy.SelfSigned = opts.selfSigned

	proxy.NewProxy(ctx, config, manifest).Run()
	return err
}

// NewRunCommand creates a new `run` command
func NewRunCommand(ctx *context.Context) *cobra.Command {
	opts := &runOptions{}

	cmd := &cobra.Command{
		Use:   "run [deployment]",
		Short: "Starts the proxy",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			servicesConfig, err := ctx.ServicesConfig()
			if err != nil {
				return err
			}

			if opts.HTTP == "" {
				opts.HTTP = servicesConfig.Proxy.HTTP
			}

			if opts.HTTPS == "" {
				opts.HTTPS = servicesConfig.Proxy.HTTPS
			}

			if !opts.selfSigned && servicesConfig.Proxy.SelfSigned {
				opts.selfSigned = true
			}

			if len(args) > 0 {
				opts.deployment = args[0]
			}

			return runRunCommand(ctx, opts)
		},
	}

	cmd.Flags().StringVar(&opts.HTTP, "http", "", "HTTP port")
	cmd.Flags().StringVar(&opts.HTTPS, "https", "", "HTTPS port")
	cmd.Flags().BoolVarP(&opts.selfSigned, "self-signed", "u", false, "Use a self-signed certificate")

	return cmd
}
