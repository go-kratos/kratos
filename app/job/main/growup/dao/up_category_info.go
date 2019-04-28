package dao

import (
	"context"
	"fmt"

	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_upInfoSQL   = "SELECT mid, nick_name FROM up_category_info WHERE mid in (%s) AND is_deleted = 0"
	_nicknameSQL = "SELECT nick_name FROM up_category_info WHERE mid = ? AND is_deleted = 0"
)

// ListUpNickname list up_category_info by mids
func (d *Dao) ListUpNickname(c context.Context, mids []int64, m map[int64]string) (err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_upInfoSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("ListUpNickname d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var mid int64
		var nickname string
		err = rows.Scan(&mid, &nickname)
		if err != nil {
			log.Error("ListUpInfo rows scan error(%v)", err)
			return
		}
		m[mid] = nickname
	}
	err = rows.Err()
	return
}

// GetNickname get nickname by mid.
func (d *Dao) GetNickname(c context.Context, mid int64) (nickname string, err error) {
	row := d.db.QueryRow(c, _nicknameSQL, mid)
	if err = row.Scan(&nickname); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}
