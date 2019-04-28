package dao

import (
	"context"
	"go-common/app/service/live/xlottery/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addPrizeData     = "INSERT INTO capsule_prize_log (prize_id,uid,day,timestamp) VALUES (?,?,?,?)"
	_getUserPrizeData = "SELECT prize_id,uid,day,`timestamp` from capsule_prize_log where prize_id = ? and uid = ? order by id desc"
	_getPrizeDayData  = "SELECT prize_id,uid,day,`timestamp` from capsule_prize_log where prize_id = ? and day = ? order by id desc"
)

// AddPrizeData 添加特殊奖品获奖记录
func (d *Dao) AddPrizeData(ctx context.Context, prizeId, uid int64, day string, timestamp int64) (status bool, err error) {
	res, err := d.db.Exec(ctx, _addPrizeData, prizeId, uid, day, timestamp)
	if err != nil {
		log.Error("[dao.prize_extra | AddWhiteUser] add(%s) error (%v)", _addPrizeData, err)
		return
	}
	rows, _ := res.RowsAffected()
	status = rows > 0
	return
}

// GetUserPrizeLog 获取特殊奖品记录
func (d *Dao) GetUserPrizeLog(ctx context.Context, prizeId int64, uid int64) (prizeLog *model.PrizeLog, err error) {
	row := d.db.QueryRow(ctx, _getUserPrizeData, prizeId, uid)
	prizeLog = &model.PrizeLog{}
	err = row.Scan(&prizeLog.PrizeId, &prizeLog.Uid, &prizeLog.Day, &prizeLog.Timestamp)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error("[dao.pool_prize | GetUserPrizeLog]query(%s) error(%v)", _getUserPrizeData, err)
		return
	}
	return prizeLog, nil
}

// GetPrizeDayLog 获取特殊奖品记录
func (d *Dao) GetPrizeDayLog(ctx context.Context, prizeId int64, day string) (prizeLog *model.PrizeLog, err error) {
	row := d.db.QueryRow(ctx, _getPrizeDayData, prizeId, day)
	prizeLog = &model.PrizeLog{}
	err = row.Scan(&prizeLog.PrizeId, &prizeLog.Uid, &prizeLog.Day, &prizeLog.Timestamp)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error("[dao.pool_prize | GetPrizeDayLog]query(%s) error(%v)", _getPrizeDayData, err)
		return
	}
	return
}
