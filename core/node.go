package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const keySplitChar = "_"

type Node struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

func (n *Node) Key() string {
	return fmt.Sprintf("%s%s%d", n.Addr, keySplitChar, n.Port)
}

func (n *Node) Endpoint() string {
	return fmt.Sprintf("%s:%d", n.Addr, n.Port)
}

func NewNodeFromKey(key string) Node {
	ss := strings.Split(key, keySplitChar)
	port, _ := strconv.Atoi(ss[1])
	return Node{
		Addr: ss[0],
		Port: port,
	}
}

func NewNodeFromJSON(data []byte) Node {
	node := Node{}
	_ = json.Unmarshal(data, node)
	return node
}
