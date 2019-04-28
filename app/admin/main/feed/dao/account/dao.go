package account

import (
	"context"

	"go-common/app/admin/main/feed/conf"
	account "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
)

// Dao is account dao.
type Dao struct {
	// rpc
	accRPC *accrpc.Service3
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		accRPC: accrpc.New3(c.AccountRPC),
	}
	return
}

// Card3 get card info by mid
func (d *Dao) Card3(c context.Context, mid int64) (res *account.Card, err error) {
	arg := &account.ArgMid{Mid: mid}
	if res, err = d.accRPC.Card3(c, arg); err != nil {
		return
	}
	return
}
