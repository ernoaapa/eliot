package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/display"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/ernoaapa/can/pkg/progress"
	"github.com/ernoaapa/can/pkg/resolve"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var createCommand = cli.Command{
	Name:        "create",
	HelpName:    "create",
	Usage:       "Create pod based on yaml spec",
	Description: "With create command, you can create new pod into the device based on yaml specification",
	UsageText: `canctl create [options] -f ./pod.yml

	 # Create pod based on pod.yml
	 canctl create -f ./pod.yml
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
			displays := map[string]*display.Line{}
			progressc := make(chan []*progress.ImageFetch)

			go func() {
				for _, fetch := range <-progressc {
					if _, ok := displays[fetch.Image]; !ok {
						displays[fetch.Image] = display.New()
					}

					if fetch.IsDone() {
						if fetch.Failed {
							displays[fetch.Image].Errorf("Failed %s", fetch.Image)
						} else {
							displays[fetch.Image].Donef("Downloaded %s", fetch.Image)
						}
					} else {
						current, total := fetch.GetProgress()
						displays[fetch.Image].WithProgress(current, total).Loadingf("Download %s", fetch.Image)
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
				log.Errorf("Error while printing pod details: %s", err)
			}
		}
		return nil
	},
}
