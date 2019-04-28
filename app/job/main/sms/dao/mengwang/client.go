package mengwang

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/job/main/sms/conf"
	"go-common/app/job/main/sms/model"
	smsmdl "go-common/app/service/main/sms/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Client .
type Client struct {
	c      conf.Provider
	client *bm.Client
}

type response struct {
	Result int   `json:"result"`
	MsgID  int64 `json:"msgid"`
}

// GetPid gets MengWang type ID.
func (v *Client) GetPid() int32 {
	return smsmdl.ProviderMengWang
}

// NewClient new MengWang.
func NewClient(c *conf.Config) *Client {
	return &Client{
		c:      *c.Provider,
		client: bm.NewClient(c.HTTPClient),
	}
}

func (v *Client) post(ctx context.Context, url string, params interface{}, res *response) (err error) {
	var bs []byte
	if bs, err = json.Marshal(params); err != nil {
		log.Error("json.Marshal param(%v) error(%v)", params, err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(bs)))
	if err != nil {
		log.Error("http.NewRequest(%s) param(%v) error(%v)", url, params, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if err = v.client.Do(ctx, req, &res); err != nil {
		log.Error("client.Do(%s,%v) error(%v)", url, params, err)
		return
	}
	log.Info("url(%s) body(%v) resp(%+v)", url, params, res)
	return
}

// SendSms sends MengWang sms.
func (v *Client) SendSms(ctx context.Context, r *smsmdl.ModelSend) (msgid string, err error) {
	params := make(map[string]string)
	pwd, ts := genPwd(v.c.MengWangSmsUser, v.c.MengWangSmsPwd)
	params["userid"] = v.c.MengWangSmsUser
	params["pwd"] = pwd
	params["timestamp"] = ts
	params["mobile"] = r.Mobile
	params["content"] = url.QueryEscape(r.Content)
	res := new(response)
	if err = v.post(ctx, v.c.MengWangSmsURL, params, res); err != nil {
		log.Error("mengwang SendSms param(%v) error(%v)", params, err)
		return
	}
	if res.Result != 0 {
		err = fmt.Errorf("mengwang SendSms param(%v) error(%v)", params, res.Result)
		return
	}
	msgid = strconv.FormatInt(res.MsgID, 10)
	return
}

// SendActSms sends MengWang act sms.
func (v *Client) SendActSms(ctx context.Context, r *smsmdl.ModelSend) (msgid string, err error) {
	params := make(map[string]string)
	pwd, ts := genPwd(v.c.MengWangActUser, v.c.MengWangActPwd)
	r.Content = r.Content + model.SmsSuffix
	params["userid"] = v.c.MengWangActUser
	params["pwd"] = pwd
	params["timestamp"] = ts
	params["mobile"] = r.Mobile
	params["content"] = url.QueryEscape(r.Content)
	res := new(response)
	if err = v.post(ctx, v.c.MengWangActURL, params, res); err != nil {
		log.Error("mengwang SendActSms param(%v) error(%v)", params, err)
		return
	}
	if res.Result != 0 {
		err = fmt.Errorf("mengwang SendActSms param(%v) error(%v)", params, res.Result)
		return
	}
	msgid = strconv.FormatInt(res.MsgID, 10)
	return
}

// SendBatchActSms sends multi MengWang act sms.
func (v *Client) SendBatchActSms(ctx context.Context, r *smsmdl.ModelSend) (msgid string, err error) {
	params := make(map[string]string)
	pwd, ts := genPwd(v.c.MengWangActUser, v.c.MengWangActPwd)
	params["userid"] = v.c.MengWangActUser
	params["pwd"] = pwd
	params["timestamp"] = ts
	params["mobile"] = r.Mobile
	params["content"] = url.QueryEscape(r.Content + model.SmsSuffix)
	res := new(response)
	if err = v.post(ctx, v.c.MengWangBatchURL, params, res); err != nil {
		log.Error("mengwang SendBatchActSms param(%v) error(%v)", params, err)
		return
	}
	if res.Result != 0 {
		err = fmt.Errorf("mengwang SendBatchActSms param(%v) error(%v)", params, res.Result)
		return
	}
	msgid = strconv.FormatInt(res.MsgID, 10)
	return
}

// SendInternationalSms sends MengWang international sms.
func (v *Client) SendInternationalSms(ctx context.Context, r *smsmdl.ModelSend) (msgid string, err error) {
	params := make(map[string]string)
	pwd, ts := genPwd(v.c.MengWangInternationUser, v.c.MengWangInternationPwd)
	params["userid"] = v.c.MengWangInternationUser
	params["pwd"] = pwd
	params["timestamp"] = ts
	params["mobile"] = "00" + r.Country + r.Mobile
	params["content"] = url.QueryEscape(r.Content)
	res := new(response)
	if err = v.post(ctx, v.c.MengWangInternationURL, params, res); err != nil {
		log.Error("mengwang SendInternationalSms param(%v) error(%v)", params, err)
		return
	}
	if res.Result != 0 {
		err = fmt.Errorf("mengwang SendInternationalSms param(%v) error(%v)", params, res.Result)
		return
	}
	msgid = strconv.FormatInt(res.MsgID, 10)
	return
}

// Callback .
type Callback struct {
	MsgID      int64  `json:"msgid"`
	Num        int    `json:"pknum"`
	Total      int    `json:"pktotal"`
	Mobile     string `json:"mobile"`
	SendTime   string `json:"stime"`
	ReportTime string `json:"rtime"`
	Status     string `json:"errcode"`
	Desc       string `json:"errdesc"`
}

type callbackResponse struct {
	Result    int         `json:"result"`
	Callbacks []*Callback `json:"rpts"`
}

// Callback .
func (v *Client) Callback(ctx context.Context, user, pwd, url string, count int) (callbacks []*Callback, err error) {
	pwd, ts := genPwd(user, pwd)
	params := make(map[string]string)
	params["userid"] = user
	params["pwd"] = pwd
	params["timestamp"] = ts
	params["retsize"] = strconv.Itoa(count)
	bs, err := json.Marshal(params)
	if err != nil {
		log.Error("json.Marshal param(%v) error(%v)", params, err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(bs)))
	if err != nil {
		log.Error("http.NewRequest(%s) param(%v) error(%v)", url, params, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	res := new(callbackResponse)
	if err = v.client.Do(ctx, req, &res); err != nil {
		log.Error("client.Do(%s,%v) error(%v)", url, params, err)
		return
	}
	if res.Result != 0 {
		err = fmt.Errorf("mengwang callback param(%v) res(%+v)", params, res)
		return
	}
	callbacks = res.Callbacks
	return
}

func genPwd(user, pwd string) (spwd, ts string) {
	ft := time.Now().Format("0102150405")
	str := fmt.Sprintf("%s%s%s%s", user, "00000000", pwd, ft)
	mh := md5.Sum([]byte(str))
	return hex.EncodeToString(mh[:]), ft
}
