package dao

import (
	"context"
	"strings"
	"time"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addDmSpecialLocationSQL = "INSERT INTO dm_special_content_location (type,oid,locations,ctime) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE locations=?"
	_getDmSpecialLocationSQL = "SELECT locations FROM dm_special_content_location WHERE oid=? AND type=?"
)

// UpsertDmSpecialLocation .
func (d *Dao) UpsertDmSpecialLocation(c context.Context, tp int32, oid int64, locations string) (err error) {
	if _, err = d.dmWriter.Exec(c, _addDmSpecialLocationSQL, tp, oid, locations, time.Now(), locations); err != nil {
		log.Error("AddDmSpecialLocation.Exec error(%v)", err)
	}
	return
}

// DMSpecialLocations .
func (d *Dao) DMSpecialLocations(c context.Context, tp int32, oid int64) (locations []string, err error) {
	row := d.dmReader.QueryRow(c, _getDmSpecialLocationSQL, oid, tp)
	var s string
	if err = row.Scan(&s); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("DMSpecialLocations.Query error(%v)", err)
		}
		return
	}
	locations = strings.Split(s, ",")
	return
}
