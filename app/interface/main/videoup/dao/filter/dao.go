package filter

import (
	"context"
	"net/url"

	"go-common/app/interface/main/videoup/conf"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_postFilter  = "/x/internal/filter/post"
	_postMFilter = "/x/internal/filter/mpost"
)

// Dao  define
type Dao struct {
	c              *conf.Config
	client         *httpx.Client
	postFilterURI  string
	postMFilterURI string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:              c,
		client:         httpx.NewClient(c.HTTPClient.FastRead),
		postFilterURI:  c.Host.APICo + _postFilter,
		postMFilterURI: c.Host.APICo + _postMFilter,
	}
	return
}

// VideoFilter fn
func (d *Dao) VideoFilter(c context.Context, msg, ip string) (resData *archive.FilterData, hit []string, err error) {
	params := url.Values{}
	params.Set("area", "video_submit")
	params.Set("msg", msg)
	var res struct {
		Code    int                 `json:"code"`
		Data    *archive.FilterData `json:"data"`
		Message string              `json:"message"`
	}
	if err = d.client.Post(c, d.postFilterURI, ip, params, &res); err != nil {
		log.Error("d.client.Post uri(%s) msg(%s) ip(%s)", d.postFilterURI+"?"+params.Encode(), msg, ip, err)
		err = ecode.VideoupFilterServiceErr
		return
	}
	log.Info("VideoupFilterService url(%s)", d.postFilterURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("d.client.Post uri(%s) msg(%s) ip(%s)", d.postFilterURI+"?"+params.Encode(), msg, ip, err)
		err = ecode.VideoupFilterServiceErr
		return
	}
	hit = res.Data.Hit
	resData = res.Data
	return
}

// VideoMultiFilter 批量过滤
func (d *Dao) VideoMultiFilter(c context.Context, msgs []string, ip string) (resDatas []*archive.FilterData, hit []string, err error) {
	if len(msgs) == 0 {
		return
	}
	params := url.Values{}
	params.Set("area", "video_submit")
	for _, v := range msgs {
		params.Add("msg", v)
	}
	var res struct {
		Code    int                   `json:"code"`
		Data    []*archive.FilterData `json:"data"`
		Message string                `json:"message"`
	}
	if err = d.client.Post(c, d.postMFilterURI, ip, params, &res); err != nil {
		log.Error("d.client.Post uri(%s) msg(%v) ip(%s)", d.postMFilterURI+"?"+params.Encode(), msgs, ip, err)
		err = ecode.VideoupFilterServiceErr
		return
	}
	log.Info("VideoupFilterService url(%s)", d.postMFilterURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("d.client.Post uri(%s) msg(%v) ip(%s)", d.postMFilterURI+"?"+params.Encode(), msgs, ip, err)
		err = ecode.VideoupFilterServiceErr
		return
	}
	resDatas = res.Data
	for _, v := range resDatas {
		if len(v.Hit) != 0 {
			hit = append(hit, v.Hit...)
		}
	}
	return
}
