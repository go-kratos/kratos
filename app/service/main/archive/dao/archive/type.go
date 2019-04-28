package archive

import (
	"context"

	"go-common/app/service/main/archive/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_tpsSQL = "SELECT id,pid,name FROM archive_type"
)

// Types get type relation.
//func (d *Dao) Types(c context.Context) (nm map[int16]string, rids []int16, sf map[int16]int16, err error) {
func (d *Dao) Types(c context.Context) (types map[int16]*archive.ArcType, err error) {
	d.infoProm.Incr("Types")
	var rows *sql.Rows
	if rows, err = d.tpsStmt.Query(c); err != nil {
		log.Error("tpsStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	types = make(map[int16]*archive.ArcType)
	for rows.Next() {
		var (
			rid, pid int16
			name     string
		)
		if err = rows.Scan(&rid, &pid, &name); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		types[rid] = &archive.ArcType{ID: rid, Pid: pid, Name: name}
	}
	return
}
