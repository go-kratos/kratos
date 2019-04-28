package gorpc

import (
	"context"

	"go-common/app/service/main/archive/api"
	model "go-common/app/service/main/archive/model/archive"
)

const (
	_types2         = "RPC.Types2"
	_videoshot2     = "RPC.Videoshot2"
	_upCount2       = "RPC.UpCount2"
	_upsPassed2     = "RPC.UpsPassed2"
	_upVideo2       = "RPC.UpVideo2"
	_delVideo2      = "RPC.DelVideo2"
	_description2   = "RPC.Description2"
	_ranksTopCount2 = "RPC.RanksTopCount2"
	_arcCache2      = "RPC.ArcCache2"
	_arcFieldCache2 = "RPC.ArcFieldCache2"
	_setStat2       = "RPC.SetStat2"
	_setStatCache2  = "RPC.SetStatCache2"
)

// Types2 get all archive types
func (s *Service2) Types2(c context.Context) (res map[int16]*model.ArcType, err error) {
	err = s.client.Call(c, _types2, _noArg, &res)
	return
}

// Videoshot2 get videoshot.
func (s *Service2) Videoshot2(c context.Context, arg *model.ArgCid2) (res *model.Videoshot, err error) {
	res = new(model.Videoshot)
	err = s.client.Call(c, _videoshot2, arg, res)
	return
}

// UpCount2 up count2
func (s *Service2) UpCount2(c context.Context, arg *model.ArgUpCount2) (count int, err error) {
	err = s.client.Call(c, _upCount2, arg, &count)
	return
}

// UpsPassed2 get UpsPassed aid and ptime
func (s *Service2) UpsPassed2(c context.Context, arg *model.ArgUpsArcs2) (res map[int64][]*model.AidPubTime, err error) {
	err = s.client.Call(c, _upsPassed2, arg, &res)
	return
}

// UpVideo2 update video cache by aid & cid
func (s *Service2) UpVideo2(c context.Context, arg *model.ArgVideo2) (err error) {
	err = s.client.Call(c, _upVideo2, arg, _noArg)
	return
}

// DelVideo2 delete video cache by aid & cid
func (s *Service2) DelVideo2(c context.Context, arg *model.ArgVideo2) (err error) {
	err = s.client.Call(c, _delVideo2, arg, _noArg)
	return
}

// Description2 add share.
func (s *Service2) Description2(c context.Context, arg *model.ArgAid) (des string, err error) {
	err = s.client.Call(c, _description2, arg, &des)
	return
}

// RanksTopCount2 get top region count.
func (s *Service2) RanksTopCount2(c context.Context, arg *model.ArgRankTopsCount2) (res map[int16]int, err error) {
	err = s.client.Call(c, _ranksTopCount2, arg, &res)
	return
}

// ArcCache2 add/update archive cache
func (s *Service2) ArcCache2(c context.Context, arg *model.ArgCache2) (err error) {
	err = s.client.Call(c, _arcCache2, arg, _noArg)
	return
}

// ArcFieldCache2  update archive field cache
func (s *Service2) ArcFieldCache2(c context.Context, arg *model.ArgFieldCache2) (err error) {
	err = s.client.Call(c, _arcFieldCache2, arg, _noArg)
	return
}

// SetStat2 set all stat info.
func (s *Service2) SetStat2(c context.Context, arg *api.Stat) (err error) {
	err = s.client.Call(c, _setStat2, arg, _noArg)
	return
}

// SetStatCache2 up stat.
func (s *Service2) SetStatCache2(c context.Context, arg *model.ArgStat2) (err error) {
	err = s.client.Call(c, _setStatCache2, arg, _noArg)
	return
}
