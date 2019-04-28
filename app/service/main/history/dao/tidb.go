package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	pb "go-common/app/service/main/history/api/grpc"
	"go-common/app/service/main/history/model"
	"go-common/library/database/sql"
	"go-common/library/database/tidb"
	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_addHistorySQL = "INSERT INTO histories(mid, kid, business_id, aid, sid, epid, sub_type, cid, device, progress, view_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)" +
		"ON DUPLICATE KEY UPDATE aid =?, sid=?, epid=?, sub_type=?, cid=?, device=?, progress=?, view_at=?"
	_businessesSQL        = "SELECT id, name, ttl FROM business"
	_historiesSQL         = "SELECT ctime, mtime, business_id, kid, aid, sid, epid, sub_type, cid, device, progress, view_at FROM histories WHERE mid=? AND view_at < ? ORDER BY view_at DESC LIMIT ?"
	_partHistoriesSQL     = "SELECT ctime, mtime, business_id, kid, aid, sid, epid, sub_type, cid, device, progress, view_at FROM histories WHERE mid=? AND business_id in (%s) AND view_at < ? ORDER BY view_at DESC LIMIT ? "
	_queryHistoriesSQL    = "SELECT ctime, mtime, business_id, kid, aid, sid, epid, sub_type, cid, device, progress, view_at FROM histories WHERE mid=? AND kid in (%s) AND business_id = ?"
	_historySQL           = "SELECT ctime, mtime, business_id, kid, aid, sid, epid, sub_type, cid, device, progress, view_at FROM histories WHERE mid=? AND kid = ? AND business_id = ?"
	_deleteHistoriesSQL   = "DELETE FROM histories WHERE mid = ? AND kid = ? AND business_id = ?"
	_clearHistoriesSQL    = "DELETE FROM histories WHERE mid = ? AND business_id in (%s)"
	_clearAllHistoriesSQL = "DELETE FROM histories WHERE mid = ?"
	_userHide             = "SELECT hide FROM users WHERE mid = ?"
	_updateUserHide       = "INSERT INTO users(mid, hide) VALUES(?,?) ON DUPLICATE KEY UPDATE hide =?"
)

