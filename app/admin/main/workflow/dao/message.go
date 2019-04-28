package dao

import (
	"context"

	"go-common/app/admin/main/workflow/model/param"
	"go-common/library/ecode"
)

const _userNotifyURI = "http://message.bilibili.co/api/notify/send.user.notify.do"

// SendMessage send message to upper.
func (d *Dao) SendMessage(c context.Context, msg *param.MessageParam) (err error) {
	uri := _userNotifyURI
	params := msg.Query()
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpWrite.Post(c, uri, "", params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		return
	}
	return
}
