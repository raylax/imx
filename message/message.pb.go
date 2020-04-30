// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package message

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type MessageType int32

const (
	MessageType_NONE  MessageType = 0
	MessageType_P2P   MessageType = 1
	MessageType_GROUP MessageType = 2
	MessageType_CTRL  MessageType = 3
)

var MessageType_name = map[int32]string{
	0: "NONE",
	1: "P2P",
	2: "GROUP",
	3: "CTRL",
}

var MessageType_value = map[string]int32{
	"NONE":  0,
	"P2P":   1,
	"GROUP": 2,
	"CTRL":  3,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}

func (MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

type MessageResponse_Status int32

const (
	MessageResponse_Ok        MessageResponse_Status = 0
	MessageResponse_NotRouted MessageResponse_Status = 1
)

var MessageResponse_Status_name = map[int32]string{
	0: "Ok",
	1: "NotRouted",
}

var MessageResponse_Status_value = map[string]int32{
	"Ok":        0,
	"NotRouted": 1,
}

func (x MessageResponse_Status) String() string {
	return proto.EnumName(MessageResponse_Status_name, int32(x))
}

func (MessageResponse_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1, 0}
}

type WsResponse_Status int32

const (
	WsResponse_NONE  WsResponse_Status = 0
	WsResponse_Ok    WsResponse_Status = 1
	WsResponse_Error WsResponse_Status = 2
)

var WsResponse_Status_name = map[int32]string{
	0: "NONE",
	1: "Ok",
	2: "Error",
}

var WsResponse_Status_value = map[string]int32{
	"NONE":  0,
	"Ok":    1,
	"Error": 2,
}

func (x WsResponse_Status) String() string {
	return proto.EnumName(WsResponse_Status_name, int32(x))
}

func (WsResponse_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2, 0}
}

type MessageRequest struct {
	Id                   []string          `protobuf:"bytes,1,rep,name=id,proto3" json:"id,omitempty"`
	Message              *WsMessageRequest `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Source               string            `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MessageRequest) Reset()         { *m = MessageRequest{} }
