package service

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/service/main/workflow/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// TagSlice get Tagslice
func (s *Service) TagSlice(business int8) []*model.Tag {
	return s.tagsCache.TagSlice[business]
}

// TagMap get TagMap
func (s *Service) TagMap(business int8, tid int32) *model.Tag {
	return s.tagsCache.TagMap[business][tid]
}

// tags get all tag list
func (s *Service) tags(c context.Context) (tagsCache *model.TagsCache, err error) {
	var (
		tags     []*model.Tag
		tagsMap  = make(map[int8]map[int32]*model.Tag)
		tagSlice = make(map[int8][]*model.Tag)
	)
	if err = s.dao.DB.Where("state=?", 1).Order("weight desc").Find(&tags).Error; err != nil {
		log.Error("d.DB.Where().Find() error(%v)", err)
		return
	}
	for _, v := range tags {
		if tagsMap[v.Business] == nil {
			tagsMap[v.Business] = make(map[int32]*model.Tag)
		}
		tagsMap[v.Business][v.Tid] = v
		tagSlice[v.Business] = append(tagSlice[v.Business], v)
	}
	tagsCache = &model.TagsCache{}
	tagsCache.TagMap = tagsMap
	tagsCache.TagSlice = tagSlice
	return
}

// tags loading...
func (s *Service) loadTags() {
	var (
		tagsCache *model.TagsCache
		err       error
	)
	if tagsCache, err = s.tags(context.TODO()); err != nil {
		log.Error("s.tag.Tags() error(%v)", err)
		return
	}
	for business, tags := range tagsCache.TagSlice {
		var tagTidList []int32
		var controls []*model.Control
		for _, tag := range tags {
			tagTidList = append(tagTidList, tag.Tid)
		}
		if err = s.dao.DB.Where("tid in (?)", tagTidList).Order("weight").Find(&controls).Error; err != nil {
			log.Error("s.tag.Controls(%d) error(%v)", tagTidList, err)
			return
		}
		for _, control := range controls {
			tagsCache.TagMap[business][control.Tid].Controls = append(tagsCache.TagMap[business][control.Tid].Controls, control)
		}
	}
	s.tagsCache = tagsCache
}

// Tag3 .
func (s *Service) Tag3(bid, rid int64) (res []*model.Tag3) {
	var (
		ok     = false
		bidRes = make(map[int64][]*model.Tag3)
	)
	res = []*model.Tag3{}
	if bidRes, ok = s.tagsCache.TagMap3[bid]; !ok {
		log.Error("nothing to get tag list by bid(%d)", bid)
		return
	}
	if rid != 0 {
		if res, ok = bidRes[rid]; !ok {
			log.Error("nothing to get tag list by bid(%d) rid(%d)", bid, rid)
		}
		return
	}
	for _, br := range bidRes {
		res = append(res, br...)
	}
	return
}

// tags3 .
func (s *Service) tags3(c context.Context) (res []*model.Tag3, err error) {
	res = []*model.Tag3{}
	tmpTag3 := model.ResponseTag3{}
	tagUV := url.Values{}
	tagUV.Set("ps", strconv.Itoa(model.ControlPageSize))
	if err = s.dao.ReadClient.Get(c, s.dao.MngTagURL, "", tagUV, &tmpTag3); err != nil {
		log.Error("s.dao.ReadClient.Get error(%v)", err)
		return
	}
	if tmpTag3.Code != ecode.OK.Code() {
		err = ecode.Int(tmpTag3.Code)
		log.Error("get tag info error(%v)", err)
		return
	}
	tmpControl3 := model.ResponseControl3{}
	controlUV := url.Values{}
	if err = s.dao.ReadClient.Get(c, s.dao.MngControlURL, "", controlUV, &tmpControl3); err != nil {
		log.Error("get control list error(%v)", err)
		return
	}
	if tmpControl3.Code != ecode.OK.Code() {
		err = ecode.Int(tmpControl3.Code)
		log.Error("get control list error(%v)", err)
		return
	}
	res = tmpTag3.Data.Data
	controls := tmpControl3.Data
	for _, r := range res {
		for _, c := range controls {
			if r.BID == c.BID && r.TagID == c.TID {
				r.Controls = append(r.Controls, c)
			}
		}
	}
	return
}

// loadTags3 .
func (s *Service) loadTags3() {
	s.tagsCache.TagMap3 = map[int64]map[int64][]*model.Tag3{}
	s.tagsCache.TagMap3Tid = map[int64]map[int64]*model.Tag3{}
	data, err := s.tags3(context.Background())
	if err != nil {
		log.Error("s.tags3 error(%v)", err)
		return
	}
	for _, d := range data {
		if _, ok := s.tagsCache.TagMap3[d.BID]; !ok {
			s.tagsCache.TagMap3[d.BID] = make(map[int64][]*model.Tag3)
		}
		if _, ok := s.tagsCache.TagMap3Tid[d.BID]; !ok {
			s.tagsCache.TagMap3Tid[d.BID] = make(map[int64]*model.Tag3)
		}
		s.tagsCache.TagMap3[d.BID][d.RID] = append(s.tagsCache.TagMap3[d.BID][d.RID], d)
		s.tagsCache.TagMap3Tid[d.BID][d.TagID] = d
	}
}
