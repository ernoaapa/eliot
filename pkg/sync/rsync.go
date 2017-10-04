package sync

import (
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

// Rsync path to target rsync server on given interval
func Rsync(done <-chan struct{}, sourceDirs []string, destination string, interval time.Duration) {
	for {
		select {
		case <-time.After(interval):
			for _, sourceDir := range sourceDirs {
				go executeRsync(sourceDir, destination)
			}
		case <-done:
			return
		}
	}
}

func executeRsync(sourceDir, destination string) {
	cmd := exec.Command("/usr/bin/rsync", "-rtp", sourceDir, destination)
	err := cmd.Run()
	if err != nil {
		log.Debugf("Error while running rsync: %s", err)
	}
}
