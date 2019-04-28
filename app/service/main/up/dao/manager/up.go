package manager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/dao"
	"go-common/app/service/main/up/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_upsWithGroupBase          = "SELECT ups.id,mid,up_group.id as group_id ,up_group.short_tag as group_tag,up_group.name as group_name,ups.note,ups.ctime, ups.mtime, ups.uid FROM ups INNER JOIN up_group on  ups.type=up_group.id "
	_upsWithGroupBaseWithColor = "SELECT ups.id,mid,up_group.id as group_id ,up_group.short_tag as group_tag,up_group.name as group_name,ups.note,ups.ctime, ups.mtime, ups.uid, up_group.colors FROM ups INNER JOIN up_group on  ups.type=up_group.id "
	_upsWithGroup              = _upsWithGroupBaseWithColor + "limit ? offset ? "
	_upsWithGroupByMtime       = _upsWithGroupBaseWithColor + "where ups.mtime >= ?"
	_delByMid                  = "DELETE FROM ups where id = ?"
	_insertMidType             = "INSERT INTO ups (mid, type, note, ctime, mtime, uid) values "
	_selectByMidTypeWithGroup  = _upsWithGroupBase + "where mid = ? and type = ?"
	_updateByID                = "update ups set type=?, note=?, uid=?, mtime=? where id = ?"
	//MysqlTimeFormat mysql time format
	MysqlTimeFormat  = "2006-01-02 15:04:05"
	_selectByID      = _upsWithGroupBase + "where ups.id = ?"
	_selectCount     = "SELECT COUNT(0) FROM ups INNER JOIN up_group on  ups.type=up_group.id "
	_upSpecialSQL    = "SELECT type FROM ups WHERE mid = ?"
	_upsSpecialSQL   = "SELECT mid, type FROM ups WHERE mid IN (%s)"
	_upGroupsMidsSQL = "SELECT id, mid, type FROM ups WHERE id > ? AND type IN (%s) LIMIT ?"
)

// UpSpecials load all ups with group info
func (d *Dao) UpSpecials(c context.Context) (ups []*model.UpSpecial, err error) {
	log.Info("start refresh ups table")
	stmt, err := d.managerDB.Prepare(_upsWithGroup)

	if err != nil {
		log.Error("d.managerDB.Prepare error(%v)", err)
		return
	}
	defer stmt.Close()
	var offset = 0
	var limit = 1000
	var cnt = 0
	var isFirst = true
	for isFirst || cnt == limit {
		isFirst = false
		rows, err1 := stmt.Query(c, limit, offset)
		if err1 != nil {
			err = err1
			log.Error("stmt.Query error(%v)", err1)
			return
		}
		cnt = 0
		for rows.Next() {
			cnt++
			var up *model.UpSpecial
			up, err = parseUpGroupWithColor(rows)
			ups = append(ups, up)
		}
		rows.Close()
		if err != nil {
			return
		}
		log.Info("reading data from table, read count=%d", cnt)
		offset += cnt
	}

	log.Info("end refresh ups table, read count=%d", offset)
	return
}

//RefreshUpSpecialIncremental refresh cache incrementally
func (d *Dao) RefreshUpSpecialIncremental(c context.Context, lastMTime time.Time) (ups []*model.UpSpecial, err error) {
	var timeStr = lastMTime.Format(MysqlTimeFormat)
	log.Info("start refresh ups table mtime>%s", timeStr)
	stmt, err := d.managerDB.Prepare(_upsWithGroupByMtime)

	if err != nil {
		log.Error("d.managerDB.Prepare error(%v)", err)
		return
	}
	defer stmt.Close()
	var cnt = 0
	rows, err1 := stmt.Query(c, timeStr)
	if err1 != nil {
		err = err1
		log.Error("stmt.Query error(%v)", err1)
		return
	}

	cnt = 0
	for rows.Next() {
		cnt++
		var up *model.UpSpecial
		up, err = parseUpGroupWithColor(rows)
		if err != nil {
			log.Error("scan row err, %v", err)
			break
		}
		ups = append(ups, up)
	}
	rows.Close()
	if err != nil {
		return
	}
	log.Info("reading data from table, read count=%d", cnt)

	return
}

