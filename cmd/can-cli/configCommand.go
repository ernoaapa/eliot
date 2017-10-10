package main

import (
	"github.com/urfave/cli"
)

var configCommand = cli.Command{
	Name:     "config",
	HelpName: "config",
	Usage:    `View and edit client configuration`,
	Description: `With this command you view and edit the client configurations
	 like device address, username, namespace, etc.`,
	ArgsUsage: "can-cli config view",
	Subcommands: []cli.Command{
		configViewCommand,
		configSetCommand,
	},
}
