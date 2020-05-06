package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/raylax/imx/client"
	"github.com/raylax/imx/core"
	pb "github.com/raylax/imx/proto"
	"github.com/raylax/imx/registry"
	"log"
	"net/http"
)

const (
	maxMessageSize  = 512  // 512kb
	readBufferSize  = 1024 // 1m
	writeBufferSize = 1024 // 1m
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBufferSize,
	WriteBufferSize: writeBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsServer struct {
	addr     string
	registry registry.Registry
}

func (s *wsServer) Serve() error {
	http.HandleFunc("/imx", s.serveWs)
	return http.ListenAndServe(s.addr, nil)
}

func (s *wsServer) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("协议升级错误:%s", err.Error())
		return
	}
	go s.handleConn(conn)
}

func (s *wsServer) handleConn(conn *websocket.Conn) {
	conn.SetReadLimit(maxMessageSize)
	wsCli, err := s.handleHandshake(conn)
	if err != nil {
		log.Printf("[%s]握手失败：%s", conn.RemoteAddr(), err)
		_ = conn.Close()
		return
	}
	user := core.User{Id: wsCli.Id()}
	err = s.registry.RegUser(user)
	if err != nil {
		log.Printf("注册用户失败：%s", err)
		_ = conn.Close()
		return
	}
	client.AddWsClient(wsCli)
	for {
		message := &pb.WsMessageRequest{}
		if err := conn.ReadJSON(message); err != nil {
			client.RemoveWsClient(wsCli)
			s.registry.UnRegUser(user)
			log.Printf("[%s]断开：%s", conn.RemoteAddr(), err)
			return
		}
		message.SourceId = wsCli.Id()
		message.MessageId = uuid.New().String()
		go func() {
			err = s.sendMessage(message)
			if err != nil {
				log.Printf("Failed to route message from [%s] to [%s], error: %s", message.SourceId, message.TargetId, err)
			}
		}()
		resp := pb.WsResponse{
			Status:  pb.WsResponse_Ok,
			Message: message.MessageId,
		}
		_ = wsCli.Send(resp)
	}

}

func (s *wsServer) sendMessage(message *pb.WsMessageRequest) error {
	switch message.Type {
	case pb.MessageType_P2P:
		return s.sendP2PMessage(message)
	case pb.MessageType_GROUP:
		return s.sendGroupMessage(message)
	default:
		return errors.New(fmt.Sprintf("unsupported message type `%s`", message.Type))
	}
}

func (s *wsServer) sendGroupMessage(message *pb.WsMessageRequest) error {
	users, err := s.registry.GetGroupUsers(message.GetTargetId())
	if err != nil {
		return err
	}

	remoteIds := make([]string, 0, len(users))
	for _, u := range users {
		// 过滤掉本地节点
		if sent, _ := s.sendMessageToLocalNode(message, u); sent {
			continue
		}
		remoteIds = append(remoteIds, u)
	}
	return s.sendMessageToRemoteNode(message, remoteIds)
}

func (s *wsServer) sendP2PMessage(message *pb.WsMessageRequest) error {
	localed, err := s.sendMessageToLocalNode(message, message.TargetId)
	if err != nil || localed {
		return err
	}
	return s.sendMessageToRemoteNode(message, []string{message.TargetId})
}

func (s *wsServer) sendMessageToLocalNode(message *pb.WsMessageRequest, targetId string) (bool, error) {
	if cli, found := client.LookupClient(targetId); found {
		return true, cli.Send(message)
	}
	return false, nil
}

// 如果多客户端在同一节点，则合并发送
func (s *wsServer) sendMessageToRemoteNode(message *pb.WsMessageRequest, targetIds []string) error {
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

func (s *wsServer) handleHandshake(conn *websocket.Conn) (*client.WsClient, error) {
	message := &pb.WsHandshakeRequest{}
	if err := conn.ReadJSON(message); err != nil {
		return nil, err
	}
	wsCli := client.NewWsClient(message.Id, conn)
	resp := &pb.WsResponse{Status: pb.WsResponse_Ok}
	if err := wsCli.Send(resp); err != nil {
		return nil, err
	}
	return wsCli, nil
}
