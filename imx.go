package main

import (
	"flag"
	"github.com/raylax/imx/core"
	"github.com/raylax/imx/registry"
	"github.com/raylax/imx/server"
	"github.com/raylax/imx/version"
	"log"
	"os"
	"strconv"
	"strings"
)

var registryAddress string
var rpcListen string
var rpcEndpoint string
var wsListen string

var h bool
var v bool

func init() {
	flag.StringVar(&registryAddress, "registry", "127.0.0.1:2379", "Registry address")
	flag.StringVar(&rpcListen, "rpc-listen", ":9321", "RPC service listen address")
	flag.StringVar(&rpcEndpoint, "rpc-endpoint", "127.0.0.1:9321", "RPC service endpoint address")
	flag.StringVar(&wsListen, "ws-listen", ":8080", "Websocket service listen address")
	flag.BoolVar(&h, "h", false, "print help")
	flag.BoolVar(&v, "v", false, "print version")
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		os.Exit(0)
	}
	if v {
		println("IM-X\n" + version.String())
		os.Exit(0)
	}
	ss := strings.Split(rpcEndpoint, ":")
	port, _ := strconv.Atoi(ss[1])
	reg := registry.NewEtcdRegistry(strings.Split(registryAddress, ","), core.Node{Addr: ss[0], Port: port})
	if err := reg.Init(); err != nil {
		log.Fatal(err)
	}
	s := server.NewServer(reg, wsListen, rpcListen)
	err := s.Serve()
	if err != nil {
		log.Fatal(err)
	}

}
