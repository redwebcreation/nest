package config

import (
	"fmt"
	"sort"

	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

func runShowCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("strategy:", pkg.Config.Strategy)
	fmt.Println("location:", pkg.Config.GetRemoteURL())
	fmt.Printf("current commit: %s\n", pkg.Config.Commit[:7])
	fmt.Println("branch:", pkg.Config.Branch)
	if pkg.Config.Dir != "" {
		fmt.Println("subdir:", pkg.Config.Dir)
	}

	repo, err := pkg.Config.LocalClone()
	if err != nil {
		return err
	}

	configFiles, err := repo.Tree()
	if err != nil {
		return err
	}

	sort.Slice(configFiles, func(i, j int) bool {
		return configFiles[i] < configFiles[j]
	})

	fmt.Println("\nfiles:")
	for _, file := range configFiles {
		fmt.Printf("- %s\n", file)
	}

	return nil
}

// NewShowCommand prints the current configuration for the config locator
func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "prints the current configuration",
		RunE:  runShowCommand,
	}

	return cmd
}
