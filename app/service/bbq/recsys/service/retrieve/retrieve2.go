package retrieve

import (
	"context"
	"fmt"
	recallv1 "go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/app/service/bbq/recsys/service/util"
	"go-common/library/log"
	"strconv"
)

//召回策略
const (
	//recall class
	HotRecall       = "HotRecall"
	RandomRecall    = "RandomRecall"
	SelectionRecall = "SelectionRecall"
	UserProfileBili = "UserProfileBili"
	UserProfileBBQ  = "UserProfileBBQ"
	LikeI2IRecall   = "LikeI2IRecall"
	LikeTagRecall   = "LikeTagRecall"
	LikeUPRecall    = "LikeUPRecall"
	PosI2IRecall    = "PosI2IRecall"
	PosTagRecall    = "PosTagRecall"
	FollowRecall    = "FollowRecall"

	_MaxRecallTagCount = 50

	//priority
	_PriorityVeryHigh = 10000
	_PriorityHigh     = 1000
	_PriorityMid      = 100
	_PriorityLow      = 10

	//recall rank method
	_RandomScorer = "default"

	_TopNLikeVideo = 10
	_TopNLikeUp    = 10
	_TopNLikeTag   = 10
	_TopNPosVideo  = 10
	_TopNPosTag    = 10
	_TopNFollow    = 10
	_TopNNegTag    = 20 // not used as recall tag
)

//RecallKey Prefixes
const (
	RecallKeyI2IPrefix   = "RECALL:I2I"
	RecallKeyTagIDPrefix = "RECALL:HOT_T"
	RecallKeyUpIDPrefix  = "RECALL:HOT_UP"
)

//Recall redis key
const (
	RecallOpVideoKey                       = "job:bbq:rec:op"
	RecallHotDefault                       = "RECALL:HOT_DEFAULT:0"
	_recRecallKeyI2I                       = "RECALL:I2I:%d"
	_recRecallKeyUP                        = "RECALL:HOT_UP:%d"
	_recRecallKeyTagIDTemplateString       = "RECALL:HOT_T:%s"
	_recRecallKeyTagIDTemplateInt          = "RECALL:HOT_T:%d"
	_recRecallKeyTagIDNewPubTemplateString = "RECALL:T:%s"

	_BloomFilter = "bloomfilter"
)

// RecallManager manages multiple retrieve functions
type RecallManager struct {
	Retrievers []Retriever2
}

//Retriever2 ...
type Retriever2 interface {
	name() (name string)

	queryRewrite(c context.Context, request *recsys.RecsysRequest, userProfile *model.UserProfile) (recallInfo *recallv1.RecallInfo, err error)
}

//RetrieverFuncV2 ...
type RetrieverFuncV2 func(c context.Context, request *recsys.RecsysRequest, userProfile *model.UserProfile, recallClient recallv1.RecsysRecallClient) (response *recsys.RecsysResponse, err error)

//NewRecallManager ...
func NewRecallManager() (m *RecallManager) {
	m = &RecallManager{
		Retrievers: make([]Retriever2, 0),
	}
	return
}

