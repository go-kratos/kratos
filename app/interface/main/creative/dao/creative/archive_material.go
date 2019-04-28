package creative

import (
	"context"

	"database/sql"
	"go-common/library/log"
)

const (
	_getMaterialDataSQL = " SELECT data from archive_material  where aid=? and cid=? and type=? limit 1"
)

// RawBgmData fn
func (d *Dao) RawBgmData(c context.Context, aid, cid, mtype int64) (data *BgmData, err error) {
	data = &BgmData{}
	row := d.creativeDb.QueryRow(c, _getMaterialDataSQL, aid, cid, mtype)
	if err = row.Scan(&data.Data); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}
