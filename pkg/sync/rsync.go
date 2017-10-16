package sync

import (
	"os/exec"
	"time"
)

// StartRsync path to target rsync server on given interval
func StartRsync(done <-chan struct{}, sourceDirs []string, destination string, interval time.Duration) {
	go func() {
		for {
			select {
			case <-time.After(interval):
				for _, sourceDir := range sourceDirs {
					executeRsync(sourceDir, destination)
				}
			case <-done:
				return
			}
		}
	}()
}

func executeRsync(sourceDir, destination string) error {
	cmd := exec.Command("/usr/bin/rsync", "--recursive", "--perms", "--times", "--links", "--devices", "--specials", "--compress", sourceDir, destination)
	// TODO: How to display sync completed /failed? Cannot print terminal because messing out terminal attachment
	return cmd.Run()
}
