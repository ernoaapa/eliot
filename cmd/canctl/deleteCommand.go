package main

import (
	"github.com/urfave/cli"
)

var deleteCommand = cli.Command{
	Name:        "delete",
	HelpName:    "delete",
	Usage:       `Delete one or more resources`,
	Description: "With this command you can delete resources",
	ArgsUsage: `canctl delete RESOURCE [options]

	 # Delete all running pods
	 canctl delete pods

	 # Delete all 'my-pod' pod
	 canctl delete pod my-pod`,
	Subcommands: []cli.Command{
		deletePodCommand,
	},
}
