package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"time"

	"go-common/app/job/main/spy/conf"
	"go-common/app/job/main/spy/model"
	"go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_reBuildMidCountSQL    = "SELECT count(mid) FROM spy_user_info_%02d WHERE state=? AND score>=30 AND mtime BETWEEN ? AND ?;"
	_reBuildMidSQL         = "SELECT mid FROM spy_user_info_%02d WHERE state=? AND score>=30 AND mtime BETWEEN ? AND ? LIMIT ?;"
	_getAllConfigSQL       = "SELECT id, property,name,val,ctime FROM spy_system_config;"
	_addEventHistorySQL    = "INSERT INTO spy_user_event_history_%02d (mid,event_id,score,base_score,event_score,remark,reason,factor_val,ctime) VALUES (?,?,?,?,?,?,?,?,?);"
	_addPunishmentSQL      = "INSERT INTO spy_punishment (mid,type,reason,batch_no,ctime) VALUES (?,?,?,?,?);"
	_updateUserStateSQL    = "UPDATE spy_user_info_%02d SET state=? WHERE mid=?"
	_getLastHistorySQL     = "SELECT id,mid,event_id,score,base_score,event_score,remark,reason,factor_val,ctime FROM spy_user_event_history_%02d WHERE mid=? ORDER BY id DESC LIMIT 1;"
	_getHistoryListSQL     = "SELECT remark,reason FROM spy_user_event_history_%02d WHERE mid= ? ORDER BY id DESC LIMIT ?;"
	_updateEventScoreSQL   = "UPDATE spy_user_info_%02d SET event_score=?, score=? WHERE mid=?;"
	_userInfoSQL           = "SELECT id,mid,score,base_score,event_score,state,ctime,mtime FROM spy_user_info_%02d WHERE mid=? LIMIT 1;"
	_punishmentCountSQL    = "SELECT COUNT(1) FROM spy_punishment where mtime > ? and mtime < ?;"
	_securityLoginCountSQL = "SELECT COUNT(1) FROM spy_user_event_history_%02d where reason = ? and mtime > ? and mtime < ?;"
	_insertReportSQL       = "INSERT INTO `spy_report`(`name`,`date_version`,`val`,`ctime`)VALUES(?,?,?,?);"
	_insertIncrStatSQL     = "INSERT INTO `spy_statistics`(`target_mid`,`target_id`,`event_id`,`state`,`type`,`quantity`,`ctime`)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE quantity=quantity+?;"
	_insertStatSQL         = "INSERT INTO `spy_statistics`(`target_mid`,`target_id`,`event_id`,`state`,`type`,`quantity`,`ctime`)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE quantity=?;"
	_allEventSQL           = "SELECT id,name,nick_name,service_id,status,ctime,mtime FROM spy_event WHERE status<>0"
)

func hitHistory(id int64) int64 {
	return id % conf.Conf.Property.HistoryShard
}

func hitInfo(id int64) int64 {
	return id % conf.Conf.Property.UserInfoShard
}

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}

// ReBuildMidCount count for need reBuild user.
func (d *Dao) ReBuildMidCount(c context.Context, i int, state int8, start, end time.Time) (res int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_reBuildMidCountSQL, i), state, start, end)
	if err = row.Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// ReBuildMidList query reBuild user mid list by page.
func (d *Dao) ReBuildMidList(c context.Context, i int, t int8, start, end time.Time, ps int64) (res []int64, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_reBuildMidSQL, i), t, start, end, ps); err != nil {
		log.Error("d.reBuildMidSQL.Query(%d, %s, %s, %d) error(%v)", t, start, end, ps, err)
		return
	}
	defer rows.Close()
	res = []int64{}
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// Configs spy system configs.
func (d *Dao) Configs(c context.Context) (res map[string]string, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getAllConfigSQL); err != nil {
		log.Error("d.getAllConfigSQL.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = map[string]string{}
	for rows.Next() {
		var r model.Config
		if err = rows.Scan(&r.ID, &r.Property, &r.Name, &r.Val, &r.Ctime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res[r.Property] = r.Val
	}
	err = rows.Err()
	return
}

// TxAddEventHistory insert user_event_history.
func (d *Dao) TxAddEventHistory(c context.Context, tx *sql.Tx, ueh *model.UserEventHistory) (err error) {
	var (
		now = time.Now()
	)
	if _, err = tx.Exec(fmt.Sprintf(_addEventHistorySQL, hitHistory(ueh.Mid)), ueh.Mid, ueh.EventID, ueh.Score, ueh.BaseScore, ueh.EventScore, ueh.Remark, ueh.Reason, ueh.FactorVal, now); err != nil {
		log.Error("db.Exec(%v) error(%v)", ueh, err)
		return
	}
	return
}

// TxAddPunishment insert punishment.
func (d *Dao) TxAddPunishment(c context.Context, tx *sql.Tx, mid int64, t int8, reason string, blockNo int64) (err error) {
	var (
		now = time.Now()
	)
	if _, err = tx.Exec(_addPunishmentSQL, mid, t, reason, blockNo, now); err != nil {
		log.Error("db.Exec(%d, %d, %s) error(%v)", mid, t, reason, err)
		return
	}
	return
}

// History get last one user history.
func (d *Dao) History(c context.Context, mid int64) (h *model.UserEventHistory, err error) {
	var (
		row *sql.Row
	)
	h = &model.UserEventHistory{}
	row = d.db.QueryRow(c, fmt.Sprintf(_getLastHistorySQL, hitHistory(mid)), mid)
	if err = row.Scan(&h.ID, &h.Mid, &h.EventID, &h.Score, &h.BaseScore, &h.EventScore, &h.Remark, &h.Reason, &h.FactorVal, &h.CTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			h = nil
			return
		}
		log.Error("History row.Scan(%d) error(%v)", mid, err)
	}
	return
}

// TxUpdateUserState insert or update  user state by mid.
func (d *Dao) TxUpdateUserState(c context.Context, tx *sql.Tx, info *model.UserInfo) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_updateUserStateSQL, hitInfo(info.Mid)), info.State, info.Mid); err != nil {
		log.Error("TxUpdateUserState db.Exec(%d, %v) error(%v)", info.Mid, info, err)
		return
	}
	return
}

