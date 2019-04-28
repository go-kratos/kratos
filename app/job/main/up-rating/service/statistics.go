package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/up-rating/model"

	"go-common/library/log"
	xtime "go-common/library/time"
	"golang.org/x/sync/errgroup"
)

const (
	// TotalType total type
	TotalType int = iota
	// CreativeType creative type
	CreativeType
	// InfluenceType influence type
	InfluenceType
	// CreditType Credit type
	CreditType
)

var (
	_offset         int64 = 10
	_sumTotal       int64
	_creativeTotal  int64
	_influenceTotal int64
	_creditTotal    int64
)

func initSection(min, max, section int64, ctype int, tagID int64, date xtime.Time) *model.RatingStatis {
	return &model.RatingStatis{
		Section: section,
		Tips:    fmt.Sprintf("\"%d-%d\"", min, max),
		TagID:   tagID,
		CDate:   date,
		CType:   ctype,
	}
}

func initSections(totalScore, tagID int64, ctype int, date xtime.Time) (statis []*model.RatingStatis) {
	var (
		idx int64
	)
	statis = make([]*model.RatingStatis, totalScore/_offset)
	for idx*_offset < totalScore {
		statis[idx] = initSection(idx*_offset, (idx+1)*_offset, idx, ctype, tagID, date)
		idx++
	}
	return
}

// RunStatistics run up rating statistics
func (s *Service) RunStatistics(c context.Context, date time.Time) (err error) {
	err = s.initTotalScore(c)
	if err != nil {
		log.Error("s.initTotalScore error(%v)", err)
		return
	}
	err = s.delStatistics(c, date)
	if err != nil {
		log.Error("s.delStatistics error(%v)", err)
		return
	}
	err = s.statistics(c, date)
	if err != nil {
		log.Error("s.scoreStatistics error(%v)", err)
	}
	return
}

func (s *Service) initTotalScore(c context.Context) (err error) {
	params, err := s.getAllParamter(c)
	if err != nil {
		log.Error("s.getAllParamter error(%v)", err)
		return
	}
	_sumTotal = params.WCSR + params.HR + params.WISR
	_creativeTotal = params.WCSR
	_influenceTotal = params.WISR
	_creditTotal = params.HR
	return
}

func (s *Service) delStatistics(c context.Context, date time.Time) (err error) {
	err = s.delRatingCom(c, "up_rating_statistics", date)
	if err != nil {
		return
	}
	err = s.delRatingCom(c, "up_rating_top", date)
	return
}

// delRatingCom del com
func (s *Service) delRatingCom(c context.Context, table string, date time.Time) (err error) {
	for {
		var rows int64
		rows, err = s.dao.DelRatingCom(c, table, date, _limit)
		if err != nil {
			return
		}
		if rows == 0 {
			break
		}
	}
	return
}

func (s *Service) statistics(c context.Context, date time.Time) (err error) {
	var (
		readGroup errgroup.Group
		sourceCh  = make(chan []*model.Rating, _limit)
		statisCh  = make(chan []*model.Rating, _limit)
		topCh     = make(chan []*model.Rating, _limit)

		sections  map[int]map[int64][]*model.RatingStatis
		topRating map[int]map[int64]*RatingHeap
	)
	baseInfo, err := s.BaseTotal(c, date)
	if err != nil {
		log.Error("s.BaseTotal error(%v)", err)
		return
	}
	// get rating info
	readGroup.Go(func() (err error) {
		err = s.RatingInfos(c, date, sourceCh)
		if err != nil {
			log.Error("s.RatingInfos error(%v)", err)
		}
		return
	})

	// dispatch
	readGroup.Go(func() (err error) {
		defer func() {
			close(topCh)
			close(statisCh)
		}()
		for rating := range sourceCh {
			statisCh <- rating
			topCh <- rating
		}
		return
	})

	// top
	readGroup.Go(func() (err error) {
		topRating, err = s.ratingTop(c, date, topCh)
		if err != nil {
			log.Error("s.RatingTop error(%v)", err)
		}
		return
	})
	// statis
	readGroup.Go(func() (err error) {
		sections, err = s.scoreStatistics(c, date, statisCh, baseInfo)
		if err != nil {
			log.Error("s.scoreStatistics error(%v)", err)
		}
		return
	})
	if err = readGroup.Wait(); err != nil {
		log.Error("run readGroup.Wait error(%v)", err)
		return
	}

	// persistent
	var writeGroup errgroup.Group

	//up_rating_statistics
	writeGroup.Go(func() (err error) {
		err = s.insertSections(c, sections)
		if err != nil {
			log.Error("s.insertSections error(%v)", err)
		}
		return
	})
	// up_rating_top
	writeGroup.Go(func() (err error) {
		_, err = s.insertTopRating(c, date, topRating, baseInfo)
		if err != nil {
			log.Error("s.insertSections error(%v)", err)
		}
		return
	})
	if err = writeGroup.Wait(); err != nil {
		log.Error("run writeGroup.Wait error(%v)", err)
	}
	return
}

