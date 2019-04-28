package data

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/web-show/conf"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// Dao struct
type Dao struct {
	client     *httpx.Client
	relatedURL string
}

// New Init
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		client:     httpx.NewClient(c.HTTPClient),
		relatedURL: "http://data.bilibili.co/recsys/related",
	}
	return
}

// Related check user bp.
func (dao *Dao) Related(c context.Context, aid int64, ip string) (aids []int64, err error) {
	params := url.Values{}
	params.Set("key", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
		Data []*struct {
			Value string `json:"value"`
		} `json:"data"`
	}
	if err = dao.client.Get(c, dao.relatedURL, ip, params, &res); err != nil {
		log.Error("realte url(%s) error(%v) ", dao.relatedURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("relate aids api failed(%d)", res.Code)
		log.Error("url(%s) res code(%d) or res.result(%v)", dao.relatedURL+"?"+params.Encode(), res.Code, res.Data)
		return
	}
	if res.Data == nil {
		err = nil
		return
	}
	// FIXME why two state
	if len(res.Data) > 0 {
		if aids, err = xstr.SplitInts(res.Data[0].Value); err != nil {
			log.Error("realte aids url(%s) error(%v)", dao.relatedURL+"?"+params.Encode(), err)
		}
	}
	return
}
