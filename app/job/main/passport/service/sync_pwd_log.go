package service

import (
	"context"
	"time"

	"go-common/library/log"
)

func (s *Service) syncPwdLog() {
	id := s.c.Sync.SyncPwdID
	for {
		pwds, err := s.d.BatchGetPwdLog(context.Background(), id)
		if err != nil {
			log.Error("failed to batch get pwd log, s.d.BatchGetPwdLog(%d), error(%v)", id, err)
			time.Sleep(1 * time.Second)
			continue
		}
		log.Info("SyncPwdID (%d), len(pwds) (%d)", id, len(pwds))
		if len(pwds) == 0 {
			break
		}
		for _, pwd := range pwds {
			if err := s.d.AddPwdLogHBase(context.Background(), pwd); err != nil {
				log.Error("failed to add pwd log to hbase, service.dao.AddLoginLogHBase(%+v) error(%v)", pwd, err)
				time.Sleep(1 * time.Second)
				continue
			}
			id = pwd.ID
		}
	}
}
