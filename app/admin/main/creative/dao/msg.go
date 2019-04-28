package dao

import (
	"context"
	"net/url"

	"go-common/library/ecode"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"go-common/library/log"

	"github.com/pkg/errors"
)

// MutliSendSysMsg Mutli send sys msg.
func (d *Dao) MutliSendSysMsg(c context.Context, allMids []int64, title string, context string) (err error) {
	var times int
	ulen := len(allMids)
	if ulen%100 == 0 {
		times = ulen / 100
	} else {
		times = ulen/100 + 1
	}
	var mids []int64
	for i := 0; i < times; i++ {
		if i == times-1 {
			mids = allMids[i*100:]
		} else {
			mids = allMids[i*100 : (i+1)*100]
		}
		if err = d.SendSysMsg(c, mids, title, context); err != nil {
			err = errors.Wrapf(err, "d.SendSysMsg(%+v,%s,%s)", mids, title, context)
			continue
		}
	}
	return
}

// SendSysMsg send sys msg.
func (d *Dao) SendSysMsg(c context.Context, mids []int64, title string, context string) (err error) {
	var ip = metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("mc", "1_14_6")
	params.Set("title", title)
	params.Set("data_type", "4")
	params.Set("context", context)
	params.Set("mid_list", xstr.JoinInts(mids))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Status int8   `json:"status"`
			Remark string `json:"remark"`
		} `json:"data"`
	}
	if err = d.client.Post(c, d.msgURL, ip, params, &res); err != nil {
		err = errors.Wrapf(err, "SendSysMsg d.client.Post(%s)", d.msgURL+"?"+params.Encode())
		return
	}
	log.Info("dao.SendSysMsg res (%+v) ", res)
	if res.Code != 0 {
		err = errors.Wrapf(ecode.Int(res.Code), "SendSysMsg d.client.Post(%s,%d)", d.msgURL+"?"+params.Encode(), res.Code)
	}
	return
}
