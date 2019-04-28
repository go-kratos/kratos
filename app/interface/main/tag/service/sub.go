package service

import (
	"context"
	"time"

	"go-common/app/interface/main/tag/model"
	rpcModel "go-common/app/service/main/tag/model"
)

var _emptyTs = []*model.Tag{}

// SubTags .
func (s *Service) SubTags(c context.Context, mid, vmid int64, pn, ps, order int) (ts []*model.Tag, total int, err error) {
	ts, _, total, err = s.subTags(c, mid, vmid, pn, ps, order)
	return
}

// SubTags .
func (s *Service) subTags(c context.Context, mid, rmid int64, pn, ps, order int) (ts []*model.Tag, tids []int64, total int, err error) {
	if rmid == 0 {
		rmid = mid
	}
	var subs *rpcModel.ResSub
	subs, err = s.subTag(c, rmid, pn, ps, order)
	if err != nil {
		return nil, nil, 0, err
	}
	for _, v := range subs.Tags {
		tids = append(tids, v.ID)
		ts = append(ts, &model.Tag{
			ID:        v.ID,
			Type:      int8(v.Type),
			Name:      v.Name,
			Cover:     v.Cover,
			Content:   v.Content,
			Attribute: int8(v.Attr),
			State:     int8(v.State),
			IsAtten:   1,
			CTime:     v.CTime,
			MTime:     v.MTime,
		})
	}
	total = subs.Total
	return ts, tids, total, nil
}

// AddSub .
func (s *Service) AddSub(c context.Context, mid int64, tids []int64, now time.Time) (err error) {
	return s.addSub(c, mid, tids)
}

// CancelSub .
func (s *Service) CancelSub(c context.Context, tid, mid int64, now time.Time) (err error) {
	return s.cancelSub(c, mid, tid)
}

// SubArcs get newest arcs of user subcribe
func (s *Service) SubArcs(c context.Context, mid, rmid int64) (as []*model.SubArcs, err error) {
	var (
		aids  []int64
		count int
		ts    []*model.Tag
	)
	// get real mid for subscribed tags
	if rmid == 0 {
		rmid = mid
	}
	// get need add sub tids
	ts, _, count, err = s.subTags(c, mid, 0, 1, model.SubTagMaxNum, -1)
	if err != nil {
		return
	}
	if len(ts) == 0 || count == 0 {
		return
	}
	// arcs num of one tid
	ps := s.c.Tag.SubArcMaxNum / count
	for _, t := range ts {
		// get new arcs of tag
		if aids, _, err = s.newArcs(c, t.ID, 0, ps); err != nil {
			continue
		}
		if len(aids) == 0 {
			continue
		}
		t.IsAtten = 1
		as = append(as, &model.SubArcs{
			Tag:  t,
			Aids: aids,
		})
	}
	return
}

func (s *Service) attened(c context.Context, mid int64) (tm map[int64]*model.Tag, err error) {
	var ts []*model.Tag
	tm = make(map[int64]*model.Tag)
	ts, _, _, err = s.subTags(c, mid, 0, 1, model.SubTagMaxNum, -1)
	if err != nil {
		return
	}
	if len(ts) == 0 {
		return
	}
	for _, t := range ts {
		tm[t.ID] = t
	}
	return
}
