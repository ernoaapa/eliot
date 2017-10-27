package display

import (
	"fmt"

	"github.com/ernoaapa/can/pkg/display/terminal"
	"github.com/fatih/color"
	"github.com/willf/pad"
)

// State of the line
type State int

const (
	// INFO is the default state which just displays the text
	INFO State = iota
	// ACTIVE state displays loading indicator, you must call Done() or Error()
	ACTIVE
	// PROGRESS state displays progress some progress
	PROGRESS
	// DONE represents that the task were done
	DONE
	// ERROR is something went wrong
	ERROR
)

// Line is single text line in the terminal output what you can update
type Line struct {
	change  chan struct{}
	state   State
	Text    string
	current int64
	total   int64
}

var (
	progressBar = terminal.NewBar()
	spinner     = terminal.NewDots()
)

// SetTextf updates the text according to provided format
func (r *Line) SetTextf(format string, args ...interface{}) {
	r.SetText(fmt.Sprintf(format, args...))
}

// SetText updates the text in the line
func (r *Line) SetText(a ...interface{}) {
	r.Text = fmt.Sprint(a...)
	r.Update()
}

// Activef mark this line to be active (displays loading indicator)
func (r *Line) Activef(format string, args ...interface{}) {
	r.Active(fmt.Sprintf(format, args...))
}

// Active mark this line to be active (displays loading indicator)
func (r *Line) Active(a ...interface{}) {
	r.state = ACTIVE
	r.SetText(a...)
}

// SetProgress mark this line to be in progress with given progress
func (r *Line) SetProgress(current, total int64) {
	r.state = PROGRESS
	r.current = current
	r.total = total
	r.Update()
}

// Errorf mark this line to be in error with given format
func (r *Line) Errorf(format string, args ...interface{}) {
	r.Error(fmt.Sprintf(format, args...))
}

// Error mark this line to be in error with given message
func (r *Line) Error(a ...interface{}) {
	r.state = ERROR
	r.SetText(a...)
}

// Donef marks this line to be done and updates the text
func (r *Line) Donef(format string, args ...interface{}) {
	r.Done(fmt.Sprintf(format, args...))
}

// Done marks this line to be done and updates the text
func (r *Line) Done(a ...interface{}) {
	r.state = DONE
	r.SetText(a...)
}

func (r *Line) render() string {
	switch r.state {
	case ACTIVE:
		return pad.Left(spinner.Rotate(), 5, " ") + " " + r.Text
	case PROGRESS:
		return pad.Left(spinner.Rotate(), 5, " ") + " " + r.Text + " " + string(progressBar.Render(70, r.current, r.total))
	case DONE:
		return color.GreenString(pad.Left("✓", 5, " ")) + " " + r.Text
	case ERROR:
		return color.RedString(pad.Left("✘", 5, " ")) + " " + r.Text
	default:
		return "    " + r.Text
	}
}

// Update triggers re-rendering
func (r *Line) Update() {
	r.change <- struct{}{}
}