//V2RetrieveFunc is default retrieve function
func (m *RecallManager) V2RetrieveFunc(c context.Context, request *recsys.RecsysRequest, userProfile *model.UserProfile, recallClient recallv1.RecsysRecallClient) (response *recsys.RecsysResponse, err error) {

	recallInfos := make([]*recallv1.RecallInfo, 0)
	for _, retriever := range m.Retrievers {
		recallInfo, err := retriever.queryRewrite(c, request, userProfile)
		if err == nil && recallInfo != nil {
			recallInfos = append(recallInfos, recallInfo)
		}
	}

	// selection 100
	selectionRecallInfo := &recallv1.RecallInfo{
		Name:     SelectionRecall,
		Tag:      RecallOpVideoKey,
		Limit:    100,
		Filter:   _BloomFilter,
		Priority: _PriorityHigh,
		Scorer:   _RandomScorer,
	}
	recallInfos = append(recallInfos, selectionRecallInfo)

	// 热门
	hotRecallInfo := &recallv1.RecallInfo{
		Name:     HotRecall,
		Tag:      RecallHotDefault,
		Limit:    200,
		Filter:   _BloomFilter,
		Priority: _PriorityHigh,
		Scorer:   _RandomScorer,
	}
	recallInfos = append(recallInfos, hotRecallInfo)

	// tags
	likeVideoMap, likeUPMap, likeTagIDMap, posVideoMap, posTagIDMap := buildSessionFeature(c, userProfile, recallClient)
	if request.DebugFlag {
		log.Info("posVideoMap = %v", posVideoMap)
	}

	// bbq实时行为 bbq like i2i
	for ID := range likeVideoMap {
		recallInfo := &recallv1.RecallInfo{
			Name:     LikeI2IRecall,
			Tag:      fmt.Sprintf(_recRecallKeyI2I, ID),
			Limit:    40,
			Filter:   _BloomFilter,
			Priority: _PriorityVeryHigh,
		}
		recallInfos = append(recallInfos, recallInfo)
	}
	// bbq实时行为 bbq like up
	for UpMID := range likeUPMap {
		recallInfo := &recallv1.RecallInfo{
			Name:     LikeUPRecall,
			Tag:      fmt.Sprintf(_recRecallKeyUP, UpMID),
			Limit:    40,
			Filter:   _BloomFilter,
			Priority: _PriorityVeryHigh,
		}
		recallInfos = append(recallInfos, recallInfo)
	}

	// bbq实时行为 bbq follow
	// x * 5
	var followMap map[int64]int64
	if len(userProfile.BBQFollowAction) > _TopNFollow {
		followUps := util.SortMapByValue(userProfile.BBQFollowAction)[0:_TopNFollow]
		followMap = make(map[int64]int64, _TopNFollow)
		for _, pair := range followUps {
			followMap[pair.Key] = pair.Value
		}
	} else {
		followMap = userProfile.BBQFollowAction
	}
	for upMID := range followMap {
		tagRecallInfo := &recallv1.RecallInfo{
			Name:     FollowRecall,
			Tag:      fmt.Sprintf(_recRecallKeyUP, upMID),
			Limit:    40,
			Filter:   _BloomFilter,
			Priority: _PriorityHigh,
		}
		recallInfos = append(recallInfos, tagRecallInfo)
	}

	for tag := range likeTagIDMap {
		tagRecallInfo := &recallv1.RecallInfo{
			Name:     LikeTagRecall,
			Tag:      fmt.Sprintf(_recRecallKeyTagIDTemplateInt, tag),
			Limit:    20,
			Filter:   _BloomFilter,
			Priority: _PriorityHigh,
		}
		recallInfos = append(recallInfos, tagRecallInfo)
	}

	for tag := range posTagIDMap {
		tagRecallInfo := &recallv1.RecallInfo{
			Name:     PosTagRecall,
			Tag:      fmt.Sprintf(_recRecallKeyTagIDTemplateInt, tag),
			Limit:    30,
			Filter:   _BloomFilter,
			Priority: _PriorityHigh,
		}
		recallInfos = append(recallInfos, tagRecallInfo)
	}

	//for tag := range userProfile.PosTags {
	//	tagRecallInfo := &recallv1.RecallInfo{
	//		Name:     LikeTagRecall,
	//		Tag:      fmt.Sprintf(_recRecallKeyTemplateTagID, tag),
	//		Limit:    10,
	//		Filter:   _BloomFilter,
	//		Priority: _PriorityHigh,
	//	}
	//	recallInfos = append(recallInfos, tagRecallInfo)
	//}

	// bbq user profile: zone
	// 20 * 5 = 100
	for tag := range userProfile.BBQZones {
		zoneRecallInfo := &recallv1.RecallInfo{
			Name:     UserProfileBBQ,
			Tag:      fmt.Sprintf(_recRecallKeyTagIDTemplateString, tag),
			Limit:    10,
			Filter:   _BloomFilter,
			Priority: _PriorityMid,
		}
		recallInfos = append(recallInfos, zoneRecallInfo)
	}
	// bbq user profile: tag
	// 30 * 5 = 150
	for tag := range userProfile.BBQTags {
		tagRecallInfo := &recallv1.RecallInfo{
			Name:     UserProfileBBQ,
			Tag:      fmt.Sprintf(_recRecallKeyTagIDTemplateString, tag),
			Limit:    20,
			Filter:   _BloomFilter,
			Priority: _PriorityMid,
		}
		recallInfos = append(recallInfos, tagRecallInfo)
	}
	// bbq user profile: follow
	// x * 5
	for upMID := range userProfile.BBQPrefUps {
		tagRecallInfo := &recallv1.RecallInfo{
			Name:     UserProfileBBQ,
			Tag:      fmt.Sprintf(_recRecallKeyUP, upMID),
			Limit:    10,
			Filter:   _BloomFilter,
			Priority: _PriorityMid,
		}
		recallInfos = append(recallInfos, tagRecallInfo)
	}

	// bili user profile: zone2
	// 20 * 5 = 100
	for tag := range userProfile.Zones2 {
		zoneRecallInfo := &recallv1.RecallInfo{
			Name:     UserProfileBili,
			Tag:      fmt.Sprintf(_recRecallKeyTagIDTemplateString, tag),
			Limit:    10,
			Filter:   _BloomFilter,
			Priority: _PriorityMid,
		}
		recallInfos = append(recallInfos, zoneRecallInfo)
	}
	// bili user profile: tag
	// 30 * 5 = 150
	tagRecallSize := 20
	if len(userProfile.BBQTags) > 0 {
		tagRecallSize = 10
	}
	for tag := range userProfile.BiliTags {
		tagRecallInfo := &recallv1.RecallInfo{
			Name:     UserProfileBili,
			Tag:      fmt.Sprintf(_recRecallKeyTagIDTemplateString, tag),
			Limit:    int32(tagRecallSize),
			Filter:   _BloomFilter,
			Priority: _PriorityMid,
		}
		recallInfos = append(recallInfos, tagRecallInfo)
	}

	// 召回标签较大的情况下，减少热门召回
	if len(recallInfos) >= 5 {
		for _, recallInfo := range recallInfos {
			if recallInfo.Name == HotRecall {
				recallInfo.Limit = 50
			}
		}
	}
	//recallTagNameMap := make(map[string][]string)
	//recallTagInfoMap := make(map[string]*recallv1.RecallInfo)
	//recallTagPriorityMap := make(map[string]int32)
	//
	//for _, recallInfo := range recallInfos {
	//	names := recallTagNameMap[recallInfo.Tag]
	//	names = append(names, recallInfo.Name)
	//	recallTagNameMap[recallInfo.Tag] = names
	//	recallTagInfoMap[recallInfo.Tag] = recallInfo
	//
	//	if priority, ok := recallTagPriorityMap[recallInfo.Tag]; ok {
	//		if recallInfo.Priority > priority {
	//			recallTagPriorityMap[recallInfo.Tag] = priority
	//		}
	//	} else {
	//		recallTagPriorityMap[recallInfo.Tag] = priority
	//	}
	//}
	//
	//newRecallInfos := make([]*recallv1.RecallInfo, 0)
	//for tag, names := range recallTagNameMap {
	//	recallInfo := recallTagInfoMap[tag]
	//	recallInfo.Name = strings.Join(names, "|")
	//	newRecallInfos = append(newRecallInfos, recallInfo)
	//}

	recallInfos = mergeRecallKey(recallInfos)
	recallKeyCount := len(recallInfos)

	// 老用户增加随机召回或者N刷之后??? TODO
	if len(userProfile.LastRecords) >= 20 && recallKeyCount < _MaxRecallTagCount && len(userProfile.BBQTags) > 0 && len(userProfile.BiliTags) > 0 {
		tagCountMap := make(map[string]int)
		for tag := range userProfile.BBQTags {
			count := tagCountMap[tag]
			tagCountMap[tag] = count + 1
		}
		for tag := range userProfile.BiliTags {
			count := tagCountMap[tag]
			tagCountMap[tag] = count + 1
		}

		randomTagCount := 0
		for tag := range tagCountMap {
			if randomTagCount > 10 {
				break
			}
			randomTagCount++
			tagRecallInfo := &recallv1.RecallInfo{
				Name:     RandomRecall,
				Tag:      fmt.Sprintf(_recRecallKeyTagIDNewPubTemplateString, tag),
				Limit:    10,
				Filter:   _BloomFilter,
				Priority: _PriorityLow,
				Scorer:   _RandomScorer,
			}
			recallInfos = append(recallInfos, tagRecallInfo)
		}
	}

	//merge random recall
	recallInfos = mergeRecallKey(recallInfos)

	recallRequest := &recallv1.RecallRequest{
		MID:        request.MID,
		BUVID:      request.BUVID,
		TotalLimit: 500,
		Infos:      recallInfos,
	}
	log.Info("recall key count is (%v), recall request: (%v)", len(recallInfos), recallRequest)

	response = new(recsys.RecsysResponse)
	response.Message = make(map[string]string)
	response.Message[model.QueryID] = request.QueryID

	recallResponse, err := recallClient.Recall(c, recallRequest)

	if err != nil || recallResponse == nil {
		log.Error("recall service error (%v) or recall response is null, traceID is %v", err, request.TraceID)
		return
	}

	if len(recallResponse.List) < int(request.Limit) {
		response.Message[model.ResponseDownGrade] = "1"
		log.Error("response size less then (%v), is (%v), , traceID is (%v)", request.Limit, recallResponse.Total, request.TraceID)
		recallResponse, err = downGradeRecall(c, recallRequest, recallClient)
		if err != nil || recallResponse == nil {
			log.Error("down grade recall service error (%v) or recall response is null, traceID is (%v)", err, request.TraceID)
			return
		}
	}

	response.Message[model.ResponseRecallCount] = strconv.Itoa(len(recallResponse.List))

	err = transform(recallResponse, response)

	deleteBlack(response, userProfile)

	return
}

