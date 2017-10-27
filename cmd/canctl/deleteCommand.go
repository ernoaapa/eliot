package main

import (
	"fmt"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/display"
	"github.com/urfave/cli"
)

var deleteCommand = cli.Command{
	Name:        "delete",
	HelpName:    "delete",
	Usage:       `Delete one or more resources`,
	Description: "With this command you can delete resources",
	ArgsUsage: `canctl delete RESOURCE [options]

	 # Delete all running pods
	 canctl delete pods

	 # Delete all 'my-pod' pod
	 canctl delete pod my-pod`,
	Subcommands: []cli.Command{
		{
			Name:    "pod",
			Aliases: []string{"pods"},
			Usage:   "Delete Pod resource(s)",
			UsageText: `canctl delete pods [options] [POD NAME]
			 
	 # Delete all Pods
	 canctl delete pods

	 # Delete all 'my-pod' pod
	 canctl delete pod my-pod`,
			Action: func(clicontext *cli.Context) error {
				config := cmd.GetConfigProvider(clicontext)
				client, err := cmd.GetClient(config)
				if err != nil {
					return err
				}

				podName := clicontext.Args().First()

				display := display.NewLine()
				display.Active("Fetch pods...")
				pods, err := client.GetPods()
				if err != nil {
					display.Error(err)
					return err
				}

				if len(pods) == 0 {
					display.Error("No pods found")
					return fmt.Errorf("No pods found")
				}

				if podName != "" {
					pods = cmd.FilterByPodName(pods, podName)

					if len(pods) == 0 {
						display.Errorf("No pod found with name %s", podName)
						return fmt.Errorf("No pod found with name %s", podName)
					}
				}

				for _, pod := range pods {
					deleted, err := client.DeletePod(pod)
					if err != nil {
						return err
					}
					display.Done("Deleted pod %s", deleted.Metadata.Name)
				}
				return nil
			},
		},
	},
}
