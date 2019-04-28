package archive

import (
	"context"
	bsql "database/sql"

	"fmt"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_inApplySQL               = `INSERT INTO archive_staff_apply (type,as_id,apply_aid,apply_up_mid,apply_staff_mid,apply_title,apply_title_id,state,deal_state) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE type=?,apply_title=?,apply_title_id=?,state=?,deal_state=?,as_id=?`
	_selApplySQL              = `SELECT  sa.id,sa.type,sa.as_id,sa.apply_aid,sa.apply_up_mid,sa.apply_staff_mid,sa.apply_title,sa.apply_title_id,sa.state,sa.deal_state,s.state as staff_state,s.staff_title  FROM  archive_staff_apply  sa LEFT JOIN  archive_staff s  on s.id=sa.as_id where sa.id=?`
	_selApplysSQL             = `SELECT  sa.id,sa.type,sa.as_id,sa.apply_aid,sa.apply_up_mid,sa.apply_staff_mid,sa.apply_title,sa.apply_title_id,sa.state,sa.deal_state,s.state as staff_state,s.staff_title  FROM  archive_staff_apply  sa LEFT JOIN  archive_staff s  on s.id=sa.as_id where sa.id IN(%s)`
	_selApplysByAIDSQL        = `SELECT  sa.id,sa.type,sa.as_id,sa.apply_aid,sa.apply_up_mid,sa.apply_staff_mid,sa.apply_title,sa.apply_title_id,sa.state,sa.deal_state,s.state as staff_state,s.staff_title  FROM  archive_staff_apply  sa LEFT JOIN  archive_staff s  on s.id=sa.as_id where sa.apply_aid =?`
	_selApplysByAIDSANDMIDSQL = `SELECT  sa.id,sa.type,sa.as_id,sa.apply_aid,sa.apply_up_mid,sa.apply_staff_mid,sa.apply_title,sa.apply_title_id,sa.state,sa.deal_state,s.state as staff_state,s.staff_title  FROM  archive_staff_apply  sa LEFT JOIN  archive_staff s  on s.id=sa.as_id where sa.apply_aid IN(%s) AND sa.apply_staff_mid=%d`
	_selApplysByMIDSTAFFSQL   = `SELECT  sa.id,sa.type,sa.as_id,sa.apply_aid,sa.apply_up_mid,sa.apply_staff_mid,sa.apply_title,sa.apply_title_id,sa.state,sa.deal_state,s.state as staff_state,s.staff_title  FROM  archive_staff_apply  sa LEFT JOIN  archive_staff s  on s.id=sa.as_id where sa.apply_up_mid=? AND  sa.apply_staff_mid =?`
	_midCountSQL              = `select count(*) as count from archive_staff_apply where apply_staff_mid=?`
)

// Apply get archive Apply
func (d *Dao) Apply(c context.Context, ID int64) (p *archive.StaffApply, err error) {
	row := d.rddb.QueryRow(c, _selApplySQL, ID)
	p = &archive.StaffApply{}
	var title bsql.NullString
	var state bsql.NullInt64
	if err = row.Scan(&p.ID, &p.Type, &p.ASID, &p.ApplyAID, &p.ApplyUpMID, &p.ApplyStaffMID, &p.ApplyTitle, &p.ApplyTitleID, &p.State, &p.DealState, &state, &title); err != nil {
		if err == sql.ErrNoRows {
			p = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	p.StaffTitle = title.String
	p.StaffState = int8(state.Int64)
	return
}

// MidCount get
func (d *Dao) MidCount(c context.Context, ID int64) (count int64, err error) {
	row := d.rddb.QueryRow(c, _midCountSQL, ID)
	if err = row.Scan(&count); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan error(%v)", err)
			return
		}
		err = nil
	}
	return
}

// Applys get .
func (d *Dao) Applys(c context.Context, ids []int64) (as []*archive.StaffApply, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_selApplysSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.db.Applys ids(%+v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var title bsql.NullString
		var state bsql.NullInt64
		s := &archive.StaffApply{}
		if err = rows.Scan(&s.ID, &s.Type, &s.ASID, &s.ApplyAID, &s.ApplyUpMID, &s.ApplyStaffMID, &s.ApplyTitle, &s.ApplyTitleID, &s.State, &s.DealState, &state, &title); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		s.StaffTitle = title.String
		s.StaffState = int8(state.Int64)
		as = append(as, s)
	}
	return
}

// FilterApplys get .
func (d *Dao) FilterApplys(c context.Context, aids []int64, mid int64) (as []*archive.StaffApply, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_selApplysByAIDSANDMIDSQL, xstr.JoinInts(aids), mid))
	if err != nil {
		log.Error("d.db.FilterApplys(%v,%d) error(%v)", aids, mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var title bsql.NullString
		var state bsql.NullInt64
		s := &archive.StaffApply{}

		if err = rows.Scan(&s.ID, &s.Type, &s.ASID, &s.ApplyAID, &s.ApplyUpMID, &s.ApplyStaffMID, &s.ApplyTitle, &s.ApplyTitleID, &s.State, &s.DealState, &state, &title); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		s.StaffTitle = title.String
		s.StaffState = int8(state.Int64)
		as = append(as, s)
	}
	return
}

// ApplysByAID get .
func (d *Dao) ApplysByAID(c context.Context, aid int64) (as []*archive.StaffApply, err error) {
	rows, err := d.db.Query(c, _selApplysByAIDSQL, aid)
	if err != nil {
		log.Error("d.db.ApplysByAID aid(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var title bsql.NullString
		var state bsql.NullInt64
		s := &archive.StaffApply{}
		if err = rows.Scan(&s.ID, &s.Type, &s.ASID, &s.ApplyAID, &s.ApplyUpMID, &s.ApplyStaffMID, &s.ApplyTitle, &s.ApplyTitleID, &s.State, &s.DealState, &state, &title); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		s.StaffTitle = title.String
		s.StaffState = int8(state.Int64)
		as = append(as, s)
	}
	return
}

// ApplysByMIDAndStaff get .
func (d *Dao) ApplysByMIDAndStaff(c context.Context, upMID, staffMID int64) (as []*archive.StaffApply, err error) {
	rows, err := d.db.Query(c, _selApplysByMIDSTAFFSQL, upMID, staffMID)
	if err != nil {
		log.Error("d.db.ApplysByAID  upMID(%d) staffMID(%d) error(%v)", upMID, staffMID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var title bsql.NullString
		var state bsql.NullInt64
		s := &archive.StaffApply{}
		if err = rows.Scan(&s.ID, &s.Type, &s.ASID, &s.ApplyAID, &s.ApplyUpMID, &s.ApplyStaffMID, &s.ApplyTitle, &s.ApplyTitleID, &s.State, &s.DealState, &state, &title); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		s.StaffTitle = title.String
		s.StaffState = int8(state.Int64)
		as = append(as, s)
	}
	return
}

// TxAddApply tx.
func (d *Dao) TxAddApply(tx *sql.Tx, param *archive.ApplyParam) (id int64, err error) {
	res, err := tx.Exec(_inApplySQL, param.Type, param.ASID, param.ApplyAID, param.ApplyUpMID, param.ApplyStaffMID, param.ApplyTitle, param.ApplyTitleID, param.State, param.DealState, param.Type, param.ApplyTitle, param.ApplyTitleID, param.State, param.DealState, param.ASID)
	if err != nil {
		log.Error("d.TxAddApply.Exec() error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}
