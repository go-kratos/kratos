package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/up-rating/model"

	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_scoreTotalSQL = "SELECT count(*) FROM up_rating_%02d WHERE cdate=? %s"
	_scoreListSQL  = "SELECT mid, tag_id, cdate, creativity_score, influence_score, credit_score, magnetic_score FROM up_rating_%02d WHERE cdate=? %s"
	_levelListSQL  = "SELECT mid, total_fans, total_avs FROM up_level_info_%02d WHERE cdate=? AND mid IN (%s)"
	_upScoresSQL   = "SELECT cdate, creativity_score, influence_score, credit_score, magnetic_score FROM up_rating_%02d WHERE mid = ?"
	_upScoreSQL    = "SELECT cdate, creativity_score, influence_score, credit_score, magnetic_score FROM up_rating_%02d WHERE mid = ? AND cdate=?"
	_taskStatusSQL = "SELECT status FROM task_status WHERE date=?"
)

// Total counts scores
func (d *Dao) Total(c context.Context, mon int, date, where string) (total int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_scoreTotalSQL, mon, where), date)
	if err = row.Scan(&total); err != nil {
		log.Error("d.Total row.Scan error(%v)", err)
	}
	return
}

// ScoreList return info
func (d *Dao) ScoreList(c context.Context, mon int, date, where string) (list []*model.RatingInfo, err error) {
	list = make([]*model.RatingInfo, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_scoreListSQL, mon, where), date)
	if err != nil {
		log.Error("d.ScoreList d.db.Query error(%v)", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		v := &model.RatingInfo{}
		err = rows.Scan(&v.Mid, &v.TagID, &v.ScoreDate, &v.CreativityScore, &v.InfluenceScore, &v.CreditScore, &v.MagneticScore)
		if err != nil {
			log.Error("d.ScoreList rows.Scan error(%v)", err)
			return
		}
		list = append(list, v)
	}
	err = rows.Err()
	return
}

// LevelList returns info
func (d *Dao) LevelList(c context.Context, mon int, date string, mids []int64) (list []*model.RatingInfo, err error) {
	list = make([]*model.RatingInfo, 0, len(mids))
	rows, err := d.db.Query(c, fmt.Sprintf(_levelListSQL, mon, xstr.JoinInts(mids)), date)
	if err != nil {
		log.Error("d.LevelList d.db.Query error(%v)", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		v := &model.RatingInfo{}
		if err = rows.Scan(&v.Mid, &v.TotalFans, &v.TotalAvs); err != nil {
			log.Error("d.LevelList rows.Scan error(%v)", err)
			return
		}
		list = append(list, v)
	}
	err = rows.Err()
	return
}

// UpScores ...
func (d *Dao) UpScores(c context.Context, mon int, mid int64) (list []*model.RatingInfo, err error) {
	list = make([]*model.RatingInfo, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_upScoresSQL, mon), mid)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := &model.RatingInfo{}
		err = rows.Scan(&v.ScoreDate, &v.CreativityScore, &v.InfluenceScore, &v.CreditScore, &v.MagneticScore)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		list = append(list, v)
	}
	err = rows.Err()
	return
}

// UpScore ...
func (d *Dao) UpScore(c context.Context, mon int, mid int64, date string) (res *model.RatingInfo, err error) {
	res = new(model.RatingInfo)
	row := d.db.QueryRow(c, fmt.Sprintf(_upScoreSQL, mon), mid, date)
	err = row.Scan(&res.ScoreDate, &res.CreativityScore, &res.InfluenceScore, &res.CreditScore, &res.MagneticScore)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("row.Scan error(%v)", err)
	}
	return
}

// TaskStatus returns status of task on date
func (d *Dao) TaskStatus(c context.Context, date string) (status int, err error) {
	row := d.db.QueryRow(c, _taskStatusSQL, date)
	if err = row.Scan(&status); err != nil {
		log.Error("d.TaskStatus row.Scan error(%v)", err)
	}
	return
}
