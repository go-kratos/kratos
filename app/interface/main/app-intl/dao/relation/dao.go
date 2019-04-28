package relation

import (
	"context"

	"go-common/app/interface/main/app-intl/conf"
	relation "go-common/app/service/main/relation/model"
	relrpc "go-common/app/service/main/relation/rpc/client"

	"github.com/pkg/errors"
)

// Dao is rpc dao.
type Dao struct {
	// relation rpc
	relRPC *relrpc.Service
}

// New new a relation dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// relation rpc
		relRPC: relrpc.New(c.RelationRPC),
	}
	return
}

// Stats fids stats
func (d *Dao) Stats(c context.Context, mids []int64) (res map[int64]*relation.Stat, err error) {
	arg := &relation.ArgMids{Mids: mids}
	if res, err = d.relRPC.Stats(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// Stat stat
func (d *Dao) Stat(c context.Context, mid int64) (stat *relation.Stat, err error) {
	if stat, err = d.relRPC.Stat(c, &relation.ArgMid{Mid: mid}); err != nil {
		err = errors.Wrapf(err, "%v", mid)
	}
	return
}
