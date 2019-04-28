package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/tag/model"
	rpcModel "go-common/app/service/main/tag/model"
	"go-common/library/ecode"
)

func (s *Service) checkTypeRole(umid, amid int64) (int8, int8, int8) {
	if umid == amid {
		return model.UpTag, model.RoleUp, model.ArcTagOpRoleUp
	}
	return model.UserTag, model.RoleUser, model.ArcTagOpRoleUser
}

// check archive tags
func (s *Service) checkArcTag(c context.Context, aid, tid int64, op int8) (atm map[int64]*model.ArcTag, err error) {
	var ats []*model.ArcTag
	if ats, err = s.arcTagsService(c, aid, 0, rpcModel.ResTypeArchive); err != nil {
		return
	}
	if op == model.ArcTagAdd && len(ats) >= s.c.Tag.ArcTagMaxNum {
		err = ecode.TagArcTagMaxNum
		return
	}
	atm = make(map[int64]*model.ArcTag)
	for _, at := range ats {
		atm[at.Tid] = at
	}
	if len(atm) > 0 {
		_, ok := atm[tid]
		if op == model.ArcTagAdd && ok {
			err = ecode.TagArcTagExist
		} else if op == model.ArcTagDel && !ok {
			err = ecode.TagArcTagNotExist
		}
	} else {
		if op == model.ArcTagDel {
			err = ecode.TagArcTagNotExist
		}
	}
	return
}

// check user level and some operator number
func (s *Service) checkUser(c context.Context, aid, authorMid, mid, tid int64, now time.Time, op int8) (err error) {
	// TODO 同一个用户在同一个视频下 删除tag的个数限制
	// 后台管理员操作判断
	switch op {
	case model.ArcTagAdd:
		var rem map[string]*rpcModel.ResourceLog
		rem, err = s.resTagLogMap(c, mid, aid, rpcModel.ResTypeArchive, 1, 20)
		if err != nil {
			return
		}
		if len(rem) > 0 {
			k := fmt.Sprintf("%d_%d_%d_%d_%d", aid, rpcModel.ResTypeArchive, tid, mid, rpcModel.ResTagLogAdd)
			if rl, ok := rem[k]; ok && rl.State == 1 {
				return ecode.TagAddNotRptPassed
			}
		}
	case model.ArcTagDel:
		if authorMid != mid {
			var rem map[string]*rpcModel.ResourceLog
			rem, err = s.resTagLogMap(c, mid, aid, rpcModel.ResTypeArchive, 1, 20)
			if err != nil {
				return
			}
			if len(rem) > 0 {
				k := fmt.Sprintf("%d_%d_%d_%d_%d", aid, rpcModel.ResTypeArchive, tid, mid, rpcModel.ResTagLogDel)
				if rl, ok := rem[k]; ok && rl.State == 1 {
					return ecode.TagDelNotRptPassed
				}
			}
		}
	}
	// if commom user
	if _, isWhite := s.whiteUser[mid]; !isWhite {
		// limit archive
		if attr, ok := s.limitArc[aid]; ok {
			if op == model.ArcTagDel {
				if attr&0x2 == 0x2 {
					err = ecode.TagArcCannotDelTag
					return
				}
			} else if op == model.ArcTagAdd {
				if attr&0x1 == 0x1 {
					err = ecode.TagArcCannotAddTag
					return
				}
			}
		}
		if op == model.ArcTagDel {
			var rtm map[int64]*rpcModel.Resource
			if rtm, err = s.resTagMap(c, 0, aid, rpcModel.ResTypeArchive); err != nil {
				return
			}
			rt, ok := rtm[tid]
			if !ok {
				return ecode.TagResTagNotExist
			}
			var tag *model.Tag
			if tag, err = s.info(c, 0, tid); err != nil {
				return
			}
			if tag.Type == model.OfficailActiveTag {
				return ecode.TagIsOfficailTag
			}
			if rt.Role == model.RoleUp && mid != authorMid {
				return ecode.TagUpTagCannotDel
			}
			if rt.Attr&0x1 == 1 {
				return ecode.TagArcTagisLocked
			}
		}
	}
	return
}

// func (s *Service) getCheckInfo2(op int8) (checkNumber int, opErr error) {
// 	switch op {
// 	case model.ArcTagDel:
// 		checkNumber = s.c.Tag.ArcTagDelSomeNum
// 		opErr = ecode.TagArcTagDelMore
// 	}
// 	return
// }

func (s *Service) policy(c context.Context, realMid, mid int64, now time.Time) (err error) {
	// 6 4政策check 十九大
	if s.c.Supervision.SixFour != nil && s.c.Supervision.SixFour.Button {
		if now.After(s.c.Supervision.SixFour.Begin) && now.Before(s.c.Supervision.SixFour.End) {
			if realMid != mid {
				return ecode.TagServiceUpdate
			}
		}
	}
	return
}

func (s *Service) realName(c context.Context, mid int64) (err error) {
	if !s.c.Supervision.RealName.Button {
		return
	}
	profile, err := s.dao.UserProfile(c, mid)
	if err != nil {
		return
	}
	if profile.Identification == 0 && profile.TelStatus == 0 {
		err = ecode.UserCheckNoPhone
		return
	}
	if profile.Identification == 0 && profile.TelStatus == 2 {
		err = ecode.UserCheckInvalidPhone
	}
	if profile.Silence != model.UserBannedNone {
		return ecode.TagArcAccountBlocked
	}
	return
}
