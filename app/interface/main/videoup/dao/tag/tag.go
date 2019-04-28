package tag

import (
	"context"
	"go-common/app/interface/main/creative/model/tag"
	"go-common/library/ecode"
	"go-common/library/log"
	"net/url"
	"strconv"
)

const (
	_upBindURI = "/x/internal/tag/archive/upbind"
	_tagCheck  = "/x/internal/tag/check"
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
	if err = d.httpW.Post(c, d.upBindURL, ip, params, &res); err != nil {
		log.Error("d.httpW.Post(%s) error(%v)", d.upBindURL+"?"+params.Encode(), err)
		err = ecode.CreativeTagErr
		return
	}
	log.Info("url(%s) code(%d)", d.upBindURL+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("url(%s) code(%d)", d.upBindURL+"?"+params.Encode(), res.Code)
		err = ecode.CreativeTagErr
	}
	return
}

// TagCheck tag check
func (d *Dao) TagCheck(c context.Context, mid int64, tagName string) (t *tag.Tag, err error) {
	var res struct {
		Code int      `json:"code"`
		Data *tag.Tag `json:"data"`
	}
	params := url.Values{}
	params.Set("tag_name", tagName)
	params.Set("mid", strconv.FormatInt(mid, 10))
	if err = d.httpW.Get(c, d.TagCheckURL, "", params, &res); err != nil {
		log.Error("TagCheck url(%s) p(%+v) response(%s) error(%v)", d.TagCheckURL, params.Encode(), res, err)
		err = ecode.CreativeTagErr
		return
	}
	log.Info("TagCheck mid(%d) url(%s) res(%v)", mid, d.TagCheckURL, res)
	if res.Code != 0 {
		log.Error("TagCheck url(%s) res(%v)", d.TagCheckURL, res)
		err = ecode.Int(res.Code)
		return
	}
	t = res.Data
	return
}
