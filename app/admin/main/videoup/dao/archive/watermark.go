package archive

import (
	"context"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const _watermark = "SELECT id, info, md5, mid, position, type, uname, url, state, mtime  FROM watermark WHERE mid=? AND state != 0"

//Watermark get watermark
func (d *Dao) Watermark(c context.Context, mid int64) (m []*archive.Watermark, err error) {
	var rows *sql.Rows
	m = []*archive.Watermark{}
	if rows, err = d.creativeDB.Query(c, _watermark, mid); err != nil {
		log.Error("Watermark d.rddb.Query error(%v) mid(%d)", err, mid)
		return
	}
	defer rows.Close()

	for rows.Next() {
		wm := new(archive.Watermark)
		if err = rows.Scan(&wm.ID, &wm.Info, &wm.MD5, &wm.MID, &wm.Position, &wm.Type, &wm.Uname, &wm.URL, &wm.State, &wm.MTime); err != nil {
			log.Error("Watermark rows.Scan error(%v) mid(%d)", err, mid)
			return
		}
		if wm.State == "0" {
			continue
		}

		m = append(m, wm)
	}
	return
}
