package discovery

import (
	"context"
	"sync"
	"time"

	node "github.com/ernoaapa/eliot/pkg/api/services/node/v1"
	"github.com/grandcat/zeroconf"
	"github.com/pkg/errors"
)

// Nodes return list of NodeInfos synchronously with given timeout
func Nodes(timeout time.Duration) (nodes []*node.Info, err error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nodes, errors.Wrapf(err, "Failed to initialize new zeroconf resolver")
	}

	// Initialize a waitgroup variable
	var wg sync.WaitGroup
	wg.Add(1)
	entries := make(chan *zeroconf.ServiceEntry)
	go func(entries <-chan *zeroconf.ServiceEntry) {
		for entry := range entries {
			nodes = append(nodes, MapToAPIModel(entry))
		}
		wg.Done()
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = resolver.Browse(ctx, ZeroConfServiceName, "", entries)
	if err != nil {
		return nodes, errors.Wrapf(err, "Failed to browse zeroconf nodes")
	}

	wg.Wait()

	return nodes, nil
}