//DelSpecialByID delete special by id
func (d *Dao) DelSpecialByID(c context.Context, id int64) (res sql.Result, err error) {
	var stmt *xsql.Stmt
	stmt, err = d.managerDB.Prepare(_delByMid)
	if err != nil {
		log.Error("stmt prepare fail, error(%v), sql=%s", err, _delByMid)
		return
	}
	defer stmt.Close()

	res, err = stmt.Exec(c, id)
	if err != nil {
		log.Error("delete mid from ups fail, id=%d, err=%v", id, err)
	}
	return
}

//InsertSpecial insert special
func (d *Dao) InsertSpecial(c context.Context, special *model.UpSpecial, mids ...int64) (res sql.Result, err error) {
	var count = len(mids)
	if count == 0 {
		err = errors.New("no data need update")
		return
	}
	var insertSchema []string
	var vals []interface{}
	var nowStr = time.Now().Format(MysqlTimeFormat)
	for _, mid := range mids {
		insertSchema = append(insertSchema, "(?,?,?,?,?,?)")
		vals = append(vals, mid, special.GroupID, special.Note, nowStr, nowStr, special.UID)
	}

	var insertSQL = _insertMidType + strings.Join(insertSchema, ",")
	var stmt *xsql.Stmt
	stmt, err = d.managerDB.Prepare(insertSQL)
	if err != nil {
		log.Error("stmt prepare fail, error(%v), sql=%s", err, insertSQL)
		return
	}
	defer stmt.Close()

	res, err = stmt.Exec(c, vals...)
	if err != nil {
		log.Error("insert ups fail, err=%v, groupid=%d, note=%s, mid=%v", err, special.GroupID, special.Note, mids)
	}
	return
}

//UpdateSpecialByID update special by id
func (d *Dao) UpdateSpecialByID(c context.Context, id int64, special *model.UpSpecial) (res sql.Result, err error) {
	var stmt *xsql.Stmt
	stmt, err = d.managerDB.Prepare(_updateByID)
	if err != nil {
		log.Error("stmt prepare fail, error(%v), sql=%s", err, _updateByID)
		return
	}
	defer stmt.Close()
	var timeStr = time.Now().Format(MysqlTimeFormat)
	res, err = stmt.Exec(c, special.GroupID, special.Note, special.UID, timeStr, id)
	if err != nil {
		log.Error("update ups fail, err=%v, groupid=%d, note=%s, id=%d", err, special.GroupID, special.Note, id)
	}
	return
}

