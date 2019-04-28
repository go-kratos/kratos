package stat

import (
	"context"
	"fmt"

	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_folderStatSharding int64 = 100

	// stat
	_statSQL = "SELECT play,fav,share from fav_folder_stat_%s WHERE fid=?"

	_upsertPlaySQL  = "INSERT INTO fav_folder_stat_%s (fid,play) VALUES(?,?) ON DUPLICATE KEY UPDATE play=?"
	_upsertFavSQL   = "INSERT INTO fav_folder_stat_%s (fid,fav) VALUES(?,?) ON DUPLICATE KEY UPDATE fav=?"
	_upsertShareSQL = "INSERT INTO fav_folder_stat_%s (fid,share) VALUES(?,?) ON DUPLICATE KEY UPDATE share=?"
)

// UpdateFav updates stat in db.
func (d *Dao) UpdateFav(c context.Context, oid, count int64) (rows int64, err error) {
	fid, table := hit(oid)
	res, err := d.db.Exec(c, fmt.Sprintf(_upsertFavSQL, table), fid, count, count)
	if err != nil {
		log.Error("UpdateFav(%d,%d) error(%+v)", oid, count, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// UpdateShare updates stat in db.
func (d *Dao) UpdateShare(c context.Context, oid, count int64) (rows int64, err error) {
	fid, table := hit(oid)
	res, err := d.db.Exec(c, fmt.Sprintf(_upsertShareSQL, table), fid, count, count)
	if err != nil {
		log.Error("UpdateShare(%d,%d) error(%+v)", oid, count, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// UpdatePlay updates stat in db.
func (d *Dao) UpdatePlay(c context.Context, oid, count int64) (rows int64, err error) {
	fid, table := hit(oid)
	res, err := d.db.Exec(c, fmt.Sprintf(_upsertPlaySQL, table), fid, count, count)
	if err != nil {
		log.Error("UpdatePlay(%d) error(%+v)", oid, err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Stat return stat count from mysql.
func (d *Dao) Stat(c context.Context, oid int64) (f *favmdl.Folder, err error) {
	fid, table := hit(oid)
	f = &favmdl.Folder{}
	row := d.db.QueryRow(c, fmt.Sprintf(_statSQL, table), fid)
	if err = row.Scan(&f.PlayCount, &f.FavedCount, &f.ShareCount); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			f = nil
			return
		}
		log.Error("Stat(%v) error(%+v)", f, err)
	}
	return
}
