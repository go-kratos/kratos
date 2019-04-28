package xuser

import (
	"context"

	"go-common/app/interface/live/app-interface/conf"
	xuser "go-common/app/service/live/xuser/api/grpc/v1"
)

// Dao dao
type Dao struct {
	c         *conf.Config
	xuserGRPC *xuser.Client
}

var _rsCli *xuser.Client

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	var err error
	if _rsCli, err = xuser.NewClient(c.XuserClient); err != nil {
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
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	// check
	return nil
}
