package archive

import (
	"context"

	"go-common/app/job/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_arcPOISQL  = "SELECT data from archive_biz   WHERE aid=? AND type= ?"
	_arcVoteSQL = "SELECT data from archive_biz   WHERE aid=? AND type=2"
)

// POI get a archive POI by avid.
func (d *Dao) POI(c context.Context, aid int64) (data []byte, err error) {
	var (
		row = d.db.QueryRow(c, _arcPOISQL, aid, archive.BIZPOI)
	)
	if err = row.Scan(&data); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}

// Vote get a archive Vote by avid.
func (d *Dao) Vote(c context.Context, aid int64) (data []byte, err error) {
	var (
		row = d.db.QueryRow(c, _arcVoteSQL, aid)
	)
	if err = row.Scan(&data); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}
