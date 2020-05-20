package server

import (
	"github.com/raylax/imx/registry"
	"github.com/raylax/imx/router"
	"log"
	"os"
)

type Server interface {
	Serve() error
	Shutdown()
}

type server struct {
	wsAddr     string
	rpcAddr    string
	apiAddr    string
	rpcServer  *rpcServer
	wsServer   *wsServer
	apiServer  *apiServer
	registry   registry.Registry
	shutdownCh chan os.Signal
}

func NewServer(registry registry.Registry, wsAddr, rpcAddr string, apiAddr string, shutdownCh chan os.Signal) *server {
	return &server{
		registry:   registry,
		wsAddr:     wsAddr,
		rpcAddr:    rpcAddr,
		apiAddr:    apiAddr,
		shutdownCh: shutdownCh,
	}
}

func (s *server) Serve() error {
	var errChain = make(chan error)

	messageRouter := router.NewMessageRouter(s.registry)

	log.Printf("启动RPC服务")
	s.rpcServer = &rpcServer{addr: s.rpcAddr, registry: s.registry}
	go func(ch chan error) {
		ch <- s.rpcServer.Serve()
	}(errChain)

	log.Printf("启动WS服务")
	s.wsServer = &wsServer{addr: s.wsAddr, registry: s.registry, messageRouter: messageRouter}
	go func(ch chan error) {
		ch <- s.wsServer.Serve()
	}(errChain)

	log.Printf("启动API服务")
	s.apiServer = &apiServer{addr: s.apiAddr, registry: s.registry, messageRouter: messageRouter}
	go func(ch chan error) {
		ch <- s.apiServer.Serve()
	}(errChain)

	if err := s.registry.Reg(); err != nil {
		return err
	}
	var err error
	select {
	case err = <-errChain:
	case <-s.shutdownCh:
		s.registry.UnReg()
		return nil
	}
	s.registry.UnReg()
	return err
}
