package dao

import (
	"context"
	"fmt"

	"go-common/library/ecode"
	"go-common/library/log"
	"net/url"
)

const (
	_frozenChange = "/internal/v1/user/frozenChange"
)

// OldFrozenChange .
func (d *Dao) OldFrozenChange(mid, status int64) (err error) {
	var (
		c = context.Background()
	)
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%v", mid))
	params.Set("status", fmt.Sprintf("%v", status))
	rel := new(struct {
		Code int64  `json:"code"`
		Data string `json:"data"`
	})
	if err = d.client.Get(c, d.c.URLConf.OldVipCoURL+_frozenChange, "127.0.0.1", params, rel); err != nil {
		log.Error("send error(%v) url(%v)", err, d.c.URLConf.OldVipCoURL+_frozenChange)
		return
	}
	if rel != nil && rel.Code == int64(ecode.OK.Code()) {
		return
	}
	err = ecode.VipJavaAPIErr
	return
}
