package service

import (
	"context"
	"sort"
	"sync"

	"go-common/app/job/main/figure-timer/model"
	"go-common/library/log"
)

var (
	rank = &scores{}
)

type scores struct {
	ss []int32
	mu sync.Mutex
}

func (s *scores) AddScore(score int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ss = append(s.ss, score)
}

func (s *scores) Init() {
	s.ss = make([]int32, 0)
}

func (s *scores) Sort() {
	sort.Sort(s)
}

func (s *scores) Get(i int) int32 {
	return s.ss[i]
}

func (s *scores) Len() int {
	return len(s.ss)
}

func (s *scores) Less(i1, i2 int) bool {
	return s.ss[i1] > s.ss[i2]
}

func (s *scores) Swap(i1, i2 int) {
	s.ss[i1], s.ss[i2] = s.ss[i2], s.ss[i1]
}

func (s *Service) calcRank(c context.Context, ver int64) {
	if rank.Len() <= 0 {
		return
	}
	rank.Sort()
	var (
		scoreFrom, scoreTo int32
		err                error
	)
	scoreTo = rank.Get(0)
	for i := int8(1); i <= 100; i++ {
		var index int
		if i == 100 {
			scoreFrom = 0
		} else {
			index = int(float64(rank.Len()) * (float64(i) / 100.0))
			scoreFrom = rank.Get(index)
		}
		if scoreFrom > scoreTo {
			scoreTo = scoreFrom
		}
		r := &model.Rank{ScoreFrom: scoreFrom, ScoreTo: scoreTo, Percentage: i, Ver: ver}
		if _, err = s.dao.InsertRankHistory(c, r); err != nil {
			log.Error("%+v", err)
			return
		}
		if _, err = s.dao.UpsertRank(c, r); err != nil {
			log.Error("%+v", err)
			return
		}
		scoreTo = scoreFrom - 1
	}
}
