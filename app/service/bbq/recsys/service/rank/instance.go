package rank

import (
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"

	"fmt"
	"go-common/app/service/bbq/recsys/model"
	"strings"
)

func (rankModel *RankModel) buildFeatures(request *recsys.RecsysRequest, response *recsys.RecsysResponse, userProfile *model.UserProfile) (featureLogs []*FeatureLog) {

	featureLogs = make([]*FeatureLog, len(response.List))

	for index, record := range response.List {

		featureLog, featureValueMap := BuildFeature(record, userProfile)
		//FIXME ...
		featureValueMap["test"] = 1
		featureLogs[index] = featureLog
	}

	return
}

func (rankModel *RankModel) buildInstancesV1(featureLogs []*FeatureLog) (instances []*Instance) {
	featureMap := map[string]int64{
		"HotRecall":        0,
		"LikeI2IRecall":    1,
		"LikeTagRecall":    2,
		"LikeUPRecall":     3,
		"NewPublishRecall": 4,
		"PosI2IRecall":     5,
		"PosTagRecall":     6,
		"SelectionRecall":  7,
		"UserProfileBili":  8,
	}
	featureSize := len(featureMap)

	instances = make([]*Instance, len(featureLogs))
	for i, instance := range instances {
		featureLog := featureLogs[i]

		featureValues := make([]float64, featureSize)

		//TODO recall feature: single recall -> multiple recall
		if index, ok := featureMap[featureLog.RetrieveName]; ok {
			featureValues[index] = 1
		}
		instance = &Instance{
			record:        featureLog.record,
			featureValues: &featureValues,
		}
		instances[i] = instance
	}

	return
}

func (rankModel *RankModel) buildInstancesV2(featureLogs []*FeatureLog) (instances []*Instance) {
	featureMap := map[string]int64{
		"HotRecall":          0,
		"LikeI2IRecall":      1,
		"LikeTagRecall":      2,
		"LikeUPRecall":       3,
		"NewPublishRecall":   4,
		"PosTagRecall":       5,
		"SelectionRecall":    6,
		"UserProfileBBQ":     7,
		"UserProfileBili":    8,
		"play_hive":          9,
		"likes_hive":         10,
		"fav_hive":           11,
		"reply_hive":         12,
		"share_hive":         13,
		"coin_hive":          14,
		"play_month_finish":  15,
		"play_month":         16,
		"likes_month":        17,
		"reply_month":        18,
		"share_month":        19,
		"has_tag_count":      20,
		"contains_tag_count": 21,
	}
	featureSize := len(featureMap)

	instances = make([]*Instance, len(featureLogs))
	for i := range instances {
		featureLog := featureLogs[i]

		//TODO
		featureValues := make([]float64, featureSize)
		if index, ok := featureMap[featureLog.RetrieveName]; ok {
			featureValues[index] = 1
		}
		featureValues[featureMap["play_hive"]] = float64(featureLog.PlayB)
		featureValues[featureMap["likes_hive"]] = float64(featureLog.LikeB)
		featureValues[featureMap["fav_hive"]] = float64(featureLog.FavB)
		featureValues[featureMap["reply_hive"]] = float64(featureLog.ReplyB)
		featureValues[featureMap["share_hive"]] = float64(featureLog.ShareB)
		featureValues[featureMap["coin_hive"]] = float64(featureLog.CoinB)
		featureValues[featureMap["play_month_finish"]] = float64(featureLog.PlayBBQFinish)
		featureValues[featureMap["play_month"]] = float64(featureLog.PlayBBQ)
		featureValues[featureMap["likes_month"]] = float64(featureLog.LikeBBQ)
		featureValues[featureMap["reply_month"]] = float64(featureLog.ReplyBBQ)
		featureValues[featureMap["share_month"]] = float64(featureLog.ShareBBQ)

		if featureLog.MatchBBQTagCount > 0 {
			featureValues[featureMap["has_tag_count"]] = 1.0
		} else {
			featureValues[featureMap["has_tag_count"]] = 0
		}
		featureValues[featureMap["contains_tag_count"]] = float64(featureLog.MatchBBQTagCount)

		instance := &Instance{
			record:        featureLog.record,
			featureValues: &featureValues,
			featureLog:    featureLog,
		}
		instances[i] = instance
	}

	return
}

