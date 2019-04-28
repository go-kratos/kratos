package income

import (
	"context"

	"go-common/library/log"
)

const (
	_businessOrderSQL = "SELECT id,av_id FROM business_order_sheet WHERE id > ? ORDER BY id LIMIT ?"
)

// BusinessOrders get business order
func (d *Dao) BusinessOrders(c context.Context, offset, limit int64) (last int64, m map[int64]bool, err error) {
	rows, err := d.rddb.Query(c, _businessOrderSQL, offset, limit)
	if err != nil {
		log.Error("d.rddb.Query BusinessOrders error(%v)", err)
		return
	}
	defer rows.Close()
	m = make(map[int64]bool)
	for rows.Next() {
		var avID int64
		err = rows.Scan(&last, &avID)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		m[avID] = true
	}
	return
}
