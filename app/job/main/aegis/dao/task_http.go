package dao

import (
	"context"
	"net/url"

	"go-common/app/job/main/aegis/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_addURL    = "/x/internal/aegis/add"
	_updateURL = "/x/internal/aegis/update"
	_cancelURL = "/x/internal/aegis/cancel"
)

//RscAdd resource add
func (d *Dao) RscAdd(c context.Context, opt *model.AddOption) error {
	uri := d.c.Host.API + _addURL
	params := opt.ToQueryURI()
	return d.commonPost(c, uri, params)
}

//RscUpdate resource update
func (d *Dao) RscUpdate(c context.Context, opt *model.UpdateOption) error {
	uri := d.c.Host.API + _updateURL
	params := opt.ToQueryURI()
	return d.commonPost(c, uri, params)
}

//RscCancel resource cancel
func (d *Dao) RscCancel(c context.Context, opt *model.CancelOption) error {
	uri := d.c.Host.API + _cancelURL
	params := opt.ToQueryURI()
	return d.commonPost(c, uri, params)
}

func (d *Dao) commonPost(c context.Context, uri string, params url.Values) error {
	res := new(model.BaseResponse)
	if err := d.httpFast.Post(c, uri, "", params, res); err != nil {
		log.Error("d.httpFast.Post(%s) params(%s) error(%v)", uri, params.Encode(), err)
		return err
	}

	if res.Code != 0 {
		log.Error("d.httpFast.Post(%s) params(%s) res(%+v)", uri, params.Encode(), res)
		return ecode.Code(res.Code)
	}
	return nil
}
