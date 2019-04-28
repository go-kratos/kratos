package service

import (
	"context"
	"html/template"
	"time"

	"go-common/app/interface/main/space/model"
	arcmdl "go-common/app/service/main/archive/api"
	filmdl "go-common/app/service/main/filter/model/rpc"
	upmdl "go-common/app/service/main/up/api/v1"
	"go-common/library/ecode"
	"go-common/library/log"

	"go-common/library/sync/errgroup.v2"
)

const (
	_reasonWarnLevel  = int8(20)
	_reasonErrLevel   = int8(30)
	_checkTypeChannel = "channel"
)

var (
	_emptyArchiveReason = make([]*model.ArchiveReason, 0)
	_emptySearchVList   = make([]*model.SearchVList, 0)
)

// UpArcStat get up all article stat.
func (s *Service) UpArcStat(c context.Context, mid int64) (data *model.UpArcStat, err error) {
	addCache := true
	if data, err = s.dao.UpArcCache(c, mid); err != nil {
		addCache = false
	} else if data != nil {
		return
	}
	dt := time.Now().AddDate(0, 0, -1).Add(-12 * time.Hour).Format("20060102")
	if data, err = s.dao.UpArcStat(c, mid, dt); data != nil && addCache {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetUpArcCache(c, mid, data)
		})
	}
	return
}

// TopArc get top archive.
func (s *Service) TopArc(c context.Context, mid, vmid int64) (res *model.ArchiveReason, err error) {
	var (
		topArc   *model.AidReason
		arcReply *arcmdl.ArcReply
	)
	if topArc, err = s.dao.TopArc(c, vmid); err != nil {
		return
	}
	if topArc == nil || topArc.Aid == 0 {
		err = ecode.SpaceNoTopArc
		return
	}
	if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: topArc.Aid}); err != nil {
		log.Error("TopArc s.arcClient.Arc(%d) error(%v)", topArc.Aid, err)
		return
	}
	arc := arcReply.Arc
	if mid != vmid && !arc.IsNormal() {
		err = ecode.SpaceNoTopArc
		return
	}
	if arc.Access >= 10000 {
		arc.Stat.View = -1
	}
	res = &model.ArchiveReason{Arc: arc, Reason: template.HTMLEscapeString(topArc.Reason)}
	return
}

// SetTopArc set top archive.
func (s *Service) SetTopArc(c context.Context, mid, aid int64, reason string) (err error) {
	var (
		arcReply *arcmdl.ArcReply
		filRes   *filmdl.FilterRes
		topArc   *model.AidReason
	)
	if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil || arcReply.Arc == nil {
		log.Error("SetTopArc s.arcClient.Arc(%d) error(%v)", aid, err)
		return
	}
	arc := arcReply.Arc
	if !arc.IsNormal() {
		err = ecode.SpaceFakeAid
		return
	}
	if arc.Author.Mid != mid {
		err = ecode.SpaceNotAuthor
		return
	}
	if reason != "" {
		if filRes, err = s.filter.FilterArea(c, &filmdl.ArgFilter{Area: "common", Message: reason}); err != nil || filRes == nil {
			log.Error("SetTopArc s.filter.FilterArea(%s) error(%v)", reason, err)
			return
		}
		if filRes.Level >= _reasonErrLevel {
			err = ecode.SpaceTextBanned
			return
		}
		if filRes.Level == _reasonWarnLevel {
			reason = filRes.Result
		}
	}
	if topArc, err = s.dao.TopArc(c, mid); err != nil {
		return
	}
	if topArc != nil && aid == topArc.Aid && reason == topArc.Reason {
		err = ecode.NotModified
		return
	}
	if err = s.dao.AddTopArc(c, mid, aid, reason); err == nil {
		s.dao.AddCacheTopArc(c, mid, &model.AidReason{Aid: aid, Reason: reason})
	}
	return
}

// DelTopArc delete top archive.
func (s *Service) DelTopArc(c context.Context, mid int64) (err error) {
	var topArc *model.AidReason
	if topArc, err = s.dao.TopArc(c, mid); err != nil {
		return
	}
	if topArc == nil {
		err = ecode.RequestErr
		return
	}
	if err = s.dao.DelTopArc(c, mid); err == nil {
		s.dao.AddCacheTopArc(c, mid, &model.AidReason{Aid: -1})
	}
	return
}

