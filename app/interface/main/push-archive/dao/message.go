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

// wechatResp 企业微信的响应
type wechatResp struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

// WechatMessage 发送企业微信消息
func (d *Dao) WechatMessage(content string) (err error) {
	params := map[string]string{
		"content":   content,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"title":     "",
		"token":     d.c.Wechat.Token,
		"type":      "wechat",
		"username":  d.c.Wechat.UserName,
		"url":       "",
	}
	params["signature"] = d.signature(params, d.c.Wechat.Secret)
	b, err := json.Marshal(params)
	if err != nil {
		log.Error("WechatMessage json.Marshal error(%v)", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, "http://bap.bilibili.co/api/v1/message/add", bytes.NewReader(b))
	if err != nil {
		log.Error("WechatMessage NewRequest error(%v), params(%s)", err, string(b))
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res := wechatResp{}
	if err = d.httpClient.Do(context.TODO(), req, &res); err != nil {
		log.Error("WechatMessage Do error(%v), params(%s)", err, string(b))
		return
	}
	if res.Status != 0 {
		err = fmt.Errorf("status(%d) msg(%s)", res.Status, res.Msg)
		log.Error("WechatMessage response error(%v), params(%s)", err, string(b))
		return
	}
	return
}

// signature 加密算法
func (d *Dao) signature(params map[string]string, secret string) string {
	// content=xxx&timestamp=xxx格式
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
	// 加密
	h := md5.New()
	io.WriteString(h, buf.String()+secret)
	return fmt.Sprintf("%x", h.Sum(nil))
}
