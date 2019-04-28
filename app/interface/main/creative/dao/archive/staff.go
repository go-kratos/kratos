package archive

import (
	"context"

	"go-common/app/interface/main/creative/model/archive"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_staffCountSQL = "SELECT count(*) as count  FROM archive_staff_apply WHERE apply_staff_mid=? AND deal_state = 1"
	_staffSQL      = "SELECT id,aid,mid,staff_mid,staff_title FROM archive_staff WHERE aid=? AND state = 1"
)

// RawStaffData get staff data from db
func (d *Dao) RawStaffData(c context.Context, aid int64) (res []*archive.Staff, err error) {
	rows, err := d.db.Query(c, _staffSQL, aid)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*archive.Staff, 0)
	res = append(res, &archive.Staff{AID: aid})
	var mid int64
	for rows.Next() {
		v := &archive.Staff{}
		if err = rows.Scan(&v.ID, &v.AID, &v.MID, &v.StaffMID, &v.StaffTitle); err != nil {
			log.Error("row.Scan error(%v)", err)
			res = res[0:0]
			return
		}
		res = append(res, v)
		mid = v.MID
	}
	if len(res) == 1 {
		return res[0:0], nil
	}
	for _, v := range res {
		if v.MID == 0 {
			v.MID = mid
			v.StaffMID = mid
			v.StaffTitle = "UP"
			break
		}
	}
	return
}

// CountByMID .
func (d *Dao) CountByMID(c context.Context, mid int64) (count int, err error) {
	row := d.db.QueryRow(c, _staffCountSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
