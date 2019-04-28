package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/library/log"
)

const (
	_codeFigureNotFound = 55001
	_figureInfoURL      = "/x/internal/figure/info"
)

func (d *Dao) figureInfoURI() string {
	return d.conf.Host.API + _figureInfoURL
}

// FigureInfo .
func (d *Dao) FigureInfo(c context.Context, mid int64) (score int32, err error) {
	var (
		res = &struct {
			Code int `json:"code"`
			Data *struct {
				Percentage int32 `json:"percentage"`
			} `json:"data"`
		}{}
		params = url.Values{}
		uri    = d.figureInfoURI()
	)
	params.Set("mid", fmt.Sprint(mid))
	if err = d.httpClient.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s,%v,%d)", uri, params, err)
		return
	}
	if res.Code != 0 && res.Code != _codeFigureNotFound {
		err = fmt.Errorf("code != 0 && code !=%d", _codeFigureNotFound)
		log.Error("d.httpClient.Get(%s,%v,%v,%d)", uri, params, err, res.Code)
		return
	}
	if res != nil && res.Data != nil {
		score = res.Data.Percentage
	}
	return
}
