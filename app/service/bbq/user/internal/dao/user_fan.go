package dao

import (
	"context"
	xsql "database/sql"
	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/database/sql"
	"go-common/library/time"
)

const (
	_addUserFan      = "insert into %s (`mid`,`fan_mid`) values (?,?) on duplicate key update state=0"
	_cancelUserFan   = "update %s set state=1 where mid = ? and fan_mid = ?"
	_getUserPartFans = "select fan_mid, mtime from user_fan_%02d where mid=? and state=0 and mtime<=? order by mtime desc, fan_mid desc limit %d"
	_isFan           = "select fan_mid from user_fan_%02d where mid = ? and state=0 and fan_mid in (%s)"
)

//TxAddUserFan .
func (d *Dao) TxAddUserFan(tx *sql.Tx, mid, fanMid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(d.userFanSQL(mid, _addUserFan), mid, fanMid); err != nil {
		return
	}
	return res.LastInsertId()
}

//TxCancelUserFan .
func (d *Dao) TxCancelUserFan(tx *sql.Tx, mid, fanMid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(d.userFanSQL(mid, _cancelUserFan), mid, fanMid); err != nil {
		return
	}
	return res.RowsAffected()
}

// IsFan 获取mid的粉丝
func (d *Dao) IsFan(c context.Context, mid int64, candidateMIDs []int64) (MIDMap map[int64]bool) {
	if len(candidateMIDs) == 0 {
		return
	}
	MIDMap = d.isMidIn(c, mid, candidateMIDs, _isFan)
	return
}

// FetchPartFanList 获取mid的粉丝列表
func (d *Dao) FetchPartFanList(c context.Context, mid int64, cursor model.CursorValue, size int) (
	MID2IDMap map[int64]time.Time, followedMIDs []int64, err error) {
	MID2IDMap, followedMIDs, err = d.fetchPartRelationUserList(c, mid, cursor, _getUserPartFans)
	return
}
