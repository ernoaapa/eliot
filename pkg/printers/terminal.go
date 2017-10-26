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

// Terminal is UI implementation which prints output
// to user terminal session
type Terminal struct {
	rows   []*Row
	change chan struct{}
	writer *goterminal.Writer

	mtx *sync.Mutex
}

type State int

const (
	INFO State = iota
	PROGRESS
	DONE
	ERROR
)

type Row struct {
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

func (t *Terminal) NewRow() *Row {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	row := &Row{change: t.change}
	t.rows = append(t.rows, row)
	return row
}

func (r *Row) SetTextf(format string, args ...interface{}) {
	r.SetText(fmt.Sprintf(format, args...))
}

func (r *Row) SetText(str string) {
	r.Text = str
	r.Update()
}

func (r *Row) SetProgress(current, total int64) {
	r.state = PROGRESS
	r.current = current
	r.total = total
	r.Update()
}

func (r *Row) Error() {
	r.state = ERROR
	r.Update()
}

func (r *Row) Done() {
	r.state = DONE
	r.Update()
}

func (r *Row) render() string {
	switch r.state {
	case PROGRESS:
		return pad.Left(spinner.Rotate(), 5, " ") + " " + r.Text + string(progressBar.Render(70, r.current, r.total))
	case DONE:
		return color.GreenString(pad.Left("✓", 5, " ") + " " + r.Text)
	case ERROR:
		return color.RedString(pad.Left("✘", 5, " ") + " " + r.Text)
	default:
		return "    " + r.Text
	}
}

// Update triggers re-rendering
func (r *Row) Update() {
	r.change <- struct{}{}
}
