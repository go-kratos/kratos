package income

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/income"

	"go-common/library/log"
)

const (
	_getUpIncomeTableSQL = "SELECT id,mid,av_count,play_count,av_income,audio_income,column_count,column_income,bgm_income,av_tax,column_tax,bgm_tax,tax_money,income,total_income,av_base_income,column_base_income,bgm_base_income,base_income,av_total_income,column_total_income,bgm_total_income,date FROM %s WHERE id > ? AND date = ? ORDER BY id LIMIT ?"

	_inUpIncomeTableSQL = "INSERT INTO %s(mid,av_count,play_count,av_income,audio_income,column_count,column_income,bgm_count,bgm_income,tax_money,income,total_income,av_base_income,av_tax,column_base_income,column_tax,bgm_base_income,bgm_tax,date,base_income,av_total_income,column_total_income,bgm_total_income) VALUES %s ON DUPLICATE KEY UPDATE mid=VALUES(mid),av_count=VALUES(av_count),play_count=VALUES(play_count),av_income=VALUES(av_income),audio_income=VALUES(audio_income),column_count=VALUES(column_count),column_income=VALUES(column_income),bgm_count=VALUES(bgm_count),bgm_income=VALUES(bgm_income),tax_money=VALUES(tax_money),income=VALUES(income),total_income=VALUES(total_income),av_base_income=VALUES(av_base_income),av_tax=VALUES(av_tax),column_base_income=VALUES(column_base_income),column_tax=VALUES(column_tax),bgm_base_income=VALUES(bgm_base_income),bgm_tax=VALUES(bgm_tax),date=VALUES(date),base_income=VALUES(base_income),av_total_income=VALUES(av_total_income),column_total_income=VALUES(column_total_income),bgm_total_income=VALUES(bgm_total_income)"
)

// GetUpIncomeTable get up_income up_income_weekly up_income_monthly
func (d *Dao) GetUpIncomeTable(c context.Context, table, date string, id int64, limit int) (ups []*model.UpIncome, err error) {
	ups = make([]*model.UpIncome, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_getUpIncomeTableSQL, table), id, date, limit)
	if err != nil {
		log.Error("GetUpIncomeTable d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpIncome{}
		err = rows.Scan(&up.ID, &up.MID, &up.AvCount, &up.PlayCount, &up.AvIncome, &up.AudioIncome, &up.ColumnCount, &up.ColumnIncome, &up.BgmIncome, &up.AvTax, &up.ColumnTax, &up.BgmTax, &up.TaxMoney, &up.Income, &up.TotalIncome, &up.AvBaseIncome, &up.ColumnBaseIncome, &up.BgmBaseIncome, &up.BaseIncome, &up.AvTotalIncome, &up.ColumnTotalIncome, &up.BgmTotalIncome, &up.Date)
		if err != nil {
			log.Error("GetUpIncomeTable rows.Scan error(%v)", err)
			return
		}
		ups = append(ups, up)
	}
	return
}

// InsertUpIncomeTable insert up_income up_income_weekly up_income_monthly
func (d *Dao) InsertUpIncomeTable(c context.Context, table, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpIncomeTableSQL, table, values))
	if err != nil {
		log.Error("InsertUpIncome error (%v)", err)
		return
	}
	return res.RowsAffected()
}
