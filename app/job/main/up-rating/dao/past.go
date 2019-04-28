package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/up-rating/model"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_pastRatingRecordSQL = "SELECT times FROM past_rating_record WHERE cdate = '%s' AND is_deleted = 0"
	_pastStatSQL         = "SELECT id,mid,creativity_score,influence_score,credit_score FROM past_score_statistics WHERE id > ? ORDER BY id LIMIT ?"

	// insert
	_inPastRecordSQL  = "INSERT INTO past_rating_record(times, cdate) VALUES(?,'%s') ON DUPLICATE KEY UPDATE times=VALUES(times)"
	_pastScoreStatSQL = "INSERT INTO past_score_statistics(mid,creativity_score,influence_score,credit_score) VALUES %s ON DUPLICATE KEY UPDATE creativity_score=creativity_score+VALUES(creativity_score),influence_score=influence_score+VALUES(influence_score),credit_score=VALUES(credit_score)"

	// delete
	_delPastStatSQL   = "DELETE FROM past_score_statistics ORDER BY ID LIMIT ?"
	_delPastRecordSQL = "DELETE FROM past_rating_record WHERE cdate='%s'"
)

// DelPastRecord del past record
func (d *Dao) DelPastRecord(c context.Context, date time.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_delPastRecordSQL, date.Format(_layout)))
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// GetPastRecord batch insert past score stat
func (d *Dao) GetPastRecord(c context.Context, cdate string) (times int, err error) {
	err = d.db.QueryRow(c, fmt.Sprintf(_pastRatingRecordSQL, cdate)).Scan(&times)
	if err == sql.ErrNoRows {
		err = nil
		times = -1
	}
	return
}

// InsertPastRecord insert past record date and times
func (d *Dao) InsertPastRecord(c context.Context, times int, cdate string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inPastRecordSQL, cdate), times)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// InsertPastScoreStat batch insert past score stat
func (d *Dao) InsertPastScoreStat(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_pastScoreStatSQL, values))
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", fmt.Sprintf(_pastScoreStatSQL, values), err)
		return
	}
	return res.RowsAffected()
}

// GetPasts get past statistics
func (d *Dao) GetPasts(c context.Context, offset, limit int64) (past []*model.Past, last int64, err error) {
	past = make([]*model.Past, 0, limit)
	rows, err := d.db.Query(c, _pastStatSQL, offset, limit)
	if err != nil {
		log.Error("d.db.Query GetPasts error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		p := &model.Past{}
		err = rows.Scan(&last, &p.MID, &p.MetaCreativityScore, &p.MetaInfluenceScore, &p.CreditScore)
		if err != nil {
			log.Error("rows.Scan GetPasts error(%v)", err)
			return
		}
		past = append(past, p)
	}
	return
}

// DelPastStat del past stat
func (d *Dao) DelPastStat(c context.Context, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delPastStatSQL, limit)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
