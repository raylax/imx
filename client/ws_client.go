package client

import (
	"github.com/gorilla/websocket"
	"sync"
)

type CodecType int32

var wsClientMap = make(map[string]*WsClient)
var wsMux = sync.RWMutex{}

func AddWsClient(cli *WsClient) {
	wsMux.Lock()
	wsClientMap[cli.Id()] = cli
	wsMux.Unlock()
}

func RemoveWsClient(cli *WsClient) {
	wsMux.Lock()
	delete(wsClientMap, cli.Id())
	wsMux.Unlock()
}

func LookupClient(id string) (*WsClient, bool) {
	wsMux.RLock()
	client, ok := wsClientMap[id]
	wsMux.RUnlock()
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
