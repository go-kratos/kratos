package service

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"go-common/app/interface/main/dm2/model"
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

// gzflate flate 压缩
func (s *Service) gzflate(input []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w, err := flate.NewWriter(buf, 4)
	if err != nil {
		return nil, err
	}
	if _, err = w.Write(input); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// Gzdeflate deflate 解码
func (s *Service) Gzdeflate(in []byte) (out []byte, err error) {
	if len(in) == 0 {
		return
	}
	out, err = ioutil.ReadAll(flate.NewReader(bytes.NewReader(in)))
	return
}

func (s *Service) dmsSeg(c context.Context, tp int32, oid int64, missed []int64) (dms []*model.DM, err error) {
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

func (s *Service) dmSegXML(c context.Context, aid int64, sub *model.Subject, seg *model.Segment) (res []byte, err error) {
	var (
		cache                   = true
		buf                     = new(bytes.Buffer)
		realname                = s.isRealname(c, sub.Pid, sub.Oid)
		dms                     []*model.DM
		dmids, normalIds, spIds []int64
	)
	if realname {
		buf.WriteString(seg.ToXMLHeader(sub.Oid, sub.State, 1))
	} else {
		buf.WriteString(seg.ToXMLHeader(sub.Oid, sub.State, 0))
	}
	defer func() {
		if err == nil {
			buf.WriteString(`</i>`)
			res, err = s.gzip(buf.Bytes())
		}
	}()

	if normalIds, err = s.dmNormalIds(c, sub.Type, sub.Oid, seg.Cnt, seg.Num, seg.Start, seg.End, 2*sub.Maxlimit); err != nil {
		return
	}
	dmids = append(dmids, normalIds...)
	if sub.Childpool > 0 {
		if spIds, err = s.dmSegSubtitlesIds(c, sub.Type, sub.Oid, seg.Start, seg.End, 2*sub.Maxlimit); err != nil {
			return
		}
		dmids = append(dmids, spIds...)
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
		if dms, err = s.dmsSeg(c, sub.Type, sub.Oid, missed); err != nil {
			return
		}
		for _, dm := range dms {
			buf.WriteString(dm.ToXMLSeg(realname))
		}
	}
	if cache && len(dms) > 0 {
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.AddIdxContentCache(ctx, sub.Type, sub.Oid, dms, realname)
		})
	}
	s.cache.Do(c, func(ctx context.Context) {
		var (
			bs  []byte
			err error
		)
		dmSeg := &model.ActionFlushDMSeg{
			Type:  sub.Type,
			Oid:   sub.Oid,
			Force: false,
			Page: &model.Page{
				Num:   seg.Num,
				Total: seg.Cnt,
				Size:  model.DefaultPageSize,
			},
		}
		if bs, err = json.Marshal(dmSeg); err != nil {
			return
		}
		s.dao.SendAction(ctx, fmt.Sprint(sub.Oid), &model.Action{
			Action: model.ActFlushDMSeg,
			Data:   bs,
		})
	})
	return
}

// DMSeg return dm content.
func (s *Service) DMSeg(c context.Context, tp, plat int32, mid, aid, oid, ps int64) (res []byte, err error) {
	var (
		sub  = &model.Subject{}
		flag = model.DefaultFlag
	)
	seg, err := s.segmentInfo(c, tp, aid, oid, ps)
	if err != nil {
		return
	}
	if sub, err = s.subject(c, tp, oid); err != nil {
		return
	}
	if sub.State == model.SubStateClosed {
		xml := []byte(seg.ToXMLHeader(sub.Oid, sub.State, 0) + `</i>`)
		if xml, err = s.gzip(xml); err != nil {
			return
		}
		res = model.Encode(flag, xml)
		return
	}
	// get from local cache first
	xml, ok := s.localCache[keySeg(tp, oid, seg.Cnt, seg.Num)]
	if ok {
		res = model.Encode(flag, xml)
		return
	}
	// NOTE 将视频弹幕上限调整为 2*maxlimit条
	data, err := s.dao.RecFlag(c, mid, aid, oid, 2*sub.Maxlimit, seg.Start, seg.End, plat)
	if err == nil {
		flag = data
	}
	// get from remote cache or database
	if xml, err = s.singleGenSegXML(c, aid, sub, seg); err != nil {
		return
	}
	res = model.Encode(flag, xml)
	return
}

func (s *Service) singleGenSegXML(c context.Context, aid int64, sub *model.Subject, seg *model.Segment) (xml []byte, err error) {
	key := keySeg(sub.Type, sub.Oid, seg.Cnt, seg.Num)
	v, err, _ := s.singleGroup.Do(key, func() (res interface{}, err error) {
		data, err := s.dao.XMLSegCache(c, sub.Type, sub.Oid, seg.Cnt, seg.Num)
		if err != nil {
			return
		}
		if len(data) > 0 {
			res = data
			return
		}
		if data, err = s.dmSegXML(c, aid, sub, seg); err != nil {
			return
		}
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.SetXMLSegCache(ctx, sub.Type, sub.Oid, seg.Cnt, seg.Num, data)
		})
		return data, err
	})
	if err != nil {
		return
	}
	xml = v.([]byte)
	return
}

// segmentInfo get segment info of oid.
func (s *Service) segmentInfo(c context.Context, tp int32, aid, oid, ps int64) (seg *model.Segment, err error) {
	var duration int64
	data, ok := s.localCache[keyDuration(tp, oid)]
	if ok {
		duration, err = strconv.ParseInt(string(data), 10, 64)
	} else {
		duration, err = s.videoDuration(c, aid, oid)
	}
	if err != nil {
		return
	}
	if duration != 0 && ps >= duration {
		log.Warn("oid:%d ps:%d larger than duration:%d", oid, ps, duration)
		err = ecode.NotModified
		return
	}
	seg = model.SegmentInfo(ps, duration)
	return
}

func (s *Service) dmNormalIds(c context.Context, tp int32, oid int64, cnt, n, ps, pe, limit int64) (dmids []int64, err error) {
	dmids, err = s.dao.DMIDCache(c, tp, oid, cnt, n, limit)
	if err != nil || len(dmids) == 0 {
		if dmids, err = s.dao.DMIDs(c, tp, oid, ps, pe, limit, model.PoolNormal); err != nil {
			return
		}
	}
	return
}

//dmSegSubtitlesIds dm Subtitles
func (s *Service) dmSegSubtitlesIds(c context.Context, tp int32, oid int64, ps, pe, limit int64) (dmids []int64, err error) {
	dmids, err = s.dao.DMIDSubtitlesCache(c, tp, oid, ps, pe, limit)
	if err != nil || len(dmids) == 0 {
		if dmids, err = s.dao.DMIDs(c, tp, oid, ps, pe, limit, model.PoolSubtitle); err != nil || len(dmids) == 0 {
			return
		}
	}
	return
}
