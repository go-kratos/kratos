package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/growup/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_columnIncomeByMIDSQL  = "SELECT aid, income, total_income, date FROM column_income WHERE mid = ? AND date >= ? AND date <= ?"
	_columnIncomeByAvIDSQL = "SELECT income, date FROM column_income WHERE aid = ? AND date <= ?"
	_columnStatisTitleSQL  = "SELECT aid, title FROM column_income_statis WHERE aid in (%s)"
)

// ListColumnIncome list column_income by mid
func (d *Dao) ListColumnIncome(c context.Context, mid int64, startTime, endTime string) (columns []*model.ArchiveIncome, err error) {
	columns = make([]*model.ArchiveIncome, 0)
	rows, err := d.db.Query(c, _columnIncomeByMIDSQL, mid, startTime, endTime)
	if err != nil {
		log.Error("ListColumnIncome d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		column := &model.ArchiveIncome{}
		err = rows.Scan(&column.ArchiveID, &column.Income, &column.TotalIncome, &column.Date)
		if err != nil {
			log.Error("ListColumnIncome rows.Scan error(%v)", err)
			return
		}
		columns = append(columns, column)
	}

	err = rows.Err()
	return
}

// ListColumnIncomeByID list column_income by aid
func (d *Dao) ListColumnIncomeByID(c context.Context, id int64, endTime string) (columns []*model.ArchiveIncome, err error) {
	columns = make([]*model.ArchiveIncome, 0)
	rows, err := d.db.Query(c, _columnIncomeByAvIDSQL, id, endTime)
	if err != nil {
		log.Error("ListColumnIncomeByID d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		column := &model.ArchiveIncome{}
		err = rows.Scan(&column.Income, &column.Date)
		if err != nil {
			log.Error("ListColumnIncomeByID rows.Scan error(%v)", err)
			return
		}
		columns = append(columns, column)
	}
	err = rows.Err()
	return
}

// GetColumnTitle get column title by id
func (d *Dao) GetColumnTitle(c context.Context, ids []int64) (titles map[int64]string, err error) {
	titles = make(map[int64]string)
	rows, err := d.db.Query(c, fmt.Sprintf(_columnStatisTitleSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("GetColumnTitle d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		var title string
		err = rows.Scan(&mid, &title)
		if err != nil {
			log.Error("GetColumnTitle rows.Scan error(%v)", err)
			return
		}
		titles[mid] = title
	}
	err = rows.Err()
	return
}
