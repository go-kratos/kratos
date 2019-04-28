package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"go-common/app/admin/main/reply/model"
	accmdl "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	rlmdl "go-common/app/service/main/relation/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

func (s *Service) reply(c context.Context, oid, rpID int64) (rp *model.Reply, err error) {
	if rp, err = s.dao.Reply(c, oid, rpID); err != nil {
		return
	}
	if rp == nil {
		err = ecode.ReplyNotExist
		return
	}
	if rp.Content, err = s.dao.ReplyContent(c, oid, rpID); err != nil {
		return
	}
	if rp.Content == nil {
		err = ecode.ReplyNotExist
	}
	return
}

func (s *Service) replies(c context.Context, oids, rpIDs []int64) (res map[int64]*model.Reply, err error) {
	res, missIDs, err := s.dao.RepliesCache(c, rpIDs)
	if err != nil {
		return
	}
	if len(missIDs) > 0 {
		var (
			rps      map[int64]*model.Reply
			rcs      map[int64]*model.ReplyContent
			miss     []*model.Reply
			missOids []int64
		)
		for _, missID := range missIDs {
			for i := range rpIDs {
				if rpIDs[i] == missID {
					missOids = append(missOids, oids[i])
					break
				}
			}
		}
		if rps, err = s.dao.Replies(c, missOids, missIDs); err != nil {
			return
		}
		if rcs, err = s.dao.ReplyContents(c, missOids, missIDs); err != nil {
			return
		}
		for id, rp := range rps {
			rp.Content = rcs[id]
			res[id] = rp
			miss = append(miss, rp)
		}
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.AddReplyCache(ctx, miss...)
		})
	}
	return
}

// ReplySearch return reply result from search.
func (s *Service) ReplySearch(c context.Context, sp *model.SearchParams, page, pageSize int64) (res *model.SearchResult, err error) {
	if res, err = s.dao.SearchReplyV3(c, sp, page, pageSize); err != nil {
		log.Error("s.dao.SearchReplyV3(%+v,%d,%d) error(%v) ", sp, page, pageSize, err)
		return
	}
	adMap := make(map[int64]*model.SearchAdminLog)
	var ids []int64
	for i := range res.Result {
		ids = append(ids, res.Result[i].ID)
	}
	if adres, err := s.dao.SearchAdminLog(c, ids); err == nil {
		for _, data := range adres {
			adMap[data.ReplyID] = data
		}
	}
	filters := make(map[int64]string, len(res.Result))
	links := make(map[int64]string, len(res.Result))
	titles := make(map[int64]string)
	var mids []int64
	for _, data := range res.Result {
		if log, ok := adMap[data.ID]; ok && log != nil {
			data.AdminID = log.AdminID
			data.AdminName = log.AdminName
			data.OpCtime = log.CTime
			data.Opremark = log.Remark
			data.Opresult = log.Result
		}
		mids = append(mids, data.Mid)
		data.OidStr = strconv.FormatInt(data.Oid, 10)
		// reply filtered
		if len(data.Attr) > 0 {
			for _, attr := range data.Attr {
				if attr == 4 {
					filters[data.ID] = data.Message
				}
			}
		}
		// show title for top reply
		if sp.Attr == "1" && data.Title == "" {
			var link string
			data.Title, link, _ = s.TitleLink(c, data.Oid, int32(data.Type))
			data.RedirectURL = fmt.Sprintf("%s#reply%d", link, data.ID)
		} else {
			links[data.Oid] = ""
			if int32(data.Type) == model.SubTypeArchive {
				titles[data.Oid] = ""
			}
		}
	}
	ip := metadata.String(c, metadata.RemoteIP)
	stasMap, err := s.relationSvc.Stats(c, &rlmdl.ArgMids{Mids: mids, RealIP: ip})
	if err == nil {
		for i, data := range res.Result {
			if stat, ok := stasMap[data.Mid]; ok {
				res.Result[i].Stat = stat
			}
		}
	} else {
		log.Error("relationSvc.Stats error(%v)", err)
	}
	s.linkByOids(c, links, sp.Type)
	s.titlesByOids(c, titles)
	s.dao.FilterContents(c, filters)
	for _, data := range res.Result {
		if content := filters[data.ID]; content != "" {
			data.Message = content
		}
		if data.RedirectURL == "" {
			if link := links[data.Oid]; link != "" {
				data.RedirectURL = fmt.Sprintf("%s#reply%d", link, data.ID)
			}
		}
		if int32(data.Type) == model.SubTypeArchive && data.Title == "" {
			if title := titles[data.Oid]; title != "" {
				data.Title = title
			}
		}
	}
	return
}

func (s *Service) titlesByOids(c context.Context, titles map[int64]string) (err error) {
	var aids []int64
	for oid := range titles {
		aids = append(aids, oid)
	}
	arg := &arcmdl.ArgAids2{
		Aids: aids,
	}
	var m map[int64]*api.Arc
	m, err = s.arcSrv.Archives3(c, arg)
	for oid := range m {
		titles[oid] = m[oid].Title
	}
	return
}

