package service

import (
	"context"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
)

// UpdateMask update mask list
func (s *Service) UpdateMask(c context.Context, cid, masktime int64, fps int32, plat int8, list string) (err error) {
	var sub *model.Subject
	if sub, err = s.subject(c, model.SubTypeVideo, cid); err != nil {
		return
	}
	if sub == nil {
		err = ecode.ArchiveNotExist
		return
	}
	if err = s.dao.UpdateMask(c, cid, masktime, fps, plat, list); err != nil {
		return
	}
	if plat == model.MaskPlatMbl {
		sub.AttrSet(model.AttrYes, model.AttrSubMblMaskReady)
	} else {
		sub.AttrSet(model.AttrYes, model.AttrSubWebMaskReady)
	}
	if _, err = s.dao.UptSubAttr(c, model.SubTypeVideo, cid, sub.Attr); err != nil {
		return
	}
	tmp := *sub
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.AddSubjectCache(ctx, &tmp)
	})
	mask, err := s.dao.MaskList(c, cid, plat)
	if err != nil || mask == nil {
		return
	}
	maskTmp := *mask
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.AddMaskCache(ctx, model.SubTypeVideo, &maskTmp)
	})
	return
}

// MaskListWithSub .
func (s *Service) MaskListWithSub(c context.Context, cid int64, plat int8, sub *model.Subject) (mask *model.Mask, err error) {
	var ok bool
	if plat == model.MaskPlatWeb {
		if sub.AttrVal(model.AttrSubMaskOpen) == model.AttrYes && sub.AttrVal(model.AttrSubWebMaskReady) == model.AttrYes {
			ok = true
		}
	} else {
		if sub.AttrVal(model.AttrSubMaskOpen) == model.AttrYes && sub.AttrVal(model.AttrSubMblMaskReady) == model.AttrYes {
			ok = true
		}
	}
	if !ok {
		return
	}
	if mask, err = s.dao.DMMaskCache(c, model.SubTypeVideo, cid, plat); err != nil {
		err = nil
		ok = false
	}
	if mask == nil {
		if mask, err = s.dao.MaskList(c, cid, plat); err != nil || mask == nil {
			return
		}
		if ok {
			tmp := *mask
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.AddMaskCache(ctx, model.SubTypeVideo, &tmp)
			})
		}
	}
	return
}

// MaskList get mask info
func (s *Service) MaskList(c context.Context, cid int64, plat int8) (mask *model.Mask, err error) {
	var sub *model.Subject
	if sub, err = s.subject(c, model.SubTypeVideo, cid); err != nil || sub == nil {
		return
	}
	return s.MaskListWithSub(c, cid, plat, sub)
}
