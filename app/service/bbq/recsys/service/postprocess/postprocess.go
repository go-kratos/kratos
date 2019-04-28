package postprocess

import (
	"context"
	"go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
)

//TagTypeZone2 二级分区
const (
	TagTypeZone2 = "2"

	// 推荐页打散一刷逻辑+相邻队列打散逻辑,
	// 因为相邻队列打散逻辑中只考虑连续情况,所以只用考虑之前队列的最后一个元素与当前队列中最前一个元素的关系
	// 所以可以不用做单独的相邻队列打散逻辑
	_RecTagTotalLimit     = 2
	_RecTagAdjacencyLimit = 2
	_RecUpTotalLimit      = 1
	_RecUpAdjacencyLimit  = 1
	_RecLastScreenCnt     = 5

	// 关注推荐打散一刷逻辑+相邻队列打散逻辑,
	// 因为相邻队列打散逻辑中只考虑连续情况,所以只用考虑之前队列的最后一个元素与当前队列中最前一个元素的关系
	// 所以可以不用做单独的相邻队列打散逻辑
	_UpsRecTagTotalLimit     = 5
	_UpsRecTagAdjacencyLimit = 2
	_UpsRecUpTotalLimit      = 1
	_UpsRecUpAdjacencyLimit  = 1
	_UpsRecLastScreenCnt     = 10

	////相邻队列打散,目前推荐页与关注推荐的逻辑是一致的
	//_AdjacentQueueTagTotalLimit     = 0
	//_AdjacentQueueTagAdjacencyLimit = 2
	//_AdjacentQueueUpTotalLimit      = 0
	//_AdjacentQueueUpAdjacencyLimit  = 1
)

//PostProcessor ...
type PostProcessor struct {
	name          string
	processors    []Processor
	ProcessRec    process
	ProcessUpsRec process
}

type process func(ctx context.Context, request *v1.RecsysRequest, response *v1.RecsysResponse, u *model.UserProfile) (err error)

//Processor ...
type Processor interface {
	name() (name string)

	process(ctx context.Context, request *v1.RecsysRequest, response *v1.RecsysResponse, u *model.UserProfile) (err error)
}

//NewPostProcessor ...
func NewPostProcessor() (p *PostProcessor) {

	processRec := p.buildProcessRec()
	processUpsRec := p.buildProcessUpsRec()

	p = &PostProcessor{
		name:          "post",
		processors:    make([]Processor, 0),
		ProcessRec:    processRec,
		ProcessUpsRec: processUpsRec,
	}

	return
}

func (p *PostProcessor) buildProcessRec() process {
	processors := make([]Processor, 0)

	weakInterventionProcessor := &WeakInterventionProcessor{}
	processors = append(processors, weakInterventionProcessor)

	downGradeProcessor := &DownGradeProcessor{}
	processors = append(processors, downGradeProcessor)

	relevantInsertProcessor := &RelevantInsertProcessor{}
	processors = append(processors, relevantInsertProcessor)

	scatterTagUpProcessor := &ScatterTagUpProcessor{lastScreenCount: _RecLastScreenCnt,
		tagTotalLimit:     _RecTagTotalLimit,
		tagAdjacencyLimit: _RecTagAdjacencyLimit,
		upTotalLimit:      _RecUpTotalLimit,
		upAdjacencyLimit:  _RecUpAdjacencyLimit,
		lastRecordsType:   "lastRecords"}
	processors = append(processors, scatterTagUpProcessor)
	process := func(ctx context.Context, request *v1.RecsysRequest, response *v1.RecsysResponse, u *model.UserProfile) (err error) {
		for _, processor := range processors {
			err = processor.process(ctx, request, response, u)
			if err != nil {
				break
			}
		}
		return
	}

	return process
}

func (p *PostProcessor) buildProcessUpsRec() process {
	processors := make([]Processor, 0)

	weakInterventionProcessor := &WeakInterventionProcessor{}
	processors = append(processors, weakInterventionProcessor)

	downGradeProcessor := &DownGradeProcessor{}
	processors = append(processors, downGradeProcessor)

	scatterTagUpProcessor := &ScatterTagUpProcessor{lastScreenCount: _UpsRecLastScreenCnt,
		tagTotalLimit:     _UpsRecTagTotalLimit,
		tagAdjacencyLimit: _UpsRecTagAdjacencyLimit,
		upTotalLimit:      _UpsRecUpTotalLimit,
		upAdjacencyLimit:  _UpsRecUpAdjacencyLimit,
		lastRecordsType:   "lastUpsRecords"}
	processors = append(processors, scatterTagUpProcessor)
	process := func(ctx context.Context, request *v1.RecsysRequest, response *v1.RecsysResponse, u *model.UserProfile) (err error) {
		for _, processor := range processors {
			err = processor.process(ctx, request, response, u)
			if err != nil {
				break
			}
		}
		return
	}

	return process
}
