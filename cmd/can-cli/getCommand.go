package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var getCommand = cli.Command{
	Name:        "get",
	HelpName:    "get",
	Usage:       `Display one or more resources`,
	Description: "With this command you can list resources",
	ArgsUsage: `can-cli get RESOURCE [options]

	 # Get table of running pods
	 can-cli get pods`,
	Subcommands: []cli.Command{
		{
			Name:    "pods",
			Aliases: []string{"pod"},
			Usage:   "Get Pod resources",
			UsageText: `can-cli get pods [options]
			 
	 # Get table of running pods
	 can-cli get pods`,
			Action: func(clicontext *cli.Context) error {
				client := cmd.GetClient(clicontext)

				pods, err := client.GetPods()
				if err != nil {
					return err
				}

				writer := printers.GetNewTabWriter(os.Stdout)
				printer := cmd.GetPrinter(clicontext)
				return printer.PrintPodsTable(pods, writer)
			},
		},
	},
}
