package main

import (
	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/cmd/log"
	"github.com/urfave/cli"
)

var deletePodCommand = cli.Command{
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

		logl := log.NewLine().Loading("Fetch pods...")
		pods, err := client.GetPods()
		if err != nil {
			logl.Fatalf("Failed to fetch pods information: %s", err)
		}
		logl.Done("Fetched pods")

		if len(pods) == 0 {
			logl.Fatal("No pods found")
		}

		if podName != "" {
			pods = cmd.FilterByPodName(pods, podName)

			if len(pods) == 0 {
				logl.Fatalf("No pod found with name %s", podName)
			}
		}

		for _, pod := range pods {
			logl = log.NewLine().Loadingf("Deleting pod %s", pod.Metadata.Name)
			deleted, err := client.DeletePod(pod)
			if err != nil {
				return err
			}
			logl.Donef("Deleted pod %s", deleted.Metadata.Name)
		}
		return nil
	},
}
