package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeLists(t *testing.T) {
	result := MergeLists([]string{"a", "b"}, []string{"b", "c"})

	assert.Equal(t, []string{"a", "b", "c"}, result, "should list without duplicates")
}
