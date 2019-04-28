package service

import (
	"context"
	"time"

	"go-common/app/interface/main/up-rating/model"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_cdateLayout  = "2006-01-02"
	_taskFinished = 1
)

// UpRating gets data from cache and falls back to db
func (s *Service) UpRating(c context.Context, mid int64) (rating *model.Rating, err error) {
	var allow bool
	if allow, err = s.allow(c, mid); err != nil {
		return
	}
	if !allow {
		err = ecode.UpRatingNoPermission
		return
	}
	var redisErr error
	if rating, redisErr = s.dao.GetUpRatingCache(c, mid); redisErr != nil {
		log.Error("s.dao.GetUpRatingCache error(%v)", redisErr)
	}
	if rating == nil {
		if rating, err = s.UpRatingDB(c, mid); err != nil {
			return
		}
		if redisErr = s.dao.SetUpRatingCache(c, mid, rating); redisErr != nil {
			log.Error("s.dao.SetUpRatingCache error(%v)", redisErr)
		}
	}
	return
}

// UpRatingDB gets data from db
func (s *Service) UpRatingDB(c context.Context, mid int64) (rating *model.Rating, err error) {
	rating = &model.Rating{}
	if rating.Score, err = s.latestScore(c, mid); err != nil {
		return
	}
	if rating.Score.Magnetic < model.LowerBoundScore {
		return nil, ecode.UpRatingScoreLimit
	}
	if rating.Rank, err = s.rank(rating.Score); err != nil {
		log.Error("s.rank error(%v)", err)
		return
	}
	if rating.Prize, err = s.prize(c, rating.Score); err != nil {
		log.Error("s.prize error(%v)", err)
		return
	}
	// privileges left for furture development
	rating.Privileges = make([]*model.Privilege, 0)
	return
}

func (s *Service) rank(score *model.Score) (rank *model.Rank, err error) {
	for _, r := range model.Ranks {
		if score.Magnetic >= r.Score() {
			rank = r.Rank()
			return
		}
	}
	rank = model.RankLevelNone.Rank()
	return
}

func (s *Service) prize(c context.Context, score *model.Score) (prize *model.Prize, err error) {
	pass := false
	m := map[model.PrizeLevel]func(){
		model.PrizeLevelOne: func() {
			if pass = score.Magnetic >= model.TotalScore*0.95; pass {
				prize = model.PrizeLevelOne.Prize()
			}
		},
		model.PrizeLevelTwo: func() {
			tmp := score
			if tmp.Magnetic < model.TotalScore*0.75 {
				return
			}
			for j := 2; j > 0; j-- {
				if tmp, err = s.score(c, tmp.MID, prevRatingTime(tmp.CDate)); err != nil || tmp == nil {
					return
				}
				if tmp.Magnetic < model.TotalScore*0.75 {
					return
				}
			}
			pass = true
			prize = model.PrizeLevelTwo.Prize()
		},
		model.PrizeLevelThree: func() {
			var prev *model.Score
			if prev, err = s.score(c, score.MID, prevRatingTime(score.CDate)); err != nil || prev == nil {
				return
			}
			v := score.Magnetic - prev.Magnetic
			if pass = v >= model.TotalScore*0.1; pass {
				prize = model.PrizeLevelThree.Prize(v)
			}
		},
		model.PrizeLevelFour: func() {
			g := func(i int) bool {
				return i == 233 || i == 223
			}
			if pass = g(score.Credit) || g(score.Creative) || g(score.Influence); pass {
				prize = model.PrizeLevelFour.Prize()
			}
		},
		model.PrizeLevelFive: func() {
			pass = true
			prize = model.PrizeLevelFive.Prize()
		},
	}

	for _, p := range model.Prizes {
		m[p]()
		if err != nil || pass {
			break
		}
	}

	return
}

func (s *Service) latestScore(c context.Context, mid int64) (score *model.Score, err error) {
	var status int
	ratingTime := prevRatingTime(time.Now())
	if status, err = s.dao.TaskStatus(c, ratingTime.Format(_cdateLayout)); err != nil {
		return
	}
	if status != _taskFinished {
		ratingTime = prevRatingTime(ratingTime)
		if status, err = s.dao.TaskStatus(c, ratingTime.Format(_cdateLayout)); err != nil {
			return
		}
		if status != _taskFinished {
			err = ecode.UpRatingNoData
			return
		}
	}
	return s.score(c, mid, ratingTime)
}

func (s *Service) score(c context.Context, mid int64, ratingTime time.Time) (score *model.Score, err error) {
	score, err = s.dao.UpScore(c, int(ratingTime.Month()), mid, ratingTime.Format(_cdateLayout))
	if err != nil {
		log.Error("s.dao.UpScore err(%v)", err)
		return
	}
	if score == nil {
		log.Error("s.dao.UpScore mid(%v) ratingTime(%v) not found", mid, ratingTime.String())
		return
	}
	score.Magnetic = score.Creative + score.Credit + score.Influence
	score.StatEnd = time.Date(score.CDate.Year(), score.CDate.Month()+1, 1, 0, 0, 0, 0, score.CDate.Location()).AddDate(0, 0, -1)
	score.StatStart = time.Date(score.CDate.Year()-1, score.CDate.Month()+1, 1, 0, 0, 0, 0, score.CDate.Location())
	return
}

func prevRatingTime(queryTime time.Time) time.Time {
	return time.Date(queryTime.Year(), queryTime.Month()-1, 1, 0, 0, 0, 0, queryTime.Location())
}

func (s *Service) allow(c context.Context, mid int64) (b bool, err error) {
	count, err := s.dao.White(c, mid)
	if err != nil {
		return
	}
	b = count != 0
	return
}
