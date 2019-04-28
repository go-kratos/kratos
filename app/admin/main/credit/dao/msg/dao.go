package msg

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/credit/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_msgURL = "/api/notify/send.user.notify.do"
)

// Dao struct info of Dao.
type Dao struct {
	// http
	client *bm.Client
	// conf
	c      *conf.Config
	msgURL string
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// http client
		client: bm.NewClient(c.HTTPClient),
	}
	d.msgURL = c.Host.Msg + _msgURL
	return
}

// SendSysMsg send sys msg.
func (dao *Dao) SendSysMsg(c context.Context, mid int64, title string, context string) (err error) {
	params := url.Values{}
	params.Set("mc", "2_1_13")
	params.Set("title", title)
	params.Set("data_type", "4")
	params.Set("context", context)
	params.Set("mid_list", fmt.Sprintf("%d", mid))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Status int8   `json:"status"`
			Remark string `json:"remark"`
		} `json:"data"`
	}
	for i := 0; i <= 5; i++ {
		if err = dao.client.Post(c, dao.msgURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
			log.Error("sendMsgURI(%s) error(%v)", dao.msgURL+"?"+params.Encode(), err)
			continue
		}
		if res.Code != 0 {
			log.Error("sendMsgURI(%s) error(%v)", dao.msgURL+"?"+params.Encode(), res.Code)
			err = ecode.Int(res.Code)
			continue
		}
		return
	}
	return
}
