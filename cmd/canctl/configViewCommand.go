package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var configViewCommand = cli.Command{
	Name:      "view",
	Usage:     "View client config",
	UsageText: "canctl config view",
	Action: func(clicontext *cli.Context) error {
		config := cmd.GetConfig(clicontext)

		writer := printers.GetNewTabWriter(os.Stdout)
		defer writer.Flush()
		printer := cmd.GetPrinter(clicontext)
		return printer.PrintConfig(config, writer)
	},
}
