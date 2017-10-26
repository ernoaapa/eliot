package humanreadable

import "github.com/ernoaapa/can/pkg/utils"

type Spinner struct {
	frames []string
}

// NewDots creates new spinner with dots
func NewDots() *Spinner {
	return &Spinner{[]string{`⠋`, `⠙`, `⠹`, `⠸`, `⠼`, `⠴`, `⠦`, `⠧`, `⠇`, `⠏`}}
}

// Rotate rotates the spinner and returns current value
func (s *Spinner) Rotate() string {
	utils.RotateL(&s.frames)
	return s.frames[0]
}
