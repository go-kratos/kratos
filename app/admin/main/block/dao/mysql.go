package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/block/model"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_user       = `SELECT id,mid,status,ctime,mtime FROM block_user WHERE mid=? LIMIT 1`
	_users      = `SELECT id,mid,status,ctime,mtime FROM block_user WHERE mid IN (%s)`
	_upsertUser = `INSERT INTO block_user (mid,status) VALUES (?,?) ON DUPLICATE KEY UPDATE status=?`

	_userDetails      = `SELECT id,mid,block_count,ctime,mtime FROM block_user_detail WHERE mid IN (%s)`
	_addAddBlockCount = `INSERT INTO block_user_detail (mid,block_count) VALUES (?,1) ON DUPLICATE KEY UPDATE block_count=block_count+1`

	_history       = `SELECT id,mid,admin_id,admin_name,source,area,reason,comment,action,start_time,duration,notify,ctime,mtime FROM block_history_%d WHERE mid=? LIMIT ?,?`
	_historyCount  = `SELECT count(*) FROM block_history_%d WHERE mid=? LIMIT 1`
	_insertHistory = `INSERT INTO block_history_%d (mid,admin_id,admin_name,source,area,reason,comment,action,start_time,duration,notify) VALUES (?,?,?,?,?,?,?,?,?,?,?)`
)

func historyIdx(mid int64) int64 {
	return mid % 10
}

// User .
func (d *Dao) User(c context.Context, mid int64) (user *model.DBUser, err error) {
	user = &model.DBUser{}
	row := d.db.QueryRow(c, _user, mid)
	if err = row.Scan(&user.ID, &user.MID, &user.Status, &user.CTime, &user.MTime); err != nil {
		err = nil
		user = nil
		return
	}
	return
}

// Users .
func (d *Dao) Users(c context.Context, mids []int64) (users []*model.DBUser, err error) {
	if len(mids) == 0 {
		return
	}
	var (
		sql  = fmt.Sprintf(_users, strings.Join(intsToStrs(mids), ","))
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, sql); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := &model.DBUser{}
		if err = rows.Scan(&user.ID, &user.MID, &user.Status, &user.CTime, &user.MTime); err != nil {
			err = errors.WithStack(err)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// TxUpdateUser .
func (d *Dao) TxUpdateUser(c context.Context, tx *xsql.Tx, mid int64, status model.BlockStatus) (err error) {
	if _, err = tx.Exec(_upsertUser, mid, status, status); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UserDetails .
func (d *Dao) UserDetails(c context.Context, mids []int64) (users []*model.DBUserDetail, err error) {
	if len(mids) == 0 {
		return
	}
	var (
		sql  = fmt.Sprintf(_userDetails, strings.Join(intsToStrs(mids), ","))
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, sql); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := &model.DBUserDetail{}
		if err = rows.Scan(&user.ID, &user.MID, &user.BlockCount, &user.CTime, &user.MTime); err != nil {
			err = errors.WithStack(err)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// UpdateAddBlockCount .
func (d *Dao) UpdateAddBlockCount(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, _addAddBlockCount, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// History 获得mid历史封禁记录
func (d *Dao) History(c context.Context, mid int64, start, limit int) (history []*model.DBHistory, err error) {
	var (
		rows *xsql.Rows
		sql  = fmt.Sprintf(_history, historyIdx(mid))
	)
	if rows, err = d.db.Query(c, sql, mid, start, limit); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		h := &model.DBHistory{}
		if err = rows.Scan(&h.ID, &h.MID, &h.AdminID, &h.AdminName, &h.Source, &h.Area, &h.Reason, &h.Comment, &h.Action, &h.StartTime, &h.Duration, &h.Notify, &h.CTime, &h.MTime); err != nil {
			err = errors.WithStack(err)
			return
		}
		history = append(history, h)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

// HistoryCount 获得历史记录总长度
func (d *Dao) HistoryCount(c context.Context, mid int64) (total int, err error) {
	var (
		row *xsql.Row
		sql = fmt.Sprintf(_historyCount, historyIdx(mid))
	)
	row = d.db.QueryRow(c, sql, mid)
	if err = row.Scan(&total); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// TxInsertHistory .
func (d *Dao) TxInsertHistory(c context.Context, tx *xsql.Tx, h *model.DBHistory) (err error) {
	var (
		sql = fmt.Sprintf(_insertHistory, historyIdx(h.MID))
	)
	if _, err = tx.Exec(sql, h.MID, h.AdminID, h.AdminName, h.Source, h.Area, h.Reason, h.Comment, h.Action, h.StartTime, h.Duration, h.Notify); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func intsToStrs(ints []int64) (strs []string) {
	for _, i := range ints {
		strs = append(strs, fmt.Sprintf("%d", i))
	}
	return
}
