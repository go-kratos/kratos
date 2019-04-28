package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"sort"

	"go-common/app/interface/main/dm2/model"
	arcMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) dmsCache(c context.Context, tp int32, oid, maxlimit int64, needTrim bool) (dms []*model.DM, err error) {
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
	if trimCnt > 0 && needTrim {
		err = s.dao.TrimDMCache(c, tp, oid, int64(trimCnt))
	}
	return
}

// 返回所有每个弹幕池对应的弹幕列表
func (s *Service) dms(c context.Context, tp int32, oid, maxlimit int64, childpool int32) (dms []*model.DM, err error) {
	var (
		key           string
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
	contents, err := s.dao.Contents(c, oid, dmids)
	if err != nil {
		return
	}
	if len(spedmids) > 0 {
		if contentSpeMap, err = s.dao.ContentsSpecial(c, spedmids); err != nil {
			return
		}
	}
	for _, content := range contents {
		if dm, ok := idxMap[content.ID]; ok {
			key = fmt.Sprint(dm.Pool)
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
	key = fmt.Sprint(model.PoolNormal)
	if normal, ok := dmMap[key]; ok {
		start := len(normal) + count - int(maxlimit)
		if start > 0 {
			normal = normal[start:]
		}
		dms = append(dms, normal...)
	}
	if childpool > 0 {
		// pool = 1 字幕弹幕
		if _, ok := dmMap[fmt.Sprint(model.PoolSubtitle)]; ok {
			dms = append(dms, dmMap[fmt.Sprint(model.PoolSubtitle)]...)
		}
		// pool =2 特殊弹幕
		key = fmt.Sprint(model.PoolSpecial)
		if _, ok := dmMap[key]; ok {
			dms = append(dms, dmMap[key]...)
		}
	}
	return
}

func (s *Service) genXML(c context.Context, sub *model.Subject) (xml []byte, err error) {
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
	realname := s.isRealname(c, sub.Pid, sub.Oid)
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
	dms, err := s.dmsCache(c, sub.Type, sub.Oid, sub.Maxlimit, true)
	if err != nil {
		return
	}
	if len(dms) > 0 {
		buf.WriteString(`<source>k-v</source>`)
	} else {
		buf.WriteString(`<source>e-r</source>`)
		if dms, err = s.dms(c, sub.Type, sub.Oid, sub.Maxlimit, sub.Childpool); err != nil {
			return
		}
		flush := &model.Flush{Type: sub.Type, Oid: sub.Oid, Force: false}
		data, err := json.Marshal(flush)
		if err != nil {
			log.Error("json.Marshal(%v) error(%v)", flush, err)
			return nil, err
		}
		s.dao.SendAction(context.TODO(), fmt.Sprint(sub.Oid), &model.Action{Action: model.ActionFlush, Data: data})
	}
	for _, dm := range dms {
		buf.WriteString(dm.ToXML(realname))
	}
	buf.WriteString("</i>")
	return buf.Bytes(), nil
}

// DMXML return dm xml.
func (s *Service) DMXML(c context.Context, tp int32, oid int64) (data []byte, err error) {
	// get from local cache
	var ok bool
	if data, ok = s.localCache[keyXML(tp, oid)]; ok {
		return
	}
	data, err = s.singleGenXML(c, tp, oid)
	return
}

func (s *Service) singleGenXML(c context.Context, tp int32, oid int64) (data []byte, err error) {
	v, err, _ := s.singleGroup.Do(keyXML(tp, oid), func() (res interface{}, err error) {
		// 从memcache获取
		if data, err = s.dao.XMLCache(c, oid); err != nil {
			return
		}
		if len(data) > 0 {
			res = data
			return
		}
		sub, err := s.subject(c, tp, oid)
		if err != nil {
			return
		}
		if data, err = s.genXML(c, sub); err != nil {
			return
		}
		if len(data) == 0 {
			err = ecode.NothingFound
			return
		}
		if data, err = s.gzflate(data); err != nil {
			log.Error("s.gzflate(oid:%d) error(%v)", oid, err)
			return
		}
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.AddXMLCache(ctx, oid, data)
		})
		res = data
		return
	})
	if err != nil {
		return
	}
	data = v.([]byte)
	return
}

// AjaxDM 返回首页弹幕列表
func (s *Service) AjaxDM(c context.Context, aid int64) (msgs []string, err error) {
	msgs = make([]string, 0)
	res, err := s.arcRPC.Archive3(c, &arcMdl.ArgAid2{Aid: aid})
	if err != nil {
		log.Error("arcRPC.Archive3(%d) error(%v)", aid, err)
		return
	}
	oid := res.FirstCid
	if msgs, err = s.dao.AjaxDMCache(c, oid); err != nil {
		return
	}
	if len(msgs) > 0 {
		return
	}
	dms, err := s.dmsCache(c, model.SubTypeVideo, oid, 20, false)
	if err != nil {
		log.Error("s.dmsCache(%d %d) error(%v)", aid, oid, err)
		return
	}
	for _, dm := range dms {
		if dm.Pool == model.PoolNormal && dm.Content.Mode != model.ModeSpecial {
			msgs = append(msgs, html.EscapeString(dm.Content.Msg))
		}
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.AddAjaxDMCache(ctx, oid, msgs)
	})
	return
}