func downGradeRecall(c context.Context, recallRequest *recallv1.RecallRequest, recallClient recallv1.RecsysRecallClient) (recallResponse *recallv1.RecallResponse, err error) {
	for _, recallInfo := range recallRequest.Infos {
		recallInfo.Filter = ""
		recallInfo.Scorer = _RandomScorer
	}
	log.Info("down grade recall request: (%v)", recallRequest)
	recallResponse, err = recallClient.Recall(c, recallRequest)
	return
}

func buildSessionFeature(c context.Context, userProfile *model.UserProfile, recallClient recallv1.RecsysRecallClient) (likeVideoMap map[int64]int64, likeUPMap map[int64]int64, likeTagIDMap map[int64]int64, posVideoMap map[int64]int64, posTagIDMap map[int64]int64) {

	if len(userProfile.LikeVideos) > 0 || len(userProfile.PosVideos) > 0 || len(userProfile.NegVideos) > 0 {
		SVIDs := make([]int64, 0)
		for SVID := range userProfile.LikeVideos {
			SVIDs = append(SVIDs, SVID)
		}
		for SVID := range userProfile.PosVideos {
			SVIDs = append(SVIDs, SVID)
		}
		for SVID := range userProfile.NegVideos {
			SVIDs = append(SVIDs, SVID)
		}
		videoIndexRequest := &recallv1.VideoIndexRequest{
			SVIDs: SVIDs,
		}
		videoIndexResponse, err := recallClient.VideoIndex(c, videoIndexRequest)
		if err != nil || videoIndexResponse == nil {
			log.Error("recall service VideoIndex error (%v) or recall response is null", err)
			return
		}
		for _, forwardIndex := range videoIndexResponse.List {
			SVIDInt := int64(forwardIndex.SVID)

			if timestamp, ok := userProfile.LikeVideos[SVIDInt]; ok {
				for _, tag := range forwardIndex.BasicInfo.Tags {
					tagID := int64(tag.TagID)
					if count, ok := userProfile.LikeTagIDs[tagID]; ok {
						userProfile.LikeTagIDs[tagID] = count + 1
					} else {
						userProfile.LikeTagIDs[tagID] = 1
					}
				}

				upMID := int64(forwardIndex.BasicInfo.MID)
				userProfile.LikeUPs[upMID] = timestamp
			}
		}

		for _, forwardIndex := range videoIndexResponse.List {
			SVIDInt := int64(forwardIndex.SVID)

			if _, ok := userProfile.PosVideos[SVIDInt]; ok {
				for _, tag := range forwardIndex.BasicInfo.Tags {
					tagID := int64(tag.TagID)
					if count, ok := userProfile.PosTagIDs[tagID]; ok {
						userProfile.PosTagIDs[tagID] = count + 1
					} else {
						userProfile.PosTagIDs[tagID] = 1
					}
				}
			}
		}

		for _, forwardIndex := range videoIndexResponse.List {
			SVIDInt := int64(forwardIndex.BasicInfo.SVID)

			if _, ok := userProfile.NegVideos[SVIDInt]; ok {
				for _, tag := range forwardIndex.BasicInfo.Tags {
					tagID := int64(tag.TagID)
					if count, ok := userProfile.NegTagIDs[tagID]; ok {
						userProfile.NegTagIDs[tagID] = count + 1
					} else {
						userProfile.NegTagIDs[tagID] = 1
					}
				}
			}
		}

		if len(userProfile.LikeVideos) > _TopNLikeVideo {
			likeVideos := util.SortMapByValue(userProfile.LikeVideos)[0:_TopNLikeVideo]
			likeVideoMap = make(map[int64]int64, _TopNLikeVideo)
			for _, pair := range likeVideos {
				likeVideoMap[pair.Key] = pair.Value
			}
		} else {
			likeVideoMap = userProfile.LikeVideos
		}

		if len(userProfile.LikeUPs) > _TopNLikeUp {
			likeUPs := util.SortMapByValue(userProfile.LikeUPs)[0:_TopNLikeUp]
			likeUPMap = make(map[int64]int64, _TopNLikeUp)
			for _, pair := range likeUPs {
				likeUPMap[pair.Key] = pair.Value
			}
		} else {
			likeUPMap = userProfile.LikeUPs
		}

		if len(userProfile.LikeTagIDs) > _TopNLikeTag {
			likeTags := util.SortMapByValue(userProfile.LikeTagIDs)[0:_TopNLikeTag]
			likeTagIDMap = make(map[int64]int64, _TopNLikeTag)
			for _, pair := range likeTags {
				likeTagIDMap[pair.Key] = pair.Value
			}
		} else {
			likeTagIDMap = userProfile.LikeTagIDs
		}

		if len(userProfile.PosVideos) > _TopNPosVideo {
			videos := util.SortMapByValue(userProfile.PosVideos)[0:_TopNPosVideo]
			posVideoMap = make(map[int64]int64, _TopNPosVideo)
			for _, pair := range videos {
				posVideoMap[pair.Key] = pair.Value
			}
		} else {
			posVideoMap = userProfile.PosVideos
		}

		//TODO 正负反馈标签考虑次数和发生时间
		if len(userProfile.PosTagIDs) > _TopNPosTag {
			posTags := util.SortMapByValue(userProfile.PosTagIDs)[0:_TopNPosTag]
			posTagIDMap = make(map[int64]int64, _TopNPosTag)
			for _, pair := range posTags {
				posTagIDMap[pair.Key] = pair.Value
			}
			userProfile.PosTagIDs = posTagIDMap
		} else {
			posTagIDMap = userProfile.PosTagIDs
		}

		if len(userProfile.NegTagIDs) > _TopNNegTag {
			tags := util.SortMapByValue(userProfile.NegTagIDs)[0:_TopNNegTag]
			tagIDMap := make(map[int64]int64, _TopNNegTag)
			for _, pair := range tags {
				tagIDMap[pair.Key] = pair.Value
			}
			userProfile.NegTagIDs = tagIDMap
		}
	}
	return
}
