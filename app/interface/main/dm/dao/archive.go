package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/dm/model"
	"go-common/library/log"
)

const (
	_cidInfo = "/videoup/cid"
)

func (d *Dao) cidInfoURI() string {
	return d.conf.Host.Archive + _cidInfo
}

// CidInfo 获取cid详细信息
func (d *Dao) CidInfo(c context.Context, cid int64) (info *model.CidInfo, err error) {
	var (
		res = &struct {
			Code int            `json:"code"`
			Data *model.CidInfo `json:"data"`
		}{}
		params = url.Values{}
		uri    = d.cidInfoURI()
	)
	params.Set("cid", fmt.Sprint(cid))
	if err = d.httpClient.Get(c, uri, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s,%v,%d)", uri, params, err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("code is:%d", res.Code)
		log.Error("d.httpClient.Get(%s,%v,%v,%d)", uri, params, err, res.Code)
		return
	}
	return res.Data, nil
}
