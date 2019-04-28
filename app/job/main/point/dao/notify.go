package dao

import (
	"context"
	"go-common/library/log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

//notify status.
const (
	_notifySuccess = "1"
)

// Notify point change .
func (d *Dao) Notify(c context.Context, notifyURL string, mid int64, orderID string, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("order_id", orderID)
	params.Set("status", _notifySuccess)
	req, err := d.client.NewRequest(http.MethodPost, notifyURL, ip, params)
	if err != nil {
		err = errors.Wrapf(err, "Notify NewRequest(%s)", notifyURL+"?"+params.Encode())
		return
	}
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		err = errors.Wrapf(err, "Notify d.client.Do(%s)", notifyURL+"?"+params.Encode())
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(err, "Notify code != 0 (%s)", notifyURL+"?"+params.Encode())
		return
	}
	return
}

// NotifyCacheDel notify cache del.
func (d *Dao) NotifyCacheDel(c context.Context, notifyURL string, mid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	req, err := d.client.NewRequest(http.MethodGet, notifyURL, ip, params)
	if err != nil {
		err = errors.Wrapf(err, "NotifyCacheDel NewRequest(%s)", notifyURL+"?"+params.Encode())
		return
	}
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		err = errors.Wrapf(err, "NotifyCacheDel d.client.Do(%s)", notifyURL+"?"+params.Encode())
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(err, "NotifyCacheDel code != 0 (%s)", notifyURL+"?"+params.Encode())
		return
	}
	log.Info("notify suc(%d) url(%s)", mid, notifyURL)
	return
}
