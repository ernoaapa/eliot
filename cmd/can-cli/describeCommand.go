package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var describeCommand = cli.Command{
	Name: "describe",
	Subcommands: []cli.Command{
		{
			Name:    "pod",
			Aliases: []string{"po", "pods"},
			Usage:   "Return details of pod",
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
