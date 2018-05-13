package server

import (
	fmt "fmt"
	"math"

	"github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RPCMessage struct {
	Namespace            string   `protobuf:"bytes,1,opt,name=Namespace" json:"Namespace,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RPCMessage) Reset()         { *m = RPCMessage{} }
func (m *RPCMessage) String() string { return proto.CompactTextString(m) }
func (*RPCMessage) ProtoMessage()    {}
func (*RPCMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_rpc_848df06cb7d89181, []int{0}
}
func (m *RPCMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RPCMessage.Unmarshal(m, b)
}
func (m *RPCMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RPCMessage.Marshal(b, m, deterministic)
}
func (dst *RPCMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RPCMessage.Merge(dst, src)
}
func (m *RPCMessage) XXX_Size() int {
	return xxx_messageInfo_RPCMessage.Size(m)
}
func (m *RPCMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_RPCMessage.DiscardUnknown(m)
}

var xxx_messageInfo_RPCMessage proto.InternalMessageInfo

func (m *RPCMessage) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *RPCMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type RPCResponse struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=Error" json:"Error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RPCResponse) Reset()         { *m = RPCResponse{} }
func (m *RPCResponse) String() string { return proto.CompactTextString(m) }
func (*RPCResponse) ProtoMessage()    {}
func (*RPCResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_rpc_848df06cb7d89181, []int{1}
}
func (m *RPCResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RPCResponse.Unmarshal(m, b)
}
func (m *RPCResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RPCResponse.Marshal(b, m, deterministic)
}
func (dst *RPCResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RPCResponse.Merge(dst, src)
}
func (m *RPCResponse) XXX_Size() int {
	return xxx_messageInfo_RPCResponse.Size(m)
}
func (m *RPCResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RPCResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RPCResponse proto.InternalMessageInfo

func (m *RPCResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *RPCResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*RPCMessage)(nil), "rpc.RPCMessage")
	proto.RegisterType((*RPCResponse)(nil), "rpc.RPCResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for MessageService service

type MessageServiceClient interface {
	Command(ctx context.Context, in *RPCMessage, opts ...grpc.CallOption) (*RPCResponse, error)
}

type messageServiceClient struct {
	cc *grpc.ClientConn
}

func NewMessageServiceClient(cc *grpc.ClientConn) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) Command(ctx context.Context, in *RPCMessage, opts ...grpc.CallOption) (*RPCResponse, error) {
	out := new(RPCResponse)
	err := grpc.Invoke(ctx, "/rpc.MessageService/Command", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for MessageService service

type MessageServiceServer interface {
	Command(context.Context, *RPCMessage) (*RPCResponse, error)
}

func RegisterMessageServiceServer(s *grpc.Server, srv MessageServiceServer) {
	s.RegisterService(&_MessageService_serviceDesc, srv)
}

func _MessageService_Command_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RPCMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).Command(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.MessageService/Command",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Command(ctx, req.(*RPCMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _MessageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Command",
			Handler:    _MessageService_Command_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor_rpc_848df06cb7d89181) }

var fileDescriptor_rpc_848df06cb7d89181 = []byte{
	// 171 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x2a, 0x48, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x2a, 0x48, 0x56, 0xb2, 0xe3, 0xe2, 0x0a, 0x0a,
	0x70, 0xf6, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15, 0x92, 0xe1, 0xe2, 0xf4, 0x4b, 0xcc, 0x4d,
	0x2d, 0x2e, 0x48, 0x4c, 0x4e, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x42, 0x08, 0x08, 0x09,
	0x71, 0xb1, 0xb8, 0x24, 0x96, 0x24, 0x4a, 0x30, 0x29, 0x30, 0x6a, 0xf0, 0x04, 0x81, 0xd9, 0x4a,
	0xe6, 0x5c, 0xdc, 0x41, 0x01, 0xce, 0x41, 0xa9, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0x08, 0x25, 0x8c,
	0x08, 0x25, 0x42, 0x22, 0x5c, 0xac, 0xae, 0x45, 0x45, 0xf9, 0x45, 0x60, 0x7d, 0x9c, 0x41, 0x10,
	0x8e, 0x91, 0x03, 0x17, 0x1f, 0xd4, 0xd6, 0xe0, 0xd4, 0xa2, 0xb2, 0xcc, 0xe4, 0x54, 0x21, 0x3d,
	0x2e, 0x76, 0xe7, 0xfc, 0xdc, 0xdc, 0xc4, 0xbc, 0x14, 0x21, 0x7e, 0x3d, 0x90, 0x33, 0x11, 0x0e,
	0x93, 0x12, 0x80, 0x09, 0xc0, 0x6c, 0x52, 0x62, 0x48, 0x62, 0x03, 0x7b, 0xc3, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0x62, 0x4f, 0xa7, 0xc1, 0xd3, 0x00, 0x00, 0x00,
}
