package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/log"
)

const (
	_doubleCheckURL = "http://passport.bilibili.co/intranet/acc/security/mid"
)

// DoubleCheck notify passport to remove login.
func (d *Dao) DoubleCheck(c context.Context, mid int64) (err error) {
	params := url.Values{}
	params.Set("mids", strconv.FormatInt(mid, 10))
	params.Set("desc", "异地风险，系统导入")
	params.Set("operator", "异地系统判断")
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, _doubleCheckURL, "", params, &res); err != nil {
		log.Error("d.Doublecheck err(%v)", err)
	}
	log.Info("d.DoubleCheck mid %d ", mid)
	return
}
