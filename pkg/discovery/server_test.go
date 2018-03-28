package discovery

import (
	"sync"
	"testing"
)

func TestServerServeStop(t *testing.T) {
	var wg sync.WaitGroup
	server := NewServer("testing", 1234, "v1.0")
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Serve()
	}()

	server.Stop()
	wg.Wait()
}
