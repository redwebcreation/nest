# Nest

[![codebeat badge](https://codebeat.co/badges/7171e9ea-53d7-4c81-82bf-a9a2f222b027)](https://codebeat.co/projects/github-com-redwebcreation-nest-next)
[![Go Report Card](https://goreportcard.com/badge/github.com/redwebcreation/nest)](https://goreportcard.com/report/github.com/redwebcreation/nest)
[![codecov](https://codecov.io/gh/redwebcreation/nest/branch/next/graph/badge.svg?token=DWSP4O0YO8)](https://codecov.io/gh/redwebcreation/nest)
![PRs not welcome](https://img.shields.io/badge/PRs-not%20welcome-red)

#### Documentation Status

The goal is to write a lot and then eventually make it more concise and improve upon it.

**VERY MUCH WIP, JUST RANDOM THINGS**

## Requirements

* docker
* git

## What is Nest?

Nest is a tool to help you manage many applications (called "services" from now on) on a single server. You can think of
it as a supercharged docker-compose.

Features:

* zero downtime deployments
* built-in reverse proxy
* versioned configuration
* powerful configuration diagnosis (if anything looks wrong in your configuration, nest will SCREAM LOUDLY)
* an api to trigger deployments automatically (CD [What's Continous Deployment (link needed)]() with a single api call)

## Why use Nest?

Nest is the perfect middle ground between messy configuration files all over your server and kubernetes.

## When not to use Nest?

* You have more than two servers

  If you have exactly two servers, you can still use nest very effectively and make your architecture redundant by
  running them in a Active-Active configuration (or Active-Passive if one is less powerful)
  . [(link needed)]()

## Concepts

### Config Locator

### Vision

Your configuration should be stored in a single place and versioned.

### Implementation

Your configuration is stored in a git repository. It can either be a local repository (not implemented yet, on roadmap)
or a remote repository. The remote option is preferred.

When first running nest, you must configure the config locator, the algorithm that will retrieve your configuration.

You may do so by running `nest configure` or alternatively `nest rcfg` (rcfg means reconfigure).

`git` must be installed on your system. If your configuration is stored remotely, you must be able to clone the
repository using SSH.

* HTTPS is not supported.
* Only GitHub, GitLab and Bitbucket are supported. (not Github Enterprise or self-hosted Gitlab, this may change in the
  future)

## Contributing

### Creating a new command

Let's say you want to create a new command called `version` that prints out Nest's current version.

```go
// command/version.go
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

// An helpful description for the command, if you don't know what to write, leave it blank
// so that it can be spotted by linters
// VersionCommand prints out the current version
func /* New[command]Command */ NewVersionCommand() *cobra.Command {
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
	// If you need to access the config or the config locator, use the following:
	// return WithConfig(cmd)
}
```

> `WithConfig` loads the config locator, downloads the config if needed and parses it. If there's any error in the config detected by the `medic`, it'll exit and ask the user to run `nest medic` to troubleshoot the problem.