package main

import (
	"fmt"

	"github.com/ernoaapa/can/cmd"
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
				client := cmd.GetClient(config)
				podName := clicontext.Args().First()

				pods, err := client.GetPods()
				if err != nil {
					return err
				}

				if podName != "" {
					pods = cmd.FilterByPodName(pods, podName)
				}

				if len(pods) == 0 {
					return fmt.Errorf("No pod found with name %s", podName)
				}

				for _, pod := range pods {
					deleted, err := client.DeletePod(pod)
					if err != nil {
						return err
					}
					fmt.Printf("Deleted pod [%s]", deleted.Metadata.Name)
				}
				return nil
			},
		},
	},
}
