package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_dmRecentLimit = 1000

	_dateTimeFormart = "2006-01-02 15:04:05"
)

// DMUpRecent recent dm of upper.
func (s *Service) DMUpRecent(c context.Context, mid, pn, ps int64) (res *model.DmRecentResponse, err error) {
	var (
		mids, aids   []int64
		aidmap       = make(map[int64]struct{})
		midMap       = make(map[int64]struct{})
		searchResult *model.SearchRecentDMResult
	)
	if ps < 0 || pn < 1 {
		err = ecode.RequestErr
		return
	}
	if (pn-1)*ps >= _dmRecentLimit {
		return
	}
	searchParam := &model.SearchRecentDMParam{
		Type:   model.SubTypeVideo,
		UpMid:  mid,
		States: []int32{model.StateNormal, model.StateHide, model.StateMonitorAfter},
		Ps:     int(ps),
		Pn:     int(pn),
		Field:  "ctime",
		Sort:   elastic.OrderDesc,
	}
	searchResult, err = s.dao.SearhcDmRecent(c, searchParam)
	if err != nil || searchResult == nil || len(searchResult.Result) == 0 || searchResult.Page == nil {
		return
	}
	for _, item := range searchResult.Result {
		if _, ok := aidmap[item.Aid]; !ok {
			aids = append(aids, item.Aid)
			aidmap[item.Aid] = struct{}{}
		}
		if _, ok := midMap[item.Mid]; !ok {
			mids = append(mids, item.Mid)
			midMap[item.Mid] = struct{}{}
		}
	}
	arcMap, err := s.archiveInfos(c, aids)
	if err != nil {
		return
	}
	infoMap, err := s.accountInfos(c, mids)
	if err != nil {
		return
	}
	memebers := make([]*model.DMMember, 0, len(searchResult.Result))
	for _, item := range searchResult.Result {
		member := &model.DMMember{
			ID:       item.ID,
			IDStr:    strconv.FormatInt(item.ID, 10),
			Type:     item.Type,
			Aid:      item.Aid,
			Oid:      item.Oid,
			Mid:      item.Mid,
			MidHash:  model.Hash(item.Mid, 0),
			Pool:     item.Pool,
			State:    item.State,
			Attrs:    model.DMAttrNtoA(item.Attr),
			Msg:      item.Msg,
			Mode:     item.Mode,
			Color:    fmt.Sprintf("%06x", item.Color),
			Progress: item.Progress,
			FontSize: item.FontSize,
		}
		if ctime, err := time.ParseInLocation(_dateTimeFormart, item.Ctime, time.Now().Location()); err == nil {
			member.Ctime = xtime.Time(ctime.Unix())
		}
		if arc, ok := arcMap[item.Aid]; ok {
			member.Title = arc.Title
		}
		if info, ok := infoMap[item.Mid]; ok {
			member.Uname = info.Name
		}
		memebers = append(memebers, member)
	}
	res = &model.DmRecentResponse{
		Data: memebers,
		Page: searchResult.Page,
	}
	if res.Page.Total > _dmRecentLimit {
		res.Page.Total = _dmRecentLimit
	}
	return
}

// DMUpSearch danmu list from search.
func (s *Service) DMUpSearch(c context.Context, mid int64, p *model.SearchDMParams) (res *model.SearchDMResult, err error) {
	var (
		mids, dmids []int64
	)
	sub, err := s.subject(c, p.Type, p.Oid)
	if err != nil {
		return
	}
	if sub.Mid != mid {
		err = ecode.AccessDenied
		return
	}
	res = &model.SearchDMResult{}
	srchData, err := s.dao.SearchDM(c, p)
	if err != nil || srchData == nil {
		return
	}
	for _, v := range srchData.Result {
		dmids = append(dmids, v.ID)
	}
	dms, err := s.dmList(c, p.Type, p.Oid, dmids)
	if err != nil {
		log.Error("s.dms(%d,%v) error(%v)", p.Oid, dmids, err)
		return
	}
	for _, dm := range dms {
		mids = append(mids, dm.Mid)
	}
	infoMap, err := s.accountInfos(c, mids)
	if err != nil {
		return
	}
	for _, dm := range dms {
		var msg string
		if dm.Content != nil {
			msg = dm.Content.Msg
		} else {
			continue
		}
		if dm.ContentSpe != nil {
			msg = dm.ContentSpe.Msg
		}
		item := &model.DMMember{
			ID:       dm.ID,
			IDStr:    strconv.FormatInt(dm.ID, 10),
			Type:     dm.Type,
			Aid:      sub.Pid,
			Oid:      dm.Oid,
			Mid:      dm.Mid,
			MidHash:  model.Hash(dm.Mid, 0),
			Pool:     dm.Pool,
			State:    dm.State,
			Attrs:    dm.AttrNtoA(),
			Msg:      msg,
			Ctime:    dm.Ctime,
			Mode:     dm.Content.Mode,
			Color:    fmt.Sprintf("%06x", dm.Content.Color),
			Progress: dm.Progress,
			FontSize: dm.Content.FontSize,
		}
		if info, ok := infoMap[dm.Mid]; ok {
			item.Uname = info.Name
		}
		res.Result = append(res.Result, item)
	}
	res.Page.Num = srchData.Page.Num
	res.Page.Size = srchData.Page.Size
	res.Page.Total = srchData.Page.Total
	return
}

// UptSearchDMState update dm search state
func (s *Service) UptSearchDMState(c context.Context, dmids []int64, oid int64, state, tp int32) (err error) {
	if err = s.dao.UptSearchDMState(c, dmids, oid, state, tp); err != nil {
		return
	}
	if err = s.dao.UptSearchRecentState(c, dmids, oid, state, tp); err != nil {
		return
	}
	return
}

// UptSearchDMPool update dm search pool
func (s *Service) UptSearchDMPool(c context.Context, dmids []int64, oid int64, pool, tp int32) (err error) {
	if err = s.dao.UptSearchDMPool(c, dmids, oid, pool, tp); err != nil {
		return
	}
	if err = s.dao.UptSearchRecentPool(c, dmids, oid, pool, tp); err != nil {
		return
	}
	return
}

// UptSearchDMAttr update dm search attr
func (s *Service) UptSearchDMAttr(c context.Context, dmids []int64, oid int64, attr, tp int32) (err error) {
	if err = s.dao.UptSearchDMAttr(c, dmids, oid, attr, tp); err != nil {
		return
	}
	if err = s.dao.UptSearchRecentAttr(c, dmids, oid, attr, tp); err != nil {
		return
	}
	return
}
