package command

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/redwebcreation/nest/global"
	"github.com/spf13/cobra"
)

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	var versionNeeded string

	if len(args) > 0 {
		versionNeeded = args[0]
	}

	release, err := global.Repository.Release(versionNeeded)
	if err != nil {
		return err
	}

	if release.TagName == global.Version {
		return fmt.Errorf("already using this version")
	}

	binary := release.Assets[0]

	if binary.State != "uploaded" {
		return fmt.Errorf("binary is being uploaded, retry in a few seconds")
	}

	fmt.Printf("Downloading %s\n", binary.BrowserDownloadUrl)

	response, err := http.Get(binary.BrowserDownloadUrl)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	err = os.WriteFile(executable+"_updated", body, 0600)
	if err != nil {
		return err
	}

	err = os.Rename(executable+"_updated", executable)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully updated to the %s.\n", release.TagName)

	return nil
}

// NewSelfUpdateCommand updates nest to its latest version
func NewSelfUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-update [version]",
		Short: "update the CLI to the latest version",
		Args:  cobra.RangeArgs(0, 1),
		RunE:  runSelfUpdate,
	}

	return cmd
}
