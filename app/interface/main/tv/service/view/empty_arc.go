package view

import (
	"time"

	"go-common/library/log"
)

func (s *Service) emptyArcproc() {
	var (
		ps        = s.conf.Cfg.EmptyArc.UnshelvePS
		emptyAids = make(map[int64]int, ps)
	)
	for {
		aid, ok := <-s.emptyArcCh
		if !ok {
			log.Warn("[emptyArcproc] channel quit")
			return
		}
		emptyAids[aid] = 1
		if len(emptyAids) < ps { // not enough cid, stay waiting
			time.Sleep(2 * time.Second)
			continue
		}
		distinctAIDs := pickKeys(emptyAids)
		emptyAids = make(map[int64]int, ps)
		if err := s.cmsDao.UnshelveArcs(ctx, distinctAIDs); err != nil {
			log.Error("emptyArc Aids %v, Err %v", distinctAIDs, err)
			continue
		}
		log.Info("emptyArc Apply %d Aids: %v", len(distinctAIDs), distinctAIDs)
	}
}

func pickKeys(q map[int64]int) (res []int64) {
	for k := range q {
		res = append(res, k)
	}
	return
}
