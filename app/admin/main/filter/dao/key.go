package dao

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/admin/main/filter/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// insert
	_conKey       = "SELECT id,`key`,mode,filter,level,comment,etime,stime FROM filter_content WHERE filter=? AND `key`=?"
	_conKeyByID   = "SELECT id,`key`,mode,filter,level,comment,etime,stime FROM filter_content WHERE id=? AND `key`=?"
	_insertConkey = "INSERT IGNORE INTO filter_content (mode,filter,`key`,comment,level,stime,etime) VALUES (?,?,?,?,?,?,?)"
	_inserKey     = "INSERT IGNORE INTO filter_key (area,`key`,filterid)VALUES(?,?,?) ON DUPLICATE KEY UPDATE state=0"

	// del
	_delCon    = "UPDATE filter_content SET state=1 WHERE id=?"
	_delKeyFid = "UPDATE filter_key SET state=1 WHERE `key`=? AND filterid=?"

	// editor
	_updateConkey = "UPDATE filter_content SET mode=?,comment=?,level=?,state=0,stime=?,etime=?,mtime=now() WHERE `key`=? AND filter=?"

	// search
	_countKey      = "SELECT count(*) FROM filter_content WHERE `key`!=''"
	_searchKey     = "SELECT id,`key`,mode,filter,level,comment,etime,stime FROM filter_content WHERE `key`!=''"
	_searchKeyArea = "SELECT area FROM filter_key WHERE `key`=? AND filterid=? AND state=0"

	//log
	_insertFkLog = "INSERT INTO filter_key_log (`key`,adid,name,comment,state)VALUES(?,?,?,?,?)"
	_fkLogs      = "SELECT `key`,adid,name,comment,state,ctime FROM filter_key_log WHERE `key`=? ORDER BY id DESC"
)

// InsertFkLog .
func (d *Dao) InsertFkLog(c context.Context, key, name, comment string, adid int64, state int8) (id int64, err error) {
	var result sql.Result
	if result, err = d.insertFkLogStmt.Exec(c, key, adid, name, comment, state); err != nil {
		log.Error("d.insertFkLogStmt.Exec(%s,%s,%s,%d,%d) error(%v)", key, adid, name, comment, state, err)
		return
	}
	return result.LastInsertId()
}

