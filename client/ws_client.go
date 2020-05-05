package client

import (
	"github.com/gorilla/websocket"
	"sync"
)

type CodecType int32

var wsClientMap = make(map[string]*WsClient)
var m = sync.RWMutex{}

func AddWsClient(cli *WsClient) {
	m.Lock()
	wsClientMap[cli.Id()] = cli
	m.Unlock()
}

func RemoveWsClient(cli *WsClient) {
	m.Lock()
	delete(wsClientMap, cli.Id())
	m.Unlock()
}

func LookupClient(id string) (*WsClient, bool) {
	m.RLock()
	client, ok := wsClientMap[id]
	m.RUnlock()
	return client, ok
}

type WsClient struct {
	id   string
	conn *websocket.Conn
}

func NewWsClient(id string, conn *websocket.Conn) *WsClient {
	return &WsClient{
		id:   id,
		conn: conn,
	}
}

func (c *WsClient) Id() string {
	return c.id
}

func (c *WsClient) Send(request interface{}) error {
	return c.conn.WriteJSON(request)
}
