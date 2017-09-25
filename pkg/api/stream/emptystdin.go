package stream

import "io"

// EmptyStdin is io.Read implementation what is empty, i.e. first read return EOF
type EmptyStdin struct{}

func (*EmptyStdin) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}
