package dao

import (
	"context"
	"net/url"

	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	couponSystemNotify = "4"
	couponMC           = "10_99_1"

	sendMessage = "/api/notify/send.user.notify.do"
)

//SendMessage send message.
func (d *Dao) SendMessage(c context.Context, mids, title, content, ip string) (err error) {
	params := url.Values{}
	params.Set("mc", couponMC)
	params.Set("title", title)
	params.Set("context", content)
	params.Set("data_type", couponSystemNotify)
	params.Set("mid_list", mids)
	if err = d.client.Post(c, d.c.Prop.MessageURL+sendMessage, ip, params, nil); err != nil {
		err = errors.WithStack(err)
	}
	log.Info("send message url:%+v params:%+v err:%+v", d.c.Prop.MessageURL+sendMessage, params, err)
	return
}
