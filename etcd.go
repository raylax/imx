package main

import (
	"flag"
	"github.com/raylax/imx/core"
	"github.com/raylax/imx/registry"
	"github.com/raylax/imx/server"
	"log"
	"strconv"
	"strings"
)

var registryAddress = flag.String("registry", "127.0.0.1:2379", "Registry address")
var rpcListen = flag.String("rpc-listen", ":9321", "RPC service listen address")
var rpcEndpoint = flag.String("rpc-endpoint", "127.0.0.1:9321", "RPC service endpoint address")
var wsListen = flag.String("ws-listen", ":8080", "Websocket service listen address")

func main() {
	flag.Parse()
	ss := strings.Split(*rpcEndpoint, ":")
	port, _ := strconv.Atoi(ss[1])
	reg := registry.NewEtcdRegistry(strings.Split(*registryAddress, ","), core.Node{Addr: ss[0], Port: port})
	if err := reg.Init(); err != nil {
		log.Fatal(err)
	}
	s := server.NewServer(reg, *wsListen, *rpcListen)
	err := s.Serve()
	if err != nil {
		log.Fatal(err)
	}

}
