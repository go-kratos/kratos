package pgc

import (
	"go-common/app/admin/main/feed/conf"
	epgrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
)

// Dao is show dao.
type Dao struct {
	// grpc
	rpcClient seasongrpc.SeasonClient
	epClient  epgrpc.EpisodeClient
}

// New new a bangumi dao.
func New(c *conf.Config) (*Dao, error) {
	var ep epgrpc.EpisodeClient
	rpcClient, err := seasongrpc.NewClient(nil)
	if err != nil {
		log.Error("seasongrpc NewClientt error(%v)", err)
		return nil, err
	}
	if ep, err = epgrpc.NewClient(nil); err != nil {
		log.Error("eprpc NewClientt error(%v)", err)
		return nil, err
	}
	d := &Dao{
		rpcClient: rpcClient,
		epClient:  ep,
	}
	return d, nil
}
