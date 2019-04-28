package service

import (
	"context"
	"sort"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_emtpyRankTop    = make([]*model.RankTop, 0)
	_emtpyRankView   = make([]*model.RankResult, 0)
	_emtpyRankFilter = make([]*model.RankFilter, 0)
	_emtpyRankNormal = make([]*model.BasicTag, 0)
)

// RegionRankList Region Rank List.
func (s *Service) RegionRankList(c context.Context, prid, rid int64, tp int32) (top []*model.RankTop,
	view []*model.RankResult, filter []*model.RankFilter, normal []*model.BasicTag, err error) {
	var (
		topMap    map[int64]*model.RankTop
		filterMap map[int64]*model.RankFilter
		resultMap map[int64]*model.RankResult
	)
	if top, topMap, _, err = s.dao.RankTop(c, prid, rid, tp); err != nil {
		return
	}
	if resultMap, err = s.dao.RankResult(c, prid, rid, tp); err != nil {
		return
	}
	if filter, filterMap, _, err = s.dao.RankFilter(c, prid, rid, tp); err != nil {
		return
	}
	hotTag, _, _ := s.dao.RegionHot(c, rid)
	for _, v := range resultMap {
		if _, ok := topMap[v.Tid]; ok {
			continue
		}
		view = append(view, v)
	}
	for _, v := range hotTag {
		if _, ok := resultMap[v.ID]; ok {
			continue
		}
		if _, ok := filterMap[v.ID]; ok {
			continue
		}
		normal = append(normal, v)
	}
	if len(view) == 0 {
		view = _emtpyRankView
	}
	if len(top) == 0 {
		top = _emtpyRankTop
	}
	if len(filter) == 0 {
		filter = _emtpyRankFilter
	}
	if len(normal) == 0 {
		normal = _emtpyRankNormal
	}
	sort.Sort(model.RankTopSort(top))
	sort.Sort(model.RankResultSort(view))
	sort.Sort(model.RankFilterSort(filter))
	return
}

// ArchiveRankList ArchiveRankList.
func (s *Service) ArchiveRankList(c context.Context, prid, rid int64, tp int32) (top []*model.RankTop,
	view []*model.RankTop, filter []*model.RankFilter, checked []*model.BasicTag, tags []*model.BasicTag, err error) {
	var (
		tagMap    map[int64]*model.Tag
		filterMap map[int64]*model.RankFilter
	)
	tops, topMap, _, err := s.dao.RankTop(c, prid, rid, tp)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if filter, filterMap, _, err = s.dao.RankFilter(c, prid, rid, tp); err != nil {
		return
	}
	checkMap, tagNames, _ := s.dao.ArchiveHot(c, rid)
	if k, ok := checkMap[rid]; ok {
		checked = k
	}
	if len(tagNames) != 0 {
		_, tagMap, _, _ = s.dao.TagByNames(c, tagNames)
	}
	for id, tag := range tagMap {
		if _, ok := topMap[id]; ok {
			continue
		}
		if _, ok := filterMap[id]; ok {
			continue
		}
		t := &model.BasicTag{
			ID:   tag.ID,
			Name: tag.Name,
		}
		tags = append(tags, t)
	}
	for _, t := range tops {
		if t.Business == model.TagOperateYes {
			top = append(top, t)
		}
		if t.Business == model.TagOperateNO {
			view = append(view, t)
		}
	}
	if len(view) == 0 {
		view = _emtpyRankTop
	}
	if len(top) == 0 {
		top = _emtpyRankTop
	}
	if len(filter) == 0 {
		filter = _emtpyRankFilter
	}
	if len(checked) == 0 {
		checked = _emtpyRankNormal
	}
	if len(tags) == 0 {
		tags = _emtpyRankNormal
	}
	return
}

