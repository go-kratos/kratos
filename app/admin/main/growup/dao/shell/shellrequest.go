package shell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/admin/main/growup/conf"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"net/http"
)

//SignInterface set signature
type SignInterface interface {
	SetSign(sign string)
	SetCustomerID(customerID string)
	SetSignType(signType string)
}

//OrderInfo order info
type OrderInfo struct {
	Brokerage    string `json:"brokerage"`
	Mid          int64  `json:"mid"`
	ThirdCoin    string `json:"thirdCoin"`
	ThirdCtime   string `json:"thirdCtime"`
	ThirdOrderNo string `json:"thirdOrderNo"`
}

//OrderRequest shell order request
type OrderRequest struct {
	CustomerID  string      `json:"customerId"`
	ProductName string      `json:"productName"`
	Data        []OrderInfo `json:"data"`
	NotifyURL   string      `json:"notifyUrl"`
	Rate        string      `json:"rate"`
	SignType    string      `json:"signType"`
	Timestamp   string      `json:"timestamp"`
	Sign        string      `json:"sign,omitempty"`
}

//SetSign set sign
func (o *OrderRequest) SetSign(sign string) {
	o.Sign = sign
}

//SetCustomerID set customId
func (o *OrderRequest) SetCustomerID(customerID string) {
	o.CustomerID = customerID
}

//SetSignType set signtype
func (o *OrderRequest) SetSignType(signType string) {
	o.SignType = signType
}

//OrderResponse shell order response
type OrderResponse struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
}

const (
	//CallbackStatusCreate 创建中状态
	CallbackStatusCreate = "CREATE"
	//CallbackStatusSuccess 成功
	CallbackStatusSuccess = "SUCCESS"
	//CallbackStatusFail 失败
	CallbackStatusFail = "FAIL"
)

//OrderCallbackJSON MsgContent in OrderCallbackParam
type OrderCallbackJSON struct {
	CustomerID   string          `json:"customerId"`
	Status       string          `json:"status"`
	ThirdOrderNo string          `json:"thirdOrderNo"`
	Mid          string          `json:"mid"`
	Ext          json.RawMessage `json:"ext"`
	Timestamp    string          `json:"timestamp"`
	SignType     string          `json:"signType"`
	Sign         string          `json:"sign"`
}

//IsSuccess success
func (o *OrderCallbackJSON) IsSuccess() bool {
	return o.Status == CallbackStatusSuccess
}

//IsFail fail
func (o *OrderCallbackJSON) IsFail() bool {
	return o.Status == CallbackStatusFail
}

//IsCreate creating
func (o *OrderCallbackJSON) IsCreate() bool {
	return o.Status == CallbackStatusCreate
}

//OrderCallbackParam call back url param
type OrderCallbackParam struct {
	MsgID      string `form:"msgId"`
	MsgContent string `form:"msgContent"`
}

//OrderCheckRequest request
type OrderCheckRequest struct {
	Sign          string `json:"sign,omitempty"`
	SignType      string `json:"signType"`
	CustomerID    string `json:"customerId"`
	Timestamp     int64  `json:"timestamp"`
	ThirdOrderNos string `json:"thirdOrderNos"`
}

//SetSign set sign
func (o *OrderCheckRequest) SetSign(sign string) {
	o.Sign = sign
}

//SetCustomerID set customer id
func (o *OrderCheckRequest) SetCustomerID(customerID string) {
	o.CustomerID = customerID
}

//SetSignType set sign type, always "MD5"
func (o *OrderCheckRequest) SetSignType(signType string) {
	o.SignType = signType
}

//OrderStatusData call back data
type OrderStatusData struct {
	ThirdOrderNo string `json:"thirdOrderNo"`
	Status       string `json:"status"`
	Mid          string `json:"mid"`
}

//IsSuccess is successful
func (o *OrderStatusData) IsSuccess() bool {
	return o.Status == CallbackStatusSuccess
}

//OrderCheckResponse response
type OrderCheckResponse struct {
	Errno  int               `json:"errno"`
	Msg    string            `json:"msg"`
	Orders []OrderStatusData `json:"data"`
}

//Client shell client
type Client struct {
	conf       conf.ShellConfig
	CustomID   string
	Token      string
	HTTPClient *blademaster.Client
	isDebug    bool
}

//New client
func New(conf *conf.ShellConfig, httpClient *blademaster.Client) *Client {
	return &Client{
		CustomID:   conf.CustomID,
		Token:      conf.Token,
		HTTPClient: httpClient,
		conf:       *conf,
	}
}

//SetDebug set debug
func (s *Client) SetDebug(isDebug bool) {
	s.isDebug = isDebug
}

//SendOrderRequest send order rquest
func (s *Client) SendOrderRequest(ctx context.Context, req *OrderRequest) (res *OrderResponse, err error) {
	var host = s.conf.PayHost
	if host == "" {
		host = "pay.bilibili.co"
	}
	var url = "http://" + host + "/bk-int/brokerage/rechargeBrokerage"
	res = &OrderResponse{}
	err = s.SendShellRequest(ctx, url, req, res)
	if err != nil {
		log.Error("send order request fail, err=%s", err)
	}
	return
}

//SendCheckOrderRequest send check order request
func (s *Client) SendCheckOrderRequest(ctx context.Context, req *OrderCheckRequest) (res *OrderCheckResponse, err error) {
	var host = s.conf.PayHost
	if host == "" {
		host = "pay.bilibili.co"
	}
	var url = "http://" + host + "/bk-int/brokerage/queryRechargeBrokerage"
	res = &OrderCheckResponse{}
	err = s.SendShellRequest(ctx, url, req, res)
	if err != nil {
		log.Error("send check order request fail, err=%s", err)
	}
	return
}

//SendShellRequest send request
func (s *Client) SendShellRequest(ctx context.Context, url string, req interface{}, res interface{}) (err error) {
	r, ok := req.(SignInterface)
	if !ok {
		err = fmt.Errorf("cast fail, req is not SignInterface")
		return
	}
	r.SetSignType("MD5")
	r.SetCustomerID(s.CustomID)
	sign, err := Sign(r, s.Token)
	if err != nil {
		return
	}
	r.SetSign(sign)
	jsonStr, err := json.Marshal(r)
	if s.isDebug {
		log.Info("send request, url=%s, req=%s", url, jsonStr)
	}
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	var buffer = bytes.NewBuffer(jsonStr)

	httpreq, _ := http.NewRequest("POST", url, buffer)
	httpreq.Header.Set("Content-Type", "application/json")
	bs, err := s.HTTPClient.Raw(ctx, httpreq)
	if s.isDebug {
		log.Info("req=%s, response=%s", jsonStr, bs)
	}
	if err != nil {
		log.Error("get response err, err=%s, response=%s", err, string(bs))
		return
	}
	if err = json.Unmarshal(bs, res); err != nil {
		log.Error("json decode err, response=%s", string(bs))
	}
	return
}
