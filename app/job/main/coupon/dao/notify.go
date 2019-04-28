package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/job/main/coupon/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_defPltform = "inner"
)

// NotifyRet notify use coupon ret.
func (d *Dao) NotifyRet(c context.Context, notifyURL string, ticketNO string, orderNO string, ip string) (data *model.CallBackRet, err error) {
	params := url.Values{}
	params.Set("ticket_no", ticketNO)
	params.Set("order_no", orderNO)
	params.Set("platform", _defPltform)
	params.Set("build", "0")
	log.Info("call service notify ret %s", notifyURL+"?"+params.Encode())
	var res struct {
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    *model.CallBackRet `json:"result"`
	}
	if err = d.client.Post(c, notifyURL, ip, params, &res); err != nil {
		err = errors.Wrapf(err, "call service(%s) error", notifyURL+"?"+params.Encode())
		return
	}
	if res.Code != 0 {
		err = errors.WithStack(fmt.Errorf("call service(%s) error, res code is not 0, resp:%v", notifyURL+"?"+params.Encode(), res))
		return
	}
	data = res.Data
	log.Info("call service notify ret suc req(%s) data(%v) ", params.Encode(), data)
	return
}
