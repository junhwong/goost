package apm

import (
	"context"
)

var (
	std  *DefaultLogger
	defi Interface
)

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	f := NewTextFormatter() // NewJsonFormatter()
	std = &DefaultLogger{
		queue:    make(chan Entry, 1024),
		handlers: []Handler{&ConsoleHandler{Formatter: f}},
		cancel:   cancel,
	}
	go std.Run(ctx.Done())
	defi = New(context.Background())
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

func GetAdapter() Adapter {
	return std
}

func Default() Interface {
	return defi
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
	r.calldepth = 1
	return r
}

type stdImpl struct {
	entryLog
	spi SpanInterface
}

func (log *stdImpl) NewSpan(ctx context.Context, options ...Option) (context.Context, Span) {
	return log.spi.NewSpan(ctx, options...)
}