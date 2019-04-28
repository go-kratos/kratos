package retrieve

import (
	"context"
	"fmt"
	"go-common/library/log"

	recallv1 "go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"strconv"
	"strings"
)

//召回策略
const (
	I2iRecall      = "I2iRecall"
	I2Zone2iRecall = "I2Zone2iRecall"
	I2tag2iRecall  = "I2tag2iRecall"
)

const (
	_recRecalli2i  = "RECALL:I2I:%d"
	_i2tag2iRecall = "RECALL:HOT_T:%d"
)

//Source Video Info
const (
	SourceTagIDs = "SourceTagIDs"
	SourceZoneID = "SourceZoneID"
	SourceUpMID  = "SourceUpMID"
)

//RelatedRec is retrieve function
func (m *RecallManager) RelatedRec(c context.Context, request *recsys.RecsysRequest, recallClient recallv1.RecsysRecallClient) (response *recsys.RecsysResponse, err error) {
	recallInfos := make([]*recallv1.RecallInfo, 0)
	SVID := request.SVID
	zoneID, tagIDs, upMID := tagIDget(c, recallClient, SVID)

	i2iRecallInfo := &recallv1.RecallInfo{
		Name:     I2iRecall,
		Tag:      fmt.Sprintf(_recRecalli2i, request.SVID),
		Limit:    50,
		Filter:   "",
		Priority: 4,
	}
	recallInfos = append(recallInfos, i2iRecallInfo)

	for _, tagID := range tagIDs {
		i2Tag2iRecallInfo := &recallv1.RecallInfo{
			Name:     I2tag2iRecall,
			Tag:      fmt.Sprintf(_i2tag2iRecall, tagID),
			Limit:    100,
			Filter:   "",
			Priority: 3,
		}
		recallInfos = append(recallInfos, i2Tag2iRecallInfo)
	}

	i2Zone2iRecallInfo := &recallv1.RecallInfo{
		Name:     I2Zone2iRecall,
		Tag:      fmt.Sprintf(_i2tag2iRecall, zoneID),
		Limit:    100,
		Filter:   "",
		Priority: 2,
	}
	recallInfos = append(recallInfos, i2Zone2iRecallInfo)

	// hotpoolRecall_Info 100
	hotRecallInfo := &recallv1.RecallInfo{
		Name:     HotRecall,
		Tag:      RecallHotDefault,
		Limit:    50,
		Filter:   "",
		Priority: 1,
	}
	recallInfos = append(recallInfos, hotRecallInfo)

	recallRequest := &recallv1.RecallRequest{
		MID:        request.MID,
		BUVID:      request.BUVID,
		TotalLimit: 100,
		Infos:      recallInfos,
	}
	log.Info("recall request: (%v)", recallRequest)

	response = new(recsys.RecsysResponse)
	response.Message = make(map[string]string)

	recallResponse, err := recallClient.Recall(c, recallRequest)

	if zoneID != 0 {
		response.Message[SourceZoneID] = strconv.Itoa(int(zoneID))
	}
	if len(tagIDs) > 0 {
		var params []string
		for _, tagID := range tagIDs {
			params = append(params, strconv.Itoa(int(tagID)))
		}
		response.Message[SourceTagIDs] = strings.Join(params, "|")
	}
	if upMID != 0 {
		response.Message[SourceUpMID] = strconv.Itoa(int(upMID))
	}

	if err != nil || recallResponse == nil {
		log.Error("recall service error (%v) or recall response is null", err)
		return
	}
	err = transform(recallResponse, response)
	return
}

func tagIDget(c context.Context, recallClient recallv1.RecsysRecallClient, SVID int64) (zoneID int64, tagIDs []int64, upMID int64) {

	SVIDs := make([]int64, 0)
	SVIDs = append(SVIDs, SVID)
	videoIndexRequest := &recallv1.VideoIndexRequest{
		SVIDs: SVIDs,
	}

	videoIndexResponse, err := recallClient.VideoIndex(c, videoIndexRequest)
	if err != nil || videoIndexResponse == nil {
		log.Error("recall service VideoIndex error (%v) or recall response is null", err)
		return
	}

	for _, forwardIndex := range videoIndexResponse.List {
		for _, tag := range forwardIndex.BasicInfo.Tags {
			if tag.TagType == 2 {
				zoneID = int64(tag.TagID)
			} else if tag.TagType == 3 {
				tagID := int64(tag.TagID)
				tagIDs = append(tagIDs, tagID)
			}
		}
		upMID = int64(forwardIndex.BasicInfo.MID)
	}

	return
}
