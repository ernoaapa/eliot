package main

import (
	"github.com/urfave/cli"
)

var getCommand = cli.Command{
	Name:        "get",
	HelpName:    "get",
	Usage:       `Display one or more resources`,
	Description: "With this command you can list resources",
	ArgsUsage: `canctl get RESOURCE [options]

	 # Get table of running pods
	 canctl get pods`,
	Subcommands: []cli.Command{
		getPodsCommand,
		getDevicesCommand,
	},
}
