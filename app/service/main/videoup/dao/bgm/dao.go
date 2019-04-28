package bgm

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/service/main/videoup/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_bgmBindURI = "/x/internal/v1/audio/song_video_relation/add"
)

// Dao is redis dao.
type Dao struct {
	c          *conf.Config
	httpW      *bm.Client
	bgmBindURL string
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		httpW:      bm.NewClient(c.HTTPClient.Write),
		bgmBindURL: c.Host.APICO + _bgmBindURI,
	}
	return d
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return
}

// Bind aid,sid,cid bind in one
func (d *Dao) Bind(c context.Context, aid, sid, cid int64) (err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("sid", strconv.FormatInt(sid, 10))
	params.Set("cid", strconv.FormatInt(cid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.bgmBindURL, "", params, &res); err != nil {
		log.Error("d.httpW.Post(%s) error(%v)", d.bgmBindURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) code(%d)", d.bgmBindURL+"?"+params.Encode(), res.Code)
	}
	return
}
