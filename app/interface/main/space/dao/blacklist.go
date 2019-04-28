package dao

import (
	"context"

	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_blacklistSQL = `SELECT mid FROM blacklist WHERE status = 0`
)

// Blacklist get blacklist from db.
func (d *Dao) Blacklist(c context.Context) (blacklist []int64, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _blacklistSQL); err != nil {
		log.Error("dao.Modules.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		if err = rows.Scan(&mid); err != nil {
			log.Error("Space dao Modules:row.Scan() error(%v)", err)
			return
		}
		blacklist = append(blacklist, mid)
	}
	if err = rows.Err(); err != nil {
		log.Error("Space dao Modules.Err() error(%v)", err)
	}
	return
}
