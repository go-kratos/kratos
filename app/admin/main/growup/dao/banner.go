package dao

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/growup/model"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	_insertBannerSQL     = "INSERT INTO banner(image,link,start_at,end_at) VALUES(?,?,?,?)"
	_bannersSQL          = "SELECT id,image,link,start_at,end_at FROM banner WHERE id > ? ORDER BY id LIMIT ?"
	_totalBannerCountSQL = "SELECT count(*) FROM banner"

	// start_at > end_at(insert) or end_at < start_at(insert)
	_bannerSQL     = "SELECT id FROM banner WHERE end_at > ? AND NOT ((start_at > ?) OR (end_at < ?))"
	_editBannerSQL = "SELECT id FROM banner WHERE id != ? AND end_at > ? AND NOT ((start_at > ?) OR (end_at < ?))"

	_updateBannerSQL = "UPDATE banner SET image=?,link=?,start_at=?,end_at=? WHERE id=?"
	_updateEndAtSQL  = "UPDATE banner SET end_at=? WHERE id=?"
)

// TotalBannerCount get total banner count
func (d *Dao) TotalBannerCount(c context.Context) (count int64, err error) {
	row := d.rddb.QueryRow(c, _totalBannerCountSQL)
	if err = row.Scan(&count); err != nil {
		log.Error("d.rddb.TotalBannerCount error(%v)", err)
	}
	return
}

// DupEditBanner find duplicate banner id where end_at >= start_at(edit)
func (d *Dao) DupEditBanner(c context.Context, startAt, endAt, now, id int64) (dup int64, err error) {
	row := d.rddb.QueryRow(c, _editBannerSQL, id, time.Time(now), time.Time(endAt), time.Time(startAt))
	if err = row.Scan(&dup); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("d.rddb.DupEditBanner error(%v)", err)
		}
	}
	return
}

// DupBanner find duplicate banner id that end_at > start_at(insert)
func (d *Dao) DupBanner(c context.Context, startAt, endAt, now int64) (dup int64, err error) {
	row := d.rddb.QueryRow(c, _bannerSQL, time.Time(now), time.Time(endAt), time.Time(startAt))
	if err = row.Scan(&dup); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("d.rddb.DupBanner error(%v)", err)
		}
	}
	return
}

// InsertBanner insert banner
func (d *Dao) InsertBanner(c context.Context, image, link string, startAt, endAt int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _insertBannerSQL, image, link, time.Time(startAt), time.Time(endAt))
	if err != nil {
		log.Error("d.db.Exec insert banner error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Banners get banners
func (d *Dao) Banners(c context.Context, offset, limit int64) (bs []*model.Banner, err error) {
	rows, err := d.rddb.Query(c, _bannersSQL, offset, limit)
	if err != nil {
		log.Error("d.db.query banners error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.Banner{}
		err = rows.Scan(&b.ID, &b.Image, &b.Link, &b.StartAt, &b.EndAt)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		bs = append(bs, b)
	}
	return
}

// UpdateBanner update banner
func (d *Dao) UpdateBanner(c context.Context, image, link string, startAt, endAt, id int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _updateBannerSQL, image, link, time.Time(startAt), time.Time(endAt), id)
	if err != nil {
		log.Error("d.db.Exec update banner error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpdateBannerEndAt update banner end at
func (d *Dao) UpdateBannerEndAt(c context.Context, endAt, id int64) (rows int64, err error) {
	res, err := d.rddb.Exec(c, _updateEndAtSQL, time.Time(endAt), id)
	if err != nil {
		log.Error("d.db.Exec update banner end_at error(%v)", err)
		return
	}
	return res.RowsAffected()
}
