package dao

import (
	"context"

	v1pb "go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/app/service/live/xlottery/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getAllCoin   = "SELECT id, title, gift_type, change_num, start_time, end_time, status FROM capsule_coin"
	_updateStatus = "UPDATE capsule_coin SET status = ? WHERE id = ?"
	_delCoin      = "DELETE FROM capsule_coin WHERE id = ?"
	_createCoin   = "INSERT INTO capsule_coin (title, gift_type, change_num, start_time, end_time, status) VALUES (?,?,?,?,?,?)"
	_updateCoin   = "UPDATE capsule_coin SET title = ?, gift_type = ?, change_num = ?, start_time = ?, end_time = ?, status = ? WHERE id = ?"
	_getOnCoin    = "SELECT id, title, gift_type, change_num, start_time, end_time, status FROM capsule_coin WHERE status = 1"
)

//GetAllCoin 获取所有扭蛋币
func (d *Dao) GetAllCoin(ctx context.Context) (coins []*model.Coin, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getAllCoin); err != nil {
		log.Error("[dao.coin | GetALLCoin] query(%s) error (%v)", _getAllCoin, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := &model.Coin{}
		if err = rows.Scan(&p.Id, &p.Title, &p.GiftType, &p.ChangeNum, &p.StartTime, &p.EndTime, &p.Status); err != nil {
			log.Error("[dao.coin | GetAllCoin] scan error, err %v", err)
			return
		}
		coins = append(coins, p)
	}
	return
}

//GetCoinMap 批量获取扭蛋币
func (d *Dao) GetCoinMap(ctx context.Context) (coinMap map[int64]*model.Coin, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getOnCoin); err != nil {
		log.Error("[dao.coin | GetCoinMap] query(%s) error(%v)", _getOnCoin, err)
		return
	}
	defer rows.Close()

	coinMap = make(map[int64]*model.Coin)
	for rows.Next() {
		p := &model.Coin{}
		if err = rows.Scan(&p.Id, &p.Title, &p.GiftType, &p.ChangeNum, &p.StartTime, &p.EndTime, &p.Status); err != nil {
			log.Error("[dao.coin | GetCoinMap] scan error, err %v", err)
			return
		}
		coinMap[p.Id] = p
	}
	return
}

//UpdateCoinStatus 更新扭蛋币状态
func (d *Dao) UpdateCoinStatus(ctx context.Context, id int64, status int64) (bool, error) {
	res, err := d.db.Exec(ctx, _updateStatus, status, id)
	if err != nil {
		log.Error("[dao.coin | UpdateCStatus] query(%s) error(%v)", _updateStatus, err)
		return false, err
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.coin | UpdateCStatus] no coin is affected, err %v", err)
		return false, err
	}
	return rows > 0, nil
}

//DeleteCoin 删除扭蛋币
func (d *Dao) DeleteCoin(ctx context.Context, id int64) (status bool, err error) {
	res, err := d.db.Exec(ctx, _delCoin, id)
	if err != nil {
		log.Error("[dao.coin | DelCoin] query(%s), error(%v)", _delCoin, err)
		return
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.coin | DelCoin] no coin is deleted, err %v", err)
		return
	}
	status = rows > 0
	return
}

// CreateCoin 更新扭蛋币
func (d *Dao) CreateCoin(ctx context.Context, data *v1pb.UpdateCoinConfigReq) (ret int64, err error) {
	title := data.Title
	giftType := data.GiftType
	changeNum := data.ChangeNum
	startTime := data.StartTime
	endTime := data.EndTime
	status := 1
	res, err := d.db.Exec(ctx, _createCoin, title, giftType, changeNum, startTime, endTime, status)
	if err != nil {
		log.Error("[dao.coin | CreateCoin] query(%s) error (%v)", _createCoin, err)
		return
	}
	ret, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.coin | CreateCoin]no coin is added, err %v", err)
		return
	}
	return
}

// UpdateCoin 更新扭蛋币
func (d *Dao) UpdateCoin(ctx context.Context, data *v1pb.UpdateCoinConfigReq) (status bool, err error) {
	id := data.Id
	title := data.Title
	giftType := data.GiftType
	changeNum := data.ChangeNum
	startTime := data.StartTime
	endTime := data.EndTime
	_, err = d.db.Exec(ctx, _updateCoin, title, giftType, changeNum, startTime, endTime, 1, id)
	if err != nil {
		log.Error("[dao.coin | UpdateCoin] query(%s) error %v", _updateCoin, err)
		return
	}
	return true, nil
}
