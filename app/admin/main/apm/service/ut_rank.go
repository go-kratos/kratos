package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// RankTen find data and avg(coverage)
func (s *Service) RankTen(c context.Context, order string) (ranks []*ut.RankResp, err error) {
	s.ranksCache.Lock()
	allRanks := s.ranksCache.Slice
	s.ranksCache.Unlock()
	ranksLen := len(allRanks)
	if ranksLen == 0 {
		log.Info("s.RankTen order(%s) ranks is empty", order)
		return
	}
	if ranksLen > 10 {
		ranksLen = 10
	}
	if order == "desc" {
		sort.Slice(allRanks, func(i, j int) bool { return allRanks[i].Score > allRanks[j].Score })
	} else {
		sort.Slice(allRanks, func(i, j int) bool { return allRanks[i].Score < allRanks[j].Score })
	}
	ranks = append(ranks, allRanks[0:ranksLen]...)
	return
}

// UserRank return one's rank
func (s *Service) UserRank(c context.Context, username string) (rank *ut.RankResp, err error) {
	s.ranksCache.Lock()
	allRanks := s.ranksCache.Map
	s.ranksCache.Unlock()
	if len(allRanks) == 0 {
		log.Info("s.RankTen(%s) ranks is empty", username)
		return
	}
	rank = allRanks[username]
	return
}

// RanksCache flush cache for ranks.
func (s *Service) RanksCache(c context.Context) (err error) {
	s.ranksCache.Lock()
	defer s.ranksCache.Unlock()
	var (
		rankSlice []*ut.RankResp
		rankMap   = make(map[string]*ut.RankResp)
		endTime   = time.Now().AddDate(0, -3, 0).Format("2006-01-02 15:04:05")
	)
	if err = s.DB.Table("ut_merge").Raw("select * from (select ut_commit.username,ROUND(avg(coverage)/100,2) AS coverage,ROUND(SUM(passed)/SUM(assertions),2)*100 AS pass_rate, SUM(assertions)AS assertions, SUM(passed) AS passed, MAX(ut_commit.mtime) as mtime from ut_merge,ut_commit,ut_pkganls where ut_merge.merge_id=ut_commit.merge_id and ut_commit.commit_id=ut_pkganls.commit_id and ut_merge.is_merged=1 and ut_merge.mtime >=? and (ut_pkganls.pkg!=substring_index(ut_pkganls.pkg,'/',5) or ut_pkganls.pkg like 'go-common/library/%') GROUP BY ut_commit.username) as t1", endTime).
		Find(&rankSlice).Error; err != nil {
		log.Error("RankResult error(%v)", err)
		return
	}
	total := len(rankSlice)
	for _, rank := range rankSlice {
		rank.Total = total
		rank.Newton = NewtonScore(rank.Mtime)
		Score := rank.Coverage * WilsonScore(rank.Passed, rank.Assertions) * rank.Newton
		rank.Score, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", Score), 64)
	}
	// desc
	sort.Slice(rankSlice, func(i, j int) bool { return rankSlice[i].Score > rankSlice[j].Score })
	for index, rank := range rankSlice {
		rank.Rank = index + 1
		rank.AvatarURL, _ = s.dao.GitLabFace(c, rank.UserName)
		rankMap[rank.UserName] = rank
		if history, ok := s.ranksCache.Map[rank.UserName]; ok {
			rank.Change = history.Rank - rank.Rank
		}
	}
	s.ranksCache.Map = rankMap
	s.ranksCache.Slice = rankSlice
	return
}

//WilsonScore  Wilson-score-interval
func WilsonScore(pos int, total int) (score float64) {
	z := float64(1.96)
	posRat := float64(pos) / float64(total)
	score = (posRat + math.Pow(z, 2)/(2*float64(total))) / (1 + math.Pow(z, 2)/float64(total))
	return
}

// NewtonScore .  Newton's Law of Cooling
// 冷却因子经我们需求计算 a = 0.02
// deltaDays<7 day -a*x^2+b  >7 NewTonScore
func NewtonScore(maxMtime xtime.Time) (score float64) {
	deltaDays := int(time.Since(maxMtime.Time()).Hours() / 24)
	if deltaDays <= int(7) {
		score = -(float64(0.0012))*math.Pow(float64(deltaDays), 2) + float64(1)
		return
	}
	score = math.Exp(-(float64(0.02) * float64(deltaDays)))
	return
}
