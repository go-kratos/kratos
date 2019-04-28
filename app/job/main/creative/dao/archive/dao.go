package archive

import (
	"context"
	"net/url"
	"strconv"

	article "go-common/app/interface/openplatform/article/rpc/client"
	"go-common/app/job/main/creative/conf"
	"go-common/app/job/main/creative/model"
	archive "go-common/app/service/main/archive/api/gorpc"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

// Dao is archive dao.
type Dao struct {
	// config
	c *conf.Config
	// rpc
	arc *archive.Service2
	art *article.Service
	// http client
	client  *httpx.Client
	viewURL string
}

// New init api url
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:   c,
		arc: archive.New2(c.ArchiveRPC),
		art: article.New(c.ArticleRPC),
		// http client
		client:  httpx.NewClient(c.HTTPClient.Normal),
		viewURL: c.Host.Videoup + "/videoup/view",
	}
	return
}

// Ping fn
func (d *Dao) Ping(c context.Context) (err error) {
	return
}

// Close fn
func (d *Dao) Close() (err error) {
	return
}

// View get archive
func (d *Dao) View(c context.Context, mid, aid int64) (av *model.ArcVideo, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int             `json:"code"`
		Data *model.ArcVideo `json:"data"`
	}
	if err = d.client.Get(c, d.viewURL, "", params, &res); err != nil {
		log.Error("archive.view url(%s) mid(%d) error(%v)", d.viewURL+"?"+params.Encode(), mid, err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.view url(%s) mid(%d) res(%v)", d.viewURL+"?"+params.Encode(), mid, res)
		err = ecode.Int(res.Code)
		return
	}
	if res.Data.Archive == nil {
		log.Error("archive.view url(%s) mid(%d) res(%v)", d.viewURL+"?"+params.Encode(), mid, res)
		return
	}
	av = res.Data
	return
}
