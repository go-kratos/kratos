package service

import (
	"context"
	"fmt"
	rpc "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/app/service/bbq/recsys/service/util"
	"go-common/library/log"
	"go-common/library/net/trace"
	"strconv"

	"strings"

	"github.com/json-iterator/go"
)

//Start this just a example
func (s *Service) Start(c context.Context, req *rpc.RecsysRequest) (res *rpc.RecsysResponse, err error) {
	return s.RecService(c, req)
}

//RecService recommend service
func (s *Service) RecService(c context.Context, req *rpc.RecsysRequest) (response *rpc.RecsysResponse, err error) {

	//请求日志
	data1, err := jsoniter.Marshal(req)
	if err == nil {
		log.Info("recsys request is %s:", data1)
	}

	// 0.0 pre process
	tracer, _ := trace.FromContext(c)
	req.TraceID = fmt.Sprintf("%s", tracer)

	response = new(rpc.RecsysResponse)
	response.Message = make(map[string]string)

	// 0.1 ab test
	s.DoABTest(req)

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

	// 2.0 query rewrite, parallel retrieve
	response, err = s.recallManager.V2RetrieveFunc(c, req, userProfile, s.dao.RecallClient)
	//is or not debug
	if req.DebugFlag {
		recallStatCountMap := make(map[string]int)
		recallTagStatCountMap := make(map[string]int)

		for index, record := range response.List {
			record.Map[model.OrderRecall] = strconv.Itoa(index)

			recallClasses := record.Map[model.RecallClasses]
			for _, recallClass := range strings.Split(recallClasses, "|") {
				recallStatCountMap[recallClass] = recallStatCountMap[recallClass] + 1
			}

			recallTags := record.Map[model.RecallTags]
			for _, recallTag := range strings.Split(recallTags, "|") {
				recallTagStatCountMap[recallTag] = recallTagStatCountMap[recallTag] + 1
			}

		}
		response.Message["DebugStatRecallCntTotal"] = strconv.Itoa(len(response.List))

		recallStatCountList := util.SortStrIntMapByValue(recallStatCountMap)
		recallCountStr, _ := jsoniter.MarshalToString(recallStatCountList)
		response.Message["DebugStatRecallCntDetail"] = recallCountStr

		recallTagStatCountList := util.SortStrIntMapByValue(recallTagStatCountMap)
		recallTagCountStr, _ := jsoniter.MarshalToString(recallTagStatCountList)
		response.Message["DebugStatRecallTagCntDetail"] = recallTagCountStr
	}

	//2.1 down grade recall
	if err != nil || len(response.List) == 0 {
		response, err = s.dao.DownGradeRecall(c)
	}

	// 3.0 merge && filter
	s.filterManager.filter(req, response, userProfile)

	// 4.0 ranker

	// 4.0.0
	s.businessInfoCount.State(model.ResponseCount, int64(len(response.List)))
	response.Message[model.ResponseCount] = strconv.Itoa(len(response.List))

	// 4.0.1 prepare feature
	// 4.0.2 do rank
	if req.Abtest == ABTestA || req.MID == 5829468 {
		err = s.rankModelManager.DoRank(req, response, userProfile)
		if err != nil {
			log.Error("rank model failed (%v)", err)
			s.rankManager.rank(c, req, response, userProfile, s.dao)
			err = nil
		}
	} else {
		s.rankManager.rank(c, req, response, userProfile, s.dao)
	}

	// 5.0 post process, apply rule, page, store results

	// 5.1 post process
	err = s.postProcessor.ProcessRec(c, req, response, userProfile)

	size := len(response.List)
	if size == 0 {
		log.Error("response is empty! request is (%v)", req)
		response = &rpc.RecsysResponse{
			Message: make(map[string]string),
			List:    make([]*rpc.RecsysRecord, 0),
		}
		response.Message["info"] = "response is empty!"
		return
	}

	for index, record := range response.List {
		record.Map[model.OrderFinal] = strconv.Itoa(index)
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
	s.dao.StoreRecResults(c, userProfile, req.MID, req.BUVID, response, s.dao.LastPageRedisKey, userProfile.LastRecords)

	// 5.4 store feature log && reduce record keys
	s.StoreLog(req, response, userProfile, "bbq-recsys")

	return
}
