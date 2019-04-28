package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/up-rating/model"

	"go-common/library/log"
)

const (
	_totalType int64 = iota
	_creativeType
	_influenceType
	_creditType
)

const (
	// select
	_ratingStatisticsSQL = "SELECT ups,section,tips,total_score,creativity_score,influence_score,credit_score,fans,avs,coin,play,tag_id,ctype,cdate FROM up_rating_statistics WHERE cdate = '%s' AND ctype = ? AND tag_id IN (%s)"

	_trendAscCountSQL  = "SELECT count(*) FROM up_rating_trend_asc WHERE date='%s' %s"
	_trendDescCountSQL = "SELECT count(*) FROM up_rating_trend_desc WHERE date='%s' %s"
	_trendAscSQL       = "SELECT mid,magnetic_score,creativity_score,influence_score,credit_score,%s_diff FROM up_rating_trend_asc WHERE date='%s' %s"
	_trendDescSQL      = "SELECT mid,magnetic_score,creativity_score,influence_score,credit_score,%s_diff FROM up_rating_trend_desc WHERE date='%s' %s"
)

// GetRatingStatis get rating statistics
func (d *Dao) GetRatingStatis(c context.Context, ctype int64, date, tags string) (statis []*model.RatingStatis, err error) {
	statis = make([]*model.RatingStatis, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_ratingStatisticsSQL, date, tags), ctype)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		s := &model.RatingStatis{}
		err = rows.Scan(
			&s.Ups, &s.Section, &s.Tips, &s.TotalScore, &s.CreativityScore, &s.InfluenceScore, &s.CreditScore, &s.Fans, &s.Avs, &s.Coin, &s.Play, &s.TagID, &s.CType, &s.CDate)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		switch ctype {
		case _totalType:
			s.Score = s.TotalScore
		case _creativeType:
			s.Score = s.CreativityScore
		case _influenceType:
			s.Score = s.InfluenceScore
		case _creditType:
			s.Score = s.CreditScore
		}
		statis = append(statis, s)
	}
	err = rows.Err()
	return
}

// AscTrendCount asc trend count
func (d *Dao) AscTrendCount(c context.Context, date string, query string) (count int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_trendAscCountSQL, date, query))
	if err = row.Scan(&count); err != nil {
		log.Error("d.db.Query error(%v)", err)
	}
	return
}

// DescTrendCount desc trend count
func (d *Dao) DescTrendCount(c context.Context, date string, query string) (count int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_trendDescCountSQL, date, query))
	if err = row.Scan(&count); err != nil {
		log.Error("d.db.Query error(%v)", err)
	}
	return
}

// GetTrendAsc get asc trend
func (d *Dao) GetTrendAsc(c context.Context, ctype string, date string, query string) (ts []*model.Trend, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_trendAscSQL, ctype, date, query))
	if err != nil {
		return
	}
	ts = make([]*model.Trend, 0)
	defer rows.Close()
	for rows.Next() {
		t := &model.Trend{}
		err = rows.Scan(&t.MID, &t.MagneticScore, &t.CreativityScore, &t.InfluenceScore, &t.CreditScore, &t.DValue)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		ts = append(ts, t)
	}
	return
}

// GetTrendDesc get desc trend
func (d *Dao) GetTrendDesc(c context.Context, ctype string, date string, query string) (ts []*model.Trend, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_trendDescSQL, ctype, date, query))
	if err != nil {
		return
	}
	ts = make([]*model.Trend, 0)
	defer rows.Close()
	for rows.Next() {
		t := &model.Trend{}
		err = rows.Scan(&t.MID, &t.MagneticScore, &t.CreativityScore, &t.InfluenceScore, &t.CreditScore, &t.DValue)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		ts = append(ts, t)
	}
	return
}
