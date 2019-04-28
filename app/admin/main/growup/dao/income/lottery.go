package income

import (
	"context"
	"fmt"

	model "go-common/app/admin/main/growup/model/income"
	"go-common/library/log"
)

const (
	// select
	_lotteryIncomeSQL = "SELECT id,av_id,mid,tag_id,upload_time,total_income,income,tax_money,date FROM lottery_av_income WHERE id > ? AND %s date >= ? AND date <= ? ORDER BY id LIMIT ?"
)

// GetLotteryIncome get lottery income by query
func (d *Dao) GetLotteryIncome(c context.Context, id int64, query string, from, to string, limit int, typ int) (avs []*model.ArchiveIncome, err error) {
	avs = make([]*model.ArchiveIncome, 0)
	if query != "" {
		query += " AND"
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_lotteryIncomeSQL, query), id, from, to, limit)
	if err != nil {
		log.Error("GetLotteryIncome d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.ArchiveIncome{}
		err = rows.Scan(&list.ID, &list.AvID, &list.MID, &list.TagID, &list.UploadTime, &list.TotalIncome, &list.Income, &list.TaxMoney, &list.Date)
		if err != nil {
			log.Error("GetLotteryIncome rows scan error(%v)", err)
			return
		}
		list.Type = typ
		avs = append(avs, list)
	}

	err = rows.Err()
	return
}
