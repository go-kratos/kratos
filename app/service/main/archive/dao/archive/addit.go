package archive

import (
	"context"

	"go-common/app/service/main/archive/model/archive"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_additSQL = "SELECT id,aid,source,redirect_url,mission_id,up_from,order_id,description FROM archive_addit WHERE aid=?"
)

// Addit get archive addit.
func (d *Dao) Addit(c context.Context, aid int64) (addit *archive.Addit, err error) {
	d.infoProm.Incr("Addit")
	row := d.additStmt.QueryRow(c, aid)
	addit = &archive.Addit{}
	if err = row.Scan(&addit.ID, &addit.Aid, &addit.Source, &addit.RedirectURL, &addit.MissionID, &addit.UpFrom, &addit.OrderID, &addit.Description); err != nil {
		if err == sql.ErrNoRows {
			addit = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
