package main

import (
	"github.com/urfave/cli"
)

var describeCommand = cli.Command{
	Name:        "describe",
	HelpName:    "describe",
	Usage:       "View details of resource",
	Description: "With this command you can get details of resources",
	ArgsUsage: `eli describe RESOURCE [options] [POD NAME]
	
	# Describe a pod
	eli describe pod my-pod-name

	# Describe all pods
	eli describe pods
`,
	Subcommands: []cli.Command{
		describePodCommand,
		describeDeviceCommand,
	},
}