// AddHistories add histories to db
func (d *Dao) AddHistories(c context.Context, hs []*model.History) (err error) {
	if len(hs) == 0 {
		return
	}
	var tx *tidb.Tx
	if tx, err = d.tidb.Begin(c); err != nil {
		log.Error("tx.BeginTran() error(%v)", err)
		return
	}
	for _, h := range hs {
		if _, err = tx.Stmts(d.insertStmt).Exec(c, h.Mid, h.Kid, h.BusinessID, h.Aid, h.Sid, h.Epid, h.SubType, h.Cid, h.Device, h.Progress, h.ViewAt,
			h.Aid, h.Sid, h.Epid, h.SubType, h.Cid, h.Device, h.Progress, h.ViewAt); err != nil {
			log.Errorv(c, log.D{Key: "mid", Value: h.Mid}, log.D{Key: "err", Value: err}, log.D{Key: "detail", Value: h})
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("add histories commit(%+v) err: %v", hs, err)
	}
	return
}

// QueryBusinesses business
func (d *Dao) QueryBusinesses(c context.Context) (res []*model.Business, err error) {
	var rows *tidb.Rows
	if rows, err = d.businessesStmt.Query(c); err != nil {
		log.Error("db.businessesStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.Business{}
		if err = rows.Scan(&b.ID, &b.Name, &b.TTL); err != nil {
			log.Error("rows.Business.Scan error(%v)", err)
			return
		}
		res = append(res, b)
	}
	err = rows.Err()
	return
}

// UserHistories get histories by time
func (d *Dao) UserHistories(c context.Context, businesses []string, mid, viewAt, ps int64) (res map[string][]*model.History, err error) {
	var rows *tidb.Rows
	if len(businesses) == 0 {
		if rows, err = d.historiesStmt.Query(c, mid, viewAt, ps); err != nil {
			log.Error("db.Histories.Query error(%v)", err)
			return
		}
	} else {
		var ids []int64
		for _, b := range businesses {
			ids = append(ids, d.BusinessNames[b].ID)
		}
		sqlStr := fmt.Sprintf(_partHistoriesSQL, xstr.JoinInts(ids))
		if rows, err = d.tidb.Query(c, sqlStr, mid, viewAt, ps); err != nil {
			log.Error("UserHistories(%v,%d,%d,%d),db.Histories.Query error(%v)", businesses, mid, viewAt, ps, err)
			return
		}
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.History{Mid: mid}
		if err = rows.Scan(&b.Ctime, &b.Mtime, &b.BusinessID, &b.Kid, &b.Aid, &b.Sid, &b.Epid, &b.SubType, &b.Cid, &b.Device, &b.Progress, &b.ViewAt); err != nil {
			log.Error("UserHistories(%v,%d,%d,%d),rows.Scan error(%v)", businesses, mid, viewAt, ps, err)
			if strings.Contains(fmt.Sprintf("%v", err), "18446744073709551615") {
				err = nil
				continue
			}
			return
		}
		b.Business = d.Businesses[b.BusinessID].Name
		if res == nil {
			res = make(map[string][]*model.History)
		}
		res[b.Business] = append(res[b.Business], b)
	}
	err = rows.Err()
	return
}

// Histories get histories by id
func (d *Dao) Histories(c context.Context, business string, mid int64, ids []int64) (res map[int64]*model.History, err error) {
	var rows *tidb.Rows
	bid := d.BusinessNames[business].ID
	if len(ids) == 1 {
		if rows, err = d.historyStmt.Query(c, mid, ids[0], bid); err != nil {
			log.Error("db.Histories.Query error(%v)", err)
			return
		}
	} else {
		sqlStr := fmt.Sprintf(_queryHistoriesSQL, xstr.JoinInts(ids))
		if rows, err = d.tidb.Query(c, sqlStr, mid, bid); err != nil {
			log.Error("tidb.Histories.Query error(%v)", err)
			return
		}
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.History{Mid: mid}
		if err = rows.Scan(&b.Ctime, &b.Mtime, &b.BusinessID, &b.Kid, &b.Aid, &b.Sid, &b.Epid, &b.SubType, &b.Cid, &b.Device, &b.Progress, &b.ViewAt); err != nil {
			log.Error("rows.Business.Scan error(%v)", err)
			return
		}
		b.Business = d.Businesses[b.BusinessID].Name
		if res == nil {
			res = make(map[int64]*model.History)
		}
		res[b.Kid] = b
	}
	err = rows.Err()
	return
}

// DeleteHistories .
func (d *Dao) DeleteHistories(c context.Context, h *pb.DelHistoriesReq) (err error) {
	var tx *tidb.Tx
	if tx, err = d.tidb.Begin(c); err != nil {
		log.Error("tx.BeginTran() error(%v)", err)
		return
	}
	for _, r := range h.Records {
		_, err = tx.Stmts(d.deleteHistoriesStmt).Exec(c, h.Mid, r.ID, d.BusinessNames[r.Business].ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("DeleteHistories(%+v),commit err:%+v", h, err)
	}
	return
}

// ClearHistory clear histories
func (d *Dao) ClearHistory(c context.Context, mid int64, businesses []string) (err error) {
	var ids []string
	for _, b := range businesses {
		ids = append(ids, strconv.FormatInt(d.BusinessNames[b].ID, 10))
	}
	sqlStr := fmt.Sprintf(_clearHistoriesSQL, strings.Join(ids, ","))
	if _, err = d.tidb.Exec(c, sqlStr, mid); err != nil {
		log.Error("mid: %d clear(%v) err: %v", mid, businesses, err)
	}
	return
}

// ClearAllHistory clear all histories
func (d *Dao) ClearAllHistory(c context.Context, mid int64) (err error) {
	if _, err = d.clearAllHistoriesStmt.Exec(c, mid); err != nil {
		log.Error("mid: %d clear all err: %v", mid, err)
	}
	return
}

// UpdateUserHide update user hide
func (d *Dao) UpdateUserHide(c context.Context, mid int64, hide bool) (err error) {
	if _, err = d.updateUserHideStmt.Exec(c, mid, hide, hide); err != nil {
		log.Error("mid: %d updateUserHide(%v) err: %v", mid, hide, err)
	}
	return
}

// UserHide get user hide
func (d *Dao) UserHide(c context.Context, mid int64) (hide bool, err error) {
	if err = d.userHideStmt.QueryRow(c, mid).Scan(&hide); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("mid: %d UserHide err: %v", mid, err)
	}
	return
}
