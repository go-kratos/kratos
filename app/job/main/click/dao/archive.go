package dao

import (
	"context"

	"go-common/library/log"
)

const (
	_maxAIDPath = "http://api.bilibili.co/x/internal/v2/archive/maxAid"
)

// MaxAID return max aid
func (d *Dao) MaxAID(c context.Context) (id int64, err error) {
	var res struct {
		Code int   `json:"code"`
		Data int64 `json:"data"`
	}
	if err = d.client.Get(c, _maxAIDPath, "", nil, &res); err != nil {
		log.Error("d.client.MaxAid error(%+v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("d.client.MaxAid Code(%d)", res.Code)
		return
	}
	log.Info("got MaxAid(%d)", res.Data)
	id = res.Data
	return
}