func (s *Service) linkByOids(c context.Context, oids map[int64]string, typ int32) (err error) {
	if len(oids) == 0 {
		return
	}
	if typ == model.SubTypeActivity {
		err = s.dao.TopicsLink(c, oids, false)
	} else {
		for oid := range oids {
			var link string
			switch typ {
			case model.SubTypeTopic:
				link = fmt.Sprintf("https://www.bilibili.com/topic/%d.html", oid)
			case model.SubTypeArchive:
				link = fmt.Sprintf("https://www.bilibili.com/video/av%d", oid)
			case model.SubTypeForbiden:
				link = fmt.Sprintf("https://www.bilibili.com/blackroom/ban/%d", oid)
			case model.SubTypeNotice:
				link = fmt.Sprintf("https://www.bilibili.com/blackroom/notice/%d", oid)
			case model.SubTypeActArc:
				_, link, err = s.dao.ActivitySub(c, oid)
				if err != nil {
					return
				}
			case model.SubTypeArticle:
				link = fmt.Sprintf("https://www.bilibili.com/read/cv%d", oid)
			case model.SubTypeLiveVideo:
				link = fmt.Sprintf("https://vc.bilibili.com/video/%d", oid)
			case model.SubTypeLiveAct:
				_, link, err = s.dao.LiveActivityTitle(c, oid)
				if err != nil {
					return
				}
			case model.SubTypeLivePicture:
				link = fmt.Sprintf("https://h.bilibili.com/ywh/%d", oid)
			case model.SubTypeCredit:
				link = fmt.Sprintf("https://www.bilibili.com/judgement/case/%d", oid)
			case model.SubTypeDynamic:
				link = fmt.Sprintf("https://t.bilibili.com/%d", oid)
			default:
				return
			}
			oids[oid] = link
		}
	}
	return
}

// AdminEditReply edit reply content by admin.
func (s *Service) AdminEditReply(c context.Context, adminID int64, adName string, oid, rpID int64, tp int32, msg, remark string) (err error) {
	rp, err := s.reply(c, oid, rpID)
	if err != nil {
		return
	}
	if rp.IsDeleted() {
		err = ecode.ReplyDeleted
		return
	}
	now := time.Now()
	if _, err = s.dao.UpReplyContent(c, oid, rpID, msg, now); err != nil {
		log.Error("s.content.UpMessage(%d, %d, %s, %v), err is (%v)", oid, rpID, msg, now, err)
		return
	}
	if err = s.dao.DelReplyCache(c, rpID); err != nil {
		log.Error("dao.AddReplyCache(%+v,%s) rpid(%d) error(%v)", rp, msg, err)
	}
	s.addAdminLog(c, rp.Oid, rp.ID, adminID, rp.Type, model.AdminIsNew, model.AdminIsNotReport, model.AdminOperEdit, "已修改评论内容", remark, now)
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReply(ctx, map[int64]*model.Reply{rp.ID: rp}, rp.State)
	})
	report.Manager(&report.ManagerInfo{
		UID:      adminID,
		Uname:    adName,
		Business: 41,
		Type:     int(tp),
		Oid:      rp.Oid,
		Ctime:    now,
		Action:   model.ReportActionReplyEdit,
		Index: []interface{}{
			rp.ID,
			rp.State,
			rp.State,
		},
		Content: map[string]interface{}{"remark": remark},
	})
	return
}

// AddTop add a top reply.
func (s *Service) AddTop(c context.Context, adid int64, adName string, oid, rpID int64, typ int32, act uint32) (err error) {
	rp, err := s.reply(c, oid, rpID)
	if err != nil {
		return
	}
	if rp.IsFolded() {
		return ecode.ReplyFolded
	}
	if rp.Root != 0 {
		log.Error("add top reply illegal reply(oid:%v,type:%v,:rpID:%v) not root", oid, typ, rpID, err)
		return
	}
	sub, err := s.subject(c, oid, typ)
	if err != nil {
		log.Error("s.subject(%d,%d),err:%v", oid, typ)
		return
	}

	if act == model.AttrYes && sub.AttrVal(model.SubAttrTopAdmin) == model.AttrYes {
		err = ecode.ReplyHaveTop
		log.Error("Repeat to add top reply(%d,%d,%d,%d) ", rp.ID, rp.Oid, typ, sub.Attr)
		return
	}
	sub.AttrSet(act, model.SubAttrTopAdmin)
	err = sub.TopSet(rpID, 0, act)
	if err != nil {
		log.Error("sub.TopSet(%d,%d,%d) failed!err:=%v ", rp.ID, rp.Oid, 0, err)
		return
	}
	rp.AttrSet(act, model.AttrTopAdmin)
	now := time.Now()
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		return
	}
	if _, err = s.dao.TxUpReplyAttr(tx, oid, rpID, rp.Attr, now); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.TxUpSubAttr(tx, oid, typ, sub.Attr, now); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.TxUpSubMeta(tx, sub.Oid, sub.Type, sub.Meta, now); err != nil {
		tx.Rollback()
		log.Error("dao.TxUpMeta(oid:%d,tp:%d) err(%v)", sub.Oid, sub.Type, err)
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	if act == model.AttrYes {
		s.dao.DelIndexBySort(c, rp, model.SortByCount)
		s.dao.DelIndexBySort(c, rp, model.SortByLike)
	} else if act == model.AttrNo && rp.IsNormal() {
		s.addReplyIndex(c, rp)
	}
	s.dao.AddTopCache(c, rp)
	s.dao.AddReplyCache(c, rp)
	s.dao.DelSubjectCache(c, rp.Oid, rp.Type)
	report.Manager(&report.ManagerInfo{
		UID:      adid,
		Uname:    adName,
		Business: 41,
		Type:     int(typ),
		Oid:      oid,
		Ctime:    now,
		Action:   model.ReportActionReplyTop,
		Index:    []interface{}{sub.Mid, act, rpID},
	})
	if act == model.AttrYes {
		s.pubEvent(c, "top", 0, sub, rp, nil)
	} else if act == model.AttrNo {
		s.pubEvent(c, "untop", 0, sub, rp, nil)
	}
	//add admin log and search log
	if act == model.AttrYes {
		s.addAdminLog(c, rp.Oid, rp.ID, adid, rp.Type, model.AdminIsNew, model.AdminIsNotReport, model.AdminOperSubTop, "管理员置顶评论", "", time.Now())
		s.cache.Do(c, func(ctx context.Context) {
			s.pubSearchReply(ctx, map[int64]*model.Reply{rp.ID: rp}, rp.State)
		})
	} else {
		s.addAdminLog(c, rp.Oid, rp.ID, adid, rp.Type, model.AdminIsNew, model.AdminIsNotReport, model.AdminOperSubTop, "管理员取消置顶评论", "", time.Now())
		s.cache.Do(c, func(ctx context.Context) {
			s.pubSearchReply(ctx, map[int64]*model.Reply{rp.ID: rp}, model.StateNormal)
		})
	}
	return
}

