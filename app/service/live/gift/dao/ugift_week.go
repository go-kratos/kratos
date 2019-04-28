package dao

import (
	"context"
	"encoding/json"
	"go-common/app/service/live/gift/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

var (
	_getWeekBagStatus = "SELECT id,uid,week,level,week_info FROM ugift_week_status WHERE uid = ? AND week = ? AND level =? ORDER BY ctime DESC LIMIT 1"
	_addWeekBag       = "INSERT INTO ugift_week_status (uid,week,level,week_info) VALUES (?,?,?,?)"
)

// GetWeekBagStatus GetWeekBagStatus
func (d *Dao) GetWeekBagStatus(ctx context.Context, uid int64, week int, level int64) (res *model.WeekGiftInfo, err error) {
	log.Info("GetWeekBagStatus,uid:%d,week:%d,level:%d", uid, week, level)
	row := d.db.QueryRow(ctx, _getWeekBagStatus, uid, week, level)
	res = &model.WeekGiftInfo{}
	if err = row.Scan(&res.ID, &res.UID, &res.Week, &res.Level, &res.WeekInfo); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("GetWeekBagStatus row.Scan error(%v)", err)
	}

	return
}

// AddWeekBag AddWeekBag
func (d *Dao) AddWeekBag(ctx context.Context, uid int64, week int, level int64, weekInfo *model.BagGiftStatus) (affected int64, err error) {
	log.Info("AddWeekBag,uid:%d,week:%d,level:%d,weekInfo:%v", uid, week, level, weekInfo)
	wi, err := json.Marshal(weekInfo)
	if err != nil {
		return
	}
	res, err := d.db.Exec(ctx, _addWeekBag, uid, week, level, wi)
	if err != nil {
		log.Error("AddUserGiftBag error(%v)", err)
		return
	}

	return res.LastInsertId()
}
