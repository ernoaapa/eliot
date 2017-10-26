package humanreadable

import (
	"bytes"
)

type Bar struct {
	// Fill is the default character representing completed progress
	Fill byte

	// Head is the default character that moves when progress is updated
	Head byte

	// Empty is the default character that represents the empty progress
	Empty byte

	// LeftEnd is the default character in the left most part of the progress indicator
	LeftEnd byte

	// RightEnd is the default character in the right most part of the progress indicator
	RightEnd byte

	// Width is the default width of the progress bar
	Width int
}

// NewBar creates new progress bar renderer
func NewBar() *Bar {
	return &Bar{
		Fill:     '=',
		Head:     '>',
		Empty:    '-',
		LeftEnd:  '[',
		RightEnd: ']',
	}
}

// Render returns the byte presentation of the progress bar
func (b *Bar) Render(width int, current, total int64) []byte {
	if total == 0 {
		return []byte{}
	}
	completedWidth := int(float64(width) * (float64(current) / float64(total)))

	// add fill and empty bits
	var buf bytes.Buffer
	for i := 0; i < completedWidth; i++ {
		buf.WriteByte(b.Fill)
	}
	for i := 0; i < width-completedWidth; i++ {
		buf.WriteByte(b.Empty)
	}

	// set head bit
	pb := buf.Bytes()
	if completedWidth > 0 && completedWidth < width {
		pb[completedWidth-1] = b.Head
	}

	// set left and right ends bits
	pb[0], pb[len(pb)-1] = b.LeftEnd, b.RightEnd

	return pb
}
