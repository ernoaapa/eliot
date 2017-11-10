package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/apoorvam/goterminal"
)

var (
	output = NewTerminal()
)

// Start starts the logging output
func Start() {
	output.Start()
}

// Stop halts updating the output
func Stop() {
	output.Stop()
}

// NewLine creates new updateable output Line
func NewLine() *Line {
	return output.NewLine()
}

// Terminal is tracks the Lines and updates all of them when needed
type Terminal struct {
	running bool
	rows    []*Line
	writer  *goterminal.Writer

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
	t.running = true
	go func() {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				t.Update()
			}
			if !t.running {
				return
			}
		}
	}()
}

// Stop updating the terminal lines
func (t *Terminal) Stop() {
	t.running = false
	t.Update()
	t.rows = []*Line{}
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

// NewLine creates new terminal output line what you can change afterward
func (t *Terminal) NewLine() *Line {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	row := &Line{terminal: t}
	t.rows = append(t.rows, row)
	return row
}
