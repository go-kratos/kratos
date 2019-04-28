package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/log"
)

// http client

const (
	smsURL = "http://api.bilibili.co/x/internal/sms/send"
)

// SendCapture is
func (d *Dao) SendCapture(c context.Context, mid int64, code int) (err error) {
	var (
		params = url.Values{}
	)
	params.Set("mid", fmt.Sprintf("%d", mid))
	params.Set("tcode", "acc_01")
	params.Set("tparam", fmt.Sprintf(`{"identify_code":"%d"}`, code))

	var resp struct {
		Code int `json:"code"`
	}
	for i := 0; i < 3; i++ {
		err = d.client.Post(c, smsURL, "", params, &resp)
		if err != nil || resp.Code != 0 {
			log.Error("d.client.Post(%s,%+v) resp.Code(%d)", smsURL, params, resp.Code)
			err = ecode.RealnameCaptureErr
		} else {
			break
		}
	}
	return
}
