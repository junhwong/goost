package zap

import (
	"context"

	"github.com/junhwong/goost/apm"
)

type impl struct {
}

func (p *impl) AddHandlers(handlers ...apm.Handler) {}
func (p *impl) Close() error                        { return nil }
func (p *impl) Flush() error                        { return nil }
func (p *impl) NewSpan(ctx context.Context, options ...apm.SpanOption) (context.Context, apm.Span) {
	return nil, nil
}

// func (p *impl) Log(ctx context.Context, calldepth int, level apm.LogLevel, args []interface{}) {

// }
// func (p *impl) Logf(ctx context.Context, calldepth int, level apm.LogLevel, format string, args []interface{}) {

// }
func (logger *impl) Log(e apm.Entry) {

}
