package service

import (
	"context"

	tag "go-common/app/interface/main/tag/model"
	"go-common/app/interface/main/web/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// TagAids gets avids by tagID from bigdata or backup cache,
// and updates the cache after getting bigdata's data.
func (s *Service) TagAids(c context.Context, tagID int64, pn, ps int) (total int, arcs []*arcmdl.Arc, err error) {
	defer func() {
		if len(arcs) == 0 {
			arcs = _emptyArc
		}
	}()
	if err = s.checkTag(c, tagID); err != nil {
		err = nil
		return
	}
	return s.tagArcs(c, tagID, pn, ps)
}

func (s *Service) tagArcs(c context.Context, tagID int64, pn, ps int) (total int, arcs []*arcmdl.Arc, err error) {
	var (
		start         = (pn - 1) * ps
		end           = start + ps - 1
		aids, allAids []int64
	)
	if allAids, err = s.dao.TagAids(c, tagID); err != nil {
		log.Error("s.dao.TagAids(%d) error(%v)", tagID, err)
		if allAids, err = s.dao.TagAidsBakCache(c, tagID); err != nil {
			log.Error("s.dao.TagAidsBakCache(%d) error(%v)", tagID, err)
			return
		}
	} else if len(allAids) > 0 {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetTagAidsBakCache(c, tagID, allAids)
		})
	}
	total = len(allAids)
	if total < start {
		err = ecode.NothingFound
		return
	}
	if total > end {
		aids = allAids[start : end+1]
	} else {
		aids = allAids[start:]
	}
	arcs, err = s.archives(c, aids)
	return
}

func (s *Service) archives(c context.Context, aids []int64) (data []*arcmdl.Arc, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	var (
		arg = &arcmdl.ArcsRequest{Aids: aids}
		res *arcmdl.ArcsReply
	)
	archivesArgLog("TagAids", aids)
	if res, err = s.arcClient.Arcs(c, arg); err != nil {
		log.Error("arcrpc.Archives3(%v,%s) error(%v)", aids, ip, err)
		return
	}
	for _, aid := range aids {
		if _, ok := res.Arcs[aid]; !ok {
			continue
		}
		data = append(data, res.Arcs[aid])
	}
	fmtArcs(data)
	return
}

func fmtArcs(arcs []*arcmdl.Arc) {
	for _, v := range arcs {
		if v.Access >= 10000 {
			v.Stat.View = -1
		}
	}
}

func (s *Service) checkTag(c context.Context, tid int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	var (
		t   *tag.Tag
		arg = &tag.ArgID{
			ID:     tid,
			RealIP: ip,
		}
	)
	if t, err = s.tag.InfoByID(c, arg); err != nil {
		log.Error("tagState(%d) error(%v)", tid, err)
		err = nil
		return
	}
	if t.State == model.TagStateDeleted || t.State == model.TagStateBlocked {
		err = ecode.TagIsSealing
	}
	return
}

// TagDetail group web tag data.
func (s *Service) TagDetail(c context.Context, tagID int64, ps int) (data *model.TagDetail, err error) {
	var tagInfo *tag.TagTop
	if tagInfo, err = s.tag.TagTop(c, &tag.ReqTagTop{Tid: tagID}); err != nil {
		return
	}
	data = &model.TagDetail{TagTop: tagInfo}
	data.Total, data.List, _ = s.tagArcs(c, tagID, _samplePn, ps)
	if len(data.List) == 0 {
		data.List = _emptyArc
	}
	if len(data.Similars) == 0 {
		data.Similars = make([]*tag.SimilarTag, 0)
	}
	return
}
