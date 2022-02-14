package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"os"
)

var config *pkg.Configuration

func runRunCommand(cmd *cobra.Command, args []string) error {
	var manifest pkg.Manifest
	contents, err := os.ReadFile(global.GetContainerManifestFile())
	if err == nil {
		err = json.Unmarshal(contents, &manifest)
		if err != nil {
			return fmt.Errorf("error parsing manifest: %s", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	pkg.NewProxy(config, &manifest).Run()

	return nil
}

// NewRunCommand starts the reverse proxy
func NewRunCommand() *cobra.Command {
	resolvedConfig, err := pkg.Locator.Resolve()

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Starts the proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err != nil {
				return err
			}

			return runRunCommand(cmd, args)
		},
	}

	cmd.Flags().StringVar(&resolvedConfig.Proxy.Http, "http", resolvedConfig.Proxy.Http, "HTTP port")
	cmd.Flags().StringVar(&resolvedConfig.Proxy.Https, "https", resolvedConfig.Proxy.Https, "HTTPS port")
	cmd.Flags().BoolVarP(&resolvedConfig.Proxy.SelfSigned, "self-signed", "s", resolvedConfig.Proxy.SelfSigned, "Use a self-signed certificate")

	config = resolvedConfig

	return cmd
}
