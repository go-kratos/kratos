package http

import (
	"context"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_filterURI = "/x/internal/filter/v3/hit"
)

// Res 筛选结果
type Res struct {
	Code int64 `json:"code"`
	Data []struct {
		Level int64  `json:"level"`
		Msg   string `json:"msg"`
	} `json:"data"`
}

// FilterMulti .批量过滤
func (d *Dao) FilterMulti(c context.Context, area string, msg string) (hits []string, err error) {
	params := url.Values{}
	params.Set("area", area)
	params.Set("msg", msg)
	params.Set("level", "10")
	res := new(Res)

	log.Info("FilterMulti area(%s) msg(%s)", area, msg)
	if err = d.clientR.Post(c, d.c.Host.API+_filterURI, "", params, res); err != nil {
		log.Error("d.clientR.Get error(%v)", err)
		return
	}

	if res.Code != 0 {
		err = ecode.Code(res.Code)
		log.Error("FilterMulti res(%+v) error(%+v)", res, err)
		return
	}

	for _, dt := range res.Data {
		hits = append(hits, dt.Msg)
	}
	log.Info("FilterMulti area(%s) msg(%s) hits(%v)", area, msg, hits)
	return
}
