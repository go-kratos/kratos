package service

import (
	"bytes"
	"context"

	"fmt"
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/dao"
	"go-common/app/service/bbq/recsys/model"
	"go-common/app/service/bbq/recsys/service/rank"
	"go-common/app/service/bbq/recsys/service/util"
	"math"
	"sort"
	"strconv"

	"github.com/json-iterator/go"
)

//RankManager ...
type RankManager struct {
	Rankers     []Ranker
	featureLogs []*rank.FeatureLog
}

//Ranker ...
type Ranker interface {
	name() (name string)

	rank(c context.Context, req *recsys.RecsysRequest, response *recsys.RecsysResponse, profile *model.UserProfile, dao *dao.Dao) (featureLogs []*rank.FeatureLog)
}

//NewRankManager ...
func NewRankManager() (m *RankManager) {
	m = &RankManager{
		Rankers:     make([]Ranker, 0),
		featureLogs: make([]*rank.FeatureLog, 0),
	}
	base := &BaseRanker{
		weights: initWeights(),
	}

	m.Rankers = append(m.Rankers, base)
	return
}

func initWeights() (weights map[string]float64) {

	weights = make(map[string]float64)

	//session feature
	weights[rank.SessionLikeTag] = 0.05
	weights[rank.LikeTagCount] = 0.05
	weights[rank.SessionPosPlayTag] = 0.01
	weights[rank.SessionNegPlayTag] = -0.04
	weights[rank.PureNegPlayTag] = -0.03
	weights[rank.LikeI2ITimeDiff] = 0.1
	weights[rank.FollowTimeDiff] = 0.05   //follow
	weights[rank.SessionBBQFollow] = 0.01 //follow

	//recall feature
	weights[rank.FollowRecall] = 0.01 //follow
	weights[rank.LikeI2IRecall] = 0.01
	weights[rank.LikeTagRecall] = 0.01
	weights[rank.PosI2IRecall] = 0.01
	weights[rank.PosTagRecall] = 0.0
	weights[rank.UserProfileBili] = 0.0
	weights[rank.UserProfileBBQ] = 0.0
	weights[rank.SelectionRecall] = 0.0
	weights[rank.BiliFollowsRecall] = 0.5
	weights[rank.HotRecall] = 0.0
	weights[rank.RandomRecall] = -0.01

	//user-item feature
	weights[rank.MatchBBQTagCountScore] = 0.02
	weights[rank.MatchBBQTagCountScore] = 0.02
	weights[rank.MatchBBQTagCount] = 0.0
	weights[rank.MatchBiliTagCount] = 0.0
	weights[rank.MatchBiliTagLevel3] = 0.015
	weights[rank.MatchBiliTagLevel2] = 0.01
	weights[rank.BiliPrefUp] = 0.005
	weights[rank.MatchBBQTagLevel3] = 0.015
	weights[rank.MatchBBQTagLevel2] = 0.01
	weights[rank.BBQPrefUp] = 0.005
	weights[rank.BBQFollow] = 0.005 //follow

	// item feature:
	weights[rank.OperationLevel] = 0.0

	// item feature: bili
	weights[rank.BiliPlayNum] = 0.05
	weights[rank.BiliFavRatio] = 0.2
	weights[rank.BiliLikeRatio] = 0.2
	weights[rank.BiliShareRatio] = 0.2
	weights[rank.BiliCoinRatio] = 0.2
	weights[rank.BiliReplyRatio] = 0.2

	// item feature: bbq
	weights[rank.BBQPlayNum] = 0.05
	weights[rank.BBQFavRatio] = 0.0
	weights[rank.BBQLikeRatio] = 0.2
	weights[rank.BBQShareRatio] = 0.2
	weights[rank.BBQCoinRatio] = 0.0
	weights[rank.BBQReplyRatio] = 0.2

	return
}

func (m *RankManager) rank(c context.Context, req *recsys.RecsysRequest, response *recsys.RecsysResponse, profile *model.UserProfile, dao *dao.Dao) {
	ranker := m.Rankers[0]
	m.featureLogs = ranker.rank(c, req, response, profile, dao)
}

//BaseRanker ...
type BaseRanker struct {
	Ranker
	weights map[string]float64
}

func (r *BaseRanker) name() (name string) {
	name = "base"
	return
}

func (r *BaseRanker) rank(c context.Context, req *recsys.RecsysRequest, response *recsys.RecsysResponse, userProfile *model.UserProfile, dao *dao.Dao) (featureLogs []*rank.FeatureLog) {

	response.Message[model.RankModelName] = "rule001"

	//dao.InitModel(c, r.weights)

	featureLogs = make([]*rank.FeatureLog, 0)

	for _, record := range response.List {

		featureLog, featureValueMap := rank.BuildFeature(record, userProfile)

		score := 0.0
		scoreMap := make(map[string]string)
		for feature, weight := range r.weights {
			featureValue := featureValueMap[feature]
			score = score + featureValue*weight

			if req.DebugFlag {
				if math.Abs(featureValue*weight) > 0.0001 {
					scoreMap[feature] = fmt.Sprintf("%.6f=%.6f*%.6f", featureValue*weight, featureValue, weight)
				}
			}
		}
		record.Score = score
		featureLog.Score = score

		//debug log
		if req.DebugFlag {
			var buffer bytes.Buffer
			buffer.WriteString(model.ScoreTotalScore)
			buffer.WriteString(":")
			buffer.WriteString(strconv.FormatFloat(score, 'f', 6, 64))
			buffer.WriteString("=")

			scoreDetailList := util.SortStrMapByValue(scoreMap)
			for _, pair := range scoreDetailList {
				buffer.WriteString(pair.Key)
				buffer.WriteString(":")
				buffer.WriteString(pair.Value)
				buffer.WriteString(",")
			}
			record.Map[model.ScoreMessage] = buffer.String()
			featureLogStr, _ := jsoniter.MarshalToString(featureLog)
			record.Map[model.FeatureString] = featureLogStr
		}

	}

	sort.Sort(sort.Reverse(util.Records(response.List)))
	for index, record := range response.List {
		record.Map[model.OrderRanker] = strconv.Itoa(index)
		record.Map[model.RankModelScore] = strconv.FormatFloat(record.Score, 'f', -1, 64)
	}

	return
}
