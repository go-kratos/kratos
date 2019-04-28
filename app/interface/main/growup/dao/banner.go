package dao

import (
	"context"
	"database/sql"

	"go-common/library/log"
	"go-common/library/time"

	"go-common/app/interface/main/growup/model"
)

const (
	_bannerSQL = "SELECT id,image,link,start_at,end_at FROM banner WHERE start_at <= ? AND end_at >= ? LIMIT 1"
)

// Banner get banner
func (d *Dao) Banner(c context.Context, t int64) (b *model.Banner, err error) {
	b = &model.Banner{}
	row := d.db.QueryRow(c, _bannerSQL, time.Time(t), time.Time(t))
	if err = row.Scan(&b.ID, &b.Image, &b.Link, &b.StartAt, &b.EndAt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}
