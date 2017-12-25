package cmd

import (
	"bufio"
	"fmt"
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

// IsPipingIn returns true if user is piping some output to `eli` command
func IsPipingIn() bool {
	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	}
	return false
}

// ReadAllStdin read os.Stdin until end and return the output as []byte
func ReadAllStdin() (result []byte, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result = append(result, []byte(fmt.Sprintln(scanner.Text()))...)
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}
