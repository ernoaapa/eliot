package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var getPodsCommand = cli.Command{
	Name:    "pods",
	Aliases: []string{"pod"},
	Usage:   "Get Pod resources",
	UsageText: `can-cli get pods [options]
			 
	 # Get table of running pods
	 can-cli get pods`,
	Action: func(clicontext *cli.Context) error {
		config := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(config)

		pods, err := client.GetPods()
		if err != nil {
			return err
		}

		writer := printers.GetNewTabWriter(os.Stdout)
		printer := cmd.GetPrinter(clicontext)
		return printer.PrintPodsTable(pods, writer)
	},
}
