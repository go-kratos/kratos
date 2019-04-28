package service

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/interface/main/reply/model/reply"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	figmdl "go-common/app/service/main/figure/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"
	xtime "go-common/library/time"
)

const (
	_reportNormalCnt = 5
	_reportAddSecs   = 5
	_reportMaxSecs   = 180
)

// AddReport report a reply.
func (s *Service) AddReport(c context.Context, mid, oid, rpID int64, tp, reason int8, cont, platform string, build int64, buvid string) (cd int, err error) {
	var (
		r   *model.Reply
		now = time.Now()
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	// check subject
	if !model.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if err = model.CheckReportReason(reason); err != nil {
		return
	}
	cnt, err := s.dao.Redis.GetUserReportCnt(c, mid, now)
	if err != nil {
		log.Error("AddReport failed, replyCacheDao.GetUserReportCnt(%d), err is (%v)", mid, err)
		return
	}
	if cnt > _reportNormalCnt {
		var ttl int
		// When report count max 5 at one day, extra add 5s to user TTL at a time, and the max user TTL is 180s.
		// The last report time, we set Redis key TTL for one day seconds, use one day seconds sub current TTL value is user TTL.
		if ttl, err = s.dao.Redis.GetUserReportTTL(c, mid, now); err != nil {
			log.Error("AddReport failed, replyCacheDao.GetUserReportTTL(%d), err is (%v)", mid, err)
			return
		}
		if ttl >= 0 {
			uttl := s.oneDaySec - ttl
			maxttl := (cnt - _reportNormalCnt) * _reportAddSecs
			if maxttl > _reportMaxSecs {
				maxttl = _reportMaxSecs
			}
			if uttl < maxttl {
				cd = maxttl - uttl
				err = ecode.ReplyReportDeniedAsCD
				return
			}
		}
	}
	if r, err = s.Reply(c, oid, tp, rpID); err != nil || r == nil || r.Content == nil || r.IsDeleted() {
		return
	}
	// upper report reply of  self in upper's arc ,treat it as del
	if s.isUpper(c, mid, oid, tp) && r.Mid == mid {
		s.dao.Databus.Delete(c, mid, oid, rpID, now.Unix(), tp, false)
		report.User(&report.UserInfo{
			Mid:      mid,
			Platform: platform,
			Build:    build,
			Buvid:    buvid,
			Business: 41,
			Type:     int(r.Type),
			Oid:      r.Oid,
			Action:   model.ReportReplyDel,
			Ctime:    time.Now(),
			IP:       ip,
			Index: []interface{}{
				r.RpID,
				r.State,
				model.ReplyStateUpDel,
			},
		})
		return
	}
	// 信用评分获取，用于优先举报排序处理
	var score int
	arg := &figmdl.ArgUserFigure{Mid: mid}
	fig, err := s.figure.UserFigure(c, arg)
	if err != nil {
		log.Error("s.figure.UserFigure(mid:%d) error(%v)", mid, err)
		err = nil
	} else {
		score = 100 - int(fig.Percentage)
	}
	ctime := xtime.Time(now.Unix())
	rpt := &model.Report{
		RpID:    rpID,
		Oid:     oid,
		Type:    tp,
		Mid:     mid,
		Reason:  reason,
		Count:   1,
		Content: cont,
		Score:   score,
		State:   model.GetReportType(reason),
		CTime:   ctime,
		MTime:   ctime,
	}
	rptUser := &model.ReportUser{
		Oid:     oid,
		Type:    tp,
		RpID:    rpID,
		Mid:     mid,
		Reason:  reason,
		Content: cont,
		State:   model.ReportUserStateNew,
		CTime:   ctime,
		MTime:   ctime,
	}
	if rpt.ID, err = s.dao.Report.InsertUser(c, rptUser); err != nil {
		return
	}
	if rpt.ID == 0 {
		err = ecode.ReplyReported
		return
	}
	if tp != model.SubTypeActArc && tp != model.SubTypePlaylist && tp != model.SubTypeComicSeason && tp != model.SubTypeComicEpisode {
		var (
			m           = make(map[int64]string)
			title, link string
			typeid      int32
		)
		m[rpID] = ""
		message := r.Content.Message
		s.dao.FilterContents(c, m)
		if m[rpID] != "" {
			message = m[rpID]
		}
		title, link, typeid, _ = s.TitleLink(c, oid, tp)
		if link != "" {
			link = fmt.Sprintf("%s#reply%d", link, rpID)
		}
		err = s.workflow.AddReport(c, oid, tp, typeid, rpID, score, reason, mid, r.Mid, r.Like, message, link, title)
		if err != nil {
			return
		}
	}
	if rpt.ID, err = s.dao.Report.Insert(c, rpt); err != nil {
		return
	}
	// set report count and set user report dao.Redis key TTL for one day seconds.
	if err = s.dao.Redis.SetUserReportCnt(c, mid, cnt+1, now); err != nil {
		log.Error("s.redis.SetUserReportCnt(%v) error(%v) or row==0", rpt, err)
		return
	}
	s.dao.Databus.AddReport(c, oid, rpID, tp)
	report.User(&report.UserInfo{
		Mid:      rptUser.Mid,
		Platform: platform,
		Build:    build,
		Buvid:    buvid,
		Business: 41,
		Type:     int(tp),
		Oid:      oid,
		Action:   model.ReportReplyReport,
		Ctime:    now,
		IP:       ip,
		Index: []interface{}{
			r.Mid,
			r.RpID,
			r.State,
			fmt.Sprint(reason),
		},
		Content: map[string]interface{}{
			"count":   cnt,
			"score":   score,
			"content": cont,
		},
	})
	log.Info("AddReport(%d) oid:%d mid:%d score:%d percentage:%v", rpt.ID, oid, mid, score, fig)
	return
}

// ReportRelated get related replies of report.
func (s *Service) ReportRelated(c context.Context, mid, oid, rpid int64, tp int8, escape bool) (sub *model.Subject, root *model.Reply, related []*model.Reply, err error) {
	if !model.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	// root reply
	sub, root, snds, err := s.ReportReply(c, mid, oid, rpid, tp, 1, s.sndDefCnt, escape)
	if err != nil {
		return
	}
	root.Replies = snds
	// related reply
	start := root.Floor - 6
	end := root.Floor + 6
	relIDs, err := s.dao.Reply.GetIDsByFloorOffset(c, oid, tp, start, end)
	if err != nil {
		return
	}
	relMap, err := s.repliesMap(c, oid, tp, relIDs)
	if err != nil {
		return
	}
	related = make([]*model.Reply, 0, len(relMap))
	bs := make([]*model.Reply, 0, len(relMap))
	for _, rpID := range relIDs {
		rp, ok := relMap[rpID]
		if ok && rp.RpID != root.RpID {
			if rp.Replies, err = s.reportReplies(c, rp, 1, s.sndDefCnt); err != nil {
				return
			}
			related = append(related, rp)
			// to build replies
			bs = append(bs, rp)
			bs = append(bs, rp.Replies...)
		}
	}
	if err = s.buildReply(c, sub, bs, mid, escape); err != nil {
		return
	}
	return
}

// ReportReply get report reply.
func (s *Service) ReportReply(c context.Context, mid, oid, rpid int64, tp int8, pn, ps int, escape bool) (sub *model.Subject, root *model.Reply, seconds []*model.Reply, err error) {
	if !model.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if sub, err = s.subject(c, oid, tp); err != nil {
		return
	}
	if root, err = s.ReplyContent(c, oid, rpid, tp); err != nil {
		return
	}
	if root.Root != 0 {
		if root, err = s.ReplyContent(c, root.Oid, root.Root, root.Type); err != nil {
			return
		}
	}
	if seconds, err = s.reportReplies(c, root, pn, ps); err != nil {
		return
	}
	bs := make([]*model.Reply, 0, len(seconds)+1)
	bs = append(bs, root)
	bs = append(bs, seconds...)
	if err = s.buildReply(c, sub, bs, mid, escape); err != nil {
		return
	}
	return
}

func (s *Service) reportReplies(c context.Context, rp *model.Reply, pn, ps int) (rs []*model.Reply, err error) {
	var (
		start = (pn - 1) * ps
	)
	rs = _emptyReplies
	if start >= rp.Count {
		return
	}
	sndIDs, err := s.dao.Reply.GetIDsByRootWithoutState(c, rp.Oid, rp.RpID, rp.Type, start, ps)
	if err != nil {
		return
	}
	sndMap, err := s.repliesMap(c, rp.Oid, rp.Type, sndIDs)
	if err != nil {
		return
	}
	rs = make([]*model.Reply, 0, len(sndMap))
	for _, rpID := range sndIDs {
		rp, ok := sndMap[rpID]
		if ok {
			rs = append(rs, rp)
		}
	}
	return
}

func (s *Service) linkByOids(c context.Context, oids map[int64]string, typ int8) (err error) {
	if len(oids) == 0 {
		return
	}
	if typ == model.SubTypeActivity {
		err = s.workflow.TopicsLink(c, oids, false)
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
				_, link, err = s.workflow.ActivitySub(c, oid)
				if err != nil {
					return
				}
			case model.SubTypeArticle:
				link = fmt.Sprintf("https://www.bilibili.com/read/cv%d", oid)
			case model.SubTypeMusic:
				link = fmt.Sprintf("https://www.bilibili.com/audio/au%d", oid)
			case model.SubTypeMusicList:
				link = fmt.Sprintf("https://www.bilibili.com/audio/am%d", oid)
			case model.SubTypeLive:
				link = fmt.Sprintf("https://vc.bilibili.com/video/%d", oid)
			case model.SubTypeLiveAct:
				_, link, err = s.workflow.LiveActivityTitle(c, oid)
				if err != nil {
					return
				}
			case model.SubTypeLivePicture:
				link = fmt.Sprintf("https://h.bilibili.com/ywh/%d", oid)
			case model.SubTypeCredit:
				link = fmt.Sprintf("https://www.bilibili.com/judgement/case/%d", oid)
			case model.SubTypeDynamic:
				link = fmt.Sprintf("https://t.bilibili.com/%d", oid)
			case model.SubTypeLiveNotice:
				link = fmt.Sprintf("http://link.bilibili.com/p/eden/news#/newsdetail?id=%d", oid)
			default:
				return
			}
			oids[oid] = link
		}
	}
	return
}

