package main

import (
	"os"
	"time"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/discovery"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var getDevicesCommand = cli.Command{
	Name:    "devices",
	Aliases: []string{"device"},
	Usage:   "Get Device resources",
	UsageText: `canctl get devices [options]
			 
	 # Get table of known devices
	 canctl get devices`,
	Action: func(clicontext *cli.Context) error {
		devices := make(chan model.DeviceInfo)
		defer close(devices)

		writer := printers.GetNewTabWriter(os.Stdout)
		defer writer.Flush()

		printer := cmd.GetPrinter(clicontext)
		printer.PrintDevicesTable(devices, writer)

		return discovery.DevicesAsync(devices, 5*time.Second)
	},
}
