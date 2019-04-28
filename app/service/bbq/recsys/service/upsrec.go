package service

import (
	"context"
	"fmt"
	"github.com/json-iterator/go"
	rpc "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/library/log"
	"go-common/library/net/trace"
)

//UpsRecService
func (s *Service) UpsRecService(c context.Context, req *rpc.RecsysRequest) (response *rpc.RecsysResponse, err error) {
	//请求日志
	data1, err := jsoniter.Marshal(req)
	if err == nil {
		log.Info("upsrec request is %s:", data1)
	}

	// 0.0 pre process: ab test
	// 0.0 pre process
	tracer, _ := trace.FromContext(c)
	req.TraceID = fmt.Sprintf("%s", tracer)

	// 1.0 get user profile
	userProfile, err := s.dao.LoadUserProfile(c, req.MID, req.BUVID)

	if req.MID != 0 {
		if err = s.dao.GetUserFollow(c, req.MID, userProfile); err != nil {
			log.Errorv(c, log.KV("userLog", "query user follow fail"), log.KV("MID", req.MID))
			err = nil
		}
		if err = s.dao.GetUserBlack(c, req.MID, userProfile); err != nil {
			log.Errorv(c, log.KV("userLog", "query user black fail"), log.KV("MID", req.MID))
			err = nil
		}
	}

	//2.0 retrieve
	response, err = s.recallManager.UpsRec(c, req, userProfile, s.dao.RecallClient, s.dao.RelationClient)
	//if err != nil {
	//
	//}

	// 3.0 filter
	s.filterManager.upsFilter(req, response, userProfile)

	// 4.0 ranker

	// 4.0.0
	// 4.0.1 prepare feature
	// 4.0.2 do rank
	// todo rank model
	s.rankManager.rank(c, req, response, userProfile, s.dao)

	// 5.0 post process, apply rule, page, store results
	err = s.postProcessor.ProcessUpsRec(c, req, response, userProfile)

	size := len(response.List)
	if size == 0 {
		log.Error("Ups response is empty!")
		response = &rpc.RecsysResponse{
			Message: make(map[string]string),
			List:    make([]*rpc.RecsysRecord, 0),
		}
		response.Message["info"] = "Ups response is empty!"
		return
	}

	//debug log
	if req.DebugFlag {
		data, _ := jsoniter.Marshal(userProfile)
		response.Message["UserInfo"] = string(data)
		return
	}

	// 5.2 page
	limit := int(req.Limit)
	if limit > size {
		limit = size
	}
	response.List = response.List[0:limit]

	// 5.3 store results
	s.dao.StoreRecResults(c, userProfile, req.MID, req.BUVID, response, s.dao.LastUpsPageRedisKey, userProfile.LastUpsRecords)

	s.StoreLog(req, response, userProfile, "upsrec")

	return
}
