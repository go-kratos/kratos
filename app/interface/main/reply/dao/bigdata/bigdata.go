package bigdata

import (
	"context"
	"net/url"

	"go-common/app/interface/main/reply/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

type result struct {
	Code int `json:"code"`
}

// Dao bigdata dao
type Dao struct {
	url        string
	topicURL   string
	httpClient *httpx.Client
}

// New return a bigdata dao
func New(c *conf.Config) *Dao {
	d := &Dao{
		httpClient: httpx.NewClient(c.HTTPClient),
		url:        c.Reply.BigdataURL,
		topicURL:   c.Reply.AiTopicURL,
	}
	return d
}

// Filter Filter
func (dao *Dao) Filter(c context.Context, msg string) (err error) {
	params := url.Values{}
	res := &result{}
	params.Set("comment", msg)
	if err = dao.httpClient.Post(c, dao.url, "", params, res); err != nil {
		log.Error("Bigdata url(%s) error(%v)", dao.url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("Bigdata url(%s) error(%v)", dao.url+"?"+params.Encode(), res.Code)
		err = ecode.ReplyDeniedAsGarbage
	}
	return
}
