package dao

import (
	"context"
	"errors"
	"net/url"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_notifyDataType = "4" // 消息类型：1、回复我的2、@我3、收到的爱4、业务通知
	_notifyURL      = "/api/notify/send.user.notify.do"
	_notifyMC       = "1_8_2"
)

func (d *Dao) notifyURI() string {
	return d.conf.Host.Message + _notifyURL
}

// SendNotify 发送站内信
func (d *Dao) SendNotify(c context.Context, title, content string, mids []int64) (err error) {
	res := struct {
		Code int `json:"code"`
		Data struct {
			TotalCount int
			ErrorCount int
		} `json:"data"`
	}{}
	params := url.Values{}
	params.Set("mc", _notifyMC)
	params.Set("title", title)
	params.Set("data_type", _notifyDataType)
	params.Set("context", content)
	params.Set("mid_list", xstr.JoinInts(mids))
	if err = d.httpClient.Post(c, d.notifyURI(), "", params, &res); err != nil {
		log.Error("d.httpClient.Post(%s,%v,%d)", d.notifyURI(), params, err)
		return
	}
	if res.Code != 0 {
		err = errors.New("code != 0")
		log.Error("d.httpClient.Post(%s,%v,%v,%d)", d.notifyURI(), params, err, res.Code)
	}
	return
}
