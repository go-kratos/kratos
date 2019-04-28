package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-common/app/interface/main/space/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_chSub            = 10
	_chSQL            = "SELECT id,mid,name,intro,modify_time FROM member_channel%d WHERE mid = ? AND id = ? LIMIT 1"
	_chListSQL        = "SELECT id,mid,name,intro,modify_time FROM member_channel%d WHERE mid = ?"
	_chCntSQL         = "SELECT COUNT(1) FROM member_channel%d WHERE mid = ? LIMIT 1"
	_chAddSQL         = "INSERT INTO member_channel%d (mid,name,intro,modify_time) VALUES (?,?,?,?)"
	_chEditSQL        = "UPDATE member_channel%d SET name = ?,intro = ?,modify_time = ? WHERE mid = ? AND id = ?"
	_chDelSQL         = "DELETE FROM member_channel%d WHERE mid = ? AND id = ?"
	_chArcSQL         = "SELECT id,mid,cid,aid,order_num,modify_time FROM member_channel_video%d WHERE mid = ? AND cid = ? ORDER BY order_num"
	_chArcCntSQL      = "SELECT COUNT(1),IFNULL(aid , 0) FROM member_channel_video%d WHERE mid = ? AND cid = ? ORDER BY order_num DESC LIMIT 1"
	_chArcAddSQL      = "INSERT INTO member_channel_video%d (mid,cid,aid,order_num,modify_time) VALUES %s"
	_chArcDelSQL      = "DELETE FROM member_channel_video%d WHERE mid = ? AND cid = ? AND aid = ?"
	_chArcsDelSQL     = "DELETE FROM member_channel_video%d WHERE mid = ? AND cid = ?"
	_chArcsSortDelSQL = "UPDATE member_channel_video%d SET order_num = order_num-1 WHERE mid = ? AND cid = ? AND order_num > ?"
	_chArcEditSQL     = "UPDATE member_channel_video%d SET order_num = CASE %s END,modify_time = ? WHERE mid = ? AND cid = ? AND aid IN (%s)"
)

func channelHit(mid int64) int64 {
	return mid % _chSub
}

