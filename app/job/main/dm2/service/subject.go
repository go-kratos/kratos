package service

import (
	"context"

	"go-common/app/job/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) subject(c context.Context, tp int32, oid int64) (sub *model.Subject, err error) {
	var cache = true
	if sub, err = s.dao.SubjectCache(c, tp, oid); err != nil {
		err = nil
		cache = false
	}
	if sub == nil {
		if sub, err = s.dao.Subject(c, tp, oid); err != nil {
			return
		}
		if sub == nil {
			sub = &model.Subject{
				Type: tp,
				Oid:  oid,
			}
		}
		if cache {
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.AddSubjectCache(ctx, sub)
			})
		}
	}
	if sub.ID == 0 {
		err = ecode.NothingFound
		log.Error("subject not exist,type:%d,oid:%d", tp, oid)
		return
	}
	return
}
