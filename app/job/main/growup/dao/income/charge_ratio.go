package income

import (
	"context"

	model "go-common/app/job/main/growup/model/income"

	"go-common/library/log"
)

const (
	_avChargeRatioSQL = "SELECT id,av_id,ratio,adjust_type,ctype FROM av_charge_ratio WHERE id > ? ORDER BY id LIMIT ?"
	_upChargeRatioSQL = "SELECT id,mid,ratio,adjust_type,ctype FROM up_charge_ratio WHERE id > ? ORDER BY id LIMIT ?"
)

// ArchiveChargeRatio map[ctype]map[archive_id]*archiveChargeRatio
func (d *Dao) ArchiveChargeRatio(c context.Context, id int64, limit int64) (m map[int]map[int64]*model.ArchiveChargeRatio, last int64, err error) {
	rows, err := d.db.Query(c, _avChargeRatioSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query AvChargeRatio error(%v)", err)
		return
	}
	m = make(map[int]map[int64]*model.ArchiveChargeRatio)
	defer rows.Close()
	for rows.Next() {
		ratio := &model.ArchiveChargeRatio{}
		err = rows.Scan(&last, &ratio.ArchiveID, &ratio.Ratio, &ratio.AdjustType, &ratio.CType)
		if err != nil {
			log.Error("AvChargeRatio scan error(%v)", err)
			return
		}
		if ac, ok := m[ratio.CType]; ok {
			ac[ratio.ArchiveID] = ratio
		} else {
			m[ratio.CType] = map[int64]*model.ArchiveChargeRatio{
				ratio.ArchiveID: ratio,
			}
		}
	}
	return
}

// UpChargeRatio get every day up charge ratio
func (d *Dao) UpChargeRatio(c context.Context, id int64, limit int64) (m map[int]map[int64]*model.UpChargeRatio, last int64, err error) {
	rows, err := d.db.Query(c, _upChargeRatioSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query UpChargeRatio error(%v)", err)
		return
	}
	m = make(map[int]map[int64]*model.UpChargeRatio)
	defer rows.Close()
	for rows.Next() {
		ratio := &model.UpChargeRatio{}
		err = rows.Scan(&last, &ratio.MID, &ratio.Ratio, &ratio.AdjustType, &ratio.CType)
		if err != nil {
			log.Error("UpChargeRatio scan error(%v)", err)
			return
		}
		if ur, ok := m[ratio.CType]; ok {
			ur[ratio.MID] = ratio
		} else {
			m[ratio.CType] = map[int64]*model.UpChargeRatio{
				ratio.MID: ratio,
			}
		}
	}
	return
}
