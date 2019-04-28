package broadcast

import (
	"context"

	"go-common/app/interface/main/app-resource/conf"
	broaddao "go-common/app/interface/main/app-resource/dao/broadcast"
	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/broadcast"
	warden "go-common/app/service/main/broadcast/api/grpc/v1"
	"go-common/library/log"
)

type Service struct {
	c   *conf.Config
	dao *broaddao.Dao
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: broaddao.New(c),
	}
	return
}

// ServerList warden server list
func (s *Service) ServerList(c context.Context, plat int8) (res *broadcast.ServerListReply, err error) {
	var (
		data     *warden.ServerListReply
		platform string
	)
	if model.IsIOS(plat) {
		platform = "ios"
	} else if model.IsAndroid(plat) {
		platform = "android"
	}
	if data, err = s.dao.ServerList(c, platform); err != nil {
		log.Error("ServerList s.dao.ServerList error(%v)", err)
		err = nil
		res = &broadcast.ServerListReply{
			Domain:       "broadcast.chat.bilibili.com",
			TCPPort:      7821,
			WsPort:       7822,
			WssPort:      7823,
			Heartbeat:    30,
			HeartbeatMax: 3,
			Nodes:        []string{"broadcast.chat.bilibili.com"},
			Backoff: &broadcast.Backoff{
				MaxDelay:  120,
				BaseDelay: 3,
				Factor:    1.6,
				Jitter:    0.2,
			},
		}
		return
	}
	res = &broadcast.ServerListReply{}
	res.ServerListChange(data)
	return
}
