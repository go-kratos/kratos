package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/workflow/model/param"
	"go-common/app/job/main/workflow/model"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	// LOGAPPID audit log appid .
	LOGAPPID = "log_audit"
	// LOGINDEX .
	LOGINDEX = "log_audit_11_all"
	// GROUPAPPID .
	GROUPAPPID = "workflow_group_common"
	// GROUPINDEX .
	GROUPINDEX = "workflow_group_common"
	// CHALLAPPID .
	CHALLAPPID = "workflow_chall_common"
	// CHALLINDEX .
	CHALLINDEX = "workflow_chall_common"
	// BUSINESS .
	BUSINESS = 1
	// MC .
	MC = "1_15_1"
)

// notifyState .
var notifyState = []int64{0, 2}

// notifyproc .
func (s *Service) notifyproc(c context.Context) {
	for {
		if err := s.sendMessage(c); err != nil {
			log.Error("s.sendMessage type(%d) error(%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		time.Sleep(time.Second * 60)
	}
}

// SendMessage .
func (s *Service) sendMessage(c context.Context) (err error) {
	var (
		group []*model.GroupRes
		chall map[int64][]*model.ChallRes
		gids  []int64
	)
	if group, err = s.NotNotifyGroupSearch(c); err != nil {
		log.Error("s.NotNotifyGroupSearch error(%v)", err)
		return
	}
	if len(group) <= 0 {
		log.Info("group search length is 0")
		return
	}
	for _, g := range group {
		gids = append(gids, g.ID)
	}
	if chall, err = s.NotNotifyChallSearch(c, gids); err != nil {
		log.Error("s.NotNotifyChallSearch error(%v)", err)
		return
	}
	if len(chall) <= 0 {
		log.Info("chall search length is 0")
		return
	}
	for _, g := range group {
		mids := []int64{}
		ok := false
		var chs []*model.ChallRes
		if chs, ok = chall[g.ID]; !ok {
			log.Error("gid(%d) in group but not in chall", g.ID)
			continue
		}
		for _, ch := range chs {
			mids = append(mids, ch.MID)
		}
		if len(mids) <= 0 {
			log.Error("gid(%d) with mid in chall is empty", g.ID)
			continue
		}
		param := &param.MessageParam{
			Type:     "json",
			Source:   1,
			DataType: 4,
			MC:       MC,
			Title:    "您的投诉已收到",
			Context:  fmt.Sprintf("您对稿件（av%d）的举报我们已经收到。感谢您对 bilibili 社区秩序的维护，哔哩哔哩 (゜-゜)つロ 干杯~ ", g.OID),
			MidList:  mids,
		}
		log.Info("send message param(%+v)", param)
		if err = s.dao.SendMessage(c, chs, param); err != nil {
			log.Error("s.dao.SendMessage error(%v)", err)
		}
	}
	return
}

// NotNotifyGroupSearch .
func (s *Service) NotNotifyGroupSearch(c context.Context) (result []*model.GroupRes, err error) {
	var (
		pn      = 1
		ps      = 1000
		tempRes *model.SearchGroup
	)
	frontTenMin, _ := time.ParseDuration("-10m")
	frontTwelveMin, _ := time.ParseDuration("-12m")
	frontTenMinFormat := time.Now().Add(frontTenMin).Format("2006-01-02 15:04:05")
	frontTwelveMinFormat := time.Now().Add(frontTwelveMin).Format("2006-01-02 15:04:05")
	e := elastic.NewElastic(nil)
	r := e.NewRequest(GROUPAPPID).Index(GROUPINDEX).Fields("id", "oid").
		WhereEq("business", BUSINESS).
		WhereIn("state", notifyState).
		WhereRange("lasttime", frontTwelveMinFormat, frontTenMinFormat, elastic.RangeScopeLoRo).
		Pn(pn).
		Ps(ps)
	if err = r.Scan(context.Background(), &tempRes); err != nil {
		log.Error("elastic search group Scan error(%v)", err)
		return
	}
	log.Info("groupparams is(%v)", r.Params())
	res := tempRes.Result
	if len(res) <= 0 {
		return
	}
	// search audit log
	var logGids map[int64][]*model.LogRes
	if logGids, err = s.searchAuditLog(c, res); err != nil {
		log.Error("s.searchAuditLog error(%v)", err)
		return
	}
	for _, r := range res {
		if _, ok := logGids[r.OID]; !ok {
			result = append(result, r)
		}
	}
	return
}

// searchAuditLog .
func (s *Service) searchAuditLog(c context.Context, grp []*model.GroupRes) (logGids map[int64][]*model.LogRes, err error) {
	pn := 1
	ps := 1000
	oids := []int64{}
	logGids = make(map[int64][]*model.LogRes)
	for _, g := range grp {
		oids = append(oids, g.OID)
	}
	e := elastic.NewElastic(nil)
	for {
		var res *model.AuditLog
		r := e.NewRequest(LOGAPPID).Index(LOGINDEX).Fields("int_1", "oid").
			WhereIn("oid", oids).
			WhereEq("type", 2).
			WhereEq("business", 11).
			WhereEq("action", "notify_users_received").
			Pn(pn).
			Ps(ps)
		if err = r.Scan(context.Background(), res); err != nil {
			log.Error("elastic search audit log Scan error(%v)", err)
			return
		}
		if res == nil {
			time.Sleep(time.Second * 3)
			continue
		}
		if res.Page.Total > 10000 {
			log.Error("elastic search audit log result too long")
			err = ecode.ServerErr
			return
		}
		log.Info("auditlogparams is(%v)", r.Params())
		logRes := res.Result
		for _, lg := range logRes {
			logGids[lg.Oid] = append(logGids[lg.Oid], lg)
		}
		if len(res.Result) < ps {
			break
		}
		pn++
	}
	return
}

// NotNotifyChallSearch .
func (s *Service) NotNotifyChallSearch(c context.Context, gids []int64) (chall map[int64][]*model.ChallRes, err error) {
	var (
		pn = 1
		ps = 1000
	)
	if len(gids) <= 0 {
		return
	}
	chall = make(map[int64][]*model.ChallRes)
	frontTenMin, _ := time.ParseDuration("-10m")
	frontTwelveMin, _ := time.ParseDuration("-12m")
	frontTenMinFormat := time.Now().Add(frontTenMin).Format("2006-01-02 15:04:05")
	frontTwelveMinFormat := time.Now().Add(frontTwelveMin).Format("2006-01-02 15:04:05")
	e := elastic.NewElastic(nil)
	for {
		res := &model.SearchChall{}
		r := e.NewRequest(CHALLAPPID).Index(CHALLINDEX).Fields([]string{"id", "gid", "mid", "oid"}...).
			WhereEq("business", BUSINESS).
			WhereEq("state", 0).
			WhereIn("gid", gids).
			WhereRange("ctime", frontTwelveMinFormat, frontTenMinFormat, elastic.RangeScopeLoRo).
			Pn(pn).
			Ps(ps)
		if err = r.Scan(context.Background(), &res); err != nil {
			log.Error("elastic search chall Scan error(%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		log.Info("challparams is(%v)", r.Params())
		for _, r := range res.Result {
			chall[r.GID] = append(chall[r.GID], r)
		}
		if len(res.Result) < ps {
			break
		}
		pn++
	}
	return
}
