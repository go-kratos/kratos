package retrieve

import (
	"context"
	"fmt"
	"go-common/app/service/main/relation/api"
	"go-common/library/log"

	recallv1 "go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
)

const BiliFollowsRecall = "BiliFollowsRecall"

func (m *RecallManager) UpsRec(c context.Context, request *recsys.RecsysRequest, userProfile *model.UserProfile, recallClient recallv1.RecsysRecallClient, relationClient api.RelationClient) (response *recsys.RecsysResponse, err error) {
	recallInfos := make([]*recallv1.RecallInfo, 0)
	// generate recall infos

	// bbq user profile: bili followups,highest priority
	// x * 50
	biliFollowings := getFollowings(c, relationClient, userProfile.Mid)
	bbqFollowings := userProfile.BBQFollow
	notFollowedBiliUpsCount := 0
	for _, upMID := range biliFollowings { //for every upMID in bili test if its in bbq
		if _, ok := bbqFollowings[upMID]; !ok {
			notFollowedBiliUpsCount++
			tagRecallInfo := &recallv1.RecallInfo{
				Name:     BiliFollowsRecall,
				Tag:      fmt.Sprintf(_recRecallKeyUP, upMID),
				Limit:    50,
				Filter:   "",
				Priority: _PriorityVeryHigh,
			}
			recallInfos = append(recallInfos, tagRecallInfo)
		}
		if notFollowedBiliUpsCount >= 48 {
			break
		}
	}

	// selection 100
	selectionRecallInfo := &recallv1.RecallInfo{
		Name:     SelectionRecall,
		Tag:      RecallOpVideoKey,
		Limit:    100,
		Filter:   "",
		Priority: _PriorityHigh,
		Scorer:   _RandomScorer,
	}
	recallInfos = append(recallInfos, selectionRecallInfo)

	// 热门
	hotRecallInfo := &recallv1.RecallInfo{
		Name:     HotRecall,
		Tag:      RecallHotDefault,
		Limit:    200,
		Filter:   "",
		Priority: _PriorityHigh,
		Scorer:   _RandomScorer,
	}
	recallInfos = append(recallInfos, hotRecallInfo)

	//recall request
	recallRequest := &recallv1.RecallRequest{
		MID:        request.MID,
		BUVID:      request.BUVID,
		TotalLimit: 500,
		Infos:      recallInfos,
	}
	log.Info("recall key count is (%v), recall request: (%v)", len(recallInfos), recallRequest)

	response = new(recsys.RecsysResponse)
	response.Message = make(map[string]string)

	// do real request action
	recallResponse, err := recallClient.Recall(c, recallRequest)

	//todo down grade recall
	//todo limit result from bili followups to 100

	if err != nil || recallResponse == nil {
		log.Error("recall service error (%v) or recall response is null, traceID is %v", err, request.TraceID)
		return
	} else if recallResponse.List == nil {
		log.Warn("recall service did not return any result")
		response.List = make([]*recsys.RecsysRecord, 0)
		return
	}
	err = transform(recallResponse, response)
	return
}

func getFollowings(c context.Context, relationClient api.RelationClient, mid int64) (followings []int64) {
	followings = make([]int64, 0)
	midReq := api.MidReq{Mid: mid}
	followingsReply, err := relationClient.Followings(c, &midReq)
	if err != nil {
		if followingsReply.FollowingList != nil {
			for _, followingReply := range followingsReply.FollowingList {
				mid := followingReply.Mid
				followings = append(followings, mid)
			}
		}
	}
	return
}
