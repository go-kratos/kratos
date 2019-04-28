package ai

import (
	"context"

	"go-common/app/interface/main/app-view/conf"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_av2GameURL = "/avid2gameid"
)

type Dao struct {
	client     *bm.Client
	av2GameURL string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:     bm.NewClient(c.HTTPGameAsync),
		av2GameURL: c.Host.AI + _av2GameURL,
	}
	return
}

func (d *Dao) Av2Game(c context.Context) (res map[int64]int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if err = d.client.Get(c, d.av2GameURL, ip, nil, &res); err != nil {
		return
	}
	return
}
