package eventbus

import (
	context "context"
	fmt "fmt"
	"sync"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type baseConn struct {
	mux   *ConnectionMgr
	conns map[string]Connection
	types []string
	mu    sync.Mutex
	local Connection
	cbs   map[string]func(context.Context, *Event)
}

func (r *baseConn) getConns(filters ...func(Connection) bool) []Connection {
	conns := []Connection{}
	for _, c := range r.conns {
		ok := true
		for _, filter := range filters {
			if !filter(c) {
				ok = false
				break
			}
		}
		if ok {
			conns = append(conns, c)
		}
	}
	return conns
}

func (*baseConn) recv(mgr *ConnectionMgr) error {
	return nil
}
func (*baseConn) isServer() bool {
	return false
}
func (*baseConn) isClient() bool {
	return false
}
func (*baseConn) IsLocal() bool {
	return true
}
func (*baseConn) Send(ctx context.Context, msg *Event) error {
	return ErrNotSupported
}
func (*baseConn) Apply(msg Event, callback func()) error {
	return ErrNotSupported
}

func (*baseConn) onConnected(conn Connection) error {
	return nil
}

func (*baseConn) close(conn Connection) error {
	return nil
}

func (r *baseConn) loop(ctx context.Context, conn Connection) error {
	for {
		select {
		case <-ctx.Done():
			// 连接断开
			return nil
		default:
		}
		err := conn.recv(r.mux)
		if err != nil {
			// server error
			// close conn
			return err
		}
	}
}

func (mgr *baseConn) Connect(cs EventBus_ConnectServer) error {
	fmt.Println("debug log")
	// md
	conn := serverConnection{inner: cs}
	if err := mgr.onConnected(&conn); err != nil {
		return status.Errorf(codes.Unimplemented, "method Connect not implemented")
	}
	defer mgr.close(&conn)
	return mgr.loop(cs.Context(), &conn)
}

func (mgr *baseConn) ClientConnect(cc EventBus_ConnectClient) error {
	fmt.Println("debug log")
	// md
	conn := clientConnection{inner: cc}
	if err := mgr.onConnected(&conn); err != nil {
		return status.Errorf(codes.Unimplemented, "method Connect not implemented")
	}
	defer mgr.close(&conn)
	return mgr.loop(cc.Context(), &conn)
}

func (mgr *baseConn) send(task sendTask) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.local != nil {
		mgr.local.queue(task)
		return
	}

	if task.callback != nil {
		mgr.cbs[task.source.GetId()] = task.callback
	}
	// mgr.mux.dispatch(mgr, task.source)

	ctx, cancel := context.WithTimeout(context.TODO(), 10)
	task.cancel = cancel
	go func() {
		defer cancel()
		<-ctx.Done()

	}()
}

type Thenable func(interface{}) interface{}

// fulfilled
// rejected

type Promise interface {
	Reject(err error) Promise
	Resolve(result interface{}) Promise
	Then(Thenable) Promise
	Catch(func(error) interface{}) Promise
}
