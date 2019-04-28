package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"time"

	"go-common/app/service/main/spy/conf"
	"go-common/app/service/main/spy/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_serviceSQL                      = "SELECT id,name,nick_name,status,ctime,mtime FROM spy_service WHERE name=? AND status<>0 LIMIT 1"
	_addServiceSQL                   = "INSERT INTO spy_service (name,nick_name,status,ctime) VALUES (?,?,?,?)"
	_eventSQL                        = "SELECT id,name,nick_name,service_id,status,ctime,mtime FROM spy_event WHERE name=? AND status<>0 LIMIT 1"
	_addEventSQL                     = "INSERT INTO spy_event (name,nick_name,service_id,status,ctime) VALUES (?,?,?,?,?)"
	_factorSQL                       = "SELECT id,nick_name,service_id,event_id,group_id,risk_level,factor_val,ctime,mtime FROM spy_factor WHERE service_id=? AND event_id=? AND risk_level=? LIMIT 1"
	_factorGroupSQL                  = "SELECT id,name,ctime FROM spy_factor_group WHERE name=? LIMIT 1"
	_userInfoSQL                     = "SELECT id,mid,score,base_score,event_score,state,relive_times,ctime,mtime FROM spy_user_info_%02d WHERE mid=? LIMIT 1"
	_addUserInfoSQL                  = "INSERT IGNORE INTO spy_user_info_%02d (mid,score,base_score,event_score,state,ctime,mtime) VALUES (?,?,?,?,?,?,?)"
	_updateEventScoreSQL             = "UPDATE spy_user_info_%02d SET event_score=?, score=? WHERE mid=?"
	_updateInfoSQL                   = "UPDATE spy_user_info_%02d SET base_score=?, event_score=?, score=?, state=? WHERE mid=?;"
	_updateBaseScoreSQL              = "UPDATE spy_user_info_%02d SET base_score=?, score=? WHERE mid=?;"
	_addEventHistorySQL              = "INSERT INTO spy_user_event_history_%02d (mid,event_id,score,base_score,event_score,remark,reason,factor_val,ctime) VALUES (?,?,?,?,?,?,?,?,?);"
	_getEventHistoryByMidAndEventSQL = "SELECT mid, event_id FROM spy_user_event_history_%02d WHERE mid = ? and event_id = ? limit 1"
	_addPunishmentSQL                = "INSERT INTO spy_punishment (mid,type,reason,ctime) VALUES (?,?,?,?);"
	_addPunishmentQueueSQL           = "INSERT INTO spy_punishment_queue (mid,batch_no,ctime) VALUES (?,?,?) ON DUPLICATE KEY UPDATE mtime=?;"
	_getAllConfigSQL                 = "SELECT id, property,name,val,ctime FROM spy_system_config;"
	_clearReliveTimesSQL             = "UPDATE spy_user_info_%02d SET relive_times=? WHERE mid=?;"
	_getHistoryListSQL               = "SELECT remark,reason FROM spy_user_event_history_%02d WHERE mid= ? ORDER BY id DESC LIMIT ?;"
	_updateEventScoreReLiveSQL       = "UPDATE spy_user_info_%02d SET event_score=?, score=?, relive_times=relive_times+1 WHERE mid=?"
	_statListByMidSQL                = "SELECT event_id,quantity FROM spy_statistics WHERE target_mid = ? AND isdel = 0  ORDER BY id desc LIMIT 100;"
	_statListByIDSQL                 = "SELECT event_id,quantity FROM spy_statistics WHERE target_id = ? AND isdel = 0 ORDER BY id desc LIMIT 100;"
	_statListByIDAndMidSQL           = "SELECT event_id,quantity FROM spy_statistics WHERE target_mid = ? AND target_id = ? AND isdel = 0 ORDER BY id desc LIMIT 100;"
	_allEventSQL                     = "SELECT id,name,nick_name,service_id,status,ctime,mtime FROM spy_event WHERE status<>0"
	_telLevelSQL                     = "SELECT id,mid,level,origin,ctime,mtime FROM spy_tel_risk_level WHERE mid = ? LIMIT 1;"
	_addTelLevelSQL                  = "INSERT INTO spy_tel_risk_level (mid,level,origin,ctime) VALUES (?,?,?,?) "
)