// JudgeDms get fengjiwei dm list
func (s *Service) JudgeDms(c context.Context, tp int8, oid int64, dmid int64) (judgeDMList *model.JudgeDMList, err error) {
	var (
		start, end, dmidIn int
		length             = 100
		dms                = make([]*model.JDM, 0)
		dmids              []int64
		contentSpec        = make(map[int64]*model.ContentSpecial)
	)
	judgeDMList, err = s.dao.DMJudgeCache(c, tp, oid, dmid)
	if err != nil {
		log.Error("DMJudge:s.dmDao.DMJudgeCache(tp:%d,oid:%d,dmid:%d) error(%v)", tp, oid, dmid, err)
		return
	}
	if judgeDMList != nil {
		return
	}
	judgeDMList = new(model.JudgeDMList)
	judgeDMList.List = make([]*model.JDM, 0)
	judgeDMList.Index = make([]int64, 0)
	idx, err := s.dao.IndexByid(c, tp, oid, dmid)
	if err != nil {
		log.Error("DMJudge:s.dmDao.Index2(type:%d,oid:%d, dmid:%d) error (%v)", tp, oid, dmid, err)
		return
	}
	if idx == nil {
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.SetDMJudgeCache(ctx, tp, oid, dmid, judgeDMList)
		})
		return
	}
	part, spart, err := s.dao.JudgeIndex(c, tp, oid, idx.Ctime-86400, idx.Ctime+86400, idx.Progress-5000, idx.Progress+5000)
	if err != nil {
		log.Error("DMJudge:s.dmDao.JudgeIndex (type:%d,oid:%d) error (%v)", tp, oid, err)
		return
	}
	sort.Sort(model.JudgeSlice(part))
	for k, d := range part {
		if d.ID == dmid {
			dmidIn = k + 1
			break
		}
	}
	if dmidIn == 0 {
		log.Error("DMJudge: cid:%d dmid:%d dm too much", oid, dmid)
		return
	}
	start = dmidIn - length/2
	end = dmidIn + length/2
	if start < 0 {
		start = 0
	}
	if end > len(part) {
		end = len(part)
	}
	part = part[start:end]
	for k, i := range part {
		if i.Mid == idx.Mid {
			judgeDMList.Index = append(judgeDMList.Index, int64(k))
		}
		dmids = append(dmids, i.ID)
	}
	ctsMap, err := s.dao.Contents(c, oid, dmids)
	if err != nil {
		log.Error("DMJudge:s.dmDao.Contents(type:%d,oid:%d dmids:%v) error (%v)", tp, oid, dmids, err)
		return
	}
	if len(spart) > 0 {
		if contentSpec, err = s.dao.ContentsSpecial(c, spart); err != nil {
			log.Error("DMJudge:s.dmDao.ContentsSpecial(type:%d,oid:%d) error (%v)", tp, oid, err)
			return
		}
	}
	for _, i := range part {
		ct, ok := ctsMap[i.ID]
		if !ok {
			continue
		}
		dm := &model.JDM{
			ID:       i.ID,
			Progress: timeStr(int64(i.Progress)),
			Msg:      ct.Msg,
			Mid:      i.Mid,
			CTime:    i.Ctime,
		}
		if i.Pool == model.PoolSpecial {
			if _, ok = contentSpec[dm.ID]; ok {
				dm.Msg = contentSpec[dm.ID].Msg
			}
		}
		dms = append(dms, dm)
	}
	judgeDMList.List = dms
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.SetDMJudgeCache(ctx, tp, oid, dmid, judgeDMList)
	})
	return
}

