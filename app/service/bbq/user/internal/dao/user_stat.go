package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"go-common/app/service/bbq/user/api"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_batchUserStatistics = "select `mid`,`av_total` - `unshelf_av_total`,`follow_total`,`fan_total`,`like_total`,`rev_like_total`,`black_total`,`play_total` from user_statistics where mid in (%s)"

	_incrUserStatisticsFollow = "update user_statistics set `follow_total` = `follow_total` + 1 where `mid` = ?"
	_incrUserStatisticsFan    = "update user_statistics set `fan_total` = `fan_total` + 1 where `mid` = ?"
	_decrUserStatisticsFollow = "update user_statistics set `follow_total` = `follow_total` - 1 where `mid` = ? and `follow_total` > 0"
	_decrUserStatisticsFan    = "update user_statistics set `fan_total` = `fan_total` - 1 where `mid` = ? and `fan_total` > 0"

	_incrUserStatisticField  = "insert into user_statistics (`mid`, `%s`) values (?, 1) on duplicate key update `%s` = `%s` + 1"
	_decrUserStatisticsField = "update user_statistics set `%s` = `%s` - 1 where `mid` = ? and `%s` > 0"

	_updateUserVideoView = "update `user_statistics` set `play_total` = ? where `mid` = ?;"
)

// RawBatchUserStatistics 从数据库获取用户基础信息
func (d *Dao) RawBatchUserStatistics(c context.Context, mids []int64) (res map[int64]*api.UserStat, err error) {
	if len(mids) == 0 {
		return
	}
	res = make(map[int64]*api.UserStat)

	midStr := xstr.JoinInts(mids)
	querySQL := fmt.Sprintf(_batchUserStatistics, midStr)
	rows, err := d.db.Query(c, querySQL)
	if err != nil {
		log.Errorv(c, log.KV("event", "mysql_query"), log.KV("error", err), log.KV("sql", querySQL))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		stat := new(api.UserStat)
		if err = rows.Scan(&mid, &stat.Sv, &stat.Follow, &stat.Fan, &stat.Like, &stat.Liked, &stat.Black, &stat.View); err != nil {
			log.Errorv(c, log.KV("event", "mysql_scan"), log.KV("error", err), log.KV("sql", querySQL))
			return
		}
		res[mid] = stat
	}
	log.Infov(c, log.KV("event", "mysql_query"), log.KV("row_num", len(res)), log.KV("sql", querySQL))
	return
}

//TxIncrUserStatisticsFollow .
func (d *Dao) TxIncrUserStatisticsFollow(tx *sql.Tx, mid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_incrUserStatisticsFollow, mid); err != nil {
		log.Error("incr user Video follow err(%v)", err)
		return
	}
	return res.RowsAffected()
}

//TxIncrUserStatisticsFan .
func (d *Dao) TxIncrUserStatisticsFan(tx *sql.Tx, mid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_incrUserStatisticsFan, mid); err != nil {
		log.Error("incr user Video fan err(%v)", err)
		return
	}
	return res.RowsAffected()
}

//TxDecrUserStatisticsFollow .
func (d *Dao) TxDecrUserStatisticsFollow(tx *sql.Tx, mid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_decrUserStatisticsFollow, mid); err != nil {
		log.Error("decr user Video follow err(%v)", err)
		return
	}
	return res.RowsAffected()
}

//TxDecrUserStatisticsFan .
func (d *Dao) TxDecrUserStatisticsFan(tx *sql.Tx, mid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(_decrUserStatisticsFan, mid); err != nil {
		log.Error("decr user Video fan err(%v)", err)
		return
	}
	return res.RowsAffected()
}

//TxIncrUserStatisticsField .
func (d *Dao) TxIncrUserStatisticsField(c context.Context, tx *sql.Tx, mid int64, field string) (rowsAffected int64, err error) {
	var res xsql.Result
	querySQL := fmt.Sprintf(_incrUserStatisticField, field, field, field)
	if res, err = tx.Exec(querySQL, mid); err != nil {
		log.Errorv(c, log.KV("event", "incr_user_statistic"), log.KV("field", field), log.KV("mid", mid), log.KV("err", err))
		return
	}
	return res.RowsAffected()
}

//TxDescUserStatisticsField .
func (d *Dao) TxDescUserStatisticsField(c context.Context, tx *sql.Tx, mid int64, field string) (rowsAffected int64, err error) {
	var res xsql.Result
	querySQL := fmt.Sprintf(_decrUserStatisticsField, field, field, field)
	if res, err = tx.Exec(querySQL, mid); err != nil {
		log.Errorv(c, log.KV("event", "incr_user_statistic"), log.KV("field", field), log.KV("mid", mid), log.KV("err", err))
		return
	}
	return res.RowsAffected()
}

// UpdateUserVideoView .
func (d *Dao) UpdateUserVideoView(c context.Context, mid int64, views int64) error {
	_, err := d.db.Exec(c, _updateUserVideoView, views, mid)
	return err
}
