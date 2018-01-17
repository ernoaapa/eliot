package main

import (
	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/cmd/ui"
	"github.com/urfave/cli"
)

var deletePodCommand = cli.Command{
	Name:    "pod",
	Aliases: []string{"pods"},
	Usage:   "Delete Pod resource(s)",
	UsageText: `eli delete pods [options] [POD NAME]
			 
	 # Delete all Pods
	 eli delete pods

	 # Delete all 'my-pod' pod
	 eli delete pod my-pod`,
	Action: func(clicontext *cli.Context) error {
		config := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(config)

		podName := clicontext.Args().First()

		uiline := ui.NewLine().Loading("Fetch pods...")
		pods, err := client.GetPods()
		if err != nil {
			uiline.Fatalf("Failed to fetch pods information: %s", err)
		}
		uiline.Done("Fetched pods")

		if len(pods) == 0 {
			uiline.Fatal("No pods found")
		}

		if podName != "" {
			pods = cmd.FilterByPodName(pods, podName)

			if len(pods) == 0 {
				uiline.Fatalf("No pod found with name %s", podName)
			}
		}

		for _, pod := range pods {
			uiline = ui.NewLine().Loadingf("Deleting pod %s", pod.Metadata.Name)
			deleted, err := client.DeletePod(pod)
			if err != nil {
				return err
			}
			uiline.Donef("Deleted pod %s", deleted.Metadata.Name)
		}
		return nil
	},
}
