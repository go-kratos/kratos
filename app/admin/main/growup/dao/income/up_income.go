package income

import (
	"context"
	"fmt"

	model "go-common/app/admin/main/growup/model/income"
	"go-common/library/log"
)

const (
	// select
	_upIncomeTableSQL     = "SELECT id,mid,%s,date FROM %s WHERE id > ? AND %s LIMIT ?"
	_upIncomeTableSortSQL = "SELECT mid,av_count,column_count,bgm_count,%s,date FROM %s WHERE %s ORDER BY date desc,%s desc LIMIT ?,? "
	_upIncomeCountSQL     = "SELECT count(*) FROM %s WHERE %s"
	_upDailyStatisSQL     = "SELECT ups,income,cdate FROM %s WHERE cdate >= '%s' AND cdate <= '%s'"
)

// UpIncomeCount count
func (d *Dao) UpIncomeCount(c context.Context, table, query string) (count int, err error) {
	if table == "" || query == "" {
		return 0, fmt.Errorf("error args table(%s), query(%s)", table, query)
	}
	err = d.db.QueryRow(c, fmt.Sprintf(_upIncomeCountSQL, table, query)).Scan(&count)
	return
}

// GetUpIncome get up_income_(weekly/monthly) from table and query
func (d *Dao) GetUpIncome(c context.Context, table, incomeType, query string, id int64, limit int) (upIncome []*model.UpIncome, err error) {
	upIncome = make([]*model.UpIncome, 0)
	if table == "" || query == "" {
		return nil, fmt.Errorf("error args table(%s), query(%s)", table, query)
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_upIncomeTableSQL, incomeType, table, query), id, limit)
	if err != nil {
		log.Error("GetUpIncome d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.UpIncome{}
		err = rows.Scan(&list.ID, &list.MID, &list.Income, &list.Date)
		if err != nil {
			log.Error("GetUpIncome rows scan error(%v)", err)
			return
		}
		upIncome = append(upIncome, list)
	}

	err = rows.Err()
	return
}

// GetUpIncomeBySort get up_income by query
func (d *Dao) GetUpIncomeBySort(c context.Context, table, typeField, sort, query string, from, limit int) (upIncome []*model.UpIncome, err error) {
	upIncome = make([]*model.UpIncome, 0)
	if table == "" || query == "" || typeField == "" {
		return nil, fmt.Errorf("error args table(%s), typeField(%s),query(%s)", table, typeField, query)
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_upIncomeTableSortSQL, typeField, table, query, sort), from, limit)
	if err != nil {
		log.Error("GetUpIncomeBySort d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpIncome{}
		err = rows.Scan(&up.MID, &up.AvCount, &up.ColumnCount, &up.BgmCount, &up.Income, &up.TaxMoney, &up.BaseIncome, &up.TotalIncome, &up.Date)
		if err != nil {
			log.Error("GetUpIncome rows scan error(%v)", err)
			return
		}
		upIncome = append(upIncome, up)
	}

	err = rows.Err()
	return
}

// GetUpDailyStatis get up income daily statis
func (d *Dao) GetUpDailyStatis(c context.Context, table, fromTime, toTime string) (s []*model.UpDailyStatis, err error) {
	if table == "" {
		return nil, fmt.Errorf("error args table(%s)", table)
	}
	s = make([]*model.UpDailyStatis, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_upDailyStatisSQL, table, fromTime, toTime))
	if err != nil {
		log.Error("GetUpDailyStatis d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.UpDailyStatis{}
		err = rows.Scan(&list.Ups, &list.Income, &list.Date)
		if err != nil {
			log.Error("GetUpIncome rows scan error(%v)", err)
			return
		}
		s = append(s, list)
	}
	err = rows.Err()
	return
}
