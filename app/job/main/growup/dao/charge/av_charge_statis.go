package charge

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/charge"

	"go-common/library/log"
)

const (
	_avChargeStatisSQL = "SELECT id,av_id,mid,total_charge,upload_time FROM av_charge_statis WHERE id > ? ORDER BY id LIMIT ?"

	_inAvChargeStatisSQL = "INSERT INTO av_charge_statis(av_id,mid,tag_id,is_original,total_charge,upload_time) VALUES %s ON DUPLICATE KEY UPDATE total_charge=VALUES(total_charge)"
)

// AvChargeStatis get av_charge_statis
func (d *Dao) AvChargeStatis(c context.Context, id int64, limit int) (data []*model.AvChargeStatis, err error) {
	rows, err := d.db.Query(c, _avChargeStatisSQL, id, limit)
	if err != nil {
		log.Error("d.Query av_charge_statis error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ads := &model.AvChargeStatis{}
		err = rows.Scan(&ads.ID, &ads.AvID, &ads.MID, &ads.TotalCharge, &ads.UploadTime)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		data = append(data, ads)
	}
	return
}

// InsertAvChargeStatisBatch add av charge statis batch
func (d *Dao) InsertAvChargeStatisBatch(c context.Context, vals string) (count int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inAvChargeStatisSQL, vals))
	if err != nil {
		log.Error("InsertAvChargeStatisBatch d.db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
