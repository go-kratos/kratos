package service

import (
	"context"

	"go-common/library/log"
)

// BaseInfo get user's base info.
// func (s *Service) BaseInfo(c context.Context, mid int64) (info *model.BaseInfo, err error) {
// 	if mid <= 0 {
// 		log.Info("s.BaseInfo(%d) mid not valid number!", mid)
// 		return
// 	}
// 	var mc = true
// 	if info, err = s.dao.BaseInfoCache(c, mid); err != nil {
// 		mc = false
// 		err = nil // ignore error
// 	}
// 	if info != nil {
// 		if info.Mid == 0 {
// 			log.Info("s.BaseInfo(%d) mid not exist!", mid)
// 		}
// 		return
// 	}
// 	if info, err = s.dao.BaseInfo(c, mid); err != nil {
// 		log.Error("s.dao.BaseInfo(%d) error(%v)", mid, err)
// 		return
// 	}
// 	if info == nil {
// 		info = &model.BaseInfo{}
// 		log.Info("s.BaseInfo(%d) mid not exist!", mid)
// 		return
// 	}
// 	if mc {
// 		s.dao.SetBaseInfoCache(context.TODO(), mid, info)
// 	}
// 	log.Info("s.BaseInfo(%d) info(%+v)", mid, info)
// 	return
// }

func (s *Service) updateAccFace(c context.Context, mid int64) error {
	base, err := s.dao.BaseInfo(c, mid)
	if err != nil {
		log.Error("updateAccFace s.dao.BaseInfoWithoutDomain(%d) error(%v)", mid, err)
		return err
	}
	if base == nil {
		return nil
	}
	return s.dao.UpdateAccFace(c, mid, base.Face)
}
