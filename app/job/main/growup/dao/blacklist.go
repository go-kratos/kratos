package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	// list all blacklist
	_listBlacklistSQL = "SELECT av_id,mid,reason,ctype,has_signed,nickname FROM av_black_list WHERE %s is_delete = 0 LIMIT ?,?"
	// query last record ctime by reason and id
	_queryOneByReasonSQL = "SELECT ctime FROM av_black_list WHERE reason = ? AND is_delete = 0 ORDER BY id DESC LIMIT 1"
	// query record from bilibili_business_up_cooperate.business_order_sheet
	_queryExecuteOrderSQL = "SELECT av_id, making_order_up_mid, ctime FROM business_order_sheet WHERE ctime >= ? AND ctime <= ? AND is_deleted = 0 ORDER BY ctime DESC"

	// add to blacklist batch
	_addBlacklistBatchSQL = "INSERT INTO av_black_list(av_id, mid, reason, ctype, has_signed, nickname) VALUES %s ON DUPLICATE KEY UPDATE mid=VALUES(mid),reason=VALUES(reason),has_signed=VALUES(has_signed),nickname=VALUES(nickname)"
	// query mid, nickname from up_info_video
	_queryHasSignUpInfoSQL = "SELECT mid, nickname FROM up_info_video where account_state = 3 AND is_deleted = 0 limit ?, ?"
)

// ListBlacklist list all blacklist
func (d *Dao) ListBlacklist(c context.Context, query string, from, limit int) (backlists []*model.Blacklist, err error) {
	if query != "" {
		query += " AND"
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_listBlacklistSQL, query), from, limit)
	if err != nil {
		log.Error("ListBlacklist d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.Blacklist{}
		err = rows.Scan(&list.AvID, &list.MID, &list.Reason, &list.CType, &list.HasSigned, &list.Nickname)
		if err != nil {
			log.Error("ListBlacklist rows scan error(%v)", err)
			return
		}
		backlists = append(backlists, list)
	}

	err = rows.Err()
	return
}

// GetExecuteOrder get execute order by date
func (d *Dao) GetExecuteOrder(c context.Context, startTime, endTime time.Time) (executeOrders []*model.ExecuteOrder, err error) {
	rows, err := d.rddb.Query(c, _queryExecuteOrderSQL, startTime, endTime)
	if err != nil {
		log.Error("GetExecuteOrder d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		order := &model.ExecuteOrder{}
		err = rows.Scan(&order.AvID, &order.MID, &order.CTime)
		if err != nil {
			log.Error("GetExecuteOrder rows.Scan error(%v)", err)
			return
		}
		if order.AvID > 0 {
			executeOrders = append(executeOrders, order)
		}
	}
	err = rows.Err()
	return
}

// GetLastCtime get last ctime by query
func (d *Dao) GetLastCtime(c context.Context, reason int) (ctime int64, err error) {
	row := d.db.QueryRow(c, _queryOneByReasonSQL, reason)
	var t time.Time
	err = row.Scan(&t)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			ctime = 0
		} else {
			log.Error("GetLastCtime row Scan error(%v)", err)
		}
	} else {
		ctime = t.Unix()
	}
	return
}

// AddBlacklistBatch add batch to blacklist
func (d *Dao) AddBlacklistBatch(c context.Context, blacklist []*model.Blacklist) (count int64, err error) {
	if len(blacklist) == 0 {
		return
	}

	var vals string
	for _, row := range blacklist {
		vals += fmt.Sprintf("(%d, %d, %d, %d, %d, '%s'),", row.AvID, row.MID, row.Reason, row.CType, row.HasSigned, row.Nickname)
	}

	if len(vals) > 0 {
		vals = vals[0 : len(vals)-1]
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_addBlacklistBatchSQL, vals))
	if err != nil {
		log.Error("AddBlacklistBatch d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// GetHasSignUpInfo get all has signed up info
func (d *Dao) GetHasSignUpInfo(c context.Context, offset, limit int, m map[int64]string) (err error) {
	rows, err := d.db.Query(c, _queryHasSignUpInfoSQL, offset, limit)
	if err != nil {
		log.Error("GetHasSignUpInfo d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		var nickname string
		err = rows.Scan(&mid, &nickname)
		if err != nil {
			log.Error("GetHasSignUpInfo rows.Scan error(%v)", err)
			return
		}
		m[mid] = nickname
	}
	err = rows.Err()
	return
}
