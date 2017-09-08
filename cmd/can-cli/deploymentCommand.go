package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var deploymentCommand = cli.Command{
	Name: "deployments",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create new deployment",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "Source file path",
				},
			},
			Action: func(clicontext *cli.Context) error {
				client := getClient(clicontext)

				deployment, err := readFromFile(clicontext.String("file"))
				if err != nil {
					log.Fatal(err)
				}
				created, err := client.CreateDeployment(deployment)
				if err != nil {
					log.Fatal(err)
				}

				log.Printf("Created: %v", created)
				return nil
			},
		},
	},
}

func readFromFile(path string) (*model.Deployment, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Cannot update state, file [%s] does not exist", path)
		}
		return nil, err
	}

	return unmarshalDeploymentJSON(data)
}

func unmarshalDeploymentJSON(data []byte) (*model.Deployment, error) {
	target := &model.Deployment{}
	unmarshalErr := json.Unmarshal(data, target)
	if unmarshalErr != nil {
		log.Debugf("Unable to parse JSON: %s", string(data[:]))
		return nil, errors.Wrapf(unmarshalErr, "Unable to parse JSON data")
	}
	return target, nil
}
