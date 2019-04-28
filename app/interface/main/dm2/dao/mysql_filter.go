package dao

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_upFilterSharding       = 10
	_userFilterSharding     = 50
	_userFltCntSharding     = 10
	_addUserFilterSQL       = "INSERT INTO dm_filter_user_%02d(mid,type,filter,comment) VALUES(?,?,?,?)"
	_getUserFilterSQL       = "SELECT id,mid,type,filter,comment,ctime,mtime FROM dm_filter_user_%02d WHERE mid=? AND type=?"
	_getUserFiltersSQL      = "SELECT id,mid,type,filter,comment,ctime,mtime FROM dm_filter_user_%02d WHERE mid=?"
	_getUserFiltersByIDSQL  = "SELECT id,mid,type,filter,comment,ctime,mtime FROM dm_filter_user_%02d WHERE mid=? AND id IN(%s)"
	_delUserFilterSQL       = "DELETE FROM dm_filter_user_%02d WHERE mid=%d AND id IN(%s)"
	_maddUpFilterSQL        = "INSERT INTO dm_filter_up_%02d(mid,type,filter,active,comment) VALUES %s"
	_getUpFilterSQL         = "SELECT id,mid,type,filter,active,comment,ctime,mtime FROM dm_filter_up_%02d WHERE mid=? AND oid=0 AND type=? AND active=1"
	_getUpFiltersSQL        = "SELECT id,mid,type,filter,active,comment,ctime,mtime FROM dm_filter_up_%02d WHERE mid=? AND oid=0 AND active=1"
	_uptUpFilterSQL         = "UPDATE dm_filter_up_%02d SET active=? WHERE type=? AND mid=? AND filter IN(%s)"
	_getGlbFilterSQL        = "SELECT id,type,filter,ctime,mtime FROM dm_filter_global WHERE type=? AND filter=?"
	_getGlbFiltersSQL       = "SELECT id,type,filter,ctime,mtime FROM dm_filter_global WHERE id>=? ORDER BY id DESC LIMIT ?"
	_addGlbFilterSQL        = "INSERT INTO dm_filter_global(type,filter) VALUES(?,?)"
	_delGlbFilterSQL        = "DELETE FROM dm_filter_global WHERE id IN(%s)"
	_getUserFilterCntSQL    = "SELECT count FROM dm_filter_user_count_%02d WHERE mid=? AND type=?"
	_getUpFilterCntSQL      = "SELECT count FROM dm_filter_up_count WHERE mid=? AND type=?"
	_addUserFilterCntSQL    = "INSERT INTO dm_filter_user_count_%02d(mid,type,count) VALUES(?,?,?)"
	_updateUserFilterCntSQL = "UPDATE dm_filter_user_count_%02d SET count=count+? WHERE mid=? AND type=? AND count<?"
	_addUpFilterCntSQL      = "INSERT INTO dm_filter_up_count(mid,type,count) VALUES(?,?,?)"
	_updateUpFilterCntSQL   = "UPDATE dm_filter_up_count SET count=count+? WHERE mid=? AND type=? AND count<?"
)

func (d *Dao) hitUpFilter(mid int64) int64 {
	return mid % _upFilterSharding
}

func (d *Dao) hitUserFilter(mid int64) int64 {
	return mid % _userFilterSharding
}

func (d *Dao) hitUserFilterCnt(mid int64) int64 {
	return mid % _userFltCntSharding
}

// addslashes 函数返回在预定义字符之前添加反斜杠的字符串。
// 预定义字符是：单引号（'）或反斜杠（\）
func addSlashes(str string) string {
	var buf bytes.Buffer
	for i := 0; i < len(str); i++ {
		if str[i] == '\'' || str[i] == '\\' {
			buf.WriteByte('\\')
		}
		buf.WriteByte(str[i])
	}
	return buf.String()
}

// AddUserFilter add filter rule
func (d *Dao) AddUserFilter(tx *sql.Tx, mid int64, fType int8, filter, comment string) (lastID int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_addUserFilterSQL, d.hitUserFilter(mid)), mid, fType, filter, comment)
	if err != nil {
		log.Error("d.AddUserFilter(%d,%d,%s) error(%v)", mid, fType, filter, err)
		return
	}
	return res.LastInsertId()
}

