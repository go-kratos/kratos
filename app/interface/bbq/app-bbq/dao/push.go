package dao

import (
	"context"
	"fmt"

	notice "go-common/app/service/bbq/notice-service/api/v1"
	"go-common/library/log"
)

// PushLogin .
func (d *Dao) PushLogin(c context.Context, req *notice.UserPushDev) (err error) {
	_, err = d.noticeClient.PushLogin(c, req)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("push login fail: req=%s", req.String())))
	}
	return
}

// PushLogout .
func (d *Dao) PushLogout(c context.Context, req *notice.UserPushDev) (err error) {
	_, err = d.noticeClient.PushLogout(c, req)
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("push logout fail: req=%s", req.String())))
	}
	return
}

// PushCallback .
func (d *Dao) PushCallback(c context.Context, tid string, nid string, mid int64, buvid string) (err error) {
	_, err = d.noticeClient.PushCallback(c, &notice.PushCallbackRequest{
		Tid:   tid,
		Nid:   nid,
		Mid:   mid,
		Buvid: buvid,
	})

	return
}
