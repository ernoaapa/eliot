package terminal

import "github.com/ernoaapa/eliot/pkg/utils"

// Spinner is loading spinner which on each rotation moves spinner around
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
