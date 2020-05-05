package client

import (
	"github.com/raylax/imx/core"
	pb "github.com/raylax/imx/proto"
	"google.golang.org/grpc"
)

type rpcClient struct {
	node core.Node
	cli  pb.MessageServiceClient
}

func NewRpcClient(node core.Node) *rpcClient {
	return &rpcClient{node: node}
}

func (c *rpcClient) Init() error {
	conn, err := grpc.Dial(c.node.Address())
	if err != nil {
		return err
	}
	c.cli = pb.NewMessageServiceClient(conn)
	return nil
}

func (c *rpcClient) MessageService() pb.MessageServiceClient {
	return c.cli
}