// Masterpiece get masterpiece.
func (s *Service) Masterpiece(c context.Context, mid, vmid int64) (res []*model.ArchiveReason, err error) {
	var (
		mps       *model.AidReasons
		arcsReply *arcmdl.ArcsReply
		aids      []int64
	)
	if mps, err = s.dao.Masterpiece(c, vmid); err != nil {
		return
	}
	if mps == nil || len(mps.List) == 0 {
		res = _emptyArchiveReason
		return
	}
	for _, v := range mps.List {
		aids = append(aids, v.Aid)
	}
	if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
		log.Error("Masterpiece s.arcClient.Arcs(%v) error(%v)", aids, err)
		return
	}
	for _, v := range mps.List {
		if arc, ok := arcsReply.Arcs[v.Aid]; ok && arc != nil {
			if !arc.IsNormal() && mid != vmid {
				continue
			}
			if arc.Access >= 10000 {
				arc.Stat.View = -1
			}
			res = append(res, &model.ArchiveReason{Arc: arc, Reason: template.HTMLEscapeString(v.Reason)})
		}
	}
	if len(res) == 0 {
		res = _emptyArchiveReason
	}
	return
}

// AddMasterpiece add masterpiece.
func (s *Service) AddMasterpiece(c context.Context, mid, aid int64, reason string) (err error) {
	var (
		mps      *model.AidReasons
		arcReply *arcmdl.ArcReply
		filRes   *filmdl.FilterRes
	)
	if mps, err = s.dao.Masterpiece(c, mid); err != nil {
		return
	}
	if mps == nil {
		mps = &model.AidReasons{}
	}
	mpLen := len(mps.List)
	if mpLen >= s.c.Rule.MaxMpLimit {
		err = ecode.SpaceMpMaxCount
		return
	}
	if mpLen > 0 {
		for _, v := range mps.List {
			if v.Aid == aid {
				err = ecode.SpaceMpExist
				return
			}
		}
	}
	if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil || arcReply.Arc == nil {
		log.Error("AddMasterpiece s.arcClient.Arc(%d) error(%v)", aid, err)
		return
	}
	arc := arcReply.Arc
	if !arc.IsNormal() {
		err = ecode.SpaceFakeAid
		return
	}
	if arc.Author.Mid != mid {
		err = ecode.SpaceNotAuthor
		return
	}
	if reason != "" {
		if filRes, err = s.filter.FilterArea(c, &filmdl.ArgFilter{Area: "common", Message: reason}); err != nil || filRes == nil {
			log.Error("SetTopArc s.filter.FilterArea(%s) error(%v)", reason, err)
			return
		}
		if filRes.Level >= _reasonErrLevel {
			err = ecode.SpaceTextBanned
			return
		}
		if filRes.Level == _reasonWarnLevel {
			reason = filRes.Result
		}
	}
	if err = s.dao.AddMasterpiece(c, mid, aid, reason); err == nil {
		mps.List = append(mps.List, &model.AidReason{Aid: aid, Reason: reason})
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddCacheMasterpiece(c, mid, mps)
		})
	}
	return
}

// EditMasterpiece edit masterpiece.
func (s *Service) EditMasterpiece(c context.Context, mid, preAid, aid int64, reason string) (err error) {
	var (
		mps      *model.AidReasons
		arcReply *arcmdl.ArcReply
		filRes   *filmdl.FilterRes
		preCheck bool
	)
	if mps, err = s.dao.Masterpiece(c, mid); err != nil {
		return
	}
	if mps == nil || len(mps.List) == 0 {
		err = ecode.SpaceMpNoArc
		return
	}
	for _, v := range mps.List {
		if v.Aid == preAid {
			preCheck = true
		}
		if v.Aid == aid {
			err = ecode.SpaceMpExist
			return
		}
	}
	if !preCheck {
		err = ecode.SpaceMpNoArc
		return
	}
	if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: aid}); err != nil || arcReply.Arc == nil {
		log.Error("AddMasterpiece s.arcClient.Arc(%d) error(%v)", aid, err)
		return
	}
	arc := arcReply.Arc
	if !arc.IsNormal() {
		err = ecode.SpaceFakeAid
		return
	}
	if arc.Author.Mid != mid {
		err = ecode.SpaceNotAuthor
		return
	}
	if reason != "" {
		if filRes, err = s.filter.FilterArea(c, &filmdl.ArgFilter{Area: "common", Message: reason}); err != nil || filRes == nil {
			log.Error("SetTopArc s.filter.FilterArea(%s) error(%v)", reason, err)
			return
		}
		if filRes.Level >= _reasonErrLevel {
			err = ecode.SpaceTextBanned
			return
		}
		if filRes.Level == _reasonWarnLevel {
			reason = filRes.Result
		}
	}
	if err = s.dao.EditMasterpiece(c, mid, aid, preAid, reason); err == nil {
		newAidReasons := &model.AidReasons{}
		for _, v := range mps.List {
			if v.Aid == preAid {
				newAidReasons.List = append(newAidReasons.List, &model.AidReason{Aid: aid, Reason: reason})
			} else {
				newAidReasons.List = append(newAidReasons.List, v)
			}
		}
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddCacheMasterpiece(c, mid, newAidReasons)
		})
	}
	return
}

