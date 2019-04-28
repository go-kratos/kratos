package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_retryAutoRenew = "/x/internal/autorenew/retry"
)

//AutoRenewPay auto renew pay.
func (d *Dao) AutoRenewPay(c context.Context, mid int64) (res *model.CommonResq, err error) {
	res = new(model.CommonResq)
	val := url.Values{}
	val.Add("mid", fmt.Sprintf("%d", mid))
	url := d.c.VipURI + _retryAutoRenew
	if err = d.client.Post(c, url, "", val, res); err != nil {
		log.Error("reques fail url %v params:%+v result:%+v, err:%+v", url, val, res, err)
		err = errors.WithStack(err)
		return
	}
	if res.Code != 0 {
		log.Error("reques fail url %v params:%+v result:%+v, err:%+v", url, val, res, err)
		return
	}
	log.Info("reques success url %v params:%+v result:%+v", url, val, res)
	return
}
