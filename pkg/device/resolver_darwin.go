package device

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/ernoaapa/elliot/pkg/model"
)

// GetInfo resolves information about the device
// Note: Darwin (OSX) implementation is just for development purpose
// For example, BootID get generated every time when process restarts
func (r *Resolver) GetInfo() *model.DeviceInfo {
	ioregOutput := runCommandOrFail("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	hostname, _ := os.Hostname()

	return &model.DeviceInfo{
		Labels:   r.labels,
		Arch:     runtime.GOARCH,
		OS:       runtime.GOOS,
		Hostname: hostname,

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
