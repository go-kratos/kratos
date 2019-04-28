package income

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/income"

	"go-common/library/log"
)

const (
	_inUpIncomeSQL = "INSERT INTO up_income(mid,av_count,play_count,av_income,audio_income,column_count,column_income,tax_money,income,total_income,av_base_income,av_tax,column_base_income,column_tax,bgm_base_income,bgm_tax,date,base_income,bgm_income,bgm_count,av_total_income,column_total_income,bgm_total_income) VALUES %s ON DUPLICATE KEY UPDATE mid=VALUES(mid),av_count=VALUES(av_count),play_count=VALUES(play_count),av_income=VALUES(av_income),audio_income=VALUES(audio_income),column_count=VALUES(column_count),column_income=VALUES(column_income),tax_money=VALUES(tax_money),income=VALUES(income),total_income=VALUES(total_income),av_base_income=VALUES(av_base_income),av_tax=VALUES(av_tax),column_base_income=VALUES(column_base_income),column_tax=VALUES(column_tax),bgm_base_income=VALUES(bgm_base_income),bgm_tax=VALUES(bgm_tax),date=VALUES(date),base_income=VALUES(base_income),bgm_income=VALUES(bgm_income),bgm_count=VALUES(bgm_count),av_total_income=VALUES(av_total_income),column_total_income=VALUES(column_total_income),bgm_total_income=VALUES(bgm_total_income)"

	_upIncomeSQL      = "SELECT id,mid,av_income,audio_income,column_income,bgm_income,date FROM up_income WHERE id > ? ORDER BY id LIMIT ?"
	_fixInUpIncomeSQL = "INSERT INTO up_income(mid,av_total_income,column_total_income,bgm_total_income,date) VALUES %s ON DUPLICATE KEY UPDATE av_total_income=VALUES(av_total_income),column_total_income=VALUES(column_total_income),bgm_total_income=VALUES(bgm_total_income)"
)

// UpIncomes up incomes
func (d *Dao) UpIncomes(c context.Context, id int64, limit int64) (last int64, us []*model.UpIncome, err error) {
	rows, err := d.db.Query(c, _upIncomeSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query UpIncomes error (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &model.UpIncome{}
		err = rows.Scan(&last, &u.MID, &u.AvIncome, &u.AudioIncome, &u.ColumnIncome, &u.BgmIncome, &u.Date)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		us = append(us, u)
	}
	return
}

// InsertUpIncome batch insert up income
func (d *Dao) InsertUpIncome(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpIncomeSQL, values))
	if err != nil {
		log.Error("d.db.Exec InsertUpIncome error (%v)", err)
		return
	}
	return res.RowsAffected()
}

// FixInsertUpIncome batch insert up income
func (d *Dao) FixInsertUpIncome(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_fixInUpIncomeSQL, values))
	if err != nil {
		log.Error("d.db.Exec InsertUpIncome error (%v)", err)
		return
	}
	return res.RowsAffected()
}
