package service

import (
	"context"
	"time"

	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/model/data"
	xtime "go-common/library/time"
)

// UpBaseStats .
func (s *Service) UpBaseStats(c context.Context, req *upgrpc.UpStatReq) (res *upgrpc.UpBaseStatReply, err error) {
	res = new(upgrpc.UpBaseStatReply)
	if req.Date.Time().IsZero() {
		// 如果没有填，则取最新的数据，如果有填，则取对应天数的数据，这里不需要做什么操作
		// 12点更新数据，数据表为昨天日期，所以在12点以前，要读前天的表
		req.Date = xtime.Time(time.Now().Add(-12*time.Hour).AddDate(0, 0, -1).Unix())
	}
	var stat *data.UpBaseStat
	if stat, err = s.Data.BaseUpStat(c, req.Mid, req.Date.Time().Format("20060102")); err != nil {
		return
	}
	stat.CopyToReply(res)
	return
}
