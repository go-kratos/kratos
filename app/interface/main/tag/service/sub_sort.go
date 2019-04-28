package service

import (
	"context"

	"go-common/app/interface/main/tag/model"
)

var (
	_emptyTagInfo = []*model.TagInfo{}
)

// UpCustomSubTags update sorted tags.
func (s *Service) UpCustomSubTags(c context.Context, mid int64, tids []int64, tp int) (err error) {
	return s.upCustomSubTags(c, mid, tp, tids)
}

// UpCustomSubChannels update sorted channels.
func (s *Service) UpCustomSubChannels(c context.Context, mid int64, tids []int64, tp int) (err error) {
	return s.addCustomSubChannels(c, mid, tp, tids)
}

// CustomSubTags custom sub tags.
func (s *Service) CustomSubTags(c context.Context, mid int64, order, tp, ps, pn int) (cst, sst []*model.TagInfo, total int, err error) {
	var tids []int64
	customTags, err := s.customSubTag(c, mid, tp, pn, ps, order)
	if err != nil {
		return _emptyTagInfo, _emptyTagInfo, 0, err
	}
	if customTags == nil {
		return _emptyTagInfo, _emptyTagInfo, 0, nil
	}
	total = customTags.Total
	for _, tag := range customTags.Sort {
		tids = append(tids, tag.ID)
	}
	for _, tag := range customTags.Tags {
		tids = append(tids, tag.ID)
	}
	countMap, _ := s.countMap(c, tids, mid)
	for _, tag := range customTags.Sort {
		t := &model.TagInfo{
			ID:        tag.ID,
			Type:      tag.Type,
			Name:      tag.Name,
			Cover:     tag.Cover,
			Content:   tag.Content,
			Verify:    tag.Verify,
			Attr:      tag.Attr,
			Attention: tag.Attention,
			State:     tag.State,
			CTime:     tag.CTime,
			MTime:     tag.MTime,
		}
		if k, ok := countMap[tag.ID]; ok {
			t.Bind = k.Bind
			t.Sub = k.Sub
		}
		cst = append(cst, t)
	}
	for _, tag := range customTags.Tags {
		t := &model.TagInfo{
			ID:        tag.ID,
			Type:      tag.Type,
			Name:      tag.Name,
			Cover:     tag.Cover,
			Content:   tag.Content,
			Verify:    tag.Verify,
			Attr:      tag.Attr,
			Attention: tag.Attention,
			State:     tag.State,
			CTime:     tag.CTime,
			MTime:     tag.MTime,
		}
		if k, ok := countMap[tag.ID]; ok {
			t.Bind = k.Bind
			t.Sub = k.Sub
		}
		sst = append(sst, t)
	}
	if cst == nil {
		cst = _emptyTagInfo
	}
	if sst == nil {
		sst = _emptyTagInfo
	}
	return
}
