package client

import (
	"encoding/json"
	"fmt"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/go-resty/resty"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Client communicates with Can server
type Client struct {
	url   string
	token string
}

// NewClient creates new Can server client
func NewClient(url, token string) *Client {
	return &Client{
		url,
		token,
	}
}

// GetDeployments lists deployments
func (c *Client) GetDeployments() (deployments []model.Deployment, err error) {
	url := fmt.Sprintf("%s/api/deployments", c.url)

	resp, err := resty.R().
		SetAuthToken(c.token).
		Get(url)
	if err != nil {
		return deployments, errors.Wrapf(err, "Error while making request to %s", url)
	}

	data := resp.Body()

	if resp.StatusCode() >= 400 {
		log.Debugf("Received error response %s %d: %s", url, resp.StatusCode, string(data[:]))
		return deployments, fmt.Errorf("Url replied with status code [%d]", resp.StatusCode())
	}

	return unmarshalDeploymentsJSON(data)
}

// CreateDeployment creates new deployment
func (c *Client) CreateDeployment(deployment *model.Deployment) (*model.Deployment, error) {
	url := fmt.Sprintf("%s/api/deployments", c.url)

	reqData, err := marshalDeploymentJSON(deployment)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create deployment")
	}

	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(c.token).
		SetBody(reqData).Post(url)

	if err != nil {
		return nil, errors.Wrapf(err, "Error while making request to %s", url)
	}

	resData := resp.Body()

	if resp.StatusCode() >= 400 {
		log.Debugf("Received error response %s %d: %s", url, resp.StatusCode, string(resData[:]))
		return nil, fmt.Errorf("Url replied with status code [%d]", resp.StatusCode())
	}

	return unmarshalDeploymentJSON(resData)
}

func unmarshalDeploymentsJSON(data []byte) ([]model.Deployment, error) {
	target := &[]model.Deployment{}
	unmarshalErr := json.Unmarshal(data, target)
	if unmarshalErr != nil {
		log.Debugf("Unable to parse JSON: %s", string(data[:]))
		return nil, errors.Wrapf(unmarshalErr, "Unable to parse JSON data")
	}
	return *target, nil
}

func marshalDeploymentJSON(deployment *model.Deployment) ([]byte, error) {
	data, unmarshalErr := json.Marshal(deployment)
	if unmarshalErr != nil {
		return nil, errors.Wrapf(unmarshalErr, "Unable to marshal to JSON")
	}
	return data, nil
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
