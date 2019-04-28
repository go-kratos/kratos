package bplus

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_checkUser = "/promo_svr/v0/promo_svr/inner_user_check"
)

// Dao bplus
type Dao struct {
	client *httpx.Client
	// url
	checkUserURL string
}

// New bplus
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPClient),
		// url
		checkUserURL: c.Host.VC + _checkUser,
	}
	return
}

// UserCheck 动态互推入口白名单
func (d *Dao) UserCheck(c context.Context, mid int64) (ok bool, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Status int `json:"status"`
		}
	}
	if err = d.client.Get(c, d.checkUserURL, "", params, &res); err != nil {
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("UserCheck url(%s) response(%s)", d.checkUserURL+"?"+params.Encode(), b)
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.checkUserURL+"?"+params.Encode())
		return
	}
	if res.Data.Status == 1 {
		ok = true
	}
	return
}
