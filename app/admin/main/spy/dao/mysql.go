package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/admin/main/spy/conf"
	"go-common/app/admin/main/spy/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_eventSQL              = "SELECT id,name,nick_name,service_id,ctime,mtime FROM spy_event WHERE name=? AND deleted = 0 LIMIT 1;"
	_allGroupSQL           = "SELECT id, name, ctime FROM spy_factor_group;"
	_updateFactorSQL       = "UPDATE spy_factor SET factor_val = ? WHERE id = ?;"
	_factorAllSQL          = "SELECT id,nick_name,service_id,event_id,group_id,risk_level,factor_val,ctime,mtime FROM spy_factor WHERE group_id = ? ORDER BY factor_val;"
	_getList               = "SELECT id,mid,event_id,score,base_score,event_score,remark,reason,factor_val,ctime FROM spy_user_event_history_%02d WHERE %s ORDER BY id DESC LIMIT %d,%d;"
	_getListTc             = "SELECT COUNT(1) FROM spy_user_event_history_%02d WHERE %s;"
	_addLogSQL             = "INSERT INTO `spy_log`(`name`,`module`,`context`,`ref_id`,`ctime`)VALUES(?,?,?,?,?);"
	_getSettingListSQL     = "SELECT id,property,name,val,ctime,mtime FROM spy_system_config"
	_updateSettingSQL      = "UPDATE spy_system_config SET val=? WHERE property=?"
	_getUserInfoSQL        = "SELECT id,mid,score,base_score,event_score,state,relive_times,mtime FROM spy_user_info_%02d WHERE mid=?"
	_addFactorSQL          = "INSERT INTO `spy_factor`(`nick_name`,`service_id`,`event_id`,`group_id`,`risk_level`,`factor_val`,`category_id`,`ctime`)VALUES(?,?,?,?,?,?,?,?);"
	_addEventSQL           = "INSERT INTO `spy_event`(`name`,`nick_name`,`service_id`,`status`,`ctime`,`mtime`)VALUES(?,?,?,?,?,?);"
	_addServiceSQL         = "INSERT INTO `spy_service`(`name`,`nick_name`,`status`,`ctime`,`mtime`)VALUES(?,?,?,?,?);"
	_addGroupSQL           = "INSERT INTO `spy_factor_group`(`name`,`ctime`)VALUES(?,?);"
	_getReportList         = "SELECT id, name, date_version, val, ctime FROM spy_report limit ?,?;"
	_getReportCount        = "SELECT COUNT(1) FROM spy_report;"
	_updateStatStateSQL    = "UPDATE spy_statistics SET state=? WHERE id=?"
	_updateStatQuantitySQL = "UPDATE spy_statistics SET quantity=? WHERE id=?"
	_updateStatDeleteSQL   = "UPDATE spy_statistics SET isdel=? WHERE id=?"
	_statByIDSQL           = "SELECT id,target_mid,target_id,event_id,state,type,quantity,isdel,ctime,mtime FROM spy_statistics WHERE id = ?;"
	_logListSQL            = "SELECT id,ref_id,name,module,context,ctime FROM spy_log WHERE ref_id = ? AND module = ?;"
	_statListByMidSQL      = "SELECT id,target_mid,target_id,event_id,state,type,quantity,isdel,ctime,mtime FROM spy_statistics WHERE target_mid = ? AND isdel = 0  ORDER BY id desc limit ?,?;"
	_statListByIDSQL       = "SELECT id,target_mid,target_id,event_id,state,type,quantity,isdel,ctime,mtime FROM spy_statistics WHERE target_id = ? AND type = ? AND isdel = 0  ORDER BY id desc limit ?,?;"
	_statCountByMidSQL     = "SELECT COUNT(1) FROM spy_statistics WHERE target_mid = ? AND isdel = 0;"
	_statCountByIDSQL      = "SELECT COUNT(1) FROM spy_statistics WHERE target_id = ? AND type = ? AND isdel = 0;"
	_allEventSQL           = "SELECT id,name,nick_name,service_id,status,ctime,mtime FROM spy_event WHERE status<>0"
	_updateEventNameSQL    = "UPDATE spy_event SET nick_name = ? WHERE id = ?;"
)

