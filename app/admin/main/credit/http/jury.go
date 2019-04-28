package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"time"

	"go-common/app/admin/main/credit/model"
	"go-common/app/admin/main/credit/model/blocked"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
)

func users(c *bm.Context) {
	v := new(blocked.ArgJurySearch)
	if err := c.Bind(v); err != nil {
		return
	}
	list, pager, err := creSvc.Jurys(c, v)
	if err != nil {
		log.Error("creSvc.Jurys error(%v)", err)
		httpCode(c, err)
		return
	}
	if len(list) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, list, pager)
}

func userAdd(c *bm.Context) {
	var err error
	v := new(blocked.ArgAddJurys)
	if err = c.Bind(v); err != nil {
		return
	}
	if err = creSvc.AddJury(c, v); err != nil {
		httpCode(c, err)
		return
	}
	log.Info("business_userAdd:%+v", v)
	httpCode(c, nil)
}

// userByID  get user by id.
func userByID(c *bm.Context) {
	var err error
	v := new(struct {
		MID int64 `form:"mid" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.Jury{}
	if err = creSvc.ReadDB.Where("mid = ?", v.MID).Find(item).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("creSvc.userByID error(%v)", err)
			httpCode(c, err)
			return
		}
		httpData(c, nil, nil)
		return
	}
	if item != nil {
		if item.VoteTotal > 0 {
			item.VoteRadio = strconv.FormatFloat(float64(item.VoteRight)/float64(item.VoteTotal)*100, 'f', 2, 64)
		} else {
			item.VoteRadio = "-1"
		}
		item.StatusDesc = blocked.JuryerStatus[item.Status]
		item.BlackDesc = blocked.JuryerStyle[item.Black]
	}
	httpData(c, item, nil)
}

func upUserStatus(c *bm.Context) {
	v := new(struct {
		MIDS   []int64 `form:"mids,split" validate:"min=1,max=100"`
		Status int8    `form:"status"  validate:"min=1,max=2" default:"2"`
		OID    int64   `form:"op_id"  validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	items := []*blocked.Jury{}
	if err := creSvc.ReadDB.Where("mid IN (?)", v.MIDS).Find(&items).Error; err != nil {
		log.Error("creSvc.juryByMIDs(%s) error(%v)", xstr.JoinInts(v.MIDS), err)
		httpCode(c, err)
		return
	}
	now := time.Now()
	if err := creSvc.DB.Model(&blocked.Jury{}).Where("mid IN(?)", v.MIDS).Updates(map[string]interface{}{"status": v.Status, "expired": now, "oper_id": v.OID}).Error; err != nil {
		httpCode(c, err)
		return
	}
	log.Info("business_upUserStatus:%+v", v)
	mids := make(map[int64]*blocked.Jury, len(v.MIDS))
	for _, v := range items {
		mids[v.UID] = v
	}
	var multiple []interface{}
	for _, mid := range v.MIDS {
		single := map[string]interface{}{
			"id":      mids[mid].ID,
			"oper_id": v.OID,
			"status":  v.Status,
			"expired": now.Format(model.TimeFormatSec),
			"black":   mids[mid].Black,
		}
		multiple = append(multiple, single)
	}
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedJury, blocked.TableBlockedJury, multiple)
	})
	httpCode(c, nil)
}

