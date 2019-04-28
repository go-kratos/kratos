package income

import (
	"context"

	"go-common/library/log"
)

const (
	_blacklistSQL = "SELECT id,av_id,ctype,is_delete FROM av_black_list WHERE id > ? ORDER BY id LIMIT ?"
)

// Blacklist map[ctype][]id
func (d *Dao) Blacklist(c context.Context, id int64, limit int64) (m map[int][]int64, last int64, err error) {
	rows, err := d.db.Query(c, _blacklistSQL, id, limit)
	if err != nil {
		return
	}
	m = make(map[int][]int64)
	for rows.Next() {
		var ctype, isDeleted int
		var avID int64
		err = rows.Scan(&last, &avID, &ctype, &isDeleted)
		if err != nil {
			log.Error("Rows Scan error(%v)", err)
			return
		}
		if isDeleted == 0 {
			if _, ok := m[ctype]; ok {
				m[ctype] = append(m[ctype], avID)
			} else {
				m[ctype] = []int64{avID}
			}
		}
	}
	return
}
