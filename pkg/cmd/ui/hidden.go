package ui

// Hidden is log output implementation what doesn't output anything
type Hidden struct {
}

// NewHidden creates new hidden log output
func NewHidden() *Hidden {
	return &Hidden{}
}

// Start is log Output implementation
func (*Hidden) Start() {}

// Stop is log Output implementation
func (*Hidden) Stop() {}

// Update is log Output implementation
func (*Hidden) Update() {}

// NewLine creates new Line what doesn't show up anywhere
func (h *Hidden) NewLine() Line {
	return &HiddenLine{}
}
