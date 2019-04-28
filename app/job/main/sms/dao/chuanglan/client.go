package chuanglan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/job/main/sms/conf"
	"go-common/app/job/main/sms/model"
	smsmdl "go-common/app/service/main/sms/model"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
)

// Client .
type Client struct {
	conf   conf.Provider
	client *xhttp.Client
}

type response struct {
	Code   string `json:"code"`
	MsgID  string `json:"msgId"`
	ErrMsg string `json:"errorMsg"`
	Time   string `json:"time"`
}

// GetPid get pid
func (v *Client) GetPid() int32 {
	return smsmdl.ProviderChuangLan
}

// NewClient new ChuangLan
func NewClient(c *conf.Config) *Client {
	return &Client{
		conf:   *c.Provider,
		client: xhttp.NewClient(c.HTTPClient),
	}
}

// SendSms send sms
func (v *Client) SendSms(ctx context.Context, r *smsmdl.ModelSend) (msgid string, err error) {
	msg := model.SmsPrefix + r.Content
	params := make(map[string]interface{})
	params["account"] = v.conf.ChuangLanSmsUser
	params["password"] = v.conf.ChuangLanSmsPwd
	params["phone"] = r.Mobile
	params["msg"] = url.QueryEscape(msg)
	params["report"] = "true"
	uri := v.conf.ChuangLanSmsURL
	msgid, err = v.post(ctx, uri, params)
	return
}

// SendActSms send act sms
func (v *Client) SendActSms(ctx context.Context, r *smsmdl.ModelSend) (msgid string, err error) {
	msg := model.SmsPrefix + r.Content + model.SmsSuffixChuangLan
	params := make(map[string]interface{})
	params["account"] = v.conf.ChuangLanActUser
	params["password"] = v.conf.ChuangLanActPwd
	params["phone"] = r.Mobile
	params["msg"] = url.QueryEscape(msg)
	params["report"] = "true"
	uri := v.conf.ChuangLanActURL
	msgid, err = v.post(ctx, uri, params)
	return
}

// SendBatchActSms send batch act sms
func (v *Client) SendBatchActSms(ctx context.Context, r *smsmdl.ModelSend) (msgid string, err error) {
	msgid, err = v.SendActSms(ctx, r)
	return
}

// SendInternationalSms send international sms
func (v *Client) SendInternationalSms(ctx context.Context, r *smsmdl.ModelSend) (msgid string, err error) {
	msg := model.SmsPrefix + r.Content
	params := make(map[string]interface{})
	params["account"] = v.conf.ChuangLanInternationUser
	params["password"] = v.conf.ChuangLanInternationPwd
	params["mobile"] = r.Country + r.Mobile
	params["msg"] = msg
	uri := v.conf.ChuangLanInternationURL
	bytesData, err := json.Marshal(params)
	if err != nil {
		log.Error("ChuangLan send international Marshal error(%v)", err)
		return
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest(http.MethodPost, uri, reader)
	if err != nil {
		log.Error("ChuangLan send international NewRequest err(%v)", err)
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	type internation struct {
		Code  string `json:"code"`
		Msgid string `json:"msgid"`
		Error string `json:"error"`
	}
	res := &internation{}
	if err = v.client.Do(ctx, request, res); err != nil {
		log.Error("ChuangLan send international client.Do err(%v)", err)
		return
	}
	if res.Code != "0" {
		err = fmt.Errorf("ChuangLan send international sms code(%v) err(%v)", res.Code, res.Error)
		return
	}
	msgid = res.Msgid
	return
}

func (v *Client) post(ctx context.Context, uri string, params map[string]interface{}) (msgid string, err error) {
	bytesData, err := json.Marshal(params)
	if err != nil {
		log.Error("ChuangLan Marshal error(%v)", err)
		return
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest(http.MethodPost, uri, reader)
	if err != nil {
		log.Error("ChuangLan NewRequest err(%v)", err)
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	res := &response{}
	if err = v.client.Do(ctx, request, res); err != nil {
		log.Error("ChuangLan client.Do err(%v)", err)
		return
	}
	if res.Code != "0" {
		err = fmt.Errorf("ChuangLan send sms code(%v) err(%v)", res.Code, res.ErrMsg)
		return
	}
	msgid = res.MsgID
	log.Info("url(%s) body(%v) resp(%+v)", uri, params, res)
	return
}

// Callback .
type Callback struct {
	MsgID      string `json:"msgId"`
	Mobile     string `json:"mobile"`
	Status     string `json:"status"`
	Desc       string `json:"statusDesc"`
	NotifyTime string `json:"notifyTime"`
	ReportTime string `json:"reportTime"`
	Length     string `json:"length"`
}

type callbackResponse struct {
	Code   int         `json:"ret"`
	Result []*Callback `json:"result"`
}

// Callback sms callbacks.
func (v *Client) Callback(ctx context.Context, account, pwd, url string, count int) (callbacks []*Callback, err error) {
	params := make(map[string]interface{})
	params["account"] = account
	params["password"] = pwd
	params["count"] = strconv.Itoa(count)
	bs, err := json.Marshal(params)
	if err != nil {
		log.Error("ChuangLan sms callback Marshal error(%v)", err)
		return
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bs))
	if err != nil {
		log.Error("ChuangLan sms callback NewRequest err(%v)", err)
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	res := &callbackResponse{}
	if err = v.client.Do(ctx, request, res); err != nil {
		log.Error("ChuangLan sms callback client.Do err(%v)", err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("ChuangLan sms callback code(%d)", res.Code)
		return
	}
	callbacks = res.Result
	return
}

// CallbackInternational sms callbacks.
func (v *Client) CallbackInternational(ctx context.Context, count int) (callbacks []*Callback, err error) {
	params := make(map[string]interface{})
	params["account"] = v.conf.ChuangLanInternationUser
	params["password"] = v.conf.ChuangLanInternationPwd
	params["count"] = strconv.Itoa(count)
	bs, err := json.Marshal(params)
	if err != nil {
		log.Error("ChuangLan international sms callback Marshal error(%v)", err)
		return
	}
	request, err := http.NewRequest(http.MethodPost, v.conf.ChuangLanInternationalCallbackURL, bytes.NewReader(bs))
	if err != nil {
		log.Error("ChuangLan international sms callback NewRequest err(%v)", err)
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	type intCallbackResponse struct {
		Code   int         `json:"code"`
		Error  string      `json:"error"`
		Result []*Callback `json:"result"`
	}
	res := &intCallbackResponse{}
	if err = v.client.Do(ctx, request, res); err != nil {
		log.Error("ChuangLan international sms callback client.Do err(%v)", err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("ChuangLan international sms callback code(%d)", res.Code)
		return
	}
	callbacks = res.Result
	return
}
