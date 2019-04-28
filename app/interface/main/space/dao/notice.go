package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/interface/main/space/model"
	"go-common/library/log"
)

const (
	_noticeKeyFmt = "spc_nt_%d"
	_noticeSQL    = `SELECT notice,is_forbid FROM member_up_notice%d WHERE mid = ?`
	_noticeSetSQL = `INSERT INTO member_up_notice%d (mid,notice) VALUES (?,?) ON DUPLICATE KEY UPDATE notice = ?`
)

func noticeHit(mid int64) int64 {
	return mid % 10
}

func noticeKey(mid int64) string {
	return fmt.Sprintf(_noticeKeyFmt, mid)
}

// RawNotice get notice from db.
func (d *Dao) RawNotice(c context.Context, mid int64) (res *model.Notice, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_noticeSQL, noticeHit(mid)), mid)
	res = new(model.Notice)
	if err = row.Scan(&res.Notice, &res.IsForbid); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("Notice row.Scan() error(%v)", err)
		}
	}
	return
}

// SetNotice change notice.
func (d *Dao) SetNotice(c context.Context, mid int64, notice string) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_noticeSetSQL, noticeHit(mid)), mid, notice, notice); err != nil {
		log.Error("SetNotice error d.db.Exec(%d,%s) error(%v)", mid, notice, err)
	}
	return
}
