package archive

import (
	"context"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

//私单业务
const (
	_inPorderSQL  = `INSERT INTO archive_porder (aid,industry_id,brand_id,brand_name,official,show_type,show_front,advertiser,agent,state) VALUES (?,?,?,?,?,?,1,?,?,0) ON DUPLICATE KEY UPDATE industry_id=?,brand_id=?,brand_name=?,official=?,show_type=?,advertiser=?,agent=?`
	_selPorderSQL = `select industry_id,brand_id,brand_name,official,show_type,advertiser,agent,state from archive_porder where aid=?`

	_pconfigSQL = `select id, type, name from porder_config where state = 0 order by rank desc,type asc`
	_parcsSQL   = `select aid,industry_id,brand_id,brand_name,official,show_type,advertiser,agent,state,show_front,ctime,mtime from archive_porder WHERE ctime BETWEEN ? AND ? order by id desc`
)

// TxUpPorder add or update archive_porder
func (d *Dao) TxUpPorder(tx *sql.Tx, aid int64, ap *archive.ArcParam) (rows int64, err error) {
	if ap.Porder.Official == 1 {
		ap.Porder.BrandName = ""
	}
	res, err := tx.Exec(_inPorderSQL, aid, ap.Porder.IndustryID, ap.Porder.BrandID, ap.Porder.BrandName, ap.Porder.Official, ap.Porder.ShowType, ap.Porder.Advertiser, ap.Porder.Agent, ap.Porder.IndustryID, ap.Porder.BrandID, ap.Porder.BrandName, ap.Porder.Official, ap.Porder.ShowType, ap.Porder.Advertiser, ap.Porder.Agent)
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

// PorderCfgList fn
func (d *Dao) PorderCfgList(c context.Context) (pcfgs []*archive.Pconfig, err error) {
	rows, err := d.rddb.Query(c, _pconfigSQL)
	if err != nil {
		log.Error("d.db.Query(%s)|error(%v)", _pconfigSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		cfg := &archive.Pconfig{}
		if err = rows.Scan(&cfg.ID, &cfg.Tp, &cfg.Name); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		pcfgs = append(pcfgs, cfg)
	}
	return
}

// PorderArcList fn
func (d *Dao) PorderArcList(c context.Context, begin, end time.Time) (res []*archive.PorderArc, err error) {
	res = []*archive.PorderArc{}
	rows, err := d.rddb.Query(c, _parcsSQL, begin, end)
	if err != nil {
		log.Error("PorderArcList error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &archive.PorderArc{}
		if err = rows.Scan(&r.AID, &r.IndustryID, &r.BrandID, &r.BrandName, &r.Official, &r.ShowType, &r.Advertiser, &r.Agent, &r.State, &r.ShowFront, &r.Ctime, &r.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}
