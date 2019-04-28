package service

import (
	"bytes"
	"compress/flate"
	"context"
	"fmt"
	"math"
	"sort"

	"go-common/app/job/main/dm2/model"
	arcMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

// Gzflate flate 压缩
func (s *Service) gzflate(in []byte, level int) (out []byte, err error) {
	if len(in) == 0 {
		return
	}
	buf := new(bytes.Buffer)
	w, err := flate.NewWriter(buf, level)
	if err != nil {
		return
	}
	if _, err = w.Write(in); err != nil {
		return
	}
	if err = w.Close(); err != nil {
		return
	}
	out = buf.Bytes()
	return
}

func (s *Service) dmsCache(c context.Context, tp int32, oid, maxlimit int64) (dms []*model.DM, err error) {
	ok, err := s.dao.ExpireDMCache(c, tp, oid)
	if err != nil || !ok {
		return
	}
	values, err := s.dao.DMCache(c, tp, oid)
	if err != nil || len(values) == 0 {
		return
	}
	var (
		start, trimCnt           int
		normal, protect, special []*model.DM
	)
	for _, value := range values {
		dm := &model.DM{}
		if err = dm.Unmarshal(value); err != nil {
			log.Error("proto.Unmarshal(%s) error(%v)", value, err)
			return
		}
		if dm.Pool == model.PoolNormal {
			if dm.AttrVal(model.AttrProtect) == model.AttrYes {
				protect = append(protect, dm)
			} else {
				normal = append(normal, dm)
			}
		} else {
			special = append(special, dm)
		}
	}
	// 保护弹幕
	if start = len(protect) - int(maxlimit); start > 0 { // 只保留maxlimit条保护弹幕
		trimCnt += start
		protect = protect[start:]
	}
	dms = append(dms, protect...)
	// 普通弹幕
	if start = len(normal) + len(protect) - int(maxlimit); start > 0 { // 保护弹幕+普通弹幕=maxlimit
		trimCnt += start
		normal = normal[start:]
	}
	dms = append(dms, normal...)
	// 追加字幕弹幕和特殊弹幕
	dms = append(dms, special...)
	if trimCnt > 0 {
		err = s.dao.TrimDMCache(c, tp, oid, int64(trimCnt))
	}
	return
}

// 返回所有每个弹幕池对应的弹幕列表
func (s *Service) dms(c context.Context, tp int32, oid, maxlimit int64, childpool int32) (dms []*model.DM, err error) {
	var (
		count         int
		keyprotect    = "kp"
		dmMap         = make(map[string][]*model.DM)
		contentSpeMap = make(map[int64]*model.ContentSpecial)
	)
	idxMap, dmids, spedmids, err := s.dao.Indexs(c, tp, oid)
	if err != nil {
		return
	}
	if len(dmids) == 0 {
		return
	}
	ctsMap, err := s.dao.Contents(c, oid, dmids)
	if err != nil {
		return
	}
	if len(spedmids) > 0 {
		if contentSpeMap, err = s.dao.ContentsSpecial(c, spedmids); err != nil {
			return
		}
	}
	for _, content := range ctsMap {
		if dm, ok := idxMap[content.ID]; ok {
			key := fmt.Sprint(dm.Pool)
			dm.Content = content
			if dm.Pool == model.PoolNormal {
				if dm.AttrVal(model.AttrProtect) == model.AttrYes {
					key = keyprotect
				}
			}
			if dm.Pool == model.PoolSpecial {
				contentSpe, ok := contentSpeMap[dm.ID]
				if ok {
					dm.ContentSpe = contentSpe
				}
			}
			dmMap[key] = append(dmMap[key], dm)
		}
	}
	// dm sort
	for _, dmsTmp := range dmMap {
		sort.Sort(model.DMSlice(dmsTmp))
	}
	// pool = 0 保护弹幕和普通弹幕总和为maxlimit
	if protect, ok := dmMap[keyprotect]; ok {
		if start := len(protect) - int(maxlimit); start > 0 { // 只保留maxlimit条保护弹幕
			protect = protect[start:]
		}
		dms = append(dms, protect...)
		count = len(protect)
	}
	if normal, ok := dmMap[fmt.Sprint(model.PoolNormal)]; ok {
		start := len(normal) + count - int(maxlimit)
		if start > 0 {
			normal = normal[start:]
		}
		dms = append(dms, normal...)
	}
	// pool = 1 字幕弹幕
	if subtitle, ok := dmMap[fmt.Sprint(model.PoolSubtitle)]; ok {
		dms = append(dms, subtitle...)
	}
	// pool =2 特殊弹幕
	if special, ok := dmMap[fmt.Sprint(model.PoolSpecial)]; ok {
		dms = append(dms, special...)
	}
	return
}

func (s *Service) genXML(c context.Context, sub *model.Subject) (xml []byte, err error) {
	realname := s.isRealname(c, sub.Pid, sub.Oid)
	buf := new(bytes.Buffer)
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?><i>`)
	buf.WriteString(`<chatserver>chat.bilibili.com</chatserver><chatid>`)
	buf.WriteString(fmt.Sprint(sub.Oid))
	buf.WriteString(`</chatid><mission>`)
	buf.WriteString(fmt.Sprint(sub.AttrVal(model.AttrSubMission)))
	buf.WriteString(`</mission><maxlimit>`)
	buf.WriteString(fmt.Sprint(sub.Maxlimit))
	buf.WriteString(`</maxlimit>`)
	buf.WriteString(fmt.Sprintf(`<state>%d</state>`, sub.State))
	if realname {
		buf.WriteString(`<real_name>1</real_name>`)
	} else {
		buf.WriteString(`<real_name>0</real_name>`)
	}
	if sub.State == model.SubStateClosed {
		buf.WriteString(`</i>`)
		xml = buf.Bytes()
		return
	}
	dms, err := s.dmsCache(c, sub.Type, sub.Oid, sub.Maxlimit)
	if err != nil {
		return
	}
	if len(dms) > 0 {
		buf.WriteString(`<source>k-v</source>`)
	} else {
		buf.WriteString(`<source>e-r</source>`)
		if dms, err = s.dms(c, sub.Type, sub.Oid, sub.Maxlimit, int32(sub.Childpool)); err != nil {
			return
		}
		if err = s.dao.SetDMCache(c, sub.Type, sub.Oid, dms); err != nil { // add redis cache
			return
		}
	}
	for _, dm := range dms {
		buf.WriteString(dm.ToXML(realname))
	}
	buf.WriteString("</i>")
	xml = buf.Bytes()
	return
}

func (s *Service) isRealname(c context.Context, aid, oid int64) (realname bool) {
	if oid == 13196688 || oid == 290932 {
		realname = true
		return
	}
	arg := &arcMdl.ArgAid2{Aid: aid}
	archive, err := s.arcRPC.Archive3(c, arg)
	if err != nil {
		log.Error("arcRPC.Archive3(%v) error(%v)", arg, err)
		return
	}
	if v, ok := s.realname[int64(archive.TypeID)]; ok && oid >= v {
		realname = true
	} else {
		realname = false
	}
	return
}

// flushXMLSegCache 刷新每个分段的缓存，NOTE:目前只是单纯删除缓存
func (s *Service) flushXMLSegCache(c context.Context, sub *model.Subject) (err error) {
	duration, err := s.videoDuration(c, sub.Pid, sub.Oid)
	if err != nil {
		return
	}
	seg := model.SegmentInfo(0, duration)
	for i := int64(1); i <= seg.Cnt; i++ {
		if err = s.dao.DelXMLSegCache(c, sub.Type, sub.Oid, seg.Cnt, i); err != nil {
			continue
		}
	}
	return
}

// rebuildDmSegCache 刷新视频每个分段弹幕缓存
func (s *Service) flushAllDmSegCache(c context.Context, oid int64, tp int32) (err error) {
	var (
		sub             *model.Subject
		duration, total int64
	)
	if sub, err = s.subject(c, tp, oid); err != nil {
		return
	}
	if duration, err = s.videoDuration(c, sub.Pid, sub.Oid); err != nil {
		return
	}
	total = int64(math.Ceil(float64(duration) / float64(model.DefaultPageSize)))
	for i := int64(1); i <= total; i++ {
		s.asyncAddFlushDMSeg(c, &model.FlushDMSeg{
			Type:  tp,
			Oid:   oid,
			Force: true,
			Page: &model.Page{
				Num:   i,
				Size:  model.DefaultPageSize,
				Total: total,
			},
		})
	}
	log.Info("flushAllDmSegCache oid:%v total:%v", oid, total)
	return
}

func (s *Service) asyncAddFlushDM(c context.Context, fc *model.Flush) {
	select {
	case s.flushMergeChan[fc.Oid%int64(s.routineSize)] <- fc:
	default:
		log.Warn("flush merge channel is full,flush(%+v)", fc)
	}
}
