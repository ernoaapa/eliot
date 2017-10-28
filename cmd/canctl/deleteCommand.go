package main

import (
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
				client := cmd.GetClient(config)

				podName := clicontext.Args().First()

				display := display.New().Loading("Fetch pods...")
				pods, err := client.GetPods()
				if err != nil {
					display.Fatalf("Failed to fetch pods information: %s", err)
				}

				if len(pods) == 0 {
					display.Fatal("No pods found")
				}

				if podName != "" {
					pods = cmd.FilterByPodName(pods, podName)

					if len(pods) == 0 {
						display.Fatalf("No pod found with name %s", podName)
					}
				}

				for _, pod := range pods {
					deleted, err := client.DeletePod(pod)
					if err != nil {
						return err
					}
					display.Donef("Deleted pod %s", deleted.Metadata.Name)
				}
				return nil
			},
		},
	},
}
