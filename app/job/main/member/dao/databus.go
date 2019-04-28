package dao

import (
	"context"
	"strconv"

	"go-common/app/job/main/member/model"
)

// DatabusAddLog add exp log messager to databus
func (d *Dao) DatabusAddLog(c context.Context, mid, exp, toExp, ts int64, oper, reason, ip string) (err error) {
	log := &model.UserLog{
		Mid: mid,
		TS:  ts,
		IP:  ip,
		Content: map[string]string{
			"from_exp": strconv.FormatInt(exp, 10),
			"to_exp":   strconv.FormatInt(toExp, 10),
			"operater": oper,
			"reason":   reason,
		},
	}
	err = d.plogDatabus.Send(c, strconv.FormatInt(mid, 10), log)
	return
}
