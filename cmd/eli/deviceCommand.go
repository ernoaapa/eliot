package main

import (
	"github.com/urfave/cli"
)

var deviceCommand = cli.Command{
	Name:        "device",
	HelpName:    "device",
	Usage:       `Device management commands`,
	Description: "Manage device OS, etc. with device commands",
	ArgsUsage: `eli device <COMMAND> [options]

	 # build device image
	 eli device pods`,
	Subcommands: []cli.Command{
		deviceBuildCommand,
	},
}