// CallbackDeleteReply delete reply by admin.
func (s *Service) CallbackDeleteReply(ctx context.Context, adminID int64, oid, rpID int64, ftime int64, typ int32, moral int32, adminName, remark string, reason, freason int32) (err error) {
	now := time.Now()
	sub, rp, err := s.delReply(ctx, oid, rpID, model.StateDelAdmin, now)
	if err != nil {
		if ecode.ReplyDeleted.Equal(err) && rp.IsDeleted() {
			err = nil
		} else {
			log.Error("delReply(%d,%d) error(%v)", oid, rpID, err)
			return err
		}
	}
	s.delCache(ctx, sub, rp)
	s.cache.Do(ctx, func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		var rpt *model.Report
		if rpt, _ = s.dao.Report(ctx, oid, rpID); rpt != nil {
			if rpt.State == model.ReportStateNew || rpt.State == model.ReportStateNew2 {
				rpt.MTime = xtime.Time(now.Unix())
				if rpt.State == model.ReportStateNew {
					rpt.State = model.ReportStateDelete1
				} else if rpt.State == model.ReportStateNew2 {
					rpt.State = model.ReportStateDelete2
				}
				if _, err = s.dao.UpReportsState(ctx, []int64{rpt.Oid}, []int64{rpt.RpID}, rpt.State, now); err != nil {
					log.Error("s.dao.UpdateReport(%+v) error(%v)", rpt, err)
				}
				state := model.StateDelAdmin
				s.pubSearchReport(ctx, map[int64]*model.Report{rpt.RpID: rpt}, &state)
			}
		}
		s.pubEvent(context.Background(), model.EventReportDel, rpt.Mid, sub, rp, rpt)
		report.Manager(&report.ManagerInfo{
			UID:      adminID,
			Uname:    adminName,
			Business: 41,
			Type:     int(typ),
			Oid:      rp.Oid,
			Ctime:    now,
			Action:   model.ReportActionReplyDel,
			Index: []interface{}{
				rp.ID,
				rp.State,
				model.StateDelAdmin,
			},
			Content: map[string]interface{}{
				"moral":   moral,
				"notify":  false,
				"ftime":   ftime,
				"freason": freason,
				"reason":  reason,
				"remark":  remark,
			},
		})
		rps := make(map[int64]*model.Reply)
		rps[rp.ID] = rp
		s.addAdminLogs(ctx, rps, adminID, typ, model.AdminIsNew, model.AdminIsReport, model.AdminOperDelete, fmt.Sprintf("已删除并封禁%s/扣除%d节操", forbidResult(ftime), moral), remark, now)
		s.pubSearchReply(ctx, rps, model.StateDelAdmin)
	})
	return
}

// AdminDeleteReply delete reply by admin.
func (s *Service) AdminDeleteReply(c context.Context, adminID int64, oids, rpIDs []int64, ftime int64, typ int32, moral int32, notify bool, adminName, remark string, reason, freason int32) (err error) {
	err = s.adminDeleteReply(c, adminID, oids, rpIDs, ftime, typ, moral, notify, adminName, remark, reason, freason)
	return
}

