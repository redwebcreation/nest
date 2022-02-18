package nest

import (
	"context"
	"fmt"
	"github.com/google/go-github/v42/github"
	"github.com/redwebcreation/nest/build"
	context2 "github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

type selfUpdateOptions struct {
	version string
}

func runSelfUpdate(ctx *context2.Context, opts *selfUpdateOptions) error {
	client := github.NewClient(nil)

	var release *github.RepositoryRelease
	var err error

	if opts.version != "" {
		release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "redwebcreation", "nest", opts.version)
	} else {
		release, _, err = client.Repositories.GetLatestRelease(context.Background(), "redwebcreation", "nest")
	}

	if err != nil {
		return err
	}

	if release.GetTagName() == build.Version {
		return fmt.Errorf("you are already using the latest version of nest")
	}

	binary := release.Assets[0]

	if binary.GetState() != "uploaded" {
		return fmt.Errorf("the binary is not available yet, try again later")
	}

	fmt.Fprintf(ctx.Out(), "Downloading %s...\n", binary.GetName())

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	err = download(binary.GetBrowserDownloadURL(), executable+".tmp")
	if err != nil {
		return err
	}

	err = os.Rename(executable+".tmp", executable)
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.Out(), "Successfully updated to version %s.\n", release.GetTagName())

	return nil
}

// NewSelfUpdateCommand creates a new `self-update` command.
func NewSelfUpdateCommand(ctx *context2.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-update [version]",
		Short: "update the CLI to the latest version",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &selfUpdateOptions{}

			if len(args) > 0 {
				opts.version = args[0]
			}

			return runSelfUpdate(ctx, opts)
		},
	}

	return cmd
}

func download(remote string, local string) error {
	response, err := http.Get(remote)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(local, body, 0600)
}
