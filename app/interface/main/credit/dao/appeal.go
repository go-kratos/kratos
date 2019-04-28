package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/credit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// AddAppeal add appeal.
func (d *Dao) AddAppeal(c context.Context, tid, btid, oid, mid, business int64, content, reason string) (err error) {
	params := url.Values{}
	params.Set("tid", strconv.FormatInt(tid, 10))
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("business_mid", strconv.FormatInt(mid, 10))
	params.Set("business_typeid", strconv.FormatInt(btid, 10))
	params.Set("business", strconv.FormatInt(business, 10))
	params.Set("business_content", content)
	params.Set("description", reason)
	var res struct {
		Code int `json:"code"`
		Data *struct {
			ChallengeNo int64 `json:"challengeNo"`
		} `json:"data"`
	}
	if err = d.client.Post(c, d.addAppealURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		err = errors.Wrapf(err, "AddAppeal url(%s) res(%v)", d.addAppealURL+"?"+params.Encode(), res)
		return
	}
	if res.Code != 0 {
		log.Warn("add appeal  url(%s) mid(%d) res(%v)", d.addAppealURL+"?"+params.Encode(), mid, res)
		err = ecode.Int(res.Code)
	}
	return
}

// AppealList appeal list .
func (d *Dao) AppealList(c context.Context, mid int64, business int) (as []*model.Appeal, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("business", strconv.Itoa(business))
	var res struct {
		Code int             `json:"code"`
		Data []*model.Appeal `json:"data"`
	}
	if err = d.client.Get(c, d.appealListURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		err = errors.Wrapf(err, "AppealList url(%s) res(%v)", d.appealListURL+"?"+params.Encode(), res)
		return
	}
	if res.Code != 0 {
		log.Warn("appeal list  url(%s) mid(%d) res(%v)", d.appealListURL+"?"+params.Encode(), mid, res)
		err = ecode.Int(res.Code)
		return
	}
	as = res.Data
	return
}
