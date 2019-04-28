package vip

import (
	"context"
	"net/url"

	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/model/vip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_emojiURL = "/internal/v1/emoji/list"
)

type result struct {
	Code int `json:"code"`
	Data []*struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		URL    string `json:"path"`
		State  int8   `json:"deleted"`
		Emojis []*struct {
			ID     int64  `json:"id"`
			Name   string `json:"name"`
			URL    string `json:"path"`
			State  int8   `json:"deleted"`
			Remark string `json:"remark"`
		} `json:"emojiList"`
	} `json:"data"`
}

// Dao dao.
type Dao struct {
	emojiURL   string
	httpClient *bm.Client
}

// New new a dao and return.
func New(c *conf.Config) *Dao {
	d := &Dao{
		httpClient: bm.NewClient(c.HTTPClient),
		emojiURL:   c.Reply.VipURL + _emojiURL,
	}
	return d
}

// Emoji return emojis to cache
func (dao *Dao) Emoji(c context.Context) (emjs []*vip.Emoji, emjM map[string]int64, err error) {
	var res result
	params := url.Values{}
	if err = dao.httpClient.Get(c, dao.emojiURL, "", params, &res); err != nil {
		log.Error("vip url(%s) error(%v)", dao.emojiURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("vip url(%s) error(%v)", dao.emojiURL+"?"+params.Encode(), res.Code)
		return
	}
	emjs = make([]*vip.Emoji, 0, len(res.Data))
	emjM = make(map[string]int64)
	for _, d := range res.Data {
		var tmp = &vip.Emoji{
			Pid:    d.ID,
			Pname:  d.Name,
			Purl:   d.URL,
			Pstate: d.State,
		}
		emoTmps := make([]*vip.Face, 0, len(d.Emojis))
		for _, e := range d.Emojis {
			var emoTmp = &vip.Face{
				ID:    e.ID,
				Name:  e.Name,
				URL:   e.URL,
				State: e.State,
			}
			emoTmps = append(emoTmps, emoTmp)
			if d.State == 0 && e.State == 0 {
				emjM[e.Name] = e.ID
			}
		}
		tmp.Emojis = emoTmps
		emjs = append(emjs, tmp)
	}
	return
}
