package archive

import (
	"context"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/log"
)

const (
	_tpsSQL = "SELECT id,pid,name,description FROM archive_type"
)

// TypeMapping is second types opposite first types.
func (d *Dao) TypeMapping(c context.Context) (tmap map[int16]*archive.Type, err error) {
	rows, err := d.rddb.Query(c, _tpsSQL)
	if err != nil {
		log.Error("d.tpsStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	tmap = make(map[int16]*archive.Type)
	for rows.Next() {
		t := &archive.Type{}
		if err = rows.Scan(&t.ID, &t.PID, &t.Name, &t.Desc); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		tmap[t.ID] = t
	}
	return
}
