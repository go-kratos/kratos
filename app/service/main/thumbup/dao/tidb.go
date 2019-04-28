package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"sync"
	"time"

	"go-common/app/service/main/thumbup/model"
	sql "go-common/library/database/tidb"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"go-common/library/sync/errgroup"
)

const (
	_tidbBulkSize          = 50
	_tidbBusinessesSQL     = "SELECT id, name, message_list_type, user_likes_limit, message_likes_limit, enable_originid, user_list_type FROM business WHERE dtime = '0000-00-00 00:00:00'"
	_tidbLikeMidSQL        = "SELECT type FROM likes WHERE mid=? AND business_id = ? AND origin_id = ? AND message_id = ?"
	_tidbUserLikeListSQL   = "SELECT message_id, mtime FROM likes WHERE mid = ? AND business_id =? AND type = ? ORDER BY mtime desc LIMIT ? OFFSET ?"
	_tidbUserHasLikeSQL    = "SELECT message_id, mtime FROM likes WHERE mid = ? AND business_id =? AND type = ? AND message_id in (%s)"
	_tidbItemLikeListSQL   = "SELECT mid, mtime FROM likes WHERE business_id =? AND origin_id =? AND message_id =? AND type = ? ORDER BY mtime desc LIMIT ? OFFSET ?"
	_tidbStatSQL           = "SELECT message_id, likes_count, dislikes_count, origin_id, likes_change, dislikes_change FROM counts WHERE business_id = ? AND origin_id = ? AND message_id = ?"
	_tidbStatsSQL          = "SELECT message_id, likes_count, dislikes_count, likes_change, dislikes_change FROM counts WHERE business_id = ? AND origin_id = 0 AND message_id in (%s)"
	_tidbStatsOriginSQL    = "SELECT message_id, likes_count, dislikes_count, likes_change, dislikes_change FROM counts WHERE business_id = ? AND origin_id = ?"
	_tidbUserLikeCountSQL  = "SELECT COUNT(*) FROM likes WHERE business_id =? AND mid = ? AND type = ?"
	_tidbUpdateLikeSQL     = "INSERT INTO likes (business_id, origin_id, message_id, mid, type, mtime) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE type=?, mtime = ?"
	_tidbUpdateCountsSQL   = "INSERT INTO counts (business_id, origin_id, message_id, likes_count, dislikes_count, up_mid) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	_tidbupdateCountChange = "UPDATE counts SET likes_change = likes_change + ?, dislikes_change = dislikes_change + ? WHERE business_id = ? AND origin_id = ? AND message_id = ?"
	_tidbItemHasLikeSQL    = "SELECT mid, mtime FROM likes WHERE business_id = ? AND origin_id = ? AND message_id = ? AND type =? AND mid in (%s)"
)

// Businesses get business list
func (d *Dao) Businesses(c context.Context) (res []*model.Business, err error) {
	var rows *sql.Rows
	if rows, err = d.businessesStmt.Query(c); err != nil {
		PromError("tidb:业务查询")
		log.Error("tidb.businessesStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.Business{}
		if err = rows.Scan(&b.ID, &b.Name, &b.MessageListType, &b.UserLikesLimit, &b.MessageLikesLimit, &b.EnableOriginID, &b.UserListType); err != nil {
			PromError("tidb:业务Scan")
			log.Error("tidb.rows.Business.Scan error(%v)", err)
			return
		}
		res = append(res, b)
	}
	err = rows.Err()
	return
}

// LikeState get like state
func (d *Dao) LikeState(c context.Context, mid, businessID, originID, messageID int64) (res int8, err error) {
	if err = d.likeStateStmt.QueryRow(c, mid, businessID, originID, messageID).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		PromError("tidb:LikeByMid")
		log.Error("tidbLikeByMid(%d,%d,%d,%d) error(%v)", mid, businessID, originID, messageID, err)
	}
	return
}

// UserLikeCount count
func (d *Dao) UserLikeCount(c context.Context, businessID, mid int64, typ int8) (res int, err error) {
	if err = d.userLikeCountStmt.QueryRow(c, businessID, mid, typ).Scan(&res); err != nil {
		PromError("tidb:UserLikeCount")
		log.Error("tidbUserLikeCount(%d,%d,%d) error(%v)", businessID, mid, typ, err)
	}
	return
}