func (s *Service) scoreStatistics(c context.Context, date time.Time, source chan []*model.Rating, baseInfo map[int64]*model.BaseInfo) (sections map[int]map[int64][]*model.RatingStatis, err error) {
	sections = make(map[int]map[int64][]*model.RatingStatis) // map[ctype][tagID][]*model.RatingStatis
	sections[TotalType] = make(map[int64][]*model.RatingStatis)
	sections[CreativeType] = make(map[int64][]*model.RatingStatis)
	sections[InfluenceType] = make(map[int64][]*model.RatingStatis)
	sections[CreditType] = make(map[int64][]*model.RatingStatis)
	for rating := range source {
		for _, r := range rating {
			statisScoreCtype(TotalType, r.CreativityScore+r.InfluenceScore+r.CreditScore, _sumTotal, sections, date, r, baseInfo[r.MID])
			statisScoreCtype(CreativeType, r.CreativityScore, _creativeTotal, sections, date, r, baseInfo[r.MID])
			statisScoreCtype(InfluenceType, r.InfluenceScore, _influenceTotal, sections, date, r, baseInfo[r.MID])
			statisScoreCtype(CreditType, r.CreditScore, _creditTotal, sections, date, r, baseInfo[r.MID])
		}
	}
	return
}

func statisScoreCtype(ctype int, score, totalScore int64, sections map[int]map[int64][]*model.RatingStatis, date time.Time, rate *model.Rating, base *model.BaseInfo) {
	if _, ok := sections[ctype][rate.TagID]; !ok {
		sections[ctype][rate.TagID] = initSections(totalScore, rate.TagID, ctype, xtime.Time(date.Unix()))
	}
	idx := score / _offset
	if idx >= int64(len(sections[ctype][rate.TagID])) {
		idx = int64(len(sections[ctype][rate.TagID]) - 1)
	}
	sections[ctype][rate.TagID][idx].Ups++
	sections[ctype][rate.TagID][idx].TotalScore += rate.CreativityScore + rate.InfluenceScore + rate.CreditScore
	sections[ctype][rate.TagID][idx].CreativityScore += rate.CreativityScore
	sections[ctype][rate.TagID][idx].InfluenceScore += rate.InfluenceScore
	sections[ctype][rate.TagID][idx].CreditScore += rate.CreditScore
	if base != nil {
		sections[ctype][rate.TagID][idx].Fans += base.TotalFans
		sections[ctype][rate.TagID][idx].Avs += base.TotalAvs
		sections[ctype][rate.TagID][idx].Coin += base.TotalCoin
		sections[ctype][rate.TagID][idx].Play += base.TotalPlay
	}
}

func (s *Service) insertSections(c context.Context, sections map[int]map[int64][]*model.RatingStatis) (err error) {
	for ctype, tags := range sections {
		for tagID, statis := range tags {
			_, err = s.insertRatingStatis(c, ctype, tagID, statis)
			if err != nil {
				return
			}
		}
	}
	return
}

func (s *Service) insertRatingStatis(c context.Context, ctype int, tagID int64, statis []*model.RatingStatis) (rows int64, err error) {
	return s.dao.InsertRatingStatis(c, assembleRatingStatis(c, ctype, tagID, statis))
}

func assembleRatingStatis(c context.Context, ctype int, tagID int64, statis []*model.RatingStatis) (vals string) {
	var buf bytes.Buffer
	for _, s := range statis {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(s.Ups, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.Section, 10))
		buf.WriteByte(',')
		buf.WriteString(s.Tips)
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.TotalScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.CreativityScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.InfluenceScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.CreditScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.Fans, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.Avs, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.Coin, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.Play, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(s.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(s.CType))
		buf.WriteByte(',')
		buf.WriteString("'" + s.CDate.Time().Format(_layout) + "'")
		buf.WriteString(")")
		buf.WriteByte(',')
	}

	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals = buf.String()
	buf.Reset()
	return
}
