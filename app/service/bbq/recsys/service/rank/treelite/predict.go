package treelite

import (
	"errors"
	"math"
	"strconv"
)

//Predict Predict single tree
func (node *Node) Predict(fvals []float64) float64 {
	if node.LeftChild == nil && node.RightChild == nil {
		return node.GetLeafValue()
	}
	for {
		val := fvals[node.GetSplitIndex()]
		if val < node.GetThreshold() {
			return node.GetLeftChild().Predict(fvals)
		}
		return node.GetRightChild().Predict(fvals)
	}
}

//ValidateModel ...
func (model *Model) ValidateModel() (valid bool, err error) {
	if model.GetNumFeature() < 1 {
		err = errors.New("number of Feature < 1")
		return false, err
	}

	if model.GetNumOutputGroup() != 1 {
		err = errors.New("number of output group != 1")
		return false, err
	}
	if model.GetRandomForestFlag() {
		err = errors.New("do not support random forest model now")
		return false, err
	}
	for _, tree := range model.Trees {
		if tree.GetHead().GetSplitType() != Node_NUMERICAL {
			err = errors.New("tree.GetHead().GetSplitType() != Node_NUMERICAL")
			return false, err
		}
	}
	return true, nil
}

//PredictSingle ...
func (model *Model) PredictSingle(fvals []float64) (predictVal float64) {

	//if predTransform, ok := model.ExtraParams["pred_transform"]; ok {
	//	switch predTransform {
	//	case "sigmoid":
	//
	//	}
	//}

	predictVal, _ = strconv.ParseFloat(model.ExtraParams["global_bias"], 64)
	for _, tree := range model.Trees {
		predictVal += tree.GetHead().Predict(fvals)
	}
	predictVal = sigmoid(predictVal)
	return predictVal
}

func sigmoid(x float64) (y float64) {
	return 1.0 / (1.0 + math.Exp(-x))
}
