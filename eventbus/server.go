package eventbus

import (
	"context"
	"errors"
	fmt "fmt"
	"sync"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Handler interface {
	Handle(conn Connection, msg *Event)
}

// func (mgr *ConnectionMgr) loop(ctx context.Context, conn Connection) error {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			// 连接断开
// 			return nil
// 		default:
// 		}
// 		err := conn.recv(mgr)
// 		if err != nil {
// 			// server error
// 			// close conn
// 			return err
// 		}
// 	}
// }

// func (mgr *ConnectionMgr) Connect(cs EventBus_ConnectServer) error {
// 	fmt.Println("debug log")
// 	// md
// 	conn := serverConnection{inner: cs}
// 	if err := mgr.onConnected(&conn); err != nil {
// 		return status.Errorf(codes.Unimplemented, "method Connect not implemented")
// 	}
// 	defer mgr.close(&conn)
// 	return mgr.loop(cs.Context(), &conn)
// }

// func (mgr *ConnectionMgr) ClientConnect(cc EventBus_ConnectClient) error {
// 	fmt.Println("debug log")
// 	// md
// 	conn := clientConnection{inner: cc}
// 	if err := mgr.onConnected(&conn); err != nil {
// 		return status.Errorf(codes.Unimplemented, "method Connect not implemented")
// 	}
// 	defer mgr.close(&conn)
// 	return mgr.loop(cc.Context(), &conn)
// }

type Connection interface {
	Send(ctx context.Context, msg *Event) error
	Apply(msg Event, callback func()) error
	IsLocal() bool

	recv(mgr *ConnectionMgr) error
	isServer() bool
	isClient() bool
	queue(task sendTask)
}

type serverConnection struct {
	inner EventBus_ConnectServer
}

func (c *serverConnection) recv(mgr *ConnectionMgr) error {
	msg, err := c.inner.Recv()
	if err == nil {
		return mgr.dispatch(c, msg)
	}
	return err
}
func (*serverConnection) isServer() bool {
	return true
}
func (*serverConnection) isClient() bool {
	return false
}
func (*serverConnection) IsLocal() bool {
	return false
}
func (*serverConnection) Send(ctx context.Context, msg *Event) error {
	return nil
}
func (*serverConnection) Apply(msg Event, callback func()) error {
	return nil
}
func (mgr *serverConnection) queue(task sendTask) {

}

type clientConnection struct {
	inner EventBus_ConnectClient
}

func (c *clientConnection) recv(mgr *ConnectionMgr) error {
	msg, err := c.inner.Recv()
	if err == nil {
		return mgr.dispatch(c, msg)
	}
	return err
}
func (*clientConnection) isServer() bool {
	return false
}
func (*clientConnection) isClient() bool {
	return true
}
func (*clientConnection) IsLocal() bool {
	return true
}
func (*clientConnection) Send(ctx context.Context, msg *Event) error {
	return nil
}
func (*clientConnection) Apply(msg Event, callback func()) error {
	return nil
}
func (mgr *clientConnection) queue(task sendTask) {

}

var ErrNotSupported = errors.New("Not support")

type runer struct {
	mux   *ConnectionMgr
	conns map[string]Connection
	types []string
	mu    sync.Mutex
	local Connection
	cbs   map[string]func(context.Context, *Event)
}

func (r *runer) getConns(filters ...func(Connection) bool) []Connection {
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

func (*runer) recv(mgr *ConnectionMgr) error {
	return nil
}
func (*runer) isServer() bool {
	return false
}
func (*runer) isClient() bool {
	return false
}
func (*runer) IsLocal() bool {
	return true
}
func (*runer) Send(ctx context.Context, msg *Event) error {
	return ErrNotSupported
}
func (*runer) Apply(msg Event, callback func()) error {
	return ErrNotSupported
}

func (*runer) onConnected(conn Connection) error {
	return nil
}

func (*runer) close(conn Connection) error {
	return nil
}

func (r *runer) loop(ctx context.Context, conn Connection) error {
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

func (mgr *runer) Connect(cs EventBus_ConnectServer) error {
	fmt.Println("debug log")
	// md
	conn := serverConnection{inner: cs}
	if err := mgr.onConnected(&conn); err != nil {
		return status.Errorf(codes.Unimplemented, "method Connect not implemented")
	}
	defer mgr.close(&conn)
	return mgr.loop(cs.Context(), &conn)
}

func (mgr *runer) ClientConnect(cc EventBus_ConnectClient) error {
	fmt.Println("debug log")
	// md
	conn := clientConnection{inner: cc}
	if err := mgr.onConnected(&conn); err != nil {
		return status.Errorf(codes.Unimplemented, "method Connect not implemented")
	}
	defer mgr.close(&conn)
	return mgr.loop(cc.Context(), &conn)
}
func (mgr *runer) queue(task sendTask) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if mgr.local != nil {
		mgr.local.queue(task)
		return
	}

	if task.callback != nil {
		mgr.cbs[task.source.GetId()] = task.callback
	}
	mgr.mux.dispatch(mgr, task.source)

}
