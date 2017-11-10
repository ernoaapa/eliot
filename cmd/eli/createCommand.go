package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/elliot/cmd"
	"github.com/ernoaapa/elliot/pkg/cmd/log"
	"github.com/ernoaapa/elliot/pkg/printers"
	"github.com/ernoaapa/elliot/pkg/progress"
	"github.com/ernoaapa/elliot/pkg/resolve"
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
			Name:  "file,f",
			Usage: "Filename, directory, or URL to files to use to create the resource",
		},
	},
	Action: func(clicontext *cli.Context) (err error) {
		sources := clicontext.StringSlice("file")
		if len(sources) == 0 {
			return fmt.Errorf("You must give --file parameter")
		}

		config := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(config)

		pods, err := resolve.Pods(sources)
		if err != nil {
			return errors.Wrapf(err, "Failed to read pod specs from path(s) %s", sources)
		}

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
