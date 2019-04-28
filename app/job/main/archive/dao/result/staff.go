package result

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/job/main/archive/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addStaffSQL = "INSERT INTO archive_staff (aid,mid,title,ctime,mtime) VALUES "
	_delStaffSQL = "DELETE FROM archive_staff WHERE aid=?"
)

// TxDelStaff del archive staff
func (d *Dao) TxDelStaff(c context.Context, tx *sql.Tx, aid int64) (err error) {
	_, err = tx.Exec(_delStaffSQL, aid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return
}

// TxAddStaff add archive staff
func (d *Dao) TxAddStaff(c context.Context, tx *sql.Tx, aid int64, staff []*archive.Staff) (err error) {
	var valSQL []string
	for _, s := range staff {
		valSQL = append(valSQL, fmt.Sprintf("(%d,%d,'%s','%s','%s')", s.Aid, s.Mid, s.Title, s.Ctime, s.Mtime))
	}
	valSQLStr := strings.Join(valSQL, ",")
	_, err = tx.Exec(_addStaffSQL + valSQLStr)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return
}