func (s *Service) mergeTags(c context.Context, top []*model.RankTop, view []*model.RankResult) (res []*model.RankResult, err error) {
	tagMap := make(map[int64]*model.RankResult)
	for k, v := range top {
		if v.Tid == 0 {
			continue
		}
		if _, ok := tagMap[v.Tid]; ok {
			continue
		}
		r := &model.RankResult{
			Tid:       v.Tid,
			TName:     v.TName,
			TagType:   v.TagType,
			Rank:      int64(k),
			HighLight: v.HighLight,
			Business:  v.Business,
			CTime:     v.CTime,
			MTime:     v.MTime,
		}
		tagMap[v.Tid] = r
		res = append(res, r)
	}
	count := int64(len(res))
	for _, v := range view {
		if v.Tid == 0 {
			continue
		}
		if _, ok := tagMap[v.Tid]; ok {
			continue
		}
		v.Rank = count
		tagMap[v.Tid] = v
		res = append(res, v)
		count = count + 1
	}
	return
}

// OperateHotTag OperateHotTag.
func (s *Service) OperateHotTag(c context.Context, tname string) (tag *model.Tag, err error) {
	if tag, err = s.dao.TagByName(c, tname); err != nil {
		return
	}
	if tag == nil || tag.State == model.StateDel {
		return nil, ecode.TagNotExist
	}
	if tag.State == model.StateShield {
		return nil, ecode.TagAlreadyShield
	}
	return
}

// UpdateRank UpdateRank.
func (s *Service) UpdateRank(c context.Context, rank *model.HotRank) (err error) {
	var (
		affect          int64
		r               *model.RankCount
		originresultMap map[int64]*model.RankResult
		diffTids        []int64
	)
	view, _ := s.mergeTags(c, rank.Top, rank.View)
	topCount := len(rank.Top)
	filterCount := len(rank.Filter)
	viewCount := len(view)
	if r, err = s.dao.RankCount(c, rank.Prid, rank.Rid, rank.Type); err != nil {
		return
	}
	if rank.Type == model.HotRegionTag {
		if originresultMap, err = s.dao.RankResult(c, rank.Prid, rank.Rid, rank.Type); err != nil {
			return
		}
		for _, v := range view {
			if _, ok := originresultMap[v.Tid]; ok {
				delete(originresultMap, v.Tid)
			}
		}
		for tid := range originresultMap {
			diffTids = append(diffTids, tid)
		}
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("tag-admin beginTran error(%v)", err)
		return
	}
	if r != nil {
		affect, err = s.dao.TxUpdateRankCount(tx, r.ID, topCount, filterCount, viewCount)
		if err != nil {
			tx.Rollback()
			return
		}
	} else {
		affect, err = s.dao.TxAddRankCount(tx, rank.Prid, rank.Rid, rank.Type, topCount, filterCount, viewCount)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	_, err = s.dao.TxRemoveRankTop(tx, rank.Prid, rank.Rid, rank.Type)
	if err != nil {
		tx.Rollback()
		return
	}
	_, err = s.dao.TxRemoveRankFilter(tx, rank.Prid, rank.Rid, rank.Type)
	if err != nil {
		tx.Rollback()
		return
	}
	_, err = s.dao.TxRemoveRankResult(tx, rank.Prid, rank.Rid, rank.Type)
	if err != nil {
		tx.Rollback()
		return
	}
	if topCount > 0 {
		affect, err = s.dao.TxInsertRankTop(tx, rank.Prid, rank.Rid, rank.Type, rank.Top)
		if err != nil || affect == 0 {
			tx.Rollback()
			return
		}
	}
	if filterCount > 0 {
		affect, err = s.dao.TxInsertRankFilter(tx, rank.Prid, rank.Rid, rank.Type, rank.Filter)
		if err != nil || affect == 0 {
			tx.Rollback()
			return
		}
	}
	if viewCount > 0 {
		affect, err = s.dao.TxInsertRankResult(tx, rank.Prid, rank.Rid, rank.Type, view)
		if err != nil || affect == 0 {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.tagArcsRefresh(ctx, rank.Rid, diffTids...)
	})
	return
}
