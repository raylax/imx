package server

import (
	"context"
	"github.com/raylax/imx/client"
	pb "github.com/raylax/imx/proto"
	"github.com/raylax/imx/registry"
	"google.golang.org/grpc"
	"net"
)

type grpcServer struct {
}

func (s *grpcServer) Route(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	for _, id := range req.TargetIds {
		wsClient, ok := client.LookupClient(id)
		// 如果客户端不在线直接跳过
		if !ok {
			continue
		}
		_ = wsClient.Send(req.Message)
	}
	return &pb.MessageResponse{Status: pb.MessageResponse_Ok}, nil
}

type rpcServer struct {
	addr     string
	listen   net.Listener
	registry registry.Registry
}

func (s *rpcServer) Serve() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterMessageServiceServer(server, &grpcServer{})
	s.listen = listen
	return server.Serve(listen)
}
