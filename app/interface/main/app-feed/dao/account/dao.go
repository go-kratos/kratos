package account

import (
	"context"

	"go-common/app/interface/main/app-feed/conf"
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

// Relations3 relations.
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

func (d *Dao) Cards3(c context.Context, mids []int64) (res map[int64]*account.Card, err error) {
	arg := &account.ArgMids{Mids: mids}
	if res, err = d.accRPC.Cards3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
