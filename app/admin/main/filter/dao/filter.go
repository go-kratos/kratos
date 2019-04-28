package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/admin/main/filter/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// filter search
	_searchCount = `SELECT count(*) FROM filter_content a
	INNER JOIN filter_area b ON b. area=? AND b.is_delete=? AND a.id =b.filterid WHERE a.state=? AND a.type IN (%s) AND a.source IN (%s) AND b.level IN (%s) AND a.filter LIKE ?`
	// TODO 目前版本先保留limit a,b ，暂时不会导致性能问题，在where语句后数据集 < 3k，下个版本会改掉(v1.1.0 以后)
	_search = `SELECT a.id,a.mode,a.filter,a.comment,a.level, b.level, a.source, a.type, a.stime,a.etime,a.ctime,b.typeid FROM filter_content a
	INNER JOIN filter_area b ON b. area=? AND b.is_delete=? AND a.id =b.filterid WHERE a.state=? AND a.type IN (%s) AND a.source IN (%s) AND b.level IN (%s) AND a.filter LIKE ? ORDER BY ID DESC LIMIT ?,? `
	_searchCountWithoutMsg = `SELECT count(*) FROM filter_content a
	INNER JOIN filter_area b ON b. area=? AND b.is_delete=? AND a.id =b.filterid WHERE a.state=? AND a.type IN (%s) AND a.source IN (%s) AND b.level IN (%s)`
	// TODO 目前版本先保留limit a,b ，暂时不会导致性能问题，在where语句后数据集 < 3k，下个版本会改掉(v1.1.0 以后)
	_searchWithoutMsg = `SELECT a.id,a.mode,a.filter,a.comment,a.level,b.level,a.source, a.type, a.stime,a.etime,a.ctime,b.typeid FROM filter_content a
	INNER JOIN filter_area b ON b. area=? AND b.is_delete=? AND a.id =b.filterid WHERE a.state=? AND a.type IN (%s) AND a.source IN (%s) AND b.level IN (%s) ORDER BY ID DESC LIMIT ?,?`

	_maxFilterID         = `SELECT id FROM filter_content ORDER BY id DESC LIMIT 1`
	_expiredRuleIDs      = `SELECT id FROM filter_content WHERE state=0 AND id BETWEEN ? AND ? AND etime<?`
	_ruleID              = "SELECT id from filter_content WHERE filter=? AND `key`='' LIMIT 1"
	_validRuleIDsByRange = "SELECT id from filter_content WHERE id BETWEEN ? AND ? AND `key`='' AND state=0"
	_validAreaRuleCount  = `SELECT count(*) FROM filter_area AS a INNER JOIN filter_content AS b ON a.filterid=b.id WHERE
	b.stime< ? AND b.etime>? AND a.area=? AND a.is_delete=0`
	_validAreaRuleCountByRuleID = `SELECT count(*) FROM filter_area WHERE filterid=? AND is_delete=0`
	_ruleAreas                  = `SELECT area,typeid,level FROM filter_area WHERE is_delete=0 AND filterid=?`
	_rule                       = `SELECT id,mode,filter,level,comment,stime,etime,source,type,state,ctime FROM filter_content WHERE id=?`
	_ruleByContent              = "SELECT id,mode,filter,level,comment,stime,etime,source,type,state,ctime FROM filter_content WHERE filter=? AND `key`=''"
	_updateRulesState           = `UPDATE filter_content SET state=? WHERE id IN(%s)`
	_updateRuleState            = `UPDATE filter_content SET state=? WHERE id=?`
	_deleteAreaRules            = `UPDATE filter_area SET is_delete=1 WHERE filterid = ?`
	_upsertRule                 = `INSERT INTO filter_content (mode,filter,comment,level,source,type,stime,etime,state) VALUES(?,?,?,?,?,?,?,?,0) ON DUPLICATE KEY UPDATE mode=?,filter=?,comment=?,level=?,source=?,type=?,stime=?,etime=?,state=0`
	_updateRule                 = "UPDATE filter_content SET mode=?,filter=?,comment=?,level=?,source=?,type=?,stime=?,etime=? WHERE id=?"
	_upsertAreaRule             = `INSERT INTO filter_area (area,typeid,filterid,level) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE level=?,is_delete=0`
	// log
	_insertLog = `INSERT INTO filter_log (filterid,adid,comment,name,state) VALUES(?,?,?,?,?)`
	_logs      = `SELECT adid,name,comment,state,ctime FROM filter_log WHERE filterid=?`
)

// InsertLog insert filter log
func (d *Dao) InsertLog(c context.Context, tx *xsql.Tx, id, adid int64, comment, name string, state int8) (newID int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_insertLog, id, adid, comment, name, state); err != nil {
		return
	}
	return res.RowsAffected()
}

