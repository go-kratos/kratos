package relation

import (
	"context"

	"go-common/app/interface/main/app-feed/conf"
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
func (d *Dao) Stats(ctx context.Context, mids []int64) (res map[int64]*relation.Stat, err error) {
	arg := &relation.ArgMids{Mids: mids}
	if res, err = d.relRPC.Stats(ctx, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
