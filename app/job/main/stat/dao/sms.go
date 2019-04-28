package dao

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
)

const (
	_maxAIDPath = "http://api.bilibili.co/x/internal/v2/archive/maxAid"
)

// MaxAID return max aid
func (d *Dao) MaxAID(c context.Context) (id int64, err error) {
	var res struct {
		Code int   `json:"code"`
		Data int64 `json:"data"`
	}
	if err = d.smsClient.Get(c, _maxAIDPath, "", nil, &res); err != nil {
		return
	}
	if res.Code != 0 {
		log.Error("d.client.MaxAid Code(%d)", res.Code)
		return
	}
	log.Info("got MaxAid(%d)", res.Data)
	id = res.Data
	return
}

// SendQiyeWX send qiye wx
func (d *Dao) SendQiyeWX(msg string) {
	type wxParams struct {
		Username  string `json:"username"`
		Content   string `json:"content"`
		Token     string `json:"token"`
		Timestamp int64  `json:"timestamp"`
		Sign      string `json:"signature"`
	}
	var resp struct {
		Status int64  `json:"status"`
		Msg    string `json:"msg"`
	}
	params := url.Values{}
	params.Set("username", d.c.Monitor.Users)
	params.Set("content", msg)
	params.Set("token", d.c.Monitor.Token)
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + d.c.Monitor.Secret))
	params.Set("signature", hex.EncodeToString(mh[:]))
	p := &wxParams{
		Username: params.Get("username"),
		Content:  params.Get("content"),
		Token:    params.Get("token"),
		Sign:     params.Get("signature"),
	}
	p.Timestamp, _ = strconv.ParseInt(params.Get("timestamp"), 10, 64)
	bs, _ := json.Marshal(p)
	payload := strings.NewReader(string(bs))
	req, _ := http.NewRequest("POST", d.c.Monitor.URL, payload)
	req.Header.Add("content-type", "application/json; charset=utf-8")
	if err := d.smsClient.Do(context.TODO(), req, &resp); err != nil {
		log.Error("d.smsClient.Do() error(%v)", err)
	}
}
