package main

import (
	"os"

	"github.com/ernoaapa/elliot/cmd"
	"github.com/ernoaapa/elliot/pkg/printers"
	"github.com/urfave/cli"
)

var describePodCommand = cli.Command{
	Name:    "pod",
	Aliases: []string{"pods"},
	Usage:   "Return details of pod",
	UsageText: `eli describe RESOURCE [options] [POD NAME]
	
	# Describe a pod
	eli describe pod my-pod-name

	# Describe all pods
	eli describe pods
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
}