func (s *Service) adminDeleteReply(c context.Context, adminID int64, oids, rpIDs []int64, ftime int64, typ int32, moral int32, notify bool, adminName, remark string, reason, freason int32) (err error) {
	var (
		lk  sync.Mutex
		rps = make(map[int64]*model.Reply)
		now = time.Now()
	)
	wg := errgroup.Group{}
	wg.GOMAXPROCS(4)
	for idx := range oids {
		i := idx
		wg.Go(func() (err error) {
			var sub *model.Subject
			var rp *model.Reply
			// 针对大忽悠事件 特殊用户删除评论不让删
			var (
				tp     int32
				ok     bool
				exsits bool
			)
			tp, ok = s.oids[oids[i]]
			if ok && tp == typ {
				_, exsits = s.ads[adminName]
				if exsits {
					return
				}
			}

			sub, rp, err = s.delReply(c, oids[i], rpIDs[i], model.StateDelAdmin, now)
			if err != nil {
				if ecode.ReplyDeleted.Equal(err) && rp.IsDeleted() {
					err = nil
				} else {
					log.Error("delReply(%d,%d) error(%v)", oids[i], rpIDs[i], err)
					return err
				}
			}
			if rp.IsFolded() {
				s.marker.Do(c, func(ctx context.Context) {
					s.handleFolded(ctx, rp)
				})
			}
			s.delCache(c, sub, rp)
			s.pubEvent(c, "reply_del", 0, sub, rp, nil)
			rpt, _ := s.dao.Report(c, oids[i], rpIDs[i])
			if rpt != nil {
				if rpt.State == model.ReportStateNew || rpt.State == model.ReportStateNew2 {
					rpt.MTime = xtime.Time(now.Unix())
					if rpt.State == model.ReportStateNew {
						rpt.State = model.ReportStateDelete1
					} else if rpt.State == model.ReportStateNew2 {
						rpt.State = model.ReportStateDelete2
					}
					if _, err = s.dao.UpReportsState(c, []int64{rpt.Oid}, []int64{rpt.RpID}, rpt.State, now); err != nil {
						log.Error("s.dao.UpdateReport(%+v) error(%v)", rpt, err)
					}
				}
			}
			s.cache.Do(c, func(ctx context.Context) {
				// 针对大忽悠事件的特殊推送
				if _, ok := s.ads[adminName]; ok {
					if e := s.NotifyTroll(ctx, rp.Mid); e != nil {
						log.Warn("notify-troll error (%v)", e)
					}
				}

				s.dao.DelReport(ctx, rp.Oid, rp.ID)
				if rpt != nil {
					rpt.ReplyCtime = rp.CTime
					state := model.StateDelAdmin
					s.pubSearchReport(ctx, map[int64]*model.Report{rpt.RpID: rpt}, &state)
				}
				s.moralAndNotify(ctx, rp, moral, notify, rp.Mid, adminID, adminName, remark, reason, freason, ftime, false)
			})
			lk.Lock()
			rps[rp.ID] = rp
			lk.Unlock()
			report.Manager(&report.ManagerInfo{
				UID:      adminID,
				Uname:    adminName,
				Business: 41,
				Type:     int(typ),
				Oid:      rp.Oid,
				Ctime:    now,
				Action:   model.ReportActionReplyDel,
				Index: []interface{}{
					rp.ID,
					rp.State,
					model.StateDelAdmin,
				},
				Content: map[string]interface{}{
					"moral":   moral,
					"notify":  notify,
					"ftime":   ftime,
					"freason": freason,
					"reason":  reason,
					"remark":  remark,
				},
			})
			return nil
		})
	}
	if err = wg.Wait(); err != nil {
		return
	}
	s.addAdminLogs(c, rps, adminID, typ, model.AdminIsNew, model.AdminIsNotReport, model.AdminOperDelete, fmt.Sprintf("已删除并封禁%s/扣除%d节操", forbidResult(ftime), moral), remark, now)
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReply(ctx, rps, model.StateDelAdmin)
	})
	return
}

// AdminRecoverReply recover reply by admin.
func (s *Service) AdminRecoverReply(c context.Context, adminID int64, adName string, oid, rpID int64, typ int32, remark string) (err error) {
	rp, err := s.reply(c, oid, rpID)
	if err != nil {
		log.Error("s.reply(%d,%d) error(%v)", oid, rpID, err)
		return
	}
	now := time.Now()
	var sub *model.Subject
	if sub, rp, err = s.recReply(c, rp.Oid, rp.ID, model.StateNormal, now); err != nil {
		log.Error("s.recReply(%d,%d) error(%v)", rp.Oid, rp.ID, err)
		return
	}
	s.addAdminLog(c, rp.Oid, rp.ID, adminID, rp.Type, model.AdminIsNew, model.AdminIsNotReport, model.AdminOperRecover, "已恢复评论", remark, now)
	s.pubEvent(c, "reply_recover", 0, sub, rp, nil)
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReply(ctx, map[int64]*model.Reply{rp.ID: rp}, model.StateNormal)
	})
	report.Manager(&report.ManagerInfo{
		UID:      adminID,
		Uname:    adName,
		Business: 41,
		Type:     int(typ),
		Oid:      rp.Oid,
		Ctime:    now,
		Action:   model.ReportActionReplyRecover,
		Index: []interface{}{
			rp.ID,
			rp.State,
			model.StateNormal,
		},
		Content: map[string]interface{}{"remark": remark},
	})
	return
}

// AdminPassReply recover reply by admin.
func (s *Service) AdminPassReply(c context.Context, adid int64, adName string, oids, rpIDs []int64, typ int32, remark string) (err error) {
	s.adminPassReply(c, adid, adName, oids, rpIDs, typ, remark)
	return
}

