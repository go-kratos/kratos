package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/job/main/mcn/model"
	xsql "go-common/library/database/sql"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_upMcnSignStateOPSQL       = "UPDATE mcn_sign SET state = ? WHERE id = ?"
	_upMcnUpStateOPSQL         = "UPDATE mcn_up SET state = ?, state_change_time = ? WHERE id = ?"
	_upMcnSignPayExpOPSQL      = "UPDATE mcn_sign SET pay_expire_state = 2 WHERE id = ?"
	_upMcnSignEmailStateSQL    = "UPDATE mcn_sign SET email_state = 2 WHERE id IN (%s)"
	_upMcnSignPayEmailStateSQL = "UPDATE mcn_sign_pay SET email_state = 2 WHERE id IN (%s)"
	_inMcnDataSummarySQL       = "INSERT mcn_data_summary(mcn_mid,sign_id,up_count,fans_count_accumulate,generate_date,data_type) VALUES (?,?,?,?,?,1)"
	_selMcnSignsSQL            = `SELECT id,begin_date,end_date,state FROM mcn_sign`
	_selMcnUpsSQL              = `SELECT id,begin_date,end_date,state FROM mcn_up LIMIT ?,?`
	_selMcnSignPayWarnsSQL     = `SELECT p.sign_id,p.due_date,p.pay_value FROM mcn_sign_pay p INNER JOIN mcn_sign s ON p.sign_id = s.id WHERE p.state = 0 AND 
	s.state = 10 AND s.end_date >= ? AND s.begin_date <= p.due_date AND p.due_date <= s.end_date AND date_sub(p.due_date,interval 7 day) <= ?`
	_selMcnSignMidsSQL  = "SELECT id,mcn_mid FROM mcn_sign WHERE state = 10"
	_selMcnUPCountSQL   = "SELECT sign_id,count(up_mid) as count FROM mcn_up WHERE sign_id IN (%s) AND state = 10 GROUP BY sign_id"
	_selMcnUPMidsSQL    = "SELECT sign_id,up_mid FROM mcn_up WHERE sign_id IN (%s) AND state = 10"
	_selCrmUpMidsSumSQL = "SELECT SUM(fans_count) as count FROM up_base_info WHERE mid IN (%s)"
	_selMcnSignPayDues  = `SELECT p.id, p.mcn_mid, p.sign_id, p.due_date, p.pay_value FROM mcn_sign_pay p LEFT JOIN mcn_sign s ON p.sign_id = s.id 
	WHERE p.due_date <= ? AND p.email_state = 1 AND p.state = 0 AND s.state = 10 AND s.end_date >= ?`
	_selMcnSignDues = "SELECT id, mcn_mid, begin_date, end_date FROM mcn_sign WHERE end_date <= ? and end_date >= ? and email_state = 1"
)

// UpMcnSignStateOP .
func (d *Dao) UpMcnSignStateOP(c context.Context, signID int64, state int8) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMcnSignStateOPSQL, state, signID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMcnUpStateOP .
func (d *Dao) UpMcnUpStateOP(c context.Context, signUpID int64, state int8) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMcnUpStateOPSQL, state, time.Now(), signUpID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMcnSignPayExpOP .
func (d *Dao) UpMcnSignPayExpOP(c context.Context, signPayID int64) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _upMcnSignPayExpOPSQL, signPayID); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMcnSignPayEmailState .
func (d *Dao) UpMcnSignPayEmailState(c context.Context, ids []int64) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_upMcnSignPayEmailStateSQL, xstr.JoinInts(ids))); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// UpMcnSignEmailState .
func (d *Dao) UpMcnSignEmailState(c context.Context, ids []int64) (rows int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_upMcnSignEmailStateSQL, xstr.JoinInts(ids))); err != nil {
		return rows, err
	}
	return res.RowsAffected()
}

// AddMcnDataSummary .
func (d *Dao) AddMcnDataSummary(c context.Context, mcnMid, signID, upCount, fansCountAccumulate int64, genDate xtime.Time) (err error) {
	_, err = d.db.Exec(c, _inMcnDataSummarySQL, mcnMid, signID, upCount, fansCountAccumulate, genDate)
	return
}

// McnSigns .
func (d *Dao) McnSigns(c context.Context) (mss []*model.MCNSignInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selMcnSignsSQL); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		ms := new(model.MCNSignInfo)
		if err = rows.Scan(&ms.SignID, &ms.BeginDate, &ms.EndDate, &ms.State); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		mss = append(mss, ms)
	}
	err = rows.Err()
	return
}

// McnUps .
func (d *Dao) McnUps(c context.Context, offset, limit int64) (ups []*model.MCNUPInfo, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selMcnUpsSQL, offset, limit); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := new(model.MCNUPInfo)
		if err = rows.Scan(&up.SignUpID, &up.BeginDate, &up.EndDate, &up.State); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		ups = append(ups, up)
	}
	err = rows.Err()
	return
}

