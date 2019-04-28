package dao

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/job/main/credit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

// SendPendant send pendant.
func (d *Dao) SendPendant(c context.Context, mid int64, pids []int64, day int64) (err error) {
	params := url.Values{}
	params.Set("pids", xstr.JoinInts(pids))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var days []int64
	for range pids {
		days = append(days, day)
	}
	params.Set("expires", xstr.JoinInts(days))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.sendPendantURL, "", params, &res); err != nil {
		log.Error("d.sendPendantURL url(%s) res(%d) error(%v)", d.sendPendantURL+"?"+params.Encode(), res.Code, err)
		return
	}
	if res.Code != 0 {
		log.Error("d.sendPendantURL url(%s) code(%d)", d.sendPendantURL+"?"+params.Encode(), res.Code)
		err = errors.New("sendPendant failed")
	}
	return
}

// SendMedal send Medal.
func (d *Dao) SendMedal(c context.Context, mid int64, nid int64) (err error) {
	params := url.Values{}
	params.Set("nid", strconv.FormatInt(nid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.sendMedalURL, "", params, &res); err != nil {
		log.Error("d.sendMedalURL url(%s) res(%d) error(%v)", d.sendMedalURL+"?"+params.Encode(), res.Code, err)
		return
	}
	if res.Code != 0 {
		log.Error("d.sendMedalURL url(%s) code(%d)", d.sendMedalURL+"?"+params.Encode(), res.Code)
		err = errors.New("sendMedalURL failed")
	}
	return
}

// DelTag del tag.
func (d *Dao) DelTag(c context.Context, tid, aid string) (err error) {
	params := url.Values{}
	params.Set("aid", aid)
	params.Set("tag_id", tid)
	params.Set("mid", "0")
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.delTagURL, "", params, &res); err != nil {
		log.Info("d.delTagURL url(%s) res(%v) error(%v)", d.delTagURL+"?"+params.Encode(), res, err)
		return
	}
	if res.Code != 0 {
		log.Error("d.delTagURL url(%s) code(%d)", d.delTagURL+"?"+params.Encode(), res.Code)
	}
	log.Info("d.delTagURL url(%s) res(%d)", d.delTagURL+"?"+params.Encode(), res.Code)
	return
}

// ReportDM report dm is delete.
func (d *Dao) ReportDM(c context.Context, cid string, dmid string, result int64) (err error) {
	params := url.Values{}
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("cid", cid)
	params.Set("dmid", dmid)
	params.Set("result", strconv.FormatInt(result, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.delDMURL, "", params, &res); err != nil {
		log.Error("d.delDMURL url(%s) error(%v)", d.delDMURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.delDMURL url(%s) code(%d)", d.delDMURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// BlockAccount block account.
func (d *Dao) BlockAccount(c context.Context, r *model.BlockedInfo) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(r.UID, 10))
	params.Set("source", "1")
	if int8(r.PunishType) == model.PunishTypeForever && r.BlockedForever == model.InBlockedForever {
		params.Set("action", "2")
	} else {
		params.Set("action", "1")
	}
	if r.CaseID == 0 {
		params.Set("notify", "1")
	} else {
		params.Set("notify", "0")
	}
	params.Set("duration", strconv.FormatInt(r.BlockedDays*86400, 10))
	params.Set("start_time", fmt.Sprintf("%d", time.Now().Unix()))
	params.Set("op_id", strconv.FormatInt(r.OPID, 10))
	params.Set("operator", r.OperatorName)
	params.Set("reason", model.ReasonTypeDesc(int8(r.ReasonType)))
	params.Set("comment", r.BlockedRemark)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.blockAccountURL, "", params, &res); err != nil {
		log.Error("d.blockAccountURL res(%d) url(%s)", res.Code, d.blockAccountURL+"?"+params.Encode())
	}
	log.Info("d.blockAccountURL res(%d) url(%s)", res.Code, d.blockAccountURL+"?"+params.Encode())
	return
}

// UnBlockAccount  unblock account.
func (d *Dao) UnBlockAccount(c context.Context, r *model.BlockedInfo) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(r.UID, 10))
	params.Set("source", "1")
	params.Set("op_id", strconv.FormatInt(r.OPID, 10))
	params.Set("operator", r.OperatorName)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.unBlockAccountURL, "", params, &res); err != nil {
		log.Error("d.unBlockAccountURL res(%d) url(%s)", res.Code, d.unBlockAccountURL+"?"+params.Encode())
	}
	log.Info("d.unBlockAccountURL res(%d) url(%s)", res.Code, d.unBlockAccountURL+"?"+params.Encode())
	return
}

