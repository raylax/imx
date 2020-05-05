package server

import (
	"github.com/gorilla/websocket"
	"github.com/raylax/imx/client"
	"github.com/raylax/imx/core"
	pb "github.com/raylax/imx/proto"
	"github.com/raylax/imx/registry"
	"log"
	"net/http"
)

const (
	maxMessageSize  = 512
	readBufferSize  = 1024
	writeBufferSize = 1024
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
		err = s.sendMessage(message)
	}

}

func (s *wsServer) sendMessage(message *pb.WsMessageRequest) error {
	// 如果在当前服务找到接收者客户端则直接发送
	if cli, found := client.LookupClient(message.TargetId); found {
		return cli.Send(message)
	}
	// 发送信息到远程节点
	return s.sendMessageToRemoteNode(message)
}

func (s *wsServer) sendMessageToRemoteNode(message *pb.WsMessageRequest) error {
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
