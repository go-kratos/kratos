package dao

import (
	"context"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_updateMask = "REPLACE INTO dm_masking(cid,plat,fps,time,list) VALUES (?,?,?,?,?)"
	_selectMask = "SELECT cid,plat,fps,time,list FROM dm_masking WHERE cid=? AND plat=?"
)

// UpdateMask replace dm_masking table for web
func (d *Dao) UpdateMask(c context.Context, cid, maskTime int64, fps int32, plat int8, list string) (err error) {
	if _, err = d.dbDM.Exec(c, _updateMask, cid, plat, fps, maskTime, list); err != nil {
		log.Error("biliDM.Exec(%v, %v %v %v %v %v) error(%v)", _updateMask, cid, plat, fps, maskTime, list, err)
	}
	return
}

// MaskList get mask linfo
func (d *Dao) MaskList(c context.Context, cid int64, plat int8) (m *model.Mask, err error) {
	m = &model.Mask{}
	var tmp string
	row := d.dbDM.QueryRow(c, _selectMask, cid, plat)
	if err = row.Scan(&m.Cid, &m.Plat, &m.FPS, &m.Time, &tmp); err != nil {
		if err == sql.ErrNoRows {
			m = nil
			err = nil
		} else {
			log.Error("MaskList.rows.Scan(cid:%d plat:%d) error(%v)", cid, plat, err)
		}
		return
	}
	if tmp == "" {
		m = nil
		return
	}
	m.MaskURL = d.conf.Host.MaskCloud + tmp
	return
}
