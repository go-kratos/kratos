package service

import (
	"context"
	"runtime/debug"

	"go-common/app/job/main/usersuit/model"
	"go-common/library/log"
)

var (
	_upTaskMedal = map[int64]int64{1: 4, 2: 3, 3: 1, 4: 7, 5: 6, 6: 5, 7: 10, 8: 9, 9: 8}
)

func (s *Service) cronUpNameplate() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("lrucleanproc panic %v : %s", x, debug.Stack())
			go s.cronUpNameplate()
		}
	}()
	var (
		err error
		res *model.UpInfo
		ctx = context.TODO()
	)
	if res, err = s.medalDao.UpInfoData(ctx); err != nil {
		log.Error("s.medalDao.UpInfoData err(%+v)", err)
		return
	}
	for _, item := range res.Data {
		var (
			nid int64
			ok  bool
		)
		if nid, ok = _upTaskMedal[item.ID]; !ok {
			continue
		}
		for _, mid := range item.Mids {
			log.Info("s.medalDao.AddMedalOwner(%d, %d)", mid, nid)
			if err = s.medalDao.Grant(ctx, mid, nid); err != nil {
				log.Error("s.medalDao.AddMedalOwner(%d, %d) err(%+v)", mid, nid, err)
			}
		}
	}

}
