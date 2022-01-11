# Nest

[![codebeat badge](https://codebeat.co/badges/7171e9ea-53d7-4c81-82bf-a9a2f222b027)](https://codebeat.co/projects/github-com-redwebcreation-nest-next)
[![Go Report Card](https://goreportcard.com/badge/github.com/redwebcreation/nest)](https://goreportcard.com/report/github.com/redwebcreation/nest)
[![codecov](https://codecov.io/gh/redwebcreation/nest/branch/next/graph/badge.svg?token=DWSP4O0YO8)](https://codecov.io/gh/redwebcreation/nest)
![PRs not welcome](https://img.shields.io/badge/PRs-not%20welcome-red)

## Contributing

### Creating a new command

Let's say you want to create a new command called `version` that prints out Nest's current version.

* Create your command in `command/version.go`

```go
package commands

import (
	"fmt"
	"github.com/redwebcreation/nest/globals"
	"github.com/spf13/cobra"
)

// Arguments must always be cmd and args even if you don't need them
// Do not underscore them.
func runVersionCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("nest@%s\n", globals.Version)

	return nil
}

func VersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		// Short must be formatted like an error, first-letter should be lowercase and without a period.
		// Keep the wording simple. Don't make it longer than a few words. Don't be fancy.
		Short: "prints nest's version",
		// Always use RunE instead of Run. Even if you don't need to return an error.
		RunE: runVersionCommand,
	}

	// Do not return directly the command as it makes adding flags harder.
	return cmd
}
```
