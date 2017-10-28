package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var describeCommand = cli.Command{
	Name:        "describe",
	HelpName:    "describe",
	Usage:       "View details of resource",
	Description: "With this command you can get details of resources",
	ArgsUsage: `canctl describe RESOURCE [options] [POD NAME]
	
	# Describe a pod
	canctl describe pod my-pod-name

	# Describe all pods
	canctl describe pods
`,
	Subcommands: []cli.Command{
		{
			Name:    "pod",
			Aliases: []string{"pods"},
			Usage:   "Return details of pod",
			UsageText: `canctl describe RESOURCE [options] [POD NAME]
	
	# Describe a pod
	canctl describe pod my-pod-name

	# Describe all pods
	canctl describe pods
`,
			Action: func(clicontext *cli.Context) error {
				config := cmd.GetConfigProvider(clicontext)
				client := cmd.GetClient(config)

				podName := clicontext.Args().First()

				pods, err := client.GetPods()
				if err != nil {
					return err
				}

				if podName != "" {
					pods = cmd.FilterByPodName(pods, podName)
				}

				writer := printers.GetNewTabWriter(os.Stdout)
				defer writer.Flush()
				printer := cmd.GetPrinter(clicontext)
				for _, pod := range pods {
					if err := printer.PrintPodDetails(pod, writer); err != nil {
						return err
					}
				}
				return nil
			},
		},
	},
}
