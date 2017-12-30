package log

import "github.com/ernoaapa/eliot/pkg/cmd"

// Output is interface for log outputs
type Output interface {
	Start()
	Stop()
	NewLine() *Line
}

var (
	output = getOutput()
)

func getOutput() Output {
	if cmd.IsPipingOut() {
		return NewHidden()
	}
	return NewTerminal()
}

// Start starts the logging output
func Start() {
	output.Start()
}

// Stop halts updating the output
func Stop() {
	output.Stop()
}

// NewLine creates new updateable output Line
func NewLine() *Line {
	return output.NewLine()
}
