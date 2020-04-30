package client

import (
	"github.com/raylax/imx/core"
	pd "github.com/raylax/imx/message"
	"google.golang.org/grpc"
)

type rpcClient struct {
	node core.Node
	cli pd.MessageServiceClient
}

func NewClient(node core.Node) *rpcClient {
	return &rpcClient{node: node}
}

func (c *rpcClient) Init() error {
	conn, err := grpc.Dial(c.node.Address())
	if err != nil {
		return err
	}
	c.cli = pd.NewMessageServiceClient(conn)
	return nil
}

func (c *rpcClient) MessageService() pd.MessageServiceClient {
	return c.cli
}