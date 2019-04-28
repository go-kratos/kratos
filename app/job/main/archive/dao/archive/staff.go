package archive

import (
	"context"

	"go-common/app/job/main/archive/model/archive"
	"go-common/library/log"
)

const (
	_staffSQL = "SELECT aid,staff_mid,staff_title,ctime,mtime FROM archive_staff WHERE aid=? AND state=1"
)

// Staff get Staff by aid.
func (d *Dao) Staff(c context.Context, aid int64) (res []*archive.Staff, err error) {
	rows, err := d.db.Query(c, _staffSQL, aid)
	if err != nil {
		log.Error("d.db.Query(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		s := &archive.Staff{}
		if err = rows.Scan(&s.Aid, &s.Mid, &s.Title, &s.Ctime, &s.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, s)
	}
	err = rows.Err()
	return
}
