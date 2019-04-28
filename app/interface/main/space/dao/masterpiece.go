package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/space/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_masterpieceFmt     = "spc_mp_%d"
	_masterpieceSQL     = `SELECT aid,recommend_reason FROM member_masterpiece%d WHERE mid = ?`
	_masterpieceAddSQL  = `INSERT INTO member_masterpiece%d (mid,aid,recommend_reason) VALUES (?,?,?)`
	_masterpieceEditSQL = `UPDATE member_masterpiece%d SET aid = ?,recommend_reason = ? WHERE mid = ? AND aid = ?`
	_masterpieceDelSQL  = `DELETE FROM member_masterpiece%d WHERE mid = ? AND aid = ?`
)

func masterpieceHit(mid int64) int64 {
	return mid % 10
}

func masterpieceKey(mid int64) string {
	return fmt.Sprintf(_masterpieceFmt, mid)
}

// RawMasterpiece get masterpiece from db.
func (d *Dao) RawMasterpiece(c context.Context, mid int64) (res *model.AidReasons, err error) {
	var (
		rows *xsql.Rows
		list []*model.AidReason
	)
	res = new(model.AidReasons)
	if rows, err = d.db.Query(c, fmt.Sprintf(_masterpieceSQL, masterpieceHit(mid)), mid); err != nil {
		log.Error("RawMasterpiece d.db.Query(%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.AidReason)
		if err = rows.Scan(&r.Aid, &r.Reason); err != nil {
			log.Error("RawMasterpiece row.Scan() error(%v)", err)
			return
		}
		list = append(list, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("RawMasterpiece rows.error(%v)", err)
		return
	}
	res.List = list
	return
}

// AddMasterpiece add masterpiece data.
func (d *Dao) AddMasterpiece(c context.Context, mid, aid int64, reason string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_masterpieceAddSQL, masterpieceHit(mid)), mid, aid, reason); err != nil {
		log.Error("AddMasterpiece error d.db.Exec(%d,%d,%s) error(%v)", mid, aid, reason, err)
	}
	return
}

// EditMasterpiece edit masterpiece data.
func (d *Dao) EditMasterpiece(c context.Context, mid, aid, preAid int64, reason string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_masterpieceEditSQL, masterpieceHit(mid)), aid, reason, mid, preAid); err != nil {
		log.Error("EditMasterpiece error d.db.Exec(%d,%d,%s) error(%v)", mid, aid, reason, err)
	}
	return
}

// DelMasterpiece delete masterpiece.
func (d *Dao) DelMasterpiece(c context.Context, mid, aid int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_masterpieceDelSQL, masterpieceHit(mid)), mid, aid); err != nil {
		log.Error("DelMasterpiece error d.db.Exec(%d,%d) error(%v)", mid, aid, err)
	}
	return
}
