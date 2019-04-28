package dao

import (
	"context"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

const (
	couponSystemNotify = 4
	couponMC           = "10_99_1"

	sendMessage = "/api/notify/send.user.notify.do"
)

//SendMessage send message.
func (d *Dao) SendMessage(c context.Context, mids, content, title string) (err error) {
	params := url.Values{}
	params.Set("mc", couponMC)
	params.Set("title", title)
	params.Set("context", content)
	params.Set("data_type", strconv.FormatInt(couponSystemNotify, 10))
	params.Set("mid_list", mids)
	defer func() {
		log.Info("send message url:%+v params:%+v err:%+v    ", d.c.Property.MessageURL+sendMessage, params, err)
	}()
	ip := metadata.String(c, metadata.RemoteIP)
	if err = d.client.Post(c, d.c.Property.MessageURL+sendMessage, ip, params, nil); err != nil {
		err = errors.WithStack(err)
	}
	return
}
