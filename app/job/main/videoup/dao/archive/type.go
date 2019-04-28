package archive

import (
	"context"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/library/log"
)

const (
	_tpsSQL    = "SELECT id,pid,name FROM archive_type WHERE pid !=0"
	_alltpsSQL = "SELECT id,pid,name FROM archive_type"
)

// TypeMapping is second types opposite first types.
func (d *Dao) TypeMapping(c context.Context) (rmap map[int16]int16, err error) {
	rows, err := d.db.Query(c, _tpsSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	rmap = map[int16]int16{}
	for rows.Next() {
		var t = &archive.Type{}
		if err = rows.Scan(&t.ID, &t.PID, &t.Name); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		rmap[t.ID] = t.PID
	}
	return
}

// TypeNaming is all type name
func (d *Dao) TypeNaming(c context.Context) (nmap map[int16]string, err error) {
	rows, err := d.db.Query(c, _alltpsSQL)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	nmap = map[int16]string{}
	for rows.Next() {
		var t = &archive.Type{}
		if err = rows.Scan(&t.ID, &t.PID, &t.Name); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		nmap[t.ID] = t.Name
	}
	return
}
