package main

import (
	"os"

	"github.com/ernoaapa/eliot/cmd"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/printers"
	"github.com/ernoaapa/eliot/pkg/progress"
	"github.com/ernoaapa/eliot/pkg/resolve"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var createPodCommand = cli.Command{
	Name:        "pod",
	HelpName:    "pod",
	Usage:       "Create new pod",
	Description: "With create pod command, you can create new pod into the node",
	UsageText: `eli create pod [options] <NAME>

	 # Create new pod 'my-pod' and create single container
	 eli create pod --image alpine my-pod
`,
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "image",
			Usage: "The container image to run. You can pass as many images you want",
		},
	},
	Action: func(clicontext *cli.Context) error {
		var (
			images = clicontext.StringSlice("image")
			name   = clicontext.Args().First()
		)

		if name == "" {
			return errors.New("You need to give name for the pod")
		}

		var pod *pods.Pod
		if len(images) > 0 {
			pod = resolve.BuildPod(name, images)
		} else {
			return errors.New("You need to give at least one --image flag")
		}

		config := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(config)

		progressc := make(chan []*progress.ImageFetch)
		go cmd.ShowDownloadProgress(progressc)

		err := client.CreatePod(progressc, pod)
		close(progressc)
		if err != nil {
			return err
		}

		result, err := client.StartPod(pod.Metadata.Name)
		if err != nil {
			return err
		}

		writer := printers.GetNewTabWriter(os.Stdout)
		defer writer.Flush()
		printer := cmd.GetPrinter(clicontext)

		return printer.PrintPod(result, writer)
	},
}
