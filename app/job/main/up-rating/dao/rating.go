package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/job/main/up-rating/model"

	"go-common/library/log"
)

const (
	_layout = "2006-01-02"
	// get up_rating start & end
	_ratingStartSQL = "SELECT id FROM up_rating_%02d WHERE cdate='%s' ORDER BY id LIMIT 1"
	_ratingEndSQL   = "SELECT id FROM up_rating_%02d WHERE cdate='%s' ORDER BY id DESC LIMIT 1"
	_ratingCountSQL = "SELECT COUNT(*) FROM up_rating_%02d WHERE cdate='%s'"

	_ratingSQL     = "SELECT id,mid,tag_id,creativity_score,influence_score,credit_score,meta_creativity_score,meta_influence_score,cdate FROM up_rating_%02d WHERE cdate= '%s' AND id > ? ORDER BY id LIMIT ?"
	_ratingByIDSQL = "SELECT id,mid,tag_id,creativity_score,influence_score,credit_score,meta_creativity_score,meta_influence_score,magnetic_score,cdate FROM up_rating_%02d WHERE id > ? AND id <= ? LIMIT ?"

	_ratingScoreSQL    = "INSERT INTO up_rating_%02d(mid,tag_id,cdate,creativity_score,influence_score,credit_score,meta_creativity_score,meta_influence_score, magnetic_score) VALUES %s ON DUPLICATE KEY UPDATE tag_id=VALUES(tag_id), cdate=VALUES(cdate), creativity_score=VALUES(creativity_score),influence_score=VALUES(influence_score),credit_score=VALUES(credit_score),meta_creativity_score=VALUES(meta_creativity_score),meta_influence_score=VALUES(meta_influence_score),magnetic_score=VALUES(magnetic_score)"
	_delRatingScoreSQL = "DELETE FROM up_rating_%02d WHERE cdate='%s' LIMIT ?"
)

// DelRatings del ratings
func (d *Dao) DelRatings(c context.Context, date time.Time, limit int) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delRatingScoreSQL, date.Month(), date.Format(_layout)), limit)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// RatingStart get start id by date
func (d *Dao) RatingStart(c context.Context, date time.Time) (start int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_ratingStartSQL, date.Month(), date.Format(_layout)))
	if err = row.Scan(&start); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// RatingEnd get end id by date
func (d *Dao) RatingEnd(c context.Context, date time.Time) (end int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_ratingEndSQL, date.Month(), date.Format(_layout)))
	if err = row.Scan(&end); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// RatingCount get end id by date
func (d *Dao) RatingCount(c context.Context, date time.Time) (count int, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_ratingCountSQL, date.Month(), date.Format(_layout)))
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// GetRatings get ratings by date
func (d *Dao) GetRatings(c context.Context, date time.Time, offset, limit int) (rs []*model.Rating, last int, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_ratingSQL, date.Month(), date.Format(_layout)), offset, limit)
	if err != nil {
		log.Error("d.db.Query Rating Info error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Rating{}
		err = rows.Scan(&last, &r.MID, &r.TagID, &r.CreativityScore, &r.InfluenceScore, &r.CreditScore, &r.MetaCreativityScore, &r.MetaInfluenceScore, &r.Date)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	return
}

// GetRatingsFast get rating fast
func (d *Dao) GetRatingsFast(c context.Context, date time.Time, start, end, limit int) (rs []*model.Rating, id int, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_ratingByIDSQL, date.Month()), start, end, limit)
	if err != nil {
		log.Error("d.db.Query Rating Info error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Rating{}
		err = rows.Scan(&id, &r.MID, &r.TagID, &r.CreativityScore, &r.InfluenceScore, &r.CreditScore, &r.MetaCreativityScore, &r.MetaInfluenceScore, &r.MagneticScore, &r.Date)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	return
}

// InsertRatingStat batch insert rating score stat
func (d *Dao) InsertRatingStat(c context.Context, month time.Month, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_ratingScoreSQL, month, values))
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", fmt.Sprintf(_ratingScoreSQL, month, values), err)
		return
	}
	return res.RowsAffected()
}
