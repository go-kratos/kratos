package service

import (
	"context"
	"time"

	"go-common/app/service/main/tag/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// Info return a tag by tid.
func (s *Service) Info(c context.Context, mid, tid int64) (res *model.Tag, err error) {
	if res, err = s.tag(c, tid); err != nil {
		return
	}
	if res == nil {
		err = ecode.TagNotExist
		return
	}
	res.Attention = 0
	if mid > 0 {
		ok, _ := s.isSubTid(c, mid, tid)
		if ok {
			res.Attention = 1
		}
	}
	return
}

// Infos return tags by tids.
func (s *Service) Infos(c context.Context, mid int64, tids []int64) (tags []*model.Tag, err error) {
	if tags, err = s.tags(c, tids); err != nil {
		return
	}
	if len(tags) == 0 {
		err = ecode.TagNotExist
		return
	}
	if mid > 0 {
		subMap, _ := s.isSubTids(c, mid, tids)
		for _, tag := range tags {
			if sub, ok := subMap[tag.ID]; ok && sub == 1 {
				tag.Attention = 1
			} else {
				tag.Attention = 0 // TODO mc cache data err： Attention = 1
			}
		}
	}
	return
}

// InfoMap return tag map by tids.
func (s *Service) InfoMap(c context.Context, mid int64, tids []int64) (tagMap map[int64]*model.Tag, err error) {
	tags, err := s.tags(c, tids)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, ecode.TagNotExist
	}
	tagMap = make(map[int64]*model.Tag)
	if mid > 0 {
		subMap, _ := s.isSubTids(c, mid, tids)
		for _, tag := range tags {
			if sub, ok := subMap[tag.ID]; ok && sub == 1 {
				tag.Attention = 1
			} else {
				tag.Attention = 0 // TODO mc cache data err： Attention = 1
			}
			tagMap[tag.ID] = tag
		}
		return
	}
	for _, tag := range tags {
		tagMap[tag.ID] = tag
	}
	return
}

// InfoByName return a tag by name.
func (s *Service) InfoByName(c context.Context, mid int64, name string) (t *model.Tag, err error) {
	if t, err = s.tagByName(c, name); err != nil {
		return
	}
	if t == nil {
		return nil, ecode.TagNotExist
	}
	t.Attention = 0
	if mid > 0 {
		ok, _ := s.isSubTid(c, mid, t.ID)
		if ok {
			t.Attention = 1
		}
	}
	return
}

// InfosByNames return tags by names.
func (s *Service) InfosByNames(c context.Context, mid int64, names []string) (tags []*model.Tag, err error) {
	if tags, err = s.tagByNames(c, names); err != nil {
		return
	}
	if len(tags) == 0 {
		err = ecode.TagNotExist
		return
	}
	if mid > 0 {
		var tids []int64
		for _, tag := range tags {
			tids = append(tids, tag.ID)
		}
		subMap, _ := s.isSubTids(c, mid, tids)
		for k, t := range tags {
			if sub, ok := subMap[t.ID]; ok && sub == 1 {
				tags[k].Attention = 1
			} else {
				tags[k].Attention = 0 // TODO mc cache data err： Attention = 1
			}
		}
	}
	return
}

func (s *Service) tagMap(c context.Context, tids []int64) (tagMap map[int64]*model.Tag, err error) {
	var (
		miss []int64
		tags []*model.Tag
	)
	if tagMap, miss, err = s.dao.TagMapCaches(c, tids); err != nil {
		return
	}
	if len(miss) == 0 {
		return
	}
	if tags, err = s.dao.Tags(c, miss); err != nil {
		return
	}
	if len(tagMap) == 0 {
		tagMap = make(map[int64]*model.Tag, len(tids))
	}
	for _, tag := range tags {
		tagMap[tag.ID] = tag
	}
	if len(tags) > 0 {
		var tc []*model.Tag
		for _, v := range tags {
			if v.State == model.TagStateDelete {
				continue
			}
			t := *v
			tc = append(tc, &t)
		}
		s.cache.Save(func() {
			s.dao.AddTagsCache(context.Background(), tc)
		})
	}
	return
}

func (s *Service) tag(c context.Context, tid int64) (tag *model.Tag, err error) {
	if tag, err = s.dao.TagCache(c, tid); err != nil {
		return
	}
	if tag != nil {
		return
	}
	if tag, err = s.dao.Tag(c, tid); err != nil {
		return
	}
	if tag != nil && tag.State != model.TagStateDelete {
		t := *tag
		s.cache.Save(func() {
			s.dao.AddTagCache(context.Background(), &t)
		})
	}
	return
}

func (s *Service) tags(c context.Context, tids []int64) (res []*model.Tag, err error) {
	var miss []int64
	if res, miss, err = s.dao.TagsCaches(c, tids); err != nil {
		return
	}
	if len(miss) == 0 {
		return
	}
	tags, err := s.dao.Tags(c, miss)
	if err != nil {
		return
	}
	if len(tags) > 0 {
		res = append(res, tags...)
		var tc []*model.Tag
		for _, v := range tags {
			if v.State == model.TagStateDelete {
				continue
			}
			t := *v
			tc = append(tc, &t)
		}
		s.cache.Save(func() {
			s.dao.AddTagsCache(context.Background(), tc)
		})
	}
	return
}

func (s *Service) tagByName(c context.Context, name string) (tag *model.Tag, err error) {
	if tag, err = s.dao.TagCacheByName(c, name); err != nil {
		return
	}
	if tag != nil {
		return
	}
	if tag, err = s.dao.TagByName(c, name); err != nil {
		return
	}
	if tag != nil && tag.State != model.TagStateDelete {
		t := *tag
		s.cache.Save(func() {
			s.dao.AddTagCache(context.Background(), &t)
		})
	}
	return
}

