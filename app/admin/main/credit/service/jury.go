package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	creditMDL "go-common/app/admin/main/credit/model"
	"go-common/app/admin/main/credit/model/blocked"
	account "go-common/app/service/main/account/api"
	relationmdl "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

// AddJury add jury to admin
func (s *Service) AddJury(c context.Context, arg *blocked.ArgAddJurys) (err error) {
	arg.MIDs = creditMDL.ArrayUnique(arg.MIDs)
	if len(arg.MIDs) == 0 {
		return
	}
	items := []*blocked.Jury{}
	if err = s.blockedDao.ReadDB.Where("mid IN (?)", arg.MIDs).Find(&items).Error; err != nil {
		log.Error("s.blockedDao.juryByMIDs(%s) error(%v)", xstr.JoinInts(arg.MIDs), err)
		err = nil
	}
	mids := make(map[int64]*blocked.Jury, len(arg.MIDs))
	for _, v := range items {
		mids[v.UID] = v
	}
	day, _ := time.ParseDuration("24h")
	users := make([]*blocked.Jury, 0, len(arg.MIDs))
	tx := s.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	now := time.Now()
	var multiple []interface{}
	for _, mid := range arg.MIDs {
		user := &blocked.Jury{
			UID:     mid,
			Status:  blocked.JuryStatusOn,
			OPID:    arg.OPID,
			Remark:  blocked.AddJuryRemark,
			Expired: xtime.Time(now.Add(day * time.Duration(arg.Day)).Unix()),
		}
		if jury, ok := mids[mid]; ok {
			user.ID = jury.ID
			switch {
			case jury.Black == blocked.JuryBlack:
				continue
			case jury.Black == blocked.JuryWhite:
				user.Black = blocked.JuryWhite
			}
			user.CTime = jury.CTime
			user.VoteTotal = jury.VoteTotal
			user.VoteRight = jury.VoteRight
			user.Total = jury.Total
			single := map[string]interface{}{
				"id":      jury.ID,
				"oper_id": arg.OPID,
				"status":  blocked.JuryStatusOn,
				"black":   user.Black,
				"expired": now.Add(day * time.Duration(arg.Day)).Format(creditMDL.TimeFormatSec),
			}
			multiple = append(multiple, single)
		}
		if err = tx.Save(&user).Error; err != nil {
			log.Error("s.DB.Create error(%v)", err)
			return
		}
		s.AddNotify(func() {
			s.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedJury, blocked.TableBlockedJury, multiple)
		})
		users = append(users, user)
	}
	if arg.Send == blocked.NeedSendMsg {
		for _, u := range users {
			msg := &blocked.SysMsg{
				Type:     blocked.MsgTypeGetJuryer,
				MID:      u.UID,
				Day:      arg.Day,
				RemoteIP: metadata.String(c, metadata.RemoteIP),
			}
			s.MsgCh <- msg
		}
	}
	return
}

