package dao

import (
	"context"
	"go-common/app/service/live/gift/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

var _getOnlinePlan = "SELECT id,list,silver_list,scene_key,scene_value,mtime,platform FROM gift_plan WHERE online_time <= ? AND offline_time >= ? ORDER BY scene_key DESC,mtime DESC"

// GetOnlinePlan GetOnlinePlan
func (d *Dao) GetOnlinePlan(ctx context.Context) (plans []*model.GiftPlan, err error) {
	log.Info("GetOnlinePlan")
	var rows *sql.Rows
	var curTime = time.Now().Format("2006-01-02 15:04:05")
	if rows, err = d.db.Query(ctx, _getOnlinePlan, curTime, curTime); err != nil {
		log.Error("query getOnlinePlan error,err %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := &model.GiftPlan{}
		if err = rows.Scan(&p.Id, &p.List, &p.SilverList, &p.SceneKey, &p.SceneValue, &p.Mtime, &p.Platform); err != nil {
			log.Error("getOnlinePlan scan error,err %v", err)
			return
		}
		plans = append(plans, p)
	}
	return
}
