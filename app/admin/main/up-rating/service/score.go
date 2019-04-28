package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-common/app/admin/main/up-rating/dao/global"
	"go-common/app/admin/main/up-rating/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

// RatingList returns rating info list
func (s *Service) RatingList(c context.Context, arg *model.RatingListArg, date time.Time) (res []*model.RatingInfo, total int64, err error) {
	var (
		cdate    = getStartMonthlyDate(date)
		cdateStr = cDateStr(cdate)
		mon      = int(cdate.Month())
	)
	res = make([]*model.RatingInfo, 0)
	q := scoreListQuery(arg)
	total = 10000
	// if total, err = s.dao.Total(c, mon, cdateStr, q); err != nil {
	// 	log.Error("s.dao.Total error(%v)", err)
	// 	return
	// }
	// if total == 0 {
	// 	return
	// }
	// scores
	q += fmt.Sprintf(" ORDER BY id ASC LIMIT %d, %d", arg.From, arg.Limit)
	if res, err = s.dao.ScoreList(c, mon, cdateStr, q); err != nil {
		log.Error("s.dao.ScoreList error(%v)", err)
		return
	}
	if len(res) == 0 {
		return
	}
	var (
		m    = make(map[int64]*model.RatingInfo, len(res))
		mids = make([]int64, 0, len(res))
	)
	for _, v := range res {
		mids = append(mids, v.Mid)
		m[v.Mid] = v
	}
	// fans and avs
	levelInfos, err := s.dao.LevelList(c, mon, cdateStr, mids)
	if err != nil {
		log.Error("s.dao.LevelList error(%v)", err)
		return
	}
	for _, v := range levelInfos {
		score := m[v.Mid]
		score.TotalAvs = v.TotalAvs
		score.TotalFans = v.TotalFans
	}
	// nicknames
	r, err := global.Names(c, mids)
	if err != nil {
		log.Error("global.Names error(%v)", err)
		return
	}
	for _, v := range res {
		v.NickName = r[v.Mid]
		v.Date = v.ScoreDate.Format("2006-01")
	}
	return
}

func scoreListQuery(arg *model.RatingListArg) (where string) {
	if arg.Mid > 0 {
		where += fmt.Sprintf(" AND mid=%d", arg.Mid)
	}
	if arg.ScoreMin > 0 {
		where += fmt.Sprintf(" AND %s>=%d", scoreField(arg.ScoreType), arg.ScoreMin)
	}
	if arg.ScoreMax > 0 {
		where += fmt.Sprintf(" AND %s<%d", scoreField(arg.ScoreType), arg.ScoreMax)
	}
	if len(arg.Tags) > 0 {
		// if len(arg.Tags) > 1 {
		// 	s := rand.NewSource(time.Now().Unix())
		// 	r := rand.New(s) // initialize local pseudorandom generator
		// 	i := r.Intn(len(arg.Tags))
		// 	where += fmt.Sprintf(" AND tag_id=%d", arg.Tags[i])
		// } else {
		// 	where += fmt.Sprintf(" AND tag_id=%d", arg.Tags[0])
		// }
		where += fmt.Sprintf(" AND tag_id IN (%s)", xstr.JoinInts(arg.Tags))
	}
	where += " AND is_deleted=0"

	return
}

// ScoreCurrent returns current rating info
func (s *Service) ScoreCurrent(c context.Context, mid int64) (res *model.ScoreCurrentResp, err error) {
	current, err := s.lastestScore(c, mid)
	if err != nil || current == nil {
		return
	}
	res = &model.ScoreCurrentResp{
		Date:       current.ScoreDate.Unix(),
		Credit:     &model.ScoreCurrent{Current: current.CreditScore, Diff: current.CreditScore},
		Influence:  &model.ScoreCurrent{Current: current.InfluenceScore, Diff: current.InfluenceScore},
		Creativity: &model.ScoreCurrent{Current: current.CreativityScore, Diff: current.CreativityScore},
	}
	prev, err := s.upScore(c, mid, prevComputation(current.ScoreDate))
	if err != nil {
		return
	}
	if prev != nil {
		res.Credit.Diff = current.CreditScore - prev.CreditScore
		res.Influence.Diff = current.InfluenceScore - prev.InfluenceScore
		res.Creativity.Diff = current.CreativityScore - prev.CreativityScore
	}
	return
}

