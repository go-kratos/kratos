package service

import (
	"context"
	"time"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
)

const (
	// identify state
	// 0: normal phone number 1: virtual phone number 2: no phone number
	_identifyPhoneVirtual = 1 // virtual phone number
	_identifyPhoneEmpty   = 2 // no phone number
)

// ApplyInfo get apply info
// 检查顺序: 是否已通过->拒绝-> 已提交-> 开放申请-> 申请已满 -> 实名认证/封禁状态/绑定手机
func (s *Service) ApplyInfo(c context.Context, mid int64) (res *model.Apply, err error) {
	if mid == 0 {
		if !s.setting.ApplyOpen {
			err = ecode.ArtApplyClose
			return
		}
		if s.checkApplyFull(c) {
			err = ecode.ArtApplyFull
			return
		}
		return
	}
	res = &model.Apply{}
	res.Forbid, _, _ = s.UserDisabled(c, mid)
	if res.Forbid {
		err = ecode.ArtApplyForbid
		return
	}
	if pass, _, _ := s.IsAuthor(c, mid); pass {
		err = ecode.ArtApplyPass
		return
	}
	var author *model.AuthorLimit
	if author, err = s.dao.RawAuthor(c, mid); err != nil {
		return
	}
	if author != nil {
		if author.State == model.AuthorStatePass {
			err = ecode.ArtApplyPass
			return
		} else if author.State == model.AuthorStateReject {
			if time.Now().Unix()-int64(author.Rtime) <= s.setting.ApplyFrozenDuration {
				err = ecode.ArtApplyReject
				return
			}
		} else if author.State == model.AuthorStatePending {
			err = ecode.ArtApplySubmit
			return
		}
	}
	if !s.setting.ApplyOpen {
		err = ecode.ArtApplyClose
		return
	}
	if s.checkApplyFull(c) {
		err = ecode.ArtApplyFull
		return
	}
	var identify *model.Identify
	if identify, err = s.dao.Identify(c, mid); err != nil {
		return
	}
	res.Verify = (identify.Identify == 0)
	res.Phone = identify.Phone
	if res.Phone == _identifyPhoneEmpty {
		err = ecode.ArtApplyPhone
	}
	return
}

func (s *Service) checkApplyFull(c context.Context) (full bool) {
	if count, err := s.dao.ApplyCount(c); err != nil {
		return
	} else if count > 0 {
		return count >= s.setting.ApplyLimit
	}
	return
}

// Apply add apply
func (s *Service) Apply(c context.Context, mid int64, content, category string) (err error) {
	var res = &model.Apply{}
	if res, err = s.ApplyInfo(c, mid); err != nil {
		return
	} else if res.Phone == _identifyPhoneVirtual {
		err = ecode.ArtApplyPhoneVirtual
		return
	}
	err = s.dao.AddApply(c, mid, content, category)
	return
}
