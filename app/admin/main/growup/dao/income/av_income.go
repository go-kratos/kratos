package income

import (
	"context"
	"fmt"

	model "go-common/app/admin/main/growup/model/income"
	"go-common/library/log"
)

var (
	_video   = 0
	_column  = 2
	_bgm     = 3
	_lottery = 5
)

const (
	// select
	_avIncomeStatisTableSQL = "SELECT avs,money_section,money_tips,income,category_id,cdate FROM %s WHERE %s LIMIT ?,?"
	_avIncomeSQL            = "SELECT id,av_id,mid,tag_id,is_original,upload_time,total_income,income,tax_money,date FROM av_income WHERE id > ? AND %s date >= ? AND date <= ? ORDER BY id LIMIT ?"
	_columnIncomeSQL        = "SELECT id,aid,mid,tag_id,upload_time,total_income,income,tax_money,date FROM column_income WHERE id > ? AND date >= ? AND date <= ? AND %s is_deleted = 0 ORDER BY id LIMIT ?"
)

// GetArchiveStatis get av/column income statis from table and query
func (d *Dao) GetArchiveStatis(c context.Context, table, query string, from, limit int) (avs []*model.ArchiveStatis, err error) {
	avs = make([]*model.ArchiveStatis, 0)
	if table == "" || query == "" {
		return nil, fmt.Errorf("error args table(%s), query(%s)", table, query)
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_avIncomeStatisTableSQL, table, query), from, limit)
	if err != nil {
		log.Error("GetArchiveStatis d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.ArchiveStatis{}
		err = rows.Scan(&list.Avs, &list.MoneySection, &list.MoneyTips, &list.Income, &list.CategroyID, &list.CDate)
		if err != nil {
			log.Error("GetArchiveStatis rows scan error(%v)", err)
			return
		}
		avs = append(avs, list)
	}

	err = rows.Err()
	return
}

// GetArchiveIncome get archive income by query
func (d *Dao) GetArchiveIncome(c context.Context, id int64, query string, from, to string, limit int, typ int) (archs []*model.ArchiveIncome, err error) {
	switch typ {
	case _video:
		return d.GetAvIncome(c, id, query, from, to, limit, typ)
	case _column:
		return d.GetColumnIncome(c, id, query, from, to, limit, typ)
	case _bgm:
		return d.GetBgmIncome(c, id, query, from, to, limit, typ)
	case _lottery:
		return d.GetLotteryIncome(c, id, query, from, to, limit, typ)
	}
	err = fmt.Errorf("get archive type error(%d)", typ)
	return
}

// GetAvIncome get av income by query
func (d *Dao) GetAvIncome(c context.Context, id int64, query string, from, to string, limit int, typ int) (avs []*model.ArchiveIncome, err error) {
	avs = make([]*model.ArchiveIncome, 0)
	if query != "" {
		query += " AND"
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_avIncomeSQL, query), id, from, to, limit)
	if err != nil {
		log.Error("GetAvIncome d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.ArchiveIncome{}
		err = rows.Scan(&list.ID, &list.AvID, &list.MID, &list.TagID, &list.IsOriginal, &list.UploadTime, &list.TotalIncome, &list.Income, &list.TaxMoney, &list.Date)
		if err != nil {
			log.Error("GetAvIncome rows scan error(%v)", err)
			return
		}
		list.Type = typ
		avs = append(avs, list)
	}

	err = rows.Err()
	return
}

// GetColumnIncome get column income by query
func (d *Dao) GetColumnIncome(c context.Context, id int64, query string, from, to string, limit int, typ int) (columns []*model.ArchiveIncome, err error) {
	columns = make([]*model.ArchiveIncome, 0)
	if query != "" {
		query += " AND"
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_columnIncomeSQL, query), id, from, to, limit)
	if err != nil {
		log.Error("GetColumnIncome d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		list := &model.ArchiveIncome{}
		err = rows.Scan(&list.ID, &list.AvID, &list.MID, &list.TagID, &list.UploadTime, &list.TotalIncome, &list.Income, &list.TaxMoney, &list.Date)
		if err != nil {
			log.Error("GetColumnIncome rows scan error(%v)", err)
			return
		}
		list.Type = typ
		columns = append(columns, list)
	}

	err = rows.Err()
	return
}
