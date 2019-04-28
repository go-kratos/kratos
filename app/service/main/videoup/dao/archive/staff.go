package archive

import (
	"context"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_upStateStaffSQL = "UPDATE archive_staff SET state =? where id=?"
	_inStaffSQL      = "INSERT into archive_staff(aid,mid,staff_mid,staff_title,staff_title_id,state)  VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE staff_title=?,staff_title_id=?,state=?"
	_staffsSQL       = "SELECT id,aid,mid,staff_mid,staff_title,staff_title_id,state FROM archive_staff WHERE  aid=? AND state=?"

	_staffByIDSQL     = "SELECT id,aid,mid,staff_mid,staff_title,staff_title_id,state FROM archive_staff WHERE id=?"
	_staffByAIdMIDSQL = "SELECT id,aid,mid,staff_mid,staff_title,staff_title_id,state FROM archive_staff WHERE aid=? AND staff_mid=? limit 1"
)

// TxAddStaff tx.
func (d *Dao) TxAddStaff(tx *sql.Tx, param *archive.Staff) (id int64, err error) {
	res, err := tx.Exec(_inStaffSQL, param.AID, param.MID, param.StaffMID, param.StaffTitle, param.StaffTitleID, param.State, param.StaffTitle, param.StaffTitleID, param.State)
	if err != nil {
		log.Error("d.TxAddStaff.Exec() error(%v)", err)
		return
	}
	id, err = res.LastInsertId()
	return
}

// TxUpStaffState tx .
func (d *Dao) TxUpStaffState(tx *sql.Tx, state int8, id int64) (rows int64, err error) {
	res, err := tx.Exec(_upStateStaffSQL, state, id)
	if err != nil {
		log.Error("d.TxUpStaffState.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Staffs get .
func (d *Dao) Staffs(c context.Context, AID int64) (fs []*archive.Staff, err error) {
	rows, err := d.db.Query(c, _staffsSQL, AID, archive.STATEON)
	if err != nil {
		log.Error("d.db.Staffs aid(%d) error(%v)", AID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		f := &archive.Staff{}
		if err = rows.Scan(&f.ID, &f.AID, &f.MID, &f.StaffMID, &f.StaffTitle, &f.StaffTitleID, &f.State); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		fs = append(fs, f)
	}
	return
}

// Staff get .
func (d *Dao) Staff(c context.Context, ID int64) (s *archive.Staff, err error) {
	row := d.db.QueryRow(c, _staffByIDSQL, ID)
	s = &archive.Staff{}
	if err = row.Scan(&s.ID, &s.AID, &s.MID, &s.StaffMID, &s.StaffTitle, &s.StaffTitleID, &s.State); err != nil {
		if err == sql.ErrNoRows {
			s = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}

// StaffByAidAndMid get .
func (d *Dao) StaffByAidAndMid(c context.Context, AID, StaffMid int64) (s *archive.Staff, err error) {
	row := d.db.QueryRow(c, _staffByAIdMIDSQL, AID, StaffMid)
	s = &archive.Staff{}
	if err = row.Scan(&s.ID, &s.AID, &s.MID, &s.StaffMID, &s.StaffTitle, &s.StaffTitleID, &s.State); err != nil {
		if err == sql.ErrNoRows {
			s = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}
