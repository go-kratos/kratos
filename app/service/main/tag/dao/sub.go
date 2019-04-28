package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

func (d *Dao) shardSub(mid int64) string {
	return fmt.Sprintf("%03d", mid%int64(_shard))
}

var (
	_addSubSQL = "INSERT IGNORE INTO subscriber_%s (mid,tid,state) VALUES (?,?,0) ON DUPLICATE KEY UPDATE state=0"
)

// TxAddSub .
func (d *Dao) TxAddSub(tx *xsql.Tx, mid, tid int64) (affected int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_addSubSQL, d.shardSub(mid)), mid, tid)
	if err != nil {
		log.Error("tx.Exec(%d,%d) error(%v)", tid, mid, err)
		return
	}
	return res.RowsAffected()
}

var (
	_delSubSQL = "UPDATE subscriber_%s SET state=1 WHERE mid=? AND tid=?"
)

// TxDelSub .
func (d *Dao) TxDelSub(tx *xsql.Tx, mid, tid int64) (affected int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_delSubSQL, d.shardSub(mid)), mid, tid)
	if err != nil {
		log.Error("tx.Exec(%d,%d), error(%v)", tid, mid, err)
		return
	}
	return res.RowsAffected()
}

const (
	_selectSubsSQL = "SELECT tid,mtime FROM subscriber_%s WHERE mid=? AND state=0"
)

// Sub .
func (d *Dao) Sub(c context.Context, mid int64) (res []*model.Sub, rem map[int64]*model.Sub, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_selectSubsSQL, d.shardSub(mid)), mid)
	if err != nil {
		log.Error("d.db.Query(%d) error(%v)", mid, err)
		return
	}
	res = make([]*model.Sub, 0)
	rem = make(map[int64]*model.Sub)
	defer rows.Close()
	for rows.Next() {
		r := &model.Sub{}
		if err = rows.Scan(&r.Tid, &r.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		if r.Tid <= 0 {
			continue
		}
		res = append(res, r)
		rem[r.Tid] = r
	}
	return
}

// SubList .
func (d *Dao) SubList(c context.Context, mid int64) (res []*model.Sub, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_selectSubsSQL, d.shardSub(mid)), mid)
	if err != nil {
		log.Error("d.db.Query(%d) error(%v)", mid, err)
		return
	}
	res = make([]*model.Sub, 0)
	defer rows.Close()
	for rows.Next() {
		r := &model.Sub{}
		if err = rows.Scan(&r.Tid, &r.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		if r.Tid <= 0 {
			continue
		}
		res = append(res, r)
	}
	return
}
