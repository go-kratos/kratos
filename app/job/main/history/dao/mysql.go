package dao

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/service/main/history/model"
	"go-common/library/database/tidb"
	"go-common/library/log"
)

var (
	_businessesSQL = "SELECT id, name, ttl FROM business"
	_addHistorySQL = "INSERT INTO histories(mid, kid, business_id, aid, sid, epid, sub_type, cid, device, progress, view_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)" +
		"ON DUPLICATE KEY UPDATE aid =?, sid=?, epid=?, sub_type=?, cid=?, device=?, progress=?, view_at=?"
	_deleteSQL       = "DELETE FROM histories WHERE business_id = ? AND mtime >= ? AND mtime < ? LIMIT ?"
	_allHisSQL       = "SELECT mtime FROM histories WHERE mid = ? AND business_id = ? ORDER BY mtime desc"
	_earlyHistorySQL = "SELECT mtime FROM histories  WHERE business_id = ?  ORDER BY mtime  LIMIT 1"
	_delUserHisSQL   = "DELETE FROM histories WHERE mid = ? AND mtime < ? and business_id = ?"
)

// Businesses business
func (d *Dao) Businesses(c context.Context) (res []*model.Business, err error) {
	var rows *tidb.Rows
	if rows, err = d.businessesStmt.Query(c); err != nil {
		log.Error("db.businessesStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.Business{}
		if err = rows.Scan(&b.ID, &b.Name, &b.TTL); err != nil {
			log.Error("rows.Business.Scan error(%v)", err)
			return
		}
		res = append(res, b)
	}
	err = rows.Err()
	return
}

// DeleteHistories delete histories
func (d *Dao) DeleteHistories(c context.Context, bid int64, beginTime, endTime time.Time) (rows int64, err error) {
	var res sql.Result
	begin := time.Now()
	if res, err = d.longDB.Exec(c, _deleteSQL, bid, beginTime, endTime, d.conf.Job.DeleteLimit); err != nil {
		log.Error("DeleteHistories(%v %v %v) err: %v", bid, beginTime, endTime, err)
		return
	}
	rows, err = res.RowsAffected()
	log.Info("clean business histories bid: %v begin: %v end: %v rows: %v, time: %v", bid, beginTime, endTime, rows, time.Since(begin))
	return
}

// AddHistories add histories to db
func (d *Dao) AddHistories(c context.Context, hs []*model.History) (err error) {
	if len(hs) == 0 {
		return
	}
	var tx *tidb.Tx
	if tx, err = d.db.Begin(c); err != nil {
		log.Error("tx.BeginTran() error(%v)", err)
		return
	}
	for _, h := range hs {
		if _, err = tx.Stmts(d.insertStmt).Exec(c, h.Mid, h.Kid, h.BusinessID, h.Aid, h.Sid, h.Epid, h.SubType, h.Cid, h.Device, h.Progress, h.ViewAt,
			h.Aid, h.Sid, h.Epid, h.SubType, h.Cid, h.Device, h.Progress, h.ViewAt); err != nil {
			log.Error("addHistories exec err mid: %v err: %+v", h.Mid, err)
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("add histories commit(%+v) err: %v", hs, err)
		return
	}
	log.Infov(c, log.D{Key: "log", Value: "addHistories db"}, log.D{Key: "len", Value: len(hs)})
	return
}

// DeleteUserHistories .
func (d *Dao) DeleteUserHistories(c context.Context, mid, bid int64, t time.Time) (rows int64, err error) {
	var res sql.Result
	if res, err = d.delUserStmt.Exec(c, mid, t, bid); err != nil {
		log.Error("DeleteUserHistories(%v %v) err: %v", bid, t, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// UserHistories .
func (d *Dao) UserHistories(c context.Context, mid, businessID int64) (res []time.Time, err error) {
	var rows *tidb.Rows
	if rows, err = d.allHisStmt.Query(c, mid, businessID); err != nil {
		log.Error("db.UserHistories.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var t time.Time
		if err = rows.Scan(&t); err != nil {
			log.Error("rows.UserHistories.Scan error(%v)", err)
			return
		}
		res = append(res, t)
	}
	err = rows.Err()
	return
}

// EarlyHistory .
func (d *Dao) EarlyHistory(c context.Context, businessID int64) (res time.Time, err error) {
	if err = d.longDB.QueryRow(c, _earlyHistorySQL, businessID).Scan(&res); err != nil {
		if err == tidb.ErrNoRows {
			res = time.Now()
			err = nil
			return
		}
		log.Error("db.EarlyHistory.Query error(%v)", err)
	}
	return
}
