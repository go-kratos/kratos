package server

import (
	"log"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/service"

	"go-common/library/net/rpc"
)

// New .
func New(config *conf.Config, s service.Service) *rpc.Server {
	rpcSvr := rpc.NewServer(config.RPC)
	if err := rpcSvr.Register(&Filter{svr: s}); err != nil {
		log.Fatalf("%+v", err)
	}
	return rpcSvr
}
