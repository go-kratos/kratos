package service

import (
	"context"
	"time"

	"go-common/app/service/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	ctime "go-common/library/time"
)

// BaseInfo get user's base info.
func (s *Service) BaseInfo(c context.Context, mid int64) (info *model.BaseInfo, err error) {
	if mid < 1 {
		err = ecode.RequestErr
		return
	}
	var mc = true
	if info, err = s.mbDao.BaseInfoCache(c, mid); err != nil {
		mc = false
		err = nil // ignore error
	}
	if info != nil {
		if info.Mid == 0 {
			err = ecode.MemberNotExist
		}
		return
	}
	if info, err = s.mbDao.BaseInfo(c, mid); err != nil {
		return
	}
	if info == nil {
		err = ecode.MemberNotExist
		info = &model.BaseInfo{Mid: 0}
	}
	if mc {
		fanoutErr := s.cache.Do(c, func(ctx context.Context) {
			s.mbDao.SetBaseInfoCache(ctx, mid, info)
		})
		if fanoutErr != nil {
			log.Error("fanout do err(%+v)", fanoutErr)
		}
	}
	return
}

// BatchBaseInfo get user's base info.
func (s *Service) BatchBaseInfo(c context.Context, mids []int64) (res map[int64]*model.BaseInfo, err error) {
	var mc = true
	if len(mids) > 100 {
		err = ecode.MemberOverLimit
		return
	}
	var (
		missed []int64
		info   *model.BaseInfo
	)
	if res, missed, err = s.mbDao.BatchBaseInfoCache(c, mids); err != nil {
		log.Error("infoDao.BatchBaseInfoCache err(%v)", err)
	}
	var missb []*model.BaseInfo
	for _, mid := range missed {
		if info, err = s.mbDao.BaseInfo(c, mid); err != nil {
			continue
		}
		if info != nil {
			res[mid] = info
			missb = append(missb, info)
		}
	}
	if len(missb) == 0 {
		return
	}
	if mc {
		fanoutErr := s.cache.Do(c, func(ctx context.Context) {
			s.mbDao.SetBatchBaseInfoCache(ctx, missb)
		})
		if fanoutErr != nil {
			log.Error("fanout do err(%+v)", fanoutErr)
		}
	}
	return
}

// loadOfficial load official into memory.
func (s *Service) loadOfficial() (err error) {
	origin := s.officials
	om, err := s.mbDao.Officials(context.Background())
	if err != nil {
		return
	}
	s.officials = om

	// compare and notify purge cache
	if origin == nil {
		// skip on first load
		return
	}
	for mid, of := range om {
		if oof := origin[mid]; oof == nil || !of.Equal(oof) {
			log.Info("Notify to purge official cache: %d", mid)
			s.mbDao.NotifyPurgeCache(context.Background(), mid, model.ActUpdateByAdmin)
		}
	}
	for mid := range origin {
		if of := om[mid]; of == nil {
			log.Info("Notify to purge official cache: %d", mid)
			s.mbDao.NotifyPurgeCache(context.Background(), mid, model.ActUpdateByAdmin)
		}
	}
	return
}

// loadOfficialproc load official into memory.
func (s *Service) loadOfficialproc() {
	for {
		if err := s.loadOfficial(); err != nil {
			time.Sleep(60 * time.Second)
			continue
		}
		time.Sleep(5 * time.Minute)
	}
}

// SetRank set user's rank
func (s *Service) SetRank(c context.Context, mid, rank int64) (err error) {
	if err = s.mbDao.SetRank(c, mid, rank); err != nil {
		return
	}
	return s.delCacheAndNotify(c, mid, model.ActUpdatePersonInfo)
}

// SetSex Set sex of user
func (s *Service) SetSex(c context.Context, mid, sex int64) (err error) {
	if err = s.mbDao.SetSex(c, mid, sex); err != nil {
		return
	}
	return s.delCacheAndNotify(c, mid, model.ActUpdatePersonInfo)
}

// SetName update user's name
func (s *Service) SetName(c context.Context, mid int64, name string) (err error) {
	if err = s.mbDao.SetName(c, mid, name); err != nil {
		return
	}
	return s.delCacheAndNotify(c, mid, model.ActUpdateUname)
}

// SetSign update user's sign
func (s *Service) SetSign(c context.Context, mid int64, sign string) (err error) {
	if err = s.mbDao.SetSign(c, mid, sign); err != nil {
		return
	}
	return s.delCacheAndNotify(c, mid, model.ActUpdatePersonInfo)
}

// SetBirthday set birthday of user
func (s *Service) SetBirthday(c context.Context, mid int64, birthday ctime.Time) (err error) {
	if err = s.mbDao.SetBirthday(c, mid, birthday); err != nil {
		return
	}
	return s.delCacheAndNotify(c, mid, model.ActUpdatePersonInfo)
}

// SetFace set face.
func (s *Service) SetFace(c context.Context, mid int64, face string) (err error) {
	if err = s.mbDao.SetFace(c, mid, face); err != nil {
		return
	}
	return s.delCacheAndNotify(c, mid, model.ActUpdateFace)
}

// SetBase set base.
func (s *Service) SetBase(c context.Context, base *model.BaseInfo) (err error) {
	base.Face = ""
	if base.Rank == 0 {
		base.Rank = model.DefaultRank
	}
	if base.Birthday == 0 {
		base.Birthday = model.DefaultTime
	}
	if err = s.mbDao.SetBase(c, base); err != nil {
		return
	}
	return s.mbDao.DelBaseInfoCache(c, base.Mid)
}

func (s *Service) delCacheAndNotify(c context.Context, mid int64, action string) (err error) {
	if err = s.mbDao.DelBaseInfoCache(c, mid); err != nil {
		return
	}
	if err = s.mbDao.NotifyPurgeCache(c, mid, action); err != nil {
		return
	}
	return
}
