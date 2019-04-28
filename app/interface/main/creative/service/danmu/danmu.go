package danmu

import (
	"context"
	"encoding/json"
	"go-common/app/interface/main/creative/model/danmu"
	pubSvc "go-common/app/interface/main/creative/service"
	dmMdl "go-common/app/interface/main/dm/model"
	"go-common/app/service/main/account/model"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"strconv"
	"time"
)

var (
	stateMap = map[int8]int8{
		0: 0, // 正常
		1: 1, // 删除
		2: 2, // 保护
		3: 3, // 取消保护
	}
	poolMap = map[int8]int8{
		0: 0, // 普通弹幕池
		1: 1, // 字幕弹幕池
		2: 2, // 特殊弹幕池
	}
)

// AdvDmPurchaseList fn
func (s *Service) AdvDmPurchaseList(c context.Context, mid int64, ip string) (danmus []*danmu.AdvanceDanmu, err error) {
	if danmus, err = s.dm.GetAdvDmPurchases(c, mid, ip); err != nil {
		log.Error("s.dm.ListAdvDmPurchases err(%v) | mid(%d), ip(%s)", err, mid, ip)
		return
	}
	return
}

// PassDmPurchase fn
func (s *Service) PassDmPurchase(c context.Context, mid, id int64, ip string) (err error) {
	if err = s.dm.PassAdvDmPurchase(c, mid, id, ip); err != nil {
		log.Error("s.dm.PassDmPurchase err(%v) | mid(%d), id(%d), ip(%s)", err, mid, id, ip)
		return
	}
	return
}

// DenyDmPurchase fn
func (s *Service) DenyDmPurchase(c context.Context, mid, id int64, ip string) (err error) {
	if err = s.dm.DenyAdvDmPurchase(c, mid, id, ip); err != nil {
		log.Error("s.dm.DenyDmPurchase err(%v) | mid(%d), id(%d), ip(%s)", err, mid, id, ip)
		return
	}
	return
}

// CancelDmPurchase fn
func (s *Service) CancelDmPurchase(c context.Context, mid, id int64, ip string) (err error) {
	if err = s.dm.CancelAdvDmPurchase(c, mid, id, ip); err != nil {
		log.Error("s.dm.CancelDmPurchase err(%v) | mid(%d), id(%d), ip(%s)", err, mid, id, ip)
		return
	}
	return
}

// Edit fn
func (s *Service) Edit(c context.Context, mid, id int64, state int8, dmids []int64, ip string) (err error) {
	if len(dmids) == 0 {
		err = ecode.CreativeDanmuFilterParamError
		return
	}
	okState := false
	for _, v := range stateMap {
		if v == state {
			okState = true
		}
	}
	if !okState {
		err = ecode.RequestErr
		return
	}
	if err = s.dm.Edit(c, mid, id, state, dmids, ip); err != nil {
		log.Error("s.dm.DmList err(%v)|mid(%d),id(%d),state(%d),dmids(%+v),ip(%s)", err, mid, id, state, dmids, ip)
		return
	}
	return
}

// Transfer fn
func (s *Service) Transfer(c context.Context, mid, fromCID, toCID int64, offset float64, ak, ck, ip string) (err error) {
	if err = s.dm.Transfer(c, mid, fromCID, toCID, offset, ak, ck, ip); err != nil {
		log.Error("s.dm.Transfer err(%v) | mid(%d),fromCID(%d),toCID(%d),offset(%+v),ak(%s),ck(%s),ip(%s)", err, mid, fromCID, toCID, offset, ak, ck, ip)
		return
	}
	return
}

// UpPool fn
func (s *Service) UpPool(c context.Context, cid int64, dmids []int64, pool int8, mid int64, ip string) (err error) {
	okPool := false
	for _, v := range poolMap {
		if v == pool {
			okPool = true
		}
	}
	if !okPool {
		err = ecode.RequestErr
		return
	}
	if err = s.dm.UpPool(c, mid, cid, dmids, pool); err != nil {
		log.Error("s.dm.UpPool err(%v)|cid(%d),dmids(%+v),pool(%+v),mid(%d),ip(%s)", err, cid, dmids, pool, mid, ip)
		return
	}
	return
}

// Distri fn
func (s *Service) Distri(c context.Context, mid, cid int64, ip string) (dmDistri map[int64]int64, err error) {
	if dmDistri, err = s.dm.Distri(c, mid, cid, ip); err != nil {
		log.Error("s.dm.Distri err(%v) | mid(%d),cid(%d)|ip(%s)", err, mid, cid, ip)
		return
	}
	return
}

