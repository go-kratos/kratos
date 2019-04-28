package message

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/message"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	_getUpListURI = "/api/notify/get.up.list.do"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// uri
	getUpListURL string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		client:       httpx.NewClient(c.HTTPClient.Normal),
		getUpListURL: c.Host.Message + _getUpListURI,
	}
	return
}

// GetUpList fn
func (d *Dao) GetUpList(c context.Context, mid int64, ak, ck, ip string) (data []*message.Message, err error) {
	data = make([]*message.Message, 0)
	var (
		res struct {
			Code int                `json:"code"`
			Data []*message.Message `json:"data"`
		}
		req *http.Request
	)
	params := url.Values{}
	params.Set("access_key", ak)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var (
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.getUpListURL + "?" + query
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Error("http.NewRequest(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeMessageAPIErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	if len(ck) > 0 {
		req.Header.Set("Cookie", ck)
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v); mid(%d), ip(%s)", url, err, mid, ip)
		err = ecode.CreativeMessageAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("GetUpList res.code!=0 url(%s) res(%v), mid(%d), ip(%s)", url, res, mid, ip)
		err = ecode.CreativeMessageAPIErr
		return
	}
	for _, v := range res.Data {
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", v.TimeAt, time.Local)
		v.TimeStamp = t.Unix()
	}
	data = res.Data
	return
}
