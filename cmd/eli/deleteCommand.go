package main

import (
	"github.com/urfave/cli"
)

var deleteCommand = cli.Command{
	Name:        "delete",
	HelpName:    "delete",
	Usage:       `Delete one or more resources`,
	Description: "With this command you can delete resources",
	ArgsUsage: `eli delete RESOURCE [options]

	 # Delete all running pods
	 eli delete pods

	 # Delete all 'my-pod' pod
	 eli delete pod my-pod`,
	Subcommands: []cli.Command{
		deletePodCommand,
	},
}
