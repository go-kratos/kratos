package dao

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"go-common/app/job/main/ugcpay/conf"
	"go-common/app/job/main/ugcpay/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// PayRechargeShell .
func (d *Dao) PayRechargeShell(c context.Context, dataJSON string) (err error) {
	resp := new(struct {
		Code int    `json:"errno"`
		Msg  string `json:"msg"`
	})

	if err = d.paySend(c, conf.Conf.Biz.Pay.RechargeShellURL, dataJSON, resp); err != nil {
		err = errors.WithStack(err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
	}
	return
}

// PayCheckRefundOrder return map[txID]*model.PayCheckOrder
func (d *Dao) PayCheckRefundOrder(c context.Context, dataJSON string) (orders map[string][]*model.PayCheckRefundOrderEle, err error) {
	resp := new(struct {
		Code int                          `json:"code"`
		Msg  string                       `json:"message"`
		Data []*model.PayCheckRefundOrder `json:"data"`
	})

	if err = d.paySend(c, conf.Conf.Biz.Pay.CheckRefundOrderURL, dataJSON, resp); err != nil {
		err = errors.WithStack(err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
		return
	}
	if resp.Data == nil {
		err = errors.Errorf("PayCheckRefundOrder got nil data, url: %s, body: %s, resp: %+v", conf.Conf.Biz.Pay.CheckOrderURL, dataJSON, resp)
		return
	}
	orders = make(map[string][]*model.PayCheckRefundOrderEle)
	for _, o := range resp.Data {
		orders[o.TXID] = append(orders[o.TXID], o.Elements...)
	}
	return
}

// PayCheckOrder return map[txID]*model.PayCheckOrder
func (d *Dao) PayCheckOrder(c context.Context, dataJSON string) (orders map[string]*model.PayCheckOrder, err error) {
	resp := new(struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"message"`
		Data []*model.PayCheckOrder `json:"data"`
	})

	if err = d.paySend(c, conf.Conf.Biz.Pay.CheckOrderURL, dataJSON, resp); err != nil {
		err = errors.WithStack(err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
		return
	}
	if resp.Data == nil {
		err = errors.Errorf("PayCheckOrder got nil data, url: %s, body: %s, resp: %+v", conf.Conf.Biz.Pay.CheckOrderURL, dataJSON, resp)
		return
	}
	orders = make(map[string]*model.PayCheckOrder)
	for _, o := range resp.Data {
		orders[o.TxID] = o
	}
	return
}

func (d *Dao) paySend(c context.Context, url string, jsonData string, respData interface{}) (err error) {
	var (
		req    *http.Request
		client = new(http.Client)
		resp   *http.Response
		bs     []byte
	)
	log.Info("paySend call url: %s, body: %s", url, jsonData)
	defer func() {
		log.Info("paySend call url: %v, body: %s, resp: %+v, err: %+v", url, jsonData, respData, err)
	}()

	if req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(jsonData)); err != nil {
		err = errors.WithStack(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	if resp, err = client.Do(req); err != nil {
		err = errors.Wrapf(err, "call url: %s, body: %s", url, jsonData)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.Errorf("d.paySend incorrect http status: %d, host: %s, url: %s", resp.StatusCode, req.URL.Host, req.URL.String())
		return
	}
	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		err = errors.Wrapf(err, "d.paySend ioutil.ReadAll")
		return
	}
	if err = json.Unmarshal(bs, respData); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
