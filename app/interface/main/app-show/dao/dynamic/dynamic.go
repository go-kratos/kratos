package dynamic

import (
	"context"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/service/main/archive/api"
	dynarc "go-common/app/service/main/dynamic/model"
	dynrpc "go-common/app/service/main/dynamic/rpc/client"
	"go-common/library/log"
)

// Dao is rpc dao.
type Dao struct {
	// dynamic rpc
	dynRpc *dynrpc.Service
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// dynamic rpc
		dynRpc: dynrpc.New(c.DynamicRPC),
	}
	return
}

// regionDynamic
func (d *Dao) RegionDynamic(ctx context.Context, rid, pn, ps int) (res []*api.Arc, aids []int64, err error) {
	arg := &dynarc.ArgRegion3{
		RegionID: int32(rid),
		Pn:       pn,
		Ps:       ps,
	}
	var as *dynarc.DynamicArcs3
	if as, err = d.dynRpc.RegionArcs3(ctx, arg); err != nil {
		log.Error("d.dynRpc.RegionArcs(%v) error(%v)", arg, err)
		return
	}
	if as != nil {
		res = as.Archives
		for _, a := range res {
			aids = append(aids, a.Aid)
		}
	}
	return
}
