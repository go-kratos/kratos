package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/up-rating/model"
)

// DelRatings del ratings
func (s *Service) DelRatings(c context.Context, date time.Time) (err error) {
	for {
		var rows int64
		rows, err = s.dao.DelRatings(c, date, _limit)
		if err != nil {
			return
		}
		if rows == 0 {
			break
		}
	}
	return s.DelPastRecord(c, date)
}

// RatingFast close chan when part finished
func (s *Service) RatingFast(c context.Context, date time.Time, start, end int, ch chan []*model.Rating) (err error) {
	defer close(ch)
	return s.Ratings(c, date, start, end, ch)
}

// Ratings chan <- ratings, close chan outside
func (s *Service) Ratings(c context.Context, date time.Time, start, end int, ch chan []*model.Rating) (err error) {
	for {
		var rs []*model.Rating
		rs, start, err = s.dao.GetRatingsFast(c, date, start, end, _limit)
		if err != nil {
			return
		}
		if len(rs) == 0 {
			break
		}
		ch <- rs
	}
	return
}

// RatingInfos rating infos
func (s *Service) RatingInfos(c context.Context, date time.Time, ch chan []*model.Rating) (err error) {
	defer close(ch)
	var id int
	for {
		var rs []*model.Rating
		rs, id, err = s.dao.GetRatings(c, date, id, _limit)
		if err != nil {
			return
		}
		ch <- rs
		if len(rs) == 0 {
			break
		}
	}
	return
}

// BatchInsertRatingStat batch insert rating stat
func (s *Service) BatchInsertRatingStat(c context.Context, wch chan []*model.Rating, date time.Time) (err error) {
	var (
		buff    = make([]*model.Rating, _limit)
		buffEnd = 0
	)
	dateStr := date.Format(_layout)
	for rs := range wch {
		for _, r := range rs {
			buff[buffEnd] = r
			buffEnd++
			if buffEnd >= _limit {
				values := ratingStatValues(buff[:buffEnd], dateStr)
				buffEnd = 0
				_, err = s.dao.InsertRatingStat(c, date.Month(), values)
				if err != nil {
					return
				}
			}
		}
		if buffEnd > 0 {
			values := ratingStatValues(buff[:buffEnd], dateStr)
			buffEnd = 0
			_, err = s.dao.InsertRatingStat(c, date.Month(), values)
		}
	}
	return
}

func ratingStatValues(rs []*model.Rating, date string) (values string) {
	var buf bytes.Buffer
	for _, r := range rs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(r.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(fmt.Sprintf("'%s'", date))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.CreativityScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.InfluenceScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.CreditScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.MetaCreativityScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.MetaInfluenceScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.MagneticScore, 10))
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
