package eventbus

import (
	"context"
	fmt "fmt"
	"net"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type ServerEventBusServer struct {
}

func (*ServerEventBusServer) Connect(srv EventBus_ConnectServer) error {
	fmt.Println("一垃圾")

OUTER:
	for {
		select {
		case <-srv.Context().Done():
			break OUTER
		default:
		}
		evt, err := srv.Recv()
		if err != nil {
			return err
		}
		fmt.Println(evt)

	}
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}

func TestServer(t *testing.T) {
	evt := Event{
		Id:   "a",
		Type: "b",
		Time: &timestamp.Timestamp{Seconds: 0},
	}
	proto.Marshal(&evt)
	conn, err := net.Listen("tcp", ":8899")
	if err != nil {
		t.Fatal(err)
	}
	grpc.WithAuthority("")
	svr := grpc.NewServer() // grpc.WithInsecure()
	RegisterEventBusServer(svr, new(ServerEventBusServer))
	t.Fatal(svr.Serve(conn))
}

func TestClient(t *testing.T) {
	// conn, err := grpc.Dial(":8899", grpc.WithInsecure())
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// cli := NewEventBusClient(conn)
	// cc, err := cli.Connect(context.TODO())
	// if err != nil {
	// 	t.Fatal(err)
	// }
	mgr := new(ConnectionMgr)
	// mgr.ClientConnect(cc)
	mgr.Bordcast(context.TODO(), &Event{
		Id:   "a",
		Type: "b",
		Time: &timestamp.Timestamp{Seconds: 0},
	})
}

func TestMux(t *testing.T) {
	mux := ConnectionMgr{
		handlers: make(map[string][]Handler),
	}
	mux.Start(nil)
	mux.handlers["b"] = []Handler{new(eachHandler)}

	mux.Bordcast(context.TODO(), &Event{
		Id:   "a",
		Type: "b",
		Time: &timestamp.Timestamp{Seconds: 0},
	})

	time.Sleep(time.Minute * 2)
}

type eachHandler struct {
}

func (*eachHandler) Accepts() []string {
	return []string{"b"}
}

func (h *eachHandler) Handle(conn Connection, msg *Event) {
	// switch msg.GetType() {
	// case "":
	// default:
	// 	return
	// }
	// conn.Send(context.TODO(), msg)

	fmt.Println(msg)
}

type MessageContext interface {
	context.Context
	ConnectionID() string
}
