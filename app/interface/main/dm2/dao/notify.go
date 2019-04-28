package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_notifyURI = "/api/notify/send.user.notify.do"
)

func (d *Dao) notifyURI() string {
	return d.conf.Host.Message + _notifyURI
}

// SendNotify 发送站内信
func (d *Dao) SendNotify(c context.Context, title, content, dataType string, mids []int64) (err error) {
	res := struct {
		Code int `json:"code"`
		Data struct {
			TotalCount int
			ErrorCount int
		} `json:"data"`
	}{}
	params := url.Values{}
	params.Set("mc", "1_8_2")
	params.Set("title", title)
	params.Set("data_type", dataType)
	params.Set("context", content)
	params.Set("mid_list", xstr.JoinInts(mids))
	if err = d.httpCli.Post(c, d.notifyURI(), "", params, &res); err != nil {
		log.Error("d.httpClient.Post(%s,%v,%d)", d.notifyURI(), params, err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("res.Code:;%d", res.Code)
		log.Error("d.httpClient.Post(%s,%v,%v,%d)", d.notifyURI(), params, err, res.Code)
	}
	return
}
