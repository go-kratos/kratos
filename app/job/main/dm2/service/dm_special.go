package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"go-common/app/job/main/dm2/model"
)

const (
	_bfsMaxSize = 16 * 1024 * 1024 // size MediumText

	_specialJSONItemSize    = 20 + 1 // {"id":,"content":""},
	_specialJSONAtLeastSize = 2      // []
)

// buildSpeicalDms build when db is no record
func (s *Service) speicalDms(c context.Context, tp int32, oid int64) (dms []*model.DM, err error) {
	var (
		dmids        []int64
		spContentMap map[int64]*model.ContentSpecial
		contentMap   map[int64]*model.Content
	)
	if dms, dmids, err = s.dao.IndexsByPool(c, tp, oid, model.PoolSpecial); err != nil {
		return
	}
	if len(dmids) == 0 {
		return
	}
	if contentMap, err = s.dao.Contents(c, oid, dmids); err != nil {
		return
	}
	if spContentMap, err = s.dao.ContentsSpecial(c, dmids); err != nil {
		return
	}
	for _, dm := range dms {
		if v, ok := contentMap[dm.ID]; ok {
			dm.Content = v
		}
		if v, ok := spContentMap[dm.ID]; ok {
			dm.ContentSpe = v
		}
	}
	sort.Slice(dms, func(i, j int) bool {
		return dms[i].Progress < dms[j].Progress
	})
	return
}

func (s *Service) buildSpecialDms(c context.Context, dms []*model.DM) (bss [][]byte, err error) {
	var (
		dmSpecialContents []*model.DmSpecialContent
		bs                []byte
		length            int
	)
	if len(dms) == 0 {
		return
	}
	dmSpecialContents = make([]*model.DmSpecialContent, 0, len(dms))
	length = _specialJSONAtLeastSize
	for _, dm := range dms {
		if len(dm.GetSpecialSeg()) == 0 {
			continue
		}
		itemSize := len(fmt.Sprint(dm.ID)) + len(dm.GetSpecialSeg()) + _specialJSONItemSize
		if length+itemSize > _bfsMaxSize {
			if bs, err = json.Marshal(dmSpecialContents); err != nil {
				return
			}
			bss = append(bss, bs)
			dmSpecialContents = make([]*model.DmSpecialContent, 0, len(dms))
			length = _specialJSONAtLeastSize
		}
		length += itemSize
		dmSpecialContents = append(dmSpecialContents, &model.DmSpecialContent{
			ID:      dm.ID,
			Content: dm.GetSpecialSeg(),
		})
	}

	if len(dmSpecialContents) > 0 {
		if bs, err = json.Marshal(dmSpecialContents); err != nil {
			return
		}
		bss = append(bss, bs)
	}
	return
}

func (s *Service) updateSpecualDms(c context.Context, tp int32, oid int64, bss [][]byte) (err error) {
	var (
		location  string
		locations []string
		ds        *model.DmSpecial
	)
	for _, bs := range bss {
		if len(bs) == 0 {
			continue
		}
		if location, err = s.dao.BfsDmUpload(c, "", bs); err != nil {
			return
		}
		locations = append(locations, location)
	}
	ds = &model.DmSpecial{
		Type: tp,
		Oid:  oid,
	}
	ds.Join(locations)
	if err = s.dao.UpsertDmSpecialLocation(c, ds.Type, ds.Oid, ds.Locations); err != nil {
		return
	}
	return
}

func (s *Service) specialLocationUpdate(c context.Context, tp int32, oid int64) (err error) {
	var (
		dms []*model.DM
		bss [][]byte
	)
	if dms, err = s.speicalDms(c, tp, oid); err != nil {
		return
	}
	if bss, err = s.buildSpecialDms(c, dms); err != nil {
		return
	}
	if err = s.updateSpecualDms(c, tp, oid, bss); err != nil {
		return
	}
	return
}