func (rankModel *RankModel) buildInstancesV3(featureLogs []*FeatureLog) (instances []*Instance) {
	featureMap := map[string]int64{
		"HotRecall":          0,
		"RandomRecall":       1,
		"SelectionRecall":    2,
		"UserProfileBili":    3,
		"UserProfileBBQ":     4,
		"LikeI2IRecall":      5,
		"LikeTagRecall":      6,
		"LikeUPRecall":       7,
		"PosI2IRecall":       8,
		"PosTagRecall":       9,
		"FollowRecall":       10,
		"play_hive":          11,
		"fav_hive":           12,
		"reply_hive":         13,
		"share_hive":         14,
		"coin_hive":          15,
		"play_month_finish":  16,
		"play_month":         17,
		"likes_month":        18,
		"reply_month":        19,
		"share_month":        20,
		"has_tag_count":      21,
		"contains_tag_count": 22,
	}
	featureSize := len(featureMap)

	instances = make([]*Instance, len(featureLogs))
	for i := range instances {
		featureLog := featureLogs[i]

		//TODO
		featureValues := make([]float64, featureSize)

		recallClasses := strings.Split(featureLog.RecallClasses, "|")
		for _, recallClass := range recallClasses {
			if index, ok := featureMap[recallClass]; ok {
				featureValues[index] = 1
			}
		}

		featureValues[featureMap["play_hive"]] = float64(featureLog.PlayB)
		//featureValues[featureMap["likes_hive"]] = float64(featureLog.LikeB)
		featureValues[featureMap["fav_hive"]] = float64(featureLog.FavB)
		featureValues[featureMap["reply_hive"]] = float64(featureLog.ReplyB)
		featureValues[featureMap["share_hive"]] = float64(featureLog.ShareB)
		featureValues[featureMap["coin_hive"]] = float64(featureLog.CoinB)
		featureValues[featureMap["play_month_finish"]] = float64(featureLog.PlayBBQFinish)
		featureValues[featureMap["play_month"]] = float64(featureLog.PlayBBQ)
		featureValues[featureMap["likes_month"]] = float64(featureLog.LikeBBQ)
		featureValues[featureMap["reply_month"]] = float64(featureLog.ReplyBBQ)
		featureValues[featureMap["share_month"]] = float64(featureLog.ShareBBQ)

		if featureLog.MatchBBQTagCount > 0 {
			featureValues[featureMap["has_tag_count"]] = 1.0
		} else {
			featureValues[featureMap["has_tag_count"]] = 0
		}
		featureValues[featureMap["contains_tag_count"]] = float64(featureLog.MatchBBQTagCount)

		instance := &Instance{
			record:        featureLog.record,
			featureValues: &featureValues,
			featureLog:    featureLog,
		}
		instances[i] = instance
	}

	return
}

func (rankModel *RankModel) buildInstancesV12(featureLogs []*FeatureLog) (instances []*Instance) {
	featureMap := map[string]int64{
		"HotRecall":          0,
		"RandomRecall":       1,
		"SelectionRecall":    2,
		"UserProfileBili":    3,
		"UserProfileBBQ":     4,
		"LikeI2IRecall":      5,
		"LikeTagRecall":      6,
		"LikeUPRecall":       7,
		"PosI2IRecall":       8,
		"PosTagRecall":       9,
		"FollowRecall":       10,
		"has_tag_count":      11,
		"contains_tag_count": 12,
		"has_zone_count":     13,
		"play_hive":          14,
		"fav_hive":           15,
		"reply_hive":         16,
		"share_hive":         17,
		"coin_hive":          18,
		"play_month_finish":  19,
		"play_month":         20,
		"likes_month":        21,
		"reply_month":        22,
		"share_month":        23,
	}
	featureSize := len(featureMap)

	instances = make([]*Instance, len(featureLogs))
	for i := range instances {
		featureLog := featureLogs[i]

		featureValues := make([]float64, featureSize)
		recallClasses := strings.Split(featureLog.RecallClasses, "|")
		for _, recallClass := range recallClasses {
			if index, ok := featureMap[recallClass]; ok {
				featureValues[index] = 1
			}
		}

		featureValues[featureMap["play_hive"]] = float64(featureLog.PlayB)
		featureValues[featureMap["fav_hive"]] = float64(featureLog.FavB)
		featureValues[featureMap["reply_hive"]] = float64(featureLog.ReplyB)
		featureValues[featureMap["share_hive"]] = float64(featureLog.ShareB)
		featureValues[featureMap["coin_hive"]] = float64(featureLog.CoinB)
		featureValues[featureMap["play_month_finish"]] = float64(featureLog.PlayBBQFinish)
		featureValues[featureMap["play_month"]] = float64(featureLog.PlayBBQ)
		featureValues[featureMap["likes_month"]] = float64(featureLog.LikeBBQ)
		featureValues[featureMap["reply_month"]] = float64(featureLog.ReplyBBQ)
		featureValues[featureMap["share_month"]] = float64(featureLog.ShareBBQ)

		if featureLog.MatchBBQTagCount > 0 {
			featureValues[featureMap["has_tag_count"]] = 1.0
		} else {
			featureValues[featureMap["has_tag_count"]] = 0
		}
		if featureLog.MatchBBQTagLevel2 > 0 {
			featureValues[featureMap["has_zone_count"]] = 1.0
		} else {
			featureValues[featureMap["has_zone_count"]] = 0
		}
		featureValues[featureMap["contains_tag_count"]] = float64(featureLog.MatchBBQTagCount)

		instance := &Instance{
			record:        featureLog.record,
			featureValues: &featureValues,
			featureLog:    featureLog,
		}
		instances[i] = instance
	}

	return
}

