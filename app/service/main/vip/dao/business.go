package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_SendUserNotify = "/api/notify/send.user.notify.do"
	_PayAPIAdd      = "/api/coupon/regular/add"
	_CleanCache     = "/notify/cleanCacheAndNotify"
	_loginout       = "/intranet/acc/security/mid"
)

//Loginout login out
func (d *Dao) Loginout(c context.Context, mid int64) (err error) {

	val := url.Values{}
	val.Add("mids", strconv.FormatInt(mid, 10))
	val.Add("operator", strconv.FormatInt(mid, 10))
	val.Add("desc", "大会员解冻")
	resp := new(struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	})
	defer func() {
		log.Info("vip Loginout url:%+v params:%+v return:%+v", d.loginOutURL, val, resp)
	}()
	if err = d.client.Post(c, d.loginOutURL, "", val, resp); err != nil {
		err = errors.Errorf("vip Loginout url:%+v params:%+v return:%+v,err:%+v", d.loginOutURL, val, resp, err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
	}
	return
}

//SendCleanCache clean cache
func (d *Dao) SendCleanCache(c context.Context, mid int64, months int16, days int64, t int, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(int64(mid), 10))
	params.Set("months", strconv.FormatInt(int64(months), 10))
	params.Set("days", strconv.FormatInt(int64(days), 10))
	params.Set("type", strconv.FormatInt(int64(t), 10))

	if err = d.client.Get(c, d.vipURI, ip, params, nil); err != nil {
		log.Error("SendCleanCache error(%v) url(%v)", err, d.vipURI)
		return
	}
	return
}

//SendMultipMsg send multip msg
func (d *Dao) SendMultipMsg(c context.Context, mids, content, title, mc, ip string, dataType int) (err error) {
	params := url.Values{}
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("context", content)
	params.Set("data_type", strconv.FormatInt(int64(dataType), 10))
	params.Set("mid_list", mids)
	if err = d.client.Post(c, d.msgURI, ip, params, nil); err != nil {
		log.Error("SendMultipMsg error(%v)", err)
		return
	}
	return
}

//SendBcoinCoupon send bcoin coupon
func (d *Dao) SendBcoinCoupon(c context.Context, mids, activityID string, money int64, dueTime time.Time) (err error) {
	params := url.Values{}
	params.Set("activity_id", activityID)
	params.Set("mids", mids)
	params.Set("money", strconv.FormatInt(int64(money), 10))
	params.Set("due_time", dueTime.Format("2006-01-02"))
	if err = d.client.Post(c, d.payURI, "127.0.0.1", params, nil); err != nil {
		log.Error("SendBcoinCoupon error(%v)", err)
		fmt.Printf("SendBcoinCoupon error(%v)", err)
		return
	}
	return
}