// Event get event from db.
func (d *Dao) Event(ctx context.Context, eventName string) (event *model.Event, err error) {
	var (
		row *xsql.Row
	)
	event = &model.Event{}
	row = d.eventStmt.QueryRow(ctx, eventName)
	if err = row.Scan(&event.ID, &event.Name, &event.NickName, &event.ServiceID, &event.CTime, &event.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// Factors get all factor info by group id.
func (d *Dao) Factors(c context.Context, gid int64) (res []*model.Factor, err error) {
	var rows *xsql.Rows
	if rows, err = d.factorAllStmt.Query(c, gid); err != nil {
		log.Error("d.allTypesStmt.Query(%d) error(%v)", gid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Factor)
		if err = rows.Scan(&r.ID, &r.NickName, &r.ServiceID, &r.EventID, &r.GroupID, &r.RiskLevel, &r.FactorVal, &r.CTime, &r.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// Groups get all group info.
func (d *Dao) Groups(c context.Context) (res []*model.FactorGroup, err error) {
	var rows *xsql.Rows
	if rows, err = d.allGroupStmt.Query(c); err != nil {
		log.Error("d.allGroupStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.FactorGroup)
		if err = rows.Scan(&r.ID, &r.Name, &r.CTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// UpdateFactor update factor.
func (d *Dao) UpdateFactor(c context.Context, factorVal float32, id int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateFactorSQL, factorVal, id); err != nil {
		log.Error("_updateFactorSQL: db.Exec(%v, %d) error(%v)", factorVal, id, err)
		return
	}
	return res.RowsAffected()
}

func hitHistory(id int64) int64 {
	return id % conf.Conf.Property.HistoryShard
}

// genListSQL get history sql
func (d *Dao) genListSQL(SQLType string, h *model.HisParamReq) (SQL string, values []interface{}) {
	values = make([]interface{}, 0, 1)
	cond := " mid = ?"
	values = append(values, h.Mid)
	switch SQLType {
	case "list":
		SQL = fmt.Sprintf(_getList, hitHistory(h.Mid), cond, (h.Pn-1)*h.Ps, h.Ps)
	case "count":
		SQL = fmt.Sprintf(_getListTc, hitHistory(h.Mid), cond)
	}
	return
}

//HistoryPage user event history.
func (d *Dao) HistoryPage(c context.Context, h *model.HisParamReq) (hs []*model.EventHistoryDto, err error) {
	SQL, values := d.genListSQL("list", h)
	rows, err := d.db.Query(c, SQL, values...)
	if err != nil {
		log.Error("dao.QuestionPage(%v,%v) error(%v)", SQL, values, err)
		return
	}
	defer rows.Close()
	hs = make([]*model.EventHistoryDto, 0, h.Ps)
	for rows.Next() {
		hdb := &model.EventHistory{}
		err = rows.Scan(&hdb.ID, &hdb.Mid, &hdb.EventID, &hdb.Score, &hdb.BaseScore, &hdb.EventScore,
			&hdb.Remark, &hdb.Reason, &hdb.FactorVal, &hdb.Ctime)
		eventMSG := &model.EventMessage{}
		if err = json.Unmarshal([]byte(hdb.Remark), eventMSG); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", hdb.Remark, err)
		} else {
			hdb.TargetID = eventMSG.TargetID
			hdb.TargetMid = eventMSG.TargetMid
		}
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		h := &model.EventHistoryDto{
			ID:         hdb.ID,
			Score:      hdb.Score,
			BaseScore:  hdb.BaseScore,
			EventScore: hdb.EventScore,
			Reason:     hdb.Reason,
			Ctime:      hdb.Ctime.Unix(),
			TargetID:   hdb.TargetID,
			TargetMid:  hdb.TargetMid,
		}
		if eventMSG.Time != 0 {
			_, offset := time.Now().Zone()
			t := time.Unix(eventMSG.Time, 0).Add(-time.Duration(offset) * time.Second)
			h.SpyTime = t.Unix()
		}
		hs = append(hs, h)
	}
	return
}

// HistoryPageTotalC user ecent history page.
func (d *Dao) HistoryPageTotalC(c context.Context, h *model.HisParamReq) (totalCount int, err error) {
	SQL, values := d.genListSQL("count", h)
	row := d.db.QueryRow(c, SQL, values...)
	if err = row.Scan(&totalCount); err != nil {
		if err == sql.ErrNoRows {
			row = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

//AddLog add log.
func (d *Dao) AddLog(c context.Context, l *model.Log) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _addLogSQL, l.Name, l.Module, l.Context, l.RefID, l.Ctime); err != nil {
		fmt.Println("add log ", err)
		log.Error("add question: d.db.Exec(%v) error(%v)", l, err)
		return
	}
	return res.RowsAffected()
}

// SettingList get all setting list
func (d *Dao) SettingList(c context.Context) (list []*model.Setting, err error) {
	var (
		rows *xsql.Rows
	)
	list = make([]*model.Setting, 0)
	if rows, err = d.db.Query(c, _getSettingListSQL); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _getSettingListSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var setting = &model.Setting{}
		if err = rows.Scan(&setting.ID, &setting.Property, &setting.Name, &setting.Val, &setting.CTime, &setting.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, setting)
	}
	return
}

// UpdateSetting update setting
func (d *Dao) UpdateSetting(c context.Context, property string, val string) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateSettingSQL, val, property); err != nil {
		log.Error("d.db.Exec(%s,%d,%s) error(%v)", _updateSettingSQL, val, property)
		return
	}
	if affected, err = res.RowsAffected(); err != nil {
		log.Error("res.RowsAffected() error(%v)", err)
		return
	}
	return
}

func hitInfo(id int64) int64 {
	return id % conf.Conf.Property.UserInfoShard
}

// Info get lastest user info by mid.
func (d *Dao) Info(c context.Context, mid int64) (res *model.UserInfo, err error) {
	res = &model.UserInfo{}
	hitIndex := hitInfo(mid)
	row := d.getUserInfoStmt[hitIndex].QueryRow(c, mid)
	if err = row.Scan(&res.ID, &res.Mid, &res.Score, &res.BaseScore, &res.EventScore, &res.State, &res.ReliveTimes,
		&res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
			return
		}
		log.Error("Info:row.Scan() error(%v)", err)
	}
	return
}

//AddFactor add factor.
func (d *Dao) AddFactor(c context.Context, f *model.Factor) (ret int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _addFactorSQL, f.NickName, f.ServiceID, f.EventID, f.GroupID, f.RiskLevel, f.FactorVal, f.CategoryID, f.CTime); err != nil {
		log.Error("d.db AddFactor: d.db.Exec(%v) error(%v)", f, err)
		return
	}
	ret, err = res.RowsAffected()
	return
}

//AddEvent add event.
func (d *Dao) AddEvent(c context.Context, f *model.Event) (ret int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _addEventSQL, f.Name, f.NickName, f.ServiceID, f.Status, f.CTime, f.MTime); err != nil {
		log.Error("d.db AddEvent: d.db.Exec(%v) error(%v)", f, err)
		return
	}
	ret, err = res.RowsAffected()
	return
}

//AddService add service.
func (d *Dao) AddService(c context.Context, f *model.Service) (ret int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _addServiceSQL, f.Name, f.NickName, f.Status, f.CTime, f.MTime); err != nil {
		log.Error("d.db AddService: d.db.Exec(%v) error(%v)", f, err)
		return
	}
	ret, err = res.RowsAffected()
	return
}

//AddGroup add group.
func (d *Dao) AddGroup(c context.Context, f *model.FactorGroup) (ret int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _addGroupSQL, f.Name, f.CTime); err != nil {
		log.Error("d.db AddGroup: d.db.Exec(%v) error(%v)", f, err)
		return
	}
	ret, err = res.RowsAffected()
	return
}

