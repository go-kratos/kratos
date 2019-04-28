package service

import (
	"context"
	"time"

	"go-common/app/job/main/usersuit/model"
	"go-common/library/log"
)

var (
	t = time.NewTimer(time.Minute * 5)
)

// startexpireproc start
func (s *Service) startexpireproc() {
	for range t.C {
		s.expiredEquip(context.TODO())
		t.Reset(time.Minute * 10)
	}
}

// expiredEquip operator equipment info
func (s *Service) expiredEquip(c context.Context) (err error) {
	var (
		mids     []int64
		expires  = time.Now().Unix()
		affected int64
	)
	if mids, err = s.pendantDao.ExpireEquipPendant(c, expires); err != nil || len(mids) == 0 {
		log.Error("s.pendantDao.ExpireEquipPendant(%d) error(%+v)", expires, err)
		return
	}
	for _, mid := range mids {
		if affected, err = s.pendantDao.UpEquipMID(c, mid); err != nil || affected == 0 {
			log.Error("s.pendantDao.UpEquipMID(%d) error(%+v)", mid, err)
			continue
		}
		s.pendantDao.DelEquipCache(c, mid)
		tid := mid
		s.addNotify(func() {
			s.accNotify(context.TODO(), tid, model.AccountNotifyUpdatePendant)
		})
	}
	return
}
