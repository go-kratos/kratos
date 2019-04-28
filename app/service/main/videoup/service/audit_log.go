package service

import (
	"context"
	"time"

	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

// AddAuditLog .
func (s *Service) AddAuditLog(c context.Context, bizID int, tp int8, action string, uid int64, uname string, oids []int64, index []interface{}, content map[string]interface{}) error {
	var err error
	for _, oid := range oids {
		userInfo := &report.UserInfo{
			Mid:      uid,
			Business: bizID,
			Type:     int(tp),
			Action:   action,
			Oid:      oid,
			Ctime:    time.Now(),
			Index:    index,
			Content:  content,
		}
		if err = report.User(userInfo); err != nil {
			log.Error("manager.User(%+v) error(%+v)", userInfo, err)
			continue
		}
		log.Info("s.AddAuditLog(%+v)", userInfo)
	}
	return err
}