func (s *Service) upScoreHistory(c context.Context, mid int64, queryAll bool, count int) ([]*model.RatingInfo, error) {
	if queryAll {
		return s.upPastScoresAll(c, mid)
	}
	return s.upPastScores(c, mid, count)
}

// ScoreHistory returns score history
func (s *Service) ScoreHistory(c context.Context, types []model.ScoreType, mid int64, queryAll bool, limit int) (res []*model.UpScoreHistory, err error) {
	res = make([]*model.UpScoreHistory, 0)
	history, err := s.upScoreHistory(c, mid, queryAll, limit)
	if err != nil || len(history) <= 0 {
		return
	}
	sort.Slice(history, func(i, j int) bool {
		return history[i].ScoreDate.Before(history[j].ScoreDate)
	})

	var (
		dates = make([]int64, 0, len(history))
		m     = map[model.ScoreType][]int64{
			model.Credit:     make([]int64, 0, len(history)),
			model.Creativity: make([]int64, 0, len(history)),
			model.Influence:  make([]int64, 0, len(history)),
		}
	)
	for _, v := range history {
		dates = append(dates, v.ScoreDate.Unix())
		m[model.Creativity] = append(m[model.Creativity], v.CreativityScore)
		m[model.Credit] = append(m[model.Credit], v.CreditScore)
		m[model.Influence] = append(m[model.Influence], v.InfluenceScore)
	}
	for _, t := range types {
		res = append(res, &model.UpScoreHistory{
			ScoreType: t,
			Date:      dates,
			Score:     m[t],
		})
	}
	return
}

// ExportScores exports scores
func (s *Service) ExportScores(ctx context.Context, arg *model.RatingListArg, date time.Time) (res []byte, err error) {
	ratings, _, err := s.RatingList(ctx, arg, date)
	if err != nil {
		log.Error("s.RatingList error(%v)", err)
		return
	}
	data := formatScores(ratings)
	res, err = formatCSV(data)
	if err != nil {
		log.Error("up-rating FormatCSV error(%v)", err)
	}
	return
}

func (s *Service) upPastScoresAll(c context.Context, mid int64) (res []*model.RatingInfo, err error) {
	res = make([]*model.RatingInfo, 0)
	for f := 1; f <= 12; f++ {
		var list []*model.RatingInfo
		if list, err = s.dao.UpScores(c, f, mid); err != nil {
			log.Error("s.dao.UpScores error(%v)", err)
			return
		}
		res = append(res, list...)
	}
	return
}

func (s *Service) upPastScores(c context.Context, mid int64, count int) (res []*model.RatingInfo, err error) {
	res = make([]*model.RatingInfo, 0)
	var lastScore *model.RatingInfo
	if lastScore, err = s.lastestScore(c, mid); err != nil || lastScore == nil {
		return
	}
	res = append(res, lastScore)
	for f := 1; f < count; f++ {
		var v *model.RatingInfo
		if v, err = s.upScore(c, mid, lastScore.ScoreDate.AddDate(0, -f, 0)); err != nil {
			return
		}
		if v != nil {
			res = append(res, v)
		}
	}

	return
}

func (s *Service) lastestScore(c context.Context, mid int64) (score *model.RatingInfo, err error) {
	cdate := prevComputation(time.Now())
	var b bool
	if b, err = s.taskFinished(c, cdate); err != nil {
		return
	}
	if !b {
		cdate = prevComputation(cdate)
		if b, err = s.taskFinished(c, cdate); err != nil {
			return
		}
	}
	if !b {
		log.Error("s.latestScore cdate(%s) no data available", cdate)
		err = ecode.ServerErr
		return
	}
	return s.upScore(c, mid, cdate)
}

func (s *Service) taskFinished(c context.Context, cdate time.Time) (bool, error) {
	str := cDateStr(cdate)
	status, err := s.dao.TaskStatus(c, str)
	if err != nil {
		log.Error("s.dao.TaskStasus date(%s) error(%v)", str, err)
		return false, err
	}
	return status == 1, nil
}

// upScore wraps dao.upScore
func (s *Service) upScore(c context.Context, mid int64, t time.Time) (*model.RatingInfo, error) {
	score, err := s.dao.UpScore(c, int(t.Month()), mid, cDateStr(t))
	if err != nil {
		log.Error("s.dao.UpScore mid(%d) date(%s) error(%v)", mid, t, err)
		return nil, err
	}
	return score, nil
}