// RawItemLikeList item like list
func (d *Dao) RawItemLikeList(c context.Context, messageID, businessID, originID int64, state int8, start, end int) (res []*model.UserLikeRecord, err error) {
	limit := end - start + 1
	var rows *sql.Rows
	if rows, err = d.itemLikeListStmt.Query(c, businessID, originID, messageID, state, limit, start); err != nil {
		PromError("tidb:项目点赞列表")
		log.Error("tidb.ItemLikeList.Query error(%v,%v,%v,%v,%v, %v)", businessID, originID, messageID, state, limit, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.UserLikeRecord{}
		var t time.Time
		if err = rows.Scan(&b.Mid, &t); err != nil {
			PromError("tidb:业务Scan")
			log.Error("tidb.rows.Business.Scan error(%v)", err)
			return
		}
		b.Time = xtime.Time(t.Unix())
		res = append(res, b)
	}
	err = rows.Err()
	return
}

// RawUserLikeList .
func (d *Dao) RawUserLikeList(c context.Context, mid, businessID int64, state int8, start, end int) (res []*model.ItemLikeRecord, err error) {
	limit := end - start
	var rows *sql.Rows
	if rows, err = d.userLikeListStmt.Query(c, mid, businessID, state, limit, start); err != nil {
		PromError("tidb:用户点赞列表")
		log.Error("tidb.UserLikeList.Query error(%v,%v,%v,%v,%v)", mid, businessID, state, limit, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.ItemLikeRecord{}
		var t time.Time
		if err = rows.Scan(&b.MessageID, &t); err != nil {
			PromError("tidb:用户Scan")
			log.Error("tidb.rows.Business.Scan error(%v)", err)
			return
		}
		b.Time = xtime.Time(t.Unix())
		res = append(res, b)
	}
	err = rows.Err()
	return
}

// MessageStats .
func (d *Dao) MessageStats(c context.Context, businessID int64, ids []int64) (res map[int64]*model.Stats, err error) {
	res = make(map[int64]*model.Stats)
	var (
		group = &errgroup.Group{}
		mutex = &sync.Mutex{}
	)
	keysLen := len(ids)
	for i := 0; i < keysLen; i += _tidbBulkSize {
		var keys []int64
		if (i + _tidbBulkSize) > keysLen {
			keys = ids[i:]
		} else {
			keys = ids[i : i+_tidbBulkSize]
		}
		group.Go(func() error {
			statsSQL := fmt.Sprintf(_tidbStatsSQL, xstr.JoinInts(keys))
			rows, e := d.tidb.Query(c, statsSQL, businessID)
			if e != nil {
				err = e
				return nil
			}
			defer rows.Close()
			for rows.Next() {
				s := &model.Stats{}
				var likeChange, dislikeChange int64
				e = rows.Scan(&s.ID, &s.Likes, &s.Dislikes, &likeChange, &dislikeChange)
				s.Likes += likeChange
				s.Dislikes += dislikeChange
				if e != nil {
					err = e
					continue
				}
				mutex.Lock()
				res[s.ID] = s
				mutex.Unlock()
			}
			return nil
		})
	}
	group.Wait()
	if err != nil {
		PromError("tidb stats Scan")
		log.Error("tidb.rows.Stats.Scan error(%v)", err)
	}
	if len(res) == 0 {
		res = nil
	}
	return
}

// OriginStats .
func (d *Dao) OriginStats(c context.Context, businessID, originID int64) (res map[int64]*model.Stats, err error) {
	var rows *sql.Rows
	if rows, err = d.statsOriginStmt.Query(c, businessID, originID); err != nil {
		PromError("tidb:计数表")
		log.Error("db.Stats.Query error(%v,%v,%v)", businessID, originID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.Stats{OriginID: originID}
		var likeChange, dislikeChange int64
		if err = rows.Scan(&b.ID, &b.Likes, &b.Dislikes, &likeChange, &dislikeChange); err != nil {
			PromError("tidb:计数Scan")
			log.Error("tidb.rows.Stats.Scan error(%v)", err)
			return
		}
		b.Likes += likeChange
		b.Dislikes += dislikeChange
		if res == nil {
			res = make(map[int64]*model.Stats)
		}
		res[b.ID] = b
	}
	err = rows.Err()
	return
}

// Stat .
func (d *Dao) Stat(c context.Context, businessID, originID, messageID int64) (res model.Stats, err error) {
	var likeChange, dislikeChange int64
	if err = d.statStmt.QueryRow(c, businessID, originID, messageID).Scan(&res.ID, &res.Likes, &res.Dislikes, &res.OriginID, &likeChange, &dislikeChange); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res.ID = messageID
			res.OriginID = originID
			return
		}
		PromError("tidb:查询计数表")
		log.Error("tidb.Stat.Query error(%v,%v,%v,%v)", businessID, originID, messageID, err)
	}
	res.Likes += likeChange
	res.Dislikes += dislikeChange
	return
}

// RawStats stat changes
func (d *Dao) RawStats(c context.Context, businessID, originID, messageID int64) (res model.RawStats, err error) {
	if err = d.statStmt.QueryRow(c, businessID, originID, messageID).Scan(&res.ID, &res.Likes, &res.Dislikes, &res.OriginID, &res.LikesChange, &res.DislikesChange); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res.ID = messageID
			res.OriginID = originID
			return
		}
		PromError("tidb:查询change计数表")
		log.Error("tidb.StatChanges.Query error(%v,%v,%v,%v)", businessID, originID, messageID, err)
	}
	return
}

