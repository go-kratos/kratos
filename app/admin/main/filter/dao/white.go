package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/admin/main/filter/model"
	xsql "go-common/library/database/sql"
)

const (
	// get white content
	_whiteContent = "SELECT id,content,mode,comment,state FROM filter_white_content WHERE content=?"

	// get white info
	_whiteInfo = "SELECT fwc.id,fwc.content,fwc.mode,fwc.comment,fwa.area,fwa.tpid FROM filter_white_content AS fwc INNER JOIN filter_white_area AS fwa ON fwc.id=fwa.content_id WHERE fwc.id=? AND fwa.state=0"

	// search
	_searchWhiteContent = "SELECT id,content,mode,comment FROM filter_white_content WHERE id IN(SELECT content_id FROM filter_white_area WHERE state=0 %s) AND state=0"
	_countWhiteContent  = "SELECT count(id) FROM filter_white_content WHERE id IN(SELECT content_id FROM filter_white_area WHERE state=0 %s) AND state=0"

	_upsertWhiteContent = "INSERT INTO filter_white_content (content,mode,comment,state) VALUES (?,?,?,0) ON DUPLICATE KEY UPDATE mode=?,comment=?,state=0"
	_updateWhiteContent = "UPDATE filter_white_content SET mode=?,comment=? WHERE content=?"
	_deleteWhiteContent = "UPDATE filter_white_content SET state=1 WHERE id=?"

	// area white
	_upsertAreaWhite = "INSERT filter_white_area (area,tpid,content_id,state)VALUES(?,?,?,0) ON DUPLICATE KEY UPDATE state=0"
	_deleteAreaWhite = "UPDATE filter_white_area SET state=1 WHERE content_id=?"
	_whiteAreas      = "SELECT area FROM filter_white_area WHERE state=0 AND content_id=?"

	// log
	_whiteLog       = "SELECT adid,name,comment,state,ctime FROM filter_white_log WHERE content_id=?"
	_insertWhiteLog = "INSERT INTO filter_white_log (content_id,adid,name,comment,state) VALUES (?,?,?,?,?)"
)

// WhiteContent get white content
func (d *Dao) WhiteContent(c context.Context, content string) (rule *model.WhiteInfo, err error) {
	rule = &model.WhiteInfo{}
	row := d.mysql.QueryRow(c, _whiteContent, content)
	if err = row.Scan(&rule.ID, &rule.Content, &rule.Mode, &rule.Comment, &rule.State); err != nil {
		if err == sql.ErrNoRows {
			rule = nil
			err = nil
			return
		}
	}
	return
}

// WhiteInfo get white info
func (d *Dao) WhiteInfo(c context.Context, id int64) (rule *model.WhiteInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.mysql.Query(c, _whiteInfo, id); err != nil {
		return
	}
	defer rows.Close()
	var (
		tpIDs = make(map[int64]struct{})
		areas = make(map[string]struct{})
	)
	for rows.Next() {
		var (
			r    = &model.WhiteInfo{}
			area string
			tpID int64
		)
		if err = rows.Scan(&r.ID, &r.Content, &r.Mode, &r.Comment, &area, &tpID); err != nil {
			return
		}
		areas[area] = struct{}{}
		tpIDs[tpID] = struct{}{}
		rule = r
	}
	if rule != nil {
		for a := range areas {
			rule.Areas = append(rule.Areas, a)
		}
		for t := range tpIDs {
			rule.TpIDs = append(rule.TpIDs, t)
		}
	}
	err = rows.Err()
	return
}

// SearchWhiteContent search white
func (d *Dao) SearchWhiteContent(c context.Context, content, area string, start, offset int64) (rs []*model.WhiteInfo, err error) {
	var (
		rows     *xsql.Rows
		querySQL = _searchWhiteContent
	)
	querySQL = fmt.Sprintf(querySQL, "AND area='"+area+"'")
	if content != "" {
		querySQL = querySQL + " AND content like \"%" + content + "%\" ORDER BY ctime DESC"
	} else {
		querySQL = querySQL + " ORDER BY ctime DESC"
	}
	querySQL = querySQL + fmt.Sprintf(" LIMIT %d,%d", start, offset)
	if rows, err = d.mysql.Query(c, querySQL); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.WhiteInfo{}
		if err = rows.Scan(&r.ID, &r.Content, &r.Mode, &r.Comment); err != nil {
			return
		}
		rs = append(rs, r)
	}
	return
}

// CountWhiteContent get total white content with area size
func (d *Dao) CountWhiteContent(c context.Context, content, area string) (total int64, err error) {
	var (
		row      *xsql.Row
		querySQL = _countWhiteContent
	)
	querySQL = fmt.Sprintf(querySQL, "AND area='"+area+"'")
	if content != "" {
		querySQL = querySQL + " AND content like \"%" + content + "%\""
	}
	if row = d.mysql.QueryRow(c, querySQL); err != nil {
		return
	}
	err = row.Scan(&total)
	return
}

// UpsertWhiteContent insert or update new white content
func (d *Dao) UpsertWhiteContent(c context.Context, tx *xsql.Tx, content, comment string, mode int8) (newID int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_upsertWhiteContent, content, mode, comment, mode, comment); err != nil {
		return
	}
	return res.LastInsertId()
}

// UpdateWhiteContent update white content
func (d *Dao) UpdateWhiteContent(c context.Context, tx *xsql.Tx, mode int8, content, comment string) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateWhiteContent, mode, comment, content); err != nil {
		return
	}
	return res.RowsAffected()
}

// DeleteWhiteContent delete white content by id
func (d *Dao) DeleteWhiteContent(c context.Context, tx *xsql.Tx, id int64) (afftected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_deleteWhiteContent, id); err != nil {
		return
	}
	return res.RowsAffected()
}

// WhiteAreas get all white by contentID areas
func (d *Dao) WhiteAreas(c context.Context, contentID int64) (areas []string, err error) {
	var (
		area string
		rows *xsql.Rows
	)
	if rows, err = d.mysql.Query(c, _whiteAreas, contentID); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&area); err != nil {
			return
		}
		areas = append(areas, area)
	}
	return
}

// UpsertAreaWhite insert or update white_area
func (d *Dao) UpsertAreaWhite(c context.Context, tx *xsql.Tx, area string, tp, contentID int64) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_upsertAreaWhite, area, tp, contentID); err != nil {
		return
	}
	return res.RowsAffected()
}

// DeleteAreaWhite delete area white by contentID
func (d *Dao) DeleteAreaWhite(c context.Context, tx *xsql.Tx, contentID int64) (affected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_deleteAreaWhite, contentID); err != nil {
		return
	}
	return res.RowsAffected()
}

// WhiteLog get white log by contentID
func (d *Dao) WhiteLog(c context.Context, contentID int64) (ls []*model.Log, err error) {
	var rows *xsql.Rows
	if rows, err = d.mysql.Query(c, _whiteLog, contentID); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		l := &model.Log{}
		if err = rows.Scan(&l.AdminID, &l.Name, &l.Comment, &l.State, &l.Ctime); err != nil {
			return
		}
		ls = append(ls, l)
	}
	return
}

// InsertWhiteLog insert new white log
func (d *Dao) InsertWhiteLog(c context.Context, tx *xsql.Tx, contentID, adid int64, name, comment string, state int8) (afftected int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_insertWhiteLog, contentID, adid, name, comment, state); err != nil {
		return
	}
	return res.RowsAffected()
}
