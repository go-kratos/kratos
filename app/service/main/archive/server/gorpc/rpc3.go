package server

import (
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/net/rpc/context"
)

// 3结尾的方法全都是pb格式的memcache

// MaxAID get max aid
func (r *RPC) MaxAID(c context.Context, a *struct{}, res *int64) (err error) {
	*res, err = r.s.MaxAID(c)
	return
}

// Archive3 receive aid, then init archive info.
func (r *RPC) Archive3(c context.Context, a *archive.ArgAid2, res *api.Arc) (err error) {
	var ar *api.Arc
	if ar, err = r.s.Archive3(c, a.Aid); err == nil {
		*res = *ar
	}
	return
}

// Archives3 receive aids, then init archives info.
func (r *RPC) Archives3(c context.Context, a *archive.ArgAids2, res *map[int64]*api.Arc) (err error) {
	if len(a.Aids) > 300 {
		log.Error("Too many Args aids(%d) caller(%s) arg(%v)", len(a.Aids), c.User(), a.Aids)
	}
	*res, err = r.s.Archives3(c, a.Aids)
	return
}

// View3 view archive.
func (r *RPC) View3(c context.Context, a *archive.ArgAid2, av *archive.View3) (err error) {
	var res *api.ViewReply
	if res, err = r.s.View3(c, a.Aid); err == nil {
		*av = *archive.BuildView3(res.Arc, res.Pages)
	}
	return
}

// Views3 view archive.
func (r *RPC) Views3(c context.Context, a *archive.ArgAids2, res *map[int64]*archive.View3) (err error) {
	if len(a.Aids) > 300 {
		log.Error("Too many Args aids(%d) caller(%s)", len(a.Aids), c.User())
	}
	var views map[int64]*api.ViewReply
	if views, err = r.s.Views3(c, a.Aids); err == nil {
		var resp = make(map[int64]*archive.View3)
		for aid, view := range views {
			v := archive.BuildView3(view.Arc, view.Pages)
			if v != nil {
				resp[aid] = v
			}
		}
		*res = resp
	}
	return
}

// Stat3 archive stat.
func (r *RPC) Stat3(c context.Context, a *archive.ArgAid2, res *api.Stat) (err error) {
	var st *api.Stat
	if st, err = r.s.Stat3(c, a.Aid); err == nil {
		*res = *st
	}
	return
}

// Page3 get videos by aid
func (r *RPC) Page3(c context.Context, a *archive.ArgAid2, res *[]*api.Page) (err error) {
	*res, err = r.s.Page3(c, a.Aid)
	return
}

// Stats3 archive stats.
func (r *RPC) Stats3(c context.Context, a *archive.ArgAids2, res *map[int64]*api.Stat) (err error) {
	if len(a.Aids) > 200 {
		log.Error("Too many Args aids(%d) caller(%s)", len(a.Aids), c.User())
		log.Error("Too many Args aids(%d) caller(%s) arg(%v)", len(a.Aids), c.User(), a.Aids)
	}
	*res, err = r.s.Stats3(c, a.Aids)
	return
}

// Click3 archive click.
func (r *RPC) Click3(c context.Context, a *archive.ArgAid2, res *api.Click) (err error) {
	var clk *api.Click
	if clk, err = r.s.Click3(c, a.Aid); err == nil {
		*res = *clk
	}
	return
}

// UpArcs3 up archives.
func (r *RPC) UpArcs3(c context.Context, a *archive.ArgUpArcs2, res *[]*api.Arc) (err error) {
	*res, err = r.s.UpperPassed3(c, a.Mid, a.Pn, a.Ps)
	return
}

// UpsArcs3 ups archives.
func (r *RPC) UpsArcs3(c context.Context, a *archive.ArgUpsArcs2, res *map[int64][]*api.Arc) (err error) {
	*res, err = r.s.UppersPassed3(c, a.Mids, a.Pn, a.Ps)
	return
}

// Recommend3 from archive_recommend by aid
func (r *RPC) Recommend3(c context.Context, a *archive.ArgAid2, res *[]*api.Arc) (err error) {
	*res, err = r.s.UpperReommend(c, a.Aid)
	return
}

// RankArcs3 Arcs by rid
func (r *RPC) RankArcs3(c context.Context, a *archive.ArgRank2, res *archive.RankArchives3) (err error) {
	if a.Type == 0 {
		res.Archives, res.Count, err = r.s.RegionArcs3(c, a.Rid, a.Pn, a.Ps)
	} else {
		res.Archives, res.Count, err = r.s.RegionOriginArcs3(c, a.Rid, a.Pn, a.Ps)
	}
	return
}

// ArchivesWithPlayer archives with player info
func (r *RPC) ArchivesWithPlayer(c context.Context, a *archive.ArgPlayer, res *map[int64]*archive.ArchiveWithPlayer) (err error) {
	var as map[int64]*archive.ArchiveWithPlayer
	if as, err = r.s.ArchivesWithPlayer(c, a, false); err != nil {
		return
	}
	*res = as
	return
}

// RanksArcs3 Arcs by rids
func (r *RPC) RanksArcs3(c context.Context, a *archive.ArgRanks2, res *map[int16]*archive.RankArchives3) (err error) {
	var (
		as    []*api.Arc
		tmp   map[int16]*archive.RankArchives3
		count int
	)
	tmp = make(map[int16]*archive.RankArchives3)
	for _, rid := range a.Rids {
		if a.Type == 0 {
			as, count, err = r.s.RegionArcs3(c, rid, a.Pn, a.Ps)
		} else {
			as, count, err = r.s.RegionOriginArcs3(c, rid, a.Pn, a.Ps)
		}
		if err != nil {
			return
		}
		tmp[rid] = &archive.RankArchives3{Archives: as, Count: count}
	}
	*res = tmp
	return
}

// RankTopArcs3 Arcs by reids
func (r *RPC) RankTopArcs3(c context.Context, a *archive.ArgRankTop2, res *[]*api.Arc) (err error) {
	*res, err = r.s.RegionTopArcs3(c, a.ReID, a.Pn, a.Ps)
	return
}

// RankAllArcs3 left 7 days all Arcs
func (r *RPC) RankAllArcs3(c context.Context, a *archive.ArgRankAll2, res *archive.RankArchives3) (err error) {
	var data *archive.RankArchives3
	data, err = r.s.RegionAllArcs3(c, a.Pn, a.Ps)
	if err == nil {
		*res = *data
	}
	return
}

// Video3 get video by aid & cid.
func (r *RPC) Video3(c context.Context, a *archive.ArgVideo2, res *api.Page) (err error) {
	var p *api.Page
	p, err = r.s.Video3(c, a.Aid, a.Cid)
	if err == nil {
		*res = *p
	}
	return
}