// List fn
func (s *Service) List(c context.Context, mid, aid, cid int64, pn, ps int, order, pool, midStr, ip string) (dmList *danmu.DmList, err error) {
	var (
		v         *api.Page
		a         *api.Arc
		mids      []int64
		elecs     map[int64]int
		followers map[int64]int
		senders   map[int64]*account.Info
	)
	a, _ = s.arc.Archive(c, aid, ip)
	v, _ = s.arc.Video(c, aid, cid, ip)
	if dmList, err = s.dm.List(c, cid, mid, pn, ps, order, pool, midStr, ip); err != nil {
		log.Error("s.dm.List err(%v) | mid(%d), cid(%d),pn(%d),ps(%d),order(%s),midStr(%s),ip(%s)", err, mid, cid, pn, ps, order, midStr, ip)
		return
	}
	for _, l := range dmList.List {
		mids = append(mids, l.Mid)
	}
	if len(mids) > 0 {
		var g, ctx = errgroup.WithContext(c)
		g.Go(func() error { //获取被充电状态
			elecs, _ = s.elec.ElecRelation(ctx, mid, mids, ip)
			return nil
		})
		g.Go(func() error { //获取用户信息
			senders, _ = s.acc.Infos(ctx, mids, ip)
			return nil
		})
		g.Go(func() error { //获取被关注状态
			followers, _ = s.acc.Followers(ctx, mid, mids, ip)
			return nil
		})
		g.Wait()
		for _, l := range dmList.List {
			l.VTitle = v.Part
			l.ArcTitle = a.Title
			l.Aid = aid
			l.Cover = pubSvc.CoverURL(a.Pic)
			if elec, ok := elecs[l.Mid]; ok { //设置充电状态
				l.IsElec = elec
			}
			if u, ok := senders[l.Mid]; ok { //设置图像和用户名
				l.Uname = u.Name
				l.Uface = u.Face
			}
			if relation, ok := followers[l.Mid]; ok { //设置关注状态
				l.Relation = relation
			}
		}
	}
	return
}

// Recent fn
func (s *Service) Recent(c context.Context, mid, pn, ps int64, ip string) (dmRecent *danmu.DmRecent, err error) {
	var (
		aids      []int64
		mids      []int64
		elecs     map[int64]int
		followers map[int64]int
		senders   map[int64]*account.Info
	)
	if dmRecent, aids, err = s.dm.Recent(c, mid, pn, ps, ip); err != nil {
		log.Error("s.dm.Recent err(%v)|(%v),(%v),(%v),(%v),", err, mid, pn, ps, ip)
		return
	}
	if len(dmRecent.List) > 0 && len(aids) > 0 {
		var avm map[int64]*api.Arc
		if avm, err = s.arc.Archives(c, aids, ip); err != nil {
			log.Error("s.arc.Archives mid(%d)|aids(%+v)|ip(%s)|err(%v)", mid, aids, ip, err)
			err = nil
		}
		for _, l := range dmRecent.List {
			mids = append(mids, l.Mid)
		}
		var g, ctx = errgroup.WithContext(c)
		g.Go(func() error { //获取被充电状态
			elecs, _ = s.elec.ElecRelation(ctx, mid, mids, ip)
			return nil
		})
		g.Go(func() error { //获取用户信息
			senders, _ = s.acc.Infos(ctx, mids, ip)
			return nil
		})
		g.Go(func() error { //获取被关注状态
			followers, _ = s.acc.Followers(ctx, mid, mids, ip)
			return nil
		})
		g.Wait()
		for _, l := range dmRecent.List {
			if av, ok := avm[l.Aid]; ok && av != nil {
				l.Cover = pubSvc.CoverURL(av.Pic)
			}
			if elec, ok := elecs[l.Mid]; ok { //设置充电状态
				l.IsElec = elec
			}
			if u, ok := senders[l.Mid]; ok { //设置图像和用户名
				l.Uname = u.Name
				l.Uface = u.Face
			}
			if relation, ok := followers[l.Mid]; ok { //设置关注状态
				l.Relation = relation
			}
		}
	}
	return
}

// DmReportCheck fn
func (s *Service) DmReportCheck(c context.Context, mid, cid, dmid, op int64, ip string) (err error) {
	if err = s.dm.ReportUpEdit(c, mid, dmid, cid, op, ip); err != nil {
		log.Error("s.dm.ReportUpEdit err(%v) | mid(%d), dmid(%d),cid(%d),op(%d),ip(%s)", err, mid, dmid, cid, op, ip)
		return
	}
	return
}

// DmProtectArchive fn
func (s *Service) DmProtectArchive(c context.Context, mid int64, ip string) (vlist []*dmMdl.Video, err error) {
	if vlist, err = s.dm.ProtectApplyVideoList(c, mid, ip); err != nil {
		log.Error("s.dm.ProtectApplyVideoList err(%v) | mid(%d),ip(%s)", err, mid, ip)
		return
	}
	return
}

