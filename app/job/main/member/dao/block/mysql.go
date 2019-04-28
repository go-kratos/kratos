package block

import (
	"context"
	"fmt"

	model "go-common/app/job/main/member/model/block"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_upsertUser     = `INSERT INTO block_user (mid,status) VALUES (?,?) ON DUPLICATE KEY UPDATE status=?`
	_userStatusList = `SELECT id,mid FROM block_user WHERE status=? AND id>? LIMIT ?`
	_userExtra      = `SELECT id,mid,credit_answer_flag,action_time FROM block_extra WHERE mid=? ORDER BY id DESC LIMIT 1`

	_userLastHistory = `SELECT id,mid,source,area,action,start_time,duration FROM block_history_%d WHERE mid=? ORDER BY id DESC LIMIT 1`

	_upsertAddBlockCount = `INSERT INTO block_user_detail (mid,block_count) VALUES (?,1) ON DUPLICATE KEY UPDATE block_count=block_count+1`

	_upsertExtra   = `INSERT INTO block_extra (mid,credit_answer_flag,action_time) VALUES (?,?,?) ON DUPLICATE KEY UPDATE credit_answer_flag=? , action_time=?`
	_insertHistory = `INSERT INTO block_history_%d (mid,admin_id,admin_name,source,area,reason,comment,action,start_time,duration,notify) VALUES (?,?,?,?,?,?,?,?,?,?,?)`
)

func historyIdx(mid int64) int64 {
	return mid % 10
}

// UserStatusList is.
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

// UserLastHistory is.
func (d *Dao) UserLastHistory(c context.Context, mid int64) (his *model.DBHistory, err error) {
	var (
		sql = fmt.Sprintf(_userLastHistory, historyIdx(mid))
	)
	row := d.db.QueryRow(c, sql, mid)
	his = &model.DBHistory{}
	if err = row.Scan(&his.ID, &his.MID, &his.Source, &his.Area, &his.Action, &his.StartTime, &his.Duration); err != nil {
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

// UserExtra is.
func (d *Dao) UserExtra(c context.Context, mid int64) (ex *model.DBExtra, err error) {
	row := d.db.QueryRow(c, _userExtra, mid)
	ex = &model.DBExtra{}
	if err = row.Scan(&ex.ID, &ex.MID, &ex.CreditAnswerFlag, &ex.ActionTime); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			ex = nil
			return
		}
		err = errors.WithStack(err)
		return
	}
	return
}

// TxUpsertUser is.
func (d *Dao) TxUpsertUser(c context.Context, tx *xsql.Tx, mid int64, status model.BlockStatus) (count int64, err error) {
	rows, err := tx.Exec(_upsertUser, mid, status, status)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return rows.RowsAffected()
}

// InsertExtra is.
func (d *Dao) InsertExtra(c context.Context, ex *model.DBExtra) (err error) {
	if _, err = d.db.Exec(c, _upsertExtra, ex.MID, ex.CreditAnswerFlag, ex.ActionTime, ex.CreditAnswerFlag, ex.ActionTime); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// TxUpsertExtra is.
func (d *Dao) TxUpsertExtra(c context.Context, tx *xsql.Tx, ex *model.DBExtra) (err error) {
	if _, err = tx.Exec(_upsertExtra, ex.MID, ex.CreditAnswerFlag, ex.ActionTime, ex.CreditAnswerFlag, ex.ActionTime); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//TxInsertHistory is.
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

//UpsertAddBlockCount is.
func (d *Dao) UpsertAddBlockCount(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, _upsertAddBlockCount, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
