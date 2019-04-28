package bnj

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"go-common/library/sync/errgroup.v2"

	"go-common/app/interface/main/activity/model/bnj"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_rewardStepOne   = 1
	_rewardStepThree = 3
	_lastCD          = 300
)

// PreviewInfo preview info
func (s *Service) PreviewInfo(c context.Context, mid int64) *bnj.PreviewInfo {
	data := &bnj.PreviewInfo{
		ActID: s.c.Bnj2019.ActID,
		SubID: s.c.Bnj2019.SubID,
	}
	// TODO del admin check
	if err := s.bnjAdminCheck(mid); err != nil {
		return data
	}
	for _, v := range s.c.Bnj2019.Reward {
		data.RewardStep = append(data.RewardStep, v.Condition)
	}
	if s.timeFinish != 0 {
		data.TimelinePic = s.c.Bnj2019.TimelinePic
		data.H5TimelinePic = s.c.Bnj2019.H5TimelinePic
	}
	if s.c.Bnj2019.GameCancel != 0 {
		data.GameCancel = 1
	}
	now := time.Now().Unix()
	group := errgroup.WithCancel(c)
	if mid > 0 && len(s.c.Bnj2019.Reward) > 0 {
		for _, v := range s.c.Bnj2019.Reward {
			if v.Step > 0 {
				step := v.Step
				group.Go(func(ctx context.Context) error {
					if check, e := s.dao.HasReward(ctx, mid, s.c.Bnj2019.SubID, step); e != nil {
						log.Error("Reward s.dao.HasReward(mid:%d,step:%d) error(%v) check(%v)", mid, step, e, check)
					} else if check {
						switch step {
						case _rewardStepOne:
							data.HasRewardFirst = 1
						case _rewardStepThree:
							data.HasRewardSecond = 1
						}
					}
					return nil
				})
			}
		}
	}
	if err := group.Wait(); err != nil {
		log.Error("PreviewInfo group wait error(%v)", err)
	}
	arcs := s.previewArcs
	for _, v := range s.c.Bnj2019.Info {
		if v.Publish.Unix() < now {
			if arc, ok := arcs[v.Aid]; ok && arc.IsNormal() {
				tmp := &bnj.Info{Nav: v.Nav, Pic: v.Pic, H5Pic: v.H5Pic, Detail: v.Detail, H5Detail: v.H5Detail, Arc: &arcmdl.Arc{Aid: v.Aid}}
				arc.Pic = strings.Replace(arc.Pic, "http://", "//", 1)
				if v.Nickname != "" {
					arc.Author.Name = arc.Author.Name + "&" + v.Nickname
				}
				tmp.Arc = arc
				data.Info = append(data.Info, tmp)
			}
		}
	}
	if len(data.Info) == 0 {
		data.Info = make([]*bnj.Info, 0)
	}
	return data
}

// Timeline only return timeline and game cancel.
func (s *Service) Timeline(c context.Context, mid int64) *bnj.Timeline {
	data := new(bnj.Timeline)
	// TODO delete admin check
	if err := s.bnjAdminCheck(mid); err != nil {
		return data
	}
	if s.timeFinish != 0 {
		data.TimelinePic = s.c.Bnj2019.TimelinePic
		data.H5TimelinePic = s.c.Bnj2019.H5TimelinePic
	}
	if s.c.Bnj2019.GameCancel != 0 {
		data.GameCancel = 1
	}
	data.LikeCount = s.likeCount
	return data
}

// TimeReset reset less time.
func (s *Service) TimeReset(c context.Context, mid int64) (ttl int64, err error) {
	if time.Now().Unix() < s.c.Bnj2019.Start.Unix() {
		err = ecode.ActivityNotStart
		return
	}
	if s.timeFinish != 0 {
		err = ecode.ActivityBnjTimeFinish
		return
	}
	var value bool
	if value, err = s.dao.CacheResetCD(c, mid, s.resetCD); err != nil {
		log.Error("TimeReset s.dao.CacheResetCD(%d) error(%v) value(%v)", mid, err, value)
		err = nil
		return
	}
	if !value {
		if ttl, err = s.dao.TTLResetCD(c, mid); err != nil {
			log.Error("TimeReset s.dao.TTLResetCD(%d) error(%v) value(%v)", mid, err)
			err = nil
		}
		return
	}
	if s.timeReset == 0 {
		atomic.StoreInt64(&s.resetMid, mid)
		atomic.StoreInt64(&s.timeReset, 1)
	}
	return
}

// DelTime .
func (s *Service) DelTime(c context.Context, key string) (err error) {
	switch key {
	case "time_finish":
		if err = s.dao.DelCacheTimeFinish(c); err != nil {
			log.Error("DelTime DelCacheTimeFinish error(%v)", err)
		}
	case "time_less":
		if err = s.dao.DelCacheTimeLess(c); err != nil {
			log.Error("DelTime DelCacheTimeLess error(%v)", err)
		}
	}
	return
}

func (s *Service) bnjAdminCheck(mid int64) (err error) {
	if s.c.Bnj2019.AdminCheck != 0 {
		if _, ok := s.bnjAdmins[mid]; !ok {
			err = ecode.AccessDenied
		}
	}
	return
}