// CancelMasterpiece delete masterpiece.
func (s *Service) CancelMasterpiece(c context.Context, mid, aid int64) (err error) {
	var (
		mps        *model.AidReasons
		existCheck bool
	)
	if mps, err = s.dao.Masterpiece(c, mid); err != nil {
		return
	}
	if mps == nil || len(mps.List) == 0 {
		err = ecode.SpaceMpNoArc
		return
	}
	for _, v := range mps.List {
		if v.Aid == aid {
			existCheck = true
			break
		}
	}
	if !existCheck {
		err = ecode.SpaceMpNoArc
		return
	}
	if err = s.dao.DelMasterpiece(c, mid, aid); err == nil {
		newAidReasons := &model.AidReasons{}
		for _, v := range mps.List {
			if v.Aid == aid {
				continue
			}
			newAidReasons.List = append(newAidReasons.List, v)
		}
		if len(newAidReasons.List) == 0 {
			newAidReasons.List = append(newAidReasons.List, &model.AidReason{Aid: -1})
		}
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddCacheMasterpiece(c, mid, newAidReasons)
		})
	}
	return
}

// UpArcs get upload archive .
func (s *Service) UpArcs(c context.Context, mid int64, pn, ps int32) (res *model.UpArc, err error) {
	res = &model.UpArc{List: []*model.ArcItem{}}
	group := errgroup.WithContext(c)
	group.Go(func(ctx context.Context) error {
		if upCount, e := s.upClient.UpCount(ctx, &upmdl.UpCountReq{Mid: mid}); e != nil {
			log.Error("UpArcs s.upClient.UpCount mid(%d) error(%v)", mid, e)
		} else {
			res.Count = upCount.Count
		}
		return nil
	})
	group.Go(func(ctx context.Context) error {
		if reply, e := s.upClient.UpArcs(ctx, &upmdl.UpArcsReq{Mid: mid, Pn: pn, Ps: ps}); e != nil {
			log.Error("UpArcs s.upClient.UpArcs mid(%d) error(%v)", mid, err)
		} else if len(reply.Archives) > 0 {
			res.List = make([]*model.ArcItem, 0, len(reply.Archives))
			for _, v := range reply.Archives {
				si := &model.ArcItem{}
				si.FromArc(v)
				res.List = append(res.List, si)
			}
		}
		return nil
	})
	if e := group.Wait(); e != nil {
		log.Error("UpArcs group.Wait mid(%d) error(%v)", mid, e)
	}
	return
}

// ArcSearch get archive from search.
func (s *Service) ArcSearch(c context.Context, mid int64, arg *model.SearchArg) (data *model.SearchRes, total int, err error) {
	if data, total, err = s.dao.ArcSearchList(c, arg); err != nil {
		return
	}
	if len(data.VList) == 0 {
		data.VList = _emptySearchVList
		return
	}
	checkAids := make(map[int64]int64)
	if arg.CheckType == _checkTypeChannel {
		if mid == 0 {
			err = ecode.RequestErr
			return
		}
		var chArcs []*model.ChannelArc
		if chArcs, err = s.dao.ChannelVideos(c, mid, arg.CheckID, false); err != nil {
			err = nil
		} else {
			for _, chArc := range chArcs {
				checkAids[chArc.Aid] = chArc.Aid
			}
		}
	}
	vlist := make([]*model.SearchVList, 0)
	for _, v := range data.VList {
		if v.HideClick {
			v.Play = "--"
		}
		if _, ok := checkAids[v.Aid]; !ok {
			vlist = append(vlist, v)
		}
	}
	data.VList = vlist
	return
}
