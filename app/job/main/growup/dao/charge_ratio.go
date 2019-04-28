package dao

import (
	"context"

	"go-common/app/job/main/growup/model"

	"go-common/library/log"
)

const (
	_avDailyIncChargeSQL = "SELECT inc_charge, tag_id FROM av_daily_charge_06 WHERE av_id = ? AND date = '2018-06-24' AND upload_time >= '2018-06-24'"
	_avChargeRatioSQL    = "SELECT id,av_id,ratio,adjust_type FROM av_charge_ratio WHERE id > ? ORDER BY id LIMIT ?"
)

// AvDailyIncCharge get av_daily_charge inc_charge
func (d *Dao) AvDailyIncCharge(c context.Context, avID int64) (incCharge, tagID int64, err error) {
	err = d.db.QueryRow(c, _avDailyIncChargeSQL, avID).Scan(&incCharge, &tagID)
	return
}

// AvChargeRatio get av_charge_ratio
func (d *Dao) AvChargeRatio(c context.Context, id int64, limit int64) (m map[int64]*model.AvChargeRatio, last int64, err error) {
	rows, err := d.db.Query(c, _avChargeRatioSQL, id, limit)
	if err != nil {
		log.Error("d.db.Query AvChargeRatio error(%v)", err)
		return
	}
	m = make(map[int64]*model.AvChargeRatio)
	defer rows.Close()
	for rows.Next() {
		ratio := &model.AvChargeRatio{}
		err = rows.Scan(&last, &ratio.AvID, &ratio.Ratio, &ratio.AdjustType)
		if err != nil {
			log.Error("AvChargeRatio scan error(%v)", err)
			return
		}
		m[ratio.AvID] = ratio
	}
	return
}