// Logs get filter logs by filterID
func (d *Dao) Logs(c context.Context, filterID int64) (logs []*model.Log, err error) {
	var rows *xsql.Rows
	logs = make([]*model.Log, 0)
	if rows, err = d.mysql.Query(c, _logs, filterID); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		log := &model.Log{}
		if err = rows.Scan(&log.AdminID, &log.Name, &log.Comment, &log.State, &log.Ctime); err != nil {
			return
		}
		logs = append(logs, log)
	}
	err = rows.Err()
	return
}

// Filter get filter by id
func (d *Dao) Filter(c context.Context, id int64) (rule *model.FilterInfo, err error) {
	rule = &model.FilterInfo{}
	row := d.mysql.QueryRow(c, _rule, id)
	if err = row.Scan(&rule.ID, &rule.Mode, &rule.Filter, &rule.Level, &rule.Comment, &rule.Stime, &rule.Etime, &rule.Source, &rule.Type, &rule.State, &rule.CTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			rule = nil
		}
		return
	}
	var rows *xsql.Rows
	if rows, err = d.mysql.Query(c, _ruleAreas, id); err != nil {
		return
	}
	defer rows.Close()
	var (
		areaM = make(map[string]*model.FilterArea)
		tpM   = make(map[int64]struct{})
	)
	for rows.Next() {
		var (
			area      string
			tpid      int64
			areaLevel int8
		)
		if err = rows.Scan(&area, &tpid, &areaLevel); err != nil {
			return
		}
		areaM[area] = &model.FilterArea{Area: area, Level: areaLevel}
		tpM[tpid] = struct{}{}
	}
	if err = rows.Err(); err != nil {
		return
	}
	for _, a := range areaM {
		rule.Areas = append(rule.Areas, a)
	}
	for t := range tpM {
		rule.TpIDs = append(rule.TpIDs, t)
	}
	return
}

// FilterByContent get filter by content
func (d *Dao) FilterByContent(c context.Context, content string) (rule *model.FilterInfo, err error) {
	rule = &model.FilterInfo{}
	row := d.mysql.QueryRow(c, _ruleByContent, content)
	if err = row.Scan(&rule.ID, &rule.Mode, &rule.Filter, &rule.Level, &rule.Comment, &rule.Stime, &rule.Etime, &rule.Source, &rule.Type, &rule.State, &rule.CTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			rule = nil
		}
		return
	}
	var rows *xsql.Rows
	if rows, err = d.mysql.Query(c, _ruleAreas, rule.ID); err != nil {
		return
	}
	defer rows.Close()
	var (
		areaM = make(map[string]*model.FilterArea)
		tpM   = make(map[int64]struct{})
	)
	for rows.Next() {
		var (
			area      string
			tpid      int64
			areaLevel int8
		)
		if err = rows.Scan(&area, &tpid, &areaLevel); err != nil {
			return
		}
		areaM[area] = &model.FilterArea{Area: area, Level: areaLevel}
		tpM[tpid] = struct{}{}
	}
	if err = rows.Err(); err != nil {
		return
	}
	for _, a := range areaM {
		rule.Areas = append(rule.Areas, a)
	}
	for t := range tpM {
		rule.TpIDs = append(rule.TpIDs, t)
	}
	return
}

