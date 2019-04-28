package service

import (
	"context"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/log"
)

// SubtitleSubject .
func (s *Service) SubtitleSubject(c context.Context, aid int64) (subtitleSubjectReply *model.SubtitleSubjectReply, err error) {
	var (
		subtitleSubject *model.SubtitleSubject
		lan, lanDoc     string
	)
	if subtitleSubject, err = s.subtitleSubject(c, aid); err != nil {
		log.Error("params(aid:%v).err(%v)", aid, err)
		err = nil
	}
	if subtitleSubject == nil {
		subtitleSubject = &model.SubtitleSubject{
			Allow: false,
		}
	}
	lan, lanDoc = s.subtitleLans.GetByID(int64(subtitleSubject.Lan))
	subtitleSubjectReply = &model.SubtitleSubjectReply{
		AllowSubmit: subtitleSubject.Allow,
		Lan:         lan,
		LanDoc:      lanDoc,
	}
	return
}

func (s *Service) subtitleSubject(c context.Context, aid int64) (sSubject *model.SubtitleSubject, err error) {
	var (
		cacheErr bool
	)
	if sSubject, err = s.dao.SubtitleSubjectCache(c, aid); err != nil {
		cacheErr = true
		err = nil
	}
	if sSubject != nil {
		if sSubject.Empty {
			err = nil
			sSubject = nil
			return
		}
		return
	}
	if sSubject, err = s.dao.GetSubtitleSubject(c, aid); err != nil {
		log.Error("params(aid:%v).err(%v)", aid, err)
		return
	}
	if sSubject == nil {
		sSubject = &model.SubtitleSubject{
			Aid:   aid,
			Empty: true,
		}
	}
	if !cacheErr {
		temp := sSubject
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.SetSubtitleSubjectCache(ctx, temp)
		})
	}
	if sSubject.Empty {
		err = nil
		sSubject = nil
		return
	}
	return
}

// SubtitleSubjectSubmit .
func (s *Service) SubtitleSubjectSubmit(c context.Context, aid int64, allow bool, lan string) (err error) {
	var (
		lanCode         int64
		subtitleSubject *model.SubtitleSubject
	)
	lanCode = s.subtitleLans.GetByLan(lan)
	subtitleSubject = &model.SubtitleSubject{
		Aid:   aid,
		Allow: allow,
		Lan:   uint8(lanCode),
	}
	if err = s.dao.AddSubtitleSubject(c, subtitleSubject); err != nil {
		log.Error("params(subtitleSubject:%+v).err(%v)", subtitleSubject, err)
		return
	}
	if err = s.dao.DelSubtitleSubjectCache(c, aid); err != nil {
		return
	}
	return
}