func blackWhite(c *bm.Context) {
	var err error
	v := new(struct {
		MIDS   []int64 `form:"mids,split"  validate:"min=1,max=100"`
		Status int8    `form:"status"  validate:"min=1,max=2"`
		Black  int8    `form:"black"`
		ReMark string  `form:"remark"`
		OID    int64   `form:"op_id"  validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	items := []*blocked.Jury{}
	if err = creSvc.ReadDB.Where("mid IN (?)", v.MIDS).Find(&items).Error; err != nil {
		log.Error("creSvc.juryByMIDs(%s) error(%v)", xstr.JoinInts(v.MIDS), err)
		httpCode(c, err)
		return
	}
	mids := make(map[int64]*blocked.Jury, len(v.MIDS))
	for _, v := range items {
		mids[v.UID] = v
	}
	now := time.Now()
	var multiple []interface{}
	switch {
	case v.Black == blocked.JuryBlack && v.Status == blocked.JuryStatusOn:
		err = creSvc.DB.Model(&blocked.Jury{}).Where("mid IN(?)", v.MIDS).Updates(
			map[string]interface{}{
				"status":  blocked.JuryStatusDown,
				"black":   blocked.JuryBlack,
				"oper_id": v.OID,
				"expired": time.Now().Unix(),
				"remark":  v.ReMark,
			}).Error
		for _, mid := range v.MIDS {
			single := map[string]interface{}{
				"id":      mids[mid].ID,
				"oper_id": v.OID,
				"status":  v.Status,
				"expired": now.Format(model.TimeFormatSec),
				"black":   blocked.JuryBlack,
			}
			multiple = append(multiple, single)
		}
	case v.Black == blocked.JuryWhite && v.Status == blocked.JuryStatusOn:
		err = creSvc.DB.Model(&blocked.Jury{}).Where("mid IN(?)", v.MIDS).Updates(
			map[string]interface{}{
				"black":   blocked.JuryWhite,
				"oper_id": v.OID,
				"remark":  v.ReMark,
			}).Error
		for _, mid := range v.MIDS {
			single := map[string]interface{}{
				"id":      mids[mid].ID,
				"oper_id": v.OID,
				"black":   blocked.JuryWhite,
			}
			multiple = append(multiple, single)
		}
	case (v.Black == blocked.JuryBlack && v.Status == blocked.JuryStatusDown) || (v.Black == blocked.JuryWhite && v.Status == blocked.JuryStatusDown):
		err = creSvc.DB.Model(&blocked.Jury{}).Where("mid IN(?)", v.MIDS).Updates(
			map[string]interface{}{
				"black":   blocked.JuryNormal,
				"oper_id": v.OID,
				"remark":  v.ReMark,
			}).Error
		for _, mid := range v.MIDS {
			single := map[string]interface{}{
				"id":      mids[mid].ID,
				"oper_id": v.OID,
				"black":   blocked.JuryNormal,
			}
			multiple = append(multiple, single)
		}
	}
	if err != nil {
		log.Error("blackWhite creSvc.DB error(%v)", err)
		httpCode(c, err)
		return
	}
	log.Info("business_blackWhite:%+v", v)
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedJury, blocked.TableBlockedJury, multiple)
	})
	httpCode(c, nil)
}

func usersEx(c *bm.Context) {
	v := new(blocked.ArgJurySearch)
	if err := c.Bind(v); err != nil {
		return
	}
	list, err := creSvc.JurysEx(c, v)
	if err != nil {
		log.Error("creSvc.InfosEx  error(%v)", err)
		httpCode(c, err)
		return
	}
	re, err := blocked.DealJury(list)
	if err != nil {
		log.Error("blocked.DealInfo  error(%v)", err)
		httpCode(c, err)
		return
	}
	outBuf(c, creSvc.FormatCSV(re), fmt.Sprintf("%s-%s", time.Now().Format(model.TimeFormatDay), "juryer"))
}

func cases(c *bm.Context) {
	v := new(blocked.ArgCaseSearch)
	if err := c.Bind(v); err != nil {
		return
	}
	list, pager, err := creSvc.Cases(c, v)
	if err != nil {
		log.Error("creSvc.Publishs(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	if len(list) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, list, pager)
}

// caseByID  get case by id.
func caseByID(c *bm.Context) {
	var err error
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.Case{}
	if err = creSvc.ReadDB.Where("id = ?", v.ID).Find(item).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("creSvc.caseByID error(%v)", err)
			httpCode(c, err)
			return
		}
		httpData(c, nil, nil)
		return
	}
	if item != nil {
		item.CaseTypeDesc = blocked.CaseTypeDesc[item.CaseType]
		item.StatusDesc = blocked.StatusDesc[item.Status]
		item.OriginTypeDesc = blocked.OriginTypeDesc[item.OriginType]
		item.ReasonTypeDesc = blocked.ReasonTypeDesc(item.ReasonType)
		item.PunishDesc = blocked.PunishDesc[item.PunishResult]
		item.VoteTotal = blocked.VoteTotal(item.VoteRule, item.VoteBreak, item.VoteDelete)
		if item.VoteRule > 0 || item.VoteBreak > 0 || item.VoteDelete > 0 {
			item.BlockedPercent = blocked.BreakPercent(item.VoteRule, item.VoteBreak, item.VoteDelete)
			item.DeletePercent = blocked.DeletePercent(item.VoteRule, item.VoteBreak, item.VoteDelete)
			item.RulePercent = blocked.RulePercent(item.VoteRule, item.VoteBreak, item.VoteDelete)
		}
	}
	httpData(c, item, nil)
}

// reasons get all case reason.
func reasons(c *bm.Context) {
	cr, err := creSvc.CaseReason(c)
	if err != nil {
		log.Error("creSvc.CaseReason error(%v)", err)
		httpCode(c, err)
		return
	}
	httpData(c, cr, nil)
}

func addCase(c *bm.Context) {
	var err error
	v := new(blocked.ArgCase)
	if err = c.Bind(v); err != nil {
		return
	}
	defTime, _ := time.ParseInLocation(model.TimeFormatSec, blocked.DefaultTime, time.Local)
	if err = creSvc.DB.Create(
		&blocked.Case{
			MID:           v.UID,
			OriginType:    v.Otype,
			ReasonType:    v.ReasonType,
			PunishResult:  v.PunishResult,
			BlockedDays:   v.BlockedDays,
			OriginTitle:   v.OriginTitle,
			OriginContent: v.OriginContent,
			OriginURL:     v.OriginURL,
			OPID:          v.OID,
			RelationID:    v.RelationID,
			Status:        blocked.CaseStatusGrantStop,
			StartTime:     xtime.Time(defTime.Unix()),
			EndTime:       xtime.Time(defTime.Unix()),
		}).Error; err != nil {
		log.Error("creSvc.addCase(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_addCase:%+v", v)
	httpCode(c, nil)
}

func upCase(c *bm.Context) {
	var err error
	v := new(blocked.ArgCase)
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.Case{}
	if err = creSvc.ReadDB.Where("id = ?", v.ID).Find(item).Error; err != nil {
		log.Error("creSvc.caseByID(%d) error(%v)", v.ID, err)
		httpCode(c, err)
		return
	}
	if err = creSvc.DB.Model(blocked.Case{}).Where("id = ?", v.ID).Updates(
		map[string]interface{}{
			"mid":            v.UID,
			"origin_type":    v.Otype,
			"reason_type":    v.ReasonType,
			"punish_result":  v.PunishResult,
			"blocked_days":   v.BlockedDays,
			"origin_title":   v.OriginTitle,
			"origin_content": v.OriginContent,
			"origin_url":     v.OriginURL,
			"relation_id":    v.RelationID,
			"oper_id":        v.OID,
		}).Error; err != nil {
		log.Error("creSvc.upCase(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_upCase:%+v", v)
	var multiple []interface{}
	single := map[string]interface{}{
		"id":          v.ID,
		"oper_id":     v.OID,
		"origin_type": v.Otype,
		"mid":         v.UID,
		"case_type":   item.CaseType,
	}
	multiple = append(multiple, single)
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedCase, blocked.TableBlockedCase, multiple)
	})
	httpCode(c, nil)
}

func addCaseVote(c *bm.Context) {
	var err error
	v := new(struct {
		ID         int64 `form:"id" validate:"required"`
		OPID       int64 `form:"op_id"  validate:"required"`
		VoteRule   int64 `form:"vote_rule" default:"0"`
		VoteBreak  int64 `form:"vote_break" default:"0"`
		VoteDelete int64 `form:"vote_delete" default:"0"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.Case{}
	if err = creSvc.ReadDB.Where("id = ?", v.ID).Find(item).Error; err != nil {
		log.Error("creSvc.caseByID(%d) error(%v)", v.ID, err)
		httpCode(c, err)
		return
	}
	if (v.VoteBreak < 0 && (float64(item.VoteBreak)-math.Abs(float64(v.VoteBreak))) < 0) ||
		(v.VoteRule < 0 && (float64(item.VoteRule)-math.Abs(float64(v.VoteRule))) < 0) ||
		(v.VoteDelete < 0 && (float64(item.VoteDelete)-math.Abs(float64(v.VoteDelete))) < 0) {
		httpCode(c, ecode.RequestErr)
		return
	}
	if err = creSvc.DB.Model(&blocked.Case{}).Where("id = ?", v.ID).UpdateColumns(
		map[string]interface{}{
			"vote_rule":   gorm.Expr("vote_rule + ?", v.VoteRule),
			"vote_break":  gorm.Expr("vote_break + ?", v.VoteBreak),
			"vote_delete": gorm.Expr("vote_delete + ?", v.VoteDelete),
			"oper_id":     v.OPID,
		}).Error; err != nil {
		log.Error("creSvc.addCaseVote(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_addCaseVote:%+v", v)
	httpCode(c, nil)
}

func upCaseStatus(c *bm.Context) {
	var err error
	v := new(blocked.ArgUpStatus)
	if err = c.Bind(v); err != nil {
		return
	}
	if err = creSvc.UpCaseStatus(c, v); err != nil {
		log.Error("creSvc.UpCaseStatus(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_upCaseStatus:%+v", v)
	httpCode(c, nil)
}

func addCaseType(c *bm.Context) {
	var err error
	v := new(struct {
		IDS  []int64 `form:"ids,split" validate:"min=1,max=200"`
		Type int8    `form:"type" validate:"min=0,max=1"`
		OID  int64   `form:"op_id" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if err = creSvc.DB.Model(blocked.Case{}).Where("id IN(?)", v.IDS).Updates(
		map[string]interface{}{
			"case_type": v.Type,
			"oper_id":   v.OID,
		}).Error; err != nil {
		log.Error("creSvc.addCaseType(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_addCaseType:%+v", v)
	var multiple []interface{}
	for _, id := range v.IDS {
		single := map[string]interface{}{
			"id":        id,
			"oper_id":   v.OID,
			"case_type": v.Type,
		}
		multiple = append(multiple, single)
	}
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedCase, blocked.TableBlockedCase, multiple)
	})
	httpCode(c, nil)
}

func opinions(c *bm.Context) {
	v := new(blocked.ArgOpinionSearch)
	if err := c.Bind(v); err != nil {
		return
	}
	list, pager, err := creSvc.Opinions(c, v)
	if err != nil {
		log.Error("creSvc.Opinions(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	if len(list) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, list, pager)
}

// opinionByID  get opinion by id.
func opinionByID(c *bm.Context) {
	var err error
	v := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.Opinion{}
	if err = creSvc.ReadDB.Where("id = ?", v.ID).Find(item).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("creSvc.opinionByID error(%v)", err)
			httpCode(c, err)
			return
		}
		httpData(c, nil, nil)
		return
	}
	if item != nil {
		item.AttrDesc = blocked.AttrDesc[item.Attr]
		item.VoteDesc = blocked.VoteDesc[item.Vote]
		item.VoteStateDesc = blocked.VoteStateDesc[item.State]
	}
	httpData(c, item, nil)
}

func delOpinions(c *bm.Context) {
	var err error
	v := new(struct {
		IDS    []int64 `form:"ids,split" validate:"min=1,max=20"`
		OID    int64   `form:"op_id"  validate:"required"`
		Status int8    `form:"status"  validate:"min=0,max=1" default:"1"`
		Send   int8    `form:"send" validate:"min=0,max=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	v.IDS = model.ArrayUnique(v.IDS)
	if len(v.IDS) == 0 {
		return
	}
	if err = creSvc.DB.Model(&blocked.Opinion{}).Where("id IN (?)", v.IDS).Updates(map[string]interface{}{"state": v.Status, "oper_id": v.OID}).Error; err != nil {
		log.Error("creSvc.delOpinions(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_delOpinions:%+v", v)
	var multiple []interface{}
	for _, id := range v.IDS {
		single := map[string]interface{}{
			"id":      id,
			"oper_id": v.OID,
			"state":   v.Status,
		}
		multiple = append(multiple, single)
	}
	creSvc.AddNotify(func() {
		creSvc.Search.SearchUpdate(context.TODO(), blocked.BusinessBlockedOpinion, blocked.TableBlockedOpinion, multiple)
	})
	if v.Send == blocked.NeedSendMsg {
		items := []*blocked.Opinion{}
		if err = creSvc.ReadDB.Where("id IN (?)", v.IDS).Find(&items).Error; err != nil {
			log.Error("creSvc.opinionByIDs(%s) error(%v)", xstr.JoinInts(v.IDS), err)
			err = nil
			return
		}
		for _, v := range items {
			msg := &blocked.SysMsg{
				Type:        blocked.MsgTypeDelOpinion,
				MID:         v.MID,
				CID:         v.CID,
				CaseContent: v.Content,
				RemoteIP:    metadata.String(c, metadata.RemoteIP),
			}
			creSvc.MsgCh <- msg
		}
	}
	httpCode(c, nil)
}

func kpis(c *bm.Context) {
	var (
		err        error
		start, end time.Time
		order      = "id"
		sort       = "DESC"
	)
	v := new(blocked.ArgKpiSearch)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Start == "" {
		if start, err = time.ParseInLocation(model.TimeFormatDay, time.Now().AddDate(0, -1, 0).Format(model.TimeFormatDay), time.Local); err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", time.Now().AddDate(0, -1, 0).Format(model.TimeFormatDay), err)
		}
	} else {
		if start, err = time.ParseInLocation(model.TimeFormatSec, v.Start, time.Local); err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", v.Start, err)
		}
	}
	if v.End == "" {
		if end, err = time.ParseInLocation(model.TimeFormatDay, time.Now().Format(model.TimeFormatDay), time.Local); err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", time.Now().Format(model.TimeFormatDay), err)
		}
	} else {
		if end, err = time.ParseInLocation(model.TimeFormatSec, v.End, time.Local); err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", v.End, err)
		}
	}
	count := int(0)
	items := []*blocked.KPI{}
	pager := &blocked.Pager{
		Total: count,
		PN:    v.PN,
		PS:    v.PS,
		Order: order,
		Sort:  sort,
	}
	if v.UID > 0 {
		err = creSvc.ReadDB.Where("mid=? and day BETWEEN ? AND ? ", v.UID, start, end).Offset((v.PN - 1) * v.PS).Limit(v.PS).Order(fmt.Sprintf("%s %s", order, sort)).Find(&items).Error
		creSvc.ReadDB.Model(&blocked.KPI{}).Where("mid=? and day BETWEEN ? AND ? ", v.UID, start, end).Count(&count)
	} else {
		err = creSvc.ReadDB.Where("day BETWEEN ? AND ? ", start, end).Offset((v.PN - 1) * v.PS).Limit(v.PS).Order(fmt.Sprintf("%s %s", order, sort)).Find(&items).Error
		creSvc.ReadDB.Model(&blocked.KPI{}).Where("day BETWEEN ? AND ? ", start, end).Count(&count)
	}
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("creSvc.kpis error(%v)", err)
			httpCode(c, err)
			return
		}
		httpData(c, nil, pager)
		return
	}
	pager.Total = count
	httpData(c, items, pager)
}

func kpiPoints(c *bm.Context) {
	v := new(blocked.ArgKpiPointSearch)
	if err := c.Bind(v); err != nil {
		return
	}
	list, pager, err := creSvc.KPIPoint(c, v)
	if err != nil {
		log.Error("creSvc.KPIPoint(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	if len(list) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, list, pager)
}

func kpisEx(c *bm.Context) {
	var (
		err         error
		start, end  time.Time
		items, list []*blocked.KPI
	)
	v := new(blocked.ArgKpiSearch)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Start == "" {
		if start, err = time.ParseInLocation(model.TimeFormatDay, time.Now().AddDate(0, -1, 0).Format(model.TimeFormatDay), time.Local); err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", time.Now().AddDate(0, -1, 0).Format(model.TimeFormatDay), err)
		}
	} else {
		if start, err = time.ParseInLocation(model.TimeFormatSec, v.Start, time.Local); err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", v.Start, err)
		}
	}
	if v.End == "" {
		if end, err = time.ParseInLocation(model.TimeFormatDay, time.Now().Format(model.TimeFormatDay), time.Local); err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", time.Now().Format(model.TimeFormatDay), err)
		}
	} else {
		if end, err = time.ParseInLocation(model.TimeFormatSec, v.End, time.Local); err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", v.End, err)
		}
	}
	ps, pn, loop, count := int64(500), int64(1), int64(0), int64(0)
	if v.UID > 0 {
		err = creSvc.ReadDB.Model(&blocked.KPI{}).Where("mid=? and day BETWEEN ? AND ?", v.UID, start, end).Count(&count).Error
	} else {
		err = creSvc.ReadDB.Model(&blocked.KPI{}).Where("day BETWEEN ? AND ?", start, end).Count(&count).Error
	}
	if err != nil {
		log.Error("creSvc.kpis count error(%v)", err)
		httpCode(c, err)
		return
	}
	if count < 1 {
		return
	}
	items = make([]*blocked.KPI, 0, count)
	list = make([]*blocked.KPI, 0, count)
	loop = (count / ps) + 1
	for pn = 1; pn <= loop; pn++ {
		if v.UID > 0 {
			err = creSvc.ReadDB.Where("mid=? and day BETWEEN ? AND ?", v.UID, start, end).Offset((pn - 1) * ps).Limit(ps).Order("id DESC").Find(&items).Error
		} else {
			err = creSvc.ReadDB.Where("day BETWEEN ? AND ?", start, end).Offset((pn - 1) * ps).Limit(ps).Order("id DESC").Find(&items).Error
		}
		if err != nil {
			log.Error("creSvc.kpis error(%v)", err)
			httpCode(c, err)
			return
		}
		list = append(list, items...)
	}
	sort.Slice(list, func(i int, j int) bool {
		return list[i].ID < list[j].ID
	})
	re, err := blocked.DealKPI(list)
	if err != nil {
		httpCode(c, err)
		return
	}
	outBuf(c, creSvc.FormatCSV(re), fmt.Sprintf("%s-%s", time.Now().Format(model.TimeFormatDay), "kpi"))
}

