package account

import (
	"context"

	"go-common/app/interface/main/app-intl/conf"
	account "go-common/app/service/main/account/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is archive dao.
type Dao struct {
	// rpc
	accRPC *accrpc.Service3
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// rpc
		accRPC: accrpc.New3(c.AccountRPC),
	}
	return
}

// Relations3 relatons
func (d *Dao) Relations3(c context.Context, owners []int64, mid int64) (follows map[int64]bool) {
	if len(owners) == 0 {
		return nil
	}
	follows = make(map[int64]bool, len(owners))
	for _, owner := range owners {
		follows[owner] = false
	}
	var (
		am  map[int64]*account.Relation
		err error
	)
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &account.ArgRelations{Owners: owners, Mid: mid, RealIP: ip}
	if am, err = d.accRPC.Relations3(c, arg); err != nil {
		log.Error("%+v", err)
		return
	}
	for i, a := range am {
		if _, ok := follows[i]; ok {
			follows[i] = a.Following
		}
	}
	return
}

// IsAttention is
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

// Card3 get card info by mid
func (d *Dao) Card3(c context.Context, mid int64) (res *account.Card, err error) {
	arg := &account.ArgMid{Mid: mid}
	if res, err = d.accRPC.Card3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// Cards3 get cards info by mids
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

// Infos3 rpc info get by mids .
func (d *Dao) Infos3(c context.Context, mids []int64) (res map[int64]*account.Info, err error) {
	arg := &account.ArgMids{Mids: mids}
	if res, err = d.accRPC.Infos3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
