package main

import (
	"github.com/urfave/cli"
)

var describeCommand = cli.Command{
	Name:        "describe",
	HelpName:    "describe",
	Usage:       "View details of resource",
	Description: "With this command you can get details of resources",
	ArgsUsage: `canctl describe RESOURCE [options] [POD NAME]
	
	# Describe a pod
	canctl describe pod my-pod-name

	# Describe all pods
	canctl describe pods
`,
	Subcommands: []cli.Command{
		describePodCommand,
	},
}
