package service

import (
	"context"
	"time"

	"go-common/app/interface/main/tag/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	rpcModel "go-common/app/service/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"
)

// UpArcBind .
func (s *Service) UpArcBind(c context.Context, aid, mid int64, tNames, regionNames []string, now time.Time) (err error) {
	var (
		addTids      []int64
		delTids      []int64
		exists       []string
		addNames     []string
		delNames     []string
		defaultNames []string
		allNames     []string
		tMap         map[string]int64
		arc          *api.Arc
	)
	if exists, err = s.tArcNamesByAid(c, aid); err != nil {
		return
	}
	allNames, addNames, delNames, defaultNames = s.getChangeTag(exists, tNames, regionNames)
	if len(addNames) == 0 && len(delNames) == 0 {
		if len(defaultNames) == 0 {
			return
		}
		var defaultTids []int64
		if _, tMap, err = s.addNewTags(c, mid, defaultNames, now); err != nil {
			return
		}
		for _, name := range defaultNames {
			if tid := tMap[name]; tid > 0 {
				defaultTids = append(defaultTids, tid)
			}
		}
		if len(defaultTids) > 0 {
			s.defaultUpBind(c, aid, mid, defaultTids, rpcModel.ResTypeArchive)
		}
		return
	}
	if _, tMap, err = s.addNewTags(c, mid, allNames, now); err != nil {
		return
	}
	var tids []int64
	for _, v := range tNames {
		if k, ok := tMap[v]; ok {
			tids = append(tids, k)
		}
	}
	if len(tids) > 0 {
		if err = s.platformUpBind(c, aid, mid, tids, rpcModel.ResTypeArchive); err != nil {
			return
		}
	}
	if len(defaultNames) != 0 {
		var defaultTids []int64
		for _, name := range defaultNames {
			if tid := tMap[name]; tid > 0 {
				defaultTids = append(defaultTids, tid)
			}
		}
		if len(defaultTids) > 0 {
			s.defaultUpBind(c, aid, mid, defaultTids, rpcModel.ResTypeArchive)
		}
	}
	for _, tn := range addNames {
		if id, ok := tMap[tn]; ok {
			addTids = append(addTids, id)
		}
	}
	for _, tn := range delNames {
		if id, ok := tMap[tn]; ok {
			delTids = append(delTids, id)
		}
	}
	arg := &archive.ArgAid2{Aid: aid, RealIP: metadata.String(c, metadata.RemoteIP)}
	if arc, err = s.arcRPC.Archive3(c, arg); err != nil {
		if ecode.Cause(err).Code() == ecode.NothingFound.Code() {
			err = nil
		}
		return
	}
	if arc != nil {
		s.auditBindTag(c, aid, arc, addTids)
		s.remTagsArc(c, aid, arc, delTids)
	}
	return
}

// ArcAdminBind .
func (s *Service) ArcAdminBind(c context.Context, aid, mid int64, tNames, regionNames []string, now time.Time) (err error) {
	var (
		addTids      []int64
		delTids      []int64
		existNames   []string
		addNames     []string
		delNames     []string
		defaultNames []string
		allNames     []string
		defaultTids  []int64
		tMap         map[string]int64
		arc          *api.Arc
	)
	if existNames, err = s.tArcNamesByAid(c, aid); err != nil {
		return
	}
	allNames, addNames, delNames, defaultNames = s.getChangeTag(existNames, tNames, regionNames)
	if len(addNames) == 0 && len(delNames) == 0 {
		if len(defaultNames) == 0 {
			return
		}
		if _, tMap, err = s.addNewTags(c, mid, defaultNames, now); err != nil {
			return
		}
		for _, name := range defaultNames {
			if tid := tMap[name]; tid > 0 {
				defaultTids = append(defaultTids, tid)
			}
		}
		if len(defaultTids) > 0 {
			s.defaultAdminBind(c, aid, mid, defaultTids, rpcModel.ResTypeArchive)
		}
		return
	}
	if _, tMap, err = s.addNewTags(c, mid, allNames, now); err != nil {
		return
	}
	var tids []int64
	for _, v := range tNames {
		if k, ok := tMap[v]; ok {
			tids = append(tids, k)
		}
	}
	if err = s.platformAdminBind(c, aid, mid, tids, rpcModel.ResTypeArchive); err != nil {
		return
	}
	for _, name := range defaultNames {
		if tid := tMap[name]; tid > 0 {
			defaultTids = append(defaultTids, tid)
		}
	}
	s.defaultAdminBind(c, aid, mid, defaultTids, rpcModel.ResTypeArchive)
	for _, tn := range addNames {
		if id, ok := tMap[tn]; ok {
			addTids = append(addTids, id)
		}
	}
	for _, tn := range delNames {
		if id, ok := tMap[tn]; ok {
			delTids = append(delTids, id)
		}
	}
	arg := &archive.ArgAid2{Aid: aid, RealIP: metadata.String(c, metadata.RemoteIP)}
	if arc, err = s.arcRPC.Archive3(c, arg); err != nil {
		if ecode.Cause(err).Code() == ecode.NothingFound.Code() {
			err = nil
		}
		return
	}
	if arc != nil {
		s.auditBindTag(c, aid, arc, addTids)
		s.remTagsArc(c, aid, arc, delTids)
	}
	return
}

func (s *Service) tArcNamesByAid(c context.Context, aid int64) (tNames []string, err error) {
	var (
		tids []int64
		ts   []*model.Tag
		res  []*rpcModel.Resource
	)
	res, err = s.resTags(c, 0, aid, rpcModel.ResTypeArchive)
	if err != nil {
		return
	}
	for _, r := range res {
		tids = append(tids, r.Tid)
	}
	if len(tids) == 0 {
		return
	}
	if ts, err = s.infos(c, 0, tids); err != nil {
		if ecode.TagNotExist.Equal(err) {
			err = nil
		}
		return
	}
	for _, t := range ts {
		tNames = append(tNames, t.Name)
	}
	return
}