// TxUpdateCounts update like info via transaction.
func (d *Dao) tidbTxUpdateCounts(c context.Context, tx *sql.Tx, businessID, originID, messageID int64, likesCount, dislikesCount, upMid int64) (err error) {
	likeSQL := _tidbUpdateCountsSQL
	if (likesCount == 0) && (dislikesCount == 0) {
		return
	}
	if (likesCount != 0) && (dislikesCount == 0) {
		if likesCount > 0 {
			likeSQL += fmt.Sprintf("likes_count = likes_count + %v", likesCount)
		} else {
			likeSQL += fmt.Sprintf("likes_count = likes_count %v", likesCount)
		}
	} else if (likesCount == 0) && (dislikesCount != 0) {
		if dislikesCount > 0 {
			likeSQL += fmt.Sprintf("dislikes_count = dislikes_count + %v", dislikesCount)
		} else {
			likeSQL += fmt.Sprintf("dislikes_count = dislikes_count %v", dislikesCount)
		}
	} else {
		if likesCount > 0 {
			likeSQL += fmt.Sprintf("likes_count = likes_count + %v", likesCount)
		} else {
			likeSQL += fmt.Sprintf("likes_count = likes_count %v", likesCount)
		}
		if dislikesCount > 0 {
			likeSQL += fmt.Sprintf(", dislikes_count = dislikes_count + %v", dislikesCount)
		} else {
			likeSQL += fmt.Sprintf(", dislikes_count = dislikes_count %v", dislikesCount)
		}
	}
	if upMid > 0 {
		likeSQL += fmt.Sprintf(", up_mid = %d", upMid)
	}
	if _, err = tx.Exec(likeSQL, businessID, originID, messageID, likesCount, dislikesCount, upMid); err != nil {
		PromError("tidb:更新计数表")
		log.Error("dao.tidbTxUpdateCounts tx.Exec(%s, %v,%v,%v) error(%v)", likeSQL, businessID, originID, messageID, err)
	}
	return
}

// tidbTxStat .
func (d *Dao) tidbTxStat(c context.Context, tx *sql.Tx, businessID, originID, messageID int64) (res model.Stats, err error) {
	var likeChange, dislikeChange int64
	if err = tx.Stmts(d.statStmt).QueryRow(c, businessID, originID, messageID).Scan(&res.ID, &res.Likes, &res.Dislikes, &res.OriginID, &likeChange, &dislikeChange); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res.ID = messageID
			res.OriginID = originID
			return
		}
		PromError("tidb:查询计数表")
		log.Error("tidb: Tx db.Stat.Query error(%v,%v,%v,%v)", businessID, originID, messageID, err)
	}
	res.Likes += likeChange
	res.Dislikes += dislikeChange
	return
}

// UpdateCount .
func (d *Dao) UpdateCount(c context.Context, businessID, originID, messageID int64, likeChange, dislikeChange int64) (err error) {
	if _, err = d.updateCountChangeStmt.Exec(c, likeChange, dislikeChange, businessID, originID, messageID); err != nil {
		PromError("tidb:更新change计数表")
		log.Error("db.tidbUpdateCountChange.Exec error(%v,%v,%v,%v)", businessID, originID, messageID, err)
	}
	return
}