// Cases  deal case data.
func (s *Service) Cases(c context.Context, arg *blocked.ArgCaseSearch) (list []*blocked.Case, pager *blocked.Pager, err error) {
	var ids []int64
	ids, pager, err = s.searchDao.Case(c, arg)
	if err != nil {
		log.Error(" s.searchDao.Case  error(%v)", err)
		return
	}
	if len(ids) == 0 {
		return
	}
	var (
		accInfoMap   map[int64]*account.Info
		relationsMap map[int64]*relationmdl.Stat
		uids         = make([]int64, 0)
	)
	ids = creditMDL.ArrayUnique(ids)
	if err = s.blockedDao.ReadDB.Where("id IN (?)", ids).Order(fmt.Sprintf("%s %s", arg.Order, arg.Sort)).Find(&list).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("s.blockedDao(%s) error(%v)", xstr.JoinInts(ids), err)
			return
		}
		log.Warn("search ids(%s) not in db", xstr.JoinInts(ids))
		err = nil
	}
	for _, v := range list {
		if v.MID > 0 {
			uids = append(uids, v.MID)
		}
	}
	if accInfoMap, err = s.accDao.RPCInfos(c, uids); err != nil {
		log.Error("s.accDao.RPCInfos(%s) error(%v)", xstr.JoinInts(uids), err)
		err = nil
	}
	if relationsMap, err = s.RelationDao.RPCStats(c, uids); err != nil {
		log.Error("s.RelationDao.RPCStats(%s) error(%v)", xstr.JoinInts(uids), err)
		err = nil
	}
	for _, v := range list {
		if accInfo, ok := accInfoMap[v.MID]; ok {
			v.Uname = accInfo.Name
		}
		if relationInfo, ok := relationsMap[v.MID]; ok {
			v.Fans = relationInfo.Follower
		}
		v.OPName = s.Managers[v.OPID]
		v.StatusDesc = blocked.StatusDesc[v.Status]
		v.OriginTypeDesc = blocked.OriginTypeDesc[v.OriginType]
		v.ReasonTypeDesc = blocked.ReasonTypeDesc(v.ReasonType)
		v.PunishDesc = blocked.PunishDesc[v.PunishResult]
		v.CaseTypeDesc = blocked.CaseTypeDesc[v.CaseType]
		v.VoteTotal = blocked.VoteTotal(v.VoteRule, v.VoteBreak, v.VoteDelete)
		if v.VoteRule > 0 || v.VoteBreak > 0 || v.VoteDelete > 0 {
			v.BlockedPercent = blocked.BreakPercent(v.VoteRule, v.VoteBreak, v.VoteDelete)
			v.DeletePercent = blocked.DeletePercent(v.VoteRule, v.VoteBreak, v.VoteDelete)
			v.RulePercent = blocked.RulePercent(v.VoteRule, v.VoteBreak, v.VoteDelete)
		}
	}
	return
}

// Opinions is block opinions
func (s *Service) Opinions(c context.Context, arg *blocked.ArgOpinionSearch) (list []*blocked.Opinion, pager *blocked.Pager, err error) {
	var ids []int64
	ids, pager, err = s.searchDao.Opinion(c, arg)
	if err != nil {
		log.Error("s.blockedDao.Search appid(%s) error(%v)", "", err)
		return
	}
	if len(ids) == 0 {
		return
	}
	var (
		accInfoMap   map[int64]*account.Info
		relationsMap map[int64]*relationmdl.Stat
		uids         = make([]int64, 0)
	)
	ids = creditMDL.ArrayUnique(ids)
	if err = s.blockedDao.ReadDB.Where("id IN (?)", ids).Order(fmt.Sprintf("%s %s", arg.Order, arg.Sort)).Find(&list).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("s.blockedDao(%s) error(%v)", xstr.JoinInts(ids), err)
			return
		}
		log.Warn("search ids(%s) not in db", xstr.JoinInts(ids))
		err = nil
	}
	for _, v := range list {
		if v.MID > 0 {
			uids = append(uids, v.MID)

		}
	}
	accInfoMap, err = s.accDao.RPCInfos(c, uids)
	if err != nil {
		log.Error("s.accDao.RPCInfos error(%v)", err)
		err = nil
	}
	if relationsMap, err = s.RelationDao.RPCStats(c, uids); err != nil {
		log.Error("s.RelationDao.RPCStats(%s) error(%v)", xstr.JoinInts(uids), err)
		err = nil
	}
	for _, v := range list {
		if opinion, ok := accInfoMap[v.MID]; ok {
			v.UName = opinion.Name
		}
		if relationInfo, ok := relationsMap[v.MID]; ok {
			v.Fans = relationInfo.Follower
		}
		v.OPName = s.Managers[v.OperID]
		v.AttrDesc = blocked.AttrDesc[v.Attr]
		v.VoteDesc = blocked.VoteDesc[v.Vote]
		v.VoteStateDesc = blocked.VoteStateDesc[v.State]
	}
	return
}

