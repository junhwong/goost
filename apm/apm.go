package apm

import (
	"context"
)

var std *DefaultLogger

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	std = &DefaultLogger{
		queue:    make(chan Entry, 1024),
		handlers: []Handler{&ConsoleHandler{Formatter: NewJsonFormatter()}},
		cancel:   cancel,
	}
	go std.Run(ctx.Done())
}

func Done() {
	std.Close()
}

// 适配接口
type Adapter interface {
	LoggerInterface
	SpanInterface
}

// 同一接口
type Interface interface {
	Logger
	SpanInterface
}

func Default() Adapter {
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

func New(ctx context.Context) Interface {
	r := &stdImpl{entryLog: entryLog{ctx: ctx, logger: std}}
	r.spi = std
	return r
}

type stdImpl struct {
	entryLog
	spi SpanInterface
}

func (log *stdImpl) NewSpan(ctx context.Context, options ...Option) (context.Context, Span) {
	return log.spi.NewSpan(ctx, options...)
}
