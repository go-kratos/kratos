package bnj

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_wechatAction = "NotifyCreate"
	_wechatType   = "wechat_message"
	_wechatURL    = "http://merak.bilibili.co"
)

// SendWechat  send wechat work message.
func (d *Dao) SendWechat(c context.Context, title, msg, user string) (err error) {
	var msgBytes []byte
	params := map[string]interface{}{
		"Action":    _wechatAction,
		"SendType":  _wechatType,
		"PublicKey": d.c.Bnj2019.WxKey,
		"UserName":  user,
		"Content": map[string]string{
			"subject": title,
			"body":    title + "\n" + msg,
		},
		"TreeId":    "",
		"Signature": "1",
		"Severity":  "P5",
	}
	if msgBytes, err = json.Marshal(params); err != nil {
		return
	}
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, _wechatURL, strings.NewReader(string(msgBytes))); err != nil {
		return
	}
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	res := &struct {
		RetCode int `json:"RetCode"`
	}{}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("SendWechat d.client.Do(title:%s,msg:%s,user:%s) error(%v)", title, msg, user, err)
		return
	}
	if res.RetCode != 0 {
		err = ecode.Int(res.RetCode)
		log.Error("SendWechat d.client.Do(title:%s,msg:%s,user:%s) error(%v)", title, msg, user, err)
		return
	}
	return
}
