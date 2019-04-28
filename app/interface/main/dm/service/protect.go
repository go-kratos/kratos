package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/dm/model"
	dm2Mdl "go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/model/oplog"
	account "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_notifyUpTitle     = "有新的弹幕保护申请"
	_notifyUpContent   = `您今天新增了一些未处理弹幕保护申请，前去 #{创作中心 - 哔哩哔哩弹幕视频网 - ( ゜- ゜)つロ 乾杯~}{"http://member.bilibili.com/v/#/danmu/report/save"} 处理吧`
	_notifyUsrTitle    = "弹幕保护申请情况更新~"
	_pa48              = "由于up主日理万机，您之前申请的%d条弹幕暂未受理，请稍后再次申请"
	_paa               = `您在视频 #{%s}{"http://www.bilibili.com/video/av%d/"} 已被全部保护`
	_pap               = `您在视频 #{%s}{"http://www.bilibili.com/video/av%d/"} 已被部分保护`
	_protectApplyLevel = 4
	_maxProtectApply   = 100
	_paExpire          = 48 * 3600
)

// AddProtectApply 批量申请保护弹幕
func (s *Service) AddProtectApply(c context.Context, uid, cid int64, dmids []int64) (err error) {
	if len(dmids) == 0 {
		err = ecode.DMNotFound
		return
	}
	count, err := s.dao.PaUsrCnt(c, uid)
	if err != nil {
		log.Error("s.dao.PaUsrCnt(%d) error(%v)", uid, err)
		return
	}
	if (count + len(dmids)) > _maxProtectApply {
		err = ecode.DMPAUserLimit
		return
	}
	cardReply, err := s.accountSvc.Card3(c, &account.MidReq{Mid: uid})
	if err != nil {
		log.Error("s.actSvc.Card3(%d) error(%v)", uid, err)
		return
	}
	if cardReply.GetCard().GetLevel() < _protectApplyLevel {
		err = ecode.DMPAUserLevel
		return
	}
	dms, err := s.dms(c, model.SubTypeVideo, cid, dmids)
	if err != nil {
		return
	}
	dc := len(dmids)
	if len(dms) < 1 && dc == 1 {
		err = ecode.DMNotFound
		return
	}
	sub, err := s.subject(c, 1, cid)
	if err != nil {
		return
	}
	aps := make([]*model.Pa, 0, len(dms))
	now := time.Now().Unix()
	var ctime time.Time
	for _, dm := range dms {
		if dm.Mid != uid && dc == 1 {
			err = ecode.DMPADMNotOwner
			return
		}
		if dm.AttrVal(model.AttrProtect) == model.AttrYes && dc == 1 {
			err = ecode.DMPADMProtected
			return
		}
		if !dm.NeedDisplay() && dc == 1 {
			err = ecode.DMNotFound
			return
		}
		ctime, err = s.dao.ProtectApplyTime(c, dm.ID)
		if err != nil {
			log.Error("dao.ProtectApplyTime(%d) error(%v)", uid, err)
			continue
		}
		if now-ctime.Unix() < _paExpire {
			if dc == 1 {
				err = ecode.DMPADMLimit
				return
			}
			continue
		}
		ap := &model.Pa{
			CID:      cid,
			UID:      sub.Mid,
			ApplyUID: dm.Mid,
			AID:      sub.Pid,
			Playtime: float32(dm.Progress) / 1000,
			DMID:     dm.ID,
			Msg:      dm.Content.Msg,
			Status:   -1,
			Ctime:    time.Now(),
			Mtime:    time.Now(),
		}
		aps = append(aps, ap)
	}
	if len(aps) < 1 {
		err = ecode.DMPAFailed
		return
	}
	affect, err := s.dao.AddProtectApply(c, aps)
	if err != nil {
		log.Error("dao.AddProtectApply(%v) error(%v)", aps, err)
		return
	}
	if err = s.dao.UptUsrPaCnt(c, uid, affect); err != nil {
		log.Error("s.dao.UptUsrPaCnt(%d,%d) error(%v)", uid, affect, err)
	}
	return
}

