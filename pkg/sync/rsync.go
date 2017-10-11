package sync

import (
	"os/exec"
	"time"
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
	cmd := exec.Command("/usr/bin/rsync", "--recursive", "--perms", "--times", "--links", "--devices", "--specials", "--compress", sourceDir, destination)
	// TODO: How to display sync completed /failed? Cannot print terminal because messing out terminal attachment
	cmd.Run()
}
