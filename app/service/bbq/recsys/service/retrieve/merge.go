package retrieve

import (
	"fmt"
	recallv1 "go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/library/log"
	"strconv"
	"strings"
)

func deleteBlack(response *recsys.RecsysResponse, userProfile *model.UserProfile) (err error) {
	records := make([]*recsys.RecsysRecord, 0)
	for _, record := range response.List {
		upMID, _ := strconv.ParseInt(record.Map[model.UperMid], 10, 64)
		if _, ok := userProfile.BBQBlack[upMID]; ok {
			continue
		}
		records = append(records, record)
	}
	response.List = records
	return
}

func transform(recallResponse *recallv1.RecallResponse, response *recsys.RecsysResponse) (err error) {

	if recallResponse == nil {
		return
	}
	response.Message[model.ResponseRecallStat] = fmt.Sprintf("%v", recallResponse.SrcInfo)

	for index, video := range recallResponse.List {
		if video.ForwardIndex == nil || video.ForwardIndex.BasicInfo == nil {
			log.Error("recall forward index null, svid: %v", video.SVID)
			continue
		}
		record := &recsys.RecsysRecord{
			Svid:  video.SVID,
			Score: 0,
			Map:   make(map[string]string),
		}

		// 视频基本信息
		record.Map[model.Title] = video.ForwardIndex.BasicInfo.Title
		record.Map[model.Content] = video.ForwardIndex.BasicInfo.Content
		record.Map[model.AVID] = strconv.Itoa(int(video.ForwardIndex.BasicInfo.AVID))
		record.Map[model.CID] = strconv.Itoa(int(video.ForwardIndex.BasicInfo.CID))
		record.Map[model.State] = strconv.Itoa(int(video.ForwardIndex.BasicInfo.State))
		record.Map[model.UperMid] = strconv.Itoa(int(video.ForwardIndex.BasicInfo.MID))
		record.Map[model.PubTime] = strconv.Itoa(int(video.ForwardIndex.BasicInfo.PubTime))
		record.Map[model.Duration] = strconv.Itoa(int(video.ForwardIndex.BasicInfo.Duration))

		tagNames := make([]string, 0)
		tagTypes := make([]string, 0)
		tagIDs := make([]string, 0)
		for _, tag := range video.ForwardIndex.BasicInfo.Tags {
			tagNames = append(tagNames, tag.TagName)
			tagTypes = append(tagTypes, strconv.Itoa(int(tag.TagType)))
			tagIDs = append(tagIDs, strconv.Itoa(int(tag.TagID)))

			if tag.TagType == 2 {
				record.Map[model.ZoneID] = strconv.Itoa(int(tag.TagID))
				record.Map[model.ZoneName] = tag.TagName
			}
		}
		record.Map[model.TagsName] = strings.Join(tagNames, "|")
		record.Map[model.TagsType] = strings.Join(tagTypes, "|")
		record.Map[model.TagsID] = strings.Join(tagIDs, "|")

		// 视频质量信息
		if video.ForwardIndex.VideoQuality != nil {
			//bili
			record.Map[model.PlayHive] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsM.PlayCnt))
			record.Map[model.FavHive] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsM.FavCnt))
			record.Map[model.LikesHive] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsM.LikeCnt))
			record.Map[model.CoinHive] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsM.CoinCnt))
			record.Map[model.ReplyHive] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsM.CommentAddCnt))
			record.Map[model.DanmuHive] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsM.DanmuCnt))
			record.Map[model.ShareHive] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsM.ShareCnt))

			record.Map[model.PlayWeekBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsW.PlayCnt))
			record.Map[model.LikesWeekBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsW.LikeCnt))
			record.Map[model.ReplyWeekBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsW.CommentAddCnt))
			record.Map[model.DanmuWeekBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsW.DanmuCnt))
			record.Map[model.ShareWeekBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsW.ShareCnt))
			record.Map[model.FavWeekBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsW.FavCnt))
			record.Map[model.CoinWeekBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsW.CoinCnt))

			record.Map[model.PlayDayBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsH.PlayCnt))
			record.Map[model.LikesDayBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsH.LikeCnt))
			record.Map[model.ReplyDayBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsH.CommentAddCnt))
			record.Map[model.DanmuDayBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsH.DanmuCnt))
			record.Map[model.ShareDayBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsH.ShareCnt))
			record.Map[model.FavDayBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsH.FavCnt))
			record.Map[model.CoinDayBili] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoMsH.CoinCnt))

			// bbq
			record.Map[model.PlayMonthTotal] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.ImpCnt))
			record.Map[model.PlayMonthFinish] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.AbsolutePlayCnt))
			record.Map[model.PlayMonth] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.PlayCnt))
			record.Map[model.FavMonth] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.FavCnt))
			record.Map[model.LikesMonth] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.LikeCnt))
			record.Map[model.ReplyMonth] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.CommentAddCnt))
			record.Map[model.DanmuMonth] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.DanmuCnt))
			record.Map[model.ShareMonth] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.ShareCnt))
			//record.Map[model.CommentLikeMonth] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.CommentLikeCnt))
			//record.Map[model.CommentReportMonth] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoM.CommentReportCnt))

			record.Map[model.PlayWeekFinish] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoW.AbsolutePlayCnt))
			record.Map[model.PlayWeek] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoW.PlayCnt))
			record.Map[model.LikesWeek] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoW.LikeCnt))
			record.Map[model.ReplyWeek] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoW.CommentAddCnt))
			record.Map[model.DanmuWeek] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoW.DanmuCnt))
			record.Map[model.ShareWeek] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoW.ShareCnt))

			record.Map[model.PlayDayFinish] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoH.AbsolutePlayCnt))
			record.Map[model.PlayDay] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoH.PlayCnt))
			record.Map[model.LikesDay] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoH.LikeCnt))
			record.Map[model.ReplyDay] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoH.CommentAddCnt))
			record.Map[model.DanmuDay] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoH.DanmuCnt))
			record.Map[model.ShareDay] = strconv.Itoa(int(video.ForwardIndex.VideoQuality.QualityInfoH.ShareCnt))
		}

		// 召回信息
		record.Map[model.RecallScore] = strconv.FormatFloat(float64(video.Score), 'f', -1, 32)
		record.Map[model.RecallOrder] = strconv.Itoa(index)

		recallTags := make([]string, 0)
		recallClasses := make([]string, 0)
		for _, invertIndex := range video.InvertedIndexes {
			recallTags = append(recallTags, invertIndex.Index)
			recallClasses = append(recallClasses, invertIndex.Name)
		}
		record.Map[model.RecallTags] = strings.Join(recallTags, "|")
		record.Map[model.RecallClasses] = strings.Join(recallClasses, "|")

		response.List = append(response.List, record)
	}
	return
}

func mergeRecallKey(recallInfos []*recallv1.RecallInfo) (newRecallInfos []*recallv1.RecallInfo) {
	recallTagNameMap := make(map[string][]string)
	recallTagInfoMap := make(map[string]*recallv1.RecallInfo)
	recallTagPriorityMap := make(map[string]int32)

	for _, recallInfo := range recallInfos {
		names := recallTagNameMap[recallInfo.Tag]
		names = append(names, recallInfo.Name)
		recallTagNameMap[recallInfo.Tag] = names
		recallTagInfoMap[recallInfo.Tag] = recallInfo

		if priority, ok := recallTagPriorityMap[recallInfo.Tag]; ok {
			if recallInfo.Priority > priority {
				recallTagPriorityMap[recallInfo.Tag] = priority
			}
		} else {
			recallTagPriorityMap[recallInfo.Tag] = priority
		}
	}

	newRecallInfos = make([]*recallv1.RecallInfo, 0)
	for tag, names := range recallTagNameMap {
		recallInfo := recallTagInfoMap[tag]
		recallInfo.Name = strings.Join(names, "|")
		newRecallInfos = append(newRecallInfos, recallInfo)
	}
	return
}
