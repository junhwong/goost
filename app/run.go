package app

import (
	"context"
	"log"
	"sync"
)

var (
	wg sync.WaitGroup
	// cancelFuncs = make(map[int64]context.CancelFunc)
)

// Run wraps the given closure function, starts a new goroutine to run it until it ends or receives a cancellation notification.
//
// param parent: must be a cancelbale Context.
func Run(run func(context.Context), parent ...context.Context) context.Context {
	if run == nil {
		log.Fatalf("launcher: run is required\n")
	}
	p := Context()
	if len(parent) > 0 {
		p = parent[0]
	}

	ctx, cancel := context.WithCancel(p)
	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()
		run(ctx)
	}()
	return ctx
}
