package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"go-common/app/admin/main/dm/model"
	"go-common/app/admin/main/dm/model/oplog"
	accountApi "go-common/app/service/main/account/api"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_pageSize = 50
)

// dms get dm list from database.
func (s *Service) dms(c context.Context, tp int32, oid int64, dmids []int64) (dms []*model.DM, err error) {
	if len(dmids) == 0 {
		return
	}
	contentSpe := make(map[int64]*model.ContentSpecial)
	idxMap, special, err := s.dao.IndexsByID(c, tp, oid, dmids)
	if err != nil || len(idxMap) == 0 {
		return
	}
	contents, err := s.dao.Contents(c, oid, dmids)
	if err != nil {
		return
	}
	if len(special) > 0 {
		if contentSpe, err = s.dao.SpecialContents(c, special); err != nil {
			return
		}
	}
	for _, content := range contents {
		if idx, ok := idxMap[content.ID]; ok {
			idx.Content = content
			if idx.Pool == model.PoolSpecial {
				if _, ok = contentSpe[idx.ID]; ok {
					idx.ContentSpe = contentSpe[idx.ID]
				}
			}
			dms = append(dms, idx)
		}
	}
	return
}

// DMSearch danmu list from search.
func (s *Service) DMSearch(c context.Context, p *model.SearchDMParams) (res *model.SearchDMResult, err error) {
	var (
		mids, dmids []int64
		sorted      []*model.DM
		protectCnt  int64
		dmMap       = make(map[int64]*model.DM)
		uidMap      = make(map[int64]bool)
	)
	res = &model.SearchDMResult{}
	sub, err := s.dao.Subject(c, p.Type, p.Oid)
	if err != nil {
		log.Error("s.dao.Subject(%d,%d) error(%v)", p.Type, p.Oid, err)
		return
	}
	if sub == nil {
		return
	}
	srchData, err := s.dao.SearchDM(c, p)
	if err != nil {
		log.Error("s.dao.SearchDM(%v) error(%v)", p, err)
		return
	}
	if srchData == nil {
		return
	}
	for _, v := range srchData.Result {
		dmids = append(dmids, v.ID)
	}
	dms, err := s.dms(c, p.Type, p.Oid, dmids)
	if err != nil {
		log.Error("s.dms(%d,%v) error(%v)", p.Oid, dmids, err)
		return
	}
	for _, dm := range dms {
		dmMap[dm.ID] = dm
		if _, ok := uidMap[dm.Mid]; !ok && dm.Mid > 0 {
			uidMap[dm.Mid] = true
			mids = append(mids, dm.Mid)
		}
	}
	for _, dmid := range dmids {
		if dm, ok := dmMap[dmid]; ok {
			sorted = append(sorted, dm)
		}
	}
	total := len(mids)
	pageNum := total / _pageSize
	if total%_pageSize != 0 {
		pageNum++
	}
	var (
		g       errgroup.Group
		lk      sync.Mutex
		infoMap = make(map[int64]*account.Info, total)
	)
	for i := 0; i < pageNum; i++ {
		start := i * _pageSize
		end := (i + 1) * _pageSize
		if end > total {
			end = total
		}
		g.Go(func() (err error) {
			var (
				arg = &accountApi.MidsReq{Mids: mids[start:end]}
				res *accountApi.InfosReply
			)
			if res, err = s.accountRPC.Infos3(c, arg); err != nil {
				log.Error("s.accRPC.Infos3(%v) error(%v)", arg, err)
			} else {
				for mid, info := range res.GetInfos() {
					lk.Lock()
					infoMap[mid] = info
					lk.Unlock()
				}
			}
			return
		})
	}
	g.Go(func() (err error) {
		if protectCnt, err = s.dao.SearchProtectCount(context.TODO(), p.Type, p.Oid); err != nil {
			log.Error("s.dao.SearchProtectCount(%d,%d) error(%v)", p.Type, p.Oid, err)
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	for _, dm := range sorted {
		msg := dm.Content.Msg
		if dm.Pool == model.PoolSpecial && dm.ContentSpe != nil {
			msg = dm.ContentSpe.Msg
		}
		item := &model.DMItem{
			IDStr:    strconv.FormatInt(dm.ID, 10),
			ID:       dm.ID,
			Type:     dm.Type,
			Oid:      dm.Oid,
			Mid:      dm.Mid,
			Pool:     dm.Pool,
			State:    dm.State,
			Attrs:    dm.AttrNtoA(),
			Msg:      msg,
			Ctime:    dm.Ctime,
			Mode:     dm.Content.Mode,
			IP:       dm.Content.IP,
			Color:    fmt.Sprintf("#%06X", dm.Content.Color),
			Progress: dm.Progress,
			Fontsize: dm.Content.FontSize,
		}
		if info, ok := infoMap[dm.Mid]; ok {
			item.Uname = info.Name
		}
		res.Result = append(res.Result, item)
	}
	res.MaxLimit = sub.Maxlimit
	res.Total = sub.ACount
	res.Page = srchData.Page.Num
	res.Pagesize = srchData.Page.Size
	res.Deleted = sub.ACount - sub.Count
	res.Protected = protectCnt
	res.Count = srchData.Page.Total
	return
}

// XMLCacheFlush 刷新弹幕缓存
func (s *Service) XMLCacheFlush(c context.Context, tp int32, oid int64) {
	v := make(map[string]interface{})
	v["type"] = tp
	v["oid"] = oid
	v["force"] = true
	data, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", v, err)
		return
	}
	action := &model.Action{Action: model.ActFlushDM, Data: data, Oid: oid}
	s.addAction(action)
}

// EditDMState multi edit dm state.
func (s *Service) EditDMState(c context.Context, tp, state int32, oid int64, reason int8, dmids []int64, moral float64, adminID int64, operator, remark string) (err error) {
	if err = s.editDmState(c, tp, state, oid, reason, dmids, moral, adminID, operator, remark); err != nil {
		log.Error("s.dao.UpSearchDMState(%d,%d,%v) err (%v)", tp, oid, dmids, err)
		return
	}
	// update dm search index
	if err = s.uptSearchDmState(c, tp, state, map[int64][]int64{oid: dmids}); err != nil {
		log.Error("s.dao.UpSearchDMState(%d,%d,%v) err (%v)", tp, oid, dmids, err)
	}
	return
}

// editDmState multi edit dm state.
func (s *Service) editDmState(c context.Context, tp, state int32, oid int64, reason int8, dmids []int64, moral float64, adminID int64, operator, remark string) (err error) {
	sub, err := s.dao.Subject(c, tp, oid)
	if err != nil || sub == nil {
		return
	}
	dms, err := s.dms(c, tp, oid, dmids)
	if err != nil {
		return
	}
	count := countDMNum(dms, state)
	affect, err := s.dao.SetStateByIDs(c, tp, oid, dmids, state)
	if err != nil || affect == 0 {
		return
	}
	if sub.Count+count < 0 {
		count = -sub.Count
	}
	// update dm_index count
	if _, err = s.dao.IncrSubjectCount(c, tp, oid, count); err != nil {
		return
	}
	// write dm admin log
	if affect > 0 {
		if remark == "" {
			remark = model.AdminRptReason[reason]
		}
		s.OpLog(c, oid, adminID, tp, dmids, "status", "", fmt.Sprint(state), remark, oplog.SourceManager, oplog.OperatorAdmin)
	}
	// update dm monitor count
	if sub.IsMonitoring() {
		s.oidLock.Lock()
		s.moniOidMap[sub.Oid] = struct{}{}
		s.oidLock.Unlock()
	}
	// flush xml cache
	s.XMLCacheFlush(c, tp, oid)
	// reduce moral
	uidMap := make(map[int64]struct{})
	for _, dm := range dms {
		if _, ok := uidMap[dm.Mid]; !ok {
			uidMap[dm.Mid] = struct{}{}
		}
	}
	if len(uidMap) > 0 && moral > 0 {
		for uid := range uidMap {
			s.reduceMoral(c, uid, int64(-moral), reason, operator, "弹幕管理")
		}
	}
	return
}

// EditDMPool edit dm pool.
func (s *Service) EditDMPool(c context.Context, tp int32, oid int64, pool int32, dmids []int64, adminID int64) (err error) {
	sub, err := s.dao.Subject(c, tp, oid)
	if err != nil || sub == nil {
		return
	}
	if sub.Childpool < pool {
		if _, err = s.dao.UpSubjectPool(c, tp, oid, pool); err != nil {
			return
		}
	}
	affect, err := s.dao.SetPoolIDByIDs(c, tp, oid, pool, dmids)
	if err != nil {
		return
	}
	if affect > 0 {
		if pool == model.PoolNormal {
			s.dao.IncrSubMoveCount(c, tp, oid, -affect) // NOTE update move_count,ignore error
		} else {
			s.dao.IncrSubMoveCount(c, tp, oid, affect) // NOTE update move_count,ignore error
		}
		s.OpLog(c, oid, adminID, tp, dmids, "pool", "", fmt.Sprint(pool), "弹幕池变更", oplog.SourceManager, oplog.OperatorAdmin)
	}
	if err = s.uptSearchDMPool(c, tp, oid, pool, dmids); err != nil {
		log.Error("s.dao.UpSearchDMPool(%d,%d,%v) err (%v)", tp, oid, dmids, err)
	}
	return
}

// EditDMAttr change dm attr
func (s *Service) EditDMAttr(c context.Context, tp int32, oid int64, dmids []int64, bit uint, value int32, adminID int64) (err error) {
	var (
		eg      = errgroup.Group{}
		attrMap = make(map[int32][]int64)
	)
	dms, err := s.dms(c, tp, oid, dmids)
	if err != nil {
		log.Error("s.dms(oid:%d ids:%v) error(%v)", oid, dmids, err)
		return
	}
	for _, dm := range dms {
		dm.AttrSet(value, bit)
		attrMap[dm.Attr] = append(attrMap[dm.Attr], dm.ID)
	}
	for k, v := range attrMap {
		attr := k
		ids := v
		eg.Go(func() (err error) {
			affect, err := s.dao.SetAttrByIDs(c, tp, oid, ids, attr)
			if err != nil {
				log.Error("s.dao.SetAttrByIDs(oid:%d ids:%v) error(%v)", oid, ids, err)
				return
			}
			if affect > 0 {
				s.OpLog(c, oid, adminID, tp, ids, "attribute", "", fmt.Sprint(attr), "弹幕保护状态变更", oplog.SourceManager, oplog.OperatorAdmin)
			}
			if err = s.uptSearchDMAttr(c, tp, oid, attr, dmids); err != nil {
				log.Error("dao.UpSearchDMAttr(oid:%d,attr:%d) error(%v)", oid, attr, err)
			}
			return
		})
	}
	return eg.Wait()
}

// DMIndexInfo get dm index info
func (s *Service) DMIndexInfo(c context.Context, cid int64) (info *model.DMIndexInfo, err error) {
	info = new(model.DMIndexInfo)
	sub, err := s.dao.Subject(c, model.SubTypeVideo, cid)
	if err != nil || sub == nil {
		return
	}
	argAid2 := &archive.ArgAid2{Aid: sub.Pid}
	arc, err := s.arcRPC.Archive3(c, argAid2)
	if err != nil {
		log.Error("s.arcRPC.Archive3(%v) error(%v)", argAid2, err)
		err = nil
	} else {
		info.Title = arc.Title
		info.Cover = arc.Pic
	}
	argVideo := &archive.ArgVideo2{Aid: sub.Pid, Cid: cid}
	video, err := s.arcRPC.Video3(c, argVideo)
	if err != nil {
		log.Error("s.arcRPC.Video3(%v) error(%v)", argVideo, err)
		err = nil
	} else {
		info.Duration = video.Duration
		info.ETitle = video.Part
	}
	argMid := &accountApi.MidReq{Mid: sub.Mid}
	uInfo, err := s.accountRPC.Info3(c, argMid)
	if err != nil {
		log.Error("s.accRPC.Info3(%v) error(%v)", argMid, err)
		err = nil
	} else {
		info.UName = uInfo.GetInfo().GetName()
	}
	info.AID = sub.Pid
	info.CID = sub.Oid
	info.MID = sub.Mid
	info.Limit = sub.Maxlimit
	if sub.State == model.SubStateOpen {
		info.Active = 1
	} else {
		info.Active = 0
	}
	info.CTime = int64(sub.Ctime)
	info.MTime = int64(sub.Mtime)
	return
}

// countDMNum count state changed dm count
func countDMNum(dms []*model.DM, state int32) (count int64) {
	for _, dm := range dms {
		if model.DMVisible(dm.State) && !model.DMVisible(state) {
			count--
		} else if !model.DMVisible(dm.State) && model.DMVisible(state) {
			count++
		}
	}
	return count
}

// FixDMCount fix dm acount,count of aid.
func (s *Service) FixDMCount(c context.Context, aid int64) (err error) {
	var (
		arg         = archive.ArgAid2{Aid: aid}
		oids        []int64
		states      = []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13} // 弹幕所有状态
		normalState = []int64{0, 2, 6}                                      // 前台可能展示的弹幕状态
	)
	pages, err := s.arcRPC.Page3(c, &arg)
	if err != nil {
		log.Error("arcRPC.Page3(%v) error(%v)", arg, err)
		return
	}
	if len(pages) == 0 {
		log.Warn("aid:%d have no pages", aid)
		return
	}
	for _, page := range pages {
		oids = append(oids, page.Cid)
	}
	subs, err := s.dao.Subjects(c, model.SubTypeVideo, oids)
	if err != nil {
		return
	}
	for _, sub := range subs {
		tp := sub.Type
		oid := sub.Oid
		s.cache.Do(c, func(ctx context.Context) {
			acount, err := s.dao.DMCount(ctx, tp, oid, states)
			if err != nil {
				return
			}
			count, err := s.dao.DMCount(ctx, tp, oid, normalState)
			if err != nil {
				return
			}
			s.dao.UpSubjectCount(ctx, tp, oid, acount, count) // 更新新库dm_subject
			log.Info("fix dm count,type:%d,oid:%d,acount:%d,count:%d", tp, oid, acount, count)
		})
	}
	return
}

func (s *Service) uptSearchDmState(c context.Context, tp int32, state int32, dmidM map[int64][]int64) (err error) {
	if err = s.dao.UpSearchDMState(c, tp, state, dmidM); err != nil {
		return
	}
	if err = s.dao.UpSearchRecentDMState(c, tp, state, dmidM); err != nil {
		return
	}
	return
}

func (s *Service) uptSearchDMPool(c context.Context, tp int32, oid int64, pool int32, dmids []int64) (err error) {
	if err = s.dao.UpSearchDMPool(c, tp, oid, pool, dmids); err != nil {
		return
	}
	if err = s.dao.UpSearchRecentDMPool(c, tp, oid, pool, dmids); err != nil {
		return
	}
	return
}

func (s *Service) uptSearchDMAttr(c context.Context, tp int32, oid int64, attr int32, dmids []int64) (err error) {
	if err = s.dao.UpSearchDMAttr(c, tp, oid, attr, dmids); err != nil {
		return
	}
	if err = s.dao.UpSearchRecentDMAttr(c, tp, oid, attr, dmids); err != nil {
		return
	}
	return
}
