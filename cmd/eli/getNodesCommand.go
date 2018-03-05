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

var getNodesCommand = cli.Command{
	Name:    "nodes",
	Aliases: []string{"node", "devices"}, // devices is deprecated command
	Usage:   "Get Node resources",
	UsageText: `eli get nodes [options]
			 
	 # Get table of known nodes
	 eli get nodes`,
	Action: func(clicontext *cli.Context) error {
		uiline := ui.NewLine().Loading("Discover from network automatically...")

		nodes, err := discovery.Nodes(5 * time.Second)
		if err != nil {
			uiline.Fatalf("Failed to auto-discover nodes in network: %s", err)
		}
		uiline.Donef("Discovered %d nodes from network", len(nodes))

		writer := printers.GetNewTabWriter(os.Stdout)
		defer writer.Flush()

		printer := cmd.GetPrinter(clicontext)
		return printer.PrintNodes(nodes, writer)
	},
}
