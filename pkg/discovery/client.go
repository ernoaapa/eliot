package discovery

import (
	"context"
	"time"

	node "github.com/ernoaapa/eliot/pkg/api/services/node/v1"
	"github.com/grandcat/zeroconf"
	"github.com/pkg/errors"
)

// Nodes return list of NodeInfos synchronously with given timeout
func Nodes(timeout time.Duration) (nodes []*node.Info, err error) {
	results := make(chan *node.Info)
	defer close(results)

	go func() {
		for node := range results {
			nodes = append(nodes, node)
		}
	}()

	err = NodesAsync(results, timeout)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// NodesAsync search for nodes in network asynchronously for given timeout
func NodesAsync(results chan<- *node.Info, timeout time.Duration) error {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return errors.Wrapf(err, "Failed to initialize new zeroconf resolver")
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(entries <-chan *zeroconf.ServiceEntry) {
		for entry := range entries {
			results <- MapToAPIModel(entry)
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = resolver.Browse(ctx, ZeroConfServiceName, "", entries)
	if err != nil {
		return errors.Wrapf(err, "Failed to browse zeroconf nodes")
	}

	<-ctx.Done()
	return nil
}
