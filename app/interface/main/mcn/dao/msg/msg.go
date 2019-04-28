package msg

import (
	"context"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// MutliSendSysMsg Mutli send sys msg.
func (d *Dao) MutliSendSysMsg(c context.Context, allUids []int64, mc, title string, context string, ip string) (err error) {
	var times int
	ulen := len(allUids)
	if ulen%100 == 0 {
		times = ulen / 100
	} else {
		times = ulen/100 + 1
	}
	var uids []int64
	for i := 0; i < times; i++ {
		if i == times-1 {
			uids = allUids[i*100:]
		} else {
			uids = allUids[i*100 : (i+1)*100]
		}
		if err = d.SendSysMsg(c, uids, mc, title, context, ip); err != nil {
			continue
		}
	}
	return
}

// SendSysMsg send sys msg.
func (d *Dao) SendSysMsg(c context.Context, uids []int64, mc, title string, context string, ip string) (err error) {
	params := url.Values{}
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("data_type", "4")
	params.Set("context", context)
	params.Set("mid_list", xstr.JoinInts(uids))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Status int8   `json:"status"`
			Remark string `json:"remark"`
		} `json:"data"`
	}
	if err = d.client.Post(c, d.msgURL, ip, params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrapf(ecode.Int(res.Code), "SendSysMsg d.client.Post(%s,%d)", d.msgURL+"?"+params.Encode(), res.Code)
	} else {
		log.Info("send msg ok, resdata=%+v", res.Data)
	}
	return
}
