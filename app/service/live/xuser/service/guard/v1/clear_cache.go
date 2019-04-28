package v1

import (
	"context"
	dahanghaiModel "go-common/app/service/live/xuser/model/dhh"
	"go-common/library/net/metadata"

	"go-common/library/log"
)

func (s *GuardService) asyncCLearExpCache(ctx context.Context, uid int64) (err error) {
	err = s.dao.DelDHHFromRedis(ctx, uid)
	if err != nil {
		log.Error(_errorServiceLogPrefix+"|clear UserExp cache(%d) error(%v)", uid, err)
	}
	return
}

func (s *GuardService) asyncSetDHHCache(ctx context.Context, req []*dahanghaiModel.DaHangHaiRedis2, uid int64) error {
	c := metadata.WithContext(ctx)
	dbreq := make([]dahanghaiModel.DaHangHaiRedis2, 0)
	for _, v := range req {
		dbreqItem := dahanghaiModel.DaHangHaiRedis2{}
		dbreqItem.Id = v.Id
		dbreqItem.Uid = v.Uid
		dbreqItem.TargetId = v.TargetId
		dbreqItem.PrivilegeType = v.PrivilegeType
		dbreqItem.StartTime = v.StartTime
		dbreqItem.ExpiredTime = v.ExpiredTime
		dbreqItem.Ctime = v.Ctime
		dbreqItem.Utime = v.Utime
		dbreq = append(dbreq, dbreqItem)
	}
	f := func(c context.Context) {
		if err := s.dao.SetDHHListCache(c, dbreq, uid); err != nil {
			log.Error(_errorServiceLogPrefix+"|asyncSetDHHCache|error(%v),missedUIDs(%v)", err)
		}
	}
	if runErr := s.async.Save(func() {
		f(c)
	}); runErr != nil {
		log.Error(_errorServiceLogPrefix+"|asyncSetDHHCache|error(%v),run cache is full(%v)", runErr)
		f(c)
	}
	return nil
}

func (s *GuardService) asyncSetDHHCacheBatch(ctx context.Context, req map[int64][]*dahanghaiModel.DaHangHaiRedis2) error {
	c := metadata.WithContext(ctx)
	dbreq := make(map[int64][]dahanghaiModel.DaHangHaiRedis2)
	dbreqList := make([]dahanghaiModel.DaHangHaiRedis2, 0)
	for k, v := range req {

		for _, vv := range v {
			dbreqItem := dahanghaiModel.DaHangHaiRedis2{}
			dbreqItem.Id = vv.Id
			dbreqItem.Uid = vv.Uid
			dbreqItem.TargetId = vv.TargetId
			dbreqItem.PrivilegeType = vv.PrivilegeType
			dbreqItem.StartTime = vv.StartTime
			dbreqItem.ExpiredTime = vv.ExpiredTime
			dbreqItem.Ctime = vv.Ctime
			dbreqItem.Utime = vv.Utime
			dbreqList = append(dbreqList, dbreqItem)
		}
		if _, exist := dbreq[k]; !exist {
			dbreq[k] = make([]dahanghaiModel.DaHangHaiRedis2, 0)
		}
		dbreq[k] = dbreqList
	}
	f := func(c context.Context) {
		for k, v := range dbreq {
			if err := s.dao.SetDHHListCache(c, v, k); err != nil {
				log.Error(_errorServiceLogPrefix+"|asyncSetDHHCache|error(%v),missedUIDs(%v)", err)
			}
		}
	}
	if runErr := s.asyncMulti.Save(func() {
		f(c)
	}); runErr != nil {
		log.Error(_errorServiceLogPrefix+"|asyncSetDHHCache|error(%v),run cache is full(%v)", runErr)
		f(c)
	}
	return nil
}

func (s *GuardService) asyncSetAnchorGuardCache(ctx context.Context, req []*dahanghaiModel.DaHangHaiRedis2, uid int64) error {
	c := metadata.WithContext(ctx)
	dbreq := make([]dahanghaiModel.DaHangHaiRedis2, 0)
	for _, v := range req {
		dbreqItem := dahanghaiModel.DaHangHaiRedis2{}
		dbreqItem.Id = v.Id
		dbreqItem.Uid = v.Uid
		dbreqItem.TargetId = v.TargetId
		dbreqItem.PrivilegeType = v.PrivilegeType
		dbreqItem.StartTime = v.StartTime
		dbreqItem.ExpiredTime = v.ExpiredTime
		dbreqItem.Ctime = v.Ctime
		dbreqItem.Utime = v.Utime
		dbreq = append(dbreq, dbreqItem)
	}
	f := func(c context.Context) {
		if err := s.dao.SetAnchorGuardListCache(c, dbreq, uid); err != nil {
			log.Error(_errorServiceLogPrefix+"|asyncSetAnchorGuardCache|error(%v),missedUIDs(%v)", err)
		}
	}
	if runErr := s.async.Save(func() {
		f(c)
	}); runErr != nil {
		log.Error(_errorServiceLogPrefix+"|asyncSetAnchorGuardCache|error(%v),run cache is full(%v)", runErr)
		f(c)
	}
	return nil
}
