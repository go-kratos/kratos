package grpc

import (
	"context"

	v1 "go-common/app/interface/main/app-resource/api/v1"
	"go-common/app/interface/main/app-resource/http"
	modulesvr "go-common/app/interface/main/app-resource/service/module"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

// Server struct
type Server struct {
	moduleSvc *modulesvr.Service
}

// New Coin warden rpc server
func New(c *warden.ServerConfig, svr *http.Server) (wsvr *warden.Server, err error) {
	wsvr = warden.NewServer(c)
	v1.RegisterAppResourceServer(wsvr.Server(), &Server{
		moduleSvc: svr.ModuleSvc,
	})
	wsvr, err = wsvr.Start()
	return
}

// ModuleUpdateCache update module cache
func (s *Server) ModuleUpdateCache(c context.Context, noArg *v1.NoArgRequest) (noReply *v1.NoReply, err error) {
	if err = s.moduleSvc.ModuleUpdateCache(); err != nil {
		log.Error("ModuleUpdateCache error(%v)", err)
		return
	}
	noReply = &v1.NoReply{}
	log.Info("ModuleUpdateCache load cache success")
	return
}
