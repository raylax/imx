// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc.proto

package proto

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

func init() { proto.RegisterFile("rpc.proto", fileDescriptor_77a6da22d6a3feb1) }

var fileDescriptor_77a6da22d6a3feb1 = []byte{
	// 103 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x2a, 0x48, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0xbc, 0xb9, 0xa9, 0xc5, 0xc5, 0x89,
	0xe9, 0xa9, 0x10, 0x51, 0x23, 0x2f, 0x2e, 0x3e, 0x5f, 0x88, 0x40, 0x70, 0x6a, 0x51, 0x59, 0x66,
	0x72, 0xaa, 0x90, 0x05, 0x17, 0x6b, 0x51, 0x7e, 0x69, 0x49, 0xaa, 0x90, 0x28, 0x44, 0x89, 0x1e,
	0x54, 0x3e, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8, 0x44, 0x4a, 0x0c, 0x5d, 0xb8, 0xb8, 0x20, 0x3f,
	0xaf, 0x38, 0x55, 0x89, 0x21, 0x89, 0x0d, 0x2c, 0x61, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x48,
	0xce, 0x30, 0x8f, 0x75, 0x00, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/proto.MessageService/route", in, out, opts...)
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
		FullMethod: "/proto.MessageService/Route",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).Route(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MessageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "route",
			Handler:    _MessageService_Route_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}
