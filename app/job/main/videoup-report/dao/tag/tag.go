package tag

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

const (
	_upBindURI    = "/x/internal/tag/archive/upbind"
	_adminBindURI = "/x/internal/tag/archive/adminbind"
)

// UpBind update bind tag.
func (d *Dao) UpBind(c context.Context, mid, aid int64, tags, regionName, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
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
		err = fmt.Errorf("UpBind response code(%d)!=0", res.Code)
	}
	return
}

// AdminBind update bind tag.
func (d *Dao) AdminBind(c context.Context, mid, aid int64, tags, regionName, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
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
		err = fmt.Errorf("AdminBind response code(%d)!=0", res.Code)
	}
	return
}