func (rankModel *RankModel) buildInstancesV13(featureLogs []*FeatureLog) (instances []*Instance) {
	featureMap := map[string]int64{
		"zone-bucket-168":        0,
		"zone-bucket-75":         1,
		"play_hive":              2,
		"zone-bucket-95":         3,
		"likes_month":            4,
		"state-bucket-3":         5,
		"share_month":            6,
		"recall-PosTagRecall":    7,
		"zone-bucket-124":        8,
		"recall-PosI2IRecall":    9,
		"state-bucket-4":         10,
		"recall-SelectionRecall": 11,
		"zone-bucket-156":        12,
		"contains_tag_count":     13,
		"zone-bucket-158":        14,
		"zone-bucket-183":        15,
		"zone-bucket-184":        16,
		"zone-bucket-21":         17,
		"zone-bucket-154":        18,
		"zone-bucket-159":        19,
		"zone-bucket-85":         20,
		"recall-LikeUPRecall":    21,
		"reply_month":            22,
		"state-bucket-1":         23,
		"zone-bucket-96":         24,
		"has_tag_count":          25,
		"zone-bucket-86":         26,
		"zone-bucket-138":        27,
		"zone-bucket-182":        28,
		"play_month_finish":      29,
		"recall-HotRecall":       30,
		"zone-bucket-157":        31,
		"zone-bucket-20":         32,
		"zone-bucket-39":         33,
		"zone-bucket-161":        34,
		"reply_hive":             35,
		"recall-LikeTagRecall":   36,
		"zone-bucket-76":         37,
		"zone-bucket-98":         38,
		"state-bucket-5":         39,
		"zone-bucket-22":         40,
		"zone-bucket-27":         41,
		"zone-bucket-122":        42,
		"zone-bucket-176":        43,
		"recall-UserProfileBBQ":  44,
		"recall-UserProfileBili": 45,
		"zone-bucket-163":        46,
		"zone-bucket-30":         47,
		"zone-bucket-31":         48,
		"zone-bucket-59":         49,
		"recall-LikeI2IRecall":   50,
		"zone-bucket-25":         51,
		"zone-bucket-28":         52,
		"zone-bucket-24":         53,
		"zone-bucket-29":         54,
		"zone-bucket-164":        55,
		"coin_hive":              56,
		"play_month":             57,
		"share_hive":             58,
		"recall-RandomRecall":    59,
		"fav_hive":               60,
		"zone-bucket-162":        61,
		"likes_hive":             62,
		"recall-FollowRecall":    63,
		"zone-bucket-47":         64,
	}
	featureSize := len(featureMap)

	instances = make([]*Instance, len(featureLogs))
	for i := range instances {
		featureLog := featureLogs[i]

		featureValues := make([]float64, featureSize)
		recallClasses := strings.Split(featureLog.RecallClasses, "|")
		for _, recallClass := range recallClasses {
			recallClass = fmt.Sprintf("recall-%s", recallClass)
			if index, ok := featureMap[recallClass]; ok {
				featureValues[index] = 1
			}
		}

		featureValues[featureMap["play_hive"]] = float64(featureLog.PlayB)
		featureValues[featureMap["likes_hive"]] = float64(featureLog.LikeB)
		featureValues[featureMap["fav_hive"]] = float64(featureLog.FavB)
		featureValues[featureMap["reply_hive"]] = float64(featureLog.ReplyB)
		featureValues[featureMap["share_hive"]] = float64(featureLog.ShareB)
		featureValues[featureMap["coin_hive"]] = float64(featureLog.CoinB)
		featureValues[featureMap["play_month_finish"]] = float64(featureLog.PlayBBQFinish)
		featureValues[featureMap["play_month"]] = float64(featureLog.PlayBBQ)
		featureValues[featureMap["likes_month"]] = float64(featureLog.LikeBBQ)
		featureValues[featureMap["reply_month"]] = float64(featureLog.ReplyBBQ)
		featureValues[featureMap["share_month"]] = float64(featureLog.ShareBBQ)

		if featureLog.MatchBBQTagCount > 0 {
			featureValues[featureMap["has_tag_count"]] = 1.0
		} else {
			featureValues[featureMap["has_tag_count"]] = 0
		}
		featureValues[featureMap["contains_tag_count"]] = float64(featureLog.MatchBBQTagCount)

		//bucket features
		zoneKey := fmt.Sprintf("zone-bucket-%d", featureLog.ZoneID)
		if index, ok := featureMap[zoneKey]; ok {
			featureValues[index] = 1
		}

		stateKey := fmt.Sprintf("state-bucket-%d", featureLog.State)
		if index, ok := featureMap[stateKey]; ok {
			featureValues[index] = 1
		}

		instance := &Instance{
			record:        featureLog.record,
			featureValues: &featureValues,
			featureLog:    featureLog,
		}
		instances[i] = instance
	}

	return
}
