package server

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/raylax/imx/client"
	"github.com/raylax/imx/core"
	pb "github.com/raylax/imx/proto"
	"github.com/raylax/imx/registry"
	"github.com/raylax/imx/router"
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
	addr          string
	registry      registry.Registry
	messageRouter router.MessageRouter
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
	// 握手失败直接断开连接
	if err != nil {
		log.Printf("[%s]握手失败：%s", conn.RemoteAddr(), err)
		_ = conn.Close()
		return
	}
	user := core.User{Id: wsCli.Id()}
	err = s.registry.RegUser(user)
	// 向注册中心注册用户失败直接断开连接
	if err != nil {
		log.Printf("注册用户失败：%s", err)
		_ = conn.Close()
		return
	}
	client.AddWsClient(wsCli)
	for {
		message := &pb.WsMessageRequest{}
		// 连接断开
		if err := conn.ReadJSON(message); err != nil {
			client.RemoveWsClient(wsCli)
			s.registry.UnRegUser(user)
			log.Printf("[%s]断开：%s", conn.RemoteAddr(), err)
			return
		}
		// 设置发送者ID
		message.SourceId = wsCli.Id()
		// 生成消息ID
		message.MessageId = uuid.New().String()
		go func() {
			err = s.messageRouter.RouteMessage(message)
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

// 握手
// TODO 待添加鉴权逻辑
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
