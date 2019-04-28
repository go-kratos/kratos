package sports

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"go-common/app/interface/main/activity/conf"
	"go-common/app/interface/main/activity/model/sports"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	xecocde "go-common/library/ecode"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/stat/prom"
)

const (
	_qqAppID    = "9"
	_qqAppKey   = "TWF0Y2hVbmlvbjpBUFBLRVk6OQ=="
	_newsAppID  = "openapi_for_bilibili"
	_newsAppKey = "d2c0d130c49baadc3d43fc731caecd43"
)

// PromError stat and log.
func PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}

// Dao dao.
type Dao struct {
	// http
	httpSports *xhttp.Client
	dClient    *http.Client
	// sports api
	sportsURI, newsURI string
	mc                 *memcache.Pool
	mcQqExpire         int32
}

// New dao new.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		httpSports: xhttp.NewClient(c.HTTPClientSports),
		dClient:    http.DefaultClient,
		sportsURI:  c.Host.Sports,
		newsURI:    c.Host.QqNews,
		mc:         memcache.NewPool(c.Memcache.Like),
		mcQqExpire: int32(time.Duration(c.Memcache.QqExpire) / time.Second),
	}
	return
}

// Qq get qq.
func (d *Dao) Qq(c context.Context, params url.Values, route string) (rs *json.RawMessage, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	var res struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	params.Del("route")
	params.Set("appId", _qqAppID)
	params.Set("appKey", _qqAppKey)
	if err = d.httpSports.Get(c, d.sportsURI+"/"+route, ip, params, &res); err != nil {
		log.Error("d.httpSports.Get(%s) err(%v)", d.sportsURI+"/"+route+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpSports.Get(%s) param(%v) ecode err(%d)", d.sportsURI, params, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	rs = &res.Data
	return
}

// QqNews get qq news.
func (d *Dao) QqNews(c context.Context, params url.Values, route string) (rs *sports.QqRes, err error) {
	var (
		req    *http.Request
		resp   *http.Response
		cancel func()
	)
	params.Set("chlid", "news_news_football")
	params.Set("appkey", _newsAppKey)
	params.Set("appid", _newsAppID)
	if req, err = http.NewRequest("GET", d.newsURI+"/"+route+"?"+params.Encode(), nil); err != nil {
		log.Error("QqNews http.NewRequest(%s) error(%v)", d.newsURI+"/"+route+"?"+params.Encode(), err)
		return
	}
	c, cancel = context.WithTimeout(c, time.Duration(conf.Conf.Rule.DTimeout))
	defer cancel()
	req = req.WithContext(c)
	if resp, err = d.dClient.Do(req); err != nil {
		log.Error("QqNews httpClient.Do(%s) error(%v)", d.newsURI+"/"+route+"?"+params.Encode(), err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		log.Error("QqNews url(%s) resp.StatusCode error(%v)", d.newsURI+"/"+route+"?"+params.Encode(), err)
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("QqNews ioutil.ReadAll() error(%v)", err)
		return
	} else if len(bs) == 0 {
		return
	}
	if e := json.Unmarshal(bs, &rs); e != nil {
		if e != io.EOF {
			log.Error("json decode body(%s) error(%v)", string(bs), e)
		}
	}
	if rs.Ret != 0 {
		err = xecocde.ActivityServerTimeout
		rs = nil
		return
	}
	return
}
