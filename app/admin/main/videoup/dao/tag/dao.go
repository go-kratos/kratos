package tag

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/admin/main/videoup/conf"
	tagr "go-common/app/interface/main/tag/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
)

// Dao is message dao.
type Dao struct {
	c            *conf.Config
	client       *xhttp.Client
	uri          string
	upBindURL    string
	adminBindURL string
	tagRPC       *tagr.Service
}

var (
	d *Dao
)

const (
	_upBindURI    = "/x/internal/tag/archive/upbind"
	_adminBindURI = "/x/internal/tag/archive/adminbind"
)

// New new a message dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		client:       xhttp.NewClient(c.HTTPClient.Write),
		tagRPC:       tagr.New2(c.TagDisRPC),
		adminBindURL: c.Host.API + _adminBindURI,
		upBindURL:    c.Host.API + _upBindURI,
	}
	return
}

// AdminBind update bind tag.
func (d *Dao) AdminBind(c context.Context, aid, mid int64, tags, regionName, ip string) (err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("tnames", tags)
	params.Set("region_name", regionName)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.adminBindURL, ip, params, &res); err != nil {
		log.Error("AdminBind d.client.Post(%s) error(%v)", d.adminBindURL+"?"+params.Encode(), err)
		return
	}
	log.Info("AdminBind url(%s) code(%d)", d.adminBindURL+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("AdminBind url(%s) code(%d)", d.adminBindURL+"?"+params.Encode(), res.Code)
		err = ecode.TagUpdateFail
	}
	return
}

// UpBind update bind tag.
func (d *Dao) UpBind(c context.Context, aid, mid int64, tags, regionName, ip string) (err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("tnames", tags)
	params.Set("region_name", regionName)
	var res struct {
		Code int `json:"code"`
	}

	if err = d.client.Post(c, d.upBindURL, ip, params, &res); err != nil {
		log.Error("UpBind d.client.Post(%s) error(%v)", d.upBindURL+"?"+params.Encode(), err)
		return
	}
	log.Info("UpBind url(%s) code(%d)", d.upBindURL+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("UpBind url(%s) code(%d)", d.upBindURL+"?"+params.Encode(), res.Code)
		err = ecode.TagUpdateFail
	}
	return
}
