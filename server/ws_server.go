package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
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
	codecType, id, err := s.handleHandshake(conn)
	if err != nil {
		log.Printf("[%s]握手失败：%s", conn.RemoteAddr(), err)
		_ = conn.Close()
		return
	}
	user := core.User{Id: id}
	err = s.registry.RegUser(user)
	if err != nil {
		log.Printf("注册用户失败：%s", err)
		_ = conn.Close()
		return
	}
	client.AddWsClient(id, conn, codecType)
	for {
		t, data, err := readBytes(conn)
		if err != nil {
			client.RemoveWsClient(id)
			s.registry.UnRegUser(user)
			log.Printf("[%s]断开：%s", conn.RemoteAddr(), err)
			return
		}
		message := &pb.WsMessageRequest{}
		switch t {
		case websocket.TextMessage:
			err = json.Unmarshal(data, message)
		case websocket.BinaryMessage:
			err = proto.Unmarshal(data, message)
		}
		message.SourceId = id
		err = s.sendMessage(message)
	}

}

func (s *wsServer) sendMessage(message *pb.WsMessageRequest) error {
	cli, ok := client.LookupClient(message.TargetId)
	if ok {
		return cli.Send(message)
	}
	return s.sendMessageToRemoteNode(message)
}

func (s *wsServer) sendMessageToRemoteNode(message *pb.WsMessageRequest) error {
	return nil
}

func (s *wsServer) handleHandshake(conn *websocket.Conn) (client.CodecType, string, error) {
	t, data, err := readBytes(conn)
	var code client.CodecType = -1
	if err != nil {
		return code, "", err
	}
	message := &pb.WsHandshakeRequest{}
	switch t {
	case websocket.TextMessage:
		err = json.Unmarshal(data, message)
		code = client.CodecTypeJSON
	case websocket.BinaryMessage:
		err = proto.Unmarshal(data, message)
		code = client.CodecTypeProtobuf
	default:
		return code, "", errors.New(fmt.Sprintf("Unsupported message type %d", t))
	}
	if err != nil {
		return code, "", err
	}
	resp := &pb.WsResponse{Status: pb.WsResponse_Ok}
	err = conn.WriteJSON(resp)
	if err != nil {
		return code, "", err
	}
	return code, message.Id, nil
}

func readBytes(conn *websocket.Conn) (int, []byte, error) {
	t, data, err := conn.ReadMessage()
	if err != nil {
		return -1, nil, err
	}
	switch t {
	case websocket.TextMessage, websocket.BinaryMessage:
		return t, data, err
	}
	return -1, nil, errors.New(fmt.Sprintf("Unknown message type '%d'", t))
}