// UptPaSwitch 保护弹幕申请开关
func (s *Service) UptPaSwitch(c context.Context, uid int64, status int) (err error) {
	if status != 1 {
		status = 0
	}
	_, err = s.dao.UptPaNoticeSwitch(c, uid, status)
	return
}

// UptPaStatus 审核保护弹幕申请
func (s *Service) UptPaStatus(c context.Context, mid int64, ids []int64, status int) (err error) {
	dmids, err := s.dao.ProtectApplyByIDs(c, mid, xstr.JoinInts(ids))
	if err != nil {
		log.Error("s.dao.ProtectApplyByIDs(%d,%s) error(%v)", mid, xstr.JoinInts(ids), err)
		return
	}
	if status != 1 {
		status = 0
	}
	if _, err = s.dao.UptPaStatus(c, mid, xstr.JoinInts(ids), status); err != nil {
		log.Error("s.dao.UptPaStatus(%d,%v,%d) error(%v)", mid, ids, status, err)
		return
	}
	if status == 0 {
		return
	}
	for oid, ids := range dmids {
		arg := &dm2Mdl.ArgEditDMAttr{
			Type:         1,
			Oid:          oid,
			Mid:          mid,
			Bit:          dm2Mdl.AttrProtect,
			Value:        dm2Mdl.AttrYes,
			Dmids:        ids,
			Source:       oplog.SourceUp,
			OperatorType: oplog.OperatorUp,
		}
		if err = s.dmRPC.EditDMAttr(c, arg); err != nil {
			log.Error("s.dmRPC.EditDMAttr(%+v) error(%v)", arg, err)
			return
		}
	}
	return
}

// ProtectApplies 保护弹幕申请列表
func (s *Service) ProtectApplies(c context.Context, uid, aid int64, page int, sort string) (res *model.ApplyListResult, err error) {
	var (
		count int
		start int
	)
	if page < 1 {
		page = 1
	}
	res = &model.ApplyListResult{
		Pager: &model.Pager{},
		List:  make([]*model.Apply, 0, model.ProtectApplyLimit),
	}
	res.List, err = s.dao.ProtectApplies(c, uid, aid, sort)
	if err != nil {
		log.Error("s.dao.PaLs(%d) error(%v)", uid, err)
		return
	}
	count = len(res.List)
	res.Pager.Current = page
	res.Pager.Total = count / model.ProtectApplyLimit
	res.Pager.Size = model.ProtectApplyLimit
	res.Pager.TotalCount = count
	if count%model.ProtectApplyLimit != 0 {
		res.Pager.Total++
	}
	if count == 0 {
		res.List = make([]*model.Apply, 0, 1)
		return
	}
	start = (page - 1) * model.ProtectApplyLimit
	if start > count {
		start = 0
	}
	end := start + model.ProtectApplyLimit
	if end > count {
		end = count
	}
	res.List = res.List[start:end]
	aids := make([]int64, 0, len(res.List))
	uids := make([]int64, 0, len(res.List))
	for _, a := range res.List {
		aids = append(aids, a.AID)
		uids = append(uids, a.ApplyUID)
	}
	infosReply, err := s.accountSvc.Infos3(c, &account.MidsReq{
		Mids: uids,
	})
	if err != nil {
		log.Error("s.actSvc.Infos2(%v) error(%v)", uids, err)
		err = nil
	}
	archives := s.archiveInfos(c, aids)
	for _, a := range res.List {
		v, ok := archives[a.AID]
		if ok {
			a.Pic = v.Pic
			a.Title = v.Title
		}
		u, ok := infosReply.GetInfos()[a.ApplyUID]
		if ok {
			a.Uname = u.GetName()
		}
	}
	return
}

