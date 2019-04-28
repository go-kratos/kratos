package service

import (
	"context"

	pb "go-common/app/service/main/resource/api/v1"
	"go-common/app/service/main/resource/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

//Relate Relate card grpc
func (s *Service) Relate(ctx context.Context, req *pb.RelateRequest) (special *pb.SpecialReply, err error) {
	special = &pb.SpecialReply{}
	if req == nil || req.Id == 0 || req.MobiApp == "" || req.Build == 0 {
		err = ecode.RequestErr
		return
	}
	var (
		relateID int64
		ok       bool
		relate   *model.Relate
		versions []*model.Version
	)
	//判断seasonID 是否有配置的相关推荐卡片
	if relateID, ok = s.relatePgcMapCache[req.Id]; !ok {
		err = ecode.NothingFound
		log.Warn("gRpc.Relate relatePgcMapCache error,req.id(%v)", req.Id)
		return
	}
	//取出seasonID 对应的 相关推荐的卡片数据
	if relate, ok = s.relateCache[relateID]; !ok {
		err = ecode.NothingFound
		log.Warn("gRpc.Relate relateCache error,req.id(%v),relateID (%v)", req.Id, relateID)
		return
	}
	p := model.Plat(req.MobiApp, req.Device)
	//判断APP是否存在
	if versions, ok = relate.Versions[p]; !ok {
		err = ecode.NothingFound
		log.Warn("gRpc.Relate relate.Versions error,req.id(%v),plat (%v)", req.Id, p)
		return
	}
	//判断APP版本是否存在
	if len(versions) == 0 {
		err = ecode.NothingFound
		log.Warn("gRpc.Relate versions error,Versions is zero,req.id(%v)", req.Id)
		return
	}
	//判断版本信息是否匹配
	for _, v := range versions {
		if model.InvalidBuild(int(req.Build), v.Build, v.Condition) {
			err = ecode.NothingFound
			log.Warn("gRpc.Relate InvalidBuild error,req.id(%v),req.Build (%v)", req.Id, req.Build)
			return
		}
	}
	var specialTmp *pb.SpecialReply
	if specialTmp, ok = s.specialCache[relate.Param]; ok && specialTmp != nil {
		*special = *specialTmp
		special.Position = relate.Position
	} else {
		special = &pb.SpecialReply{}
		log.Warn("gRpc.Relate specialCache error,req.id(%v),relate.Param (%v)", req.Id, relate.Param)
	}
	return
}

//loadSpecialCache load special card cache
func (s *Service) loadSpecialCache() {
	special, err := s.manager.Specials(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.specialCache = special
}