func (s *Service) adminPassReply(c context.Context, adid int64, adName string, oids, rpIDs []int64, typ int32, remark string) (err error) {
	now := time.Now()
	rps, err := s.replies(c, oids, rpIDs)
	if err != nil {
		return
	}
	wg, ctx := errgroup.WithContext(c)
	for _, m := range rps {
		rp := m
		wg.Go(func() (err error) {
			if rp.State == model.StatePending {
				var sub *model.Subject
				if sub, rp, err = s.recReply(ctx, rp.Oid, rp.ID, model.StateNormal, now); err != nil {
					return
				}
				s.dao.DelAuditIndex(ctx, rp)
				s.pubEvent(c, "reply_recover", 0, sub, rp, nil)
			} else {
				var (
					tx   *sql.Tx
					rows int64
				)
				if tx, err = s.dao.BeginTran(ctx); err != nil {
					return
				}
				if rows, err = s.dao.TxUpdateReplyState(tx, rp.Oid, rp.ID, model.StateNormal, now); err != nil || rows == 0 {
					log.Error("dao.Reply.TxUpdateReplyState(%v,%d) error(%v)", rp, model.StateNormal, err)
					tx.Rollback()
					return
				}
				if rp.State == model.StateMonitor {
					if _, err = s.dao.TxSubDecrMCount(tx, rp.Oid, rp.Type, now); err != nil {
						log.Error("dao.Reply.TxSubDecrMCount(%v) error(%v)", rp, err)
						tx.Rollback()
						return
					}
				}
				if err = tx.Commit(); err != nil {
					log.Error("tx.Commit error(%v)", err)
					return
				}
				if err = s.dao.DelReplyCache(ctx, rp.ID); err != nil {
					log.Error("s.dao.DelReplyCache(%d,%d) error(%v)", rp.Oid, rp.ID, err)
				}
			}
			report.Manager(&report.ManagerInfo{
				UID:      adid,
				Uname:    adName,
				Business: 41,
				Type:     int(typ),
				Oid:      rp.Oid,
				Ctime:    now,
				Action:   model.ReportActionReplyPass,
				Index: []interface{}{
					rp.ID,
					rp.State,
					model.StateNormal,
				},
				Content: map[string]interface{}{"remark": remark},
			})
			return
		})
	}
	if err = wg.Wait(); err != nil {
		return
	}
	s.addAdminLogs(c, rps, adid, typ, model.AdminIsNew, model.AdminIsNotReport, model.AdminOperPass, "已通过评论", remark, now)
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReply(ctx, rps, model.StateNormal)
	})
	return
}

// addReplyIndex add reply index to redis.
func (s *Service) addReplyIndex(c context.Context, rp *model.Reply) (err error) {
	var ok bool
	if rp.IsRoot() {
		if ok, err = s.dao.ExpireIndex(c, rp.Oid, rp.Type, model.SortByFloor); err == nil && ok {
			if err = s.dao.AddFloorIndex(c, rp); err != nil {
				log.Error("d.AddFloorIndex(%d,%d) error(%v)", rp.Oid, rp.Type, err)
			}
		}
		if ok, err = s.dao.ExpireIndex(c, rp.Oid, rp.Type, model.SortByCount); err == nil && ok {
			if err = s.dao.AddCountIndex(c, rp); err != nil {
				log.Error("s.AddCountIndex(%d,%d) error(%v)", rp.Oid, rp.Type, err)
			}
		}
		if ok, err = s.dao.ExpireIndex(c, rp.Oid, rp.Type, model.SortByLike); err == nil && ok {
			rpt, _ := s.dao.Report(c, rp.Oid, rp.ID)
			if err = s.dao.AddLikeIndex(c, rp, rpt); err != nil {
				log.Error("d.AddLikeIndex(%d,%d) error(%v)", rp.Oid, rp.Type, err)
			}
		}
	} else {
		if ok, err = s.dao.ExpireNewChildIndex(c, rp.Root); err == nil && ok {
			if err = s.dao.AddNewChildIndex(c, rp); err != nil {
				log.Error("d.AddFloorRootIndex(%d) error(%v)", rp.Root, err)
			}
		}
	}
	return
}

func (s *Service) recReply(c context.Context, oid, rpID int64, state int32, now time.Time) (sub *model.Subject, rp *model.Reply, err error) {
	if sub, rp, err = s.tranRecover(c, oid, rpID, state, now); err != nil {
		return
	}
	if rp.Content, err = s.dao.ReplyContent(c, oid, rpID); err != nil {
		return
	}
	if rp.Content == nil {
		err = ecode.ReplyNotExist
		return
	}
	if !rp.IsRoot() {
		if err = s.dao.DelReplyCache(c, rp.Root); err != nil {
			log.Error("s.dao.DelReplyCache(%d,%d) error(%v)", oid, rpID, err)
		}
	}
	if err = s.dao.DelReplyCache(c, rpID); err != nil {
		log.Error("s.dao.DelReplyCache(%d,%d) error(%v)", oid, rpID, err)
	}
	if err = s.addReplyIndex(c, rp); err != nil {
		log.Error("s.dao.DelReplyIndex(%d,%d) error(%v)", oid, rpID, err)
	}
	if err = s.dao.AddSubjectCache(c, sub); err != nil {
		log.Error("s.dao.DelSubjectCache(%+v) error(%v)", sub, err)
	}
	s.dao.SendStats(c, sub.Type, sub.Oid, sub.ACount)
	return
}

func (s *Service) delReply(c context.Context, oid, rpID int64, state int32, now time.Time) (sub *model.Subject, rp *model.Reply, err error) {
	if sub, rp, err = s.tranDel(c, oid, rpID, state, now); err != nil {
		if ecode.ReplyDeleted.Equal(err) && rp.IsDeleted() {
			if rp.Content, err = s.dao.ReplyContent(c, oid, rpID); err != nil {
				return
			} else if rp.Content == nil {
				err = ecode.ReplyNotExist
				return
			}
			err = ecode.ReplyDeleted
		}
		return
	}
	if rp.Content, err = s.dao.ReplyContent(c, oid, rpID); err != nil {
		return
	}
	if rp.Content == nil {
		err = ecode.ReplyNotExist
		return
	}
	return
}

