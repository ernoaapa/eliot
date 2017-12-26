package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ernoaapa/eliot/pkg/cmd"
	"github.com/ernoaapa/eliot/pkg/cmd/log"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var deviceBuildCommand = cli.Command{
	Name:    "build",
	Aliases: []string{"b"},
	Usage:   "Build device image",
	UsageText: `eli device build [options]
	
	 # Build device image
	 eli device build
	 
	 # Create Linuxkit file but don't build it
	 eli device build --dry-run
	 eli device build --dry-run > linuxkit.yml`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "Linuxkit build source file",
		},
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
		log := log.NewLine().Loading("Get Linuxkit config...")
		var (
			file       = clicontext.String("file")
			dryRun     = clicontext.Bool("dry-run")
			serverURL  = clicontext.String("build-server")
			outputFile = clicontext.String("output")
			outputType = clicontext.String("type")
			linuxkit   []byte
		)

		if file != "" {
			linuxkit, err = ioutil.ReadFile(file)
			if err != nil {
				log.Errorf("Failed to read Linuxkit file: %s", err)
				return err
			}
		} else if cmd.IsPipingIn() {
			linuxkit, err = cmd.ReadAllStdin()
			if err != nil {
				log.Errorf("Failed to read Linuxkit config from stdin: %s", err)
			}
		} else {
			log.Errorf("You must define --file and give path to Linuxkit config file!")
			return errors.New("No Linuxkit config defined")
		}

		if len(linuxkit) == 0 {
			log.Errorf("Invalid Linuxkit config!")
		}

		log.Infof("Resolved Linuxkit config!")

		if dryRun {
			fmt.Println(string(linuxkit))
			return nil
		}

		log.Loadingf("Building RaspberryPI3 Linuxkit image in remote build server...")
		res, err := http.Post(fmt.Sprintf("%s/linuxkit/%s/build/%s", serverURL, "eli-cli", outputType), "application/yml", bytes.NewReader(linuxkit))
		if err != nil {
			return errors.Wrap(err, "Error while making request to Linuxkit build server")
		}
		defer res.Body.Close()

		outFile, err := os.Create(outputFile)
		defer outFile.Close()
		if err != nil {
			log.Errorf("Error, cannot create target output file %s", outputFile)
			return fmt.Errorf("Cannot create target output file %s", outputFile)
		}

		log.Loadingf("Write Linuxkit image to target file...")
		_, err = io.Copy(outFile, res.Body)
		if err != nil {
			log.Errorf("Error while copying image to file [%s]: %s", outFile.Name(), err)
			return errors.New("Unable to copy image to output file")
		}

		log.Donef("Build complete!")
		return nil
	},
}
