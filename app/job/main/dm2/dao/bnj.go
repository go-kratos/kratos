package dao

import (
	"context"

	"go-common/app/job/main/dm2/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_bnjLiveConfig = "/activity/v0/bainian/config"
)

func (d *Dao) bnjConfigURI() string {
	return d.conf.Host.APILive + _bnjLiveConfig
}

// BnjConfig .
func (d *Dao) BnjConfig(c context.Context) (bnjConfig *model.BnjLiveConfig, err error) {
	var (
		res struct {
			Code    int64                `json:"code"`
			Message string               `json:"message"`
			Data    *model.BnjLiveConfig `json:"data"`
		}
	)
	if err = d.httpCli.Get(c, d.bnjConfigURI(), metadata.String(c, metadata.RemoteIP), nil, &res); err != nil {
		log.Error("bnjLiveConfig BnjConfig(url:%v) error(%v)", d.bnjConfigURI(), err)
		return
	}
	bnjConfig = res.Data
	return
}