// Service get service from db.
func (d *Dao) Service(c context.Context, serviceName string) (service *model.Service, err error) {
	var (
		row *sql.Row
	)
	service = &model.Service{}
	row = d.db.QueryRow(c, _serviceSQL, serviceName)
	if err = row.Scan(&service.ID, &service.Name, &service.NickName, &service.Status, &service.CTime, &service.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			service = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// Event get event from db.
func (d *Dao) Event(c context.Context, eventName string) (event *model.Event, err error) {
	var row *sql.Row
	event = &model.Event{}
	row = d.db.QueryRow(c, _eventSQL, eventName)
	if err = row.Scan(&event.ID, &event.Name, &event.NickName, &event.ServiceID, &event.Status, &event.CTime, &event.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			event = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// AddService insert service to db.
func (d *Dao) AddService(c context.Context, service *model.Service) (id int64, err error) {
	var (
		res xsql.Result
	)
	if res, err = d.db.Exec(c, _addServiceSQL, service.Name, service.NickName, service.Status, time.Now()); err != nil {
		log.Error("d.db.Exec(%s, %s, %d) error(%v)", _addServiceSQL, service.Name, service.Status, err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		log.Error("res.LastInsertId() error(%v)", err)
	}
	return
}

// AddEvent insert service to db.
func (d *Dao) AddEvent(c context.Context, event *model.Event) (id int64, err error) {
	var (
		res xsql.Result
	)
	if res, err = d.db.Exec(c, _addEventSQL, event.Name, event.NickName, event.ServiceID, event.Status, time.Now()); err != nil {
		log.Error("d.db.Exec(%s, %s, %s, %d, %d) error(%v)", _addEventSQL, event.Name, event.NickName, event.ServiceID, event.Status, err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		log.Error("res.LastInsertId() error(%v)", err)
	}
	return
}

// Factor get factor from db.
func (d *Dao) Factor(c context.Context, serviceID, eventID int64, riskLevel int8) (factor *model.Factor, err error) {
	var (
		row *sql.Row
	)
	factor = &model.Factor{}
	row = d.db.QueryRow(c, _factorSQL, serviceID, eventID, riskLevel)
	if err = row.Scan(&factor.ID, &factor.NickName, &factor.ServiceID, &factor.EventID, &factor.GroupID, &factor.RiskLevel, &factor.FactorVal, &factor.CTime, &factor.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			factor = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// FactorGroup get factor group from db.
func (d *Dao) FactorGroup(c context.Context, groupName string) (factorGroup *model.FactorGroup, err error) {
	var (
		row *sql.Row
	)
	factorGroup = &model.FactorGroup{}
	row = d.db.QueryRow(c, _factorGroupSQL, groupName)
	if err = row.Scan(&factorGroup.ID, &factorGroup.Name, &factorGroup.CTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			factorGroup = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

func hitInfo(id int64) int64 {
	return id % conf.Conf.Property.UserInfoShard
}

func hitHistory(id int64) int64 {
	return id % conf.Conf.Property.HistoryShard
}

// BeginTran begin transaction.
func (d *Dao) BeginTran(c context.Context) (*sql.Tx, error) {
	return d.db.Begin(c)
}

// UserInfo get info by mid.
func (d *Dao) UserInfo(c context.Context, mid int64) (res *model.UserInfo, err error) {
	var (
		row *sql.Row
	)
	res = &model.UserInfo{}
	row = d.db.QueryRow(c, fmt.Sprintf(_userInfoSQL, hitInfo(mid)), mid)
	if err = row.Scan(&res.ID, &res.Mid, &res.Score, &res.BaseScore, &res.EventScore, &res.State, &res.ReliveTimes, &res.CTime, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// TxUpdateInfo insert or update  user info by mid.
func (d *Dao) TxUpdateInfo(c context.Context, tx *sql.Tx, info *model.UserInfo) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_updateInfoSQL, hitInfo(info.Mid)), info.BaseScore, info.EventScore, info.Score, info.State, info.Mid); err != nil {
		log.Error("db.Exec(%d, %+v) error(%v)", info.Mid, info, err)
		return
	}
	return
}

// TxAddInfo add user info.
func (d *Dao) TxAddInfo(c context.Context, tx *sql.Tx, info *model.UserInfo) (id int64, err error) {
	var (
		res xsql.Result
		now = time.Now()
	)
	if res, err = tx.Exec(fmt.Sprintf(_addUserInfoSQL, hitInfo(info.Mid)), info.Mid, info.Score, info.BaseScore, info.EventScore, info.State, now, now); err != nil {
		log.Error("db.Exec(%d, %v) error(%v)", info.Mid, *info, err)
		return
	}
	return res.LastInsertId()
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

// EventHistoryByMidAndEvent get one event history with mid and eventID from db.
func (d *Dao) EventHistoryByMidAndEvent(c context.Context, mid int64, eventID int64) (res *model.UserEventHistory, err error) {
	var (
		row *sql.Row
	)
	res = &model.UserEventHistory{}
	row = d.db.QueryRow(c, fmt.Sprintf(_getEventHistoryByMidAndEventSQL, hitHistory(mid)), mid, eventID)
	if err = row.Scan(&res.Mid, &res.EventID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("row.Scan() error:(%v)", err)
	}
	return
}

// TxAddPunishment insert punishment.
func (d *Dao) TxAddPunishment(c context.Context, tx *sql.Tx, mid int64, t int8, reason string) (err error) {
	var (
		now = time.Now()
	)
	if _, err = tx.Exec(_addPunishmentSQL, mid, t, reason, now); err != nil {
		log.Error("db.Exec(%d, %d, %s) error(%v)", mid, t, reason, err)
		return
	}
	return
}

// TxAddPunishmentQueue insert punishment queue.
func (d *Dao) TxAddPunishmentQueue(c context.Context, tx *sql.Tx, mid int64, blockNo int64) (err error) {
	var (
		now = time.Now()
	)
	if _, err = tx.Exec(_addPunishmentQueueSQL, mid, blockNo, now, now); err != nil {
		log.Error("TxAddPunishmentQueue:db.Exec(%d) error(%v)", mid, err)
		return
	}
	return
}

// AddPunishmentQueue insert punishment queue.
func (d *Dao) AddPunishmentQueue(c context.Context, mid int64, blockNo int64) (err error) {
	var (
		now = time.Now()
	)
	if _, err = d.db.Exec(c, _addPunishmentQueueSQL, mid, blockNo, now, now); err != nil {
		log.Error("AddPunishmentQueue:db.Exec(%d) error(%v)", mid, err)
		return
	}
	return
}

// TxUpdateEventScore update user event score.
func (d *Dao) TxUpdateEventScore(c context.Context, tx *sql.Tx, mid int64, escore, score int8) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_updateEventScoreSQL, hitInfo(mid)), escore, score, mid); err != nil {
		log.Error("db.TxUpdateEventScore(%s, %d, %d, %d) error(%v)", _updateEventScoreSQL, escore, score, mid, err)
		return
	}
	return
}

//TxUpdateBaseScore do update user base score.
func (d *Dao) TxUpdateBaseScore(c context.Context, tx *sql.Tx, ui *model.UserInfo) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_updateBaseScoreSQL, hitInfo(ui.Mid)), ui.BaseScore, ui.Score, ui.Mid); err != nil {
		log.Error("db.TxUpdateBaseScore(%d, %d, %d) error(%v)", ui.BaseScore, ui.Score, ui.Mid, err)
		return
	}
	return
}

//TxClearReliveTimes do clear user relivetimes.
func (d *Dao) TxClearReliveTimes(c context.Context, tx *sql.Tx, ui *model.UserInfo) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_clearReliveTimesSQL, hitInfo(ui.Mid)), ui.ReliveTimes, ui.Mid); err != nil {
		log.Error("db.TxClearReliveTimes(%d, %d) error(%v)", ui.ReliveTimes, ui.Mid, err)
		return
	}
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
	res = make(map[string]string)
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

// TxUpdateEventScoreReLive update event score and times
func (d *Dao) TxUpdateEventScoreReLive(c context.Context, tx *sql.Tx, mid int64, escore, score int8) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_updateEventScoreReLiveSQL, hitInfo(mid)), escore, score, mid); err != nil {
		log.Error("db.TxUpdateEventScoreReLive(%s, %d, %d, %d) error(%v)", _updateEventScoreReLiveSQL, escore, score, mid, err)
		return
	}
	return
}

