package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/up-rating/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_upScoreSQL    = "SELECT mid, creativity_score, influence_score, credit_score, cdate FROM up_rating_%02d WHERE mid=? AND cdate=? AND is_deleted=0"
	_taskStatusSQL = "SELECT status FROM task_status WHERE date=? "

	_whitelistSQL = "SELECT count(*) FROM up_white_list WHERE mid=? AND is_deleted=0"
)

// White will del later
func (d *Dao) White(c context.Context, mid int64) (count int64, err error) {
	err = d.rddb.QueryRow(c, _whitelistSQL, mid).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("UpScore row scan error(%v)", err)
		}
	}
	return
}

// TaskStatus ...
func (d *Dao) TaskStatus(c context.Context, date string) (status int, err error) {
	err = d.rddb.QueryRow(c, _taskStatusSQL, date).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("UpScore row scan error(%v)", err)
		}
	}
	return
}

// UpScore gets score data of UP
func (d *Dao) UpScore(c context.Context, mon int, mid int64, date string) (score *model.Score, err error) {
	score = new(model.Score)
	row := d.rddb.QueryRow(c, fmt.Sprintf(_upScoreSQL, mon), mid, date)
	err = row.Scan(&score.MID, &score.Creative, &score.Influence, &score.Credit, &score.CDate)
	if err != nil {
		if err == sql.ErrNoRows {
			score, err = nil, nil
		} else {
			log.Error("UpScore row scan error(%v)", err)
		}
	}
	return
}
