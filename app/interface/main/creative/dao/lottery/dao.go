package lottery

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/lottery"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"net/url"
	"strconv"
)

const (
	_userCheckURI = "/lottery_svr/v0/lottery_svr/user_check"
	_noticeURI    = "/lottery_svr/v0/lottery_svr/lottery_notice"
)

// Dao  define
type Dao struct {
	c                       *conf.Config
	client                  *bm.Client
	UserCheckURL, NoticeURL string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		client:       bm.NewClient(c.HTTPClient.Normal),
		UserCheckURL: c.Host.Dynamic + _userCheckURI,
		NoticeURL:    c.Host.Dynamic + _noticeURI,
	}
	return
}

// UserCheck fn
func (d *Dao) UserCheck(c context.Context, mid int64, ip string) (ret int, err error) {
	params := url.Values{}
	params.Set("sender_uid", strconv.FormatInt(mid, 10))
	params.Set("business_type", "8")
	var res struct {
		Code int `json:"code"`
		Data struct {
			Result int `json:"result"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.UserCheckURL, ip, params, &res); err != nil {
		log.Error("UserCheck url(%s) response(%v) error(%v)", d.UserCheckURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLotteryAPIErr
		return
	}
	log.Info("UserCheck d.UserCheckURL url(%s) code(%d)", d.UserCheckURL+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("UserCheck url(%s) res(%v)", d.UserCheckURL, res)
		err = ecode.CreativeLotteryAPIErr
		return
	}
	ret = res.Data.Result
	return
}

// Notice fn
func (d *Dao) Notice(c context.Context, aid, mid int64, ip string) (ret *lottery.Notice, err error) {
	ret = &lottery.Notice{}
	params := url.Values{}
	params.Set("sender_uid", strconv.FormatInt(mid, 10))
	params.Set("business_type", "8")
	params.Set("business_id", strconv.FormatInt(aid, 10))
	var res struct {
		Code int             `json:"code"`
		Data *lottery.Notice `json:"data"`
	}
	if err = d.client.Get(c, d.NoticeURL, ip, params, &res); err != nil {
		log.Error("Notice url(%s) response(%v) error(%v)", d.NoticeURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLotteryAPIErr
		return
	}
	log.Info("Notice d.NoticeURL url(%s) code(%d)", d.NoticeURL+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("Notice url(%s) res(%v)", d.NoticeURL, res)
		if res.Code == -9999 {
			err = ecode.NothingFound
		} else {
			err = ecode.CreativeLotteryAPIErr
		}
		return
	}
	ret = res.Data
	return
}
