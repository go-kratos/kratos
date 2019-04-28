package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/tv/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_selectSQL    = "SELECT page_id,title,index_type,index_tid FROM tv_pages WHERE deleted=0 AND valid=1 ORDER BY rank ASC"
	_findMaxmTime = "SELECT max(mtime) FROM tv_pages WHERE deleted=0 AND valid=1"
)

// Regions .
func (d *Dao) Regions(c context.Context) (res []*model.Region, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _selectSQL); err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Region)
		if err = rows.Scan(&r.PageID, &r.Title, &r.IndexType, &r.IndexTid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// FindLastMtime .
func (d *Dao) FindLastMtime(c context.Context) (res int64, err error) {
	var m time.Time
	if err = d.db.QueryRow(c, _findMaxmTime).Scan(&m); err != nil {
		log.Error("d.db.QueryRow error(%v)", err)
		return
	}
	res = m.Unix()
	return
}
