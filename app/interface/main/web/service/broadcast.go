package service

import (
	"context"

	"go-common/app/interface/main/web/model"
	warden "go-common/app/service/main/broadcast/api/grpc/v1"
	"go-common/library/log"
)

// BroadServers broadcast server list.
func (s *Service) BroadServers(c context.Context, platform string) (res *warden.ServerListReply, err error) {
	if res, err = s.broadcastClient.ServerList(c, &warden.ServerListReq{Platform: platform}); err != nil {
		log.Error("s.broadCastClient.ServerList(%s) error(%v)", platform, err)
		res = model.DefaultServer
		err = nil
	}
	return
}
