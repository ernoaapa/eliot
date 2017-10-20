package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/printers"
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
	UsageText: `can-cli create [options] -f ./pod.yml

	 # Create pod based on pod.yml
	 can-cli create -f ./pod.yml
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

		writer := printers.GetNewTabWriter(os.Stdout)
		defer writer.Flush()
		printer := cmd.GetPrinter(clicontext)
		config := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(config)

		pods, err := resolve.Pods(sources)
		if err != nil {
			return errors.Wrapf(err, "Failed to read pod specs from path(s) %s", sources)
		}

		for _, pod := range pods {
			if err := client.CreatePod(pod); err != nil {
				return err
			}

			result, err := client.StartPod(pod.Metadata.Name)
			if err != nil {
				return err
			}

			if err := printer.PrintPodDetails(result, writer); err != nil {
				log.Errorf("Error while printing pod details: %s", err)
			}
		}
		return nil
	},
}