//StatListByMid stat list by mid.
func (d *Dao) StatListByMid(c context.Context, mid int64) (list []*model.Statistics, err error) {
	var (
		rows *sql.Rows
	)
	list = make([]*model.Statistics, 0)
	if rows, err = d.db.Query(c, _statListByMidSQL, mid); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _statListByMidSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.Statistics{}
		if err = rows.Scan(&r.EventID, &r.Quantity); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, &model.Statistics{
			EventID:  r.EventID,
			Quantity: r.Quantity,
		})
	}
	return
}

//StatListByIDAndMid stat list by id.
func (d *Dao) StatListByIDAndMid(c context.Context, mid, id int64) (list []*model.Statistics, err error) {
	var (
		rows *sql.Rows
	)
	list = make([]*model.Statistics, 0)
	if rows, err = d.db.Query(c, _statListByIDAndMidSQL, mid, id); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _statListByIDAndMidSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.Statistics{}
		if err = rows.Scan(&r.EventID, &r.Quantity); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, &model.Statistics{
			EventID:  r.EventID,
			Quantity: r.Quantity,
		})
	}
	return
}

//StatListByID stat list by id.
func (d *Dao) StatListByID(c context.Context, id int64) (list []*model.Statistics, err error) {
	var (
		rows *sql.Rows
	)
	list = make([]*model.Statistics, 0)
	if rows, err = d.db.Query(c, _statListByIDSQL, id); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _statListByIDSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.Statistics{}
		if err = rows.Scan(&r.EventID, &r.Quantity); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, &model.Statistics{
			EventID:  r.EventID,
			Quantity: r.Quantity,
		})
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
		if err = rows.Scan(&event.ID, &event.Name, &event.NickName, &event.ServiceID, &event.Status, &event.CTime, &event.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, &model.Event{
			ID:        event.ID,
			Name:      event.Name,
			NickName:  event.NickName,
			ServiceID: event.ServiceID,
			Status:    event.Status,
			CTime:     event.CTime,
			MTime:     event.MTime,
		})
	}
	return
}

// TelLevel tel level.
func (d *Dao) TelLevel(c context.Context, mid int64) (res *model.TelRiskLevel, err error) {
	var (
		row *sql.Row
	)
	res = &model.TelRiskLevel{}
	row = d.db.QueryRow(c, _telLevelSQL, mid)
	if err = row.Scan(&res.ID, &res.Mid, &res.Level, &res.Origin, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// AddTelLevelInfo add tel level info.
func (d *Dao) AddTelLevelInfo(c context.Context, t *model.TelRiskLevel) (id int64, err error) {
	var (
		res xsql.Result
	)
	if res, err = d.db.Exec(c, _addTelLevelSQL, t.Mid, t.Level, t.Origin, time.Now()); err != nil {
		log.Error("d.db.Exec(%s, %v) error(%v)", _addTelLevelSQL, t, err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		log.Error("res.LastInsertId() error(%v)", err)
	}
	return
}
