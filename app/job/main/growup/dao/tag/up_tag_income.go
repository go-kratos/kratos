package tag

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/app/job/main/growup/model/tag"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// select
	_upTagIncomeSQL   = "SELECT id, mid, tag_id FROM up_tag_income WHERE id > ? AND date >= '%s' ORDER BY id LIMIT ?"
	_upTagIncomeAvSQL = "SELECT av_id, total_income FROM up_tag_income WHERE id > ? ORDER BY id LIMIT ?"
	_tagAvInfoSQL     = "SELECT id, tag_id, income, av_id, is_deleted FROM up_tag_income WHERE date = ? AND id > ? ORDER BY id LIMIT ?"
	_tagUpByID        = "SELECT id, mid FROM up_tag_income WHERE id > ? AND tag_id = ? ORDER BY id LIMIT ?"

	// insert
	_insertUpTagIncomeSQL = "INSERT INTO up_tag_income(tag_id,mid,av_id,income,base_income,total_income,tax_money,date) VALUES %s ON DUPLICATE KEY UPDATE income = VALUES(income),total_income = VALUES(total_income),tax_money=VALUES(tax_money)"
)

// UpTagIncomeByDate get up tag income
func (d *Dao) UpTagIncomeByDate(c context.Context, date string, id int64, limit int) (vals []*tag.UpTagIncome, err error) {
	vals = make([]*tag.UpTagIncome, 0, limit)
	rows, err := d.db.Query(c, fmt.Sprintf(_upTagIncomeSQL, date), id, limit)
	if err != nil {
		log.Error("d.db.Query(%v), error(%v)", _upTagIncomeSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		val := &tag.UpTagIncome{}
		err = rows.Scan(&val.ID, &val.MID, &val.TagID)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		vals = append(vals, val)
	}
	return
}

// GetUpTagIncomeMap get up_tag_income map
func (d *Dao) GetUpTagIncomeMap(c context.Context, id int64, limit int, avs map[int64]int64) (lastID int64, count int, err error) {
	rows, err := d.db.Query(c, _upTagIncomeAvSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query(%v), error(%v)", _upTagIncomeAvSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var avID, totalIncome int64
		err = rows.Scan(&avID, &totalIncome)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		count++
		avs[avID] = totalIncome
		lastID = avID
	}
	return
}

// GetTagUpByID get up tag
func (d *Dao) GetTagUpByID(c context.Context, tagID int64, id int64, limit int) (vals []*tag.UpTagIncome, err error) {
	vals = make([]*tag.UpTagIncome, 0, limit)
	rows, err := d.db.Query(c, _tagUpByID, id, tagID, limit)
	if err != nil {
		log.Error("d.db.Query(%v), error(%v)", _tagUpByID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		val := &tag.UpTagIncome{}
		err = rows.Scan(&val.ID, &val.MID)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		vals = append(vals, val)
	}
	return
}

// TxInsertUpTagIncome insert up_tag_income.
func (d *Dao) TxInsertUpTagIncome(tx *sql.Tx, vals string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := tx.Exec(fmt.Sprintf(_insertUpTagIncomeSQL, vals))
	if err != nil {
		log.Error("dao.TxInsertUpTagIncome exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// GetTagAvInfo get tag av infos by date.
func (d *Dao) GetTagAvInfo(c context.Context, date time.Time, from, limit int64) (infos []*model.AvIncome, err error) {
	rows, err := d.db.Query(c, _tagAvInfoSQL, date, from, limit)
	if err != nil {
		log.Error("dao.GetTagAvInfo query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.AvIncome{}
		if err = rows.Scan(&a.ID, &a.TagID, &a.Income, &a.AvID, &a.IsDeleted); err != nil {
			log.Error("dao.GetTagAvInfo scan error(%v)", err)
			return
		}
		infos = append(infos, a)
	}
	err = rows.Err()
	return
}
