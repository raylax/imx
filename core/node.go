package core

import (
	"encoding/json"
	"fmt"
	"github.com/raylax/imx/random"
	"os"
)

const keySplitChar = "_"

type Node struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
	Name string `json:"name"`
}

func NewNode(addr string, port int) Node {
	name, err := os.Hostname()
	if err != nil {
		name = fmt.Sprintf("%s_%d", addr, port)
	}
	name = fmt.Sprintf("%s-%s", name, random.Letter(8))
	return Node{
		Addr: addr,
		Port: port,
		Name: name,
	}
}

func (n *Node) Key() string {
	return fmt.Sprintf("%s%s%d", n.Addr, keySplitChar, n.Port)
}

func (n *Node) Endpoint() string {
	return fmt.Sprintf("%s:%d", n.Addr, n.Port)
}

func NewNodeFromJSON(data []byte) Node {
	node := Node{}
	_ = json.Unmarshal(data, &node)
	return node
}
