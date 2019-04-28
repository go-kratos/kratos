package up

import (
	"context"
	"go-common/app/interface/main/creative/model/up"
	upmdl "go-common/app/service/main/up/model"
	"go-common/library/log"
)

const (
	StaffWhiteGroupID = 24
)

// ArcUpInfo for main app submit.
func (s *Service) ArcUpInfo(c context.Context, mid int64, ip string) (isAuthor int32, err error) {
	res, err := s.up.UpInfo(c, mid, 1, ip)
	if err != nil {
		log.Error("s.acc.ArcUpInfo(%d) error(%v)", mid, err)
		return
	}
	isAuthor = res.IsAuthor
	return
}

// UpSwitch get switch.
func (s *Service) UpSwitch(c context.Context, mid int64, ip string) (res *up.Switch, err error) {
	var ups *upmdl.PBUpSwitch
	ups, err = s.up.UpSwitch(c, mid, 0, ip)
	if err != nil {
		log.Error("s.up.UpSwitch mid(%d)|from(%d)|error(%v)", mid, err)
	}
	if ups == nil {
		return
	}
	pf, err := s.acc.Profile(c, mid, ip)
	if err != nil {
		log.Error("s.acc.Profile mid(%d)|error(%v)", mid, err)
		return
	}
	show := 0
	res = &up.Switch{
		State: ups.State,
		Show:  show,
		Face:  pf.Face,
	}
	return
}

// SetUpSwitch set switch.
func (s *Service) SetUpSwitch(c context.Context, mid int64, state, from int, ip string) (res *upmdl.PBSetUpSwitchRes, err error) {
	res, err = s.up.SetUpSwitch(c, mid, state, from, ip)
	if err != nil {
		log.Error("s.up.SetUpSwitch mid(%d)|state(%d)|from(%d)|error(%v)", mid, state, from, err)
	}
	return
}

// ShowStaff 用户是否能看到联合投稿
func (s *Service) ShowStaff(c context.Context, mid int64) (show bool, err error) {
	//如果关了灰度，则展示
	if !s.c.StaffConf.IsGray {
		show = true
		return
	}
	var (
		groups map[int64]*up.SpecialGroup
	)
	groups = make(map[int64]*up.SpecialGroup)
	if groups, err = s.up.UpSpecialGroups(c, mid); err != nil {
		log.Error("s.up.UpSpecialGroups(%d) error(%v)", mid, err)
		return
	}
	if _, ok := groups[StaffWhiteGroupID]; ok {
		show = true
		return
	}
	return
}
