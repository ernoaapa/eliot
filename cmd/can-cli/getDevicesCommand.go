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
	UsageText: `can-cli get devices [options]
			 
	 # Get table of known devices
	 can-cli get devices`,
	Action: func(clicontext *cli.Context) error {
		devices := make(chan model.DeviceInfo)

		writer := printers.GetNewTabWriter(os.Stdout)
		printer := cmd.GetPrinter(clicontext)
		printer.PrintDevicesTable(devices, writer)

		return discovery.Devices(devices, 4*time.Second)
	},
}
