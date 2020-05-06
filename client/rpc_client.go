package client

import (
	"github.com/raylax/imx/core"
	pb "github.com/raylax/imx/proto"
	"google.golang.org/grpc"
	"sync"
)

var rpcClientMap = make(map[string]*RpcClient)
var rpcMux = sync.RWMutex{}

type RpcClient struct {
	node core.Node
	cli  pb.MessageServiceClient
}

func newRpcClient(node core.Node) *RpcClient {
	return &RpcClient{node: node}
}

func (c *RpcClient) init() error {
	conn, err := grpc.Dial(c.node.Endpoint(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.cli = pb.NewMessageServiceClient(conn)
	return nil
}

func (c *RpcClient) MessageService() pb.MessageServiceClient {
	return c.cli
}

func AddRpcClient(node core.Node) *RpcClient {
	client := newRpcClient(node)
	err := client.init()
	if err != nil {
		return nil
	}
	rpcMux.Lock()
	rpcClientMap[node.Endpoint()] = client
	rpcMux.Unlock()
	return client
}

func RemoveRpcClient(node core.Node) {
	rpcMux.Lock()
	delete(rpcClientMap, node.Endpoint())
	rpcMux.Unlock()
}

func GetRpcClient(node core.Node) *RpcClient {
	rpcMux.RLock()
	rpcCli, ok := rpcClientMap[node.Endpoint()]
	rpcMux.RUnlock()
	if !ok {
		return AddRpcClient(node)
	}
	return rpcCli
}
