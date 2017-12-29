package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ernoaapa/eliot/pkg/cmd"
	"github.com/ernoaapa/eliot/pkg/cmd/build"
	"github.com/ernoaapa/eliot/pkg/cmd/log"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var deviceBuildCommand = cli.Command{
	Name:    "build",
	Aliases: []string{"b"},
	Usage:   "Build device image",
	UsageText: `eli device build [options] [FILE | URL]
	
	 # Build default device image
	 eli device build
	 
	 # Create Linuxkit file but don't build it
	 eli device build --dry-run
	 eli device build --dry-run > custom-linuxkit.yml
	 
	 # Build from custom config and unpack to directory
	 mkdir dist
	 eli device build custom-linuxkit.yml | tar xv -C dist
	 `,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Print the final Linuxkit config and don't actually build it",
		},
		cli.StringFlag{
			Name:   "build-server",
			Usage:  "Linuxkit build server (github.com/ernoaapa/linuxkit-server) base url",
			Value:  "http://build.eliot.run",
			EnvVar: "ELIOT_BUILD_SERVER",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "Target output file",
			Value: "image.tar",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "Target build type, one of Linuxkit output types",
			Value: "rpi3",
		},
	},
	Action: func(clicontext *cli.Context) (err error) {
		var (
			source     = clicontext.Args().First()
			dryRun     = clicontext.Bool("dry-run")
			serverURL  = clicontext.String("build-server")
			outputFile = clicontext.String("output")
			outputType = clicontext.String("type")

			logline  *log.Line
			linuxkit []byte
			output   io.Writer
		)

		logline = log.NewLine().Loading("Get Linuxkit config...")
		linuxkit, err = build.ResolveLinuxkitConfig(source)
		if err != nil {
			logline.Errorf("Failed to resolve Linuxkit config: %s", err)
			return errors.Wrap(err, "Cannot resolve Linuxkit config")
		}
		logline.Done("Resolved Linuxkit config!")

		logline = log.NewLine().Loading("Resolve output...")
		if outputFile != "" {
			outFile, err := os.Create(outputFile)
			if err != nil {
				logline.Errorf("Error, cannot create target output file %s", outputFile)
				return fmt.Errorf("Cannot create target output file %s", outputFile)
			}
			defer outFile.Close()
			output = outFile
		} else if cmd.IsPipingOut() {
			output = os.Stdout
		} else {
			logline.Errorf("You must give target path with --output or pipe output!")
			return errors.New("No output defined")
		}
		logline.Done("Resolved output!")

		if dryRun {
			fmt.Println(string(linuxkit))
			return nil
		}

		logline = log.NewLine().Loadingf("Building RaspberryPI3 Linuxkit image in remote build server...")
		image, err := build.BuildImage(serverURL, outputType, linuxkit)
		if err != nil {
			logline.Errorf("Failed to build Linuxkit image: %s", err)
			return errors.Wrap(err, "Failed to build Linuxkit image")
		}

		logline.Loadingf("Write Linuxkit image to output...")
		_, err = io.Copy(output, image)
		if err != nil {
			logline.Errorf("Error while writing output: %s", err)
			return errors.New("Unable to copy image to output")
		}

		logline.Donef("Build complete!")
		return nil
	},
}
