package service

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/up-rating/model"

	"go-common/library/log"
	"golang.org/x/sync/errgroup"
)

var (
	_layout = "2006-01-02"
	_limit  = 2000
)

// RunPastScore run past score by date
func (s *Service) RunPastScore(c context.Context, date time.Time) (err error) {
	date = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	times, err := s.getPastRecord(c, date.Format(_layout))
	if err != nil {
		log.Error("s.getPastRecord error(%v)", err)
		return
	}
	if times < 0 {
		log.Info("This month's calculation did not start")
		return
	}
	// 创作力需要计算前22个月的数据
	if times >= 22 {
		log.Info("Last month's calculation has end")
		return
	}
	var (
		readGroup errgroup.Group
		cw        float64 // 创作力当月权重
		iw        int64   // 影响力当月权重
		pastScore []*model.Past
		pastCh    = make(chan []*model.Rating, _limit)
	)
	// 获取前n个月的数据
	pastDate := date.AddDate(0, -1*(22-times), 0)
	times++ // update calculate times

	//csr = csm0 + csm1 + ... + csm11 + 11/12 * csm12 + 10/12 * csm13 + ... 1/12 * csm22
	cw = float64(times) / float64(12)
	if cw > 1.0 {
		cw = 1.0
	}
	// isr = mfans0 + mfans1 + ... + mfans12
	iw = int64(float64(times) / float64(12))
	// get past month data
	readGroup.Go(func() (err error) {
		err = s.RatingInfos(c, pastDate, pastCh)
		if err != nil {
			log.Error("s.RatingInfos error(%v)", err)
		}
		return
	})
	// cal past month data
	readGroup.Go(func() (err error) {
		pastScore, err = s.calPastScores(c, pastCh, cw, iw)
		if err != nil {
			log.Error("s.calPastScores error(%v)", err)
		}
		return
	})
	if err = readGroup.Wait(); err != nil {
		log.Error("run readGroup.Wait error(%v)", err)
		return
	}

	err = s.insertPastRecord(c, times, date.Format(_layout))
	if err != nil {
		log.Error("s.upPastRecord error(%v)", err)
		return
	}

	err = s.batchInsertPastScore(c, pastScore)
	if err != nil {
		log.Error("s.batchInsertPastScore error(%v)", err)
	}
	return
}

// InsertPastRecord insert past record
func (s *Service) InsertPastRecord(c context.Context, date string) (err error) {
	return s.insertPastRecord(c, 0, date)
}

func (s *Service) calPastScores(c context.Context, pastRating chan []*model.Rating, cw float64, iw int64) (pastScore []*model.Past, err error) {
	pastScore = make([]*model.Past, 0)
	for rating := range pastRating {
		p := calPastScore(rating, cw, iw)
		pastScore = append(pastScore, p...)
	}
	return
}

func calPastScore(rating []*model.Rating, cw float64, iw int64) (pastScore []*model.Past) {
	pastScore = make([]*model.Past, 0, len(rating))
	for _, r := range rating {
		pastScore = append(pastScore, &model.Past{
			MID:                 r.MID,
			MetaCreativityScore: int64(float64(r.MetaCreativityScore) * cw),
			MetaInfluenceScore:  r.MetaInfluenceScore * iw,
			CreditScore:         r.CreditScore,
		})
	}
	return
}

// get past calculate record
func (s *Service) getPastRecord(c context.Context, date string) (times int, err error) {
	return s.dao.GetPastRecord(c, date)
}

func (s *Service) insertPastRecord(c context.Context, times int, date string) (err error) {
	_, err = s.dao.InsertPastRecord(c, times, date)
	return err
}

func (s *Service) pastInfos(c context.Context) (past map[int64]*model.Past, err error) {
	past = make(map[int64]*model.Past)
	var id int64
	for {
		var p []*model.Past
		p, id, err = s.dao.GetPasts(c, id, int64(_limit))
		if err != nil {
			return
		}
		for i := 0; i < len(p); i++ {
			past[p[i].MID] = p[i]
		}
		if len(p) < _limit {
			break
		}
	}
	return
}

func (s *Service) batchInsertPastScore(c context.Context, past []*model.Past) (err error) {
	var (
		buff    = make([]*model.Past, 2000)
		buffEnd = 0
	)
	for _, p := range past {
		buff[buffEnd] = p
		buffEnd++
		if buffEnd >= 2000 {
			values := assemblePastValues(buff[:buffEnd])
			buffEnd = 0
			_, err = s.dao.InsertPastScoreStat(c, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := assemblePastValues(buff[:buffEnd])
		buffEnd = 0
		_, err = s.dao.InsertPastScoreStat(c, values)
	}
	return
}

func assemblePastValues(past []*model.Past) (values string) {
	var buf bytes.Buffer
	for _, p := range past {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(p.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(p.MetaCreativityScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(p.MetaInfluenceScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(p.CreditScore, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}

func (s *Service) delOldPastInfo(c context.Context, limit int64) (err error) {
	var rows int64
	for {
		rows, err = s.dao.DelPastStat(c, limit)
		if err != nil {
			return
		}
		if rows < limit {
			break
		}
	}
	return
}

// DelPastRecord del past record
func (s *Service) DelPastRecord(c context.Context, date time.Time) (err error) {
	_, err = s.dao.DelPastRecord(c, date)
	return
}