// ReportList report list.
func (d *Dao) ReportList(c context.Context, ps, pn int) (list []*model.ReportDto, err error) {
	var (
		rows *xsql.Rows
	)
	if ps == 0 || pn == 0 {
		ps = 8
		pn = 1
	}
	list = make([]*model.ReportDto, 0)
	if rows, err = d.db.Query(c, _getReportList, (pn-1)*ps, ps); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _getReportList, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.Report{}
		if err = rows.Scan(&r.ID, &r.Name, &r.DateVersion, &r.Val, &r.Ctime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, &model.ReportDto{
			ID:          r.ID,
			Name:        r.Name,
			DateVersion: r.DateVersion,
			Val:         r.Val,
			Ctime:       r.Ctime.Unix(),
		})
	}
	return
}

// ReportCount get repoet total count.
func (d *Dao) ReportCount(c context.Context) (totalCount int, err error) {
	var row = d.db.QueryRow(c, _getReportCount)
	if err = row.Scan(&totalCount); err != nil {
		if err == sql.ErrNoRows {
			row = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// UpdateStatState update stat state.
func (d *Dao) UpdateStatState(c context.Context, state int8, id int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateStatStateSQL, state, id); err != nil {
		log.Error("d.db.Exec(%s,%d,%s) error(%v)", _updateStatStateSQL, state, id)
		return
	}
	if affected, err = res.RowsAffected(); err != nil {
		log.Error("res.RowsAffected() error(%v)", err)
		return
	}
	return
}

// UpdateStatQuantity update stat quantity
func (d *Dao) UpdateStatQuantity(c context.Context, count int64, id int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateStatQuantitySQL, count, id); err != nil {
		log.Error("d.db.Exec(%s,%d,%s) error(%v)", _updateStatQuantitySQL, count, id)
		return
	}
	if affected, err = res.RowsAffected(); err != nil {
		log.Error("res.RowsAffected() error(%v)", err)
		return
	}
	return
}

// DeleteStat delete stat.
func (d *Dao) DeleteStat(c context.Context, isdel int8, id int64) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateStatDeleteSQL, isdel, id); err != nil {
		log.Error("d.db.Exec(%s,%d,%s) error(%v)", _updateStatDeleteSQL, isdel, id)
		return
	}
	if affected, err = res.RowsAffected(); err != nil {
		log.Error("res.RowsAffected() error(%v)", err)
		return
	}
	return
}

