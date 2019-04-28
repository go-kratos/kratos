package service

import (
	"context"
	"encoding/binary"

	"go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	"go-common/app/service/bbq/recsys-recall/dao"
	"go-common/app/service/bbq/recsys-recall/model"
	"go-common/app/service/bbq/recsys-recall/service/index"
	"go-common/library/log"

	jsoniter "github.com/json-iterator/go"
)

// RecallTask .
type RecallTask struct {
	ctx    *context.Context
	d      *dao.Dao
	mid    int64
	buvid  string
	info   *v1.RecallInfo
	sc     *ScorerManager
	ranker *RankerManager
	filter *FilterManager
	debug  bool
}

func newRecallTask(ctx *context.Context, d *dao.Dao, mid int64, buvid string, info *v1.RecallInfo) *RecallTask {
	return &RecallTask{
		ctx:   ctx,
		d:     d,
		mid:   mid,
		buvid: buvid,
		info:  info,
		debug: false,
	}
}

// SetScorerManager .
func (t *RecallTask) SetScorerManager(sc *ScorerManager) {
	t.sc = sc
}

// SetRankerManager .
func (t *RecallTask) SetRankerManager(ranker *RankerManager) {
	t.ranker = ranker
}

// SetFilterManager .
func (t *RecallTask) SetFilterManager(filter *FilterManager) {
	t.filter = filter
}

// SetDebug .
func (t *RecallTask) SetDebug(d bool) {
	t.debug = d
}

// Run .
func (t *RecallTask) Run() *[]byte {
	// 获取倒排
	raw, err := t.d.GetInvertedIndex(*t.ctx, t.info.Tag, false)
	if err != nil {
		log.Errorv(*t.ctx, log.KV("Tag", t.info.Tag), log.KV("redis", err))
		return nil
	}
	var recallList, svidList []uint64
	if len(raw) > 13 && binary.BigEndian.Uint64(raw[:8]) == 0xffffffffdeadbeef {
		ii := new(index.InvertedIndex)
		if err = ii.Load(raw); err != nil {
			log.Errorv(*t.ctx, log.KV("Tag", t.info.Tag), log.KV("inverted index load", err))
			return nil
		}
		recallList = ii.Data
	} else {
		if err = jsoniter.Unmarshal(raw, &recallList); err != nil {
			log.Errorv(*t.ctx, log.KV("Tag", t.info.Tag), log.KV("jsoninter", err))
			return nil
		}
	}

	// filter
	for _, svid := range recallList {
		if !t.filter.DoFilter(svid, "default", t.info.Filter) {
			// if !t.filter.DoFilter(svid, t.info.Filter) {
			svidList = append(svidList, svid)
		}
	}

	tuples := make([]*model.Tuple, len(svidList))
	for i, v := range svidList {
		// score
		score := t.sc.DoScore(v, t.info.Scorer)
		tuples[i] = &model.Tuple{
			Svid:  v,
			Score: score,
		}
	}

	// rank
	t.ranker.DoRank(&tuples, t.info.Ranker, defaultCompare)

	// truncate
	size := int(t.info.Limit)
	if size > len(tuples) {
		size = len(tuples)
	}

	result := &Result{
		TotalHit:    int32(len(recallList)),
		FilterCount: int32(len(svidList)),
		FinalCount:  int32(size),
		Tuples:      tuples[:size],
	}

	log.Infov(*t.ctx, log.KV("req_tag", t.info.Tag), log.KV("req_name", t.info.Name), log.KV("recall", len(recallList)), log.KV("filter", len(svidList)), log.KV("result", len(result.Tuples)))

	return result.ToBytes()
}
