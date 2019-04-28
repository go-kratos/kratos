package income

import (
	"context"
)

const (
	_getBubbleMetaSQL = "SELECT id,av_id,b_type FROM lottery_av_info WHERE id > ? ORDER BY id LIMIT ?"
)

// GetBubbleMeta get bubble meta
func (d *Dao) GetBubbleMeta(c context.Context, id int64, limit int64) (data map[int64][]int, last int64, err error) {
	data = make(map[int64][]int)
	rows, err := d.db.Query(c, _getBubbleMetaSQL, id, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			avID  int64
			bType int
		)
		err = rows.Scan(&last, &avID, &bType)
		if err != nil {
			return
		}
		data[avID] = append(data[avID], bType)
	}
	return
}
