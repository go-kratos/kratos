package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// credit score record
	_inCreditRecordSQL = "INSERT INTO credit_score_record (mid,operate_at,operator,reason,deducted,remaining) VALUES(?,?,?,?,?,?)"
	_creditRecordsSQL  = "SELECT id,mid,operate_at,operator,reason,deducted,remaining,is_deleted FROM credit_score_record WHERE mid=?"
	_creditRecordSQL   = "SELECT deducted FROM credit_score_record WHERE id=? AND is_deleted=0"

	// credit score
	_inCreditScoreSQL = "INSERT INTO credit_score (mid) VALUES %s ON DUPLICATE KEY UPDATE mid=VALUES(mid)"
	_creditScoreSQL   = "SELECT score FROM credit_score WHERE mid=?"
	_creditScoresSQL  = "SELECT mid,score FROM credit_score WHERE mid IN (%s)"

	// update credit score
	_upCreditScoreSQL = "UPDATE credit_score SET score=? WHERE mid=?"
	_reCreditScoreSQL = "UPDATE credit_score SET score=score+%d WHERE mid=?"
)

// InsertCreditRecord insert credit record
func (d *Dao) InsertCreditRecord(c context.Context, cr *model.CreditRecord) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _inCreditRecordSQL, cr.MID, cr.OperateAt, cr.Operator, cr.Reason, cr.Deducted, cr.Remaining)
	if err != nil {
		log.Error("db.inCreditRecordSQL.Exec(%s) error(%v)", _inCreditRecordSQL, err)
		return
	}
	return res.RowsAffected()
}

// TxInsertCreditRecord tx insert credit record
func (d *Dao) TxInsertCreditRecord(tx *sql.Tx, cr *model.CreditRecord) (rows int64, err error) {
	res, err := tx.Exec(_inCreditRecordSQL, cr.MID, cr.OperateAt, cr.Operator, cr.Reason, cr.Deducted, cr.Remaining)
	if err != nil {
		log.Error("tx.inCreditRecordSQL.Exec(%s) error(%v)", _inCreditRecordSQL, err)
		return
	}
	return res.RowsAffected()
}

// CreditRecords get credit records by mid
func (d *Dao) CreditRecords(c context.Context, mid int64) (crs []*model.CreditRecord, err error) {
	rows, err := d.rddb.Query(c, _creditRecordsSQL, mid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		cr := &model.CreditRecord{}
		err = rows.Scan(&cr.ID, &cr.MID, &cr.OperateAt, &cr.Operator, &cr.Reason, &cr.Deducted, &cr.Remaining, &cr.IsDeleted)
		if err != nil {
			log.Error("rows Scan error(%v)", err)
			return
		}
		crs = append(crs, cr)
	}
	return
}

// DeductedScore get deducted credit score from credit_score_record by id
func (d *Dao) DeductedScore(c context.Context, id int64) (deducted int, err error) {
	row := d.rddb.QueryRow(c, _creditRecordSQL, id)
	if err = row.Scan(&deducted); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// InsertCreditScore insert credit score
func (d *Dao) InsertCreditScore(c context.Context, values string) (rows int64, err error) {
	res, err := d.rddb.Exec(c, fmt.Sprintf(_inCreditScoreSQL, values))
	if err != nil {
		log.Error("db.inCreditScoreSQL.Exec(%s) error(%v)", _inCreditScoreSQL, err)
		return
	}
	return res.RowsAffected()
}

// UpdateCreditScore update credit score
func (d *Dao) UpdateCreditScore(c context.Context, mid int64, score int) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _upCreditScoreSQL, score, mid)
	if err != nil {
		log.Error("db.upCreditScoreSQL.Exec(%s) error(%v)", _inCreditScoreSQL, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateCreditScore update credit score
func (d *Dao) TxUpdateCreditScore(tx *sql.Tx, mid int64, score int) (rows int64, err error) {
	res, err := tx.Exec(_upCreditScoreSQL, score, mid)
	if err != nil {
		log.Error("tx.upCreditScoreSQL.Exec(%s) error(%v)", _inCreditScoreSQL, err)
		return
	}
	return res.RowsAffected()
}

// CreditScore credit score by mid
func (d *Dao) CreditScore(c context.Context, mid int64) (score int, err error) {
	row := d.rddb.QueryRow(c, _creditScoreSQL, mid)
	if err = row.Scan(&score); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// CreditScores scores map[mid]score
func (d *Dao) CreditScores(c context.Context, mids []int64) (scores map[int64]int, err error) {
	scores = make(map[int64]int)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_creditScoresSQL, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		var score int
		err = rows.Scan(&mid, &score)
		if err != nil {
			log.Error("rows Scan error(%v)", err)
			return
		}
		scores[mid] = score
	}
	return
}

// TxRecoverCreditScore recover credit score
func (d *Dao) TxRecoverCreditScore(tx *sql.Tx, deducted int, mid int64) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_reCreditScoreSQL, deducted), mid)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