// FkLogs .
func (d *Dao) FkLogs(c context.Context, key string) (ls []*model.Log, err error) {
	rows, err := d.fkLogsStmt.Query(c, key)
	if err != nil {
		log.Error("d.fkLogsStmt.Query(%s) err(%v)", key, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		l := &model.Log{}
		if err = rows.Scan(&l.Key, &l.AdminID, &l.Name, &l.Comment, &l.State, &l.Ctime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			continue
		}
		ls = append(ls, l)
	}
	return
}

// ConKey .
func (d *Dao) ConKey(c context.Context, filter, key string) (r *model.KeyInfo, err error) {
	var row = d.ConKeyStmt.QueryRow(c, filter, key)
	r = &model.KeyInfo{}
	if err = row.Scan(&r.ID, &r.Key, &r.Mode, &r.Filter, &r.Level, &r.Comment, &r.Etime, &r.Stime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// ConKeyByID .
func (d *Dao) ConKeyByID(c context.Context, id int64, key string) (r *model.KeyInfo, err error) {
	var row = d.ConKeyIDStmt.QueryRow(c, id, key)
	r = &model.KeyInfo{}
	if err = row.Scan(&r.ID, &r.Key, &r.Mode, &r.Filter, &r.Level, &r.Comment, &r.Etime, &r.Stime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// InsertConkey .
func (d *Dao) InsertConkey(c context.Context, filter, key, comment string, mode, level int8, stime, etime int64) (id int64, err error) {
	var result sql.Result
	if result, err = d.insertConkeyStmt.Exec(c, mode, filter, key, comment, level, time.Unix(stime, 0), time.Unix(etime, 0)); err != nil {
		log.Error("d.insertConkeyStmt.Exec(%d,%s,%s,%s,%d,%d,%d) error(%v)", mode, filter, key, comment, level, stime, etime, err)
		return
	}
	return result.LastInsertId()
}

// InsertKey .
func (d *Dao) InsertKey(c context.Context, area, key string, filterID int64) (id int64, err error) {
	var result sql.Result
	if result, err = d.inserKeyStmt.Exec(c, area, key, filterID); err != nil {
		log.Error("d.inserKeyStmt.Exec(%s,%s,%d) error(%v)", area, key, filterID, err)
		return
	}
	return result.LastInsertId()
}

// TxDelCon .
func (d *Dao) TxDelCon(tx *xsql.Tx, id int64) (rows int64, err error) {
	var result sql.Result
	if result, err = tx.Exec(_delCon, id); err != nil {
		log.Error("tx.Exec(%d) error(%v)", id, err)
		return
	}
	return result.RowsAffected()
}

// TxDelKeyFid .
func (d *Dao) TxDelKeyFid(tx *xsql.Tx, key string, fid int64) (rows int64, err error) {
	var result sql.Result
	if result, err = tx.Exec(_delKeyFid, key, fid); err != nil {
		log.Error("tx.Exec(%s,%d) error(%v)", key, fid, err)
		return
	}
	return result.RowsAffected()
}

// DelKeyFid .
func (d *Dao) DelKeyFid(c context.Context, key string, fid int64) (rows int64, err error) {
	var result sql.Result
	if result, err = d.delKeyFidStmt.Exec(c, key, fid); err != nil {
		log.Error("d.delKeyFidStmt(%s,%d) error(%v)", key, fid, err)
		return
	}
	return result.RowsAffected()
}

// UpdateConkey .
func (d *Dao) UpdateConkey(c context.Context, filter, key, comment string, mode, level int8, stime, etime int64) (affected int64, err error) {
	var result sql.Result
	if result, err = d.updateConKeyStmt.Exec(c, mode, comment, level, time.Unix(stime, 0), time.Unix(etime, 0), key, filter); err != nil {
		log.Error("d.updateConKeyStmt.Exec(%d,%s,%d,%d,%d,%s,%s)", mode, comment, level, stime, etime, filter, key)
		return
	}
	return result.RowsAffected()
}

// CountKey .
func (d *Dao) CountKey(c context.Context, key, comment string, state int8) (total int64, err error) {
	var (
		querySQL string
		row      *xsql.Row
	)
	querySQL = _countKey
	if key != "" {
		querySQL += " AND `key`='" + key + "'"
	}
	querySQL += fmt.Sprintf(" AND state=%d", state)
	if comment != "" {
		querySQL += " AND comment like '%" + comment + "%'"
	}
	if row = d.mysql.QueryRow(c, querySQL); err != nil {
		log.Error("d.mysql.QueryRow(%s) error(%v)", key, err)
		return
	}
	if err = row.Scan(&total); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// SearchKey .
func (d *Dao) SearchKey(c context.Context, key, comment string, start, limit int64, state int8) (rs []*model.KeyInfo, err error) {
	var (
		querySQL string
		rows     *xsql.Rows
	)
	querySQL = _searchKey
	if key != "" {
		querySQL += " AND `key`='" + key + "'"
	}
	querySQL += fmt.Sprintf(" AND state=%d", state)
	if comment != "" {
		querySQL += " AND comment like \"%" + comment + "%\""
	}
	querySQL += " ORDER BY mtime DESC" + " LIMIT " + fmt.Sprintf("%d", start) + "," + fmt.Sprintf("%d", limit)
	if rows, err = d.mysql.Query(c, querySQL); err != nil {
		log.Error("d.mysql.Query(%s) error(%v)", querySQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.KeyInfo{}
		if err = rows.Scan(&r.ID, &r.Key, &r.Mode, &r.Filter, &r.Level, &r.Comment, &r.Etime, &r.Stime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	return
}

// KeyArea .
func (d *Dao) KeyArea(c context.Context, key string, fid int64) (areas []string, err error) {
	var rows *xsql.Rows
	if rows, err = d.searchKeyAreaStmt.Query(c, key, fid); err != nil {
		log.Error("d.searchKeyAreaStmt.Query(%s,%d) error(%v)", key, fid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var area string
		if err = rows.Scan(&area); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		areas = append(areas, area)
	}
	return
}

// func (d *Dao) KeyAreaRule(c context.Context, key string, areas []string) (rs []*model.KeyInfo, err error) {
// 	var (
// 		querySQL string
// 		rows     *xsql.Rows
// 	)
// 	if rows, err = d.mysql.Query(c, fmt.Sprintf(_areaKeyRule, d.CoverStr(areas)), key); err != nil {
// 		log.Error("d.mysql.Query(%s,%s) error(%v)", querySQL, key, err)
// 		return
// 	}
// 	for rows.Next() {
// 		r := &model.KeyInfo{}
// 		if err = rows.Scan(&r.FkID, &r.Fid, &r.Area, &r.Key, &r.Filter, &r.Mode, &r.Level, &r.Comment, &r.Etime, &r.Stime); err != nil {
// 			log.Error("rows.Scan() error(%v)", err)
// 			return
// 		}
// 		rs = append(rs, r)
// 	}
// 	return
// }

// CoverStr .
func (d *Dao) CoverStr(strs []string) string {
	var buf = bytes.NewBuffer(nil)
	for _, str := range strs {
		buf.WriteString("'")
		buf.WriteString(str)
		buf.WriteString("'")
		buf.WriteString(",")
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}
