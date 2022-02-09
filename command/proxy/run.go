package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"os"
)

var httpPort string
var httpsPort string

func runRunCommand(cmd *cobra.Command, args []string) error {
	config, err := pkg.Locator.Resolve()
	if err != nil {
		return err
	}

	contents, err := os.ReadFile(global.ContainerManifestFile)
	if err != nil {
		return err
	}

	var manifest pkg.Manifest
	err = json.Unmarshal(contents, &manifest)
	if err != nil {
		return fmt.Errorf("error parsing manifest: %s", err)
	}

	pkg.NewProxy(httpPort, httpsPort, config.Services, &manifest).Run()

	return nil
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
