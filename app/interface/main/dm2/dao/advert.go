package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_adURL = "/bce/api/bce/wise"
)

func (d *Dao) adURI() string {
	return d.conf.Host.Advert + _adURL
}

// DMAdvert dm advert.
func (d *Dao) DMAdvert(c context.Context, aid, cid, mid, build int64, buvid, mobiApp, adExtra string) (data *model.AD, err error) {
	var (
		res *struct {
			Code int       `json:"code"`
			Data *model.AD `json:"data"`
		}
	)
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("aid", fmt.Sprint(aid))
	params.Set("cid", fmt.Sprint(cid))
	params.Set("buvid", buvid)
	params.Set("resource", model.Resource(mobiApp))
	params.Set("mobi_app", mobiApp)
	params.Set("build", fmt.Sprint(build))
	params.Set("ip", ip)
	params.Set("ad_extra", adExtra)
	if mid != 0 {
		params.Set("mid", fmt.Sprint(mid))
	}
	if err = d.httpCli.Get(c, d.adURI(), ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.adURI()+"?"+params.Encode())
		return
	}
	data = res.Data
	return
}
