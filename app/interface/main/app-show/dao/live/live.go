package live

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	clive "go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model/live"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_live       = "/appIndex/recommendFeedList"
	_rec        = "/appIndex/recommendList"
	_topic      = "/topic/v0/Topic/hots"
	_dynamichot = "/dynamic_detail/v0/Dynamic/hot"
)

// Dao is live dao
type Dao struct {
	client     *httpx.Client
	clientAsyn *httpx.Client
	live       string
	rec        string
	topic      string
	dynamichot string
}

// New live dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:     httpx.NewClient(c.HTTPClient),
		clientAsyn: httpx.NewClient(c.HTTPClientAsyn),
		live:       c.Host.ApiLiveCo + _live,
		rec:        c.Host.ApiLiveCo + _rec,
		topic:      c.Host.Dynamic + _topic,
		dynamichot: c.Host.Dynamic + _dynamichot,
	}
	return
}

// Live feed
func (d *Dao) Feed(c context.Context, mid int64, ak, ip string, now time.Time) (r *live.Feed, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("access_key", ak)
	var res struct {
		Code int        `json:"code"`
		Data *live.Feed `json:"data"`
	}
	if err = d.client.Get(c, d.live, ip, params, &res); err != nil {
		log.Error("Feed url(%s) error(%v)", d.live+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("Feed url(%s) error(%v)", d.live+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("feed send failed")
		return
	}
	r = res.Data
	return
}

// Recommend get live Recommend data.
func (d *Dao) Recommend(now time.Time) (r *live.Recommend, err error) {
	params := url.Values{}
	params.Set("count", "60")
	var res struct {
		Code int             `json:"code"`
		Data *live.Recommend `json:"data"`
	}
	if err = d.clientAsyn.Get(context.TODO(), d.rec, "", params, &res); err != nil { // TODO context arg, service context.TODO
		log.Error("live recommend url(%s) error(%v)", d.rec+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("live recommend url(%s) error(%v)", d.rec+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("recommend send failed")
		return
	}
	r = res.Data
	return
}

// TopicHots get live topic hots
func (d *Dao) TopicHots(c context.Context) (topics []*clive.TopicHot, err error) {
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*clive.TopicHot `json:"list"`
		} `json:"data"`
	}
	if err = d.clientAsyn.Get(c, d.topic, "", nil, &res); err != nil {
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("topichots list url(%v) response(%s)", d.topic, b)
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.topic)
		return
	}
	for _, t := range res.Data.List {
		tmp := &clive.TopicHot{}
		*tmp = *t
		if err = tmp.TopicJSONChange(); err != nil {
			log.Error("TopicJSONChange error(%v)", err)
			return
		}
		topics = append(topics, tmp)
	}
	return
}

// DynamicHot get dynamic hot all
func (d *Dao) DynamicHot(c context.Context) (list []*clive.DynamicHot, err error) {
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*clive.DynamicHot `json:"list"`
		} `json:"data"`
	}
	if err = d.clientAsyn.Get(c, d.dynamichot, "", nil, &res); err != nil {
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("dynamichot list url(%v) response(%s)", d.dynamichot, b)
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.topic)
		return
	}
	list = res.Data.List
	return
}
