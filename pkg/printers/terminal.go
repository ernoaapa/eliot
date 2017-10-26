package printers

import (
	"fmt"
	"os"
	"sync"

	"github.com/ernoaapa/can/pkg/printers/humanreadable"

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
	r.update()
}

func (r *Row) SetProgress(current, total int64) {
	r.state = PROGRESS
	r.current = current
	r.total = total
	r.update()
}

func (r *Row) Done() {
	r.state = DONE
	r.update()
}

func (r *Row) render() string {
	switch r.state {
	case PROGRESS:
		return r.Text + string(progressBar.Render(70, r.current, r.total))
	default:
		return r.Text
	}
}

func (r *Row) update() {
	r.change <- struct{}{}
}
