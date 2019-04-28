package message

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/job/main/videoup/conf"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"strings"
	"time"
)

// Dao is message dao.
type Dao struct {
	c       *conf.Config
	client  *xhttp.Client
	uri     string
	pushURI string
}

// New new a message dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:       c,
		client:  xhttp.NewClient(c.HTTPClient),
		uri:     c.Host.Message + "/api/notify/send.user.notify.do",
		pushURI: c.Host.API + "/x/internal/push/single",
	}
	return
}

// Send send message to upper.
func (d *Dao) Send(c context.Context, mc, title, msg string, mid int64, ts int64) (err error) {
	params := url.Values{}
	source := strings.Split(mc, "_")
	params.Set("type", "json")
	params.Set("source", source[0])
	params.Set("data_type", "4")
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("context", msg)
	params.Set("mid_list", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}

	if err = d.client.Post(c, d.uri, "", params, &res); err != nil {
		log.Error("message url(%s) error(%v)", d.uri+"?"+params.Encode(), err)
		return
	}
	log.Info("send.user.notify.do send mid(%d) message(%s) url(%s) code(%v)", mid, msg, d.uri+"?"+params.Encode(), res.Code)
	if res.Code != 0 {
		log.Error("message url(%s) error(%v)", d.uri+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("message send failed")
	}
	return
}

//PushMsg 发送推送消息给创作姬APP TODO deprecated
func (d *Dao) PushMsg(c context.Context, mid int64, title string, msg string) (err error) {
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	params := url.Values{}
	params.Set("appkey", d.c.HTTPClient.Key)
	params.Set("appsecret", d.c.HTTPClient.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("app_id", d.c.Host.Push.AppID)
	params.Set("business_id", d.c.Host.Push.BusinessID)
	params.Set("token", d.c.Host.Push.Token)
	params.Set("link_type", "8")
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("alert_title", title)
	params.Set("alert_body", msg)

	log.Info("start send PushMsg：url(%s)", d.pushURI+"?"+params.Encode())
	if err = d.client.Post(c, d.pushURI, "", params, &res); err != nil {
		log.Error("message PushMsg error(%v), url(%s)", err, d.pushURI+"?"+params.Encode())
		return
	}
	log.Info("after send PushMsg：response(%v) url(%s)", res, d.pushURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Info("message PushMsg response is failed, res message(%+v), url(%s)", res, d.pushURI+"?"+params.Encode())
		err = fmt.Errorf("message PushMsg send failed")
	}
	return
}
