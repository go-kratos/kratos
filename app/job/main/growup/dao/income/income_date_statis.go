package income

import (
	"context"
	"fmt"

	"go-common/library/log"
)

const (
	_inIncomeStatisTableSQL = "INSERT INTO %s(avs, money_section, money_tips, income, category_id, cdate) VALUES %s ON DUPLICATE KEY UPDATE avs=VALUES(avs),income=VALUES(income),cdate=VALUES(cdate)"

	_delIncomeStatisTableSQL = "DELETE FROM %s WHERE cdate = ?"

	_inUpIncomeDailyStatisSQL = "INSERT INTO %s(ups, money_section, money_tips, income, cdate) VALUES %s ON DUPLICATE KEY UPDATE ups=VALUES(ups),income=VALUES(income),cdate=VALUES(cdate)"
)

// InsertIncomeStatisTable add av_income_date_statis batch
func (d *Dao) InsertIncomeStatisTable(c context.Context, table, vals string) (rows int64, err error) {
	if table == "" {
		err = fmt.Errorf("InsertIncomeStatisTable table(%s) val(%s) error", table, vals)
		return
	}
	if vals == "" {
		return
	}

	res, err := d.db.Exec(c, fmt.Sprintf(_inIncomeStatisTableSQL, table, vals))
	if err != nil {
		log.Error("incomeStatisTableBatch d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertUpIncomeDailyStatis add up_income_daily_statis batch
func (d *Dao) InsertUpIncomeDailyStatis(c context.Context, table string, vals string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpIncomeDailyStatisSQL, table, vals))
	if err != nil {
		log.Error("InsertUpIncomeDailyStatis d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// DelIncomeStatisTable del income statis by table
func (d *Dao) DelIncomeStatisTable(c context.Context, table, date string) (rows int64, err error) {
	if table == "" || date == "" {
		err = fmt.Errorf("DelIncomeStatisTable table(%s) date(%s) error", table, date)
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_delIncomeStatisTableSQL, table), date)
	if err != nil {
		log.Error("DelIncomeStatisTable d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
