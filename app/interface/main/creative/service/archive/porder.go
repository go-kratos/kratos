package archive

import (
	"context"

	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/log"
)

// Porder fn
func (s *Service) Porder(c context.Context, aid int64) (pd *archive.Porder, err error) {
	log.Warn("Porder with aid (%+v)", aid)
	cache := true
	if pd, err = s.arc.POrderCache(c, aid); err != nil {
		err = nil
		cache = false
	} else if pd != nil {
		s.pCacheHit.Incr("porder_cache")
		return
	}
	s.pCacheMiss.Incr("porder_cache")
	if pd, err = s.arc.Porder(c, aid); err != nil {
		log.Error("s.porder.Porder aid(%d) err(%v)", aid, err)
		return
	}
	if cache {
		s.addCache(func() {
			if pd == nil {
				pd = &archive.Porder{}
			}
			s.arc.AddPOrderCache(context.Background(), aid, pd)
		})
	}
	return
}

// FlowJudge fn
func (s *Service) FlowJudge(c context.Context, business, groupID int64, oids []int64) (hitOids []int64, err error) {
	if hitOids, err = s.arc.FlowJudge(c, business, groupID, oids); err != nil {
		log.Error("s.porder.FlowJudge business(%d) groupID(%d) err(%v)", business, groupID, err)
		return
	}
	return
}
