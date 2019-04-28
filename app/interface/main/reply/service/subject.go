package service

import (
	"context"

	model "go-common/app/interface/main/reply/model/reply"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) subject(c context.Context, oid int64, tp int8) (sub *model.Subject, err error) {
	if sub, err = s.dao.Mc.GetSubject(c, oid, tp); err != nil {
		log.Error("replyCacheDao.GetSubject(%d,%d) error(%v)", oid, tp, err)
		return
	} else if sub == nil {
		if sub, err = s.dao.Subject.Get(c, oid, tp); err != nil {
			log.Error("s.subject.Get(%d,%d) error(%v)", oid, tp, err)
			return
		}
		if sub == nil {
			sub = &model.Subject{ID: -1, Oid: oid, Type: tp, State: model.SubStateForbid} // empty cache
		}
		s.cache.Do(c, func(ctx context.Context) {
			if err = s.dao.Mc.AddSubject(ctx, sub); err != nil {
				log.Error("s.dao.Mc.AddSubject(%d,%d) error(%v)", oid, tp, err)
			}
		})
	}
	return
}

// ReplyCount get all reply count.
func (s *Service) ReplyCount(c context.Context, oid int64, tp int8) (count int, err error) {
	if !model.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	sub, err := s.subject(c, oid, tp)
	if err != nil {
		return
	}
	if sub != nil && sub.State != model.SubStateForbid {
		count = sub.ACount
	}
	return
}

// GetReplyCounts get reply counts.
func (s *Service) GetReplyCounts(ctx context.Context, oids []int64, otyp int8) (map[int64]*model.Counts, error) {
	caches, missIDs, err := s.dao.Mc.GetMultiSubject(ctx, oids, otyp)
	if err != nil {
		return nil, err
	}
	if len(missIDs) > 0 {
		miss, err := s.dao.Subject.Gets(ctx, missIDs, otyp)
		if err != nil {
			log.Error("s.subject.Gets(%v,%d) error(%v)", missIDs, otyp, err)
			return nil, err
		}
		if caches == nil {
			caches = make(map[int64]*model.Subject)
		}
		var subs []*model.Subject
		for _, oid := range missIDs {
			sub, ok := miss[oid]
			if !ok {
				sub = &model.Subject{ID: -1, Oid: oid, Type: otyp, State: model.SubStateForbid}
			}
			caches[oid] = sub
			subs = append(subs, sub)
		}
		s.cache.Do(ctx, func(ctx context.Context) { s.dao.Mc.AddSubject(ctx, subs...) })
	}
	counts := make(map[int64]*model.Counts)
	for _, sub := range caches {
		counts[sub.Oid] = &model.Counts{
			SubjectState: sub.State,
			Counts:       int64(sub.ACount),
		}
	}
	return counts, nil
}

// ReplyCounts get all reply count.
func (s *Service) ReplyCounts(c context.Context, oids []int64, tp int8) (counts map[int64]int, err error) {
	caches, missIDs, err := s.dao.Mc.GetMultiSubject(c, oids, tp)
	if err != nil {
		return
	}
	if len(missIDs) > 0 {
		var miss map[int64]*model.Subject
		if miss, err = s.dao.Subject.Gets(c, missIDs, tp); err != nil {
			log.Error("s.subject.Gets(%v,%d) error(%v)", missIDs, tp, err)
			return
		}
		if caches == nil {
			caches = make(map[int64]*model.Subject, len(miss))
		}
		subs := make([]*model.Subject, 0, len(miss))
		for _, oid := range missIDs {
			sub, ok := miss[oid]
			if !ok {
				sub = &model.Subject{ID: -1, Oid: oid, Type: tp, State: model.SubStateForbid} // empty cache
			}
			caches[oid] = sub
			subs = append(subs, sub)
		}
		s.cache.Do(c, func(ctx context.Context) { s.dao.Mc.AddSubject(ctx, subs...) })
	}
	counts = make(map[int64]int, len(caches))
	for _, sub := range caches {
		if sub.State == model.SubStateForbid {
			counts[sub.Oid] = 0
		} else {
			counts[sub.Oid] = sub.ACount
		}
	}
	return
}