// Jurys select jury lists.
func (s *Service) Jurys(c context.Context, arg *blocked.ArgJurySearch) (list []*blocked.Jury, pager *blocked.Pager, err error) {
	var ids []int64
	ids, pager, err = s.searchDao.Jury(c, arg)
	if err != nil {
		log.Error("s.searchDao.Jury error(%v)", err)
		return
	}
	if len(ids) == 0 {
		return
	}
	var (
		accInfoMap map[int64]*account.Info
		uids       = make([]int64, len(ids))
	)
	ids = creditMDL.ArrayUnique(ids)
	if err = s.blockedDao.ReadDB.Where("id IN (?)", ids).Order(fmt.Sprintf("%s %s", arg.Order, arg.Sort)).Find(&list).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("s.blockedDao(%s) error(%v)", xstr.JoinInts(ids), err)
			return
		}
		log.Warn("search ids(%s) not in db", xstr.JoinInts(ids))
		err = nil
	}
	for _, v := range list {
		uids = append(uids, v.UID)
	}
	accInfoMap, err = s.accDao.RPCInfos(c, uids)
	if err != nil {
		log.Error("s.accDao.RPCInfos error(%v)", err)
		err = nil
	}
	for _, v := range list {
		if jury, ok := accInfoMap[v.UID]; ok {
			v.UName = jury.Name
		}
		v.OPName = s.Managers[v.OPID]
		v.EffectDay = xtime.Time(v.Expired.Time().AddDate(0, 0, -blocked.JuryDay).Unix())
		if v.EffectDay < v.CTime {
			v.EffectDay = v.CTime
		}
		if v.VoteTotal > 0 {
			v.VoteRadio = strconv.FormatFloat(float64(v.VoteRight)/float64(v.VoteTotal)*100, 'f', 2, 64)
		} else {
			v.VoteRadio = "-1"
		}
		v.StatusDesc = blocked.JuryerStatus[v.Status]
		v.BlackDesc = blocked.JuryerStyle[v.Black]
	}
	return
}

// JurysEx export jury list.
func (s *Service) JurysEx(c context.Context, arg *blocked.ArgJurySearch) (list []*blocked.Jury, err error) {
	var (
		ps    = 500
		count int
		pager *blocked.Pager
		g     errgroup.Group
	)
	if list, pager, err = s.Jurys(c, arg); err != nil {
		log.Error("s.Infos(%+v) error(%v)", arg, err)
		return
	}
	if pager == nil {
		log.Warn("arg(%+v) jury search data empty!", arg)
		return
	}
	count = pager.Total / ps
	if pager.Total%ps != 0 {
		count++
	}
	lCh := make(chan []*blocked.Jury, count)
	for pn := 1; pn <= count; pn++ {
		tmpPn := pn
		g.Go(func() (err error) {
			var gJury []*blocked.Jury
			gArg := &blocked.ArgJurySearch{
				UID:         arg.UID,
				Status:      arg.Status,
				Black:       arg.Black,
				ExpiredFrom: arg.ExpiredFrom,
				ExpiredTo:   arg.ExpiredTo,
				Order:       arg.Order,
				Sort:        arg.Sort,
				PN:          tmpPn,
				PS:          ps,
			}
			gJury, _, err = s.Jurys(c, gArg)
			if err != nil {
				log.Error("s.Jurys(%+v) error(%v)", gJury, err)
				err = nil
				return
			}
			lCh <- gJury
			return
		})
	}
	g.Wait()
	close(lCh)
	for jInfo := range lCh {
		list = append(list, jInfo...)
	}
	sort.Slice(list, func(i int, j int) bool {
		return list[i].ID < list[j].ID
	})
	return
}

