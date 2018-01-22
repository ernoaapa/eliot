package printers

import (
	"io"

	device "github.com/ernoaapa/eliot/pkg/api/services/device/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type YamlPrinter struct {
}

func NewYamlPrinter() *YamlPrinter {
	return &YamlPrinter{}
}

// PrintPodsTable implementation
func (p *YamlPrinter) PrintPodsTable(pods []*pods.Pod, w io.Writer) error {
	if err := writeAsYml(pods, w); err != nil {
		return errors.Wrap(err, "Failed to write pods yaml")
	}
	return nil
}

// PrintDevicesTable implementation
func (p *YamlPrinter) PrintDevicesTable(devices []model.DeviceInfo, w io.Writer) error {
	if err := writeAsYml(devices, w); err != nil {
		return errors.Wrap(err, "Failed to write devices yaml")
	}
	return nil
}

// PrintDeviceDetails implementation
func (p *YamlPrinter) PrintDeviceDetails(device *device.Info, w io.Writer) error {
	if err := writeAsYml(device, w); err != nil {
		return errors.Wrap(err, "Failed to write device yaml")
	}
	return nil
}

// PrintPodDetails implementation
func (p *YamlPrinter) PrintPodDetails(pod *pods.Pod, w io.Writer) error {
	if err := writeAsYml(pod, w); err != nil {
		return errors.Wrap(err, "Failed to write pod yaml")
	}
	return nil
}

// PrintConfig implementation
func (p *YamlPrinter) PrintConfig(config *config.Config, w io.Writer) error {
	if err := writeAsYml(config, w); err != nil {
		return errors.Wrap(err, "Failed to write config yaml")
	}
	return nil
}

func writeAsYml(in interface{}, w io.Writer) error {
	data, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}
