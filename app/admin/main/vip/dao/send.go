package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	_SendUserNotify = "/api/notify/send.user.notify.do"

	_payRefund = "/payplatform/refund/request"

	_minRead = 1024 * 64
)

// SendMultipMsg send multip msg
func (d *Dao) SendMultipMsg(c context.Context, mids, content, title, mc, ip string, dataType int) (err error) {

	params := url.Values{}
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("context", content)
	params.Set("data_type", strconv.FormatInt(int64(dataType), 10))
	params.Set("mid_list", mids)
	if err = d.client.Post(c, d.c.Property.MsgURI+_SendUserNotify, "127.0.0.1", params, nil); err != nil {
		log.Error("SendMultipMsg error(%v)", err)
		return
	}
	return
}

//PayRefund .
func (d *Dao) PayRefund(c context.Context, arg *model.PayOrder, refundAmount float64, refundID string) (err error) {
	params := make(map[string]string)
	params["customerId"] = strconv.FormatInt(d.c.PayConf.CustomerID, 10)
	params["notifyUrl"] = d.c.PayConf.RefundURL
	params["version"] = d.c.PayConf.Version
	params["signType"] = "MD5"
	params["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	params["traceId"] = uuid.New().String()
	params["refundDesc"] = "大会员退款"
	params["customerRefundId"] = refundID
	params["txId"] = arg.ThirdTradeNo
	params["totalAmount"] = strconv.Itoa(int(arg.Money * 100))
	params["refundAmount"] = strconv.Itoa(int(refundAmount * 100))
	sign := d.paySign(params, d.c.PayConf.Token)
	params["sign"] = sign
	resq := new(struct {
		Code int    `json:"errno"`
		Msg  string `json:"msg"`
	})
	if err = d.doPaySend(c, d.c.PayConf.BaseURL, _payRefund, "", nil, nil, params, resq); err != nil {
		err = errors.WithStack(err)
		return
	}
	if resq.Code != ecode.OK.Code() {
		err = ecode.Int(resq.Code)
	}
	return
}

func (d *Dao) paySign(params map[string]string, token string) (sign string) {
	tmp := d.sortParamsKey(params)
	var b bytes.Buffer
	b.WriteString(tmp)
	b.WriteString(fmt.Sprintf("&token=%s", token))
	log.Info("pay sign params:(%s) \n", b.String())
	mh := md5.Sum(b.Bytes())
	// query
	sign = hex.EncodeToString(mh[:])
	log.Info("pay sign (%v)", sign)
	return
}

func (d *Dao) sortParamsKey(v map[string]string) string {
	if v == nil {
		return ""
	}

	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		buf.WriteString(vs)
	}
	return buf.String()
}

func (d *Dao) doPaySend(c context.Context, basePath, path, IP string, cookie []*http.Cookie, header map[string]string, params map[string]string, data interface{}) (err error) {
	var (
		req    *http.Request
		client = new(http.Client)
		resp   *http.Response
		bs     []byte
	)
	url := basePath + path
	marshal, _ := json.Marshal(params)
	if req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(string(marshal))); err != nil {
		err = errors.WithStack(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-backend-bili-real-ip", IP)
	for _, v := range cookie {
		req.AddCookie(v)
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	if resp, err = client.Do(req); err != nil {
		log.Error("call url:%v params:(%+v)", basePath+path, params)
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, req.URL.String())
		return
	}
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = errors.Wrapf(err, "host:%s, url:%s", req.URL.Host, req.URL.String())
		return
	}
	if err = json.Unmarshal(bs, data); err != nil {
		err = errors.WithStack(err)
		return
	}
	log.Info("call url:%v params:%+v result:%+v", url, params, data)
	return
}

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