// PaVideoLs 被申请保护弹幕的视频列表
func (s *Service) PaVideoLs(c context.Context, uid int64) (res []*model.Video, err error) {
	aids, err := s.dao.ProtectAids(c, uid)
	if err != nil {
		log.Error("s.dao.ProtectArchives(%d) error(%v)", uid, err)
		return
	}
	archives := s.archiveInfos(c, aids)
	res = make([]*model.Video, 0, len(aids))
	for _, aid := range aids {
		a := new(model.Video)
		v, ok := archives[aid]
		a.Aid = aid
		if ok {
			a.Title = v.Title
		} else {
			a.Title = ""
		}
		res = append(res, a)
	}
	return
}

// sendProtectNotifyToUp 发送申请保护弹幕通知给up主
func (s *Service) sendProtectNotifyToUp(c context.Context) (err error) {
	if time.Now().Format("15") != "20" {
		return
	}
	lk, err := s.dao.PaLock(c, "up")
	if err != nil {
		log.Error("s.dao.PaLock() error(%v)", err)
		return
	}
	if lk != 1 {
		return
	}
	uids, err := s.dao.ProtectApplyStatistics(c)
	if err != nil {
		log.Error("s.dao.PaStat() error(%v)", err)
		return
	}
	if len(uids) < 1 {
		return
	}
	m, err := s.dao.PaNoticeClose(c, uids)
	if err != nil {
		log.Error("s.dao.PaNoticeClose(%v) error(%v)", uids, err)
		return
	}
	if len(m) > 0 {
		for k, v := range uids {
			if _, ok := m[v]; ok {
				uids = append(uids[:k], uids[k+1:]...)
			}
		}
	}
	s.dao.SendNotify(c, _notifyUpTitle, _notifyUpContent, uids)
	return
}

// sendProtectNotifyToUser 发送申请保护弹幕处理结果给申请用户
func (s *Service) sendProtectNotifyToUser(c context.Context) {
	if time.Now().Format("15") != "22" {
		return
	}
	incr, err := s.dao.PaLock(c, "user")
	if err != nil {
		log.Error("s.dao.PaLock() error(%v)", err)
		return
	}
	if incr != 1 {
		return
	}
	stats, err := s.dao.PaUsrStat(c)
	if err != nil {
		log.Error("s.dao.PaStat() error(%v)", err)
		return
	}
	aids := make([]int64, 0, len(stats))
	for _, stat := range stats {
		aids = append(aids, stat.Aid)
	}
	archives := s.archiveInfos(c, aids)
	userStats := make(map[int64]map[int64]*model.ApplyUserNotify)
	now := time.Now().Unix()
	untreated := make(map[int64]int)
	for _, stat := range stats {
		m, ok := userStats[stat.Aid]
		if !ok {
			m = make(map[int64]*model.ApplyUserNotify)
			userStats[stat.Aid] = m
		}
		n, ok := m[stat.UID]
		if !ok {
			n = &model.ApplyUserNotify{}
			m[stat.UID] = n
		}
		if stat.Status == 1 {
			n.Protect++
		} else {
			n.Unprotect++
			if stat.Status == -1 && (now-stat.Ctime.Unix()) > 2*24*3600 {
				untreated[stat.UID]++
			}
		}
	}
	for k, v := range untreated {
		s.dao.SendNotify(c, _notifyUsrTitle, fmt.Sprintf(_pa48, v), []int64{k})
	}
	for aid, m := range userStats {
		archive, ok := archives[aid]
		if !ok {
			continue
		}
		for uid, stat := range m {
			var content string
			if stat.Protect > 0 && stat.Unprotect == 0 {
				content = fmt.Sprintf(_paa, archive.Title, archive.Aid)
			}
			if stat.Protect > 0 && stat.Unprotect > 0 {
				content = fmt.Sprintf(_pap, archive.Title, archive.Aid)
			}
			if content == "" {
				continue
			}
			s.dao.SendNotify(c, _notifyUsrTitle, content, []int64{uid})
		}
	}
}
