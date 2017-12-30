package main

import (
	"github.com/urfave/cli"
)

var buildCommand = cli.Command{
	Name:        "build",
	HelpName:    "build",
	Usage:       `Build Device image, etc.`,
	Description: "With build command, you can build device image, etc.",
	ArgsUsage: `eli build <RESOURCE> [options]

	 # build device image
	 eli build device`,
	Subcommands: []cli.Command{
		buildDeviceCommand,
	},
}
