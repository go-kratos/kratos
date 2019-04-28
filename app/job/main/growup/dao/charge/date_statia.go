package charge

import (
	"context"
	"fmt"

	"go-common/library/log"
)

const (
	_inStatisTableSQL = "INSERT INTO %s(avs, money_section, money_tips, charge, category_id, cdate) VALUES %s ON DUPLICATE KEY UPDATE avs=VALUES(avs),charge=VALUES(charge),cdate=VALUES(cdate)"

	_delStatisTableSQL = "DELETE FROM %s WHERE cdate = ?"
)

// InsertStatisTable add archive_charge_date_statis batch
func (d *Dao) InsertStatisTable(c context.Context, table, vals string) (rows int64, err error) {
	if table == "" {
		err = fmt.Errorf("InsertStatisTable table(%s) val(%s) error", table, vals)
		return
	}
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inStatisTableSQL, table, vals))
	if err != nil {
		log.Error("InsertStatisTable d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// DelStatisTable delete av_charge_statis
func (d *Dao) DelStatisTable(c context.Context, table, date string) (rows int64, err error) {
	if table == "" || date == "" {
		err = fmt.Errorf("DelStatisTable table(%s) date(%s) error", table, date)
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_delStatisTableSQL, table), date)
	if err != nil {
		log.Error("DelStatisTable d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