func (m *MessageRequest) String() string { return proto.CompactTextString(m) }
func (*MessageRequest) ProtoMessage()    {}
func (*MessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *MessageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageRequest.Unmarshal(m, b)
}
func (m *MessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageRequest.Marshal(b, m, deterministic)
}
func (m *MessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageRequest.Merge(m, src)
}
func (m *MessageRequest) XXX_Size() int {
	return xxx_messageInfo_MessageRequest.Size(m)
}
func (m *MessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MessageRequest proto.InternalMessageInfo

func (m *MessageRequest) GetId() []string {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *MessageRequest) GetMessage() *WsMessageRequest {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *MessageRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

type MessageResponse struct {
	Status               MessageResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=message.MessageResponse_Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *MessageResponse) Reset()         { *m = MessageResponse{} }
func (m *MessageResponse) String() string { return proto.CompactTextString(m) }
func (*MessageResponse) ProtoMessage()    {}
func (*MessageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *MessageResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageResponse.Unmarshal(m, b)
}
func (m *MessageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageResponse.Marshal(b, m, deterministic)
}
func (m *MessageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageResponse.Merge(m, src)
}
func (m *MessageResponse) XXX_Size() int {
	return xxx_messageInfo_MessageResponse.Size(m)
}
func (m *MessageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MessageResponse proto.InternalMessageInfo

func (m *MessageResponse) GetStatus() MessageResponse_Status {
	if m != nil {
		return m.Status
	}
	return MessageResponse_Ok
}

type WsResponse struct {
	Status               WsResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=message.WsResponse_Status" json:"status,omitempty"`
	Message              string            `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *WsResponse) Reset()         { *m = WsResponse{} }
func (m *WsResponse) String() string { return proto.CompactTextString(m) }
func (*WsResponse) ProtoMessage()    {}
func (*WsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}

func (m *WsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WsResponse.Unmarshal(m, b)
}
func (m *WsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WsResponse.Marshal(b, m, deterministic)
}
func (m *WsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WsResponse.Merge(m, src)
}
func (m *WsResponse) XXX_Size() int {
	return xxx_messageInfo_WsResponse.Size(m)
}
func (m *WsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WsResponse proto.InternalMessageInfo

func (m *WsResponse) GetStatus() WsResponse_Status {
	if m != nil {
		return m.Status
	}
	return WsResponse_NONE
}

func (m *WsResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type WsHandshakeRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Secret               string   `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WsHandshakeRequest) Reset()         { *m = WsHandshakeRequest{} }
func (m *WsHandshakeRequest) String() string { return proto.CompactTextString(m) }
func (*WsHandshakeRequest) ProtoMessage()    {}
func (*WsHandshakeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{3}
}

func (m *WsHandshakeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WsHandshakeRequest.Unmarshal(m, b)
}
func (m *WsHandshakeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WsHandshakeRequest.Marshal(b, m, deterministic)
}
func (m *WsHandshakeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WsHandshakeRequest.Merge(m, src)
}
func (m *WsHandshakeRequest) XXX_Size() int {
	return xxx_messageInfo_WsHandshakeRequest.Size(m)
}
func (m *WsHandshakeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WsHandshakeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WsHandshakeRequest proto.InternalMessageInfo

func (m *WsHandshakeRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *WsHandshakeRequest) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

type WsMessageRequest struct {
	TargetId             string      `protobuf:"bytes,1,opt,name=targetId,proto3" json:"targetId,omitempty"`
	Type                 MessageType `protobuf:"varint,2,opt,name=type,proto3,enum=message.MessageType" json:"type,omitempty"`
	Data                 string      `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	SourceId             string      `protobuf:"bytes,4,opt,name=sourceId,proto3" json:"sourceId,omitempty"`
	Ext                  int32       `protobuf:"varint,5,opt,name=ext,proto3" json:"ext,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *WsMessageRequest) Reset()         { *m = WsMessageRequest{} }
func (m *WsMessageRequest) String() string { return proto.CompactTextString(m) }
func (*WsMessageRequest) ProtoMessage()    {}
func (*WsMessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4}
}

func (m *WsMessageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WsMessageRequest.Unmarshal(m, b)
}
func (m *WsMessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WsMessageRequest.Marshal(b, m, deterministic)
}
func (m *WsMessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WsMessageRequest.Merge(m, src)
}
func (m *WsMessageRequest) XXX_Size() int {
	return xxx_messageInfo_WsMessageRequest.Size(m)
}
func (m *WsMessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WsMessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WsMessageRequest proto.InternalMessageInfo

func (m *WsMessageRequest) GetTargetId() string {
	if m != nil {
		return m.TargetId
	}
	return ""
}

func (m *WsMessageRequest) GetType() MessageType {
	if m != nil {
		return m.Type
	}
	return MessageType_NONE
}

func (m *WsMessageRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *WsMessageRequest) GetSourceId() string {
	if m != nil {
		return m.SourceId
	}
	return ""
}

func (m *WsMessageRequest) GetExt() int32 {
	if m != nil {
		return m.Ext
	}
	return 0
}

func init() {
	proto.RegisterEnum("message.MessageType", MessageType_name, MessageType_value)
	proto.RegisterEnum("message.MessageResponse_Status", MessageResponse_Status_name, MessageResponse_Status_value)
	proto.RegisterEnum("message.WsResponse_Status", WsResponse_Status_name, WsResponse_Status_value)
	proto.RegisterType((*MessageRequest)(nil), "message.MessageRequest")
	proto.RegisterType((*MessageResponse)(nil), "message.MessageResponse")
	proto.RegisterType((*WsResponse)(nil), "message.WsResponse")
	proto.RegisterType((*WsHandshakeRequest)(nil), "message.WsHandshakeRequest")
	proto.RegisterType((*WsMessageRequest)(nil), "message.WsMessageRequest")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0x4d, 0x6b, 0xdb, 0x40,
	0x10, 0xf5, 0xea, 0x2b, 0xd1, 0x84, 0xa8, 0xcb, 0x50, 0xda, 0xad, 0x2f, 0x11, 0x0b, 0x05, 0xd1,
	0x83, 0x0f, 0x0a, 0xa5, 0x97, 0xdc, 0x4a, 0x68, 0x03, 0xad, 0x6c, 0x36, 0x29, 0x3e, 0xab, 0xd6,
	0x92, 0x1a, 0x63, 0x4b, 0xdd, 0x5d, 0x95, 0xfa, 0xd8, 0xbf, 0xd1, 0x5f, 0x1b, 0x24, 0xad, 0x84,
	0x2d, 0xe3, 0xdb, 0x3c, 0xe9, 0xed, 0x7b, 0x33, 0x6f, 0x06, 0xae, 0xb7, 0x52, 0xeb, 0xfc, 0x59,
	0xce, 0x2a, 0x55, 0x9a, 0x12, 0x2f, 0x2c, 0xe4, 0x5b, 0x88, 0xbe, 0x77, 0xa5, 0x90, 0xbf, 0x6b,
	0xa9, 0x0d, 0x46, 0xe0, 0xac, 0x0b, 0x46, 0x62, 0x37, 0x09, 0x85, 0xb3, 0x2e, 0xf0, 0x16, 0x7a,
	0x32, 0x73, 0x62, 0x92, 0x5c, 0xa5, 0xef, 0x66, 0xbd, 0xd6, 0x52, 0x1f, 0xbf, 0x15, 0x3d, 0x13,
	0xdf, 0x40, 0xa0, 0xcb, 0x5a, 0xad, 0x24, 0x73, 0x63, 0x92, 0x84, 0xc2, 0x22, 0xbe, 0x81, 0x57,
	0xc3, 0x13, 0x5d, 0x95, 0x3b, 0x2d, 0xf1, 0x13, 0x04, 0xda, 0xe4, 0xa6, 0xd6, 0x8c, 0xc4, 0x24,
	0x89, 0xd2, 0x9b, 0x41, 0x7e, 0xc4, 0x9c, 0x3d, 0xb6, 0x34, 0x61, 0xe9, 0xfc, 0x06, 0x82, 0xee,
	0x0b, 0x06, 0xe0, 0xcc, 0x37, 0x74, 0x82, 0xd7, 0x10, 0x66, 0xa5, 0x11, 0x65, 0x6d, 0x64, 0x41,
	0x09, 0xff, 0x47, 0x00, 0x96, 0x7a, 0x30, 0x4a, 0x47, 0x46, 0xd3, 0x83, 0x39, 0xce, 0x78, 0x20,
	0x3b, 0x1e, 0x3e, 0x1c, 0x26, 0xe4, 0xef, 0x07, 0xf7, 0x4b, 0xf0, 0xb2, 0x79, 0x76, 0x4f, 0x27,
	0xb6, 0x0f, 0x82, 0x21, 0xf8, 0xf7, 0x4a, 0x95, 0x8a, 0x3a, 0xfc, 0x0e, 0x70, 0xa9, 0xbf, 0xe6,
	0xbb, 0x42, 0xff, 0xca, 0x37, 0x27, 0x19, 0x13, 0x9b, 0x71, 0x13, 0x97, 0x5c, 0x29, 0x69, 0xac,
	0x8b, 0x45, 0xfc, 0x3f, 0x01, 0x3a, 0x0e, 0x19, 0xa7, 0x70, 0x69, 0x72, 0xf5, 0x2c, 0xcd, 0x43,
	0x2f, 0x31, 0x60, 0x4c, 0xc0, 0x33, 0xfb, 0xaa, 0x6b, 0x36, 0x4a, 0x5f, 0x8f, 0xa3, 0x7c, 0xda,
	0x57, 0x52, 0xb4, 0x0c, 0x44, 0xf0, 0x8a, 0xdc, 0xe4, 0x76, 0x3f, 0x6d, 0xdd, 0x28, 0x77, 0x7b,
	0x7a, 0x28, 0x98, 0xd7, 0x29, 0xf7, 0x18, 0x29, 0xb8, 0xf2, 0xaf, 0x61, 0x7e, 0x4c, 0x12, 0x5f,
	0x34, 0xe5, 0x87, 0x8f, 0x70, 0x75, 0x20, 0x7b, 0x10, 0xc3, 0x05, 0xb8, 0x8b, 0x74, 0xd1, 0xe5,
	0xf0, 0x45, 0xcc, 0x7f, 0x2c, 0xa8, 0xd3, 0xfc, 0xfd, 0xfc, 0x24, 0xbe, 0x51, 0x37, 0xcd, 0x86,
	0x8b, 0x7b, 0x94, 0xea, 0xcf, 0x7a, 0x25, 0xf1, 0x0e, 0x7c, 0xd5, 0xec, 0x0c, 0xdf, 0x9e, 0xae,
	0xbe, 0x1d, 0x79, 0xca, 0xce, 0xdd, 0x04, 0x9f, 0xfc, 0x0c, 0xda, 0x8b, 0xbe, 0x7d, 0x09, 0x00,
	0x00, 0xff, 0xff, 0x93, 0xb3, 0x48, 0xf9, 0xe2, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MessageServiceClient interface {
	Route(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error)
}

type messageServiceClient struct {
	cc *grpc.ClientConn
}

func NewMessageServiceClient(cc *grpc.ClientConn) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) Route(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, "/message.MessageService/route", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServiceServer is the server API for MessageService service.
type MessageServiceServer interface {
	Route(context.Context, *MessageRequest) (*MessageResponse, error)
}

func RegisterMessageServiceServer(s *grpc.Server, srv MessageServiceServer) {
	s.RegisterService(&_MessageService_serviceDesc, srv)
}

func _MessageService_Route_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).Route(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.MessageService/Route",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Route(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MessageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "message.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "route",
			Handler:    _MessageService_Route_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "message.proto",
}
