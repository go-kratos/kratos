package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/growup/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_canEnrol = iota // 未报名,可以报名
	_hasEnrol        // 已报名
	_win             // 中奖
	_notEnrol        // 不能报名
	_enrolEnd        // 报名已结束
)

// ShowActivity show creative_activity
func (s *Service) ShowActivity(c context.Context, mid, activityID int64) (ac *model.CActivity, err error) {
	ac, err = s.dao.GetActivity(c, activityID)
	if err != nil {
		log.Error("s.dao.GetActivity error(%v)", err)
		return
	}
	enrolNum, winNum, mids, upState, err := s.handleUpActivity(c, mid, activityID)
	if err != nil {
		log.Error("s.handleUpActivity error(%v)", err)
		return
	}
	// 当前时间是否可以展示
	now := xtime.Time(time.Now().Unix())
	if ac.UpdatePage == 1 && now >= ac.ProgressStart && now <= ac.ProgressEnd {
		ac.ProgressState = 1
		ac.Enrollment = enrolNum
		ac.WinNum = winNum
		// 2: 中奖类型为排序型
		if ac.WinType == 2 {
			ac.Ranking, err = s.getActUpInfo(c, mids)
			if err != nil {
				log.Error("s.getActUpInfo error(%v)", err)
				return
			}
		}
	}
	if !(ac.BonusQuery == 1 && now >= ac.BonusQueryStart && now <= ac.BonusQueryEnd) {
		ac.BonusQuery = 0
	}
	// 获取up主当前状态
	ac.SignUpState, err = s.getSignUpState(c, mid, upState, ac)
	if err != nil {
		log.Error("s.getSignUpState error(%v)", err)
	}
	return
}

// get mid name and face
func (s *Service) getActUpInfo(c context.Context, mids []int64) (upInfos []*model.ActUpInfo, err error) {
	upInfoMap, err := s.dao.AccountInfos(c, mids)
	if err != nil {
		return
	}
	upInfos = make([]*model.ActUpInfo, len(upInfoMap))
	for i := 0; i < len(mids); i++ {
		upInfos[i] = upInfoMap[mids[i]]
	}
	return
}

func (s *Service) getSignUpState(c context.Context, mid int64, upState int, ac *model.CActivity) (signUpState int, err error) {
	now := xtime.Time(time.Now().Unix())
	// 报名未开始,不能报名
	if now < ac.SignUpStart {
		signUpState = _notEnrol
		return
	}
	// 报名已结束并且未中奖
	if now > ac.SignUpEnd && upState != _win {
		upState = _enrolEnd
	}
	// 已报名
	if upState >= _hasEnrol {
		signUpState = upState
		return
	}

	// 签约结束时间 >= 报名结束时间, 任何人都可以报名
	if ac.SignedEnd >= ac.SignUpEnd {
		signUpState = _canEnrol
	} else {
		var signedAt xtime.Time
		signedAt, err = s.dao.GetUpSignedAt(c, "up_info_video", mid)
		if err != nil {
			return
		}
		if signedAt >= ac.SignedStart && signedAt <= ac.SignedEnd {
			signUpState = _canEnrol
		} else {
			signUpState = _notEnrol
		}
	}
	return
}

func (s *Service) handleUpActivity(c context.Context, mid, activityID int64) (enrol, win int, mids []int64, upState int, err error) {
	ups, err := s.dao.ListUpActivity(c, activityID)
	if err != nil {
		return
	}
	sort.Slice(ups, func(i, j int) bool {
		return ups[i].Rank < ups[j].Rank
	})
	mids = make([]int64, 0)
	for _, up := range ups {
		if up.State >= _win {
			win++
			enrol++
			mids = append(mids, up.MID)
		} else if up.State == _hasEnrol {
			enrol++
		}
		if up.MID == mid {
			upState = up.State
			// 已发奖
			if upState == 3 {
				upState = 2
			}
		}
	}
	return
}

// SignUpActivity up sign up activity
func (s *Service) SignUpActivity(c context.Context, mid, activityID int64) (err error) {
	nickname, _, err := s.dao.CategoryInfo(c, mid)
	if err != nil {
		log.Error("s.dao.CategoryInfo error(%v)", err)
		return
	}
	upBonus := &model.UpBonus{
		MID:        mid,
		ActivityID: activityID,
		Nickname:   nickname,
		State:      1,
		SignUpTime: xtime.Time(time.Now().Unix()),
	}
	if _, err = s.dao.SignUpActivity(c, upBonus); err != nil {
		log.Error("s.dao.SignUpActivity error(%v)", err)
	}
	return
}
