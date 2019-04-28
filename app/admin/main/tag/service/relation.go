package service

import (
	"context"
	"time"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_emtpyResource = make([]*model.Resource, 0)
	_emtpyRelation = make([]*model.Resource, 0)
)

// RelationListByTag RelationListByTag.
func (s *Service) RelationListByTag(c context.Context, tname string, pn, ps int32) (total int64, res []*model.Resource, err error) {
	var (
		oids, authorIDs []int64
		tag             *model.Tag
		tagCount        *model.TagCount
	)
	start := (pn - 1) * ps
	end := ps
	arcMap := make(map[int64]*model.SearchRes)
	authorMap := make(map[int64]*model.UserInfo)
	if tag, err = s.dao.TagByName(c, tname); err != nil {
		return
	}
	if tag == nil {
		return 0, nil, ecode.TagNotExist
	}
	total, _ = s.dao.TagResCount(c, tag.ID)
	if res, oids, err = s.dao.RelationsByTid(c, tag.ID, start, end); err != nil {
		return 0, _emtpyRelation, err
	}
	if res == nil {
		return 0, _emtpyResource, nil
	}
	if len(oids) != 0 {
		arcMap, authorIDs, _ = s.arcInfos(c, oids)
	}
	if len(authorIDs) != 0 {
		authorMap, _ = s.userInfos(c, authorIDs)
	}
	tagCount, _ = s.dao.TagCount(c, tag.ID)
	if tagCount == nil {
		tagCount = _emtpyTagCount
	}
	for _, v := range res {
		v.Tag = tag
		v.TagCount = tagCount
		if arc, ok := arcMap[v.Oid]; ok {
			v.Title = arc.Title
		}
		if u, ok := authorMap[v.Mid]; ok {
			v.Author = u.Name
		}
	}
	return
}

// RelationListByOid RelationListByOid.
func (s *Service) RelationListByOid(c context.Context, oid int64, tp, pn, ps int32) (total int64, res []*model.Resource, err error) {
	var (
		tids []int64
		arc  *model.SearchRes
	)
	start := (pn - 1) * ps
	end := ps
	tagMap := make(map[int64]*model.Tag)
	tagCountMap := make(map[int64]*model.TagCount)
	total, _ = s.dao.ResTagCount(c, oid, tp)
	if arc, err = s.arcInfo(c, oid); err != nil {
		return
	}
	if arc == nil {
		return 0, nil, ecode.ArchiveNotExist
	}
	if res, tids, err = s.dao.RelationsByOid(c, oid, tp, start, end); err != nil {
		return 0, _emtpyRelation, err
	}
	if len(tids) != 0 {
		_, tagMap, _ = s.dao.Tags(c, tids)
		tagCountMap, _ = s.dao.TagCounts(c, tids)
	}
	for _, v := range res {
		if k, ok := tagMap[v.Tid]; ok {
			if k != nil {
				v.Tag = k
			} else {
				v.Tag = _emtpyTag
			}
		}
		if k, ok := tagCountMap[v.Tid]; ok {
			if k != nil {
				v.TagCount = k
			} else {
				v.TagCount = _emtpyTagCount
			}
		}
	}
	return
}

// RelationAdd RelationAdd.
func (s *Service) RelationAdd(c context.Context, tname string, oid, mid int64, tp int32) (err error) {
	var (
		affect int64
		tag    *model.Tag
		r      *model.Resource
		arc    *model.SearchRes
	)
	if err = s.dao.Filter(c, tname); err != nil {
		return
	}
	if tag, err = s.dao.TagByName(c, tname); err != nil {
		return
	}
	if tag == nil {
		return ecode.TagNotExist
	}
	if arc, err = s.arcInfo(c, oid); err != nil {
		return
	}
	if arc == nil {
		return ecode.ArchiveNotExist
	}
	if r, err = s.dao.Relation(c, oid, tag.ID, tp, model.RelationStateNormal); err != nil {
		return
	}
	if r != nil {
		return ecode.TagResTagExist
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("relation add beginTran error(%v)", err)
		return
	}
	relation := &model.Relation{
		Oid:   oid,
		Type:  tp,
		Tid:   tag.ID,
		Mid:   mid,
		Role:  model.ResRoleAdmin,
		Enjoy: 0,
		Hate:  0,
		Attr:  model.AttrLockNone,
		State: model.RelationStateNormal,
	}
	affect, err = s.dao.TxInsertResTag(tx, relation)
	if err != nil || affect == 0 {
		tx.Rollback()
		return
	}
	affect, err = s.dao.TxInsertTagRes(tx, relation)
	if err != nil || affect == 0 {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.dao.DelResMemCache(ctx, oid, tag.ID, tp)
		s.dao.DelResTagCache(ctx, oid, tp)
		s.dao.DelRelationCache(ctx, oid, tag.ID, tp)
	})
	return
}

