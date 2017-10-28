package display

import (
	"fmt"
	"os"

	"github.com/ernoaapa/can/pkg/display/terminal"
	"github.com/fatih/color"
	"github.com/willf/pad"
)

// State of the line
type State int

const (
	// BLANK is the default state which just displays the text
	BLANK State = iota
	// LOADING state displays loading indicator, you must call Done() or Error()
	LOADING
	// DONE represents that the task were done
	DONE
	// WARN is something should be warned
	WARN
	// ERROR is something went wrong
	ERROR
)

// Line is single text line in the terminal output what you can update
type Line struct {
	terminal     *Terminal
	state        State
	Text         string
	showProgress bool
	current      int64
	total        int64
}

var (
	progressBar = terminal.NewBar()
	spinner     = terminal.NewDots()
)

// WithProgress display progress bar when line is in loading state
func (r *Line) WithProgress(current, total int64) *Line {
	r.showProgress = true
	r.current = current
	r.total = total
	r.Update()
	return r
}

func (r *Line) GetProgress() (int64, int64) {
	return r.current, r.total
}

// SetTextf updates the text according to provided format
func (r *Line) SetTextf(format string, args ...interface{}) {
	r.SetText(fmt.Sprintf(format, args...))
}

// SetText updates the text in the line
func (r *Line) SetText(a ...interface{}) {
	r.Text = fmt.Sprint(a...)
	r.Update()
}

// Infof mark this line to be just blank info line
func (r *Line) Infof(format string, args ...interface{}) *Line {
	return r.Info(fmt.Sprintf(format, args...))
}

// Info mark this line to be just blank info line
func (r *Line) Info(a ...interface{}) *Line {
	r.state = BLANK
	r.SetText(a...)
	return r
}

// Loadingf mark this line to be loading (displays loading indicator)
func (r *Line) Loadingf(format string, args ...interface{}) *Line {
	return r.Loading(fmt.Sprintf(format, args...))
}

// Loading mark this line to be loading (displays loading indicator)
func (r *Line) Loading(a ...interface{}) *Line {
	r.state = LOADING
	r.SetText(a...)
	return r
}

// Donef marks this line to be done and updates the text
func (r *Line) Donef(format string, args ...interface{}) *Line {
	return r.Done(fmt.Sprintf(format, args...))
}

// Done marks this line to be done and updates the text
func (r *Line) Done(a ...interface{}) *Line {
	r.state = DONE
	r.SetText(a...)
	return r
}

// Warnf mark this line to be in warning with given format
func (r *Line) Warnf(format string, args ...interface{}) *Line {
	return r.Warn(fmt.Sprintf(format, args...))
}

// Warn mark this line to be in warning with given message
func (r *Line) Warn(a ...interface{}) *Line {
	r.state = WARN
	r.SetText(a...)
	return r
}

// Errorf mark this line to be in error with given format
func (r *Line) Errorf(format string, args ...interface{}) *Line {
	return r.Error(fmt.Sprintf(format, args...))
}

// Error mark this line to be in error with given message
func (r *Line) Error(a ...interface{}) *Line {
	r.state = ERROR
	r.SetText(a...)
	return r
}

// Fatalf mark this line to be in fatal with given format
// Will exit(1) after rerendering the lines
func (r *Line) Fatalf(format string, args ...interface{}) {
	r.Fatal(fmt.Sprintf(format, args...))
}

// Fatal mark this line to be in error with given message
// Will exit(1) after rerendering the lines
func (r *Line) Fatal(a ...interface{}) {
	r.state = ERROR
	r.SetText(a...)
	os.Exit(1)
}

func (r *Line) render() string {
	switch r.state {
	case LOADING:
		if r.showProgress {
			return pad.Left(spinner.Rotate(), 5, " ") + " " + r.Text + " " + string(progressBar.Render(70, r.current, r.total))
		}
		return pad.Left(spinner.Rotate(), 5, " ") + " " + r.Text
	case DONE:
		return color.GreenString(pad.Left("✓", 5, " ")) + " " + r.Text
	case WARN:
		return color.YellowString(pad.Left("⚠", 5, " ")) + " " + r.Text
	case ERROR:
		return color.RedString(pad.Left("✘", 5, " ")) + " " + r.Text
	case BLANK:
		return pad.Left("•", 5, " ") + " " + r.Text
	default:
		return "    " + r.Text
	}
}

// Update triggers re-rendering
func (r *Line) Update() {
	r.terminal.Update()
}