// Statistics get stat info by id from db.
func (d *Dao) Statistics(c context.Context, id int64) (stat *model.Statistics, err error) {
	var (
		row *xsql.Row
	)
	stat = &model.Statistics{}
	row = d.db.QueryRow(c, _statByIDSQL, id)
	if err = row.Scan(&stat.ID, &stat.TargetMid, &stat.TargetID, &stat.EventID, &stat.State, &stat.Type, &stat.Quantity, &stat.Isdel, &stat.Ctime, &stat.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			stat = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

//LogList log list.
func (d *Dao) LogList(c context.Context, refID int64, module int8) (list []*model.Log, err error) {
	var (
		rows *xsql.Rows
	)
	list = make([]*model.Log, 0)
	if rows, err = d.db.Query(c, _logListSQL, refID, module); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _logListSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = &model.Log{}
		if err = rows.Scan(&r.ID, &r.RefID, &r.Name, &r.Module, &r.Context, &r.Ctime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, &model.Log{
			ID:        r.ID,
			RefID:     r.RefID,
			Name:      r.Name,
			Module:    r.Module,
			Context:   r.Context,
			CtimeUnix: r.Ctime.Unix(),
		})
	}
	return
}

//StatListByMid stat list by mid.
func (d *Dao) StatListByMid(c context.Context, mid int64, pn, ps int) (list []*model.Statistics, err error) {
	var (
		rows *xsql.Rows
	)
	if ps == 0 || pn == 0 {
		ps = 8
		pn = 1
	}
	list = make([]*model.Statistics, 0)
	if rows, err = d.db.Query(c, _statListByMidSQL, mid, (pn-1)*ps, ps); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _statListByMidSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var stat = &model.Statistics{}
		if err = rows.Scan(&stat.ID, &stat.TargetMid, &stat.TargetID, &stat.EventID, &stat.State, &stat.Type, &stat.Quantity, &stat.Isdel, &stat.Ctime, &stat.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, &model.Statistics{
			ID:        stat.ID,
			TargetMid: stat.TargetMid,
			TargetID:  stat.TargetID,
			EventID:   stat.EventID,
			Type:      stat.Type,
			State:     stat.State,
			Quantity:  stat.Quantity,
			Isdel:     stat.Isdel,
			Ctime:     stat.Ctime,
			Mtime:     stat.Mtime,
			CtimeUnix: stat.Ctime.Unix(),
			MtimeUnix: stat.Mtime.Unix(),
		})
	}
	return
}

