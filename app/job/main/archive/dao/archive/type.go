package archive

import (
	"context"

	"go-common/library/log"
)

const (
	_tpsSQL = "SELECT id,pid FROM archive_type WHERE pid !=0"
)

// TypeMapping is second types opposite first types.
func (d *Dao) TypeMapping(c context.Context) (rmap map[int16]int16, err error) {
	rows, err := d.db.Query(c, _tpsSQL)
	if err != nil {
		log.Error("d.tpsStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	rmap = map[int16]int16{}
	for rows.Next() {
		var id, pid int16
		if err = rows.Scan(&id, &pid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		rmap[id] = pid
	}
	return
}
