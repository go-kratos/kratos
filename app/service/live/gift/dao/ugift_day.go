package dao

import (
	"context"
	"encoding/json"
	"go-common/app/service/live/gift/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

var (
	_getBagStatus = "SELECT id,uid,day,day_info FROM ugift_day_status WHERE uid = ? AND day = ? LIMIT 1"
	_addDayBag    = "INSERT INTO ugift_day_status (uid,day,day_info) VALUES (?,?,?)"
)

// GetDayBagStatus GetDayBagStatus
func (d *Dao) GetDayBagStatus(ctx context.Context, uid int64, date string) (res *model.DayGiftInfo, err error) {
	log.Info("GetDayBagStatus,%d,%s", uid, date)
	row := d.db.QueryRow(ctx, _getBagStatus, uid, date)
	res = &model.DayGiftInfo{}
	if err = row.Scan(&res.ID, &res.UID, &res.Day, &res.DayInfo); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("GetUserGiftBagStatus row.Scan error(%v)", err)
	}

	return
}

// AddDayBag AddDayBag
func (d *Dao) AddDayBag(ctx context.Context, uid int64, date string, dayInfo *model.BagGiftStatus) (affected int64, err error) {
	log.Info("AddDayBag,%d,%s,%v", uid, date, dayInfo)
	di, err := json.Marshal(dayInfo)
	if err != nil {
		return
	}
	res, err := d.db.Exec(ctx, _addDayBag, uid, date, di)
	if err != nil {
		log.Error("AddUserGiftBag error(%v)", err)
		return
	}

	return res.LastInsertId()
}