//StatListByID stat list by id.
func (d *Dao) StatListByID(c context.Context, id int64, t int8, pn, ps int) (list []*model.Statistics, err error) {
	var (
		rows *xsql.Rows
	)
	if ps == 0 || pn == 0 {
		ps = 8
		pn = 1
	}
	list = make([]*model.Statistics, 0)
	if rows, err = d.db.Query(c, _statListByIDSQL, id, t, (pn-1)*ps, ps); err != nil {
		log.Error("d.db.Query(%s) error(%v)", _statListByIDSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var stat = &model.Statistics{}
		if err = rows.Scan(&stat.ID, &stat.TargetMid, &stat.TargetID, &stat.EventID, &stat.State, &stat.Type, &stat.Quantity, &stat.Isdel, &stat.Ctime, &stat.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		list = append(list, &model.Statistics{
			ID:        stat.ID,
			TargetMid: stat.TargetMid,
			TargetID:  stat.TargetID,
			EventID:   stat.EventID,
			Type:      stat.Type,
			State:     stat.State,
			Quantity:  stat.Quantity,
			Isdel:     stat.Isdel,
			Ctime:     stat.Ctime,
			Mtime:     stat.Mtime,
			CtimeUnix: stat.Ctime.Unix(),
			MtimeUnix: stat.Mtime.Unix(),
		})
	}
	return
}

// StatCountByMid count by mid.
func (d *Dao) StatCountByMid(c context.Context, mid int64) (totalCount int64, err error) {
	row := d.db.QueryRow(c, _statCountByMidSQL, mid)
	if err = row.Scan(&totalCount); err != nil {
		if err == sql.ErrNoRows {
			row = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// StatCountByID count by id.
func (d *Dao) StatCountByID(c context.Context, id int64, t int8) (totalCount int64, err error) {
	row := d.db.QueryRow(c, _statCountByIDSQL, id, t)
	if err = row.Scan(&totalCount); err != nil {
		if err == sql.ErrNoRows {
			row = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

//AllEvent all event.
func (d *Dao) AllEvent(c context.Context) (list []*model.Event, err error) {
	var (
		rows *xsql.Rows
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

// UpdateEventName update event name.
func (d *Dao) UpdateEventName(c context.Context, e *model.Event) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateEventNameSQL, e.NickName, e.ID); err != nil {
		log.Error("_updateEventNameSQL: db.Exec(%v, %d) error(%v)", e.NickName, e.ID, err)
		return
	}
	return res.RowsAffected()
}
