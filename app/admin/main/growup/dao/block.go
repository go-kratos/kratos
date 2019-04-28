package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/admin/main/growup/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	// insert
	_inBlockedSQL = "INSERT INTO up_blocked (mid,nickname,original_archive_count,category_id,fans,apply_at,is_deleted) VALUES (?,?,?,?,?,?,0) ON DUPLICATE KEY UPDATE mid=?,nickname=?,original_archive_count=?,category_id=?,fans=?,apply_at=?,is_deleted=0"
	// select
	_upInfoVideoSQL  = "SELECT apply_at FROM up_info_video WHERE mid=?"
	_blockedCountSQL = "SELECT count(*) FROM up_blocked WHERE %s"
	_blockedSQL      = "SELECT mid,nickname,original_archive_count,category_id,fans,apply_at FROM up_blocked WHERE %s"
	_isBlockedSQL    = "SELECT mid FROM up_blocked WHERE mid=? AND is_deleted=0"
	// update
	_upBlockedStateSQL = "UPDATE up_blocked SET is_deleted=? WHERE mid=?"
)

// ApplyAt query apply_at from up_info_video by mid
func (d *Dao) ApplyAt(c context.Context, mid int64) (applyAt xtime.Time, err error) {
	row := d.rddb.QueryRow(c, _upInfoVideoSQL, mid)
	if err = row.Scan(&applyAt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// InsertBlocked insert block up to blacklist
func (d *Dao) InsertBlocked(c context.Context, v *model.Blocked) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _inBlockedSQL, v.MID, v.Nickname, v.OriginalArchiveCount, v.MainCategory, v.Fans, v.ApplyAt, v.MID, v.Nickname, v.OriginalArchiveCount, v.MainCategory, v.Fans, v.ApplyAt)
	if err != nil {
		log.Error("db.inBlockedStmt.Exec(%s) error(%v)", _inBlockedSQL, err)
		return
	}
	return res.RowsAffected()
}

// UpdateBlockedState update blocked is_deleted
func (d *Dao) UpdateBlockedState(c context.Context, mid int64, del int) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _upBlockedStateSQL, del, mid)
	if err != nil {
		log.Error("db.upBlockedState.Exec(%s) error(%v)", _upBlockedStateSQL, err)
		return
	}
	return res.RowsAffected()
}

// DelFromBlocked del blocked
func (d *Dao) DelFromBlocked(c context.Context, mid int64) (rows int64, err error) {
	return d.UpdateBlockedState(c, mid, 1)
}

// BlockCount get blocked count
func (d *Dao) BlockCount(c context.Context, query string) (count int, err error) {
	row := d.rddb.QueryRow(c, fmt.Sprintf(_blockedCountSQL, query))
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("rows.Scan error(%v)", err)
		}
	}
	return
}

// QueryFromBlocked query blocked user in black list
func (d *Dao) QueryFromBlocked(c context.Context, query string) (ups []*model.Blocked, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_blockedSQL, query))
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		v := &model.Blocked{}
		err = rows.Scan(&v.MID, &v.Nickname, &v.OriginalArchiveCount, &v.MainCategory, &v.Fans, &v.ApplyAt)
		if err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		ups = append(ups, v)
	}
	return
}

// Blocked check user is blocked
func (d *Dao) Blocked(c context.Context, mid int64) (id int64, err error) {
	row := d.rddb.QueryRow(c, _isBlockedSQL, mid)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row Scan error(%v)", err)
		}
	}
	return
}