// timeStr microTime to string.
func timeStr(t int64) string {
	sec := t / 1000 % 60
	min := t / 1000 / 60 % 60
	hour := t / 1000 / 60 / 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
}

func (s *Service) isRealname(c context.Context, aid, oid int64) (realname bool) {
	for _, v := range s.conf.Realname.Whitelist {
		if oid == v {
			realname = true
			return
		}
	}
	if !s.conf.Realname.SwitchOn {
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

// DMDistribution get dm distribution from dm list.
func (s *Service) DMDistribution(c context.Context, typ int32, oid int64, interval int32) (res map[int64]int64, err error) {
	res = make(map[int64]int64)
	sub, err := s.subject(c, typ, oid)
	if err != nil {
		return
	}
	dms, err := s.dmsCache(c, sub.Type, sub.Oid, sub.Maxlimit, true)
	if err != nil {
		return
	}
	if len(dms) == 0 {
		if dms, err = s.dms(c, sub.Type, sub.Oid, sub.Maxlimit, sub.Childpool); err != nil {
			return
		}
		flush := &model.Flush{Type: sub.Type, Oid: sub.Oid, Force: false}
		data, err := json.Marshal(flush)
		if err != nil {
			log.Error("json.Marshal(%v) error(%v)", flush, err)
			return nil, err
		}
		s.dao.SendAction(context.TODO(), fmt.Sprint(sub.Oid), &model.Action{Action: model.ActionFlush, Data: data})
	}
	for _, dm := range dms {
		x := int64(dm.Progress/1000.0/interval) + 1
		if _, ok := res[x]; !ok {
			res[x] = 1
		} else {
			res[x]++
		}
	}
	return
}

func (s *Service) loadLocalcache(oids []int64) {
	var (
		err      error
		duration int64
		sub      *model.Subject
		tp       = model.SubTypeVideo
		c        = context.Background()
		tmp      = make(map[string][]byte)
		bs       []byte
	)
	for _, oid := range oids {
		if sub, err = s.dao.Subject(c, tp, oid); err != nil || sub == nil {
			continue
		}
		if bs, err = json.Marshal(sub); err != nil {
			continue
		}
		tmp[keySubject(sub.Type, sub.Oid)] = bs
		// local cache video duration
		if duration, err = s.videoDuration(c, sub.Pid, oid); err != nil {
			continue
		}
		tmp[keyDuration(tp, oid)] = []byte(fmt.Sprint(duration))
		// local cache segment dm xml
		seg := model.SegmentInfo(0, duration)
		for i := int64(1); i <= seg.Cnt; i++ { // flush every segment cache
			seg.Num = i
			var xml []byte
			if xml, err = s.singleGenSegXML(c, sub.Pid, sub, seg); err != nil {
				continue
			}
			tmp[keySeg(tp, oid, seg.Cnt, seg.Num)] = xml
		}
		// local cache dm xml
		data, err := s.singleGenXML(c, tp, oid)
		if err != nil {
			continue
		}
		tmp[keyXML(tp, oid)] = []byte(data)
	}
	if len(tmp) > 0 {
		s.localCache = tmp
	}
}

// dmList get dm list from database.
func (s *Service) dmList(c context.Context, tp int32, oid int64, dmids []int64) (dms []*model.DM, err error) {
	if len(dmids) == 0 {
		return
	}
	dms = make([]*model.DM, 0, len(dmids))
	contentSpe := make(map[int64]*model.ContentSpecial)
	idxMap, special, err := s.dao.IndexsByid(c, tp, oid, dmids)
	if err != nil || len(idxMap) == 0 {
		return
	}
	ctsMap, err := s.dao.Contents(c, oid, dmids)
	if err != nil {
		return
	}
	if len(special) > 0 {
		if contentSpe, err = s.dao.ContentsSpecial(c, special); err != nil {
			return
		}
	}
	for _, dmid := range dmids {
		if dm, ok := idxMap[dmid]; ok {
			var content *model.Content
			if content, ok = ctsMap[dmid]; ok {
				dm.Content = content
			} else {
				log.Error("dm content not exist,tp:%d,oid:%d,dmid:%d", tp, oid, dmid)
				continue
			}
			if dm.Pool == model.PoolSpecial {
				if _, ok = contentSpe[dmid]; ok {
					dm.ContentSpe = contentSpe[dmid]
				} else {
					log.Error("dm special content not exist,tp:%d,oid:%d,dmid:%d", tp, oid, dmid)
					continue
				}
			}
			dms = append(dms, dm)
		}
	}
	return
}
