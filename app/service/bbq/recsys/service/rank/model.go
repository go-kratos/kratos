package rank

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	xgb "go-common/app/service/bbq/recsys/service/rank/treelite"
	"go-common/library/log"

	"github.com/gogo/protobuf/proto"
)

func (m *RankModelManager) loadModels() (models []*RankModel, err error) {

	modelFileNamesStr := os.Getenv(EnvPredictModelDirs)
	log.Info("modelFileNamesStr: %s", modelFileNamesStr)
	if modelFileNamesStr == "" {
		return nil, fmt.Errorf("env variable %s is empty", EnvPredictModelDirs)
	}

	models = make([]*RankModel, 0)
	modelFileNames := strings.Split(modelFileNamesStr, ",")
	for _, modelFileDir := range modelFileNames {
		modelFileName := fmt.Sprintf("%s/model.proto", modelFileDir)
		if modelFileName == "" {
			err = fmt.Errorf("model name is empty: %s", modelFileNamesStr)
			continue
		}
		var model *xgb.Model
		model, err = m.readModel(modelFileName)
		if model == nil || err != nil {
			log.Error("read model error: (%v)", err)
			continue
		}
		log.Info("xgb NumOutputGroup:%d, NFeatures:%d, NEstimators:%d\n", model.GetNumOutputGroup(), model.GetNumFeature(), len(model.Trees))

		// model & feature conf
		if strings.Contains(modelFileName, "0.0.13") {
			rankModel := &RankModel{
				name:  fmt.Sprintf(ModelNameTemplate, "0.0.13"),
				model: model,
				score: model.PredictSingle,
			}
			rankModel.buildInstances = rankModel.buildInstancesV13
			models = append(models, rankModel)
		} else if strings.Contains(modelFileName, "0.0.12") {
			rankModel := &RankModel{
				name:  fmt.Sprintf(ModelNameTemplate, "0.0.12"),
				model: model,
				score: model.PredictSingle,
			}
			rankModel.buildInstances = rankModel.buildInstancesV12
			models = append(models, rankModel)
		} else if strings.Contains(modelFileName, "0.0.11") {
			rankModel := &RankModel{
				name:  fmt.Sprintf(ModelNameTemplate, "0.0.11"),
				model: model,
				score: model.PredictSingle,
			}
			rankModel.buildInstances = rankModel.buildInstancesV3
			models = append(models, rankModel)
		} else if strings.Contains(modelFileName, "0.0.5") {
			rankModel := &RankModel{
				name:  fmt.Sprintf(ModelNameTemplate, "0.0.5"),
				model: model,
				score: model.PredictSingle,
			}
			rankModel.buildInstances = rankModel.buildInstancesV2
			models = append(models, rankModel)
		} else if strings.Contains(modelFileName, "0.0.4") {
			rankModel := &RankModel{
				name:  fmt.Sprintf(ModelNameTemplate, "0.0.4"),
				model: model,
				score: model.PredictSingle,
			}
			rankModel.buildInstances = rankModel.buildInstancesV1
			models = append(models, rankModel)
		}
	}
	return models, err
}

func (m *RankModelManager) readModel(modelFileName string) (model *xgb.Model, err error) {
	model = &xgb.Model{}
	data, err := ioutil.ReadFile(modelFileName)
	if err != nil {
		return nil, err
	}
	log.Info("read model success: ", modelFileName)
	err = proto.Unmarshal(data, model)
	if err != nil {
		return nil, err
	}

	valid, err := model.ValidateModel()
	if !valid || err != nil {
		log.Error("model is not valid %v", err)
		return nil, err
	}

	// test
	vals := []float64{0, 0, 1, 0, 0, 0}
	fvals := make([]float64, model.GetNumFeature())
	copy(vals, fvals)

	p := model.PredictSingle(fvals)
	log.Info("xgb Test Prediction for %v: %f\n", fvals, p)
	return model, err
}
