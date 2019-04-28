package feed

import (
	"go-common/app/interface/main/app-card/model/card/rank"
	"go-common/app/interface/main/app-intl/model"
)

// RankCard is.
func (s *Service) RankCard(plat int8) (ranks []*rank.Rank, aids []int64) {
	var limit int
	if !model.IsIPad(plat) {
		limit = 3
	} else {
		limit = 4
	}
	ranks = make([]*rank.Rank, 0, limit)
	aids = make([]int64, 0, limit)
	for _, rank := range s.rankCache {
		ranks = append(ranks, rank)
		aids = append(aids, rank.Aid)
		if len(ranks) == limit {
			break
		}
	}
	return
}
