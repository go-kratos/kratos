package archive

import (
	"context"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

// Dao is archive dao.
type Dao struct {
	c *conf.Config
	// rpc
	arcRpc *arcrpc.Service2
}

// New new a archive dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// rpc
		arcRpc: arcrpc.New2(c.ArchiveRPC),
	}
	return
}

// Archive get archive by aid.
func (d *Dao) Archive(ctx context.Context, aid int64) (a *api.Arc, err error) {
	arg := &archive.ArgAid2{Aid: aid}
	if a, err = d.arcRpc.Archive3(ctx, arg); err != nil {
		log.Error("d.arcRpc.Archive3(%v) error(%v)", arg, err)
		return
	}
	return
}

// ArchivesPB multi get archives.
func (d *Dao) ArchivesPB(ctx context.Context, aids []int64) (as map[int64]*api.Arc, err error) {
	arg := &archive.ArgAids2{Aids: aids}
	return d.arcRpc.Archives3(ctx, arg)
}

// RanksArcs
func (d *Dao) RanksArcs(ctx context.Context, rid, pn, ps int) (res []*api.Arc, aids []int64, err error) {
	arg := &archive.ArgRank2{
		Rid: int16(rid),
		Pn:  pn,
		Ps:  ps,
	}
	var as *archive.RankArchives3
	if as, err = d.arcRpc.RankArcs3(ctx, arg); err != nil {
		log.Error("d.arcRpc.RankArcs3(%v) error(%v)", arg, err)
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

// RankTopArcs
func (d *Dao) RankTopArcs(ctx context.Context, rid, pn, ps int) (res []*api.Arc, err error) {
	arg := &archive.ArgRankTop2{
		ReID: int16(rid),
		Pn:   pn,
		Ps:   ps,
	}
	if res, err = d.arcRpc.RankTopArcs3(ctx, arg); err != nil {
		log.Error("d.arcRpc.RankTopArcs3(%v) error(%v)", arg, err)
		return
	}
	return
}
