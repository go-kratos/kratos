package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/library/log"
)

const (
	_inRaingStatisticsSQL = "INSERT INTO up_rating_statistics(ups,section,tips,total_score,creativity_score,influence_score,credit_score,fans,avs,coin,play,tag_id,ctype,cdate) VALUES %s"
	_inAscSQL             = "INSERT INTO up_rating_trend_%s(mid,tag_id,creativity_score,creativity_diff,influence_score,influence_diff,credit_score,credit_diff,magnetic_score,magnetic_diff,date,ctype,section,tips) VALUES %s ON DUPLICATE KEY UPDATE tag_id=VALUES(tag_id),creativity_score=VALUES(creativity_score),influence_score=VALUES(influence_score),credit_score=VALUES(credit_score)"
	_inRatingTopSQL       = "INSERT INTO up_rating_top(mid,ctype,tag_id,score,fans,play,cdate) VALUES %s"

	_delTrendSQL     = "DELETE FROM up_rating_trend_%s LIMIT ?"
	_delRatingComSQL = "DELETE FROM %s WHERE cdate='%s' LIMIT ?"
)

// DelTrend del trend limit x
func (d *Dao) DelTrend(c context.Context, table string, limit int) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delTrendSQL, table), limit)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", fmt.Sprintf(_delTrendSQL, table), err)
		return
	}
	return res.RowsAffected()
}

// InsertRatingStatis batch insert rating statistics
func (d *Dao) InsertRatingStatis(c context.Context, values string) (rows int64, err error) {
	if values == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inRaingStatisticsSQL, values))
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", fmt.Sprintf(_inRaingStatisticsSQL, values), err)
		return
	}
	return res.RowsAffected()
}

// InsertTrend insert asc values
func (d *Dao) InsertTrend(c context.Context, table string, values string) (rows int64, err error) {
	if values == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inAscSQL, table, values))
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", fmt.Sprintf(_inAscSQL, table, values), err)
		return
	}
	return res.RowsAffected()
}

// InsertTopRating insert rating top
func (d *Dao) InsertTopRating(c context.Context, values string) (rows int64, err error) {
	if values == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inRatingTopSQL, values))
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", fmt.Sprintf(_inRatingTopSQL, values), err)
		return
	}
	return res.RowsAffected()
}

// DelRatingCom del rating common by date
func (d *Dao) DelRatingCom(c context.Context, table string, date time.Time, limit int) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delRatingComSQL, table, date.Format(_layout)), limit)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
