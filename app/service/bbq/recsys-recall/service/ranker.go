package service

import (
	"go-common/app/service/bbq/recsys-recall/dao"
	"go-common/app/service/bbq/recsys-recall/model"
)

// Ranker interface
type Ranker interface {
	doRank(tuples *[]*model.Tuple, comp func(interface{}, interface{}) bool)
}

// RankerManager .
type RankerManager struct {
	rankers map[string]Ranker
}

// NewRankerManager .
func NewRankerManager(d *dao.Dao) *RankerManager {
	r := make(map[string]Ranker)

	r["default"] = &DefaultRanker{
		d: d,
	}

	return &RankerManager{
		rankers: r,
	}
}

// DoRank .
func (rm *RankerManager) DoRank(tuples *[]*model.Tuple, name string, comp func(interface{}, interface{}) bool) {
	if r, ok := rm.rankers[name]; ok && r != nil {
		r.doRank(tuples, comp)
	}
}

// DefaultRanker .
type DefaultRanker struct {
	d *dao.Dao
}

func defaultCompare(a, b interface{}) bool {
	t1 := a.(model.Tuple)
	t2 := b.(model.Tuple)

	return t1.Score > t2.Score
}

func (dr *DefaultRanker) doRank(tuples *[]*model.Tuple, comp func(interface{}, interface{}) bool) {
	if len(*tuples) <= 0 {
		return
	}
	for i, u := range *tuples {
		for j, v := range *tuples {
			if comp(u, v) {
				tmp := (*tuples)[i]
				(*tuples)[i] = (*tuples)[j]
				(*tuples)[j] = tmp
			}
		}
	}
}
