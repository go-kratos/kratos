package relation

import (
	"context"

	"go-common/app/admin/main/credit/conf"
	relationmdl "go-common/app/service/main/relation/model"
	relationrpc "go-common/app/service/main/relation/rpc/client"
	"go-common/library/log"
)

// Dao is account dao.
type Dao struct {
	// rpc
	relationRPC *relationrpc.Service
}

// New is initial for account .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		relationRPC: relationrpc.New(c.RPCClient.Relation),
	}
	return
}

// RPCStats rpc info get by  muti mid .
func (d *Dao) RPCStats(c context.Context, mids []int64) (res map[int64]*relationmdl.Stat, err error) {
	arg := &relationmdl.ArgMids{Mids: mids}
	if res, err = d.relationRPC.Stats(c, arg); err != nil {
		log.Error("d.relationRPC.Stats error(%v)", err)
	}
	return
}
