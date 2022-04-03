package apm

import (
	"context"
)

var std *DefaultLogger

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	std = &DefaultLogger{
		queue:    make(chan Entry, 1000),
		handlers: []Handler{&ConsoleHandler{Formatter: &JsonFormatter{}}},
		cancel:   cancel,
	}
	go std.Run(ctx.Done())
}

func Done() {
	std.Close()
}

func Default() LoggerInterface {
	return std
}

func AddHandlers(handlers ...Handler) {
	std.mu.Lock()
	defer std.mu.Unlock()
	for _, it := range handlers {
		if it == nil {
			continue
		}
		std.handlers = append(std.handlers, it)
	}
	std.handlers.Sort()
}
