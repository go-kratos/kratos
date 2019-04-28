package archive

import (
	"context"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_slFirstPassByAID = "SELECT `id` FROM `archive_first_pass` WHERE `aid`=? LIMIT 1"
)

// GetFirstPassByAID is
func (d *Dao) GetFirstPassByAID(c context.Context, aid int64) (id int64, err error) {
	row := d.db.QueryRow(c, _slFirstPassByAID, aid)
	if err = row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("GetFirstPassByAID error(%v) aid(%d)", err, aid)
		}
	}
	return
}
