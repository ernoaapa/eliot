package yaml

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

type TestStruct struct {
	Foo string
}

var exampleData = []byte(`
foo: bar
---
foo: baz
`)

func TestSplitYamlDocument(t *testing.T) {
	scanner := bufio.NewScanner(bytes.NewReader(exampleData))
	scanner.Split(SplitYAMLDocument)

	result := []*TestStruct{}
	for scanner.Scan() {
		target := &TestStruct{}
		err := yaml.Unmarshal(scanner.Bytes(), target)
		assert.NoError(t, err)

		result = append(result, target)
	}

	assert.Equal(t, "bar", result[0].Foo)
	assert.Equal(t, "baz", result[1].Foo)
}
