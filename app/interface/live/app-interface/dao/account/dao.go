package account

import (
	"context"

	"go-common/app/interface/live/app-interface/conf"
	account "go-common/app/service/main/account/rpc/client"
)

// Dao dao
type Dao struct {
	c          *conf.Config
	accountRPC *account.Service3
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:          c,
		accountRPC: account.New3(nil),
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
