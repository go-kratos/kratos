package dao

import (
	"context"
	xsql "database/sql"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/time"

	"go-common/app/service/bbq/user/internal/model"
)

const (
	_addUserFollow      = "insert into %s (`mid`,`followed_mid`)values(?,?) on duplicate key update state=0"
	_cancelUserFollow   = "update %s set state=1 where mid=? and followed_mid=?"
	_getUserFollows     = "select followed_mid from %s where mid = ? and state=0"
	_getUserPartFollows = "select followed_mid, mtime from user_follow_%02d where mid=? and state=0 and mtime<=? order by mtime desc, followed_mid desc limit %d"
	_isFollow           = "select followed_mid from user_follow_%02d where mid=? and state=0 and followed_mid in (%s)"
)

//TxAddUserFollow .
func (d *Dao) TxAddUserFollow(c context.Context, tx *sql.Tx, mid, upMid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(d.userFollowSQL(mid, _addUserFollow), mid, upMid); err != nil {
		log.Errorv(c, log.KV("event", "user_follow"), log.KV("err", err), log.KV("mid", mid), log.KV("up_mid", upMid))
		return
	}
	return res.RowsAffected()
}

//TxCancelUserFollow .
func (d *Dao) TxCancelUserFollow(c context.Context, tx *sql.Tx, mid, upMid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(d.userFollowSQL(mid, _cancelUserFollow), mid, upMid); err != nil {
		log.Errorv(c, log.KV("event", "cancel_user_follow"), log.KV("err", err), log.KV("mid", mid), log.KV("up_mid", upMid))
		return
	}
	return res.RowsAffected()
}

// FetchFollowList 获取mid的所有关注up主
func (d *Dao) FetchFollowList(c context.Context, mid int64) (upMid []int64, err error) {
	querySQL := d.userFollowSQL(mid, _getUserFollows)
	rows, err := d.db.Query(c, querySQL, mid)
	if err != nil {
		log.Errorv(c, log.KV("event", "mysql_select"), log.KV("table", "user_follow"))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var m int64
		if err = rows.Scan(&m); err != nil {
			log.Errorv(c, log.KV("event", "mysql_scan"), log.KV("table", "user_follow"))
			return
		}
		upMid = append(upMid, m)
	}
	log.Infov(c, log.KV("event", "mysql_select"), log.KV("table", "user_follow"),
		log.KV("mid", mid), log.KV("follow_num", len(upMid)))
	return
}

// FetchPartFollowList 获取mid的关注up主
// 		cursor_id代表相应的up_mid，cursor_time表示mtime
func (d *Dao) FetchPartFollowList(c context.Context, mid int64, cursor model.CursorValue, size int) (
	MID2IDMap map[int64]time.Time, followedMIDs []int64, err error) {
	MID2IDMap, followedMIDs, err = d.fetchPartRelationUserList(c, mid, cursor, _getUserPartFollows)
	return
}

// IsFollow 获取mid的关注up主
func (d *Dao) IsFollow(c context.Context, mid int64, candidateMIDs []int64) (MIDMap map[int64]bool) {
	if len(candidateMIDs) == 0 {
		return
	}
	MIDMap = d.isMidIn(c, mid, candidateMIDs, _isFollow)
	return
}
