package archive

import (
	"context"
	"go-common/app/job/main/videoup/model/archive"
	"go-common/library/log"
)

const (
	_staffsSQL = "SELECT id,aid,mid,staff_mid,staff_title,staff_title_id,state FROM archive_staff WHERE  aid=? AND state=?"
)

// Staffs get .
func (d *Dao) Staffs(c context.Context, AID int64) (fs []*archive.Staff, err error) {
	rows, err := d.db.Query(c, _staffsSQL, AID, archive.STATESTAFFON)
	if err != nil {
		log.Error("d.db.Staffs(%d) error(%v)", AID, err)
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
