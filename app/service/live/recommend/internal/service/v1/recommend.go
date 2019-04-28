package v1

import (
	"context"

	v1pb "go-common/app/service/live/recommend/api/grpc/v1"
	"go-common/app/service/live/recommend/internal/conf"
	"go-common/app/service/live/recommend/internal/dao"
	"go-common/library/log"
)

// RecommendService struct
type RecommendService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewRecommendService init
func NewRecommendService(c *conf.Config) (s *RecommendService) {
	s = &RecommendService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

// RandomRecsByUser implementation
// 随机获取n个推荐
func (s *RecommendService) RandomRecsByUser(ctx context.Context, req *v1pb.GetRandomRecReq) (resp *v1pb.GetRandomRecResp, err error) {
	resp = &v1pb.GetRandomRecResp{}

	ids, err := s.dao.GetRandomRoomIds(ctx, req.Uid, int(req.Count), req.ExistIds)
	if err != nil {
		log.Error("RandomRecsByUser err: %+v req:%+v", err, req)
		return
	}
	resp.RoomIds = ids
	resp.Count = uint32(len(ids))
	log.Info("RandomRecsByUser uid:%d, reqCount:%d count:%d roomIds: %v", req.Uid, req.Count, resp.Count, resp.RoomIds)
	return
}

// ClearRecommendCache implementation
// 清空推荐缓存，清空推荐过的集合
func (s *RecommendService) ClearRecommendCache(ctx context.Context, req *v1pb.ClearRecommendRequest) (resp *v1pb.ClearRecommendResponse, err error) {
	resp = &v1pb.ClearRecommendResponse{}
	err = s.dao.ClearRecommend(ctx, req.Uid)
	return
}
