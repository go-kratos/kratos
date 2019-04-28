package dao

import (
	"context"
	"net/http"
	"net/url"

	"go-common/library/log"
)

const _smsURL = "http://ops-mng.bilibili.co/api/sendsms"

func (d *Dao) SendSMS(msg string) (err error) {
	var (
		req *http.Request
		res struct {
			Result bool `json:"result"`
		}
	)
	params := url.Values{}
	params.Set("phone", d.c.SMS.Phone)
	params.Set("message", msg)
	params.Set("token", d.c.SMS.Token)
	if req, err = d.smsClient.NewRequest("GET", _smsURL+"?"+params.Encode(), "", nil); err != nil {
		return
	}

	if err = d.smsClient.Do(context.TODO(), req, &res); err != nil {
		log.Error("ops-mng sendsms url(%s) error(%v)", _smsURL+"?"+params.Encode(), err)
		return
	}
	if !res.Result {
		log.Error("ops-mng sendsms url(%s) error(%v)", _smsURL+"?"+params.Encode(), res.Result)
	}
	return
}
