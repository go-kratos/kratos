package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/log"
)

// consts
const (
	notifyURL = "http://message.bilibili.co/api/notify/send.user.notify.do"
	dataType  = "4"
	source    = "2"
)

// SendMessage is
func (d *Dao) SendMessage(c context.Context, mid int64, title, msg, mc string) (err error) {
	var (
		params = url.Values{}
	)
	params.Set("mid_list", fmt.Sprintf("%d", mid))
	params.Set("mc", mc)
	params.Set("data_type", dataType)
	params.Set("source", source)
	params.Set("context", msg)
	params.Set("title", title)
	var resp struct {
		Code int `json:"code"`
	}
	log.Info("SendMessage() params(%v)", params)
	if err = d.client.Post(c, notifyURL, "", params, &resp); err != nil {
		log.Error("d.client.Post(%s,%+v) error(%v)", notifyURL, params, err)
		return
	}
	if resp.Code != 0 {
		err = ecode.Int(resp.Code)
		log.Error("d.client.Post(%s,%+v) resp.Code(%d)", notifyURL, params, resp.Code)
	}
	return
}