// UpCaseStatus update case status.
func (s *Service) UpCaseStatus(c context.Context, arg *blocked.ArgUpStatus) (err error) {
	items := []blocked.Case{}
	if err = s.blockedDao.ReadDB.Where("id IN(?)", arg.IDS).Find(&items).Error; err != nil {
		log.Error("s.blockedDao.caseByIDs(%s) error(%v)", xstr.JoinInts(arg.IDS), err)
		return
	}
	var multiple []interface{}
	if arg.Status == blocked.CaseStatusGranting {
		hour := s.CaseConfig(blocked.ConfigCaseGiveHours)
		now := time.Now()
		if err = s.DB.Model(blocked.Case{}).Where("id IN(?)", arg.IDS).Updates(
			&blocked.Case{
				Status:    arg.Status,
				OPID:      arg.OID,
				StartTime: xtime.Time(now.Unix()),
				EndTime:   xtime.Time(now.Add(time.Duration(hour) * time.Hour).Unix()),
			}).Error; err != nil {
			return
		}
		for _, item := range items {
			single := map[string]interface{}{
				"id":         item.ID,
				"oper_id":    arg.OID,
				"start_time": now.Format(creditMDL.TimeFormatSec),
				"status":     arg.Status,
				"case_type":  item.CaseType,
			}
			multiple = append(multiple, single)
		}
	} else {
		if err = s.DB.Model(blocked.Case{}).Where("id IN(?)", arg.IDS).Updates(
			&blocked.Case{
				Status: arg.Status,
				OPID:   arg.OID,
			}).Error; err != nil {
			return
		}
		for _, item := range items {
			single := map[string]interface{}{
				"id":        item.ID,
				"oper_id":   arg.OID,
				"status":    arg.Status,
				"case_type": item.CaseType,
			}
			multiple = append(multiple, single)
		}
	}
	s.AddNotify(func() {
		s.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedCase, blocked.TableBlockedCase, multiple)
	})
	return
}

// CaseReason reason info.
func (s *Service) CaseReason(c context.Context) (cr map[int]string, err error) {
	var reason []*blocked.Reason
	cr = make(map[int]string)
	if err = s.ReadDB.Select([]string{"id", "reason"}).Where("status = 0").Find(&reason).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
			return
		}
		return
	}
	for _, v := range reason {
		cr[v.ID] = v.Reason
	}
	return
}

// CaseConf caseconfig list.
func (s *Service) CaseConf(c context.Context) (cc map[string]string, err error) {
	var configs []*blocked.Config
	cc = make(map[string]string)
	err = s.ReadDB.Select([]string{"config_key", "content"}).Find(&configs).Error
	if err != nil || configs == nil {
		return
	}
	for _, v := range configs {
		cc[v.ConfigKey] = v.Content
	}
	return
}

// SetCaseConf update case config.
func (s *Service) SetCaseConf(c context.Context, cc *blocked.ArgCaseConf) (err error) {
	v := reflect.ValueOf(cc).Elem()
	tx := s.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	for i := 0; i < v.NumField(); i++ {
		fieldValue := strconv.Itoa(v.Field(i).Interface().(int))
		tag := v.Type().Field(i).Tag.Get("form")
		if err = tx.Model(&blocked.Config{}).Where("config_key = ?", tag).Updates(map[string]interface{}{"content": fieldValue, "oper_id": cc.OID}).Error; err != nil {
			log.Error("caseConfig(%s,%d) update error(%v)", v.Field(i).Type().Name(), v.Field(i).Interface().(int64), err)
			return
		}
	}
	return
}

// CaseConfig get config by key.
func (s *Service) CaseConfig(ck string) (cv int64) {
	if v, ok := s.caseConfCache[ck]; ok {
		cv, _ = strconv.ParseInt(v, 10, 64)
	}
	return
}