func caseConf(c *bm.Context) {
	cc, err := creSvc.CaseConf(c)
	if err != nil {
		log.Error("caseConfig error(%v)", err)
		httpCode(c, err)
		return
	}
	httpData(c, cc, nil)
}

func setCaseConf(c *bm.Context) {
	var err error
	v := new(blocked.ArgCaseConf)
	if err = c.Bind(v); err != nil {
		return
	}
	if err = creSvc.SetCaseConf(c, v); err != nil {
		log.Error("creSvc.SetCaseConf(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_setCaseConf:%+v", v)
	httpCode(c, nil)
}

func autoCaseConfig(c *bm.Context) {
	var err error
	v := new(struct {
		Platform int8 `form:"platform"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.AutoCase{}
	if err = creSvc.ReadDB.Where("platform=?", v.Platform).Find(&item).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("http.autoCaseConfig error(%v)", err)
			httpCode(c, err)
			return
		}
		httpData(c, nil, nil)
		return
	}
	if item.ReasonStr != "" {
		if item.Reasons, err = xstr.SplitInts(item.ReasonStr); err != nil {
			log.Error("xstr.SplitInts(%s) err(%v)", item.ReasonStr, err)
			httpCode(c, err)
			return
		}
	}
	httpData(c, item, nil)
}

func setAutoCaseConfig(c *bm.Context) {
	var err error
	v := new(blocked.ArgAutoCaseConf)
	if err = c.Bind(v); err != nil {
		return
	}
	item := &blocked.AutoCase{}
	creSvc.ReadDB.Where("id = ?", v.ID).Find(item)
	var autoCase *blocked.AutoCase
	if item != nil && item.ID != 0 {
		autoCase = &blocked.AutoCase{
			ID:          item.ID,
			Platform:    item.Platform,
			Likes:       v.Likes,
			ReasonStr:   xstr.JoinInts(v.Reasons),
			ReportScore: v.ReportScore,
			OPID:        v.OID,
			CTime:       item.CTime,
		}
	} else {
		autoCase = &blocked.AutoCase{
			Platform:    v.Platform,
			Likes:       v.Likes,
			ReasonStr:   xstr.JoinInts(v.Reasons),
			ReportScore: v.ReportScore,
			OPID:        v.OID,
		}
	}
	if err = creSvc.DB.Save(&autoCase).Error; err != nil {
		log.Error("http.setAutoCaseConfig error(%v)", err)
		httpCode(c, err)
		return
	}
	log.Info("business_setAutoCaseConfig:%+v", v)
	httpCode(c, nil)
}

func votenumConf(c *bm.Context) {
	content, err := creSvc.VotenumConf(c)
	if err != nil {
		log.Error("caseConfig error(%v)", err)
		httpCode(c, err)
		return
	}
	vc := blocked.VoteNum{}
	if err = json.Unmarshal([]byte(content), &vc); err != nil {
		log.Error("votenumConf_Unmarshal err:%+v", err)
		httpCode(c, err)
		return
	}
	httpData(c, vc, nil)
}

func setVotenumConf(c *bm.Context) {
	var err error
	v := new(blocked.ArgVoteNum)
	if err = c.Bind(v); err != nil {
		return
	}
	if err = creSvc.SetVotenumConf(c, v); err != nil {
		log.Error("creSvc.setVotenumConf(%+v) error(%v)", v, err)
		httpCode(c, err)
		return
	}
	log.Info("business_SetVotenumConf:%+v", v)
	httpCode(c, nil)
}

func webHook(c *bm.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("read body fail error(%v)", err)
		httpCode(c, ecode.RequestErr)
		return
	}
	if err = creSvc.WebHook(c, body); err != nil {
		log.Error("creSvc.WebHook error(%v)", err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}
