package device

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/ernoaapa/eliot/pkg/version"
)

// GetInfo resolves information about the device
// Note: Darwin (OSX) implementation is just for development purpose
// For example, BootID get generated every time when process restarts
func (r *Resolver) GetInfo(grpcPort int) *model.DeviceInfo {
	ioregOutput := runCommandOrFail("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	hostname, _ := os.Hostname()

	return &model.DeviceInfo{
		Version:   version.VERSION,
		Labels:    r.labels,
		Arch:      runtime.GOARCH,
		OS:        runtime.GOOS,
		Hostname:  hostname,
		Addresses: getAddresses(),
		GrpcPort:  grpcPort,

		MachineID: parseFieldFromIoregOutput(ioregOutput, "IOPlatformSerialNumber"),

		SystemUUID: parseFieldFromIoregOutput(ioregOutput, "IOPlatformUUID"),

		BootID: runCommandOrFail("/usr/bin/uuidgen"),
	}
}

func runCommandOrFail(name string, arg ...string) string {
	bytes, err := exec.Command(name, arg...).Output()
	if err != nil {
		log.Fatalf("Failed to resolve device info: %s", err)
	}
	return strings.TrimSpace(string(bytes))
}

func parseFieldFromIoregOutput(output, field string) string {
	exp := regexp.MustCompile(fmt.Sprintf(".*\"%s\".*\"(.*)\"", field))
	return exp.FindStringSubmatch(output)[1]
}
