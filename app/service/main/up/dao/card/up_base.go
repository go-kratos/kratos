package card

import (
	"context"
	"fmt"
	"go-common/library/log"
)

const _listUpBaseSQL = "SELECT id, mid FROM up_base_info WHERE id > ? AND business_type = 1 %s LIMIT ?"

// ListUpBase list <id, mid> k-v pairs
func (d *Dao) ListUpBase(c context.Context, size int, lastID int64, where string) (idMids map[int64]int64, err error) {
	idMids = make(map[int64]int64)
	rows, err := d.db.Query(c, fmt.Sprintf(_listUpBaseSQL, where), lastID, size)
	if err != nil {
		log.Error("ListUpBase d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var mid int64
		var id int64

		err = rows.Scan(&id, &mid)
		if err != nil {
			log.Error("ListUpBase rows.Scan error(%v)", err)
			return
		}
		idMids[id] = mid
	}

	return
}
