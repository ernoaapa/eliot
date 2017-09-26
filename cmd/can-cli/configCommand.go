package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/printers"
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
		{
			Name:      "view",
			Usage:     "View client config",
			UsageText: "can-cli config view",
			Action: func(clicontext *cli.Context) error {
				config := cmd.GetConfig(clicontext)

				writer := printers.GetNewTabWriter(os.Stdout)
				printer := cmd.GetPrinter(clicontext)
				return printer.PrintConfig(config, writer)
			},
		},
	},
}
