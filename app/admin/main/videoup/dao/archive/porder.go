package archive

import (
	"context"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inPorderSQL = `INSERT INTO archive_porder (aid,industry_id,brand_id,brand_name,official,show_type,show_front,advertiser,agent,state) VALUES (?,?,?,?,?,?,?,?,?,1) ON DUPLICATE KEY UPDATE
					industry_id=?,brand_id=?,brand_name=?,official=?,show_type=?,show_front=?,advertiser=?,agent=?`
	_selPorderSQL       = `select industry_id,brand_id,brand_name,official,show_type,advertiser,agent,state from archive_porder where aid=?`
	_selPorderConfigSQL = `select id,name,rank,type,state from porder_config`
)

// TxUpPorder  archive_porder
func (d *Dao) TxUpPorder(tx *sql.Tx, aid int64, ap *archive.ArcParam) (rows int64, err error) {
	res, err := tx.Exec(_inPorderSQL, aid, ap.IndustryID, ap.BrandID, ap.BrandName, ap.Official, ap.ShowType, ap.ShowFront, ap.Advertiser, ap.Agent, ap.IndustryID, ap.BrandID, ap.BrandName, ap.Official, ap.ShowType, ap.ShowFront, ap.Advertiser, ap.Agent)
	if err != nil {
		log.Error("d.TxUpPorder.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Porder get archive Proder
func (d *Dao) Porder(c context.Context, aid int64) (p *archive.Porder, err error) {
	row := d.rddb.QueryRow(c, _selPorderSQL, aid)
	p = &archive.Porder{}
	if err = row.Scan(&p.IndustryID, &p.BrandID, &p.BrandName, &p.Official, &p.ShowType, &p.Advertiser, &p.Agent, &p.State); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan error(%v)", err)
			return
		}
		err = nil
	}
	return
}

// PorderConfig get archive ProderConfigs
func (d *Dao) PorderConfig(c context.Context) (pc map[int64]*archive.PorderConfig, err error) {
	rows, err := d.rddb.Query(c, _selPorderConfigSQL)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", err)
		return
	}
	defer rows.Close()
	pc = make(map[int64]*archive.PorderConfig)
	for rows.Next() {
		ap := &archive.PorderConfig{}
		if err = rows.Scan(&ap.ID, &ap.Name, &ap.Rank, &ap.Type, &ap.State); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		pc[ap.ID] = ap
	}
	return
}
