package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_typesURL = "/videoup/types"
)

func (d *Dao) typesURI() string {
	return d.conf.Host.Archive + _typesURL
}

// TypeMapping is second types opposite first types.
func (d *Dao) TypeMapping(c context.Context) (rmap map[int16]int16, err error) {
	var res struct {
		Code    int                     `json:"code"`
		Message string                  `json:"message"`
		Data    map[int16]*archive.Type `json:"data"`
	}
	if err = d.httpCli.Get(c, d.typesURI(), "", nil, &res); err != nil {
		log.Error("d.httpCli.Get() error(%v) typesURI(%s)", err, d.typesURI())
		return
	}
	if res.Code != ecode.OK.Code() {
		err = fmt.Errorf("bangumi seasons api failed(%d)", res.Code)
		log.Error("url(%s) res code(%d)", d.typesURI(), res.Code)
		return
	}
	rmap = make(map[int16]int16, len(res.Data))
	for _, v := range res.Data {
		if v.PID != 0 {
			rmap[v.ID] = v.PID
		}
	}
	return
}
