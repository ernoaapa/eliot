package main

import (
	"os"

	"github.com/ernoaapa/eliot/cmd"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/cmd/log"
	"github.com/ernoaapa/eliot/pkg/printers"
	"github.com/ernoaapa/eliot/pkg/progress"
	"github.com/ernoaapa/eliot/pkg/resolve"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var createCommand = cli.Command{
	Name:        "create",
	HelpName:    "create",
	Usage:       "Create pod based on yaml spec",
	Description: "With create command, you can create new pod into the device based on yaml specification",
	UsageText: `eli create [options] -f ./pod.yml

	 # Create pod based on pod.yml
	 eli create -f ./pod.yml
`,
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name: "file, f",

			Usage: "Filename, directory, or URL to files to use to create the resource",
		},
	},
	Subcommands: []cli.Command{
		createPodCommand,
	},
	Action: func(clicontext *cli.Context) (err error) {
		pods := []*pods.Pod{}
		if len(clicontext.StringSlice("file")) > 0 {
			pods, err = resolve.Pods(clicontext.StringSlice("file"))
			if err != nil {
				return err
			}
		} else {
			return errors.New("You need to give --file flag")
		}

		config := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(config)

		for _, pod := range pods {
			logs := map[string]*log.Line{}
			progressc := make(chan []*progress.ImageFetch)

			go func() {
				for fetches := range progressc {
					for _, fetch := range fetches {
						if _, ok := logs[fetch.Image]; !ok {
							logs[fetch.Image] = log.NewLine().Loadingf("Download %s", fetch.Image)
						}

						if fetch.IsDone() {
							if fetch.Failed {
								logs[fetch.Image].Errorf("Failed %s", fetch.Image)
							} else {
								logs[fetch.Image].Donef("Downloaded %s", fetch.Image)
							}
						} else {
							current, total := fetch.GetProgress()
							logs[fetch.Image].WithProgress(current, total)
						}
					}
				}
			}()
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

			if err := printer.PrintPodDetails(result, writer); err != nil {
				return err
			}
		}
		return nil
	},
}
