package conf

import (
	"go-common/library/log"
	"time"
)

// RecordTimeCostLog 记录日志
func RecordTimeCostLog(nowTime int64, desc string) {
	log.Info(desc+"|%d", nowTime)
}

// RecordTimeCost ...
// 记录时间
func RecordTimeCost() (NowTime int64) {
	return time.Now().UnixNano() / 1000000 // 用毫秒
}
