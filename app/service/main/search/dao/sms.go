package dao

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/library/log"
)

const _smsURL = "http://ops-mng.bilibili.co/api/sendsms"

type sms struct {
	d *Dao

	client   *http.Client
	lastTime int64
	interval int64
	params   *url.Values
}

func newSMS(d *Dao) (s *sms) {
	s = &sms{
		d:        d,
		client:   &http.Client{},
		lastTime: time.Now().Unix() - d.c.SMS.Interval, //如果不想让初始化的时候告警，把减号去掉
		interval: d.c.SMS.Interval,
		params: &url.Values{
			"phone": []string{d.c.SMS.Phone},
			"token": []string{d.c.SMS.Token},
		},
	}
	return
}

func (d *Dao) SendSMS(msg string) (err error) {
	if !d.sms.IntervalCheck() {
		log.Error("发短信太频繁啦, msg：%s", msg)
		return
	}
	if err = d.sms.Send(msg); err != nil {
		log.Error("发短信失败, msg：%s, error(%v)", msg, err)
	}
	return
}

func (sms *sms) Send(msg string) (err error) {
	var req *http.Request
	sms.params.Set("message", msg)
	if req, err = http.NewRequest("GET", _smsURL+"?"+sms.params.Encode(), nil); err != nil {
		return
	}
	req.Header.Set("x1-bilispy-timeout", strconv.FormatInt(int64(time.Duration(1)/time.Millisecond), 10))
	if _, err = sms.client.Do(req); err != nil {
		log.Error("ops-mng sendsms url(%s) error(%v)", _smsURL+"?"+sms.params.Encode(), err)
	}
	return
}

// IntervalCheck accessible or not to send msg at present time
func (sms *sms) IntervalCheck() (send bool) {
	now := time.Now().Unix()
	if (now - sms.lastTime) >= sms.interval {
		send = true
		sms.lastTime = now
	} else {
		send = false
	}
	return
}
