package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/me/nest/global"
	"github.com/spf13/cobra"
)

func SelfUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-update [version]",
		Short: "update the CLI to the latest version",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var versionNeeded string

			if len(args) > 0 {
				versionNeeded = args[0]
			}

			releases, err := global.Repository.Releases(versionNeeded)
			if err != nil {
				return err
			}

			if len(releases) == 0 {
				return fmt.Errorf("no releases found")
			}

			latestRelease := releases[0]

			if latestRelease.TagName == global.Version {
				fmt.Printf("You are already using the latest version of the CLI.\n")
				return nil
			}

			binary := latestRelease.Assets[0]

			if binary.State != "uploaded" {
				fmt.Println("The binary for this release is still being uploaded.")
				fmt.Println("Please try again in a few seconds.")
				return nil
			}

			fmt.Printf("Downloading %s;\n", binary.BrowserDownloadUrl)

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

			err = os.WriteFile(executable+"_latest", body, 0755)
			if err != nil {
				return err
			}

			err = os.Rename(executable+"_latest", executable)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully updated to the %s.\n", latestRelease.TagName)

			return nil
		},
	}

	return cmd
}
