package service

import (
	"bytes"
	"context"
	"fmt"
	"math"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/log"
)

const (
	_hisPagesize = 5000
)

// SearchDMHisIndex get history date index.
func (s *Service) SearchDMHisIndex(c context.Context, tp int32, oid int64, month string) (dates []string, err error) {
	if dates, err = s.dao.HistoryIdxCache(c, tp, oid, month); err == nil && len(dates) > 0 {
		return
	}
	if dates, err = s.dao.SearchDMHisIndex(c, tp, oid, month); err != nil {
		log.Error("dao.SearchDMHisIndex(%d,%d,%s) error(%v)", tp, oid, month, err)
		return
	}
	if len(dates) > 0 {
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.AddHisIdxCache(ctx, tp, oid, month, dates)
		})
	}
	return
}

// SearchDMHistory get history dm list from search.
func (s *Service) SearchDMHistory(c context.Context, tp int32, oid, ctimeTo int64) (xml []byte, err error) {
	var (
		sub           *model.Subject
		dmids         []int64
		contentSpeMap = make(map[int64]*model.ContentSpecial)
	)
	if xml, err = s.dao.HistoryCache(c, tp, oid, ctimeTo); err == nil && len(xml) > 0 {
		return
	}
	if sub, err = s.subject(c, tp, oid); err != nil {
		return
	}
	buf := new(bytes.Buffer)
	defer func() {
		if err == nil {
			buf.WriteString(`</i>`)
			xml, err = s.gzflate(buf.Bytes())
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.AddHistoryCache(ctx, tp, oid, ctimeTo, xml)
			})
		}
	}()
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?><i>`)
	buf.WriteString(`<chatserver>chat.bilibili.com</chatserver><chatid>`)
	buf.WriteString(fmt.Sprint(sub.Oid))
	buf.WriteString(`</chatid><mission>`)
	buf.WriteString(fmt.Sprint(sub.AttrVal(model.AttrSubMission)))
	buf.WriteString(`</mission><maxlimit>`)
	buf.WriteString(fmt.Sprint(sub.Maxlimit))
	buf.WriteString(`</maxlimit>`)
	buf.WriteString(fmt.Sprintf(`<state>%d</state>`, sub.State))
	realname := s.isRealname(c, sub.Pid, sub.Oid)
	if realname {
		buf.WriteString(`<real_name>1</real_name>`)
	} else {
		buf.WriteString(`<real_name>0</real_name>`)
	}
	if sub.State == model.SubStateClosed {
		return
	}
	totalPage := int(math.Ceil(float64(sub.Maxlimit) / float64(_hisPagesize)))
	dmids, err = s.dao.SearchDMHistory(c, tp, oid, ctimeTo, totalPage, _hisPagesize)
	if err != nil {
		return
	}
	if len(dmids) == 0 {
		return
	}
	if int64(len(dmids)) > sub.Maxlimit {
		dmids = dmids[:sub.Maxlimit]
	}
	idxMap, special, err := s.dao.IndexsByid(c, tp, oid, dmids)
	if err != nil {
		return
	}
	ctsMap, err := s.dao.Contents(c, oid, dmids)
	if err != nil {
		return
	}
	if len(special) > 0 {
		if contentSpeMap, err = s.dao.ContentsSpecial(c, special); err != nil {
			return
		}
	}
	for _, dmid := range dmids {
		if idx, ok := idxMap[dmid]; ok {
			dm := &model.DM{
				ID:       idx.ID,
				Type:     idx.Type,
				Oid:      idx.Oid,
				Mid:      idx.Mid,
				Progress: idx.Progress,
				Pool:     idx.Pool,
				Attr:     idx.Attr,
				State:    idx.State,
				Ctime:    idx.Ctime,
				Mtime:    idx.Mtime,
			}
			content, ok := ctsMap[dmid]
			if !ok {
				continue
			}
			dm.Content = content
			if idx.Pool == model.PoolSpecial {
				if ctsSpec, ok := contentSpeMap[dmid]; ok {
					dm.ContentSpe = ctsSpec
				}
			}
			buf.WriteString(dm.ToXML(realname))
		}
	}
	return
}

// SearchDMHistoryV2 get history dm list from search.
func (s *Service) SearchDMHistoryV2(c context.Context, tp int32, oid, ctimeTo int64) (res *model.DMSeg, err error) {
	sub, err := s.subject(c, tp, oid)
	if err != nil {
		return
	}
	res = &model.DMSeg{Elems: make([]*model.Elem, 0, sub.Maxlimit)}
	if sub.State == model.SubStateClosed {
		return
	}
	totalPage := int(math.Ceil(float64(sub.Maxlimit) / float64(_hisPagesize)))
	dmids, err := s.dao.SearchDMHistory(c, tp, oid, ctimeTo, totalPage, _hisPagesize)
	if err != nil {
		return
	}
	if len(dmids) == 0 {
		fmt.Println("dmids from search is empty")
		return
	}
	if int64(len(dmids)) > sub.Maxlimit {
		dmids = dmids[:sub.Maxlimit]
	}
	// TODO special dm
	idxMap, _, err := s.dao.IndexsByid(c, tp, oid, dmids)
	if err != nil {
		return
	}
	ctsMap, err := s.dao.Contents(c, oid, dmids)
	if err != nil {
		return
	}
	for _, dmid := range dmids {
		if idx, ok := idxMap[dmid]; ok {
			dm := &model.DM{
				ID:       idx.ID,
				Type:     idx.Type,
				Oid:      idx.Oid,
				Mid:      idx.Mid,
				Progress: idx.Progress,
				Pool:     idx.Pool,
				Attr:     idx.Attr,
				State:    idx.State,
				Ctime:    idx.Ctime,
				Mtime:    idx.Mtime,
			}
			content, ok := ctsMap[dmid]
			if !ok {
				continue
			}
			dm.Content = content
			if e := dm.ToElem(); e != nil {
				res.Elems = append(res.Elems, e)
			}
		}
	}
	return
}
