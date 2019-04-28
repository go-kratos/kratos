package service

import (
	"context"
	"time"

	"go-common/app/service/main/resource/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

//loadRelateCache load relate card cache
func (s *Service) loadRelateCache() {
	relate, err := s.show.Relate(context.Background(), time.Now())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	pgcMapRelate := make(map[int64]int64)
	relateCache := make(map[int64]*model.Relate)
	for _, r := range relate {
		var pgcIDs []int64
		if pgcIDs, err = xstr.SplitInts(r.PgcIDs); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", r.PgcIDs, err)
			return
		}
		if len(pgcIDs) > 0 {
			for _, pgcID := range pgcIDs {
				pgcMapRelate[pgcID] = r.ID
			}
		}
		relateCache[r.ID] = r
	}
	s.relatePgcMapCache = pgcMapRelate
	s.relateCache = relateCache
}
