package service

import (
	"context"

	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/model"
)

//SetSwitch insert or update switch。
func (s *Service) SetSwitch(c context.Context, mid int64, state int, from uint8) (row int64, err error) {
	var res, mdlUp *model.UpSwitch
	// 查询db
	if res, err = s.up.RawUpSwitch(c, mid); err != nil {
		return
	}
	if res != nil {
		mdlUp = &model.UpSwitch{
			MID:       res.MID,
			Attribute: res.Attribute,
		}
	} else {
		mdlUp = &model.UpSwitch{MID: mid}
	}
	mdlUp.AttrSet(state, from)
	if row, err = s.up.SetSwitch(c, mdlUp); err != nil {
		return
	}
	if row > 0 {
		s.up.DelCacheUpSwitch(context.Background(), mid)
	}
	return
}

// SetUpSwitch .
func (s *Service) SetUpSwitch(c context.Context, req *upgrpc.UpSwitchReq) (res *upgrpc.NoReply, err error) {
	res = new(upgrpc.NoReply)
	var (
		row int64
		us  *model.UpSwitch
	)
	// 查询db
	if us, err = s.up.RawUpSwitch(c, req.Mid); err != nil {
		return
	}
	if us == nil {
		us = &model.UpSwitch{MID: req.Mid}
	}
	us.AttrSet(int(req.State), req.From)
	if row, err = s.up.SetSwitch(c, us); err != nil {
		return
	}
	if row > 0 {
		s.up.DelCacheUpSwitch(context.Background(), req.Mid)
	}
	return
}

// UpSwitch for app with cache.
func (s *Service) UpSwitch(c context.Context, req *upgrpc.UpSwitchReq) (res *upgrpc.UpSwitchReply, err error) {
	res = new(upgrpc.UpSwitchReply)
	var us *model.UpSwitch
	if us, err = s.up.UpSwitch(c, req.Mid); err != nil {
		return
	}
	if us == nil {
		if req.From == 0 { // 播放器开关，默认为打开
			res.State = 1
		}
		return
	}
	res.State = uint8(us.AttrVal(req.From))
	return
}

//UpSwitchs for app with cache.
func (s *Service) UpSwitchs(c context.Context, mid int64, from uint8) (state int, err error) {
	var res *model.UpSwitch
	if res, err = s.up.UpSwitch(c, mid); err != nil {
		return
	}
	if res == nil {
		if from == 0 { // 播放器开关，默认为打开
			state = 1
		}
		return
	}
	state = res.AttrVal(from)
	return
}

//RawUpSwitch for creative with no cache.
func (s *Service) RawUpSwitch(c context.Context, mid int64, from uint8) (state int, err error) {
	var res *model.UpSwitch
	if res, err = s.up.RawUpSwitch(c, mid); err != nil {
		return
	}
	if res == nil {
		return
	}
	state = res.AttrVal(from)
	return
}
