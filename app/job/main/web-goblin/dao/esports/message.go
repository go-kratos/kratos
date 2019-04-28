package esports

import (
	"context"
	"net/url"

	mdlesp "go-common/app/job/main/web-goblin/model/esports"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

var _notify = "4"

// SendMessage send system notify.
func (d *Dao) SendMessage(mids []int64, msg string, contest *mdlesp.Contest) (err error) {
	params := url.Values{}
	params.Set("mid_list", xstr.JoinInts(mids))
	params.Set("title", d.c.Rule.AlertTitle)
	params.Set("mc", d.c.Message.MC)
	params.Set("data_type", _notify)
	params.Set("context", msg)
	var res struct {
		Code int `json:"code"`
	}
	err = d.messageHTTPClient.Post(context.Background(), d.c.Message.URL, "", params, &res)
	if err != nil {
		log.Error("SendMessage d.messageHTTPClient.Post(%s) error(%+v)", d.c.Message.URL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("SendMessage url(%s) res code(%d)", d.c.Message.URL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

//Batch 批量处理
func (d *Dao) Batch(list []int64, msg string, contest *mdlesp.Contest, batchSize int, f func(users []int64, msg string, contest *mdlesp.Contest) error) {
	if msg == "" {
		log.Warn("Batch msg is empty")
		return
	}
	retry := d.c.Push.RetryTimes
	for {
		var (
			mids []int64
			err  error
		)
		l := len(list)
		if l == 0 {
			break
		} else if l <= batchSize {
			mids = list[:l]
		} else {
			mids = list[:batchSize]
			l = batchSize
		}
		list = list[l:]

		for i := 0; i < retry; i++ {
			if err = f(mids, msg, contest); err == nil {
				break
			}
		}
		if err != nil {
			log.Error("Batch error(%v), params(%s)", err, msg)
		}
	}
}
