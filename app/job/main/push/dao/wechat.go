package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"go-common/library/log"
)

type wechatResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

const (
	_url = "http://bap.bilibili.co/api/v1/message/add"
)

// SendWechat 发送企业微信消息
func (d *Dao) SendWechat(content string) (err error) {
	params := map[string]string{
		"content":   content,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"token":     d.c.Wechat.Token,
		"type":      "wechat",
		"username":  d.c.Wechat.Username,
		"url":       "",
	}
	params["signature"] = d.sign(params)
	b, err := json.Marshal(params)
	if err != nil {
		log.Error("SendWechat json.Marshal error(%v)", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, _url, bytes.NewReader(b))
	if err != nil {
		log.Error("SendWechat NewRequest error(%v), params(%s)", err, string(b))
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res := wechatResp{}
	if err = d.httpClient.Do(context.TODO(), req, &res); err != nil {
		log.Error("SendWechat Do error(%v), params(%s)", err, string(b))
		return
	}
	if res.Status != 0 {
		err = fmt.Errorf("status(%d) msg(%s)", res.Status, res.Msg)
		log.Error("SendWechat response error(%v), params(%s)", err, string(b))
	}
	return
}

func (d *Dao) sign(params map[string]string) string {
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buf := bytes.Buffer{}
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k) + "=")
		buf.WriteString(url.QueryEscape(params[k]))
	}
	h := md5.New()
	io.WriteString(h, buf.String()+d.c.Wechat.Secret)
	return fmt.Sprintf("%x", h.Sum(nil))
}
