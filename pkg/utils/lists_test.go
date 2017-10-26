package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeLists(t *testing.T) {
	result := MergeLists([]string{"a", "b"}, []string{"b", "c"})

	assert.Equal(t, []string{"a", "b", "c"}, result, "should list without duplicates")
}

func TestRotateL(t *testing.T) {
	list := []string{"a", "b", "c"}
	RotateL(&list)
	assert.Equal(t, []string{"b", "c", "a"}, list)
	RotateLBy(&list, 2)
	assert.Equal(t, []string{"a", "b", "c"}, list)
}
func TestRotateR(t *testing.T) {
	list := []string{"a", "b", "c"}
	RotateR(&list)
	assert.Equal(t, []string{"c", "a", "b"}, list)
	RotateRBy(&list, 2)
	assert.Equal(t, []string{"a", "b", "c"}, list)
}