// TitleLink TitleLink
func (s *Service) TitleLink(c context.Context, oid int64, typ int8) (title, link string, typeId int32, err error) {
	switch typ {
	case model.SubTypeArchive:
		arg := &arcmdl.ArgAid2{
			Aid: oid,
		}
		var m *api.Arc
		m, err = s.arcSrv.Archive3(c, arg)
		if err != nil || m == nil {
			log.Error("s.arcSrv.Archive3(%v) ret:%v error(%v)", arg, m, err)
			return
		}
		link = fmt.Sprintf("http://www.bilibili.com/video/av%d/", oid)
		title = m.Title
		typeId = m.TypeID
	case model.SubTypeTopic:
		if title, link, err = s.workflow.TopicTitle(c, oid); err != nil {
			log.Error("s.noticeDao.Topic(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeMusic:
		link = fmt.Sprintf("https://www.bilibili.com/audio/au%d", oid)
	case model.SubTypeMusicList:
		link = fmt.Sprintf("https://www.bilibili.com/audio/am%d", oid)
	case model.SubTypeActivity:
		if title, link, err = s.workflow.TopicTitle(c, oid); err != nil {
			log.Error("s.noticeDao.Activity(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeForbiden:
		title, link, err = s.workflow.BanTitle(c, oid)
		if err != nil {
			return
		}
	case model.SubTypeNotice:
		title, link, err = s.workflow.NoticeTitle(c, oid)
		if err != nil {
			return
		}
	case model.SubTypeActArc:
		if title, link, err = s.workflow.TopicTitle(c, oid); err != nil {
			log.Error("s.noticeDao.ActivitySub(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeArticle:
		arg := &artmdl.ArgAids{
			Aids: []int64{oid},
		}
		var m map[int64]*artmdl.Meta
		m, err = s.articleSrv.ArticleMetas(c, arg)
		if err != nil || m == nil {
			log.Error("s.articleSrv.ArticleMetas(%v) ret:%v error(%v)", arg, m, err)
			return
		}
		if meta, ok := m[oid]; ok {
			title = meta.Title
			link = fmt.Sprintf("http://www.bilibili.com/read/cv%d", oid)
		}
	case model.SubTypeLive:
		if title, link, err = s.workflow.LiveVideoTitle(c, oid); err != nil {
			log.Error("s.noticeDao.LiveSmallVideo(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeLiveAct:
		if title, link, err = s.workflow.LiveActivityTitle(c, oid); err != nil {
			log.Error("s.noticeDao.LiveActivity(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeLiveNotice:
		if title, err = s.workflow.LiveNotice(c, oid); err != nil {
			log.Error("s.noticeDao.LiveNotice(%d) error(%v)", oid, err)
			return
		}
		link = fmt.Sprintf("http://link.bilibili.com/p/eden/news#/newsdetail?id=%d", oid)
		return
	case model.SubTypeLivePicture:
		if title, link, err = s.workflow.LivePictureTitle(c, oid); err != nil {
			log.Error("s.noticeDao.LivePiture(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeCredit:
		if title, link, err = s.workflow.CreditTitle(c, oid); err != nil {
			log.Error("s.noticeDao.Credit(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeDynamic:
		if title, link, err = s.workflow.DynamicTitle(c, oid); err != nil {
			log.Error("s.noticeDao.Dynamic(%d) error(%v)", oid, err)
			return
		}
	case model.SubTypeHuoniao:
		if title, link, err = s.workflow.HuoniaoTitle(c, oid); err != nil {
			log.Error("s.workflow.HuoniaoTitle(%d) error(%v)", oid, err)
			return
		}
	default:
		return
	}
	return
}
