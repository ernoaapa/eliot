package log

// Debug is log implementation which just outputs line by line
type Debug struct {
	running bool
}

// NewDebug creates new UI which just prints debug output
func NewDebug() *Debug {
	debug := &Debug{}
	debug.Start()
	return debug
}

func (t *Debug) Start() {
	t.running = true
}

// Stop updating the terminal lines
func (t *Debug) Stop() {
	t.running = false
}

// NewLine creates new terminal output line what you can change afterward
func (t *Debug) NewLine() Line {
	return &DebugLine{}
}