// Channel get channel simple info
func (d *Dao) Channel(c context.Context, mid, cid int64) (channel *model.Channel, err error) {
	row := d.channelStmt[channelHit(mid)].QueryRow(c, mid, cid)
	channel = new(model.Channel)
	if err = row.Scan(&channel.Cid, &channel.Mid, &channel.Name, &channel.Intro, &channel.Mtime); err != nil {
		if err == sql.ErrNoRows {
			channel = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// ChannelCnt get channel count by mid
func (d *Dao) ChannelCnt(c context.Context, mid int64) (count int, err error) {
	row := d.channelCntStmt[channelHit(mid)].QueryRow(c, mid)
	if err = row.Scan(&count); err != nil {
		log.Error("row.Scan error(%v)", err)
	}
	return
}

// ChannelList get channel list by mid
func (d *Dao) ChannelList(c context.Context, mid int64) (res []*model.Channel, err error) {
	var rows *xsql.Rows
	if rows, err = d.channelListStmt[channelHit(mid)].Query(c, mid); err != nil {
		log.Error("d.channelListStmt[%d].Query(%d) error(%v)", channelHit(mid), mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Channel)
		if err = rows.Scan(&r.Cid, &r.Mid, &r.Name, &r.Intro, &r.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// ChannelVideos get channel video list
func (d *Dao) ChannelVideos(c context.Context, mid, cid int64, order bool) (res []*model.ChannelArc, err error) {
	var (
		orderSQL string
		rows     *xsql.Rows
	)
	orderSQL = fmt.Sprintf(_chArcSQL, channelHit(mid))
	if !order {
		orderSQL = orderSQL + " DESC"
	}
	if rows, err = d.db.Query(c, orderSQL, mid, cid); err != nil {
		log.Error("d.channelVideoStmt[%d].Query(%d,%d) error(%v)", channelHit(mid), mid, cid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ChannelArc)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Cid, &r.Aid, &r.OrderNum, &r.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// ChannelExtra get channel count and cover aid
func (d *Dao) ChannelExtra(c context.Context, mid, cid int64) (res *model.ChannelExtra, err error) {
	row := d.channelArcCntStmt[channelHit(mid)].QueryRow(c, mid, cid)
	res = &model.ChannelExtra{Cid: cid}
	if err = row.Scan(&res.Count, &res.Aid); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// AddChannel add channel
func (d *Dao) AddChannel(c context.Context, mid int64, name, intro string, ts time.Time) (lastID int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_chAddSQL, channelHit(mid)), mid, name, intro, ts.Format("2006-01-02 15:04:05")); err != nil {
		log.Error("AddChannel: db.Exec(mid:%d,name:%s,intro:%s) error(%v)", mid, name, intro, err)
		return
	}
	return res.LastInsertId()
}

// EditChannel edit channel name or intro
func (d *Dao) EditChannel(c context.Context, mid, cid int64, name, intro string, ts time.Time) (affected int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_chEditSQL, channelHit(mid)), name, intro, ts.Format("2006-01-02 15:04:05"), mid, cid); err != nil {
		log.Error("EditChannel: db.Exec(mid:%d,cid:%d,name:%s,intro:%s) error(%v)", mid, cid, name, intro, err)
		return
	}
	return res.RowsAffected()
}

// DelChannel delete channel
func (d *Dao) DelChannel(c context.Context, mid, cid int64) (affected int64, err error) {
	var (
		res sql.Result
		tx  *xsql.Tx
	)
	if tx, err = d.db.Begin(c); err != nil {
		log.Error("d.db.Begin error(%v)", err)
		return
	}
	if res, err = tx.Exec(fmt.Sprintf(_chDelSQL, channelHit(mid)), mid, cid); err != nil {
		tx.Rollback()
		log.Error("DelChannel: db.Exec(mid:%d,cid:%d) error(%v)", mid, cid, err)
		return
	}
	if _, err = tx.Exec(fmt.Sprintf(_chArcsDelSQL, channelHit(mid)), mid, cid); err != nil {
		tx.Rollback()
		log.Error("DelChannelArcs: db.Exec(mid:%d,cid:%d) error(%v)", mid, cid, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// AddChannelArc add archives to channel
func (d *Dao) AddChannelArc(c context.Context, mid, cid int64, ts time.Time, chs []*model.ChannelArcSort) (lastID int64, err error) {
	var (
		res    sql.Result
		values []string
	)
	for _, v := range chs {
		values = append(values, fmt.Sprintf("(%d,%d,%d,%d,'%s')", mid, cid, v.Aid, v.OrderNum, ts.Format("2006-01-02 15:04:05")))
	}
	if res, err = d.db.Exec(c, fmt.Sprintf(_chArcAddSQL, channelHit(mid), strings.Join(values, ","))); err != nil {
		log.Error("AddChannelArc: db.Exec(%d,%d) error(%v)", mid, cid, err)
		return
	}
	return res.LastInsertId()
}

// DelChannelArc delete channel archives
func (d *Dao) DelChannelArc(c context.Context, mid, cid, aid int64, orderNum int) (affected int64, err error) {
	var (
		res sql.Result
		tx  *xsql.Tx
	)
	if tx, err = d.db.Begin(c); err != nil {
		log.Error("d.db.Begin error(%v)", err)
		return
	}
	if res, err = d.db.Exec(c, fmt.Sprintf(_chArcDelSQL, channelHit(mid)), mid, cid, aid); err != nil {
		tx.Rollback()
		log.Error("DelChannelArc: db.Exec(mid:%d,cid:%d,aid:%d) error (%v)", mid, cid, aid, err)
		return
	}
	if _, err = tx.Exec(fmt.Sprintf(_chArcsSortDelSQL, channelHit(mid)), mid, cid, orderNum); err != nil {
		tx.Rollback()
		log.Error("DelChannelArcs: db.Exec(mid:%d,cid:%d) error(%v)", mid, cid, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// EditChannelArc edit channel archive order num
func (d *Dao) EditChannelArc(c context.Context, mid, cid int64, ts time.Time, chSort []*model.ChannelArcSort) (affected int64, err error) {
	var (
		caseStr string
		aids    []int64
		res     sql.Result
	)
	for _, v := range chSort {
		caseStr = fmt.Sprintf("%s WHEN aid = %d THEN %d", caseStr, v.Aid, v.OrderNum)
		aids = append(aids, v.Aid)
	}
	if res, err = d.db.Exec(c, fmt.Sprintf(_chArcEditSQL, channelHit(mid), caseStr, xstr.JoinInts(aids)), ts.Format("2006-01-02 15:04:05"), mid, cid); err != nil {
		log.Error("EditChannel: db.Exec(mid:%d,cid:%d,%s) error(%v)", mid, cid, caseStr, err)
		return
	}
	return res.RowsAffected()
}
