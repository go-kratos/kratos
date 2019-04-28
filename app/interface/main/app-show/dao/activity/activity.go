package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model/activity"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
	"net/url"
)

const (
	_activitys = "/activity/pages"
)

// Dao is activity dao.
type Dao struct {
	// http client
	client *httpx.Client
	// activitys
	activitys string
}

// New new a activity dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client:    httpx.NewClient(c.HTTPClientAsyn),
		activitys: c.Host.Activity + _activitys,
	}
	return d
}

func (d *Dao) Activitys(c context.Context, ids []int64, mold int, ip string) (actm map[int64]*activity.Activity, err error) {
	params := url.Values{}
	params.Set("pids", xstr.JoinInts(ids))
	params.Set("http", "2")
	params.Set("platform", "pegasus")
	params.Set("mold", strconv.Itoa(mold))
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*activity.Activity `json:"list"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.activitys, ip, params, &res); err != nil {
		log.Error("activitys url(%s) error(%v)", d.activitys+"?"+params.Encode(), err)
		return
	}
	b, _ := json.Marshal(&res)
	log.Info("activitys url(%s) response(%s)", d.activitys+"?"+params.Encode(), b)
	if res.Code != 0 {
		err = fmt.Errorf("activitys api failed(%d)", res.Code)
		log.Error("url(%s) res code(%d)", d.activitys+"?"+params.Encode(), res.Code)
		return
	}
	actm = make(map[int64]*activity.Activity, len(res.Data.List))
	for _, act := range res.Data.List {
		actm[act.ID] = act
	}
	return
}
