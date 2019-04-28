package monitor

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/model"
	"go-common/library/log"
)

const (
	_uri      = "/api/v1/message/add"
	_method   = "POST"
	_fileType = "application/json"
)

// Dao is message dao.
type Dao struct {
	c      *conf.Config
	client *http.Client
	url    string
}

// New new a message dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		client: &http.Client{
			Timeout: time.Duration(time.Second * 1),
		},
		url: c.Monitor.Host + _uri,
	}
	return
}

// Send send exception message to owner.
func (d *Dao) Send(c context.Context, username, msg string) (err error) {
	params := url.Values{}
	now := time.Now().Unix()
	params.Set("username", username)
	params.Set("content", msg)
	params.Set("title", "test")
	params.Set("url", "")
	params.Set("type", "wechat")
	params.Set("token", d.c.Monitor.AppToken)
	params.Set("timestamp", strconv.FormatInt(now, 10))
	bap := &model.BAP{
		UserName:  params.Get("username"),
		Content:   params.Get("content"),
		Title:     params.Get("title"),
		URL:       params.Get("url"),
		Ty:        params.Get("type"),
		Token:     params.Get("token"),
		TimeStamp: now,
		Signature: d.getSign(params),
	}
	jsonStr, err := json.Marshal(bap)
	if err != nil {
		log.Error("monitor json.Marshal error (%v)", err)
		return
	}
	req, err := http.NewRequest(_method, d.url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error("monitor http.NewRequest error (%v)", err)
		return
	}
	req.Header.Add("Content-Type", _fileType)
	// timeout
	ctx, cancel := context.WithTimeout(c, 800*time.Millisecond)
	req = req.WithContext(ctx)
	defer cancel()
	response, err := d.client.Do(req)
	if err != nil {
		log.Error("monitor  d.client.Post error(%v)", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Error("monitor http.StatusCode nq http.StatusOK (%d) | url(%s)", response.StatusCode, d.url)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("monitor ioutil.ReadAll error(%v)", err)
		return
	}
	var result struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	}
	if err = json.Unmarshal(body, &result); err != nil {
		log.Error("monitor json.Unmarshal error(%v)", err)
	}
	if result.Status != 0 {
		log.Error("monitor get status(%d) msg(%s)", result.Status, result.Msg)
	}
	return
}

func (d *Dao) getSign(params url.Values) (sign string) {
	for k, v := range params {
		if len(v) == 0 {
			params.Del(k)
		}
	}
	h := md5.New()
	io.WriteString(h, params.Encode()+d.c.Monitor.AppSecret)
	sign = fmt.Sprintf("%x", h.Sum(nil))
	return
}
