package dao

import (
	"context"
	xsql "database/sql"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/user/api"
	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"
)

//常量
const (
	_addUserLike    = "insert into %s (`mid`, `opid`) values (?,?) on duplicate key update state=0"
	_cancelUserLike = "update %s set state = 1 where mid=? and opid=?"
	_selectUserLike = "select opid from user_like_%02d where mid = ? and state = 0 and opid in (%s)"
	_spaceUserLike  = "select opid, mtime from user_like_%02d where mid=%d and state=0 and mtime %s ? order by mtime %s, opid %s limit %d"
)

//TxAddUserLike .
func (d *Dao) TxAddUserLike(tx *sql.Tx, mid, svid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(d.userLikeSQL(mid, _addUserLike), mid, svid); err != nil {
		return
	}
	return res.RowsAffected()
}

//TxCancelUserLike .
func (d *Dao) TxCancelUserLike(tx *sql.Tx, mid, svid int64) (num int64, err error) {
	var res xsql.Result
	if res, err = tx.Exec(d.userLikeSQL(mid, _cancelUserLike), mid, svid); err != nil {
		return
	}
	return res.RowsAffected()
}

// CheckUserLike 检测用户是否点赞
func (d *Dao) CheckUserLike(c context.Context, mid int64, svids []int64) (res []int64, err error) {
	log.V(1).Info("user like mid(%d) svids(%v)", mid, svids)
	ls := len(svids)
	if ls == 0 || mid == 0 {
		return
	}
	idStr := xstr.JoinInts(svids)
	querySQL := fmt.Sprintf(_selectUserLike, d.getTableIndex(mid), idStr)
	rows, err := d.db.Query(c, querySQL, mid)
	log.V(1).Infov(c, log.KV("log", fmt.Sprintf("user like rows(%v) err(%v)", rows, err)))
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	for rows.Next() {
		opid := int64(0)
		rows.Scan(&opid)
		res = append(res, opid)
	}
	log.V(1).Infov(c, log.KV("log", fmt.Sprintf("user like res(%v)", res)))
	return
}

// GetUserLikeList 返回用户点赞列表
// 		当前cursorID表示opid
func (d *Dao) GetUserLikeList(c context.Context, mid int64, cursorNext bool, cursor model.CursorValue, size int) (
	likeSvs []*api.LikeSv, err error) {

	compareSymbol := string(">=")
	orderDirection := "asc"
	if cursorNext {
		compareSymbol = "<="
		orderDirection = "desc"
	}
	querySQL := fmt.Sprintf(_spaceUserLike, d.getTableIndex(mid), mid, compareSymbol, orderDirection, orderDirection, size)
	log.V(1).Infov(c, log.KV("like_list_sql", querySQL))
	rows, err := d.db.Query(c, querySQL, cursor.CursorTime)
	if err != nil {
		log.Errorv(c, log.KV("event", "mysql_select"), log.KV("table", "user_like"),
			log.KV("mid", mid), log.KV("sql", querySQL))
		return
	}
	defer rows.Close()
	var svID int64
	var curMtime time.Time
	conflict := bool(true)
	for rows.Next() {
		if err = rows.Scan(&svID, &curMtime); err != nil {
			log.Errorv(c, log.KV("event", "mysql_scan"), log.KV("table", "user_like"),
				log.KV("sql", querySQL))
			return
		}
		// 为了解决同一个mtime的冲突问题
		if curMtime == cursor.CursorTime && conflict {
			if svID == cursor.CursorID {
				conflict = false
			}
			continue
		}
		cursorValue := model.CursorValue{CursorID: svID, CursorTime: curMtime}
		jsonStr, _ := json.Marshal(cursorValue) // marshal的时候相信库函数，不做err判断
		likeSvs = append(likeSvs, &api.LikeSv{Svid: svID, CursorValue: string(jsonStr)})
	}
	log.Infov(c, log.KV("event", "mysql_select"), log.KV("table", "user_like"),
		log.KV("mid", mid), log.KV("id", cursor.CursorID), log.KV("size", len(likeSvs)))
	return
}
