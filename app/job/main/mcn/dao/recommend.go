package dao

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/job/main/mcn/model"
	xsql "go-common/library/database/sql"
)

const (
	_inMcnUpRecommendPoolSQL = `INSERT INTO mcn_up_recommend_pool (up_mid, fans_count, fans_count_increase_month, archive_count, play_count_accumulate, play_count_average, active_tid, last_archive_time, generate_time) VALUES (?,?,?,?,?,?,?,?,?) 
ON DUPLICATE KEY UPDATE fans_count=?, fans_count_increase_month=?, archive_count=?, play_count_accumulate=?, play_count_average=?, active_tid=?, last_archive_time=?, generate_time=?`
	_delMcnUpRecommendPoolSQL    = "UPDATE mcn_up_recommend_pool SET state = 100  WHERE generate_time < ? AND state IN (1,2) AND source = 1"
	_delMcnUpRecommendSourceSQL  = "DELETE FROM mcn_up_recommend_source WHERE id = ?"
	_selMcnUpRecommendSourcesSQL = "SELECT id,up_mid,fans_count,fans_count_increase_month,archive_count,play_count_accumulate,play_count_average,active_tid,last_archive_time,ctime,mtime FROM mcn_up_recommend_source LIMIT ?"
)

// AddMcnUpRecommend .
func (d *Dao) AddMcnUpRecommend(c context.Context, arg *model.McnUpRecommendPool) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _inMcnUpRecommendPoolSQL, arg.UpMid, arg.FansCount, arg.FansCountIncreaseMonth, arg.ArchiveCount, arg.PlayCountAccumulate, arg.PlayCountAverage, arg.ActiveTid, arg.LastArchiveTime, arg.GenerateTime,
		arg.FansCount, arg.FansCountIncreaseMonth, arg.ArchiveCount, arg.PlayCountAccumulate, arg.PlayCountAverage, arg.ActiveTid, arg.LastArchiveTime, arg.GenerateTime); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// DelMcnUpRecommendPool .
func (d *Dao) DelMcnUpRecommendPool(c context.Context) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _delMcnUpRecommendPoolSQL, time.Now().AddDate(0, 0, -3)); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// DelMcnUpRecommendSource .
func (d *Dao) DelMcnUpRecommendSource(c context.Context, id int64) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _delMcnUpRecommendSourceSQL, id); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// McnUpRecommendSources .
func (d *Dao) McnUpRecommendSources(c context.Context, limit int) (rps []*model.McnUpRecommendPool, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selMcnUpRecommendSourcesSQL, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		rp := new(model.McnUpRecommendPool)
		if err = rows.Scan(&rp.ID, &rp.UpMid, &rp.FansCount, &rp.FansCountIncreaseMonth, &rp.ArchiveCount, &rp.PlayCountAccumulate, &rp.PlayCountAverage, &rp.ActiveTid, &rp.LastArchiveTime, &rp.Ctime, &rp.Mtime); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		rps = append(rps, rp)
	}
	err = rows.Err()
	return
}
