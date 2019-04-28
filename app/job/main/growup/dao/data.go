package dao

import (
	"context"
	"fmt"

	"go-common/library/database/sql"
	"go-common/library/xstr"

	"go-common/app/job/main/growup/model"
)

const (
	_avIDs            = "SELECT av_id,mid,tag_id,income,tax_money FROM av_income WHERE date=? AND mid=?"
	_avCharges        = "SELECT av_id,inc_charge FROM av_daily_charge_04 WHERE av_id IN (%s) AND date=?"
	_upChargeRatio    = "SELECT mid,ratio FROM up_charge_ratio WHERE tag_id = ?"
	_upIncomeStatis   = "SELECT mid,total_income FROM up_income_statis WHERE mid in (%s)"
	_upIncomeDate     = "SELECT mid, income FROM %s WHERE date = ? AND mid in (%s)"
	_upIncome         = "SELECT id,mid,date,base_income FROM %s WHERE id > ? ORDER BY id LIMIT ?"
	_inUpIncome       = "INSERT INTO %s(mid,date,av_base_income) VALUES %s ON DUPLICATE KEY UPDATE av_base_income=VALUES(av_base_income)"
	_creditScoreSQL   = "SELECT id,mid,credit_score FROM up_info_%s WHERE id > ? ORDER BY id LIMIT ?"
	_inCreditScoreSQL = "INSERT INTO credit_score(mid,score) VALUES %s ON DUPLICATE KEY UPDATE score=VALUES(score)"
	_bgmIncome        = "SELECT sid,income FROM bgm_income"
	_bgmIncomeStatis  = "INSERT INTO bgm_income_statis(sid,total_income) VALUES(?,?)"
)

// InsertBGMIncomeStatis fix bgm income statis
func (d *Dao) InsertBGMIncomeStatis(c context.Context, sid int64, income int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _bgmIncomeStatis, sid, income)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// GetBGMIncome  map[sid]totalIncome
func (d *Dao) GetBGMIncome(c context.Context) (statis map[int64]int64, err error) {
	rows, err := d.db.Query(c, _bgmIncome)
	if err != nil {
		return
	}
	statis = make(map[int64]int64)
	defer rows.Close()
	for rows.Next() {
		var sid, income int64
		err = rows.Scan(&sid, &income)
		if err != nil {
			return
		}
		if _, ok := statis[sid]; ok {
			statis[sid] += income
		} else {
			statis[sid] = income
		}
	}
	return
}

// GetCreditScore get credit scores
func (d *Dao) GetCreditScore(c context.Context, table string, id int64, limit int64) (scores map[int64]int, last int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_creditScoreSQL, table), id, limit)
	if err != nil {
		return
	}
	scores = make(map[int64]int)
	defer rows.Close()
	for rows.Next() {
		var mid int64
		var score int
		err = rows.Scan(&last, &mid, &score)
		if err != nil {
			return
		}
		scores[mid] = score
	}
	return
}

// SyncCreditScore sync credit score
func (d *Dao) SyncCreditScore(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inCreditScoreSQL, values))
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// GetAvBaseIncome get up av baseincome
func (d *Dao) GetAvBaseIncome(c context.Context, table string, id, limit int64) (abs []*model.AvBaseIncome, last int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_upIncome, table), id, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		ab := &model.AvBaseIncome{}
		err = rows.Scan(&last, &ab.MID, &ab.Date, &ab.AvBaseIncome)
		if err != nil {
			return
		}
		abs = append(abs, ab)
	}
	return
}

// BatchUpdateUpIncome batch update av_base_income
func (d *Dao) BatchUpdateUpIncome(c context.Context, table, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpIncome, table, values))
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// GetAvs avs map[av_id]*model.Av
func (d *Dao) GetAvs(c context.Context, date string, mid int64) (avs map[int64]*model.Av, err error) {
	rows, err := d.db.Query(c, _avIDs, date, mid)
	if err != nil {
		return
	}
	avs = make(map[int64]*model.Av)
	defer rows.Close()
	for rows.Next() {
		av := &model.Av{}
		err = rows.Scan(&av.AvID, &av.MID, &av.TagID, &av.Income, &av.TaxMoney)
		if err != nil {
			return
		}
		avs[av.AvID] = av
	}
	return
}

// GetAvCharges get av charges
func (d *Dao) GetAvCharges(c context.Context, avIds []int64, date string) (charges map[int64]int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_avCharges, xstr.JoinInts(avIds)), date)
	if err != nil {
		return
	}
	charges = make(map[int64]int64)
	defer rows.Close()
	for rows.Next() {
		var avID, charge int64
		err = rows.Scan(&avID, &charge)
		if err != nil {
			return
		}
		charges[avID] = charge
	}
	return
}

// GetUpChargeRatio get up_charge_ratio
func (d *Dao) GetUpChargeRatio(c context.Context, tagID int64) (ups map[int64]int64, err error) {
	rows, err := d.db.Query(c, _upChargeRatio, tagID)
	if err != nil {
		return
	}
	defer rows.Close()
	ups = make(map[int64]int64)
	for rows.Next() {
		var mid, ratio int64
		err = rows.Scan(&mid, &ratio)
		if err != nil {
			return
		}
		ups[mid] = ratio
	}
	return
}

// GetUpIncomeStatis get up_income_statis
func (d *Dao) GetUpIncomeStatis(c context.Context, mids []int64) (ups map[int64]int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_upIncomeStatis, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	defer rows.Close()
	ups = make(map[int64]int64)
	for rows.Next() {
		var mid, income int64
		err = rows.Scan(&mid, &income)
		if err != nil {
			return
		}
		ups[mid] = income
	}
	return
}

// GetUpIncomeDate get up_income by date
func (d *Dao) GetUpIncomeDate(c context.Context, mids []int64, table, date string) (ups map[int64]int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_upIncomeDate, table, xstr.JoinInts(mids)), date)
	if err != nil {
		return
	}
	defer rows.Close()
	ups = make(map[int64]int64)
	for rows.Next() {
		var mid, income int64
		err = rows.Scan(&mid, &income)
		if err != nil {
			return
		}
		ups[mid] = income
	}
	return
}

// UpdateDate update date
func (d *Dao) UpdateDate(tx *sql.Tx, stmt string) (count int64, err error) {
	res, err := tx.Exec(stmt)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
