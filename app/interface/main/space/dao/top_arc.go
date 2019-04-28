package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/interface/main/space/model"
	"go-common/library/log"
)

const (
	_topArcFmt    = "spc_ta_%d"
	_topArcSQL    = `SELECT aid,recommend_reason FROM member_top%d WHERE mid = ? LIMIT 1`
	_topArcAddSQL = `INSERT INTO member_top%d(mid,aid,recommend_reason) VALUES (?,?,?) ON DUPLICATE KEY UPDATE aid = ?,recommend_reason = ?`
	_topArcDelSQL = `DELETE FROM member_top%d WHERE mid = ?`
)

func topArcHit(mid int64) int64 {
	return mid % 10
}

func topArcKey(mid int64) string {
	return fmt.Sprintf(_topArcFmt, mid)
}

// RawTopArc get top aid from db.
func (d *Dao) RawTopArc(c context.Context, mid int64) (res *model.AidReason, err error) {
	var row = d.db.QueryRow(c, fmt.Sprintf(_topArcSQL, topArcHit(mid)), mid)
	res = new(model.AidReason)
	if err = row.Scan(&res.Aid, &res.Reason); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		} else {
			log.Error("RawTopArc row.Scan() error(%v)", err)
		}
	}
	return
}

// AddTopArc add top archive.
func (d *Dao) AddTopArc(c context.Context, mid, aid int64, reason string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_topArcAddSQL, topArcHit(mid)), mid, aid, reason, aid, reason); err != nil {
		log.Error("AddTopArc error d.db.Exec(%d,%d,%s) error(%v)", mid, aid, reason, err)
	}
	return
}

// DelTopArc delete top archive.
func (d *Dao) DelTopArc(c context.Context, mid int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_topArcDelSQL, topArcHit(mid)), mid); err != nil {
		log.Error("DelTopArc error d.db.Exec(%d) error(%v)", mid, err)
	}
	return
}
