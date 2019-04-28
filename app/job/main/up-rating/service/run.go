package service

import (
	"context"
	"time"

	"go-common/app/job/main/up-rating/model"

	"go-common/library/log"
	"golang.org/x/sync/errgroup"
)

// Run run scores
func (s *Service) Run(c context.Context, date time.Time) (err error) {
	date = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	// get weight parameters
	params, err := s.getAllParamter(c)
	if err != nil {
		log.Error("s.getAllParamter error(%v)", err)
		return
	}
	// get past scores
	past, err := s.pastInfos(c)
	if err != nil {
		log.Error("s.pastInfos error(%v)", err)
		return
	}
	err = s.PrepareData(c, date, past, params)
	if err != nil {
		log.Error("s.PrepareDate error(%v)", err)
		return
	}

	err = s.CalScores(c, date, past, params)
	if err != nil {
		log.Error("s.CalScores error(%v)", err)
		return
	}

	err = s.delOldPastInfo(c, int64(_limit))
	if err != nil {
		log.Error("s.delOldPastInfo error(%v)", err)
		return
	}

	err = s.insertPastRecord(c, 0, date.AddDate(0, 1, 0).Format(_layout))
	if err != nil {
		log.Error("s.insertPastRecord error(%v)", err)
		return
	}

	log.Info("run data read finished")
	err = s.InsertTaskStatus(c, 0, 1, date.Format(_layout), "ok")
	return
}

// PrepareData prepare old data
func (s *Service) PrepareData(c context.Context, date time.Time, past map[int64]*model.Past, params *model.RatingParameter) (err error) {
	var (
		g         errgroup.Group
		routines  = s.conf.Con.Concurrent
		lastMonth = time.Date(date.Year(), date.Month()-1, 1, 0, 0, 0, 0, time.Local)
	)

	offset, end, _, err := s.RatingOffEnd(c, lastMonth)
	if err != nil {
		return
	}
	total := end - offset
	section := (total - total%routines) / routines
	for i := 0; i < routines; i++ {
		begin := section*i + offset
		over := begin + section
		if i == routines-1 {
			over = end
		}

		// read chan: for last ratings
		rch := make(chan []*model.Rating, _limit)
		g.Go(func() (err error) {
			err = s.RatingFast(c, lastMonth, begin, over, rch)
			if err != nil {
				log.Error("s.RatingFast error(%v)", err)
			}
			return
		})

		wch := make(chan []*model.Rating, _limit)
		g.Go(func() (err error) {
			s.Copy(rch, wch, past, params)
			return
		})

		g.Go(func() (err error) {
			err = s.BatchInsertRatingStat(c, wch, date)
			if err != nil {
				log.Error("s.BatchInsertRatingStat error(%v)", err)
			}
			return
		})
	}

	if err = g.Wait(); err != nil {
		log.Error("run g.Wait error(%v)", err)
		return
	}
	return
}

// CalScores cal scores
func (s *Service) CalScores(c context.Context, date time.Time, past map[int64]*model.Past, params *model.RatingParameter) (err error) {
	var (
		routines = s.conf.Con.Concurrent
		_limit   = s.conf.Con.Limit
	)
	t := time.Now().UnixNano()
	var g errgroup.Group
	// get id start:end by date
	offset, end, err := s.BaseInfoOffEnd(c, date)
	if err != nil {
		return
	}
	total := end - offset
	section := (total - total%routines) / routines
	// parallelization and pipeling
	for i := 0; i < routines; i++ {
		begin := section*i + offset
		over := begin + section
		if i == routines-1 {
			over = end
		}
		// read chan: for origin datas
		rch := make(chan []*model.BaseInfo, _limit)
		g.Go(func() (err error) {
			err = s.BaseInfo(c, date, begin, over, rch)
			if err != nil {
				log.Error("s.BaseInfo error(%v)", err)
			}
			return
		})

		// write chan: for calculated results
		wch := make(chan []*model.Rating, _limit)
		g.Go(func() (err error) {
			s.CalScore(rch, wch, params, past, date)
			return
		})

		g.Go(func() (err error) {
			err = s.BatchInsertRatingStat(c, wch, date)
			if err != nil {
				log.Error("s.BatchInsertRatingStat error(%v)", err)
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("run g.Wait error(%v)", err)
		return
	}

	log.Info("cal time cost:", time.Now().UnixNano()-t)
	return
}
