package main

import (
	"os"
	"time"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/cmd/ui"
	"github.com/ernoaapa/eliot/pkg/discovery"
	"github.com/ernoaapa/eliot/pkg/printers"
	"github.com/urfave/cli"
)

var getDevicesCommand = cli.Command{
	Name:    "devices",
	Aliases: []string{"device"},
	Usage:   "Get Device resources",
	UsageText: `eli get devices [options]
			 
	 # Get table of known devices
	 eli get devices`,
	Action: func(clicontext *cli.Context) error {
		uiline := ui.NewLine().Loading("Discover from network automatically...")

		devices, err := discovery.Devices(5 * time.Second)
		if err != nil {
			uiline.Fatalf("Failed to auto-discover devices in network: %s", err)
		}
		uiline.Donef("Discovered %d devices from network", len(devices))

		writer := printers.GetNewTabWriter(os.Stdout)
		defer writer.Flush()

		printer := cmd.GetPrinter(clicontext)
		return printer.PrintDevicesTable(devices, writer)
	},
}
