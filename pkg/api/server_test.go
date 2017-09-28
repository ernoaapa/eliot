package api

import (
	"testing"

	assert "gopkg.in/go-playground/assert.v1"
)

func TestGetMetadataValue(t *testing.T) {
	md := map[string][]string{
		"foo": []string{
			"bar",
		},
		"crazy": []string{
			"first",
			"second",
		},
	}

	assert.Equal(t, "bar", getMetadataValue(md, "foo"))
	assert.Equal(t, "first", getMetadataValue(md, "crazy"))
	assert.Equal(t, "", getMetadataValue(md, "dontexist"))
}
