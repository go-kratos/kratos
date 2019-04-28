package dao

import (
	"context"
	"go-common/app/service/live/gift/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

var _getDiscountPlan = "SELECT id,scene_key,scene_value,platform FROM discount_plan WHERE online_time < ? AND offline_time > ? ORDER BY ctime"

// GetDiscountPlan GetDiscountPlan
func (d *Dao) GetDiscountPlan(ctx context.Context, now time.Time) (plans []*model.DiscountPlan, err error) {
	log.Info("GetDiscountPlan")
	var rows *sql.Rows
	var curTime = now.Format("2006-01-02 15:04:05")
	curTime = "2018-07-20 00:00:00"
	if rows, err = d.db.Query(ctx, _getDiscountPlan, curTime, curTime); err != nil {
		log.Error("query GetDiscountPlan error,err %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		d := &model.DiscountPlan{}
		if err = rows.Scan(&d.Id, &d.SceneKey, &d.SceneValue, &d.Platform); err != nil {
			log.Error("GetDiscountPlan scan error,err %v", err)
			return
		}
		plans = append(plans, d)
	}
	return
}
