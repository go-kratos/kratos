package service

import (
	"context"

	"go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	"go-common/app/service/bbq/recsys-recall/model"
	"go-common/app/service/bbq/recsys-recall/service/index"
	"go-common/library/log"

	"github.com/Dai0522/workpool"
)

const (
	_recallLimit = 50
)

// Recall recsys recall video id list by tag
func (s *Service) Recall(ctx context.Context, request *v1.RecallRequest) (*v1.RecallResponse, error) {
	var response *v1.RecallResponse

	// recall from redis
	ftasks := s.parallelRecall(&ctx, request)
	recallResult := s.wait(ctx, ftasks)

	totalRecall := 0
	srcInfo := make([]*v1.RecallSrc, len(recallResult))
	for i, v := range recallResult {
		totalRecall += len(v.Result.Tuples)
		srcInfo[i] = &v1.RecallSrc{
			TotalHit: v.TotalHit,
			Filter:   v.FilterCount,
			Final:    v.FinalCount,
			Tag:      v.Tag,
			Name:     v.Name,
		}
	}
	offset := 0
	priorityTuples := make([]*model.PriorityTuple, totalRecall)
	for _, v := range recallResult {
		for _, u := range v.Result.Tuples {
			priorityTuples[offset] = &model.PriorityTuple{
				Tuple:    *u,
				Tag:      v.Tag,
				Name:     v.Name,
				Priority: v.Priority,
			}
			offset++
		}
	}

	// merge
	videos := s.merge(ctx, priorityTuples, request.TotalLimit)

	response = &v1.RecallResponse{
		Total:   int32(totalRecall),
		List:    videos,
		SrcInfo: srcInfo,
	}

	log.Infov(ctx, log.KV("total_recall", totalRecall), log.KV("result", len(videos)))

	return response, nil
}

func (s *Service) parallelRecall(ctx *context.Context, request *v1.RecallRequest) []*workpool.FutureTask {
	size := len(request.Infos)
	if size > _recallLimit {
		log.Errorv(*ctx, log.KV("RecallTag", size))
		size = _recallLimit
	}

	sc := NewScorerManager(s.dao)
	ranker := NewRankerManager(s.dao)
	filter := NewFilterManager(s.dao)
	filter.SetFilter("default", &DefaultFilter{})
	filter.SetFilter("bloomfilter", &BloomFilter{bf: s.loadBloomFilter(ctx, request.MID, request.BUVID)})

	var list []*v1.RecallInfo
	if len(request.Infos) > _recallLimit {
		// 优先级排序取前50
		list = recallSrcSortByPriority(request.Infos, size)
	} else {
		list = request.Infos
	}

	tasks := make([]*RecallTask, size)
	for i, v := range list {
		t := newRecallTask(ctx, s.dao, request.MID, request.BUVID, v)
		t.SetScorerManager(sc)
		t.SetRankerManager(ranker)
		t.SetFilterManager(filter)
		tasks[i] = t
	}
	return s.parallel(ctx, tasks)
}

// merge
func (s *Service) merge(c context.Context, tuples []*model.PriorityTuple, limit int32) []*v1.Video {
	list := sortByScore(tuples)
	list = sortByPriority(list)

	count := int32(0)
	videos := make(map[uint64]*v1.Video)
	for _, v := range list {
		if _, ok := videos[v.Svid]; !ok {
			if count > limit {
				continue
			}
			fi := index.Index.Get(v.Svid)
			if fi == nil {
				log.Errorv(c, log.KV("forward_index", nil), log.KV("svid", v.Svid))
			}
			videos[v.Svid] = &v1.Video{
				SVID:          int64(v.Svid),
				Score:         v.Score,
				Name:          v.Name,
				InvertedIndex: v.Tag,
				ForwardIndex:  fi,
				InvertedIndexes: []*v1.InvertedIndex{
					{
						Index: v.Tag,
						Name:  v.Name,
						Score: v.Score,
					},
				},
			}
			count++
		} else {
			videos[v.Svid].InvertedIndexes = append(videos[v.Svid].InvertedIndexes, &v1.InvertedIndex{
				Index: v.Tag,
				Name:  v.Name,
				Score: v.Score,
			})
		}
	}

	result := make([]*v1.Video, count)
	i := 0
	for _, v := range videos {
		result[i] = v
		i++
	}

	return result
}

func sortByScore(tuples []*model.PriorityTuple) []*model.PriorityTuple {
	for i := range tuples {
		for j := range tuples {
			if (*tuples[i]).Score > (*tuples[j]).Score {
				tmp := tuples[i]
				tuples[i] = tuples[j]
				tuples[j] = tmp
			}
		}
	}

	return tuples
}

func sortByPriority(tuples []*model.PriorityTuple) []*model.PriorityTuple {
	for i := range tuples {
		for j := range tuples {
			if (*tuples[i]).Priority > (*tuples[j]).Priority {
				tmp := tuples[i]
				tuples[i] = tuples[j]
				tuples[j] = tmp
			}
		}
	}

	return tuples
}

func recallSrcSortByPriority(list []*v1.RecallInfo, limit int) []*v1.RecallInfo {
	for i := range list {
		for j := range list {
			if list[i].Priority > list[j].Priority {
				tmp := list[i]
				list[i] = list[j]
				list[j] = tmp
			}
		}
	}
	return list[:limit]
}
