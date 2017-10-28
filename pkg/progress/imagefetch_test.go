package progress

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProgress(t *testing.T) {
	fetch := CreateImageFetch("containerID", "imageref", true, map[string]*Status{
		"1": {Offset: 20, Total: 100},
		"2": {Offset: 50, Total: 200},
	})
	current, total := fetch.GetProgress()

	assert.Equal(t, int64(70), current)
	assert.Equal(t, int64(300), total)
}
func TestIsDone(t *testing.T) {
	fetch := CreateImageFetch("containerID", "imageref", true, map[string]*Status{
		"1": {Offset: 100, Total: 100},
		"2": {Offset: 200, Total: 200},
	})

	assert.True(t, fetch.IsDone())
}
