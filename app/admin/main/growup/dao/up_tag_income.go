package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/growup/model"
	"go-common/library/log"
)

const (
	// select
	_upTagIncomeSQL = "SELECT id, mid, av_id, income, base_income, total_income, date, is_deleted FROM up_tag_income WHERE date = ? AND tag_id = ? %s"
)

// GetUpTagIncome get up_tag_income by query
func (d *Dao) GetUpTagIncome(c context.Context, date string, tagID int64, query string) (infos []*model.UpTagIncome, err error) {
	if query != "" {
		query = fmt.Sprintf("AND %s", query)
	}
	rows, err := d.rddb.Query(c, fmt.Sprintf(_upTagIncomeSQL, query), date, tagID)
	if err != nil {
		log.Error("d.rddb.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.UpTagIncome{}
		if err = rows.Scan(&a.ID, &a.MID, &a.AvID, &a.Income, &a.BaseIncome, &a.TotalIncome, &a.Date, &a.IsDeleted); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		infos = append(infos, a)
	}
	return
}