// ValidFilterAreaCountByArea get area rule count which still working (not deleted , not expired)
func (d *Dao) ValidFilterAreaCountByArea(c context.Context, area string) (count int64, err error) {
	row := d.mysql.QueryRow(c, _validAreaRuleCount, time.Now(), time.Now(), area)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// ValidFilterAreaCountByRuleID get area rule count which still working (just check is_deleted)
func (d *Dao) ValidFilterAreaCountByRuleID(c context.Context, ruleID int64) (count int64, err error) {
	row := d.mysql.QueryRow(c, _validAreaRuleCountByRuleID, ruleID)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// Search .
func (d *Dao) Search(c context.Context, msg, area, sourceStr, typeStr, levelStr string, stage, deleted int, start int64, end int64) (rules []*model.FilterInfo, err error) {
	var (
		rows *xsql.Rows
	)
	if msg == "" {
		rows, err = d.mysql.Query(c, fmt.Sprintf(_searchWithoutMsg, typeStr, sourceStr, levelStr), area, deleted, stage, start, end)
	} else {
		msg = fmt.Sprintf("%%%s%%", msg)
		rows, err = d.mysql.Query(c, fmt.Sprintf(_search, typeStr, sourceStr, levelStr), area, deleted, stage, msg, start, end)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			rule      = &model.FilterInfo{}
			tpID      int64
			areaLevel int8
		)
		if err = rows.Scan(&rule.ID, &rule.Mode, &rule.Filter, &rule.Comment, &rule.Level, &areaLevel, &rule.Source, &rule.Type, &rule.Stime, &rule.Etime, &rule.CTime, &tpID); err != nil {
			return
		}
		rule.TpIDs = []int64{tpID}
		rule.Areas = append(rule.Areas, &model.FilterArea{Area: area, Level: areaLevel})
		rules = append(rules, rule)
	}
	err = rows.Err()
	return
}

// SearchCount .
func (d *Dao) SearchCount(c context.Context, msg string, area string, sourceStr string, typeStr string, levelStr string, stage, deleted int) (count int64, err error) {
	var (
		row *xsql.Row
	)
	if msg == "" {
		row = d.mysql.QueryRow(c, fmt.Sprintf(_searchCountWithoutMsg, typeStr, sourceStr, levelStr), area, deleted, stage)
	} else {
		msg = fmt.Sprintf("%%%s%%", msg)
		row = d.mysql.QueryRow(c, fmt.Sprintf(_searchCount, typeStr, sourceStr, levelStr), area, deleted, stage, msg)
	}
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// FilterID get filter_content id by filter content .
func (d *Dao) FilterID(c context.Context, filter string) (id int64, err error) {
	row := d.mysql.QueryRow(c, _ruleID, filter)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// ValidFilterIDsByRange .
func (d *Dao) ValidFilterIDsByRange(c context.Context, begin, end int64) (ids []int64, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.mysql.Query(c, _validRuleIDsByRange, begin, end); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.Scan err(%+v)", err)
			return
		}
		ids = append(ids, id)
	}
	err = rows.Err()
	return
}

// UpsertRule insert or update rule (filter_content)
func (d *Dao) UpsertRule(c context.Context, tx *xsql.Tx, filter, comment string, level, mode, source, keyType int8, stime, etime time.Time) (newID int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_upsertRule, mode, filter, comment, level, source, keyType, stime, etime, mode, filter, comment, level, source, keyType, stime, etime); err != nil {
		return
	}
	return res.LastInsertId()
}

// UpdateRule update rule(filter_content)
func (d *Dao) UpdateRule(c context.Context, tx *xsql.Tx, id int64, filter, comment string, mode, level, source, keyType int8, stime, etime time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateRule, mode, filter, comment, level, source, keyType, stime, etime, id); err != nil {
		return
	}
	return res.RowsAffected()
}

// UpsertAreaRule insert or update a filter_area
func (d *Dao) UpsertAreaRule(c context.Context, tx *xsql.Tx, area string, tpid int64, filterID int64, level int8) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_upsertAreaRule, area, tpid, filterID, level, level); err != nil {
		return
	}
	return res.RowsAffected()
}

// DeleteAreaRules delete area rules (filter_area) by filterID
func (d *Dao) DeleteAreaRules(c context.Context, tx *xsql.Tx, filterID int64) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_deleteAreaRules, filterID); err != nil {
		return
	}
	return res.RowsAffected()
}

// MaxFilterID get max filter id.
func (d *Dao) MaxFilterID(c context.Context) (res int64, err error) {
	row := d.mysql.QueryRow(c, _maxFilterID)
	if err = row.Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// ExpiredRuleIDs get expired filter ids
func (d *Dao) ExpiredRuleIDs(c context.Context, startID, endID int) (ids []int64, err error) {
	var (
		cur  = time.Now()
		rows *xsql.Rows
	)
	if rows, err = d.mysql.Query(c, _expiredRuleIDs, startID, endID, cur); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return
		}
		ids = append(ids, id)
	}
	err = rows.Err()
	return
}

// UpdateRulesState update all rule ids' state
func (d *Dao) UpdateRulesState(c context.Context, ids []int64, state int8) (affected int64, err error) {
	var (
		res    sql.Result
		sqlStr = fmt.Sprintf(_updateRulesState, xstr.JoinInts(ids))
	)
	if res, err = d.mysql.Exec(c, sqlStr, state); err != nil {
		return
	}
	return res.RowsAffected()
}

// TxUpdateRuleState update all rule ids' state
func (d *Dao) TxUpdateRuleState(c context.Context, tx *xsql.Tx, id int64, state int8) (affected int64, err error) {
	var (
		res sql.Result
	)
	if res, err = tx.Exec(_updateRuleState, state, id); err != nil {
		return
	}
	return res.RowsAffected()
}

// UpdateRuleState update all rule ids' state
func (d *Dao) UpdateRuleState(c context.Context, id int64, state int8) (affected int64, err error) {
	var (
		res sql.Result
	)
	if res, err = d.mysql.Exec(c, _updateRuleState, state, id); err != nil {
		return
	}
	return res.RowsAffected()
}
