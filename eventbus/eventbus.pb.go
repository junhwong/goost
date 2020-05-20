// Code generated by protoc-gen-go. DO NOT EDIT.
// source: eventbus.proto

package eventbus

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// 事件
type Event struct {
	// 当前事件标识。应该在生产者端唯一，通过它可以确定某个事件的处理结果(如果需要)。
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// 事件类型。
	// 如：{apiVerion}/{kind}/{verb} = v1/Pod/POST  v1/loging.PUSH
	Type                 string               `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Time                 *timestamp.Timestamp `protobuf:"bytes,4,opt,name=time,proto3" json:"time,omitempty"`
	Header               map[string]string    `protobuf:"bytes,5,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Data                 []byte               `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf5d638f5c8c3cc4, []int{0}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Event) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Event) GetTime() *timestamp.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *Event) GetHeader() map[string]string {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Event) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Event)(nil), "eventbus.Event")
	proto.RegisterMapType((map[string]string)(nil), "eventbus.Event.HeaderEntry")
}

func init() {
	proto.RegisterFile("eventbus.proto", fileDescriptor_cf5d638f5c8c3cc4)
}

var fileDescriptor_cf5d638f5c8c3cc4 = []byte{
	// 249 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x8f, 0x4f, 0x4f, 0x83, 0x40,
	0x10, 0xc5, 0xb3, 0x14, 0xb0, 0x0e, 0xa6, 0x9a, 0x89, 0x87, 0x0d, 0x1e, 0x24, 0x9e, 0xf6, 0xb4,
	0x18, 0x7a, 0xf1, 0xcf, 0xcd, 0xa4, 0x89, 0x67, 0xe2, 0x17, 0x58, 0xca, 0x58, 0x89, 0x2d, 0x4b,
	0x60, 0x69, 0xc2, 0xa7, 0xf5, 0xab, 0x18, 0x86, 0x62, 0x4c, 0x6f, 0x6f, 0xde, 0xfe, 0xf2, 0xde,
	0x5b, 0x58, 0xd1, 0x91, 0x6a, 0x57, 0xf4, 0x9d, 0x6e, 0x5a, 0xeb, 0x2c, 0x2e, 0xe7, 0x3b, 0xbe,
	0xdf, 0x59, 0xbb, 0xdb, 0x53, 0xca, 0x7e, 0xd1, 0x7f, 0xa6, 0xae, 0x3a, 0x50, 0xe7, 0xcc, 0xa1,
	0x99, 0xd0, 0x87, 0x1f, 0x01, 0xc1, 0x66, 0xa4, 0x71, 0x05, 0x5e, 0x55, 0x4a, 0x91, 0x08, 0x75,
	0x99, 0x7b, 0x55, 0x89, 0x08, 0xbe, 0x1b, 0x1a, 0x92, 0x1e, 0x3b, 0xac, 0x51, 0x83, 0x3f, 0x06,
	0x48, 0x3f, 0x11, 0x2a, 0xca, 0x62, 0x3d, 0xa5, 0xeb, 0x39, 0x5d, 0x7f, 0xcc, 0xe9, 0x39, 0x73,
	0xb8, 0x86, 0xf0, 0x8b, 0x4c, 0x49, 0xad, 0x0c, 0x92, 0x85, 0x8a, 0xb2, 0x3b, 0xfd, 0xb7, 0x94,
	0x4b, 0xf5, 0x3b, 0xbf, 0x6e, 0x6a, 0xd7, 0x0e, 0xf9, 0x09, 0x1d, 0x8b, 0x4b, 0xe3, 0x8c, 0x0c,
	0x13, 0xa1, 0xae, 0x72, 0xd6, 0xf1, 0x33, 0x44, 0xff, 0x50, 0xbc, 0x81, 0xc5, 0x37, 0x0d, 0xa7,
	0xb1, 0xa3, 0xc4, 0x5b, 0x08, 0x8e, 0x66, 0xdf, 0xcf, 0x73, 0xa7, 0xe3, 0xc5, 0x7b, 0x12, 0xd9,
	0x2b, 0x2c, 0xb9, 0xeb, 0xad, 0xef, 0x30, 0x85, 0x8b, 0xad, 0xad, 0x6b, 0xda, 0x3a, 0xbc, 0x3e,
	0x9b, 0x12, 0x9f, 0x1b, 0x4a, 0x3c, 0x8a, 0x22, 0xe4, 0xaf, 0xad, 0x7f, 0x03, 0x00, 0x00, 0xff,
	0xff, 0xd1, 0xea, 0xd7, 0xfd, 0x62, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// EventBusClient is the client API for EventBus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EventBusClient interface {
	// 连接并开启双向事件流
	Connect(ctx context.Context, opts ...grpc.CallOption) (EventBus_ConnectClient, error)
}

type eventBusClient struct {
	cc grpc.ClientConnInterface
}

func NewEventBusClient(cc grpc.ClientConnInterface) EventBusClient {
	return &eventBusClient{cc}
}

func (c *eventBusClient) Connect(ctx context.Context, opts ...grpc.CallOption) (EventBus_ConnectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_EventBus_serviceDesc.Streams[0], "/eventbus.EventBus/connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventBusConnectClient{stream}
	return x, nil
}

type EventBus_ConnectClient interface {
	Send(*Event) error
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventBusConnectClient struct {
	grpc.ClientStream
}

func (x *eventBusConnectClient) Send(m *Event) error {
	return x.ClientStream.SendMsg(m)
}

func (x *eventBusConnectClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventBusServer is the server API for EventBus service.
type EventBusServer interface {
	// 连接并开启双向事件流
	Connect(EventBus_ConnectServer) error
}

// UnimplementedEventBusServer can be embedded to have forward compatible implementations.
type UnimplementedEventBusServer struct {
}

func (*UnimplementedEventBusServer) Connect(srv EventBus_ConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}

func RegisterEventBusServer(s *grpc.Server, srv EventBusServer) {
	s.RegisterService(&_EventBus_serviceDesc, srv)
}

func _EventBus_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventBusServer).Connect(&eventBusConnectServer{stream})
}

type EventBus_ConnectServer interface {
	Send(*Event) error
	Recv() (*Event, error)
	grpc.ServerStream
}

type eventBusConnectServer struct {
	grpc.ServerStream
}

func (x *eventBusConnectServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func (x *eventBusConnectServer) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _EventBus_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eventbus.EventBus",
	HandlerType: (*EventBusServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "connect",
			Handler:       _EventBus_Connect_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "eventbus.proto",
}
