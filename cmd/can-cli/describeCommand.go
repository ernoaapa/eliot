package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var describeCommand = cli.Command{
	Name:        "describe",
	HelpName:    "describe",
	Usage:       "View details of resource",
	Description: "With this command you can get details of resources",
	ArgsUsage: `can-cli describe RESOURCE [options] [POD NAME]
	
	# Describe a pod
	can-cli describe pod my-pod-name

	# Describe all pods
	can-cli describe pods
`,
	Subcommands: []cli.Command{
		{
			Name:    "pod",
			Aliases: []string{"pods"},
			Usage:   "Return details of pod",
			UsageText: `can-cli describe RESOURCE [options] [POD NAME]
	
	# Describe a pod
	can-cli describe pod my-pod-name

	# Describe all pods
	can-cli describe pods
`,
			Action: func(clicontext *cli.Context) error {
				client := cmd.GetClient(clicontext)
				podName := clicontext.Args().First()

				pods, err := client.GetPods()
				if err != nil {
					return err
				}

				if podName != "" {
					pods = filterByPodName(pods, podName)
				}

				writer := printers.GetNewTabWriter(os.Stdout)
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

func filterByPodName(pods []*pb.Pod, podName string) []*pb.Pod {
	for _, pod := range pods {
		if pod.Metadata.Name == podName {
			return []*pb.Pod{
				pod,
			}
		}
	}
	return []*pb.Pod{}
}
