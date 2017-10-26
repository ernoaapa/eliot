package printers

import (
	"fmt"
	"os"
	"sync"

	"github.com/ernoaapa/can/pkg/printers/humanreadable"
	"github.com/fatih/color"
	"github.com/willf/pad"

	"github.com/apoorvam/goterminal"
)

var (
	terminal = NewTerminal()
)

// NewOutputLine creates new updateable output Line
func NewOutputLine() *Line {
	return terminal.NewLine()
}

// Terminal is tracks the Lines and updates all of them when needed
type Terminal struct {
	rows   []*Line
	change chan struct{}
	writer *goterminal.Writer

	mtx *sync.Mutex
}

// State of the line
type State int

const (
	// INFO is the default state which just displays the text
	INFO State = iota
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
	progressBar = humanreadable.NewBar()
	spinner     = humanreadable.NewDots()
)

// NewTerminal creates new Terminal UI which prints
// output to the
func NewTerminal() *Terminal {
	terminal := &Terminal{
		change: make(chan struct{}),
		writer: goterminal.New(os.Stdout),
		mtx:    &sync.Mutex{},
	}
	terminal.start()
	return terminal
}

func (t *Terminal) start() {
	go func() {
		for range t.change {
			t.update()
		}
	}()
}

func (t *Terminal) update() {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.writer.Clear()
	for _, row := range t.rows {
		fmt.Fprintln(t.writer, row.render())
	}
	t.writer.Print()
}

// NewLine creates new terminal output line what you can update
func (t *Terminal) NewLine() *Line {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	row := &Line{change: t.change}
	t.rows = append(t.rows, row)
	return row
}

// SetTextf updates the text according to provided format
func (r *Line) SetTextf(format string, args ...interface{}) {
	r.SetText(fmt.Sprintf(format, args...))
}

// SetText updates the text in the line
func (r *Line) SetText(str string) {
	r.Text = str
	r.Update()
}

// SetProgress mark this line to be in progress with given progress
func (r *Line) SetProgress(current, total int64) {
	r.state = PROGRESS
	r.current = current
	r.total = total
	r.Update()
}

// Error mark this line to be in error with given message
func (r *Line) Error(text string) {
	r.state = ERROR
	r.SetText(text)
}

// Done marks this line to be done and updates the text
func (r *Line) Done(text string) {
	r.state = DONE
	r.SetText(text)
}

func (r *Line) render() string {
	switch r.state {
	case PROGRESS:
		return pad.Left(spinner.Rotate(), 5, " ") + " " + r.Text + string(progressBar.Render(70, r.current, r.total))
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
