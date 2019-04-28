package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/member/model"
)

func notifyKey(mid int64) string {
	return fmt.Sprintf("MemberJob-AccountNotify-T%d", mid)
}

// NotifyPurgeCache is
func (d *Dao) NotifyPurgeCache(c context.Context, mid int64, action string) error {
	msg := &model.NotifyInfo{
		Mid:    mid,
		Action: action,
	}
	key := notifyKey(mid)
	return d.accNotify.Send(c, key, msg)
}
