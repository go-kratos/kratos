package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/main/member/model"
)

func notifyKey(mid int64) string {
	return fmt.Sprintf("MemberService-AccountNotify-T%d", mid)
}

// AddExplog add exp log with databus
func (d *Dao) AddExplog(c context.Context, mid, exp, toExp int64, oper, reason, ip string) (err error) {
	log := &model.UserLog{
		Mid:   mid,
		IP:    ip,
		TS:    time.Now().Unix(),
		LogID: model.UUID4(),
		Content: map[string]string{
			"from_exp": strconv.FormatInt(exp, 10),
			"to_exp":   strconv.FormatInt(toExp, 10),
			"operater": oper,
			"reason":   reason,
		},
	}
	err = d.logDatabus.Send(c, strconv.FormatInt(mid, 10), log)
	return
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
