package dynamic

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/videoup/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_userCheckURI   = "/lottery_svr/v0/lottery_svr/user_check"
	_lotteryBindURI = "/lottery_svr/v0/lottery_svr/bind"
)

// Dao  define
type Dao struct {
	c              *conf.Config
	client         *bm.Client
	LotteryBindURL string
	UserCheckURL   string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:              c,
		client:         bm.NewClient(c.HTTPClient.Write),
		LotteryBindURL: c.Host.Dynamic + _lotteryBindURI,
		UserCheckURL:   c.Host.Dynamic + _userCheckURI,
	}
	return
}

// LotteryBind fn
func (d *Dao) LotteryBind(c context.Context, lotteryID, aid, mid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("lottery_id", strconv.FormatInt(lotteryID, 10))
	params.Set("business_type", "8")
	params.Set("business_id", strconv.FormatInt(aid, 10))
	params.Set("sender_uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.LotteryBindURL, ip, params, &res); err != nil {
		log.Error("LotteryBind url(%s) response(%s) error(%v)", d.LotteryBindURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLotteryAPIErr
		return
	}
	log.Info("LotteryBind d.LotteryBindURL url(%s)", d.LotteryBindURL+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("LotteryBind url(%s) res(%v)", d.LotteryBindURL, res)
		err = ecode.CreativeLotteryAPIErr
		return
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
		log.Error("UserCheck url(%s) response(%s) error(%v)", d.UserCheckURL+"?"+params.Encode(), res, err)
		err = ecode.CreativeLotteryAPIErr
		return
	}
	log.Info("UserCheck d.UserCheckURL url(%s)", d.UserCheckURL+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("UserCheck url(%s) res(%v)", d.UserCheckURL, res)
		err = ecode.CreativeLotteryAPIErr
		return
	}
	ret = res.Data.Result
	return
}
