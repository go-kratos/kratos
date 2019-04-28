package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	_addUserBlack      = "insert into user_black_%02d (`mid`,`black_mid`) values (?,?) on duplicate key update state=0"
	_cancelUserBlack   = "update user_black_%02d set state=1 where mid=? and black_mid=?"
	_getUserPartBlacks = "select black_mid, mtime from user_black_%02d where mid=? and state=0 and mtime<=? order by mtime desc, black_mid desc limit %d"
	_getUserBlacks     = "select black_mid from user_black_%02d where mid=? and state=0"
	_isBlack           = "select black_mid from user_black_%02d where mid=? and state=0 and black_mid in (%s)"
)

// TxAddUserBlack .
func (d *Dao) TxAddUserBlack(c context.Context, tx *sql.Tx, mid, upMid int64) (num int64, err error) {
	var res xsql.Result
	// sql中含有ignore，该情况为了简化是否已存在的判断
	querySQL := fmt.Sprintf(_addUserBlack, d.getTableIndex(mid))
	if res, err = tx.Exec(querySQL, mid, upMid); err != nil {
		log.Errorv(c, log.KV("event", "add_user_black"), log.KV("err", err), log.KV("mid", mid), log.KV("up_mid", upMid), log.KV("sql", querySQL))
		return
	}
	return res.RowsAffected()
}

// TxCancelUserBlack .
func (d *Dao) TxCancelUserBlack(c context.Context, tx *sql.Tx, mid, upMid int64) (num int64, err error) {
	var res xsql.Result
	querySQL := fmt.Sprintf(_cancelUserBlack, d.getTableIndex(mid))
	if res, err = tx.Exec(querySQL, mid, upMid); err != nil {
		log.Errorv(c, log.KV("event", "cancel_user_black"), log.KV("err", err), log.KV("mid", mid), log.KV("up_mid", upMid), log.KV("sql", querySQL))
		return
	}
	return res.RowsAffected()
}

// FetchBlackList 获取mid的所有拉黑up主
func (d *Dao) FetchBlackList(c context.Context, mid int64) (upMid []int64, err error) {
	querySQL := fmt.Sprintf(_getUserBlacks, d.getTableIndex(mid))
	rows, err := d.db.Query(c, querySQL, mid)
	if err != nil {
		log.Errorv(c, log.KV("event", "mysql_select"), log.KV("table", "user_black"))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var m int64
		if err = rows.Scan(&m); err != nil {
			log.Errorv(c, log.KV("event", "mysql_scan"), log.KV("table", "user_black"))
			return
		}
		upMid = append(upMid, m)
	}
	log.Infov(c, log.KV("event", "mysql_select"), log.KV("sql", querySQL),
		log.KV("mid", mid), log.KV("black_num", len(upMid)))
	return
}

// FetchPartBlackList 获取mid的拉黑up主
func (d *Dao) FetchPartBlackList(c context.Context, mid int64, cursor model.CursorValue, size int) (
	MID2IDMap map[int64]time.Time, blackMIDs []int64, err error) {
	MID2IDMap, blackMIDs, err = d.fetchPartRelationUserList(c, mid, cursor, _getUserPartBlacks)
	return
}

// IsBlack 获取mid的拉黑用户
func (d *Dao) IsBlack(c context.Context, mid int64, candidateMIDs []int64) (MIDMap map[int64]bool) {
	if len(candidateMIDs) == 0 {
		return
	}
	MIDMap = d.isMidIn(c, mid, candidateMIDs, _isBlack)
	return
}
