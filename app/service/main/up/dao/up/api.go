package up

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_picUpInfoURI   string = "/link_draw_ex/v0/doc/check"
	_blinkUpInfoURI string = "/clip_ext/v0/video/have"
)

// Pic pic return value
type Pic struct {
	Has int `json:"has_doc"`
}

// Blink blink return value
type Blink struct {
	Has int `json:"has"`
}

// Pic get pic up info.
func (d *Dao) Pic(c context.Context, mid int64, ip string) (has int, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data Pic `json:"data"`
	}
	err = d.client.Get(c, d.picUpInfoURL, ip, params, &res)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.picUpInfoURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("Pic url(%s) error(%v)", d.picUpInfoURL+"?"+params.Encode(), err)
		err = ecode.Int(res.Code)
		return
	}
	has = res.Data.Has
	return
}

// Blink get BLink up info.
func (d *Dao) Blink(c context.Context, mid int64, ip string) (has int, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int   `json:"code"`
		Data Blink `json:"data"`
	}
	err = d.client.Get(c, d.blinkUpInfoURL, ip, params, &res)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", d.blinkUpInfoURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("Blink url(%s) error(%v)", d.blinkUpInfoURL+"?"+params.Encode(), err)
		err = ecode.Int(res.Code)
		return
	}
	has = res.Data.Has
	return
}

// IsAuthor checks that whether user has permission to write article.
func (d *Dao) IsAuthor(c context.Context, mid int64, ip string) (isArt int, err error) {
	var (
		arg = &model.ArgMid{
			Mid:    mid,
			RealIP: ip,
		}
		res bool
	)
	if res, err = d.art.IsAuthor(c, arg); err != nil {
		if _, er := strconv.ParseInt(err.Error(), 10, 64); er != nil {
			log.Error("d.art.IsAuthor (%v) error(%v)", arg, err)
			err = ecode.CreativeArticleRPCErr
		}
		if ecode.Cause(err) == ecode.ArtCreationNoPrivilege {
			log.Error("d.art.IsAuthor(%d) error(%v)", mid, err)
		}
		return
	}
	if res {
		isArt = 1
	}
	return
}
