package sync

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// StartRsync path to target rsync server on given interval
func StartRsync(done <-chan struct{}, host string, port int, syncs []Sync, interval time.Duration) {
	go func() {
		for {
			select {
			case <-time.After(interval):
				for _, sync := range syncs {
					err := executeRsync(sync.Source, fmt.Sprintf("rsync://%s:%d/%s/", host, port, strings.Replace(sync.Destination, "/", "_", -1)))
					if err != nil {
						log.Errorf("Error while running rsync: %s", err)
					}
				}
			case <-done:
				return
			}
		}
	}()
}

func executeRsync(sourceDir, destination string) error {
	cmd := exec.Command("/usr/bin/rsync", "--recursive", "--times", "--links", "--devices", "--specials", "--compress", sourceDir, destination)
	// TODO: How to display sync completed /failed? Cannot print terminal because messing out terminal attachment. Desktop notification?
	return cmd.Run()
}
