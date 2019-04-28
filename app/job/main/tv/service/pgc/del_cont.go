package pgc

import (
	"database/sql"
	"time"

	"go-common/app/job/main/tv/dao/lic"
	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
)

// sync the deleted EP data to the license owner
func (s *Service) delCont() {
	var (
		sign   = s.c.Sync.Sign
		prefix = s.c.Sync.AuditPrefix
	)
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("delCont DB closed!")
			return
		}
		// pick data
		delCont, err := s.dao.DelCont(ctx)
		if err == sql.ErrNoRows || len(delCont) == 0 {
			log.Info("No deleted data to pick from Cont to sync")
			time.Sleep(time.Duration(s.c.Sync.Frequency.FreModSeason))
			continue
		}
		delEpids := []int{}
		for _, v := range delCont {
			delEpids = append(delEpids, v.EPID)
		}
		s.dao.DelaySync(ctx, delCont) // avoid always be stuck by one error data
		body := lic.DelEpLic(prefix, sign, delEpids)
		// call API
		var res *model.Document
		res, err = s.licDao.CallRetry(ctx, s.c.Sync.API.DelEPURL, body)
		// 3 times still error
		if err != nil {
			log.Error("DelEPURL interface not available! %v", err)
			time.Sleep(time.Duration(s.c.Sync.Frequency.ErrorWait))
			continue
		}
		// update the state
		if err == nil && res != nil {
			for _, v := range delCont {
				_, err := s.dao.SyncCont(ctx, v.EPID)
				if err != nil {
					log.Error("SyncCont EP %v to auditing fail!", v.ID)
					continue
				}
			}
		}
		// break after each loop
		time.Sleep(1 * time.Second)
	}
}
