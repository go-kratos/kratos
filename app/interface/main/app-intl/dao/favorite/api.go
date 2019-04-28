package favorite

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// IsFavDefault faorite count
func (d *Dao) IsFavDefault(c context.Context, mid, aid int64) (is bool, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Default bool `json:"default"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.isFavDef, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.isFavDef+"?"+params.Encode())
		return
	}
	is = res.Data.Default
	return
}

// IsFav is
func (d *Dao) IsFav(c context.Context, mid, aid int64) (is bool, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			Favorite bool `json:"favoured"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.isFav, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.isFav+"?"+params.Encode())
		return
	}
	is = res.Data.Favorite
	return
}

// AddFav add fav video
func (d *Dao) AddFav(c context.Context, mid, aid int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.addFav, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.isFav+"?"+params.Encode())
		return
	}
	return
}