func (s *Service) tagByNames(c context.Context, names []string) (res []*model.Tag, err error) {
	var miss []string
	if res, miss, err = s.dao.TagCachesByNames(c, names); err != nil {
		return
	}
	if len(miss) == 0 {
		return
	}
	tags, err := s.dao.TagsByNames(c, miss)
	if err != nil {
		return
	}
	if len(tags) > 0 {
		res = append(res, tags...)
		var tc []*model.Tag
		for _, v := range tags {
			if v.State == model.TagStateDelete {
				continue
			}
			t := *v
			tc = append(tc, &t)
		}
		s.cache.Save(func() {
			s.dao.AddTagsCache(context.Background(), tc)
		})
	}
	return
}

//RecommandTag get recomand tags from cache, and where cache not exist, only allow one goroutine top update cache from mysql
func (s *Service) RecommandTag(c context.Context) (res map[int64]map[string][]*model.UploadTag, err error) {
	var (
		rids                             []int64
		recommandTagFil, recommandTagTop []*model.UploadTag
	)
	recommandTagFil, err = s.dao.RecommandTagFilter(c)
	if err != nil {
		log.Error("s.dao.RecommandTagFilter error(%v)", err)
		return nil, nil
	}
	recommandTagTop, err = s.dao.RecommandTagTop(c)
	if err != nil {
		log.Error("s.dao.RecommandTagTop error(%v)", err)
		return
	}
	if _, _, rids, err = s.dao.Rids(c); err != nil {
		log.Error("s.tag.Rids error(%v)", err)
		return
	}
	res = make(map[int64]map[string][]*model.UploadTag)
	for _, rid := range rids {
		tempMap := make(map[string][]*model.UploadTag)
		for _, t := range recommandTagFil {
			if rid == t.Rid {
				tempMap["filter"] = append(tempMap["filter"], t)
			}
		}
		for _, t := range recommandTagTop {
			if rid == t.Rid {
				switch t.IsBusiness {
				case 0:
					tempMap["recommend"] = append(tempMap["recommend"], t)
				case 1:
					tempMap["business"] = append(tempMap["business"], t)
				}
			}
		}
		res[rid] = tempMap
	}
	return
}

// TagGroup .
func (s *Service) TagGroup(c context.Context) (rs []*model.Synonym, err error) {
	var resMap map[int64][]int64
	resMap, err = s.dao.TagGroup(c)
	for parent, child := range resMap {
		res := &model.Synonym{
			Parent: parent,
			Childs: child,
		}
		rs = append(rs, res)
	}
	return
}

// Count return tag count .
func (s *Service) Count(c context.Context, tid int64) (res *model.Count, err error) {
	if res, err = s.dao.CountCache(c, tid); err != nil {
		return
	}
	if res == nil {
		if res, err = s.dao.Count(c, tid); err != nil {
			return
		}
		s.cache.Save(func() {
			s.dao.AddCountCache(context.Background(), res)
		})
	}
	return
}

// Counts return tags`count .`
func (s *Service) Counts(c context.Context, tids []int64) (res map[int64]*model.Count, err error) {
	var (
		miss   []int64
		counts []*model.Count
	)
	if len(tids) == 0 {
		res = make(map[int64]*model.Count)
		return
	}
	if res, miss, err = s.dao.CountMapCache(c, tids); err != nil {
		return
	}
	if res == nil {
		res = make(map[int64]*model.Count)
	}
	if len(miss) > 0 {
		if counts, err = s.dao.Counts(c, miss); err != nil {
			return
		}
		for _, v := range counts {
			res[v.Tid] = v
		}
		s.cache.Save(func() {
			s.dao.AddCountsCache(context.Background(), counts)
		})
	}
	return
}

// HideTag .
func (s *Service) HideTag(c context.Context, tid int64, state int32) (err error) {
	var row int64
	if row, err = s.dao.UpTagState(c, tid, state); err != nil || row == 0 {
		return
	}
	s.dao.DelTagCache(c, tid)
	return
}

// CreateTag .
func (s *Service) CreateTag(c context.Context, t *model.Tag) (err error) {
	return s.dao.CreateTag(c, t)
}

// CreateTags .
func (s *Service) CreateTags(c context.Context, ts []*model.Tag) (err error) {
	return s.dao.CreateTags(c, ts)
}

// CheckTag .
func (s *Service) CheckTag(c context.Context, name string, tp int32, now time.Time, ip string) (t *model.Tag, err error) {
	var (
		tx   *xsql.Tx
		xnow = xtime.Time(now.Unix())
	)
	if t, err = s.tagByName(c, name); err != nil {
		return
	}
	if t != nil {
		if t.State != model.TagStateNormal {
			t = nil
			err = ecode.TagIsSealing
		}
		return
	}
	t = &model.Tag{
		Name:  name,
		Type:  tp,
		State: model.TagStateNormal,
		CTime: xnow,
		MTime: xnow,
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	if t.ID, err = s.dao.TxAddTag(tx, t); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.TxUpTagBindCount(tx, t.ID, 1); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

// Hots .
func (s *Service) Hots(c context.Context, rid, hotType int64) ([]*model.HotTag, error) {
	return s.dao.Hots(c, rid, hotType)
}

// HotMap .
func (s *Service) HotMap(c context.Context) (res map[int16][]int64, err error) {
	return s.dao.HotMap(c)
}

// Prids .
func (s *Service) Prids(c context.Context) (res []int16, err error) {
	return s.dao.Prids(c)
}

// Rids .
func (s *Service) Rids(c context.Context) (pridMap map[int64]int64, err error) {
	_, pridMap, _, err = s.dao.Rids(c)
	return
}
