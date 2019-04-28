package topic

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/favorite/conf"
	"go-common/app/interface/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

const _topic = "http://matsuri.bilibili.co/activity/pages"

// Dao defeine fav Dao
type Dao struct {
	httpClient *httpx.Client
}

// New return fav dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		httpClient: httpx.NewClient(c.HTTPClient),
	}
	return
}

// TopicMap return the user favorited topic's map data(all state).
func (d *Dao) TopicMap(c context.Context, tpIDs []int64, isNomal bool, appInfo *model.AppInfo) (data map[int64]*model.Topic, err error) {
	params := url.Values{}
	params.Set("mold", "1")
	if !isNomal {
		params.Set("all", "isOne")
	}
	params.Set("pids", xstr.JoinInts(tpIDs))
	if appInfo != nil {
		params.Set("http", strconv.Itoa(model.HttpMode4Https))
	} else {
		params.Set("http", strconv.Itoa(model.HttpMode4Both))
	}
	res := new(model.TopicsResult)
	ip := metadata.String(c, metadata.RemoteIP)
	if err = d.httpClient.Get(c, _topic, ip, params, res); err != nil {
		log.Error("d.HTTPClient.Get(%s?%s) error(%v)", _topic, params.Encode())
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.HTTPClient.Get(%s?%s) code:%d msg:%s", _topic, params.Encode(), res.Code)
		err = model.ErrTopicRequest
		return
	}
	data = make(map[int64]*model.Topic, len(res.Data.List))
	for _, r := range res.Data.List {
		data[r.ID] = r
	}
	return
}
