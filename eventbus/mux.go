package eventbus

import (
	context "context"
	"net"
	"time"

	grpc "google.golang.org/grpc"
)

type ConnectionMgr struct {
	handlers map[string][]Handler
	runer    *runer
}

func (mgr *ConnectionMgr) dispatch(conn Connection, msg *Event) error {
	handlers := map[Handler]bool{}

	if hls, ok := mgr.handlers[msg.GetType()]; ok {
		for _, handler := range hls {
			handlers[handler] = true
		}
	}

	for handler := range handlers {
		go handler.Handle(conn, msg)
	}

	return nil
}

func (mgr *ConnectionMgr) Bordcast(ctx context.Context, msg *Event, filters ...func(Connection) bool) error {
	var runer *runer

	runer = mgr.runer // lock

	incLocl := true
	for _, filter := range filters {
		if !filter(runer) {
			incLocl = false
			break
		}
	}
	if incLocl {
		if err := mgr.dispatch(runer, msg); err != nil {
			return err
		}
	}

	for _, c := range runer.getConns(filters...) {
		if err := c.Send(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

type sendTask struct {
	callback        func(context.Context, *Event)
	callbackTimeout time.Duration
	callbackOnce    bool
	localOnly       bool
	mustSend        bool
	source          *Event
	filters         []func(ConnectionMeta) bool
	cancel          context.CancelFunc
}

//https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/Promise/resolve
func (*sendTask) Reject(err error) {

}
func (*sendTask) Resolve(result interface{}) {

}

type SendOption interface {
	filter(task *sendTask) bool
}

func WithLocal() SendOption {
	return nil
}

type sendOption struct {
	do func(task *sendTask) bool
}

func (opt *sendOption) filter(task *sendTask) bool {
	return opt.do(task)
}
func WithCallback(cb func(context.Context, *Event), timeout ...time.Duration) SendOption {
	if cb == nil {
		panic("eventbus: cb cannot be nil")
	}
	var t time.Duration = -1
	if n := len(timeout); n > 0 {
		t = timeout[n-1]
	}
	// todo(pool)
	return &sendOption{do: func(task *sendTask) bool {
		if task.callback != nil {
			panic("eventbus: already set callback")
		}
		task.callback = cb
		task.callbackTimeout = t
		return true
	}}
}

func (mux *ConnectionMgr) Send(msg *Event, opts ...SendOption) error {
	task := sendTask{} // todo(pool)
	task.source = msg
	task.callback = nil
	task.callbackOnce = false

	for _, opt := range opts {
		if opt != nil {
			opt.filter(&task)
		}
	}

	conns := mux.runer.getConns()
	for _, c := range conns {
		if c == nil {
			// filter type
			// filter label
		}
		c.Send(nil, msg)
	}
	msgType := msg.GetType()
	redy := map[string]map[string]Connection{} // map[msg-type][conn-id]conn
	sent := false
OUTER:
	for _, conn := range redy[msgType] {
		if conn == nil {
			continue
		}
		if task.localOnly && !conn.IsLocal() {
			continue
		}
		for _, filter := range task.filters {
			if !filter(nil) {
				continue OUTER
			}
		}
		conn.queue(task)
		sent = true
	}
	if task.mustSend && !sent {

	}

	return nil
}

func (mux *ConnectionMgr) Start(config interface{}, opts ...interface{}) error {
	mux.runer = &runer{
		mux:   mux,
		conns: make(map[string]Connection),
		types: make([]string, 0),
	}

	isRunServer := true
	if isRunServer {
		if err := mux.startServer(":8899"); err != nil {
			return err
		}
	}
	return nil
}

func (mux *ConnectionMgr) startServer(addr string) error {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	svr := grpc.NewServer() // grpc.WithInsecure()
	RegisterEventBusServer(svr, mux.runer)
	go func() {
		svr.Serve(conn) // todo(retry)
	}()
	return nil
}

type ConnectionMeta interface {
	ID() string
	Types() []string
	Tags() []string
}

type connectionMeta struct {
	id    string
	types []string
	tags  []string
}

func (m *connectionMeta) ID() string {
	return m.id
}
func (m *connectionMeta) Types() []string {
	return m.types
}
func (m *connectionMeta) Tags() []string {
	return m.tags
}
