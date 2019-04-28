package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/thumbup/model"
	sql "go-common/library/database/tidb"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_statSQL            = "SELECT message_id, likes_count, dislikes_count, origin_id, likes_change, dislikes_change FROM counts WHERE business_id = ? AND origin_id = ? AND message_id = ?"
	_itemLikesSQL       = "SELECT mid, mtime FROM likes WHERE business_id =? AND origin_id =? AND message_id =? AND type = ? ORDER BY mtime desc LIMIT ?"
	_likeStateSQL       = "SELECT type FROM likes WHERE mid=? AND business_id=? AND origin_id=? AND message_id=?"
	_businessesSQL      = "SELECT id, name, message_list_type, user_likes_limit, message_likes_limit, enable_originid, user_list_type FROM business WHERE dtime = '0000-00-00 00:00:00'"
	_updateCountsSQL    = "INSERT INTO counts (business_id, origin_id, message_id, likes_count, dislikes_count, up_mid) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE "
	_userLikeListSQL    = "SELECT message_id, mtime FROM likes WHERE mid = ? AND business_id =? AND type = ? ORDER BY mtime desc LIMIT ?"
	_updateLikeStateSQL = "INSERT INTO likes (business_id, origin_id, message_id, mid, type, mtime) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE type=?, mtime=?"
)

// Business .
func (d *Dao) Business(c context.Context) (res []*model.Business, err error) {
	var rows *sql.Rows
	if rows, err = d.businessesStmt.Query(c); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.Business{}
		if err = rows.Scan(&b.ID, &b.Name, &b.MessageListType, &b.UserLikesLimit, &b.MessageLikesLimit, &b.EnableOriginID, &b.UserListType); err != nil {
			return
		}
		res = append(res, b)
	}
	err = rows.Err()
	return
}

// UserLikes .
func (d *Dao) UserLikes(c context.Context, mid, businessID int64, typ int8, limit int) (res []*model.ItemLikeRecord, err error) {
	var rows *sql.Rows
	if rows, err = d.userLikesStmt.Query(c, mid, businessID, typ, limit); err != nil {
		log.Error("d.tidb.userList error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.ItemLikeRecord{}
		var t time.Time
		if err = rows.Scan(&b.MessageID, &t); err != nil {
			log.Error("tidb.rows.userList error(%v)", err)
			return
		}
		b.Time = xtime.Time(t.Unix())
		res = append(res, b)
	}
	err = rows.Err()
	return
}

// ItemLikes .
func (d *Dao) ItemLikes(c context.Context, businessID, originID, messageID int64, typ int8, limit int) (res []*model.UserLikeRecord, err error) {
	var rows *sql.Rows
	if rows, err = d.itemLikesStmt.Query(c, businessID, originID, messageID, typ, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.UserLikeRecord{}
		var t time.Time
		if err = rows.Scan(&b.Mid, &t); err != nil {
			return
		}
		b.Time = xtime.Time(t.Unix())
		res = append(res, b)
	}
	err = rows.Err()
	return
}

// LikeState .
func (d *Dao) LikeState(c context.Context, mid, businessID, originID, messageID int64) (res int8, err error) {
	if err = d.likeStateStmt.QueryRow(c, mid, businessID, originID, messageID).Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("tidbLikeState(%d,%d,%d,%d) error(%v)", mid, businessID, originID, messageID, err)
	}
	return
}

// UpdateLikeState .
func (d *Dao) UpdateLikeState(c context.Context, mid, businessID, originID, messageID int64, state int8, likeTime time.Time) (res model.Stats, err error) {
	var tx *sql.Tx
	if tx, err = d.tidb.Begin(c); err != nil {
		log.Error("d.tidb.Begin() error(%v)", err)
		return
	}
	if _, err = tx.Stmts(d.updateLikeStateStmt).Exec(c, businessID, originID, messageID, mid, state, likeTime, state, likeTime); err != nil {
		log.Error("d.tidbUpdateLikeState error(%v)", err)
		if err1 := tx.Rollback(); err1 != nil {
			log.Error("tx.Rollback error(%v)", err)
		}
		return
	}
	if res, err = d.tidbTxStat(c, tx, businessID, originID, messageID); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			log.Error("tx.Rollback() error(%v)", err1)
		}
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
	}
	return
}

// UpdateCounts .
func (d *Dao) UpdateCounts(c context.Context, businessID, originID, messageID, likeCounts, dislikeCounts, upMid int64) (err error) {
	var tx *sql.Tx
	if tx, err = d.tidb.Begin(c); err != nil {
		log.Error("tx.tidb.Begin() error(%v)", err)
		return
	}
	var stat model.Stats
	if stat, err = d.tidbTxStat(c, tx, businessID, originID, messageID); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			log.Error("tx.Rollback() error(%v)", err1)
		}
		return
	}
	if stat.Likes+likeCounts < 0 {
		likeCounts = -stat.Likes
	}
	if stat.Dislikes+dislikeCounts < 0 {
		dislikeCounts = -stat.Dislikes
	}
	if err = d.tidbTxUpdateCounts(c, tx, businessID, originID, messageID, likeCounts, dislikeCounts, upMid); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			log.Error("tx.Rollback() error(%v)", err1)
		}
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
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
		log.Error("tidb: Tx db.Stat.Query error(%v,%v,%v,%v)", businessID, originID, messageID, err)
		return
	}
	res.Likes += likeChange
	res.Dislikes += dislikeChange
	return
}

// TxUpdateCounts .
func (d *Dao) tidbTxUpdateCounts(c context.Context, tx *sql.Tx, businessID, originID, messageID int64, likesCount, dislikesCount, upMid int64) (err error) {
	likeSQL := _updateCountsSQL
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
		log.Error("dao.tidbTxUpdateCounts tx.Exec(%s, %v,%v,%v) error(%v)", likeSQL, businessID, originID, messageID, err)
	}
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
		log.Error("tidb.Stat.Query error(%v)", err)
	}
	res.Likes += likeChange
	res.Dislikes += dislikeChange
	return
}
