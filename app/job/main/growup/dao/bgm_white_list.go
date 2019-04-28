package dao

import (
	"context"
)

const (
	// insert
	_inBgmWhiteSQL = "INSERT INTO bgm_white_list(mid) VALUES (?) ON DUPLICATE KEY UPDATE mid = VALUES(mid)"
)

// InsertBgmWhiteList insert into bgm_white_list
func (d *Dao) InsertBgmWhiteList(c context.Context, mid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _inBgmWhiteSQL, mid)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