func (s *Service) arcInfo(c context.Context, oid int64) (*model.SearchRes, error) {
	arcMap, err := s.dao.ESearchArchives(c, []int64{oid})
	if err != nil {
		return nil, err
	}
	if k, ok := arcMap[oid]; ok && k != nil {
		return k, nil
	}
	return nil, nil
}

func (s *Service) arcInfos(c context.Context, oids []int64) (arcMap map[int64]*model.SearchRes, authors []int64, err error) {
	if arcMap, err = s.dao.ESearchArchives(c, oids); err != nil {
		return
	}
	authors = make([]int64, 0, len(arcMap))
	for _, arc := range arcMap {
		authors = append(authors, arc.Mid)
	}
	return
}

// RelationLock RelationLock.
func (s *Service) RelationLock(c context.Context, tid, oid int64, tp int32) (err error) {
	var (
		affect int64
		tag    *model.Tag
	)
	if tag, err = s.dao.Tag(c, tid); err != nil {
		return
	}
	if tag == nil {
		return ecode.TagNotExist
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("relation lock beginTran error(%v)", err)
		return
	}
	affect, err = s.dao.TxUpdateAttrResTag(tx, tid, oid, tp, model.AttrLocked)
	if err != nil || affect == 0 {
		tx.Rollback()
		return
	}
	affect, err = s.dao.TxUpdateAttrTagRes(tx, tid, oid, tp, model.AttrLocked)
	if err != nil || affect == 0 {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.dao.DelResMemCache(ctx, oid, tid, tp)
		s.dao.DelResTagCache(ctx, oid, tp)
	})
	return
}

// RelationUnLock RelationUnLock.
func (s *Service) RelationUnLock(c context.Context, tid, oid int64, tp int32) (err error) {
	var (
		affect int64
		tag    *model.Tag
	)
	if tag, err = s.dao.Tag(c, tid); err != nil {
		return
	}
	if tag == nil {
		return ecode.TagNotExist
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("relation unlock beginTran error(%v)", err)
		return
	}
	affect, err = s.dao.TxUpdateAttrResTag(tx, tid, oid, tp, model.AttrLockNone)
	if err != nil || affect == 0 {
		tx.Rollback()
		return
	}
	affect, err = s.dao.TxUpdateAttrTagRes(tx, tid, oid, tp, model.AttrLockNone)
	if err != nil || affect == 0 {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.dao.DelResMemCache(ctx, oid, tid, tp)
		s.dao.DelResTagCache(ctx, oid, tp)
	})
	return
}

// RelationDelete RelationDelete.
func (s *Service) RelationDelete(c context.Context, tid, oid int64, tp int32) (err error) {
	var (
		affect int64
		tag    *model.Tag
	)
	if tag, err = s.dao.Tag(c, tid); err != nil {
		return
	}
	if tag == nil {
		return ecode.TagNotExist
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("relation delete beginTran error(%v)", err)
		return
	}
	affect, err = s.dao.TxUpdateStateResTag(tx, tid, oid, tp, model.RelationStateDelete)
	if err != nil || affect == 0 {
		tx.Rollback()
		return
	}
	affect, err = s.dao.TxUpdateStateTagRes(tx, tid, oid, tp, model.RelationStateDelete)
	if err != nil || affect == 0 {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.dao.DelResMemCache(ctx, oid, tid, tp)
		s.dao.DelResTagCache(ctx, oid, tp)
		s.dao.DelRelationCache(ctx, oid, tid, tp)
	})
	return
}

// tagArcsRefresh tag archives to redis cache.
func (s *Service) tagArcsRefresh(c context.Context, rid int64, tids ...int64) (err error) {
	for _, tid := range tids {
		var (
			batchSize, size = int(1000), int(1000)
			start           = int32(0)
		)
		for size == batchSize {
			var aids = make([]int64, 0, batchSize)
			if _, aids, err = s.dao.RelationsByTid(c, tid, start, int32(batchSize)); err != nil {
				return
			}
			log.Warn("s.tagArcsRefresh(rid:%d tid:%d) this is mysql aids(%v)", rid, tid, aids)
			size = len(aids)
			start = start + int32(batchSize)
			if len(aids) == 0 {
				continue
			}
			var arcMap map[int64]*model.SearchRes
			if arcMap, _, err = s.arcInfos(c, aids); err != nil {
				return
			}
			if len(arcMap) == 0 {
				continue
			}
			oids := make([]int64, 0, len(arcMap))
			for _, v := range arcMap {
				if v.TypeID == rid && v.State == 0 {
					oids = append(oids, v.ID)
				}
			}
			log.Warn("s.tagArcsRefresh(rid:%d tid:%d) this is es resource aids(%v)", rid, tid, oids)
			if err = s.dao.AddTagArcs(c, tid, arcMap); err != nil {
				return
			}
			time.Sleep(time.Second * 1)
		}
	}
	return
}

// RegionTagArcsRefresh tag arcs refresh.
func (s *Service) RegionTagArcsRefresh(c context.Context, rid int64, tid int64) (err error) {
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.tagArcsRefresh(ctx, rid, tid)
	})
	return
}
