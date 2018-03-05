package node

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
)

func readFile(file string, handler func(string) error) error {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if err := handler(string(line)); err != nil {
			return err
		}
	}

	return nil
}