// DmProtectList fn
func (s *Service) DmProtectList(c context.Context, mid int64, page int64, aidStr, sort, ip string) (list *danmu.ApplyList, err error) {
	if list, err = s.dm.ProtectApplyList(c, mid, page, aidStr, sort, ip); err != nil {
		log.Error("s.dm.ProtectApplyList err(%v) | mid(%d),aidStr(%s),page(%d),sort(%s),ip(%s)", err, mid, aidStr, page, sort, ip)
		return
	}
	return
}

// DmProtectOper fn
func (s *Service) DmProtectOper(c context.Context, mid, status int64, idsStr, ip string) (err error) {
	if err = s.dm.ProtectOper(c, mid, status, idsStr, ip); err != nil {
		log.Error("s.dm.ProtectOper err(%v) | mid(%d), status(%d),ids(%s),ip(%s)", err, mid, status, idsStr, ip)
		return
	}
	return
}

// UserMid for dm search by name
func (s *Service) UserMid(c context.Context, name, ip string) (mid int64, err error) {
	if mid, err = s.acc.MidByName(c, name); err != nil {
		log.Error("UserMid err(%v) | name(%s), ip(%s)", err, name, ip)
		if err == ecode.AccountInexistence {
			err = nil
		}
		return
	}
	return
}

// DmReportList fn
func (s *Service) DmReportList(c context.Context, mid, pn, ps int64, aidStr, ip string) (res map[string]interface{}, err error) {
	var (
		pager = map[string]int64{
			"total":   0,
			"current": pn,
			"size":    ps,
		}
		list   = []*danmu.DmReport{}
		srchs  = []*dmMdl.RptSearch{}
		ars    = []*danmu.DmArc{}
		total  = int64(0)
		owners = []int64{}
	)
	res = map[string]interface{}{
		"pager":    pager,
		"list":     list,
		"archives": ars,
	}
	if srchs, total, err = s.dm.ReportUpList(c, mid, pn, ps, aidStr, ip); err != nil {
		log.Error("s.dm.ReportUpList err(%v) | mid(%d),pn(%d),ps(%d),aidStr(%s),ip(%s)", err, mid, pn, ps, aidStr, ip)
		return
	}
	if len(srchs) > 0 {
		for _, v := range srchs {
			owners = append(owners, v.Owner)
		}
		var ownerProfiles map[int64]*model.Info
		if ownerProfiles, err = s.acc.Infos(c, owners, ip); err != nil {
			log.Error("s.cc.Infos err(%v) | mid(%d),pn(%d),ps(%d),aidStr(%s),ip(%s)", err, mid, pn, ps, aidStr, ip)
			return
		}
		for _, s := range srchs {
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", s.RpTime, time.Local)
			report := &danmu.DmReport{
				RpID:       s.ID,
				DmInID:     s.Cid,
				AID:        s.Aid,
				Pic:        pubSvc.CoverURL(s.Cover),
				ReportTime: t.Unix(),
				Title:      s.Title,
				Reason:     s.Content,
				DmID:       s.Did,
				DmIDStr:    strconv.FormatInt(s.Did, 10),
				UpUID:      s.UPUid,
				Content:    s.Msg,
				UID:        s.Owner,
			}
			if profile, ok := ownerProfiles[s.Owner]; ok {
				report.UserName = profile.Name
			}
			list = append(list, report)
		}
		if ars, err = s.dm.ReportUpArchives(c, mid, ip); err != nil {
			log.Error("s.dm.ReportUpArchives err(%v) | mid(%d),ip(%s)", err, mid, ip)
			return
		}
	}
	pager["total"] = total
	res["pager"] = pager
	res["list"] = list
	res["archives"] = ars
	return
}

// EditBatch fn
func (s *Service) EditBatch(c context.Context, mid int64, paramsJSON, ip string) (err error) {
	type P struct {
		CID   int64 `json:"cid"`
		DmID  int64 `json:"dmid"`
		State int8  `json:"state"`
	}
	var filtersJSONData []*P
	if err = json.Unmarshal([]byte(paramsJSON), &filtersJSONData); err != nil {
		err = ecode.RequestErr
		return
	}
	if len(filtersJSONData) == 0 {
		err = ecode.CreativeDanmuFilterParamError
		return
	}
	var (
		g errgroup.Group
	)
	for _, v := range filtersJSONData {
		dmids := []int64{}
		dmids = append(dmids, v.DmID)
		cid := v.CID
		state := v.State
		g.Go(func() (err error) {
			if err = s.dm.Edit(c, mid, cid, state, dmids, ip); err != nil {
				log.Error("s.d.Edit v(%+v)|dmids(%+v)|cid(%+v)|state(%+v)|err(%+v)", v, dmids, cid, state, err)
			}
			log.Info("filtersJSONData v(%+v)|dmids(%+v)|cid(%+v)|state(%+v)|err(%+v)", v, dmids, cid, state, err)
			return
		})
	}
	g.Wait()
	return
}
