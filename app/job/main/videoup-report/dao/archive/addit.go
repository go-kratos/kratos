package archive

import (
	"context"

	"fmt"
	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_additSQL  = "SELECT id,aid,source,redirect_url,mission_id,up_from,order_id,dynamic FROM archive_addit WHERE aid=?"
	_additsSQL = "SELECT id,aid,source,redirect_url,mission_id,up_from,order_id,dynamic FROM archive_addit WHERE aid IN (%s)"
)

// Addit get archive addit.
func (d *Dao) Addit(c context.Context, aid int64) (addit *archive.Addit, err error) {
	row := d.db.QueryRow(c, _additSQL, aid)
	addit = &archive.Addit{}
	if err = row.Scan(&addit.ID, &addit.Aid, &addit.Source, &addit.RedirectURL, &addit.MissionID, &addit.UpFrom, &addit.OrderID, &addit.Dynamic); err != nil {
		if err == sql.ErrNoRows {
			addit = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Addits batch get archive addit.
func (d *Dao) Addits(c context.Context, aids []int64) (addits map[int64]*archive.Addit, err error) {
	addits = make(map[int64]*archive.Addit)
	if len(aids) < 1 {
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_additsSQL, xstr.JoinInts(aids)))
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", fmt.Sprintf(_additsSQL, xstr.JoinInts(aids)), err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		addit := &archive.Addit{}
		if err = rows.Scan(&addit.ID, &addit.Aid, &addit.Source, &addit.RedirectURL, &addit.MissionID, &addit.UpFrom, &addit.OrderID, &addit.Dynamic); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		addits[addit.Aid] = addit
	}
	return
}
