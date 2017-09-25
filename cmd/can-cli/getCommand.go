package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var getCommand = cli.Command{
	Name: "get",
	Subcommands: []cli.Command{
		{
			Name:  "pods",
			Usage: "List pods",
			Action: func(clicontext *cli.Context) error {
				client := cmd.GetClient(clicontext)

				pods, err := client.GetPods()
				if err != nil {
					return err
				}

				writer := printers.GetNewTabWriter(os.Stdout)
				writerErr := printers.NewHumanReadablePrinter().PrintPods(pods, writer)
				return writerErr
			},
		},
	},
}