// McnSignPayWarns .
func (d *Dao) McnSignPayWarns(c context.Context) (sps []*model.SignPayInfo, err error) {
	var (
		rows     *xsql.Rows
		now      time.Time
		template = time.Now().Format(model.TimeFormatDay)
	)
	if now, err = time.ParseInLocation(model.TimeFormatDay, template, time.Local); err != nil {
		err = errors.Errorf("time.ParseInLocation(%s) error(%+v)", template, err)
		return
	}
	if rows, err = d.db.Query(c, _selMcnSignPayWarnsSQL, now, now); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		sp := new(model.SignPayInfo)
		if err = rows.Scan(&sp.SignID, &sp.DueDate, &sp.PayValue); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		sps = append(sps, sp)
	}
	err = rows.Err()
	return
}

// McnSignMids .
func (d *Dao) McnSignMids(c context.Context) (msid map[int64]int64, sids []int64, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _selMcnSignMidsSQL); err != nil {
		return
	}
	defer rows.Close()
	msid = make(map[int64]int64)
	for rows.Next() {
		var signID, mcnMid int64
		if err = rows.Scan(&signID, &mcnMid); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		msid[signID] = mcnMid
		sids = append(sids, signID)
	}
	err = rows.Err()
	return
}

// McnUPCount .
func (d *Dao) McnUPCount(c context.Context, signIDs []int64) (mmc map[int64]int64, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selMcnUPCountSQL, xstr.JoinInts(signIDs))); err != nil {
		return
	}
	defer rows.Close()
	mmc = make(map[int64]int64)
	for rows.Next() {
		var signID, count int64
		if err = rows.Scan(&signID, &count); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		mmc[signID] = count
	}
	err = rows.Err()
	return
}

// McnUPMids .
func (d *Dao) McnUPMids(c context.Context, signIDs []int64) (mup map[int64][]int64, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selMcnUPMidsSQL, xstr.JoinInts(signIDs))); err != nil {
		return
	}
	defer rows.Close()
	mup = make(map[int64][]int64)
	for rows.Next() {
		var signID, upMid int64
		if err = rows.Scan(&signID, &upMid); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		mup[signID] = append(mup[signID], upMid)
	}
	err = rows.Err()
	return
}

// CrmUpMidsSum .
func (d *Dao) CrmUpMidsSum(c context.Context, upMids []int64) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selCrmUpMidsSumSQL, xstr.JoinInts(upMids)))
	var countNull sql.NullInt64
	if err = row.Scan(&countNull); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
	}
	count = countNull.Int64
	return
}

// McnSignPayDues .
func (d *Dao) McnSignPayDues(c context.Context) (sps []*model.SignPayInfo, err error) {
	var (
		rows     *xsql.Rows
		now, future   time.Time
		nowDate        = time.Now()
		date     = nowDate.AddDate(0, 0, 7)
		template = date.Format(model.TimeFormatDay)
		nowTemplate    = nowDate.Format(model.TimeFormatDay)
	)
	if now, err = time.ParseInLocation(model.TimeFormatDay, nowTemplate, time.Local); err != nil {
		err = errors.Errorf("time.ParseInLocation(%s) now error(%+v)", nowTemplate, err)
		return
	}
	if future, err = time.ParseInLocation(model.TimeFormatDay, template, time.Local); err != nil {
		err = errors.Errorf("time.ParseInLocation(%s) error(%+v)", template, err)
		return
	}
	if rows, err = d.db.Query(c, _selMcnSignPayDues, future, now); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		sp := new(model.SignPayInfo)
		if err = rows.Scan(&sp.SignPayID, &sp.McnMid, &sp.SignID, &sp.DueDate, &sp.PayValue); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		sps = append(sps, sp)
	}
	err = rows.Err()
	return
}

// McnSignDues .
func (d *Dao) McnSignDues(c context.Context) (mss []*model.MCNSignInfo, err error) {
	var (
		rows           *xsql.Rows
		now, future    time.Time
		nowDate        = time.Now()
		nowTemplate    = nowDate.Format(model.TimeFormatDay)
		futureDate     = nowDate.AddDate(0, 0, 30)
		futureTemplate = futureDate.Format(model.TimeFormatDay)
	)
	if now, err = time.ParseInLocation(model.TimeFormatDay, nowTemplate, time.Local); err != nil {
		err = errors.Errorf("time.ParseInLocation(%s) now error(%+v)", nowTemplate, err)
		return
	}
	if future, err = time.ParseInLocation(model.TimeFormatDay, futureTemplate, time.Local); err != nil {
		err = errors.Errorf("time.ParseInLocation(%s) future error(%+v)", futureTemplate, err)
		return
	}
	if rows, err = d.db.Query(c, _selMcnSignDues, future, now); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		ms := new(model.MCNSignInfo)
		if err = rows.Scan(&ms.SignID, &ms.McnMid, &ms.BeginDate, &ms.EndDate); err != nil {
			if err == xsql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		mss = append(mss, ms)
	}
	err = rows.Err()
	return
}
