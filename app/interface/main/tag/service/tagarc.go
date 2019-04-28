package service

import (
	"context"
	"strconv"

	"go-common/app/interface/main/tag/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var (
	_emptyArcs    = []*api.Arc{}
	_emptySimilar = []*model.SimilarTag{}
)

// NewArcs tid new arcs.
func (s *Service) NewArcs(c context.Context, tid int64, ps, pn int) (as []*api.Arc, count int, err error) {
	var (
		aids        []int64
		invalidArcs []*api.Arc
		start       = (pn - 1) * ps
		end         = start + ps - 1
	)
	if end > s.c.Tag.MaxArcsLimit {
		end = s.c.Tag.MaxArcsLimit
	}
	if aids, count, err = s.newArcs(c, tid, start, end); err != nil {
		return
	}
	if len(aids) == 0 {
		as = _emptyArcs
		return
	}
	var (
		res map[int64]*api.Arc
		arg = &archive.ArgAids2{Aids: aids, RealIP: metadata.String(c, metadata.RemoteIP)}
	)
	if res, err = s.arcRPC.Archives3(c, arg); err != nil {
		return
	}
	for _, aid := range aids {
		k, ok := res[aid]
		if !ok {
			log.Warn("aid: %d do not received from arcRPC.Archives3.", aid)
			continue
		}
		if k.IsNormal() {
			as = append(as, k)
		} else {
			invalidArcs = append(invalidArcs, k)
		}
	}
	if len(as) == 0 {
		as = _emptyArcs
	}
	if len(invalidArcs) > 0 {
		s.invalidArcCh.Do(c, func(ctx context.Context) {
			s.delInvalidArc(ctx, tid, 0, invalidArcs)
		})
	}
	return
}

// RegionNewArcs rid-tid newArcs.
func (s *Service) RegionNewArcs(c context.Context, rid int32, tid int64, tp int8, ps, pn int) (as []*api.Arc, count int, err error) {
	var (
		aids        []int64
		invalidArcs []*api.Arc
		mArchive    map[int64]*api.Arc
		start       = (pn - 1) * ps
		end         = start + ps - 1
	)
	if tp == archive.CopyrightOriginal {
		if aids, count, err = s.dao.OriginRegionNewArcsCache(c, rid, tid, start, end); err != nil {
			log.Error("s.ta.OriginRankedArcsCache(tid:%d, start:%d, end:%d) error(%v)", tid, start, end, err)
			return
		}
	} else {
		if aids, count, err = s.dao.RegionNewArcsCache(c, rid, tid, start, end); err != nil {
			log.Error("s.ta.RankedArcsCache(tid:%d, start:%d, end:%d) error(%v)", tid, start, end, err)
			return
		}
	}
	if len(aids) == 0 {
		as = _emptyArcs
		return
	}
	var arg = &archive.ArgAids2{Aids: aids, RealIP: metadata.String(c, metadata.RemoteIP)}
	if mArchive, err = s.arcRPC.Archives3(c, arg); err != nil {
		log.Error("s.arcRPC.Archives3(%v) error(%v)", aids, err)
		return
	}
	as = make([]*api.Arc, 0, len(aids))
	for _, aid := range aids {
		k, ok := mArchive[aid]
		if !ok {
			log.Warn("aid: %d do not received from arcRPC.Archives3.", aid)
			continue
		}
		if k.IsNormal() {
			as = append(as, k)
		} else {
			invalidArcs = append(invalidArcs, k)
		}
	}
	if len(as) == 0 {
		as = _emptyArcs
	}
	if len(invalidArcs) > 0 {
		s.invalidArcCh.Do(c, func(ctx context.Context) {
			s.delInvalidArc(ctx, tid, rid, invalidArcs)
		})
	}
	return
}

// DetailRankArc .
func (s *Service) DetailRankArc(c context.Context, tid, prid int64, pn, ps int) (as []*api.Arc, count int, err error) {
	var (
		aids  []int64
		mas   map[int64]*api.Arc
		start = (pn - 1) * ps
		end   = start + ps - 1
	)
	if aids, count, err = s.tagPridAids(c, tid, prid, start, end); err != nil {
		return
	}
	if len(aids) == 0 {
		as = _emptyArcs
		return
	}
	if len(aids) == 1 && aids[0] == -1 {
		count = 0
		as = _emptyArcs
		return
	}
	arg := &archive.ArgAids2{Aids: aids, RealIP: metadata.String(c, metadata.RemoteIP)}
	if mas, err = s.arcRPC.Archives3(c, arg); err != nil {
		log.Error("s.arcRPC.Archives3(%v) error(%v)", aids, err)
		return
	}
	as = make([]*api.Arc, 0, len(aids))
	for _, aid := range aids {
		if a, ok := mas[aid]; ok && a.IsNormal() {
			as = append(as, a)
		}
	}
	if len(as) == 0 {
		as = _emptyArcs
	}
	return
}

// Detail .
func (s *Service) Detail(c context.Context, tid, mid int64, pn, ps int) (detail *model.Detail, err error) {
	var tag *model.Tag
	if tag, err = s.info(c, mid, tid); err != nil {
		return
	}
	if tag == nil {
		return nil, ecode.TagNotExist
	}
	if tag.State == model.TagStateDel || tag.State == model.TagStateHide {
		return nil, ecode.TagIsSealing
	}
	detail = &model.Detail{Info: tag}
	var (
		tids        []int64
		tags        []*model.Tag
		similarTags []*model.SimilarTag
	)
	tids, _ = s.similarsTids(c, tid)
	if len(tids) > 0 {
		if tags, err = s.infos(c, 0, tids); err != nil {
			return
		}
		for _, tag := range tags {
			similarTag := &model.SimilarTag{
				Tid:   tag.ID,
				Tname: tag.Name,
			}
			similarTags = append(similarTags, similarTag)
		}
		detail.Similar = similarTags
	} else {
		detail.Similar = _emptySimilar
	}
	// get newArc
	var (
		start = (pn - 1) * ps
		end   = start + ps - 1
		aids  []int64
		count int
	)
	if aids, count, err = s.newArcs(c, tid, start, end); err != nil {
		return
	}
	if len(aids) > 0 {
		var (
			mas map[int64]*api.Arc
			arg = &archive.ArgAids2{Aids: aids, RealIP: metadata.String(c, metadata.RemoteIP)}
		)
		if mas, err = s.arcRPC.Archives3(c, arg); err != nil {
			log.Error("s.arcRPC.Archives3(%v) error(%v)", aids, err)
			return
		}
		var arcs = make([]*api.Arc, 0, len(mas))
		for _, aid := range aids {
			if a, ok := mas[aid]; ok && a.IsNormal() {
				arcs = append(arcs, a)
			}
		}
		detail.News.Archives = arcs
		detail.News.Count = count
		return
	}
	detail.News.Archives = _emptyArcs
	return
}

func (s *Service) delInvalidArc(c context.Context, tid int64, rid int32, invalidArcs []*api.Arc) (err error) {
	var str string
	for _, a := range invalidArcs {
		str = strconv.FormatInt(a.Aid, 10) + " "
	}
	s.dao.DeleteNewArcCache(c, tid, str)
	s.dao.DeleteRegionNewArcsCache(c, tid, rid, invalidArcs)
	return
}
