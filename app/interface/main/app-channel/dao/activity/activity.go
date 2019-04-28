package activity

import (
	"context"
	"net/url"

	"go-common/app/interface/main/app-channel/conf"
	"go-common/app/interface/main/app-channel/model/activity"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_activitys = "/activity/pages"
)

// Dao is activity dao.
type Dao struct {
	// http client
	client *bm.Client
	// activitys
	activitys string
}

// New new a activity dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		client:    bm.NewClient(c.HTTPClient),
		activitys: c.Host.Activity + _activitys,
	}
	return d
}

// Activitys activity or tpoci
func (d *Dao) Activitys(c context.Context, ids []int64) (actm map[int64]*activity.Activity, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("pids", xstr.JoinInts(ids))
	params.Set("http", "2")
	params.Set("platform", "pegasus")
	var res struct {
		Code int `json:"code"`
		Data struct {
			List []*activity.Activity `json:"list"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.activitys, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.activitys+"?"+params.Encode())
		return
	}
	actm = make(map[int64]*activity.Activity, len(res.Data.List))
	for _, act := range res.Data.List {
		actm[act.ID] = act
	}
	return
}
