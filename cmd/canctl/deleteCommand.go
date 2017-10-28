package main

import (
	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/cmd/log"
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

				log := log.NewLine().Loading("Fetch pods...")
				pods, err := client.GetPods()
				if err != nil {
					log.Fatalf("Failed to fetch pods information: %s", err)
				}

				if len(pods) == 0 {
					log.Fatal("No pods found")
				}

				if podName != "" {
					pods = cmd.FilterByPodName(pods, podName)

					if len(pods) == 0 {
						log.Fatalf("No pod found with name %s", podName)
					}
				}

				for _, pod := range pods {
					deleted, err := client.DeletePod(pod)
					if err != nil {
						return err
					}
					log.Donef("Deleted pod %s", deleted.Metadata.Name)
				}
				return nil
			},
		},
	},
}
