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
	ArgsUsage: `can-cli delete RESOURCE [options]

	 # Delete all running pods
	 can-cli delete pods

	 # Delete all 'my-pod' pod
	 can-cli delete pod my-pod`,
	Subcommands: []cli.Command{
		{
			Name:    "pod",
			Aliases: []string{"pods"},
			Usage:   "Delete Pod resource(s)",
			UsageText: `can-cli delete pods [options] [POD NAME]
			 
	 # Delete all Pods
	 can-cli delete pods

	 # Delete all 'my-pod' pod
	 can-cli delete pod my-pod`,
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
