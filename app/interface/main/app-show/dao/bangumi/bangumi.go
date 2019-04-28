package bangumi

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model/bangumi"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_rcmmd       = "/api/get_season_by_tag"
	_seasonidURL = "/api/inner/archive/aid2seasonid"
	_bannerURL   = "/jsonp/slideshow/%d.ver"
)

// Dao is bangumi dao
type Dao struct {
	client      *httpx.Client
	clientAsyn  *httpx.Client
	rcmmd       string
	seasonidURL string
	bannerURL   string
	// grpc
	rpcClient seasongrpc.SeasonClient
}

// New bangumi dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:      httpx.NewClient(c.HTTPClient),
		clientAsyn:  httpx.NewClient(c.HTTPClientAsyn),
		rcmmd:       c.Host.Bangumi + _rcmmd,
		seasonidURL: c.Host.Bangumi + _seasonidURL,
		bannerURL:   c.Host.Bangumi + _bannerURL,
	}
	var err error
	if d.rpcClient, err = seasongrpc.NewClient(c.PGCRPC); err != nil {
		panic(fmt.Sprintf("seasongrpc NewClientt error (%+v)", err))
	}
	return
}

// Recommend get bangumi's recommend.
func (d *Dao) Recommend(now time.Time) (bgms []*bangumi.Bangumi, err error) {
	params := url.Values{}
	params.Set("tag_id", "109")
	params.Set("page", "1")
	params.Set("pagesize", "50")
	params.Set("indexType", "0")
	params.Set("build", "app-api")
	params.Set("platform", "Golang")
	var res struct {
		Code   int                `json:"code"`
		Result []*bangumi.Bangumi `json:"result"`
	}
	if err = d.clientAsyn.Get(context.TODO(), d.rcmmd, "", params, &res); err != nil {
		log.Error("bangumi url(%s) error(%v)", d.rcmmd+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("bangumi recommend api failed(%d)", res.Code)
		log.Error("url(%s) res code(%d) or res.result(%v)", d.rcmmd+"?"+params.Encode(), res.Code, res.Result)
		return
	}
	bgms = res.Result
	return
}

// Seasonid
func (d *Dao) Seasonid(aids []int64, now time.Time) (data map[int64]*bangumi.SeasonInfo, err error) {
	var (
		aidStr string
		msg1   = []byte(`,`)
		buf    bytes.Buffer
	)
	if len(aids) == 0 {
		log.Error("aids is null")
		return
	}
	for _, aid := range aids {
		buf.WriteString(strconv.FormatInt(aid, 10))
		buf.Write(msg1)
	}
	buf.Truncate(buf.Len() - 1)
	aidStr = buf.String()
	buf.Reset()
	params := url.Values{}
	params.Set("build", "app-api")
	params.Set("platform", "Golang")
	params.Set("aids", aidStr)
	var res struct {
		Code   int                           `json:"code"`
		Result map[int64]*bangumi.SeasonInfo `json:"result"`
	}
	if err = d.client.Get(context.TODO(), d.seasonidURL, "", params, &res); err != nil {
		log.Error("bangumi seasonid url(%s) error(%v)", d.seasonidURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("bangumi seasonid api failed(%d)", res.Code)
		log.Error("url(%s) res code(%d) or res.result(%v)", d.seasonidURL+"?"+params.Encode(), res.Code, res.Result)
		return
	}
	data = res.Result
	return
}

// Banners pgc banners
func (d *Dao) Banners(c context.Context, pgcID int) (data []*bangumi.Banner, err error) {
	var res struct {
		Code   int               `json:"code"`
		Result []*bangumi.Banner `json:"result"`
	}
	api := fmt.Sprintf(d.bannerURL, pgcID)
	if err = d.client.Get(c, api, "", nil, &res); err != nil {
		log.Error("bangumi banner url(%s) error(%v)", api, err)
	}
	if res.Code != 0 {
		log.Error("bangumi banner url(%s) error(%v)", api, res.Code)
		err = fmt.Errorf("bangumi banner api response code(%v)", res)
		return
	}
	data = res.Result
	return
}
