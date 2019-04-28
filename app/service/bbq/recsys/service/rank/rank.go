package rank

import (
	"errors"
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	xgb "go-common/app/service/bbq/recsys/service/rank/treelite"
	"go-common/app/service/bbq/recsys/service/util"
	"go-common/library/log"
	"sort"
	"strconv"
)

//Rank Model Const
const (
	DefaultModelKey   = "xgb_model_v0.0.13"
	ModelNameTemplate = "xgb_model_v%s"

	EnvPredictModelDirs = "ENV_PREDICT_MODLE_DIR_LIST"
)

//RankModelManager ...
type RankModelManager struct {
	RankModels map[string]*RankModel
}

//NewRankModelManager ...
func NewRankModelManager() (m *RankModelManager) {
	m = &RankModelManager{}

	models, err := m.loadModels()
	if len(models) == 0 && err != nil {
		log.Error("load model error %v", err)
		return
	}

	m.RankModels = make(map[string]*RankModel, len(models))
	for _, model := range models {
		m.RankModels[model.name] = model
	}

	return
}

//RankModel ...
type RankModel struct {
	model          *xgb.Model
	name           string
	score          func([]float64) float64
	buildInstances func(featureLogs []*FeatureLog) (instances []*Instance)
}

//FeatureConf ...
type FeatureConf struct {
}

//Instance ...
type Instance struct {
	record        *recsys.RecsysRecord
	featureLog    *FeatureLog
	featureValues *[]float64
}

//DoRank ...
func (m *RankModelManager) DoRank(request *recsys.RecsysRequest, response *recsys.RecsysResponse, userProfile *model.UserProfile) (err error) {
	// 1.0 choose model
	rankModel, ok := m.RankModels[DefaultModelKey]
	if !ok {
		return errors.New("rank model is missing")
	}

	response.Message[model.RankModelName] = rankModel.name
	return rankModel.rank(request, response, userProfile)
}

func (rankModel *RankModel) rank(request *recsys.RecsysRequest, response *recsys.RecsysResponse, userProfile *model.UserProfile) (err error) {
	// 2.0 init/load model & feature conf
	// 3.0 build feature
	// 3.0 build instances (feature + conf + operator -> instance)
	// 4.0 score each record
	// 5.0 rank

	//build features
	features := rankModel.buildFeatures(request, response, userProfile)

	instances := rankModel.buildInstances(features)

	for _, instance := range instances {
		//rankModel.score(instance)
		score := rankModel.model.PredictSingle(*instance.featureValues)
		instance.record.Score = score
		instance.featureLog.Score = score
	}

	sort.Sort(sort.Reverse(util.Records(response.List)))

	for index, record := range response.List {
		record.Map[model.OrderRanker] = strconv.Itoa(index)
		record.Map[model.RankModelScore] = strconv.FormatFloat(record.Score, 'f', -1, 64)
	}
	return
}
