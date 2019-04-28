package service

import (
	"context"
	"time"

	"go-common/library/log"
	manager "go-common/library/queue/databus/report"
)

// AddAuditLog .
func (s *Service) AddAuditLog(c context.Context, bizID int, tp int8, action string, uid int64, uname string, oids []int64, index []interface{}, content map[string]interface{}) error {
	var err error
	for _, oid := range oids {
		managerInfo := &manager.ManagerInfo{
			UID:      uid,
			Uname:    uname,
			Business: bizID,
			Type:     int(tp),
			Action:   action,
			Oid:      oid,
			Ctime:    time.Now(),
			Index:    index,
			Content:  content,
		}
		if err = manager.Manager(managerInfo); err != nil {
			log.Error("manager.Manager(%+v) error(%+v)", managerInfo, err)
			continue
		}
		log.Info("s.managerSendLog(%+v)", managerInfo)
	}
	return err
}
