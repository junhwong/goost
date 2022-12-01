package apm

import (
	"context"

	"github.com/junhwong/goost/pkg/field"
)

var (
	std  *DefaultLogger
	defi Interface
)

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	f := NewTextFormatter() // NewJsonFormatter() //
	std = &DefaultLogger{
		queue:    make(chan Entry, 1024),
		inqueue:  make(chan Entry, 1024),
		handlers: []Handler{&ConsoleHandler{Formatter: f}},
		cancel:   cancel,
	}
	go std.Run(ctx.Done())
	defi = New(context.Background())
}

func Done() {
	std.Close()
}
func Flush() {
	std.Flush()
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

	old := std.gethandlers()
	for _, it := range handlers {
		if it == nil {
			continue
		}
		old = append(old, it)
	}
	old.Sort()
	std.handlers = old
}

// type appendFields struct {
// 	fields []field.Field
// }

type appendFields func() []field.Field

func (fn appendFields) applyInterface(impl *stdImpl) {
	fs := fn()
	impl.fields = append(impl.fields, fs...)
}

type Option interface {
	applyInterface(*stdImpl)
}

func WithFields(fs ...field.Field) appendFields {
	return func() []field.Field {
		return fs
	}
}

func New(ctx context.Context, options ...Option) Interface {
	r := &stdImpl{entryLog: entryLog{ctx: ctx, logger: std, calldepth: 1}}
	r.spi = std
	return r
}

type stdImpl struct {
	entryLog
	fields []field.Field
	spi    SpanInterface
}

func (log *stdImpl) NewSpan(ctx context.Context, calldepth int, options ...SpanOption) (context.Context, Span) {
	return log.spi.NewSpan(ctx, calldepth+1, options...)
}
