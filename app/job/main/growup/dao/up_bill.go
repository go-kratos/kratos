package dao

import (
	"context"
	"fmt"

	"go-common/library/log"
)

const (
	_upSignedAvsSQL = "SELECT id, mid, avs FROM up_signed_avs WHERE id > ? ORDER BY id LIMIT ?"

	_inUpBillSQL = "INSERT INTO up_bill(mid,first_income,max_income,total_income,av_count,av_max_income,av_id,quality_value,defeat_num,title,share_items,first_time,max_time,signed_at,end_at) VALUES %s"
)

// InsertUpBillBatch insert up_bill
func (d *Dao) InsertUpBillBatch(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpBillSQL, values))
	if err != nil {
		log.Error("d.db.Exec InsertUpBill error (%v)", err)
		return
	}
	return res.RowsAffected()
}

// ListUpSignedAvs list up_signed_avs
func (d *Dao) ListUpSignedAvs(c context.Context, id int64, limit int) (ups map[int64]int64, last int64, err error) {
	ups = make(map[int64]int64)
	rows, err := d.db.Query(c, _upSignedAvsSQL, id, limit)
	if err != nil {
		log.Error("ListUpSignedAvs d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid, avs int64
		err = rows.Scan(&last, &mid, &avs)
		if err != nil {
			log.Error("ListUpSignedAvs rows.Scan error(%v)", err)
			return
		}
		ups[mid] = avs
	}
	return
}
