package service

import (
	"context"
	"encoding/json"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) subject(c context.Context, tp int32, oid int64) (sub *model.Subject, err error) {
	var (
		cache = true
		bs    []byte
		ok    bool
	)
	if bs, ok = s.localCache[keySubject(tp, oid)]; ok {
		sub = &model.Subject{}
		if err = json.Unmarshal(bs, &sub); err == nil {
			return
		}
	}
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
		return
	}
	return
}

func (s *Service) subjects(c context.Context, tp int32, oids []int64) (res map[int64]*model.Subject, err error) {
	var (
		cache       = true
		missed      []int64
		missedCache map[int64]*model.Subject
		hitedCache  map[int64]*model.Subject
	)
	res = make(map[int64]*model.Subject, len(oids))
	if hitedCache, missed, err = s.dao.SubjectsCache(c, tp, oids); err != nil {
		cache = false
	}
	if len(hitedCache) == 0 {
		missed = oids
	}
	if len(missed) > 0 {
		if missedCache, err = s.dao.Subjects(c, tp, missed); err != nil {
			return
		}
		for _, oid := range missed {
			sub, ok := missedCache[oid]
			if ok {
				res[sub.Oid] = sub
			} else {
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
	}
	for _, hit := range hitedCache {
		if hit.ID > 0 {
			res[hit.Oid] = hit
		}
	}
	return
}

// SubjectInfos get dm subject info by oids.
func (s *Service) SubjectInfos(c context.Context, tp int32, plat int8, oids []int64) (res map[int64]*model.SubjectInfo, err error) {
	subs, err := s.subjects(c, tp, oids)
	if err != nil {
		log.Error("s.subjects(%v) error(%v)", oids, err)
		return
	}
	res = make(map[int64]*model.SubjectInfo, len(oids))
	for _, sub := range subs {
		subInfo := new(model.SubjectInfo)
		if sub.Count > sub.Maxlimit {
			subInfo.Count = sub.ACount
		} else {
			subInfo.Count = sub.Count
		}
		if s.isRealname(c, sub.Pid, sub.Oid) {
			subInfo.Realname = true
		}
		if sub.State == model.SubStateClosed {
			subInfo.Closed = true
		}
		res[sub.Oid] = subInfo
	}
	return
}
