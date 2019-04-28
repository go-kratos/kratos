package xanchor

import (
	"context"
	xanchor "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/app/service/live/xuser/conf"
)

// Dao dao
type Dao struct {
	c         *conf.Config
	xuserGRPC *xanchor.Client
}

var _rsCli *xanchor.Client

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	var err error
	if _rsCli, err = xanchor.NewClient(c.XanchorClient); err != nil {
		panic(err)
	}
	dao = &Dao{
		c:         c,
		xuserGRPC: _rsCli,
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	// check
	return nil
}
