package router

import (
	"context"
	"fmt"
	"github.com/raylax/imx/client"
	"github.com/raylax/imx/core"
	pb "github.com/raylax/imx/proto"
	"github.com/raylax/imx/registry"
	"log"
)

type MessageRouter struct {
	registry registry.Registry
}

func NewMessageRouter(registry registry.Registry) MessageRouter {
	return MessageRouter{registry: registry}
}

func (s *MessageRouter) RouteMessage(message *pb.WsMessageRequest) error {
	switch message.Type {
	// 路由点到点消息
	case pb.MessageType_P2P:
		return s.routeP2PMessage(message)
	// 路由群消息
	case pb.MessageType_GROUP:
		return s.routeGroupMessage(message)
	default:
		return fmt.Errorf("unsupported message type `%s`", message.Type)
	}
}

func (s *MessageRouter) routeGroupMessage(message *pb.WsMessageRequest) error {
	users, err := s.registry.GetGroupUsers(message.GetTargetId())
	if err != nil {
		return err
	}
	remoteIds := make([]string, 0, len(users))
	// 优先发送本地节点
	for _, u := range users {
		if sent, _ := s.routeMessageToLocalNode(message, u); sent {
			continue
		}
		remoteIds = append(remoteIds, u)
	}
	// 剩余的发送到远程节点
	return s.routeMessageToRemoteNode(message, remoteIds)
}

func (s *MessageRouter) routeP2PMessage(message *pb.WsMessageRequest) error {
	localed, err := s.routeMessageToLocalNode(message, message.TargetId)
	// 如果发生错误或通过本地节点发送成功则直接返回
	if err != nil || localed {
		return err
	}
	return s.routeMessageToRemoteNode(message, []string{message.TargetId})
}

func (s *MessageRouter) routeMessageToLocalNode(message *pb.WsMessageRequest, targetId string) (bool, error) {
	if cli, found := client.LookupClient(targetId); found {
		return true, cli.Send(message)
	}
	return false, nil
}

// 如果多客户端在同一节点，则合并发送
func (s *MessageRouter) routeMessageToRemoteNode(message *pb.WsMessageRequest, targetIds []string) error {
	var targetMap = make(map[string][]string)
	var nodeMap = make(map[string]core.Node)
	for _, target := range targetIds {
		nodes, err := s.registry.LookupNodes(target)
		if err != nil {
			continue
		}
		for _, n := range nodes {
			targets, ok := targetMap[n.Endpoint()]
			if !ok {
				targets = make([]string, 0, len(targetIds)*3)
			}
			// 合并节点
			targetMap[n.Endpoint()] = append(targets, target)
			_, ok = nodeMap[n.Endpoint()]
			if !ok {
				nodeMap[n.Endpoint()] = n
			}
		}
	}
	for endpoint, targetIds := range targetMap {
		rpcClient := client.GetRpcClient(nodeMap[endpoint])
		service := rpcClient.MessageService()
		go func() {
			_, err := service.Route(context.Background(), &pb.MessageRequest{TargetIds: targetIds, Message: message})
			if err != nil {
				log.Printf("failed to route message, error:%s", err)
			}
		}()
	}
	return nil
}
