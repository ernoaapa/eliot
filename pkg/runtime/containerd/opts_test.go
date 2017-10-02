package containerd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceOrAppendEnvValues(t *testing.T) {
	result := replaceOrAppendEnvValues([]string{
		"FOO=bar",
		"OTHER=keep",
		"UNSET=this",
	}, []string{
		"FOO=baz",
		"UNSET",
	})

	assert.Equal(t, []string{
		"FOO=baz",
		"OTHER=keep",
	}, result)
}
