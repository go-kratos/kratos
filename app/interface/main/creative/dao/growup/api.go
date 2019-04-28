package growup

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/model/growup"
	"go-common/library/ecode"
	"go-common/library/log"
)

// UpStatus get up info.
func (d *Dao) UpStatus(c context.Context, mid int64, ip string) (us *growup.UpStatus, err error) {
	us = &growup.UpStatus{}
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data *growup.UpStatus
	}
	if err = d.client.Get(c, d.upStatusURL, ip, params, &res); err != nil {
		log.Error("growup url(%s) response(%v) error(%v)", d.upStatusURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("growup url(%s) res(%v)", d.upStatusURL, res)
		err = ecode.CreativeOrderAPIErr
		return
	}
	us = res.Data
	return
}

// UpInfo get up info.
func (d *Dao) UpInfo(c context.Context, mid int64, ip string) (ui *growup.UpInfo, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data *growup.UpInfo
	}
	if err = d.client.Get(c, d.upInfoURL, ip, params, &res); err != nil {
		log.Error("growup url(%s) response(%v) error(%v)", d.upInfoURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("growup url(%s) res(%v)", d.upInfoURL, res)
		err = ecode.CreativeOrderAPIErr
		return
	}
	ui = res.Data
	return
}

//Join join growup.
func (d *Dao) Join(c context.Context, mid int64, accTy, signTy int, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("account_type", strconv.Itoa(accTy))
	params.Set("sign_type", strconv.Itoa(signTy))
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.client.Post(c, d.joinURL, ip, params, &res); err != nil {
		log.Error("growup url(%s) mid(%d) error(%v)", d.joinURL+"?"+params.Encode(), mid, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("growup url(%s) mid(%d) res.code(%d) error(%v)", d.joinURL+"?"+params.Encode(), mid, res.Code, err)
		err = ecode.CreativeOrderAPIErr
	}
	return
}

//Quit quit growup.
func (d *Dao) Quit(c context.Context, mid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.client.Post(c, d.quitURL, ip, params, &res); err != nil {
		log.Error("growup url(%s) mid(%d) error(%v)", d.quitURL+"?"+params.Encode(), mid, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("growup url(%s) mid(%d) res.code(%d) error(%v)", d.quitURL+"?"+params.Encode(), mid, res.Code, err)
		err = ecode.CreativeOrderAPIErr
	}
	return
}

// Summary income.
func (d *Dao) Summary(c context.Context, mid int64, ip string) (sm *growup.Summary, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	var res struct {
		Code    int `json:"code"`
		Data    *growup.Summary
		Message string `json:"message"`
	}
	if err = d.client.Get(c, d.summaryURL, ip, params, &res); err != nil {
		log.Error("growup url(%s) response(%v) error(%v)", d.summaryURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("growup url(%s) res(%v)", d.summaryURL, res)
		err = ecode.CreativeOrderAPIErr
		return
	}
	sm = res.Data
	return
}

// Stat income by month.
func (d *Dao) Stat(c context.Context, mid int64, ip string) (st *growup.Stat, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	var res struct {
		Code    int `json:"code"`
		Data    *growup.Stat
		Message string `json:"message"`
	}
	if err = d.client.Get(c, d.statURL, ip, params, &res); err != nil {
		log.Error("growup url(%s) response(%v) error(%v)", d.statURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("growup url(%s) res(%v)", d.statURL, res)
		err = ecode.CreativeOrderAPIErr
		return
	}
	st = res.Data
	return
}

// IncomeList income by video/article/music.
func (d *Dao) IncomeList(c context.Context, mid int64, ty, pn, ps int, ip string) (il *growup.IncomeList, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("type", strconv.Itoa(ty))
	params.Set("page", strconv.Itoa(pn))
	params.Set("size", strconv.Itoa(ps))
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	var res struct {
		Code    int `json:"code"`
		Data    *growup.IncomeList
		Message string `json:"message"`
	}
	if err = d.client.Get(c, d.arcURL, ip, params, &res); err != nil {
		log.Error("growup url(%s) response(%v) error(%v)", d.arcURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("growup url(%s) res(%v)", d.arcURL, res)
		err = ecode.CreativeOrderAPIErr
		return
	}
	il = res.Data
	return
}

// BreachList breach list.
func (d *Dao) BreachList(c context.Context, mid int64, pn, ps int, ip string) (bl *growup.BreachList, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("page", strconv.Itoa(pn))
	params.Set("size", strconv.Itoa(ps))
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	var res struct {
		Code    int `json:"code"`
		Data    *growup.BreachList
		Message string `json:"message"`
	}
	if err = d.client.Get(c, d.breachURL, ip, params, &res); err != nil {
		log.Error("growup url(%s) response(%v) error(%v)", d.breachURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeOrderAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("growup url(%s) res(%v)", d.breachURL, res)
		err = ecode.CreativeOrderAPIErr
		return
	}
	bl = res.Data
	return
}
