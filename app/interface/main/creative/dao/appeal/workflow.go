package appeal

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/creative/model/appeal"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// appeal url
const (
	_list           = "/x/internal/workflow/appeal/list"
	_detail         = "/x/internal/workflow/appeal/info"
	_addappeal      = "/x/internal/workflow/appeal/add"
	_addreply       = "/x/internal/workflow/appeal/reply/add"
	_appealstar     = "/x/internal/workflow/appeal/extra/up"
	_appealstate    = "/x/internal/workflow/appeal/state"
	_appealStarInfo = "/x/internal/workflow/appeal/extra/info"
)

// AppealList appeal list .
func (d *Dao) AppealList(c context.Context, mid int64, business int, ip string) (as []*appeal.AppealMeta, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("business", strconv.Itoa(business))
	var res struct {
		Code int                  `json:"code"`
		Data []*appeal.AppealMeta `json:"data"`
	}
	if err = d.client.Get(c, d.list, ip, params, &res); err != nil {
		err = errors.Wrap(err, "appeal list api")
		return
	}
	if res.Code != 0 {
		log.Error("appeal list  url(%s) mid(%d) res(%v)", d.list+"?"+params.Encode(), mid, res)
		err = ecode.Int(res.Code)
		return
	}
	as = res.Data
	return
}

// AppealDetail appeal detail.
func (d *Dao) AppealDetail(c context.Context, mid, id int64, business int, ip string) (a *appeal.AppealMeta, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("id", strconv.FormatInt(id, 10))
	params.Set("business", strconv.Itoa(business))
	var res struct {
		Code int                `json:"code"`
		Data *appeal.AppealMeta `json:"data"`
	}
	if err = d.client.Get(c, d.detail, ip, params, &res); err != nil {
		err = errors.Wrap(err, "appeal api")
		return
	}
	if res.Code != 0 {
		log.Error("appeal url(%s) mid(%d) res(%v)", d.detail+"?"+params.Encode(), mid, res)
		err = ecode.Int(res.Code)
		return
	}
	a = res.Data
	return
}

// AddAppeal add appeal.
func (d *Dao) AddAppeal(c context.Context, tid, aid, mid, business int64, qq, phone, email, desc, attachments, ip string, ap *appeal.BusinessAppeal) (apID int64, err error) {
	params := url.Values{}
	params.Set("tid", strconv.FormatInt(tid, 10))
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("business", strconv.FormatInt(business, 10))
	params.Set("qq", qq)
	params.Set("phone", phone)
	params.Set("email", email)
	params.Set("description", desc)
	params.Set("attachments", attachments)
	params.Set("business_typeid", strconv.FormatInt(ap.BusinessTypeID, 10))
	params.Set("business_mid", strconv.FormatInt(ap.BusinessMID, 10))
	params.Set("business_title", ap.BusinessTitle)
	params.Set("business_content", ap.BusinessContent)
	params.Set("business_state", strconv.FormatInt(appeal.StateCreate, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			ChallengeNo int64 `json:"challengeNo"`
		} `json:"data"`
	}
	log.Info("AddAppeal params(%v)", params)
	if err = d.client.Post(c, d.addappeal, ip, params, &res); err != nil {
		err = errors.Wrap(err, "appeal api")
		return
	}
	if res.Code != 0 {
		log.Error("add appeal  url(%s) mid(%d) res(%v)", d.addappeal+"?"+params.Encode(), mid, res)
		err = ecode.Int(res.Code)
	}
	apID = res.Data.ChallengeNo
	return
}

// AddReply add appeal reply .
func (d *Dao) AddReply(c context.Context, cid, event int64, content, attachments, ip string) (err error) {
	params := url.Values{}
	params.Set("cid", strconv.FormatInt(cid, 10))
	params.Set("content", content)
	params.Set("attachments", attachments)
	params.Set("event", strconv.FormatInt(event, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.addreply, ip, params, &res); err != nil {
		err = errors.Wrap(err, "appeal reply api")
		return
	}
	if res.Code != 0 {
		log.Error("add appeal  url(%s)  res(%v)", d.addreply+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
	}
	return
}

// AppealExtra modify appeal extra.
func (d *Dao) AppealExtra(c context.Context, mid, cid, business, val int64, key, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("cid", strconv.FormatInt(cid, 10))
	params.Set("business", strconv.FormatInt(business, 10))
	params.Set("key", key)
	params.Set("val", strconv.FormatInt(val, 10))
	var res struct {
		Code int `json:"code"`
	}
	log.Info("AppealExtra params(%v)", params)
	if err = d.client.Post(c, d.appealstar, ip, params, &res); err != nil {
		err = errors.Wrap(err, "appeal star api")
		return
	}
	if res.Code != 0 {
		log.Error("appeal start url(%s)  res(%v)", d.appealstar+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
	}
	return
}

// AppealState modify appeal state .
func (d *Dao) AppealState(c context.Context, mid, id, business, state int64, ip string) (err error) {
	params := url.Values{}
	params.Set("id", strconv.FormatInt(id, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("business", strconv.FormatInt(business, 10))
	params.Set("business_state", strconv.FormatInt(state, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.appealstate, ip, params, &res); err != nil {
		err = errors.Wrap(err, "appeal state api")
		return
	}
	if res.Code != 0 {
		log.Error("appeal state url(%s)  res(%v)", d.appealstate+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
	}
	return
}

// AppealStarInfo appeal star detail.
func (d *Dao) AppealStarInfo(c context.Context, mid, cid int64, business int, ip string) (star, etime string, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("cid", strconv.FormatInt(cid, 10))
	params.Set("business", strconv.Itoa(business))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Star  string `json:"star"`
			ETime string `json:"etime"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.appealStarInfo, ip, params, &res); err != nil {
		err = errors.Wrap(err, "appeal api")
		return
	}
	if res.Code != 0 {
		log.Error("appeal url(%s) mid(%d) res(%v)", d.detail+"?"+params.Encode(), mid, res)
		err = ecode.Int(res.Code)
		return
	}
	star = res.Data.Star
	etime = res.Data.ETime
	return
}