func (s *Service) delCache(c context.Context, sub *model.Subject, rp *model.Reply) (err error) {
	if !rp.IsRoot() {
		if err = s.dao.DelReplyCache(c, rp.Root); err != nil {
			log.Error("s.dao.DelReplyCache(%d,%d) error(%v)", rp.Oid, rp.ID, err)
		}
	}
	if err = s.dao.DelReplyCache(c, rp.ID); err != nil {
		log.Error("s.dao.DelReplyCache(%d,%d) error(%v)", rp.Oid, rp.ID, err)
	}
	if err = s.dao.DelReplyIndex(c, rp); err != nil {
		log.Error("s.dao.DelReplyIndex(%d,%d) error(%v)", rp.Oid, rp.ID, err)
	}
	if err = s.dao.AddSubjectCache(c, sub); err != nil {
		log.Error("s.dao.DelSubjectCache(%+v) error(%v)", sub, err)
	}
	if rp.AttrVal(model.AttrTopAdmin) == model.AttrYes {
		s.dao.DelTopCache(c, rp.Oid, model.SubAttrTopAdmin)
	}
	if rp.AttrVal(model.AttrTopUpper) == model.AttrYes {
		s.dao.DelTopCache(c, rp.Oid, model.SubAttrTopUpper)
	}
	s.dao.SendStats(c, sub.Type, sub.Oid, sub.ACount)
	return
}

func (s *Service) tranRecover(c context.Context, oid, rpID int64, state int32, now time.Time) (sub *model.Subject, rp *model.Reply, err error) {
	var (
		rootRp *model.Reply
		count  int32
	)
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		return
	}
	if rp, err = s.dao.TxReplyForUpdate(tx, oid, rpID); err != nil {
		tx.Rollback()
		err = fmt.Errorf("s.dao.Reply(%d,%d) error(%v) ", oid, rpID, err)
		return
	}
	if rp == nil {
		tx.Rollback()
		err = ecode.ReplyNotExist
		return
	} else if rp.IsNormal() {
		tx.Rollback()
		err = ecode.ReplyActioned
		return
	}
	rows, err := s.dao.TxUpdateReplyState(tx, rp.Oid, rp.ID, state, now)
	if err != nil || rows == 0 {
		tx.Rollback()
		err = fmt.Errorf("error(%v) or rows(%d)", err, rows)
		return
	}
	rp.MTime = xtime.Time(now.Unix())
	if rp.IsRoot() {
		count = rp.RCount + 1
	} else {
		if rootRp, err = s.dao.TxReply(tx, rp.Oid, rp.Root); err != nil {
			tx.Rollback()
			return
		}
		count = 1
	}
	if rp.IsRoot() {
		rows, err = s.dao.TxIncrSubRCount(tx, rp.Oid, rp.Type, now)
	} else {
		rows, err = s.dao.TxIncrReplyRCount(tx, rp.Oid, rp.Root, now)
	}
	if err != nil || rows == 0 {
		tx.Rollback()
		err = fmt.Errorf("tranRecover increase count error(%v) or rows(%d)", err, rows)
		return
	}
	if rp.IsRoot() || (rootRp != nil && rootRp.IsNormal()) {
		if rows, err = s.dao.TxIncrSubACount(tx, rp.Oid, rp.Type, count, now); err != nil || rows == 0 {
			tx.Rollback()
			err = fmt.Errorf("TxIncrSubACount error(%v) or rows(%d)", err, rows)
			return
		}
	}
	if rp.State == model.StatePending {
		if rows, err = s.dao.TxSubDecrMCount(tx, rp.Oid, rp.Type, now); err != nil {
			tx.Rollback()
			err = fmt.Errorf("TxSubDecrMCount error(%v)", err)
			return
		}
	}
	if sub, err = s.dao.TxSubject(tx, rp.Oid, rp.Type); err != nil || sub == nil {
		tx.Rollback()
		err = fmt.Errorf(" s.dao.TxSubject(%d,%d) or rows(%d)", rp.Oid, rp.Type, rows)
		return
	}
	err = tx.Commit()
	return
}

