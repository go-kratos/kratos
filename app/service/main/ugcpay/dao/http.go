package dao

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"go-common/app/service/main/ugcpay/conf"
	"go-common/app/service/main/ugcpay/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// PayRefund 调用支付平台退款接口
func (d *Dao) PayRefund(c context.Context, dataJSON string) (err error) {
	resp := new(struct {
		Code int    `json:"errno"`
		Msg  string `json:"msg"`
	})

	if err = d.paySend(c, conf.Conf.Biz.Pay.URLRefund, dataJSON, resp); err != nil {
		err = errors.WithStack(err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
	}
	return
}

// PayCancel 调用支付平台订单取消接口
func (d *Dao) PayCancel(c context.Context, dataJSON string) (err error) {
	resp := new(struct {
		Code int    `json:"errno"`
		Msg  string `json:"msg"`
	})

	if err = d.paySend(c, conf.Conf.Biz.Pay.URLCancel, dataJSON, resp); err != nil {
		err = errors.WithStack(err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
	}
	return
}

// PayQuery 调用支付平台订单查询接口 return map[orderID]*model.PayOrder
func (d *Dao) PayQuery(c context.Context, dataJSON string) (orders map[string][]*model.PayOrder, err error) {
	resp := new(struct {
		Code int             `json:"errno"`
		Msg  string          `json:"msg"`
		Data *model.PayQuery `json:"data"`
	})

	if err = d.paySend(c, conf.Conf.Biz.Pay.URLQuery, dataJSON, resp); err != nil {
		err = errors.WithStack(err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = ecode.Int(resp.Code)
		return
	}
	if resp.Data == nil {
		err = errors.Errorf("PayQuery got nil data, resp: %+v", resp)
		return
	}
	orders = make(map[string][]*model.PayOrder)
	for _, o := range resp.Data.Orders {
		orders[o.OrderID] = append(orders[o.OrderID], o)
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
	log.Info("paySend call url: %s, body: %s, resp: %s", url, jsonData, bs)
	if err = json.Unmarshal(bs, respData); err != nil {
		err = errors.WithStack(err)
		return
	}
	log.Info("paySend call url: %v, body: %s, resp: %+v", url, jsonData, respData)
	return
}
