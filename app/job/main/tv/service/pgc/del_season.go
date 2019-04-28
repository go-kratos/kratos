package pgc

import (
	"database/sql"
	"time"

	"go-common/app/job/main/tv/dao/lic"
	"go-common/library/log"
)

// sync the deleted season data to the license owner
func (s *Service) delSeason() {
	var (
		sign   = s.c.Sync.Sign
		prefix = s.c.Sync.AuditPrefix
	)
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("delSeason DB closed!")
			return
		}
		delSeason, err := s.dao.DelSeason(ctx)
		if err == sql.ErrNoRows || len(delSeason) == 0 {
			log.Info("No deleted data to pick from Season to sync")
			time.Sleep(time.Duration(s.c.Sync.Frequency.FreModSeason))
			continue
		}
		for _, v := range delSeason {
			data := lic.DelLic(sign, prefix, v.ID)
			// ignore the program part during modified season sync
			body := lic.PrepareXML(data)
			res, err := s.licDao.CallRetry(ctx, s.c.Sync.API.DelSeasonURL, body)
			// 3 times still error
			if err != nil {
				log.Error("DelSeasonURL interface not available!Sid: %v, Err: %v", v.ID, err)
				s.dao.DelaySeason(ctx, v.ID)
				time.Sleep(time.Duration(s.c.Sync.Frequency.ErrorWait))
				// avoid always be stuck by one error data
				break
			}
			if err == nil && res != nil {
				_, err := s.dao.RejectSeason(ctx, int(v.ID))
				if err != nil {
					log.Error("DelSeasonSync season %v to rejected fail!", v.ID)
					// sync next one
					continue
				}
			}
		}
		// break after each loop
		time.Sleep(1 * time.Second)
	}
}