func (s *Service) tranDel(c context.Context, oid, rpID int64, state int32, now time.Time) (sub *model.Subject, rp *model.Reply, err error) {
	var (
		count     int32
		rootReply *model.Reply
	)
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		return
	}
	rp, err = s.dao.TxReplyForUpdate(tx, oid, rpID)
	if err != nil {
		tx.Rollback()
		err = fmt.Errorf("s.dao.TxReplyForUpdate(%d,%d) error(%v) ", oid, rpID, err)
		return
	}
	if rp == nil {
		err = ecode.NothingFound
		return
	}
	if rp.AttrVal(model.AttrTopAdmin) == 1 || rp.IsDeleted() {
		if sub, err = s.dao.TxSubject(tx, rp.Oid, rp.Type); err != nil || sub == nil {
			tx.Rollback()
			err = fmt.Errorf("s.dao.TxSubject(%d,%d)  error(%v)", rp.Oid, rp.Type, err)
			return
		}
		tx.Rollback()
		err = ecode.ReplyDeleted
		return
	}
	rows, err := s.dao.TxUpdateReplyState(tx, oid, rpID, state, now)
	if err != nil || rows == 0 {
		tx.Rollback()
		err = fmt.Errorf("s.dao.TxUpdateReplyState(%+v) rows:%d error(%v)", rp, rows, err)
		return
	}
	rp.MTime = xtime.Time(now.Unix())
	if rp.IsNormal() {
		if rp.IsRoot() {
			count = rp.RCount + 1
			if rows, err = s.dao.TxSubDecrACount(tx, rp.Oid, rp.Type, count, now); err != nil || rows == 0 {
				tx.Rollback()
				err = fmt.Errorf("s.dao.TxSubDecrACount(%+v) rows:%d error(%v)", rp, rows, err)
				return
			}
			rows, err = s.dao.TxDecrSubRCount(tx, rp.Oid, rp.Type, now)
			if err != nil || rows == 0 {
				tx.Rollback()
				err = fmt.Errorf("TxDecrReplyRCount(%+v) rows:%d error(%v)", rp, rows, err)
				return
			}
		} else {
			if rootReply, err = s.dao.TxReplyForUpdate(tx, rp.Oid, rp.Root); err != nil {
				tx.Rollback()
				return
			}
			if rootReply != nil {
				if rootReply.IsNormal() {
					if rows, err = s.dao.TxSubDecrACount(tx, rp.Oid, rp.Type, 1, now); err != nil || rows == 0 {
						tx.Rollback()
						err = fmt.Errorf("s.dao.TxSubDecrACount(%+v) rows:%d error(%v)", rp, rows, err)
						return
					}
				}
				_, err = s.dao.TxDecrReplyRCount(tx, rp.Oid, rp.Root, now)
				if err != nil {
					tx.Rollback()
					err = fmt.Errorf("TxDecrReplyRCount(%+v) error(%v)", rp, err)
					return
				}
			}
		}
	}
	if sub, err = s.dao.TxSubject(tx, rp.Oid, rp.Type); err != nil || sub == nil {
		tx.Rollback()
		err = fmt.Errorf("s.dao.TxSubject(%d,%d) rows:%d error(%v)", rp.Oid, rp.Type, rows, err)
		return
	}
	if rp.State == model.StatePending || rp.State == model.StateMonitor {
		if _, err = s.dao.TxSubDecrMCount(tx, rp.Oid, rp.Type, now); err != nil {
			log.Error("dao.Reply.TxSubDecrMCount(%v) error(%v)", rp, err)
			tx.Rollback()
			return
		}
	}
	if rp.AttrVal(model.AttrTopUpper) == model.AttrYes {
		rp.AttrSet(model.AttrNo, model.AttrTopUpper)
		sub.AttrSet(model.AttrNo, model.SubAttrTopUpper)
		err = sub.TopSet(0, 1, 0)
		if err != nil {
			tx.Rollback()
			log.Error("sub.TopSet(%d,%d,%d) failed!err:=%v ", rp.ID, rp.Oid, 0, err)
			return
		}
		if _, err = s.dao.TxUpSubMeta(tx, sub.Oid, sub.Type, sub.Meta, now); err != nil {
			tx.Rollback()
			log.Error("dao.TxUpMeta(oid:%d,tp:%d) err(%v) rows(%d)", sub.Oid, sub.Type, err)
			return
		}
		if rows, err = s.dao.TxUpSubAttr(tx, sub.Oid, sub.Type, sub.Attr, now); err != nil || rows == 0 {
			tx.Rollback()
			err = fmt.Errorf("s.dao.TxUpSubAttr(%+v) rows:%d error(%v)", sub, rows, err)
			return
		}
		if rows, err = s.dao.TxUpReplyAttr(tx, rp.Oid, rp.ID, rp.Attr, now); err != nil || rows == 0 {
			tx.Rollback()
			err = fmt.Errorf("s.dao.TxUpReplyAttr(%+v) rows:%d error(%v)", rp, rows, err)
			return
		}
	}
	err = tx.Commit()
	return
}

// ReplyTopLog ReplyTopLog
func (s *Service) ReplyTopLog(c context.Context, sp model.LogSearchParam) (result *model.ReplyTopLogResult, err error) {
	result = &model.ReplyTopLogResult{
		Logs: []*model.ReplyTopLog{},
	}
	sp.Action = "top"
	reportData, err := s.dao.ReportLog(c, sp)
	if err != nil {
		return
	}
	result.Page = reportData.Page
	result.Sort = reportData.Sort
	result.Order = reportData.Order
	var mids []int64
	for _, data := range reportData.Result {
		mid := data.Index0
		action := data.Index1
		rpid := data.Index2
		title, link, _ := s.TitleLink(c, data.Oid, data.Type)
		var extra map[string]string
		if data.Content != "" {
			err = json.Unmarshal([]byte(data.Content), &extra)
			if err != nil {
				log.Error("MointorLog unmarshal failed!err:=%v", err)
				return
			}
		}
		result.Logs = append(result.Logs, &model.ReplyTopLog{
			Mid:         mid,
			AdminID:     data.AdminID,
			AdminName:   data.AdminName,
			Oid:         data.Oid,
			Type:        data.Type,
			Remark:      extra["remark"],
			CTime:       data.Ctime,
			RpID:        rpid,
			Action:      action,
			Title:       title,
			RedirectURL: link,
		})
		mids = append(mids, mid)

	}

	if len(mids) > 0 {
		var res *accmdl.InfosReply
		res, err = s.accSrv.Infos3(c, &accmdl.MidsReq{Mids: mids})
		if err != nil {
			log.Error(" s.accSrv.Infos3 (%v) error(%v)", mids, err)
			err = nil
			return
		}
		for _, log := range result.Logs {
			if user, ok := res.Infos[log.Mid]; ok {
				log.UserName = user.GetName()
			}
		}
	}
	return
}

