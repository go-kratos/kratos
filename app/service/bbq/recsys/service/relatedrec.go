package service

import (
	"context"
	"fmt"
	rpc "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/app/service/bbq/recsys/service/retrieve"
	"go-common/app/service/bbq/recsys/service/util"
	"go-common/library/log"
	"go-common/library/net/trace"
	"math"
	"sort"
	"strconv"
	"strings"
)

//RelatedRecService ...
func (s *Service) RelatedRecService(c context.Context, req *rpc.RecsysRequest) (response *rpc.RecsysResponse, err error) {

	// 0.0 pre process: ab test
	// 0.0 pre process
	tracer, _ := trace.FromContext(c)
	req.TraceID = fmt.Sprintf("%s", tracer)

	// 1.0 get user profile
	userProfile := &model.UserProfile{Mid: req.MID, Buvid: req.BUVID}
	//if userProfile, err = s.dao.GetUserProfile(c, req.MID, req.BUVID); err != nil {
	//	log.Warn("get user profile failed, mid: ", req.MID)
	//}

	// 2.0 query rewrite, retrieve

	response, err = s.recallManager.RelatedRec(c, req, s.dao.RecallClient)

	// 3.0 filter
	s.filterManager.relatedFilter(req, response, userProfile)

	//fmt.Println("response size after filter:", len(response.List))

	// 4.0 ranker

	// 4.0.0
	// 4.0.1 prepare feature
	// 4.0.2 do rank

	rankRelated(response)

	// 5.0 post process, apply rule, page, store results

	//s.postProcessor.process(c, response)

	size := len(response.List)
	if size == 0 {
		log.Error("Related_response is empty!")
		response = &rpc.RecsysResponse{
			Message: make(map[string]string),
			List:    make([]*rpc.RecsysRecord, 0),
		}
		response.Message["info"] = "Related_response is empty!"
		return
	}

	// 5.2 page
	limit := int(req.Limit)
	if limit > size {
		limit = size
	}
	response.List = response.List[0:limit]

	s.StoreLog(req, response, userProfile, "relatedrec")

	return
}

func rankRelated(response *rpc.RecsysResponse) {

	sourceZoneID := response.Message[retrieve.SourceZoneID]
	sourceTagIDs := strings.Split(response.Message[retrieve.SourceTagIDs], "|")

	for _, record := range response.List {
		retriever := record.Map[model.Retriever]
		if retriever == retrieve.I2iRecall {
			record.Score = 1.2
		} else if retriever == retrieve.I2tag2iRecall {
			record.Score = 1.0
		} else if retriever == retrieve.I2tag2iRecall {
			record.Score = 0.5
		} else if retriever == retrieve.HotRecall {
			record.Score = 0
		}

		zoneScore := 0.0
		tagCount := 0
		tagCommonCount := 0
		if itemTagIDs, ok := record.Map[model.TagsID]; ok {
			tagIDs := strings.Split(itemTagIDs, "|")
			tagCount = len(tagIDs)
			for _, tagID := range tagIDs {
				if tagID == sourceZoneID {
					zoneScore = 0.5
				}
				for _, sourceTagID := range sourceTagIDs {
					if tagID == sourceTagID {
						tagCommonCount++
					}
				}
			}
		}
		record.Score += zoneScore

		tagScore := 0.0
		if tagCommonCount > 0 {
			tagScore += 0.5
		}
		tagScore += 0.5 * (float64(tagCommonCount) + 1) / (float64(tagCount) + 1)

		stateStr := record.Map[model.State]
		state, _ := strconv.ParseInt(stateStr, 10, 64)
		if state == model.State5 {
			record.Score += 0.3
		} else if state == model.State4 {
			record.Score += 0.1
		}

		if play, ok := record.Map[model.PlayHive]; ok {
			playNum, _ := strconv.ParseFloat(play, 64)
			playNumScore := math.Log10(math.Min(playNum+1.0, 1000000.0)) / math.Log10(1000000.0)
			record.Score += playNumScore

			if fav, ok := record.Map[model.FavHive]; ok {
				favNum, _ := strconv.ParseFloat(fav, 64)
				favScore := (math.Min(favNum, playNum) + 1.0) / (playNum + 200.0)
				favScore = math.Min(favScore, 0.1)
				record.Score += 5 * favScore
			}

			if likes, ok := record.Map[model.LikesHive]; ok {
				likesNum, _ := strconv.ParseFloat(likes, 64)
				likesScore := (math.Min(likesNum, playNum) + 1.0) / (playNum + 200.0)
				likesScore = math.Min(likesScore, 0.1)
				record.Score += 5 * likesScore
			}

			if shares, ok := record.Map[model.ShareHive]; ok {
				shareNum, _ := strconv.ParseFloat(shares, 64)
				shareScore := (math.Min(shareNum, playNum) + 1.0) / (playNum + 200.0)
				shareScore = math.Min(shareScore, 0.1)
				record.Score += 5 * shareScore
			}
		}

		// bbq video feature
		if play, ok := record.Map[model.PlayMonth]; ok {
			playNum, _ := strconv.ParseFloat(play, 64)
			playNumScore := math.Log10(math.Min(playNum+1.0, 1000000.0)) / math.Log10(1000000.0)
			record.Score += playNumScore

			if likes, ok := record.Map[model.LikesMonth]; ok {
				likesNum, _ := strconv.ParseFloat(likes, 64)
				likesScore := (math.Min(likesNum, playNum) + 1.0) / (playNum + 200.0)
				likesScore = math.Min(likesScore, 0.1)
				record.Score += 5 * likesScore
			}

			if share, ok := record.Map[model.ShareMonth]; ok {
				shareNum, _ := strconv.ParseFloat(share, 64)
				shareScore := (math.Min(shareNum, playNum) + 1.0) / (playNum + 200.0)
				shareScore = math.Min(shareScore, 0.1)
				record.Score += 10 * shareScore
			}

			if reply, ok := record.Map[model.ReplyMonth]; ok {
				replyNum, _ := strconv.ParseFloat(reply, 64)
				replyScore := (math.Min(replyNum, playNum) + 1.0) / (playNum + 200.0)
				replyScore = math.Min(replyScore, 0.1)
				record.Score += 10 * replyScore
			}
		}
	}

	sort.Sort(sort.Reverse(util.Records(response.List)))
}
