package archive

import (
	"context"
	"fmt"

	"go-common/library/database/sql"
	"go-common/library/xstr"
)

const (
	_staffSQL    = "SELECT aid FROM archive_staff WHERE staff_mid = ? AND state = 1"
	_staffsSQL   = "SELECT aid, staff_mid FROM archive_staff WHERE staff_mid IN (%s) AND state = 1"
	_staffAidSQL = "SELECT staff_mid FROM archive_staff WHERE aid = ? AND state = 1"
)

// Staff get upper staff aids by mid.
func (d *Dao) Staff(c context.Context, mid int64) (aids []int64, err error) {
	d.infoProm.Incr("Staff")
	rows, err := d.archiveDB.Query(c, _staffSQL, mid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		aids = append(aids, aid)
	}
	return
}

// Staffs get uppers staff aids by mids.
func (d *Dao) Staffs(c context.Context, mids []int64) (aidm map[int64][]int64, err error) {
	d.infoProm.Incr("Staffs")
	rows, err := d.archiveDB.Query(c, fmt.Sprintf(_staffsSQL, xstr.JoinInts(mids)))
	if err != nil {
		return
	}
	defer rows.Close()
	aidm = make(map[int64][]int64, len(mids))
	for rows.Next() {
		var (
			aid      int64
			staffMid int64
		)
		if err = rows.Scan(&aid, &staffMid); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		aidm[staffMid] = append(aidm[staffMid], aid)
	}
	return
}

// StaffAid get uppers staff mid-list by aid.
func (d *Dao) StaffAid(c context.Context, aid int64) (mids []int64, err error) {
	d.infoProm.Incr("StaffAid")
	rows, err := d.archiveDB.Query(c, _staffAidSQL, aid)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		if err = rows.Scan(&mid); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			return
		}
		mids = append(mids, mid)
	}
	return
}
