package service

import (
	"context"
	"sync"

	tagmdl "go-common/app/interface/main/tag/model"
	"go-common/app/interface/main/web/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_tagBlkSize = 50
	_tagArcType = 3
)

var (
	_emptyWxArc    = make([]*model.WxArchive, 0)
	_emptyWxArcTag = make([]*model.WxArcTag, 0)
)

// WxHot get wx hot archives.
func (s *Service) WxHot(c context.Context, pn, ps int) (res []*model.WxArchive, count int, err error) {
	var (
		addCache       = true
		aids           []int64
		wxArcs, allRes []*model.WxArchive
		arcs           map[int64]*model.WxArchive
	)
	defer func() {
		res, count = wxHotPage(allRes, pn, ps)
	}()
	if wxArcs, err = s.dao.WxHotCache(c); err != nil {
		err = nil
		addCache = false
	} else if len(wxArcs) > 0 {
		allRes = s.fmtWxArcs(c, wxArcs)
		return
	}
	if aids, err = s.dao.WxHot(c); err != nil {
		err = nil
	} else if len(aids) > s.c.Rule.MinWxHotCount {
		if arcs, err = s.archiveWithTag(c, aids); err != nil {
			err = nil
		} else if len(arcs) > 0 {
			for _, aid := range aids {
				if arc, ok := arcs[aid]; ok && arc != nil {
					allRes = append(allRes, arc)
				}
			}
			if addCache {
				s.cache.Do(c, func(c context.Context) {
					s.dao.SetWxHotCache(c, allRes)
				})
			}
			return
		}
	} else {
		log.Error("s.dao.WxHot len(%d)", len(allRes))
	}
	allRes, err = s.dao.WxHotBakCache(c)
	return
}

func wxHotPage(list []*model.WxArchive, pn, ps int) (res []*model.WxArchive, count int) {
	count = len(list)
	start := (pn - 1) * ps
	end := start + ps - 1
	if count == 0 || count < start {
		res = _emptyWxArc
		return
	}
	if count > end {
		res = list[start : end+1]
	} else {
		res = list[start:]
	}
	return
}

func (s *Service) archiveWithTag(c context.Context, aids []int64) (list map[int64]*model.WxArchive, err error) {
	var (
		arcErr, tagErr error
		arcsReply      *arcmdl.ArcsReply
		tags           map[int64][]*tagmdl.Tag
		mutex          = sync.Mutex{}
		ip             = metadata.String(c, metadata.RemoteIP)
	)
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if arcsReply, arcErr = s.arcClient.Arcs(errCtx, &arcmdl.ArcsRequest{Aids: aids}); arcErr != nil {
			log.Error("s.arcClient.Arcs(%d) error %v", aids, err)
			return arcErr
		}
		return nil
	})
	aidsLen := len(aids)
	tags = make(map[int64][]*tagmdl.Tag, aidsLen)
	for i := 0; i < aidsLen; i += _tagBlkSize {
		var partAids []int64
		if i+_tagBlkSize > aidsLen {
			partAids = aids[i:]
		} else {
			partAids = aids[i : i+_tagBlkSize]
		}
		group.Go(func() (err error) {
			var tmpRes map[int64][]*tagmdl.Tag
			arg := &tagmdl.ArgResTags{Oids: partAids, Type: _tagArcType, RealIP: ip}
			if tmpRes, tagErr = s.tag.ResTags(errCtx, arg); tagErr != nil {
				log.Error("s.tag.ResTag(%+v) error(%v)", arg, tagErr)
				return
			}
			mutex.Lock()
			for aid, tmpTags := range tmpRes {
				tags[aid] = tmpTags
			}
			mutex.Unlock()
			return nil
		})
	}
	if err = group.Wait(); err != nil {
		return
	}
	list = make(map[int64]*model.WxArchive, len(aids))
	for _, aid := range aids {
		if arc, ok := arcsReply.Arcs[aid]; ok && arc.IsNormal() {
			wxArc := new(model.WxArchive)
			wxArc.FromArchive(arc)
			if tag, ok := tags[aid]; ok {
				for _, v := range tag {
					wxArc.Tags = append(wxArc.Tags, &model.WxArcTag{ID: v.ID, Name: v.Name})
				}
			}
			if len(wxArc.Tags) == 0 {
				wxArc.Tags = _emptyWxArcTag
			}
			list[aid] = wxArc
		}
	}
	return
}

func (s *Service) fmtWxArcs(c context.Context, arcs []*model.WxArchive) (res []*model.WxArchive) {
	var (
		aids    []int64
		newArcs map[int64]*model.WxArchive
		err     error
	)
	for _, arc := range arcs {
		aids = append(aids, arc.Aid)
	}
	if newArcs, err = s.archiveWithTag(c, aids); err != nil {
		log.Error("fmtWxArcStat archiveWithTag aids(%+v) error(%v)", aids, err)
		return
	}
	for _, arc := range arcs {
		if newArc, ok := newArcs[arc.Aid]; ok {
			res = append(res, newArc)
		} else {
			res = append(res, arc)
		}
	}
	return
}
