package service

import (
	"context"
	"fmt"

	pb "go-common/app/service/main/history/api/grpc"
	"go-common/app/service/main/history/model"
	"go-common/library/stat/prom"

	"golang.org/x/sync/singleflight"
)

var cacheSingleFlight = &singleflight.Group{}

// UserHide 查询是否记录播放历史
func (s *Service) UserHide(c context.Context, arg *pb.UserHideReq) (reply *pb.UserHideReply, err error) {
	reply = &pb.UserHideReply{}
	addCache := true
	var value int64
	value, err = s.dao.UserHideCache(c, arg.Mid)
	reply.Hide = value == model.HideStateON
	if err != nil {
		addCache = false
		err = nil
	}
	if value != model.HideStateNotFound {
		prom.CacheHit.Incr("UserHide")
		return
	}
	var rr interface{}
	sf := fmt.Sprintf("sf_u%d", arg.Mid)
	rr, err, _ = cacheSingleFlight.Do(sf, func() (r interface{}, e error) {
		prom.CacheMiss.Incr("UserHide")
		r, e = s.dao.UserHide(c, arg.Mid)
		return
	})
	reply.Hide = rr.(bool)
	if err != nil || !addCache {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.SetUserHideCache(ctx, arg.Mid, hideToState(reply.Hide))
	})
	return
}

// UpdateUserHide 修改是否记录播放历史
func (s *Service) UpdateUserHide(c context.Context, arg *pb.UpdateUserHideReq) (reply *pb.UpdateUserHideReply, err error) {
	reply = &pb.UpdateUserHideReply{}
	if err = s.dao.UpdateUserHide(c, arg.Mid, arg.Hide); err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.SetUserHideCache(ctx, arg.Mid, hideToState(arg.Hide))
	})
	return
}

func hideToState(hide bool) int64 {
	if hide {
		return model.HideStateON
	}
	return model.HideStateOFF
}
