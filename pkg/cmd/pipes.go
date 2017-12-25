package cmd

import (
	"os"
)

// IsPipingOut returns true if user is piping the output to file or some other command
func IsPipingOut() bool {
	stat, _ := os.Stdout.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	}
	return false
}
