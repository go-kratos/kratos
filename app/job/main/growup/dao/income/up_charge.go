package income

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/income"

	"go-common/library/log"
)

const (
	_insertUpChargeSQL = "INSERT INTO %s (mid,inc_charge,total_charge,date) VALUES %s ON DUPLICATE KEY UPDATE inc_charge=VALUES(inc_charge),total_charge=VALUES(total_charge)"
	_upChargeSQL       = "SELECT id,mid,inc_charge,total_charge,date FROM %s WHERE date=? AND id > ? ORDER BY id LIMIT ?"
)

// GetUpCharges get up charges
func (d *Dao) GetUpCharges(c context.Context, table string, date string, offset, limit int64) (last int64, charges map[int64]*model.UpCharge, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_upChargeSQL, table), date, offset, limit)
	if err != nil {
		log.Error("d.db.Query GetUpCharges error(%v)", err)
		return
	}
	charges = make(map[int64]*model.UpCharge)
	defer rows.Close()
	for rows.Next() {
		c := &model.UpCharge{}
		err = rows.Scan(&last, &c.MID, &c.IncCharge, &c.TotalCharge, &c.Date)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		charges[c.MID] = c
	}
	return
}

// InsertUpCharge batch insert up charge
func (d *Dao) InsertUpCharge(c context.Context, table string, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insertUpChargeSQL, table, values))
	if err != nil {
		log.Error("d.db.Exec InsertUpCharge error(%v)", err)
		return
	}
	return res.RowsAffected()
}
