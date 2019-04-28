package data

import (
	"context"
	"errors"
	"go-common/app/job/main/videoup-report/model/monitor"
)

func (d *Dao) MonitorNotify(c context.Context) (data []*monitor.RuleResultData, err error) {
	var (
		res = &monitor.RuleResultRes{}
	)
	if err = d.client.Get(c, d.moniNotifyURL, "", nil, &res); err != nil {
		return
	}
	if res == nil || res.Data == nil {
		err = errors.New("监控结果获取失败")
		return
	}
	data = res.Data
	return
}
