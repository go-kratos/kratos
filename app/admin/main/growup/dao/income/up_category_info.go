package income

import (
	"context"
	"fmt"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_upInfoSQL = "SELECT mid, nick_name FROM up_category_info WHERE mid in (%s) AND is_deleted = 0"
)

// ListUpInfo list up_category_info by mids
func (d *Dao) ListUpInfo(c context.Context, mids []int64) (upInfo map[int64]string, err error) {
	upInfo = make(map[int64]string)
	if len(mids) == 0 {
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_upInfoSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("ListUpInfo d.db.Query error(%v)", err)
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
		upInfo[mid] = nickname
	}
	err = rows.Err()
	return
}
