package postprocess

import (
	"context"
	"go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/app/service/bbq/recsys/service/retrieve"
	"strconv"
	"strings"
)

//RelevantInsertProcessor ..
type RelevantInsertProcessor struct {
	Processor
}

func (p *RelevantInsertProcessor) name() (name string) {
	name = "RelevantInsert"
	return
}

func (p *RelevantInsertProcessor) process(ctx context.Context, request *v1.RecsysRequest, response *v1.RecsysResponse, u *model.UserProfile) (err error) {

	insertPosition := 0
	if insertPosition > len(response.List) {
		return
	}

	targetIndex := -1
	for index, record := range response.List {
		if targetIndex != -1 {
			break
		}
		for _, recallClass := range strings.Split(record.Map[model.RecallClasses], "|") {
			if recallClass == retrieve.LikeI2IRecall || recallClass == retrieve.LikeTagRecall || recallClass == retrieve.LikeUPRecall || recallClass == retrieve.FollowRecall {
				if actionTimeToNow, ok := record.Map[model.SourceTimeToNow]; ok {
					timeInSceonds, _ := strconv.ParseInt(actionTimeToNow, 10, 64)
					if timeInSceonds <= 2*3600 {
						if index <= insertPosition {
							break
						} else {
							targetIndex = index
							break
						}
					}
				}
			}
		}
	}

	if targetIndex != -1 {
		record := response.List[targetIndex]
		record.Map[model.OrderPostProcess] = p.name()
		response.List = append(response.List[:targetIndex], response.List[targetIndex+1:]...)
		tmpList := append(response.List[:insertPosition], record)
		response.List = append(tmpList, response.List[insertPosition+1:]...)
	}
	return
}
