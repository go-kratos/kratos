package income

import (
	"context"
	"fmt"

	model "go-common/app/admin/main/growup/model/income"

	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// insert
	_inAvBreachSQL = "INSERT INTO av_breach_record(av_id,mid,cdate,money,ctype,reason,upload_time) VALUES %s"

	// select
	_avBreachByMIDsSQL = "SELECT av_id, mid, cdate, money FROM av_breach_record WHERE mid in (%s) AND ctype in (%s)"
	_breachSQL         = "SELECT av_id,mid,cdate,money,ctype,reason,upload_time FROM av_breach_record WHERE %s"
	_breachCountSQL    = "SELECT count(*) FROM av_breach_record WHERE %s"

	// update
	_upAvBreachPreSQL = "UPDATE av_breach_pre SET state = 2 WHERE aid IN (%s) AND ctype = 0 AND cdate <= '%s'"
)

// BreachCount breach count
func (d *Dao) BreachCount(c context.Context, query string) (total int, err error) {
	err = d.db.QueryRow(c, fmt.Sprintf(_breachCountSQL, query)).Scan(&total)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// ListArchiveBreach list av_breach_record by query
func (d *Dao) ListArchiveBreach(c context.Context, query string) (breachs []*model.AvBreach, err error) {
	breachs = make([]*model.AvBreach, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_breachSQL, query))
	if err != nil {
		log.Error("ListArchiveBreach d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.AvBreach{}
		err = rows.Scan(&b.AvID, &b.MID, &b.CDate, &b.Money, &b.CType, &b.Reason, &b.UploadTime)
		if err != nil {
			log.Error("ListArchiveBreach rows.Scan error(%v)", err)
			return
		}
		breachs = append(breachs, b)
	}
	err = rows.Err()
	return
}

// TxInsertAvBreach insert av_breach_record
func (d *Dao) TxInsertAvBreach(tx *sql.Tx, val string) (rows int64, err error) {
	if val == "" {
		return
	}
	res, err := tx.Exec(fmt.Sprintf(_inAvBreachSQL, val))
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// GetAvBreachByMIDs  get av_breach_record by mids
func (d *Dao) GetAvBreachByMIDs(c context.Context, mids []int64, types []int64) (breachs []*model.AvBreach, err error) {
	if len(mids) == 0 {
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_avBreachByMIDsSQL, xstr.JoinInts(mids), xstr.JoinInts(types)))
	if err != nil {
		log.Error("GetAvBreachByMIDs d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.AvBreach{}
		err = rows.Scan(&b.AvID, &b.MID, &b.CDate, &b.Money)
		if err != nil {
			log.Error("GetAvBreachByMIDs rows scan error(%v)", err)
			return
		}
		breachs = append(breachs, b)
	}

	err = rows.Err()
	return
}

// TxUpdateBreachPre update av_breach_pre state = 2
func (d *Dao) TxUpdateBreachPre(tx *sql.Tx, aids []int64, cdate string) (rows int64, err error) {
	if len(aids) == 0 {
		return
	}
	res, err := tx.Exec(fmt.Sprintf(_upAvBreachPreSQL, xstr.JoinInts(aids), cdate))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
