package service

import (
	"time"

	mlog "go-common/app/admin/main/apm/model/log"
	"go-common/library/log"
	context "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus/report"
)

// SQLLog log
type SQLLog struct {
	SQLType string
	Content interface{}
}

// LogAdd add log
func (s *Service) LogAdd(c context.Context, lg *mlog.Log) (err error) {
	l := &mlog.Log{
		UserName: lg.UserName,
		Business: lg.Business,
		Info:     lg.Info,
	}
	if err = s.dao.DB.Create(&l).Error; err != nil {
		log.Error("s.LogAdd create error(%v)", err)
	}
	return
}

// SendLog log
func (s *Service) SendLog(c context.Context, username string, uid int64, tp int, oid int64, action string, context interface{}) (err error) {
	report.Manager(&report.ManagerInfo{
		Uname:    username,
		UID:      uid,
		Business: 71,
		Type:     tp, // 1 add 2 update 3 delete 4 soft delete 5 Transaction 6 kafka
		Oid:      oid,
		Action:   action,
		Ctime:    time.Now(),
		// Index:    []interface{}{0, 0},
		Content: map[string]interface{}{
			"content": context,
		},
	})
	return
}
