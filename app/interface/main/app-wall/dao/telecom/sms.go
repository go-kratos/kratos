package telecom

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

const (
	_smsSendURL = "/x/internal/sms/send"
)

// SendSMS
func (d *Dao) SendSMS(c context.Context, phone int, smsTemplate, dataJSON string) (err error) {
	params := url.Values{}
	params.Set("mobile", strconv.Itoa(phone))
	params.Set("country", "086")
	params.Set("tcode", smsTemplate)
	params.Set("tparam", dataJSON)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.smsSendURL, "", params, &res); err != nil {
		log.Error("SendSMS hots url(%s) error(%v)", d.smsSendURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("SendSMS hots url(%s) error(%v)", d.smsSendURL+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("SendSMS api response code(%v)", res)
		return
	}
	return
}

// SendTelecomSMS
func (d *Dao) SendTelecomSMS(c context.Context, phone int, smsTemplate string) (err error) {
	params := url.Values{}
	params.Set("mobile", strconv.Itoa(phone))
	params.Set("country", "086")
	params.Set("tcode", smsTemplate)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.smsSendURL, "", params, &res); err != nil {
		log.Error("SendSMS hots url(%s) model(%v) error(%v)", d.smsSendURL+"?"+params.Encode(), smsTemplate, err)
		return
	}
	if res.Code != 0 {
		log.Error("SendSMS hots url(%s) error(%v)", d.smsSendURL+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("SendSMS api response code(%v)", res)
		return
	}
	return
}