// HistoryList query .
func (d *Dao) HistoryList(c context.Context, mid int64, size int) (res []*model.UserEventHistory, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_getHistoryListSQL, hitHistory(mid)), mid, size); err != nil {
		log.Error("d.HistoryList.Query(%d, %d) error(%v)", mid, size, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.UserEventHistory{}
		if err = rows.Scan(&r.Remark, &r.Reason); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// UserInfo get info by mid.
func (d *Dao) UserInfo(c context.Context, mid int64) (res *model.UserInfo, err error) {
	var (
		row *sql.Row
	)
	res = &model.UserInfo{}
	row = d.db.QueryRow(c, fmt.Sprintf(_userInfoSQL, hitInfo(mid)), mid)
	if err = row.Scan(&res.ID, &res.Mid, &res.Score, &res.BaseScore, &res.EventScore, &res.State, &res.CTime, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// TxUpdateEventScore update event score.
func (d *Dao) TxUpdateEventScore(c context.Context, tx *sql.Tx, mid int64, escore, score int8) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_updateEventScoreSQL, hitInfo(mid)), escore, score, mid); err != nil {
		log.Error("db.TxUpdateEventScore(%s, %d, %d, %d) error(%v)", _updateEventScoreSQL, escore, score, mid, err)
		return
	}
	return
}

// AddReport add report info.
func (d *Dao) AddReport(c context.Context, r *model.Report) (affected int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _insertReportSQL, r.Name, r.DateVersion, r.Val, r.Ctime); err != nil {
		log.Error("AddReport: db.Exec(%v) error(%v)", r, err)
		return
	}
	return res.RowsAffected()
}

// PunishmentCount punishment count.
func (d *Dao) PunishmentCount(c context.Context, start, end time.Time) (res int64, err error) {
	row := d.db.QueryRow(c, _punishmentCountSQL, start, end)
	if err = row.Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// SecurityLoginCount security login count.
func (d *Dao) SecurityLoginCount(c context.Context, index int64, reason string, stime, etime time.Time) (res int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_securityLoginCountSQL, hitHistory(index)), reason, stime, etime)
	if err = row.Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// AddStatistics add statistics info.
func (d *Dao) AddStatistics(c context.Context, s *model.Statistics) (id int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _insertStatSQL, s.TargetMid, s.TargetID, s.EventID, s.State, s.Type, s.Quantity, s.Ctime, s.Quantity); err != nil {
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// AddIncrStatistics add increase statistics info.
func (d *Dao) AddIncrStatistics(c context.Context, s *model.Statistics) (id int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, _insertIncrStatSQL, s.TargetMid, s.TargetID, s.EventID, s.State, s.Type, s.Quantity, s.Ctime, s.Quantity); err != nil {
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//AllEvent all event.
func (d *Dao) AllEvent(c context.Context) (list []*model.Event, err error) {
	var (
		rows *sql.Rows
	)
	list = make([]*model.Event, 0)
	if rows, err = d.db.Query(c, _allEventSQL); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _allEventSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var event = &model.Event{}
		if err = rows.Scan(&event.ID, &event.Name, &event.NickName, &event.ServiceID, &event.Status, &event.Ctime, &event.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, &model.Event{
			ID:        event.ID,
			Name:      event.Name,
			NickName:  event.NickName,
			ServiceID: event.ServiceID,
			Status:    event.Status,
			Ctime:     event.Ctime,
			Mtime:     event.Mtime,
		})
	}
	return
}
