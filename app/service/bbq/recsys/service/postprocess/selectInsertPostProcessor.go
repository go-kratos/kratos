package postprocess

import (
	"context"
	"go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/app/service/bbq/recsys/service/retrieve"
)

//SelectInsertProcessor ..
type SelectInsertProcessor struct {
	Processor
}

func (p *SelectInsertProcessor) name() (name string) {
	name = "SelectInsert"
	return
}

func (p *SelectInsertProcessor) process(ctx context.Context, request *v1.RecsysRequest, response *v1.RecsysResponse, u *model.UserProfile) (err error) {

	//response.List[100].Map[model.Retriever] = retrieve.Test

	insertPosition := 0
	if insertPosition > len(response.List) {
		return
	}

	targetIndex := -1
	for index, record := range response.List {
		if retriever, ok := record.Map[model.Retriever]; ok {
			if index > insertPosition && retriever == retrieve.SelectionRecall {
				targetIndex = index
				break
			}
		}
	}

	if targetIndex != -1 {
		record := response.List[targetIndex]
		//标记
		record.Map[model.OrderPostProcess] = p.name()
		response.List = append(response.List[:targetIndex], response.List[targetIndex+1:]...)
		tmpList := append(response.List[:insertPosition], record)
		response.List = append(tmpList, response.List[insertPosition+1:]...)
	}

	return
}
