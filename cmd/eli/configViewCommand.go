package main

import (
	"os"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/printers"
	"github.com/urfave/cli"
)

var configViewCommand = cli.Command{
	Name:      "view",
	Usage:     "View client config",
	UsageText: "eli config view",
	Action: func(clicontext *cli.Context) error {
		config := cmd.GetConfig(clicontext)

		writer := printers.GetNewTabWriter(os.Stdout)
		defer writer.Flush()
		printer := cmd.GetPrinter(clicontext)
		return printer.PrintConfig(config, writer)
	},
}
