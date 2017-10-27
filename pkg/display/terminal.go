package display

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/apoorvam/goterminal"
)

var (
	display = NewTerminal()
)

func Start() {
	display.Start()
}

func Stop() {
	display.Stop()
}

// NewLine creates new updateable output Line
func NewLine() *Line {
	return display.NewLine()
}

// Terminal is tracks the Lines and updates all of them when needed
type Terminal struct {
	rows   []*Line
	change chan struct{}
	writer *goterminal.Writer

	mtx *sync.Mutex
}

// NewTerminal creates new Terminal UI which prints
// output to the
func NewTerminal() *Terminal {
	terminal := &Terminal{
		writer: goterminal.New(os.Stdout),
		mtx:    &sync.Mutex{},
	}
	terminal.Start()
	return terminal
}

func (t *Terminal) Start() {
	t.change = make(chan struct{})
	go func() {
		for {
			select {
			case <-t.change:
				t.Update()
			case <-time.After(100 * time.Millisecond):
				t.Update()
			}
			if t.change == nil {
				return
			}
		}
	}()
}

// Stop updating the terminal lines
func (t *Terminal) Stop() {
	if t.change != nil {
		close(t.change)
		t.change = nil
		t.Update()
		t.rows = []*Line{}
	}
}

// Update will re-render the output
func (t *Terminal) Update() {
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
