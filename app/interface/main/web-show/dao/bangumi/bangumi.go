package bangumi

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/web-show/conf"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

// Dao struct
type Dao struct {
	client  *httpx.Client
	isbpURL string
}

// New Init
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		client:  httpx.NewClient(c.HTTPClient),
		isbpURL: c.Host.Bangumi + "/api/bp",
	}
	return
}

// IsBp check user bp.
func (dao *Dao) IsBp(c context.Context, mid, aid int64, ip string) (is bool) {
	params := url.Values{}
	params.Set("build", "web-show")
	params.Set("platform", "Golang")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code   int  `json:"code"`
		Result bool `json:"result"`
	}
	if err := dao.client.Get(c, dao.isbpURL, ip, params, &res); err != nil {
		log.Error("bangumi url(%s) error(%v) ", dao.isbpURL+"?"+params.Encode(), err)
		return
	}
	// FIXME why two state
	if res.Code != 0 && res.Code != 1 {
		log.Error("bangumi Isbp api fail(%d)", res.Code)
		return
	}
	is = res.Result
	return
}
