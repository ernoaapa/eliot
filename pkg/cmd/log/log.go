package log

// Output is interface for log outputs
type Output interface {
	Start()
	Stop()
	NewLine() Line
}

var (
	output Output = NewTerminal()
)

// SetOutput updates logging output
func SetOutput(o Output) {
	output = o
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
func NewLine() Line {
	return output.NewLine()
}
