package service

import (
	"context"

	"go-common/app/admin/main/dm/model"
	"go-common/library/log"
)

// SubtitleSwitch .
func (s *Service) SubtitleSwitch(c context.Context, aid int64, allow bool, closed bool) (err error) {
	var (
		subtitleSubject *model.SubtitleSubject
		attr            = model.AttrNo
	)
	if subtitleSubject, err = s.getSubtitleSubject(c, aid); err != nil {
		log.Error("SubtitleSwitch(aid:%v) error(%v)", aid, err)
		return
	}
	if subtitleSubject == nil {
		subtitleSubject = &model.SubtitleSubject{
			Aid: aid,
		}
	}
	subtitleSubject.Allow = allow
	if closed {
		attr = model.AttrYes
	}
	subtitleSubject.AttrSet(attr, model.AttrSubtitleClose)
	if err = s.addSubtitleSubject(c, subtitleSubject); err != nil {
		log.Error("SubtitleSwitch(subtitleSubject:%+v) error(%v)", subtitleSubject, err)
		return
	}
	return
}

func (s *Service) getSubtitleSubject(c context.Context, aid int64) (subtitleSubject *model.SubtitleSubject, err error) {
	if subtitleSubject, err = s.dao.GetSubtitleSubject(c, aid); err != nil {
		return
	}
	return
}

func (s *Service) addSubtitleSubject(c context.Context, subtitleSubject *model.SubtitleSubject) (err error) {
	if err = s.dao.AddSubtitleSubject(c, subtitleSubject); err != nil {
		return
	}
	if err = s.dao.DelSubtitleSubjectCache(c, subtitleSubject.Aid); err != nil {
		return
	}
	return
}
