package account

import (
	"context"

	"go-common/app/interface/main/app-view/conf"
	account "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
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
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

func (d *Dao) Cards3(c context.Context, mids []int64) (res map[int64]*account.Card, err error) {
	arg := &account.ArgMids{Mids: mids}
	if res, err = d.accRPC.Cards3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// Following3 following.
func (d *Dao) Following3(c context.Context, mid, owner int64) (follow bool, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &account.ArgRelation{Mid: mid, Owner: owner, RealIP: ip}
	rl, err := d.accRPC.Relation3(c, arg)
	if err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	if rl != nil {
		follow = rl.Following
	}
	return
}

// IsAttention is attention
func (d *Dao) IsAttention(c context.Context, owners []int64, mid int64) (isAtten map[int64]int8) {
	if len(owners) == 0 || mid == 0 {
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &account.ArgRelations{Owners: owners, Mid: mid, RealIP: ip}
	res, err := d.accRPC.Relations3(c, arg)
	if err != nil {
		log.Error("%+v", err)
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