//GetSpecialByMidGroup get special by mid and group
func (d *Dao) GetSpecialByMidGroup(c context.Context, mid int64, groupID int64) (res *model.UpSpecial, err error) {
	var stmt *xsql.Stmt
	stmt, err = d.managerDB.Prepare(_selectByMidTypeWithGroup)
	if err != nil {
		log.Error("stmt prepare fail, error(%v), sql=%s", err, _selectByMidTypeWithGroup)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(c, mid, groupID)
	var note sql.NullString
	var up = model.UpSpecial{}
	switch err = row.Scan(&up.ID, &up.Mid, &up.GroupID, &up.GroupTag, &up.GroupName, &note, &up.CTime, &up.MTime, &up.UID); err {
	case sql.ErrNoRows:
		err = nil
		return
	case nil:
		up.Note = note.String
		res = &up
	default:
		log.Error("rows.Scan error(%v)", err)
		return
	}
	return
}

//GetSpecialByID get special by id
func (d *Dao) GetSpecialByID(c context.Context, id int64) (res *model.UpSpecial, err error) {
	rows, err := prepareAndQuery(c, d.managerDB, _selectByID, id)
	if err != nil {
		return
	}
	defer rows.Close()

	var note sql.NullString
	for rows.Next() {
		var up = &model.UpSpecial{}
		err = rows.Scan(&up.ID, &up.Mid, &up.GroupID, &up.GroupTag, &up.GroupName, &note, &up.CTime, &up.MTime, &up.UID)
		up.Note = note.String
		res = up
		break
	}

	return
}

//GetSepcialCount get special count
func (d *Dao) GetSepcialCount(c context.Context, conditions ...dao.Condition) (count int, err error) {
	var conditionStr, args, hasOperator = dao.ConcatCondition(conditions...)
	var where = " WHERE "
	if !hasOperator {
		where = ""
	}
	rows, err := prepareAndQuery(c, d.managerDB, _selectCount+where+conditionStr, args...)
	if err != nil {
		log.Error("get special db fail, err=%+v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&count)
	}
	return
}

func parseUpGroupWithColor(rows *xsql.Rows) (up *model.UpSpecial, err error) {
	var note sql.NullString
	var colors sql.NullString
	up = &model.UpSpecial{}
	err = rows.Scan(&up.ID, &up.Mid, &up.GroupID, &up.GroupTag, &up.GroupName, &note, &up.CTime, &up.MTime, &up.UID, &colors)
	if err != nil {
		log.Error("scan row err, %v", err)
		return
	}
	up.Note = note.String
	if colors.Valid {
		var colors = strings.Split(colors.String, "|")
		if len(colors) >= 2 {
			up.FontColor = colors[0]
			up.BgColor = colors[1]
		}
	}
	return
}

//GetSpecial get special from db
func (d *Dao) GetSpecial(c context.Context, conditions ...dao.Condition) (res []*model.UpSpecial, err error) {
	var conditionStr, args, hasOperator = dao.ConcatCondition(conditions...)
	var where = " WHERE "
	if !hasOperator {
		where = ""
	}
	rows, err := prepareAndQuery(c, d.managerDB, _upsWithGroupBaseWithColor+where+conditionStr, args...)
	if err != nil {
		log.Error("get special db fail, err=%+v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var up *model.UpSpecial
		up, err = parseUpGroupWithColor(rows)
		if err != nil {
			log.Error("scan row err, %v", err)
			break
		}
		res = append(res, up)
	}
	return
}

//GetSpecialByMid get speical by mid
func (d *Dao) GetSpecialByMid(c context.Context, mid int64) (res []*model.UpSpecial, err error) {
	var condition = dao.Condition{
		Key:      "ups.mid",
		Operator: "=",
		Value:    mid,
	}

	return d.GetSpecial(c, condition)
}

// RawUpSpecial get up special propertys
func (d *Dao) RawUpSpecial(c context.Context, mid int64) (us *upgrpc.UpSpecial, err error) {
	rows, err := d.managerDB.Query(c, _upSpecialSQL, mid)
	if err != nil {
		return
	}
	defer rows.Close()
	us = new(upgrpc.UpSpecial)
	for rows.Next() {
		var groupID int64
		if err = rows.Scan(&groupID); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			return
		}
		us.GroupIDs = append(us.GroupIDs, groupID)
	}
	return
}

// RawUpsSpecial get mult up special propertys
func (d *Dao) RawUpsSpecial(c context.Context, mids []int64) (mu map[int64]*upgrpc.UpSpecial, err error) {
	rows, err := d.managerDB.Query(c, fmt.Sprintf(_upsSpecialSQL, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	defer rows.Close()
	mu = make(map[int64]*upgrpc.UpSpecial, len(mids))
	for rows.Next() {
		var mid, groupID int64
		if err = rows.Scan(&mid, &groupID); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			return
		}
		if mu[mid] == nil {
			mu[mid] = new(upgrpc.UpSpecial)
		}
		mu[mid].GroupIDs = append(mu[mid].GroupIDs, groupID)
	}
	return
}

// UpGroupsMids get mids in one group.
func (d *Dao) UpGroupsMids(c context.Context, groupIDs []int64, lastID int64, ps int) (lid int64, gmids map[int64][]int64, err error) {
	rows, err := d.managerDB.Query(c, fmt.Sprintf(_upGroupsMidsSQL, xstr.JoinInts(groupIDs)), lastID, ps)
	if err != nil {
		return
	}
	defer rows.Close()
	gmids = make(map[int64][]int64)
	for rows.Next() {
		var id, gid, mid int64
		if err = rows.Scan(&id, &mid, &gid); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
			}
			return
		}
		if id > lid {
			lid = id
		}
		gmids[gid] = append(gmids[gid], mid)
	}
	return
}
