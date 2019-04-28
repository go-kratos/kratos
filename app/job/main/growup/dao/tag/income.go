package tag

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/tag"
	"go-common/library/log"
)

const (
	// select
	_avIncomeSQL       = "SELECT id,av_id,mid,income,base_income,total_income,tax_money FROM av_income WHERE id > ? %s ORDER BY id LIMIT ?"
	_cmIncomeSQL       = "SELECT id,aid,mid,income,base_income,total_income,tax_money FROM column_income WHERE id > ? %s ORDER BY id LIMIT ?"
	_bgmIncomeSQL      = "SELECT id,sid,mid,income,base_income,total_income,tax_money FROM bgm_income WHERE id > ? %s ORDER BY id LIMIT ?"
	_upIncomeSQL       = "SELECT id,mid,income,base_income,total_income,av_income,av_base_income,av_total_income,column_income,column_base_income,column_total_income,bgm_income,bgm_base_income,bgm_total_income,tax_money,av_tax,column_tax,bgm_tax FROM up_income WHERE id > ? %s ORDER BY id LIMIT ?"
	_avIncomeStatisSQL = "SELECT id,av_id,mid,tag_id,upload_time FROM av_income_statis WHERE id > ? ORDER BY id LIMIT ?"
	_cmIncomeStatisSQL = "SELECT id,aid,mid,tag_id,upload_time FROM column_income_statis WHERE id > ? ORDER BY id LIMIT ?"
)

var (
	_video  = 1
	_column = 2
	_bgm    = 3
)

// GetArchiveIncome get archive income
func (d *Dao) GetArchiveIncome(c context.Context, id int64, query string, limit int, ctype int) (archives []*model.ArchiveIncome, err error) {
	if query != "" {
		query = "AND " + query
	}
	sql := ""
	switch ctype {
	case _video:
		sql = _avIncomeSQL
	case _column:
		sql = _cmIncomeSQL
	case _bgm:
		sql = _bgmIncomeSQL
	}

	rows, err := d.db.Query(c, fmt.Sprintf(sql, query), id, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		archive := &model.ArchiveIncome{}
		err = rows.Scan(&archive.ID, &archive.AID, &archive.MID, &archive.Income, &archive.BaseIncome, &archive.TotalIncome, &archive.TaxMoney)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		archives = append(archives, archive)
	}
	return
}

// GetUpIncome get up_income by query
func (d *Dao) GetUpIncome(c context.Context, id int64, query string, limit int) (ups []*model.UpIncome, err error) {
	if query != "" {
		query = "AND " + query
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_upIncomeSQL, query), id, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.UpIncome{}
		err = rows.Scan(&up.ID, &up.MID, &up.Income, &up.BaseIncome, &up.TotalIncome, &up.AvIncome, &up.AvBaseIncome, &up.AvTotalIncome, &up.ColumnIncome, &up.ColumnBaseIncome, &up.ColumnTotalIncome, &up.BgmIncome, &up.BgmBaseIncome, &up.BgmTotalIncome, &up.TaxMoney, &up.AvTax, &up.ColumnTax, &up.BgmTax)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		ups = append(ups, up)
	}
	return
}

// GetAvIncomeStatis get av_income_statis
func (d *Dao) GetAvIncomeStatis(c context.Context, id int64, limit int) (avs []*model.ArchiveCharge, err error) {
	rows, err := d.db.Query(c, _avIncomeStatisSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		av := &model.ArchiveCharge{}
		err = rows.Scan(&av.ID, &av.AID, &av.MID, &av.CategoryID, &av.UploadTime)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		avs = append(avs, av)
	}
	return
}

// GetCmIncomeStatis get av_income_statis
func (d *Dao) GetCmIncomeStatis(c context.Context, id int64, limit int) (cms []*model.ArchiveCharge, err error) {
	rows, err := d.db.Query(c, _cmIncomeStatisSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		cm := &model.ArchiveCharge{}
		err = rows.Scan(&cm.ID, &cm.AID, &cm.MID, &cm.CategoryID, &cm.UploadTime)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		cms = append(cms, cm)
	}
	return
}
