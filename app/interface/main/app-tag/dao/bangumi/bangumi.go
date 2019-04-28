package bangumi

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-tag/conf"
	"go-common/app/interface/main/app-tag/model/bangumi"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_seasonidURL = "/api/inner/archive/aid2seasonid"
)

// Dao is bangumi dao
type Dao struct {
	client      *httpx.Client
	seasonidURL string
}

// New bangumi dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:      httpx.NewClient(c.HTTPClient),
		seasonidURL: c.Host.Bangumi + _seasonidURL,
	}
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