// tidbBeginTran begin transaction.
func (d *Dao) tidbBeginTran(c context.Context) (*sql.Tx, error) {
	return d.tidb.Begin(c)
}

// UpdateCounts .
func (d *Dao) UpdateCounts(c context.Context, businessID, originID, messageID int64, likeCount, dislikeCount int64, upMid int64) (err error) {
	var tx *sql.Tx
	if tx, err = d.tidbBeginTran(c); err != nil {
		log.Error("tx.BeginTran() error(%v)", err)
		return
	}
	var stat model.Stats
	if stat, err = d.tidbTxStat(c, tx, businessID, originID, messageID); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			log.Error("tx.Rollback() error(%v)", err1)
		}
		return
	}
	if stat.Likes+likeCount < 0 {
		likeCount = -stat.Likes
	}
	if stat.Dislikes+dislikeCount < 0 {
		dislikeCount = -stat.Dislikes
	}
	if err = d.tidbTxUpdateCounts(c, tx, businessID, originID, messageID, likeCount, dislikeCount, upMid); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			log.Error("tx.Rollback() error(%v)", err1)
		}
		return
	}
	if err = tx.Commit(); err != nil {
		PromError("like:更新喜欢计数")
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	return
}

// UpdateUpMids .
func (d *Dao) UpdateUpMids(c context.Context, businessID int64, data []*model.UpMidsReq) (rows int64, err error) {
	var tx *sql.Tx
	if tx, err = d.tidbBeginTran(c); err != nil {
		log.Error("UpdateUpMids tx.BeginTran() error(%v)", err)
		return
	}
	var res xsql.Result
	for _, u := range data {
		if res, err = tx.Exec("UPDATE counts SET up_mid = ? WHERE business_id = ? AND origin_id = ? AND message_id = ?", u.UpMid, businessID, u.OriginID, u.MessageID); err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("UpdateUpMids tx.Rollback() error(%v)", err1)
			}
			return
		}
		aff, _ := res.RowsAffected()
		rows += aff
	}
	if err = tx.Commit(); err != nil {
		PromError("like:更新UpMids")
		log.Error("UpdateUpMids tx.Commit() error(%v)", err)
	}
	return
}

// ItemHasLike .
func (d *Dao) ItemHasLike(c context.Context, businessID int64, originID, messageID int64, mids []int64, typ int8) (res map[int64]int64, err error) {
	var rows *sql.Rows
	sqlStr := fmt.Sprintf(_tidbItemHasLikeSQL, xstr.JoinInts(mids))
	if rows, err = d.tidb.Query(c, sqlStr, businessID, originID, messageID, typ); err != nil {
		PromError("tidb:itemHasLike")
		log.Error("tidb.ItemHasLike.Query asc error(%v,%v,%v,%v,%v)", businessID, originID, messageID, typ, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var t time.Time
		var mid int64
		if res == nil {
			res = make(map[int64]int64)
		}
		if err = rows.Scan(&mid, &t); err != nil {
			PromError("tidb:itemHasLikeScan")
			log.Error("tidb.rows.itemHasLike.Scan error(%v)", err)
			return
		}
		res[mid] = t.Unix()
	}
	err = rows.Err()
	return
}

// UserHasLike .
func (d *Dao) UserHasLike(c context.Context, businessID, mid int64, messageIDs []int64, typ int8) (res []*model.ItemLikeRecord, err error) {
	var rows *sql.Rows
	sqlStr := fmt.Sprintf(_tidbUserHasLikeSQL, xstr.JoinInts(messageIDs))
	if rows, err = d.tidb.Query(c, sqlStr, mid, businessID, typ); err != nil {
		PromError("tidb:userHasLike")
		log.Error("tidb.UserHasLike.Query asc error(%v,%v,%v,%v,%v)", mid, businessID, messageIDs, typ, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var t xtime.Time
		var messageID int64
		if err = rows.Scan(&messageID, &t); err != nil {
			PromError("tidb:userHasLikeScan")
			log.Error("tidb.rows.userHasLike.Scan error(%v)", err)
			return
		}
		tmp := &model.ItemLikeRecord{
			MessageID: messageID,
			Time:      t,
		}
		res = append(res, tmp)
	}
	err = rows.Err()
	return
}
