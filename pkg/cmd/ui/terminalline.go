package ui

import (
	"fmt"
	"os"

	"github.com/ernoaapa/eliot/pkg/cmd/ui/terminal"
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

// TerminalLine is single text line in the terminal output what you can change afterward
type TerminalLine struct {
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
func (r *TerminalLine) WithProgress(current, total int64) Line {
	r.showProgress = true
	r.current = current
	r.total = total
	r.Update()
	return r
}

// setText updates the text in the line
func (r *TerminalLine) setText(a ...interface{}) {
	r.Text = fmt.Sprint(a...)
	r.Update()
}

// Infof mark this line to be just blank info line
func (r *TerminalLine) Infof(format string, args ...interface{}) Line {
	return r.Info(fmt.Sprintf(format, args...))
}

// Info mark this line to be just blank info line
func (r *TerminalLine) Info(a ...interface{}) Line {
	r.state = BLANK
	r.setText(a...)
	return r
}

// Loadingf mark this line to be loading (displays loading indicator)
func (r *TerminalLine) Loadingf(format string, args ...interface{}) Line {
	return r.Loading(fmt.Sprintf(format, args...))
}

// Loading mark this line to be loading (displays loading indicator)
func (r *TerminalLine) Loading(a ...interface{}) Line {
	r.state = LOADING
	r.setText(a...)
	return r
}

// Donef marks this line to be done and updates the text
func (r *TerminalLine) Donef(format string, args ...interface{}) Line {
	return r.Done(fmt.Sprintf(format, args...))
}

// Done marks this line to be done and updates the text
func (r *TerminalLine) Done(a ...interface{}) Line {
	r.state = DONE
	r.setText(a...)
	return r
}

// Warnf mark this line to be in warning with given format
func (r *TerminalLine) Warnf(format string, args ...interface{}) Line {
	return r.Warn(fmt.Sprintf(format, args...))
}

// Warn mark this line to be in warning with given message
func (r *TerminalLine) Warn(a ...interface{}) Line {
	r.state = WARN
	r.setText(a...)
	return r
}

// Errorf mark this line to be in error with given format
func (r *TerminalLine) Errorf(format string, args ...interface{}) Line {
	return r.Error(fmt.Sprintf(format, args...))
}

// Error mark this line to be in error with given message
func (r *TerminalLine) Error(a ...interface{}) Line {
	r.state = ERROR
	r.setText(a...)
	return r
}

// Fatalf mark this line to be in fatal with given format
// Will exit(1) after rerendering the lines
func (r *TerminalLine) Fatalf(format string, args ...interface{}) {
	r.Fatal(fmt.Sprintf(format, args...))
}

// Fatal mark this line to be in error with given message
// Will exit(1) after rerendering the lines
func (r *TerminalLine) Fatal(a ...interface{}) {
	r.state = ERROR
	r.setText(a...)
	os.Exit(1)
}

func (r *TerminalLine) render() string {
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
func (r *TerminalLine) Update() {
	r.terminal.Update()
}
