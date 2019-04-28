package util

import (
	"time"

	"go-common/library/queue/databus/report"
)

//AddLogs add action logs
func AddLogs(logtype int, uname string, uid int64, oid int64, action string, obj interface{}) (err error) {
	report.Manager(&report.ManagerInfo{
		Uname:    uname,
		UID:      uid,
		Business: 1,
		Type:     logtype,
		Oid:      oid,
		Action:   action,
		Ctime:    time.Now(),
		// extra
		Index: []interface{}{},
		Content: map[string]interface{}{
			"json": obj,
		},
	})
	return
}
