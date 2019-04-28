package postprocess

import (
	"context"
	"go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/app/service/bbq/recsys/service/util"
	"sort"
	"strconv"
)

//WeakInterventionProcessor ...
type WeakInterventionProcessor struct {
	Processor
}

func (p *WeakInterventionProcessor) name() (name string) {
	name = "WeakIntervention"
	return
}

func (p *WeakInterventionProcessor) process(ctx context.Context, request *v1.RecsysRequest, response *v1.RecsysResponse, u *model.UserProfile) (err error) {
	if response == nil || len(response.List) == 0 {
		return
	}
	for _, record := range response.List {
		if state, ok := record.Map[model.State]; ok {
			stateID, _ := strconv.ParseInt(state, 10, 64)
			switch stateID {
			case model.State4:
				record.Score = 1.02 * record.Score
				continue
			case model.State5:
				record.Score = 1.05 * record.Score
				continue
			default:
				continue
			}
		}
	}
	//sort
	sort.Sort(sort.Reverse(util.Records(response.List)))
	for index, record := range response.List {
		record.Map[model.OrderWeakIntervention] = strconv.Itoa(index)
	}
	return
}
