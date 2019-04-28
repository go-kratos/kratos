package service

import (
	"fmt"
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/library/log"
	"time"

	"github.com/json-iterator/go"
)

//StoreLog ...stores request and response log
func (s *Service) StoreLog(request *recsys.RecsysRequest, response *recsys.RecsysResponse, user *model.UserProfile, business string) {

	reqString, _ := jsoniter.MarshalToString(request)
	responseBytes, _ := jsoniter.Marshal(response)

	ext := make(map[string]string)
	ext["queryid"] = request.QueryID
	ext["traceid"] = request.TraceID
	ext["rankmodel"] = response.Message[model.RankModelName]
	ext[model.ResponseDownGrade] = response.Message[model.ResponseDownGrade]
	extString, _ := jsoniter.MarshalToString(ext)

	s.infoc.Info(request.MID, request.BUVID, time.Now().Unix(), reqString, string(responseBytes), request.Abtest, business, extString)

	// 5.4 reduce record keys
	records := make([]*recsys.RecsysRecord, 0)
	for _, record := range response.List {
		newRecord := &recsys.RecsysRecord{
			Svid:  record.Svid,
			Score: record.Score,
			Map:   make(map[string]string),
		}
		newRecord.Map[model.RecallClasses] = record.Map[model.RecallClasses]
		newRecord.Map[model.RecallTags] = record.Map[model.RecallTags]
		newRecord.Map[model.AVID] = record.Map[model.AVID]
		newRecord.Map[model.CID] = record.Map[model.CID]
		newRecord.Map[model.State] = record.Map[model.State]
		newRecord.Map[model.ScatterTag] = record.Map[model.ScatterTag]
		newRecord.Map[model.UperMid] = record.Map[model.UperMid]
		newRecord.Map[model.Title] = record.Map[model.Title]
		records = append(records, newRecord)
	}
	response.List = records
	responseStr, _ := jsoniter.MarshalToString(response)
	log.Info(fmt.Sprintf("response log: mid:%v, buvid:%v, %v, %v, %v, %v, %v, %v", request.MID, request.BUVID, time.Now().Unix(), reqString, responseStr, request.Abtest, business, extString))
}
