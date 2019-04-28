package service

import (
	"context"

	"go-common/library/log"
)

// kfcActionDeal.
func (s *Service) kfcActionDeal(j int) {
	defer s.waiter.Done()
	var (
		ch = s.kfcActionCh[j]
		c  = context.Background()
	)
	log.Info("kfcActionDeal goroutine(%d) start", j)
	for {
		ms, ok := <-ch
		if !ok {
			log.Warn("kfcActionDeal(%d): quit", j)
			return
		}
		if err := s.kfcDao.KfcDelver(c, ms.CouponID, ms.UID); err != nil {
			log.Error("kfcActionDeal(%d):s.kfcDao.KfcDelver(%d %d) error(%v)", j, ms.CouponID, ms.UID, err)
			return
		}
		log.Info("kfcActionDeal(%d) success id(%d) uid(%d)", j, ms.CouponID, ms.UID)
	}
}
