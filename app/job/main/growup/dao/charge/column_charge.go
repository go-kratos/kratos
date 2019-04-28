package charge

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/job/main/growup/model/charge"

	"go-common/library/log"
)

const (
	_columnChargeSQL = "SELECT id,aid,title,mid,tag_id,words,upload_time,inc_charge,view_c,date FROM %s WHERE id > ? AND date = ? AND inc_charge > 0 ORDER BY id LIMIT ?"
	_columnStatisSQL = "SELECT id,aid,total_charge FROM column_charge_statis WHERE id > ? ORDER BY id LIMIT ?"
	_countCmDailySQL = "SELECT COUNT(*) FROM column_daily_charge WHERE date = '%s'"

	_inCmChargeTableSQL = "INSERT INTO %s(aid,mid,tag_id,inc_charge,date,upload_time) VALUES %s ON DUPLICATE KEY UPDATE inc_charge=VALUES(inc_charge)"
	_inCmStatisSQL      = "INSERT INTO column_charge_statis(aid,mid,tag_id,total_charge,upload_time) VALUES %s ON DUPLICATE KEY UPDATE total_charge=VALUES(total_charge)"
)

// ColumnCharge get column charge by date
func (d *Dao) ColumnCharge(c context.Context, date time.Time, id int64, limit int, table string) (columns []*model.Column, err error) {
	columns = make([]*model.Column, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_columnChargeSQL, table), id, date, limit)
	if err != nil {
		log.Error("ColumnCharge d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		column := &model.Column{}
		err = rows.Scan(&column.ID, &column.AID, &column.Title, &column.MID, &column.TagID, &column.Words, &column.UploadTime, &column.IncCharge, &column.IncViewCount, &column.Date)
		if err != nil {
			log.Error("ColumnCharge rows.Scan error(%v)", err)
			return
		}
		columns = append(columns, column)
	}
	return
}

// CmStatis column statis
func (d *Dao) CmStatis(c context.Context, id int64, limit int) (columns []*model.ColumnStatis, err error) {
	columns = make([]*model.ColumnStatis, 0)
	rows, err := d.db.Query(c, _columnStatisSQL, id, limit)
	if err != nil {
		log.Error("CmStatis d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		column := &model.ColumnStatis{}
		err = rows.Scan(&column.ID, &column.AID, &column.TotalCharge)
		if err != nil {
			log.Error("CmStatis rows.Scan error(%v)", err)
			return
		}
		columns = append(columns, column)
	}
	return
}

// InsertCmChargeTable insert column charge
func (d *Dao) InsertCmChargeTable(c context.Context, vals, table string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inCmChargeTableSQL, table, vals))
	if err != nil {
		log.Error("InsertCmChargeTable(%s) tx.Exec error(%v)", table, err)
		return
	}
	return res.RowsAffected()
}

// InsertCmStatisBatch insert column statis
func (d *Dao) InsertCmStatisBatch(c context.Context, vals string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inCmStatisSQL, vals))
	if err != nil {
		log.Error("InsertCmStatisBatch tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// CountCmDailyCharge get column_daily_charge count
func (d *Dao) CountCmDailyCharge(c context.Context, date string) (count int64, err error) {
	err = d.db.QueryRow(c, fmt.Sprintf(_countCmDailySQL, date)).Scan(&count)
	return
}
