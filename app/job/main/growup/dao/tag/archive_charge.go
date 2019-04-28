package tag

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/tag"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	// select
	_avDailyChargeSQL = "SELECT id,av_id,mid,tag_id,inc_charge,upload_time,is_deleted FROM av_daily_charge_%s WHERE %s LIMIT ?,?"
	_cmDailyChargeSQL = "SELECT id,aid,mid,tag_id,inc_charge,upload_time,is_deleted FROM column_daily_charge WHERE %s LIMIT ?,?"
	_getBGMSQL        = "SELECT id,sid,mid,join_at FROM background_music WHERE id > ? ORDER BY id LIMIT ?"
)

// GetAvDailyCharge get av_charge by tagID and date
func (d *Dao) GetAvDailyCharge(c context.Context, month string, query string, from, limit int) (avs []*model.ArchiveCharge, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_avDailyChargeSQL, month, query), from, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		av := &model.ArchiveCharge{}
		err = rows.Scan(&av.ID, &av.AID, &av.MID, &av.CategoryID, &av.IncCharge, &av.UploadTime, &av.IsDeleted)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		avs = append(avs, av)
	}
	return
}

// GetCmDailyCharge get av_charge by tagID and date
func (d *Dao) GetCmDailyCharge(c context.Context, query string, from, limit int) (cms []*model.ArchiveCharge, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_cmDailyChargeSQL, query), from, limit)
	if err != nil {
		log.Error("d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		cm := &model.ArchiveCharge{}
		var uploadTime int64
		err = rows.Scan(&cm.ID, &cm.AID, &cm.MID, &cm.CategoryID, &cm.IncCharge, &uploadTime, &cm.IsDeleted)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		cm.UploadTime = xtime.Time(uploadTime)
		cms = append(cms, cm)
	}
	return
}

// GetBgm get bgms
func (d *Dao) GetBgm(c context.Context, id int64, limit int64) (bs []*model.ArchiveCharge, last int64, err error) {
	rows, err := d.db.Query(c, _getBGMSQL, id, limit)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		b := &model.ArchiveCharge{}
		err = rows.Scan(&last, &b.AID, &b.MID, &b.UploadTime)
		if err != nil {
			return
		}
		bs = append(bs, b)
	}
	return
}
