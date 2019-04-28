package service

import (
	"context"

	"go-common/app/service/main/msm/model"
	"go-common/library/ecode"
)

// Limits get limits.
func (s *Service) Limits(c context.Context, family, hmd5 string) (*model.Limits, error) {
	// var (
	// 	id     int64
	// 	bytes  []byte
	// 	num    int
	// 	err    error
	// 	limits map[string]*model.Limit
	// )
	// if num, err = s.dao.HostAmount(c, family); err != nil {
	// 	log.Error("HostAmount() error(%v)", err)
	// 	return nil, err
	// }
	// if id, err = s.dao.LimitID(c, family); err != nil {
	// 	log.Error("LimitID() error(%v)", err)
	// 	return nil, err
	// }
	// if limits, err = s.dao.LimitBusiness(c, id, num); err != nil {
	// 	log.Error("LimitBusiness() error(%v)", err)
	// 	return nil, err
	// }
	// if bytes, err = json.Marshal(limits); err != nil {
	// 	log.Error("json.Marshal(%v) error(%v)", limits, err)
	// 	return nil, err
	// }
	// mb := md5.Sum(bytes)
	// if md5 := hex.EncodeToString(mb[:]); md5 != hmd5 {
	// 	return &model.Limits{Apps: limits, MD5: md5}, nil
	// }
	return nil, ecode.NotModified
}
