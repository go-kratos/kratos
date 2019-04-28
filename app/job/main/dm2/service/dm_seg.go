package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"math"

	"go-common/app/job/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) gzip(input []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	zw := gzip.NewWriter(buf)
	if _, err := zw.Write(input); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Service) dmsByid(c context.Context, tp int32, oid int64, missed []int64) (dms []*model.DM, err error) {
	idxMap, special, err := s.dao.IndexsByid(c, tp, oid, missed)
	if err != nil || len(idxMap) == 0 {
		return
	}
	ctsMap, err := s.dao.Contents(c, oid, missed)
	if err != nil {
		return
	}
	ctsSpeMap := make(map[int64]*model.ContentSpecial)
	if len(special) > 0 {
		if ctsSpeMap, err = s.dao.ContentsSpecial(c, special); err != nil {
			return
		}
	}
	for _, content := range ctsMap {
		if idx, ok := idxMap[content.ID]; ok {
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
				Content:  content,
			}
			if idx.Pool == model.PoolSpecial {
				if _, ok = ctsSpeMap[dm.ID]; ok {
					dm.ContentSpe = ctsSpeMap[dm.ID]
				}
			}
			dms = append(dms, dm)
		}
	}
	return
}

func (s *Service) dmSeg(c context.Context, tp int32, oid, limit int64, childpool int32, p *model.Page) (res *model.DMSeg, err error) {
	var (
		ids   []int64
		cache = true
		dmids = make([]int64, 0, limit)
		elems = make([]*model.Elem, 0, limit)
		ps    = (p.Num - 1) * p.Size
		pe    = p.Num * p.Size
	)
	res = new(model.DMSeg)
	if ids, err = s.dmidsSeg(c, tp, oid, p.Total, p.Num, ps, pe, limit); err != nil {
		return
	}
	dmids = append(dmids, ids...)
	if childpool > 0 {
		if ids, err = s.dmidSubtitle(c, tp, oid, ps, pe, limit); err != nil {
			return
		}
		dmids = append(dmids, ids...)
	}
	if len(dmids) <= 0 {
		return
	}
	elemsCache, missed, err := s.dao.IdxContentCacheV2(c, tp, oid, dmids)
	if err != nil {
		missed = dmids
		cache = false
	} else {
		elems = append(elems, elemsCache...)
	}
	if len(missed) == 0 {
		res.Elems = elems
		return
	}
	dms, err := s.dmsByid(c, tp, oid, missed)
	if err != nil {
		return
	}
	for _, dm := range dms {
		if e := dm.ToElem(); e != nil {
			elems = append(elems, e)
		}
	}
	res.Elems = elems
	if cache && len(dms) > 0 {
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.AddIdxContentCaches(ctx, tp, oid, dms...)
		})
	}
	return
}

func (s *Service) dmidsSeg(c context.Context, tp int32, oid, total, num, ps, pe, limit int64) (dmids []int64, err error) {
	if dmids, err = s.dao.DMIDCache(c, tp, oid, total, num, limit); err != nil || len(dmids) == 0 {
		if dmids, err = s.dao.IndexsSegID(c, tp, oid, ps, pe, limit, model.PoolNormal); err != nil {
			return
		}
		if len(dmids) > 0 {
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.AddDMIDCache(ctx, tp, oid, total, num, dmids...)
			})
		}
	}
	return
}

func (s *Service) dmidSubtitle(c context.Context, tp int32, oid, ps, pe, limit int64) (dmids []int64, err error) {
	if dmids, err = s.dao.DMIDSubtitleCache(c, tp, oid, ps, pe, limit); err != nil || len(dmids) == 0 {
		var dms []*model.DM
		if dms, dmids, err = s.dao.IndexsSeg(c, tp, oid, ps, pe, limit, model.PoolSubtitle); err != nil {
			return
		}
		if len(dms) > 0 {
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.AddDMIDSubtitleCache(ctx, tp, oid, dms...)
			})
		}
	}
	return
}

// add flush dm segment action to flush channel.
func (s *Service) asyncAddFlushDMSeg(c context.Context, fc *model.FlushDMSeg) (err error) {
	select {
	case s.flushSegChan[fc.Oid%int64(s.routineSize)] <- fc:
	default:
		log.Warn("segment flush merge channel is full,fc:%+v page:%+v", fc, fc.Page)
	}
	return
}

func (s *Service) pageinfo(c context.Context, pid int64, dm *model.DM) (p *model.Page, err error) {
	duration, err := s.videoDuration(c, pid, dm.Oid)
	if err != nil {
		return
	}
	if duration != 0 {
		p = &model.Page{
			Num:   int64(math.Ceil(float64(dm.Progress) / float64(model.DefaultPageSize))),
			Size:  model.DefaultPageSize,
			Total: int64(math.Ceil(float64(duration) / float64(model.DefaultPageSize))),
		}
		if p.Num == 0 { // fix progress == 0
			p.Num = 1
		}
	} else { // duration not exist
		p = model.DefaultPage
	}
	// NOTE PoolSpecial store in the first segment
	if dm.Pool == model.PoolSpecial {
		p.Num = 1
	}
	return
}

func (s *Service) dmSegXML(c context.Context, sub *model.Subject, seg *model.Segment) (res []byte, err error) {
	var (
		cache                         = true
		buf                           = new(bytes.Buffer)
		dms                           []*model.DM
		dmids, normalIds, subtitleIds []int64
	)
	buf.WriteString(seg.ToXMLHeader(sub.Oid, sub.State, 0))
	defer func() {
		if err == nil {
			buf.WriteString(`</i>`)
			res, err = s.gzip(buf.Bytes())
		}
	}()

	if normalIds, err = s.dmidsSeg(c, sub.Type, sub.Oid, seg.Cnt, seg.Num, seg.Start, seg.End, 2*sub.Maxlimit); err != nil {
		return
	}
	dmids = append(dmids, normalIds...)
	if sub.Childpool > 0 {
		if subtitleIds, err = s.dmidSubtitle(c, sub.Type, sub.Oid, seg.Start, seg.End, 2*sub.Maxlimit); err != nil {
			return
		}
		dmids = append(dmids, subtitleIds...)
	}
	if len(dmids) <= 0 {
		return
	}
	content, missed, err := s.dao.IdxContentCache(c, sub.Type, sub.Oid, dmids)
	if err != nil {
		missed = dmids
		cache = false
	} else {
		buf.Write(content)
	}
	if len(missed) > 0 {
		if dms, err = s.dmsByid(c, sub.Type, sub.Oid, missed); err != nil {
			return
		}
		for _, dm := range dms {
			buf.WriteString(dm.ToXMLSeg())
		}
	}
	if cache && len(dms) > 0 {
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.AddIdxContentCaches(ctx, sub.Type, sub.Oid, dms...)
		})
	}
	return
}

// segmentInfo get segment info of oid.
func (s *Service) segmentInfo(c context.Context, aid, oid, ps int64, duration int64) (seg *model.Segment, err error) {
	if duration != 0 && ps >= duration {
		log.Warn("oid:%d ps:%d larger than duration:%d", oid, ps, duration)
		err = ecode.NotModified
		return
	}
	seg = model.SegmentInfo(ps, duration)
	return
}
