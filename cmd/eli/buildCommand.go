package main

import (
	"github.com/urfave/cli"
)

var buildCommand = cli.Command{
	Name:        "build",
	HelpName:    "build",
	Usage:       `Build node image, etc.`,
	Description: "With build command, you can build node image, etc.",
	ArgsUsage: `eli build <RESOURCE> [options]

	 # build node image
	 eli build node`,
	Subcommands: []cli.Command{
		buildNodeCommand,
	},
}
