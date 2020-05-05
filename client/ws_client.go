package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	pb "github.com/raylax/imx/proto"
	"sync"
)

type CodecType int32

const (
	CodecTypeProtobuf CodecType = 1
	CodecTypeJSON     CodecType = 2
)

var wsClientMap = make(map[string]*wsClient)
var m = sync.RWMutex{}

func AddWsClient(id string, conn *websocket.Conn, codecType CodecType) {
	m.Lock()
	wsClientMap[id] = &wsClient{
		id:        id,
		conn:      conn,
		codecType: codecType,
	}
	m.Unlock()
}

func RemoveWsClient(id string) {
	m.Lock()
	delete(wsClientMap, id)
	m.Unlock()
}

func LookupClient(id string) (*wsClient, bool) {
	m.RLock()
	client, ok := wsClientMap[id]
	m.RUnlock()
	return client, ok
}

type wsClient struct {
	id        string
	conn      *websocket.Conn
	codecType CodecType
}

func (c *wsClient) Send(request *pb.WsMessageRequest) error {
	var data []byte
	var err error
	var messageType int
	switch c.codecType {
	case CodecTypeJSON:
		data, err = json.Marshal(request)
		messageType = websocket.TextMessage
	case CodecTypeProtobuf:
		data, err = proto.Marshal(request)
		messageType = websocket.BinaryMessage
	default:
		return errors.New(fmt.Sprintf("Unsupported codec type %d", c.codecType))
	}
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(messageType, data)
}
