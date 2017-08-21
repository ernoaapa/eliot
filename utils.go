package main

import (
	"context"

	"github.com/containerd/containerd/namespaces"
	"github.com/urfave/cli"
)

// appContext returns the context for a command. Should only be called once per
// command, near the start.
//
// This will ensure the namespace is picked up and set the timeout, if one is
// defined.
func appContext(clicontext *cli.Context) (context.Context, context.CancelFunc) {
	var (
		ctx       = context.Background()
		timeout   = clicontext.GlobalDuration("timeout")
		namespace = clicontext.GlobalString("namespace")
		cancel    context.CancelFunc
	)

	ctx = namespaces.WithNamespace(ctx, namespace)

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	return ctx, cancel
}
