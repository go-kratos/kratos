package dao

import (
	"context"

	"go-common/app/admin/main/dm/model"
	"go-common/library/log"
)

const (
	_countMaskUp = "SELECT count(*) FROM dm_mask_up WHERE state=1"
	_getMaskUp   = "SELECT id,mid,state,comment,ctime,mtime from dm_mask_up where state=1 ORDER BY mtime DESC limit ?,? "
	_maskUpOpen  = "INSERT INTO dm_mask_up(mid,state,comment) VALUES (?,?,?) ON DUPLICATE KEY UPDATE state=?,comment=?"
)

// MaskUps get mask up info from db.
func (d *Dao) MaskUps(c context.Context, pn, ps int64) (maskUps []*model.MaskUp, total int64, err error) {
	maskUps = make([]*model.MaskUp, 0)
	countRow := d.biliDM.QueryRow(c, _countMaskUp)
	if err = countRow.Scan(&total); err != nil {
		log.Error("row.ScanCount error(%v)", err)
		return
	}
	rows, err := d.biliDM.Query(c, _getMaskUp, (pn-1)*ps, ps)
	if err != nil {
		log.Error("biliDM.Query(%s) error(%v)", _getMaskUp, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		maskUp := &model.MaskUp{}
		if err = rows.Scan(&maskUp.ID, &maskUp.Mid, &maskUp.State, &maskUp.Comment, &maskUp.CTime, &maskUp.MTime); err != nil {
			log.Error("biliDM.Scan(%s) error(%v)", _getMaskUp, err)
			return
		}
		maskUps = append(maskUps, maskUp)
	}
	if err = rows.Err(); err != nil {
		log.Error("biliDM.rows.Err() error(%v)", err)
	}
	return
}

// MaskUpOpen mask up open.
func (d *Dao) MaskUpOpen(c context.Context, mid int64, state int32, comment string) (affect int64, err error) {
	res, err := d.biliDM.Exec(c, _maskUpOpen, mid, state, comment, state, comment)
	if err != nil {
		log.Error("d.biliDM.Exec(%s,%d,%v,%v) error(%v)", _maskUpOpen, mid, state, comment, err)
		return
	}
	return res.RowsAffected()
}
