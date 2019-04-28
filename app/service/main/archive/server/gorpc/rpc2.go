package server

import (
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/net/rpc/context"
)

// Types2 get all types
func (r *RPC) Types2(c context.Context, a *struct{}, res *map[int16]*archive.ArcType) (err error) {
	*res = r.s.AllTypes(c)
	return
}

// Videoshot2 get vidoshot info data.
func (r *RPC) Videoshot2(c context.Context, a *archive.ArgCid2, res *archive.Videoshot) (err error) {
	var (
		v *archive.Videoshot
	)
	if v, err = r.s.Videoshot(c, a.Aid, a.Cid); err == nil {
		*res = *v
	}
	return
}

// UpCount2 up count.
func (r *RPC) UpCount2(c context.Context, a *archive.ArgUpCount2, count *int) (err error) {
	*count, err = r.s.UpperCount(c, a.Mid)
	return
}

// UpsPassed2 ups pass aid and ptime
func (r *RPC) UpsPassed2(c context.Context, a *archive.ArgUpsArcs2, res *map[int64][]*archive.AidPubTime) (err error) {
	*res, err = r.s.UppersAidPubTime(c, a.Mids, a.Pn, a.Ps)
	return
}

// UpVideo2 up video by aid & cid.
func (r *RPC) UpVideo2(c context.Context, a *archive.ArgVideo2, res *struct{}) (err error) {
	return r.s.UpVideo(c, a.Aid, a.Cid)
}

// DelVideo2 delete video by aid & cid.
func (r *RPC) DelVideo2(c context.Context, a *archive.ArgVideo2, res *struct{}) (err error) {
	return r.s.DelVideo(c, a.Aid, a.Cid)
}

// Description2 get description by aid.
func (r *RPC) Description2(c context.Context, a *archive.ArgAid, reDes *string) (err error) {
	*reDes, err = r.s.Description(c, a.Aid)
	return
}

// RanksTopCount2 top region count.
func (r *RPC) RanksTopCount2(c context.Context, a *archive.ArgRankTopsCount2, res *map[int16]int) (err error) {
	*res, err = r.s.RegionTopCount(c, a.ReIDs)
	return
}

// ArcsStat2 archive stat.
// func (r *RPC) ArcsStat2(c context.Context, a *archive.ArgAids2, res *map[int64]*archive.Stat) (err error) {
// 	*res, err = r.s.Stats(c, a.Aids)
// 	return
// }

// ArcCache2 update archive cache.
func (r *RPC) ArcCache2(c context.Context, a *archive.ArgCache2, res *struct{}) (err error) {
	err = r.s.CacheUpdate(c, a.Aid, a.Tp, a.OldMid)
	return
}

// ArcFieldCache2 update archive by field changed.
func (r *RPC) ArcFieldCache2(c context.Context, a *archive.ArgFieldCache2, res *struct{}) (err error) {
	err = r.s.FieldCacheUpdate(c, a.Aid, a.OldTypeID, a.TypeID)
	return
}

// SetStat2 set all stat cache(redis)
func (r *RPC) SetStat2(c context.Context, a *api.Stat, res *struct{}) (err error) {
	err = r.s.SetStat(c, a)
	return
}