// UserFilter select user filter by mid and type.
func (d *Dao) UserFilter(c context.Context, mid int64, fType int8) (res []*model.UserFilter, err error) {
	rows, err := d.dbDM.Query(c, fmt.Sprintf(_getUserFilterSQL, d.hitUserFilter(mid)), mid, fType)
	if err != nil {
		log.Error("dbDM.Query(mid:%d,type:%d) error(%v)", mid, fType, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &model.UserFilter{}
		if err = rows.Scan(&f.ID, &f.Mid, &f.Type, &f.Filter, &f.Comment, &f.Ctime, &f.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, f)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UserFilters return all filter.
func (d *Dao) UserFilters(c context.Context, mid int64) (res []*model.UserFilter, err error) {
	rows, err := d.dbDM.Query(c, fmt.Sprintf(_getUserFiltersSQL, d.hitUserFilter(mid)), mid)
	if err != nil {
		log.Error("dbDM.Query(mid:%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &model.UserFilter{}
		if err = rows.Scan(&f.ID, &f.Mid, &f.Type, &f.Filter, &f.Comment, &f.Ctime, &f.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, f)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UserFiltersByID return all filters specified by ids
func (d *Dao) UserFiltersByID(c context.Context, mid int64, ids []int64) (res []*model.UserFilter, err error) {
	rows, err := d.dbDM.Query(c, fmt.Sprintf(_getUserFiltersByIDSQL, d.hitUserFilter(mid), xstr.JoinInts(ids)), mid)
	if err != nil {
		log.Error("dbDM.Quer(mid:%d, ids:%v) error(%v)", mid, ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &model.UserFilter{}
		if err = rows.Scan(&f.ID, &f.Mid, &f.Type, &f.Filter, &f.Comment, &f.Ctime, &f.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, f)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// DelUserFilter batch delete filter rules by mid
func (d *Dao) DelUserFilter(tx *sql.Tx, mid int64, ids []int64) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_delUserFilterSQL, d.hitUserFilter(mid), mid, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("tx.Exec(mid:%d, ids:%v) error(%v)", mid, ids, err)
		return
	}
	return res.RowsAffected()
}

// MultiAddUpFilter add filter rule,the key of fltMap is filter content,value is comment.
// TODO add comment field in table:dm_filter_up_xx and insert comment
func (d *Dao) MultiAddUpFilter(tx *sql.Tx, mid int64, fType int8, fltMap map[string]string) (affect int64, err error) {
	var buf bytes.Buffer
	for filter, comment := range fltMap {
		buf.WriteString(fmt.Sprintf(`(%d,%d,'%s',%d,'%s'),`, mid, fType, addSlashes(filter), model.FilterActive, addSlashes(comment)))
	}
	buf.Truncate(buf.Len() - 1)
	res, err := tx.Exec(fmt.Sprintf(_maddUpFilterSQL, d.hitUpFilter(mid), buf.String()))
	if err != nil {
		log.Error("d.AddUpFilter(mid:%d,type:%d,filters:%v) error(%v)", mid, fType, fltMap, err)
		return
	}
	return res.RowsAffected()
}

// UpFilter return filter rules by mid and filter type.
func (d *Dao) UpFilter(c context.Context, mid int64, fType int8) (res []*model.UpFilter, err error) {
	res = make([]*model.UpFilter, 0)
	rows, err := d.dbDM.Query(c, fmt.Sprintf(_getUpFilterSQL, d.hitUpFilter(mid)), mid, fType)
	if err != nil {
		log.Error("dbDM.Query(mid:%d,type:%d) error(%v)", mid, fType, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &model.UpFilter{}
		if err = rows.Scan(&f.ID, &f.Mid, &f.Type, &f.Filter, &f.Active, &f.Comment, &f.Ctime, &f.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, f)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UpFilters return all filter rules
func (d *Dao) UpFilters(c context.Context, mid int64) (res []*model.UpFilter, err error) {
	rows, err := d.dbDM.Query(c, fmt.Sprintf(_getUpFiltersSQL, d.hitUpFilter(mid)), mid)
	if err != nil {
		log.Error("dbDM.Query(mid:%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &model.UpFilter{}
		if err = rows.Scan(&f.ID, &f.Mid, &f.Type, &f.Filter, &f.Active, &f.Comment, &f.Ctime, &f.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, f)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UpdateUpFilter batch edit filter.
func (d *Dao) UpdateUpFilter(tx *sql.Tx, mid int64, fType, active int8, filters []string) (affect int64, err error) {
	sli := make([]string, len(filters))
	for i, ss := range filters {
		sli[i] = "'" + addSlashes(ss) + "'"
	}
	filter := strings.Join(sli, ",")
	res, err := tx.Exec(fmt.Sprintf(_uptUpFilterSQL, d.hitUpFilter(mid), filter), active, fType, mid)
	if err != nil {
		log.Error("tx.Exec(mid:%d,filters:%v) error(%v)", mid, filters, err)
		return
	}
	return res.RowsAffected()
}

// AddGlobalFilter add filter rule
func (d *Dao) AddGlobalFilter(c context.Context, fType int8, filter string) (lastID int64, err error) {
	res, err := d.dbDM.Exec(c, _addGlbFilterSQL, fType, filter)
	if err != nil {
		log.Error("dbDM.Exec(%d,%s) error(%v)", fType, filter, err)
		return
	}
	return res.LastInsertId()
}

// GlobalFilter select global filters by type and filter.
func (d *Dao) GlobalFilter(c context.Context, fType int8, filter string) (res []*model.GlobalFilter, err error) {
	rows, err := d.dbDM.Query(c, _getGlbFilterSQL, fType, filter)
	if err != nil {
		log.Error("d.dbDM(type:%d, filter:%s) error(%v)", fType, filter, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &model.GlobalFilter{}
		if err = rows.Scan(&f.ID, &f.Type, &f.Filter, &f.Ctime, &f.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, f)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// GlobalFilters return all filter rules
func (d *Dao) GlobalFilters(c context.Context, sid, limit int64) (res []*model.GlobalFilter, err error) {
	var rows *sql.Rows
	if rows, err = d.dbDM.Query(c, _getGlbFiltersSQL, sid, limit); err != nil {
		log.Error("dbDM.Query(start id:%d, limit:%d) error(%v)", sid, limit, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.GlobalFilter{}
		if err = rows.Scan(&r.ID, &r.Type, &r.Filter, &r.Ctime, &r.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// DelGlobalFilters batch delete filter rules
func (d *Dao) DelGlobalFilters(c context.Context, ids []int64) (affect int64, err error) {
	res, err := d.dbDM.Exec(c, fmt.Sprintf(_delGlbFilterSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("dbDM.Exec(ids:%v) error(%v)", ids, err)
		return
	}
	return res.RowsAffected()
}

// UserFilterCnt get count by mid and type
func (d *Dao) UserFilterCnt(c context.Context, tx *sql.Tx, mid int64, tp int8) (count int, err error) {
	row := tx.QueryRow(fmt.Sprintf(_getUserFilterCntSQL, d.hitUserFilterCnt(mid)), mid, tp)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			count = model.FilterNotExist
			err = nil
		} else {
			log.Error("rows.Scan() error(%v)", err)
		}
	}
	return
}

// InsertUserFilterCnt add a new row
func (d *Dao) InsertUserFilterCnt(c context.Context, tx *sql.Tx, mid int64, tp int8, count int) (id int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_addUserFilterCntSQL, d.hitUserFilterCnt(mid)), mid, tp, count)
	if err != nil {
		log.Error("d.InsertUserFilterCnt(mid:%d, type:%d, count:%d) error(%v)", mid, tp, count, err)
		return
	}
	return res.LastInsertId()
}

// UpdateUserFilterCnt set count
func (d *Dao) UpdateUserFilterCnt(c context.Context, tx *sql.Tx, mid int64, tp int8, count, limit int64) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateUserFilterCntSQL, d.hitUserFilterCnt(mid)), count, mid, tp, limit)
	if err != nil {
		log.Error("d.UpdateUserFilterCnt(mid:%d, type:%d, count:%d) error(%v)", mid, tp, count, err)
		return
	}
	return res.RowsAffected()
}

// UpFilterCnt get count by mid and type
func (d *Dao) UpFilterCnt(c context.Context, tx *sql.Tx, mid int64, tp int8) (count int, err error) {
	row := tx.QueryRow(_getUpFilterCntSQL, mid, tp)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			count = model.FilterNotExist
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// InsertUpFilterCnt insert up rule count.
func (d *Dao) InsertUpFilterCnt(c context.Context, tx *sql.Tx, mid int64, tp int8, count int) (id int64, err error) {
	res, err := tx.Exec(_addUpFilterCntSQL, mid, tp, count)
	if err != nil {
		log.Error("d.InsertUpFilterCnt(mid:%d, type:%d, count:%d) error(%v)", mid, tp, count, err)
		return
	}
	return res.LastInsertId()
}

// UpdateUpFilterCnt set count
func (d *Dao) UpdateUpFilterCnt(c context.Context, tx *sql.Tx, mid int64, tp int8, count, limit int) (affect int64, err error) {
	res, err := tx.Exec(_updateUpFilterCntSQL, count, mid, tp, limit)
	if err != nil {
		log.Error("d.UpdateUpFilterCnt(mid:%d, type:%d, count:%d) error(%v)", mid, tp, count, err)
		return
	}
	return res.RowsAffected()
}