// KPIPoint get kpi_point data
func (s *Service) KPIPoint(c context.Context, arg *blocked.ArgKpiPointSearch) (list []*blocked.KPIPoint, pager *blocked.Pager, err error) {
	var ids []int64
	ids, pager, err = s.searchDao.KPIPoint(c, arg)
	if err != nil {
		log.Error("s.searchDao.KPIPoint error (%v)", err)
		return
	}
	if len(ids) == 0 {
		return
	}
	ids = creditMDL.ArrayUnique(ids)
	if err = s.blockedDao.ReadDB.Where("id IN (?)", ids).Order(fmt.Sprintf("%s %s", arg.Order, arg.Sort)).Find(&list).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("s.blockedDao(%s) error(%v)", xstr.JoinInts(ids), err)
			return
		}
		log.Warn("search ids(%s) not in db", xstr.JoinInts(ids))
		err = nil
	}
	return
}

// VotenumConf votenumConf list.
func (s *Service) VotenumConf(c context.Context) (cc string, err error) {
	var bc = &blocked.Config{}
	err = s.ReadDB.Select("content").Where("config_key = 'vote_num'").Find(&bc).Error
	if err != nil || bc.Content == "" {
		log.Error("s.VotenumConf(%s) error(%v)", cc, err)
		return
	}
	cc = bc.Content
	return
}

// SetVotenumConf .
func (s *Service) SetVotenumConf(c context.Context, vn *blocked.ArgVoteNum) (err error) {
	v := blocked.VoteNum{}
	v.RateS = vn.RateS
	v.RateA = vn.RateA
	v.RateB = vn.RateB
	v.RateC = vn.RateC
	v.RateD = vn.RateD
	vnb, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal(%+v) error(%v)", v, err)
	}
	config := &blocked.Config{}
	s.ReadDB.Where("config_key = 'vote_num'").Find(config)
	var vconfig *blocked.Config
	if config != nil && config.ID != 0 {
		vconfig = &blocked.Config{
			ID:        config.ID,
			ConfigKey: "vote_num",
			Content:   string(vnb),
			OperID:    vn.OID,
		}
	} else {
		vconfig = &blocked.Config{
			ConfigKey: "vote_num",
			Content:   string(vnb),
			OperID:    vn.OID,
		}
	}
	if err = s.DB.Save(&vconfig).Error; err != nil {
		log.Error("http.setAutoCaseConfig error(%v)", err)
		return
	}
	return
}

// WebHook a webhook from work flow.
func (s *Service) WebHook(c context.Context, data []byte) (err error) {
	wh := &blocked.WebHook{}
	if err = json.Unmarshal(data, &wh); err != nil {
		err = ecode.RequestErr
		log.Error("webhook json rawmessage(%s) error(%v)", string(data), err)
		return
	}
	if wh.Verb == "chall.SetResult" || wh.Verb == "chall.BatchSetResult" {
		if wh.Target == nil || wh.Object == nil {
			log.Warn("wh.Target or wh.Object is nil %v,%v", wh.Target, wh.Object)
			return
		}
		// appeal state not changed .
		if wh.Target.State == wh.Object.State {
			log.Warn("appeal state not changed target=%d object=%d", wh.Target.State, wh.Object.State)
			return
		}
		switch wh.Object.State {
		case blocked.AppealStateSucc:
			if err = s.blockedDao.DB.Model(&blocked.Info{}).Where("case_id =?", wh.Target.OID).Update("status", blocked.BlockStateClose).Error; err != nil {
				log.Error("s.blockedDao.DB error(%v)", err)
				return
			}
			s.MsgCh <- &blocked.SysMsg{Type: blocked.MsgTypeAppealSucc, MID: wh.Target.Mid, CID: wh.Target.OID, RemoteIP: metadata.String(c, metadata.RemoteIP)}
		case blocked.AppealStateFail:
			s.MsgCh <- &blocked.SysMsg{Type: blocked.MsgTypeAppealFail, MID: wh.Target.Mid, CID: wh.Target.OID, RemoteIP: metadata.String(c, metadata.RemoteIP)}
		default:
			log.Warn("unknown webhook state(%d) ", wh.Object.State)
		}
	}
	return
}
