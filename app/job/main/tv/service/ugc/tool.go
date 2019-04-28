package ugc

import (
	appDao "go-common/app/job/main/tv/dao/app"
	arccli "go-common/app/service/main/archive/api"
	"go-common/library/log"
)

// call ArcRPC for types data
func (s *Service) loadTypes() {
	var (
		resp *arccli.TypesReply
		err  error
	)
	if resp, err = s.arcClient.Types(ctx, &arccli.NoArgRequest{}); err != nil {
		log.Error("arcRPC loadType Error %v", err)
		return
	}
	s.arcTypes = resp.Types
}

func (s *Service) hitPGC(tid int32) (hit bool) {
	_, hit = s.pgcTypes[s.getPTypeName(tid)]
	return
}

func (s *Service) delPGC(tid int32, aid int64) (hit bool, err error) {
	if hit = s.hitPGC(tid); !hit { // if not hit, do nothing
		appDao.PromInfo("HitPGC:FdSucc")
		return
	}
	log.Info("delPGC Aid %d, Tid %d", aid, tid)
	appDao.PromInfo("HitPGC:DelSucc")
	if err = s.delArc(aid); err != nil { // if hit, delete it if exist
		appDao.PromInfo("HitPGC:DelErr")
		log.Error("HitPGC DelArc %d, Err %v", aid, err)
	}
	return
}

func pickKeys(q map[int64]int) (res []int64) {
	for k := range q {
		res = append(res, k)
	}
	return
}
