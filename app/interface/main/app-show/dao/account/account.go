package account

import (
	"context"

	"go-common/app/interface/main/app-show/conf"
	account "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// Dao is rpc dao.
type Dao struct {
	// account rpc
	accRPC *accrpc.Service3
}

// New new a account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// account rpc
		accRPC: accrpc.New3(c.AccountRPC),
	}
	return
}

// Cards3 users card
func (d *Dao) Cards3(ctx context.Context, mids []int64) (res map[int64]*account.Card, err error) {
	arg := &account.ArgMids{
		Mids: mids,
	}
	if res, err = d.accRPC.Cards3(ctx, arg); err != nil {
		log.Error("d.accRPC.Infos(%v) error(%v)", arg, err)
		res = nil
		return
	}
	return
}

// Relations3 users info
func (d *Dao) Relations3(ctx context.Context, mid int64, owners []int64) (res map[int64]*account.Relation, err error) {
	arg := &account.ArgRelations{
		Mid:    mid,
		Owners: owners,
	}
	if res, err = d.accRPC.Relations3(ctx, arg); err != nil {
		log.Error("d.accRPC.Relations2(%v) error(%v)", arg, err)
		res = nil
		return
	}
	return
}

func (d *Dao) IsAttention(c context.Context, owners []int64, mid int64) (isAtten map[int64]int8) {
	if len(owners) == 0 || mid == 0 {
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &account.ArgRelations{Owners: owners, Mid: mid, RealIP: ip}
	res, err := d.accRPC.Relations3(c, arg)
	if err != nil {
		return
	}
	isAtten = make(map[int64]int8, len(res))
	for mid, rel := range res {
		if rel.Following {
			isAtten[mid] = 1
		}
	}
	return
}