// SendMsg send message.
func (d *Dao) SendMsg(c context.Context, mid int64, title string, context string) (err error) {
	params := url.Values{}
	params.Set("mc", "2_1_13")
	params.Set("title", title)
	params.Set("data_type", "4")
	params.Set("context", context)
	params.Set("mid_list", fmt.Sprintf("%d", mid))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Status int8   `json:"status"`
			Remark string `json:"remark"`
		} `json:"data"`
	}
	if err = d.client.Post(c, d.sendMsgURL, "", params, &res); err != nil {
		log.Error("sendMsgURL(%s) error(%v)", d.sendMsgURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("sendMsgURL(%s) res(%d)", d.sendMsgURL+"?"+params.Encode(), res.Code)
	}
	log.Info("d.sendMsgURL url(%s) res(%d)", d.sendMsgURL+"?"+params.Encode(), res.Code)
	return
}

// Sms send monitor sms.
func (d *Dao) Sms(c context.Context, phone, token, msg string) (err error) {
	var (
		urlStr = "http://ops-mng.bilibili.co/api/sendsms"
		res    struct {
			Result bool `json:"result"`
		}
	)
	params := url.Values{}
	params.Set("phone", phone)
	params.Set("message", msg)
	params.Set("token", token)
	if err = d.client.Get(c, urlStr, "", params, &res); err != nil {
		log.Error("d.Sms url(%s) res(%v) err(%v)", urlStr+"?"+params.Encode(), res, err)
		return
	}
	if !res.Result {
		log.Error("ops-mng sendsms url(%s) result(%v)", urlStr+"?"+params.Encode(), res.Result)
	}
	log.Info("d.Sms url(%s) res(%v)", urlStr+"?"+params.Encode(), res)
	return
}

// AddMoral add or reduce moral to user
func (d *Dao) AddMoral(c context.Context, mid int64, moral float64, reasonType int8, oper, reason, remark, remoteIP string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("addMoral", strconv.FormatFloat(moral, 'f', -1, 64))
	if moral > 0 {
		params.Set("origin", "1")
	} else {
		params.Set("origin", "2")
	}
	if oper == "" {
		params.Set("operater", "系统")
	} else {
		params.Set("operater", oper)
	}
	params.Set("reason", reason)
	params.Set("reason_type", strconv.FormatInt(int64(reasonType), 10))
	params.Set("remark", remark)
	params.Set("is_notify", "1")
	var res struct {
		Code int `json:"code"`
	}
	err = d.client.Get(c, d.addMoralURL, remoteIP, params, &res)
	if err != nil {
		log.Error("d.addMoralURL url(%s) error(%v)", d.addMoralURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.addMoralURL url(%s) error(%v)", d.addMoralURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	log.Info("d.addMoralURL url(%s) res(%v)", d.addMoralURL+"?"+params.Encode(), res)
	return
}

// AddMoney  Modify user Coins（old）.
func (d *Dao) AddMoney(c context.Context, mid int64, money float64, reason string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("count", strconv.FormatFloat(money, 'f', -1, 64))
	params.Set("reason", reason)
	var res struct {
		Code int `json:"code"`
	}
	err = d.client.Post(c, d.modifyCoinsURL, "", params, &res)
	log.Info("d.modifyCoinsURL url(%s) res(%v)", d.modifyCoinsURL+"?"+params.Encode(), res)
	if err != nil {
		log.Error("d.modifyCoinsURL url(%s) error(%v)", d.modifyCoinsURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.modifyCoinsURL url(%s) error(%v)", d.modifyCoinsURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// CheckFilter check content filter.
func (d *Dao) CheckFilter(c context.Context, area, msg, ip string) (need bool, err error) {
	params := url.Values{}
	params.Set("area", area)
	params.Set("msg", msg)
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Level int    `json:"level"`
			Limit int    `json:"limit"`
			Msg   string `json:"msg"`
		} `json:"data"`
	}
	err = d.client.Get(c, d.filterURL, ip, params, &res)
	if err != nil {
		log.Error("d.filterURL url(%s) error(%v)", d.filterURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.filterURL url(%s) error(%v)", d.filterURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data == nil {
		log.Warn("d.filterURL url(%s) res(%+v)", d.filterURL+"?"+params.Encode(), res.Data)
		return
	}
	need = res.Data.Level >= 20
	return
}

// UpAppealState .
func (d *Dao) UpAppealState(c context.Context, bid, oid, eid int64) (err error) {
	params := url.Values{}
	params.Set("business", strconv.FormatInt(bid, 10))
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("eid", strconv.FormatInt(eid, 10))
	var res struct {
		Code int `json:"code"`
	}
	err = d.client.Post(c, d.upAppealStateURL, "", params, &res)
	log.Info("d.upAppealStateURL url(%s) res(%v)", d.upAppealStateURL+"?"+params.Encode(), res)
	if err != nil {
		log.Error("d.upAppealStateURL url(%s) error(%v)", d.upAppealStateURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("d.upAppealStateURL url(%s) error(%v)", d.upAppealStateURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	return
}
