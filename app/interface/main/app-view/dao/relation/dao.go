package relation

import (
	"context"

	"go-common/app/interface/main/app-view/conf"
	relation "go-common/app/service/main/relation/model"
	relrpc "go-common/app/service/main/relation/rpc/client"

	"github.com/pkg/errors"
)

type Dao struct {
	relationRPC *relrpc.Service
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		relationRPC: relrpc.New(c.RelationRPC),
	}
	return
}

// Prompt prompt
func (d *Dao) Prompt(c context.Context, mid, vmid int64, btype int8) (prompt bool, err error) {
	arg := &relation.ArgPrompt{Mid: mid, Fid: vmid, Btype: btype}
	if prompt, err = d.relationRPC.Prompt(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// Stat stat
func (d *Dao) Stat(c context.Context, mid int64) (stat *relation.Stat, err error) {
	if stat, err = d.relationRPC.Stat(c, &relation.ArgMid{Mid: mid}); err != nil {
		err = errors.Wrapf(err, "%v", mid)
	}
	return
}

// Stats fids stats
func (d *Dao) Stats(ctx context.Context, mids []int64) (res map[int64]*relation.Stat, err error) {
	arg := &relation.ArgMids{Mids: mids}
	if res, err = d.relationRPC.Stats(ctx, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
