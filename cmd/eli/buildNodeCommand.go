package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ernoaapa/eliot/pkg/cmd"
	"github.com/ernoaapa/eliot/pkg/cmd/build"
	"github.com/ernoaapa/eliot/pkg/cmd/ui"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var buildNodeCommand = cli.Command{
	Name:    "node",
	Aliases: []string{"nodes", "device"}, // device is deprecated command
	Usage:   "Build node image",
	UsageText: `eli build node [options] [FILE | URL]
	
	 # Build default node disk img -file
	 eli build node
	 
	 # Create Linuxkit file but don't build it
	 eli build node --dry-run
	 eli build node --dry-run > custom-linuxkit.yml
	 
	 # Build from custom config and unpack to directory
	 mkdir dist
	 eli build node custom-linuxkit.yml --format tar | tar xv -C dist
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
Usage: "Target output file. (default: image.tar or image.img)",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "Target build type, one of Linuxkit output types",
			Value: "rpi3",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "Target output format. One of [tar, img] (default: img)",
			Value: "img",
		},
	},
	Action: func(clicontext *cli.Context) (err error) {
		var (
			source       = clicontext.Args().First()
			dryRun       = clicontext.Bool("dry-run")
			serverURL    = clicontext.String("build-server")
			outputType   = clicontext.String("type")
			outputFormat = clicontext.String("format")
			outputFile   = getBuildOutputFile(outputFormat, clicontext.String("output"))

			uiline   ui.Line
			linuxkit []byte
			output   io.Writer
		)

		uiline = ui.NewLine().Loading("Get Linuxkit config...")
		linuxkit, err = build.ResolveLinuxkitConfig(source)
		if err != nil {
			uiline.Errorf("Failed to resolve Linuxkit config: %s", err)
			return errors.Wrap(err, "Cannot resolve Linuxkit config")
		}
		uiline.Done("Resolved Linuxkit config!")

		uiline = ui.NewLine().Loading("Resolve output...")
		if cmd.IsPipingOut() {
			output = os.Stdout
			uiline.Done("Resolved output to stdout!")
		} else {
			outFile, err := os.Create(outputFile)
			if err != nil {
				uiline.Errorf("Error, cannot create target output file %s", outputFile)
				return fmt.Errorf("Cannot create target output file %s", outputFile)
			}
			defer outFile.Close()
			output = outFile
			uiline.Donef("Resolved output: %s!", outFile.Name())
		}

		if dryRun {
			fmt.Println(string(linuxkit))
			return nil
		}

		uiline = ui.NewLine().Loadingf("Building Linuxkit image in remote build server...")
		image, err := build.BuildImage(serverURL, outputType, outputFormat, linuxkit)
		if err != nil {
			uiline.Errorf("Failed to build Linuxkit image: %s", err)
			return errors.Wrap(err, "Failed to build Linuxkit image")
		}

		uiline.Loadingf("Write Linuxkit image to output...")
		_, err = io.Copy(output, image)
		if err != nil {
			uiline.Errorf("Error while writing output: %s", err)
			return errors.New("Unable to copy image to output")
		}

		uiline.Donef("Build complete!")
		return nil
	},
}

func getBuildOutputFile(format, outputFile string) string {
	if outputFile != "" {
		return outputFile
	}

	switch format {
	case "img":
		return "image.img"
	case "tar":
		return "image.tar"
	default:
		log.Fatalf("Unknown output format %s, cannot resolve default output file", format)
		return ""
	}
}
