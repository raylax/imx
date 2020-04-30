package server

import (
	"context"
	pd "github.com/raylax/imx/message"
	"github.com/raylax/imx/registry"
	"google.golang.org/grpc"
	"log"
	"net"
)

type grpcServer struct {

}

func (s *grpcServer) Route(ctx context.Context, req *pd.MessageRequest) (*pd.MessageResponse, error) {
	log.Printf("Received: %v", req.Id)
	return &pd.MessageResponse{Status: pd.MessageResponse_Ok}, nil
}

type rpcServer struct {
	addr string
	listen net.Listener
	registry registry.Registry
}

func (s *rpcServer) Serve() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pd.RegisterMessageServiceServer(server, &grpcServer{})
	s.listen = listen
	return server.Serve(listen)
}
