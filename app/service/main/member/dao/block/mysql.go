package block

import (
	"context"
	"fmt"

	model "go-common/app/service/main/member/model/block"
	xsql "go-common/library/database/sql"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_user                   = `SELECT id,mid,status,ctime,mtime FROM block_user WHERE mid=? LIMIT 1`
	_upsertUser             = `INSERT INTO block_user (mid,status) VALUES (?,?) ON DUPLICATE KEY UPDATE status=?`
	_userStatusList         = `SELECT id,mid FROM block_user WHERE status=? AND id >= 0 LIMIT ?`
	_userStatusListWithMIDs = `SELECT mid FROM block_user WHERE status=? AND mid IN (%s)`

	_userLastHistory = `SELECT id,mid,action,start_time,duration FROM block_history_%d WHERE mid=? ORDER BY id DESC LIMIT 1`

	_addAddBlockCount = `INSERT INTO block_user_detail (mid,block_count) VALUES (?,1) ON DUPLICATE KEY UPDATE block_count=block_count+1`

	_insertHistory = `INSERT INTO block_history_%d (mid,admin_id,admin_name,source,area,reason,comment,action,start_time,duration,notify) VALUES (?,?,?,?,?,?,?,?,?,?,?)`

	_userDetails = `SELECT id,mid,block_count,ctime,mtime FROM block_user_detail WHERE mid IN (%s)`
)

func historyIdx(mid int64) int64 {
	return mid % 10
}

// User get block user from db
func (d *Dao) User(c context.Context, mid int64) (user *model.DBUser, err error) {
	user = &model.DBUser{}
	row := d.db.QueryRow(c, _user, mid)
	if err = row.Scan(&user.ID, &user.MID, &user.Status, &user.CTime, &user.MTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			user = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// UserLastHistory get lastest block history
func (d *Dao) UserLastHistory(c context.Context, mid int64) (his *model.DBHistory, err error) {
	var (
		sql = fmt.Sprintf(_userLastHistory, historyIdx(mid))
	)
	row := d.db.QueryRow(c, sql, mid)
	his = &model.DBHistory{}
	if err = row.Scan(&his.ID, &his.MID, &his.Action, &his.StartTime, &his.Duration); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			his = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// UserStatusList get user status list
func (d *Dao) UserStatusList(c context.Context, status model.BlockStatus, startID int64, limit int) (maxID int64, mids []int64, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _userStatusList, status, startID, limit); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id, mid int64
		)
		if err = rows.Scan(&id, &mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		if maxID < id {
			maxID = id
		}
		mids = append(mids, mid)
	}

	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UserStatusMapWithMIDs get user block status map by mids
func (d *Dao) UserStatusMapWithMIDs(c context.Context, status model.BlockStatus, mids []int64) (midMap map[int64]struct{}, err error) {
	var (
		rows *xsql.Rows
	)
	url := fmt.Sprintf(_userStatusListWithMIDs, xstr.JoinInts(mids))
	if rows, err = d.db.Query(c, url, status); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	midMap = make(map[int64]struct{})
	for rows.Next() {
		var (
			mid int64
		)
		if err = rows.Scan(&mid); err != nil {
			err = errors.WithStack(err)
			return
		}
		midMap[mid] = struct{}{}
	}

	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// TxUpdateUser is
func (d *Dao) TxUpdateUser(c context.Context, tx *xsql.Tx, mid int64, status model.BlockStatus) (err error) {
	if _, err = tx.Exec(_upsertUser, mid, status, status); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateAddBlockCount is
func (d *Dao) UpdateAddBlockCount(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, _addAddBlockCount, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// TxInsertHistory is
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

// UserDetails .
func (d *Dao) UserDetails(c context.Context, mids []int64) (midMap map[int64]*model.DBUserDetail, err error) {
	var (
		rows *xsql.Rows
	)
	sql := fmt.Sprintf(_userDetails, xstr.JoinInts(mids))
	if rows, err = d.db.Query(c, sql); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer rows.Close()
	midMap = make(map[int64]*model.DBUserDetail, len(mids))
	for rows.Next() {
		detail := &model.DBUserDetail{}
		if err = rows.Scan(&detail.ID, &detail.MID, &detail.BlockCount, &detail.CTime, &detail.MTime); err != nil {
			err = errors.WithStack(err)
			return
		}
		midMap[detail.MID] = detail
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
