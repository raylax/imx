package server

import (
	"github.com/raylax/imx/registry"
	"log"
)

type Server interface {
	Serve() error
}

type server struct {
	wsAddr    string
	rpcAddr   string
	rpcServer *rpcServer
	wsServer  *wsServer
	registry  registry.Registry
}

func NewServer(registry registry.Registry, wsAddr, rpcAddr string) *server {
	return &server{
		registry: registry,
		wsAddr:   wsAddr,
		rpcAddr:  rpcAddr,
	}
}

func (s *server) Serve() error {
	var errChain = make(chan error)
	log.Printf("启动RPC服务")
	s.rpcServer = &rpcServer{addr: s.rpcAddr, registry: s.registry}
	go func(ch chan error) {
		ch <- s.rpcServer.Serve()
	}(errChain)
	log.Printf("启动WS服务")
	s.wsServer = &wsServer{addr: s.wsAddr, registry: s.registry}
	go func(ch chan error) {
		ch <- s.wsServer.Serve()
	}(errChain)
	if err := s.registry.Reg(); err != nil {
		return err
	}
	err := <-errChain
	s.registry.UnReg()
	return err
}