// MarkAsSpam mark a reply(normal state) as spam.
func (s *Service) MarkAsSpam(c context.Context, oids, rpIDs []int64, adminID int64, adminName, remark string) (err error) {
	rps, err := s.replies(c, oids, rpIDs)
	if err != nil {
		return
	}
	now := time.Now()
	for rpID, rp := range rps {
		if rp.State == model.StateNormal && rp.AttrVal(model.AttrGarbage) == model.AttrNo {
			var (
				rows int64
				tx   *sql.Tx
			)
			tx, err = s.dao.BeginTran(c)
			if err != nil {
				return
			}
			_, err = s.dao.TxUpdateReplyState(tx, rp.Oid, rpID, model.StateGarbage, now)
			if err != nil {
				tx.Rollback()
				log.Error("s.dao.TxUpdateReplyState(%+v) error(%v)", rp, err)
				return
			}
			rp.AttrSet(model.AttrYes, model.AttrGarbage)
			_, err = s.dao.TxUpReplyAttr(tx, rp.Oid, rpID, rp.Attr, now)
			if err != nil {
				tx.Rollback()
				log.Error("s.dao.TxUpReplyAttr(%+v) rows:%d error(%v)", rp, rows, err)
				return
			}
			if err = tx.Commit(); err != nil {
				log.Error("tx.Commit(%+v) error(%v)", rp, err)
				return
			}
			if err = s.dao.DelReplyCache(c, rpID); err != nil {
				log.Error("s.dao.DelReplyCache(%+v) error(%v)", rp, err)
			}
			report.Manager(&report.ManagerInfo{
				UID:      adminID,
				Uname:    adminName,
				Business: 41,
				Type:     int(rp.Type),
				Oid:      rp.Oid,
				Ctime:    now,
				Action:   model.ReportActionReplyGarbage,
				Index: []interface{}{
					rp.ID,
					rp.State,
					model.StateGarbage,
				},
			})
			s.cache.Do(c, func(ctx context.Context) {
				s.addAdminLog(ctx, rp.Oid, rp.ID, adminID, rp.Type, model.AdminIsNew, model.AdminIsNotReport, model.AdminOperMarkSpam, "标记为垃圾", remark, now)
			})
		}
	}
	return
}

// ExportReply get exported replies by query
func (s *Service) ExportReply(c context.Context, oid, mid int64, tp int8, state string, startTime, endTime time.Time) (data [][]string, err error) {
	if data, err = s.dao.ExportReplies(c, oid, mid, tp, state, startTime, endTime); err != nil {
		log.Error("s.dao.ExportReplies(%d,%d,%d,%s,%v,%v) error(%v)", oid, mid, tp, state, startTime, endTime)
		return
	}
	return
}

// ReplyList ReplyList
func (s *Service) ReplyList(c context.Context, oids, rpids []int64) (res map[int64]*model.ReplyEx, err error) {
	res = make(map[int64]*model.ReplyEx, 0)
	replies, err := s.replies(c, oids, rpids)
	if err != nil {
		return
	}
	subjects := make(map[int32]map[int64]*model.Subject, 0)
	var roots []int64
	var rootoids []int64
	for _, data := range replies {
		sub := subjects[data.Type]
		if sub == nil {
			sub = make(map[int64]*model.Subject, 0)
			subjects[data.Type] = sub
		}
		sub[data.Oid] = nil
		if data.Root != 0 {
			rootoids = append(rootoids, data.Oid)
			roots = append(roots, data.Root)
		}
	}
	rootreplies, err := s.replies(c, rootoids, roots)
	if err != nil {
		return
	}
	for typ, data := range subjects {
		var ids []int64
		for oid := range data {
			ids = append(ids, oid)
		}
		sub, err := s.subjects(c, ids, typ)
		if err != nil {
			return res, err
		}
		subjects[typ] = sub
	}
	for _, data := range replies {
		var isUp bool
		var rootFloor int32
		sub := (subjects[data.Type])[data.Oid]
		if sub != nil && sub.Mid == data.Mid {
			isUp = true
		}
		if data.Root != 0 && rootreplies[data.Root] != nil {
			rootFloor = rootreplies[data.Root].Floor
		}
		res[data.ID] = &model.ReplyEx{*data, isUp, rootFloor}
	}
	return
}

// TopChildReply ...
func (s *Service) TopChildReply(c context.Context, rootID, childID, oid int64) (err error) {
	var (
		root  *model.Reply
		child *model.Reply
		ok    bool
	)
	rps, err := s.dao.Replies(c, []int64{oid, oid}, []int64{rootID, childID})
	if err != nil {
		return ecode.ReplyNotExist
	}
	if root, ok = rps[rootID]; !ok {
		return ecode.ReplyNotExist
	}
	if child, ok = rps[childID]; !ok {
		return ecode.ReplyNotExist
	}
	if root.Root != 0 || child.Root != root.ID {
		return ecode.ReplyNotExist
	}
	if ok, err = s.dao.ExpireNewChildIndex(c, rootID); !ok || err != nil {
		return ecode.ReplyNotExist
	}
	if err = s.dao.TopChildReply(c, rootID, childID); err != nil {
		return ecode.ReplyNotExist
	}
	return
}
