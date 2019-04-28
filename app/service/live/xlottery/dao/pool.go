package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	v1pb "go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/app/service/live/xlottery/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getPool       = "SELECT id, coin_id, title, rule, start_time, end_time, status, is_bottom FROM capsule_pool"
	_getOnePool    = "SELECT id, coin_id,title,rule,start_time,end_time,status, is_bottom FROM capsule_pool WHERE id = ?"
	_updatePStatus = "UPDATE capsule_pool SET status = ? WHERE id = ?"
	_delPool       = "DELETE FROM capsule_pool WHERE id = ?"
	_createPool    = "INSERT INTO capsule_pool (coin_id,title,rule,start_time,end_time,status) VALUES (?,?,?,?,?,0)"
	_updatePool    = "UPDATE capsule_pool SET coin_id = ?, title = ?, rule = ?, start_time = ?, end_time = ? WHERE id = ?"
	_getPoolMap    = "SELECT id, coin_id, title, rule, start_time, end_time, status, is_bottom FROM capsule_pool WHERE status = 1 and coin_id in (%v)"
)

//GetAllPool 获取奖池信息
func (d *Dao) GetAllPool(ctx context.Context) (pools []*model.Pool, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getPool); err != nil {
		log.Error("[dao.pool | GetAllPool] query(%s) error(%v)", _getPool, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := &model.Pool{}
		if err = rows.Scan(&p.Id, &p.CoinId, &p.Title, &p.Description, &p.StartTime, &p.EndTime, &p.Status, &p.IsBottom); err != nil {
			log.Error("[dao.pool | GetAllPool] scan error, err %v", err)
			return
		}
		pools = append(pools, p)
	}
	return
}

//GetPoolMap 批量奖池信息
func (d *Dao) GetPoolMap(ctx context.Context, coinIds []int64) (poolMap map[int64][]*model.Pool, err error) {
	var rows *sql.Rows
	stringCoinIds := make([]string, 0)
	for _, coinId := range coinIds {
		stringCoinIds = append(stringCoinIds, strconv.FormatInt(coinId, 10))
	}
	coinString := strings.Join(stringCoinIds, ",")
	if rows, err = d.db.Query(ctx, fmt.Sprintf(_getPoolMap, coinString)); err != nil {
		log.Error("[dao.pool | GetPoolMap]query(%s) error(%v)", _getPoolMap, err)
		return
	}
	defer rows.Close()

	poolMap = make(map[int64][]*model.Pool)
	for rows.Next() {
		d := &model.Pool{}
		if err = rows.Scan(&d.Id, &d.CoinId, &d.Title, &d.Description, &d.StartTime, &d.EndTime, &d.Status, &d.IsBottom); err != nil {
			log.Error("[dao.pool |GetPoolMap] scan error, err %v", err)
			return
		}
		if _, ok := poolMap[d.CoinId]; !ok {
			poolMap[d.CoinId] = make([]*model.Pool, 0)
		}
		poolMap[d.CoinId] = append(poolMap[d.CoinId], d)
	}
	return
}

// GetPoolById 通过ID获取奖池信息
func (d *Dao) GetPoolById(ctx context.Context, poolID int64) (pool *model.Pool, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getOnePool, poolID); err != nil {
		log.Error("[dao.pool | GetPoolById]query(%s) error(%v)", _getOnePool, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		p := &model.Pool{}
		if err = rows.Scan(&p.Id, &p.CoinId, &p.Title, &p.Description, &p.StartTime, &p.EndTime, &p.Status, &p.IsBottom); err != nil {
			log.Error("[dao.pool | GetPoolById] scan error, err %v", err)
			return
		}
		pool = p
	}
	return
}

//UpdatePoolStatus 更新奖池状态
func (d *Dao) UpdatePoolStatus(ctx context.Context, id int64, status int64) (bool, error) {
	//s := fmt.Sprintf(_updatePStatus, status, id)
	res, err := d.db.Exec(ctx, _updatePStatus, status, id)
	if err != nil {
		log.Error("[dao.pool | UpdatePoolStatus] query(%s) error(%v)", _updatePStatus, err)
		return false, err
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.pool | UpdatePoolStatus] no coin is affected, err %v", err)
		return false, err
	}
	return rows > 0, nil
}

//DelPool 删除奖池
func (d *Dao) DelPool(ctx context.Context, id int64) (rows int64, err error) {
	res, err := d.db.Exec(ctx, _delPool, id)
	if err != nil {
		log.Error("[dao.pool | DelPool]query(%s) error(%v)", _delPool, err)
		return
	}
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.pool | DelPool]no pool is deleted, err %v", err)
		return
	}
	return
}

//CreatePool 新建奖池
func (d *Dao) CreatePool(ctx context.Context, data *v1pb.UpdatePoolReq) (status bool, err error) {
	coinId := data.CoinId
	title := data.Title
	startTime := data.StartTime
	endTime := data.EndTime
	rule := data.Rule
	//s := fmt.Sprintf(_createPool,coinTitle,title,rule,startTime,endTime)
	res, err := d.db.Exec(ctx, _createPool, coinId, title, rule, startTime, endTime)
	if err != nil {
		log.Error("[dao.pool | CreatePool] query(%s) error(%v)", _createPool, err)
		return
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.pool | CreatePool] no pool is created, err %v", err)
		return
	}
	status = rows > 0
	return
}

//UpdatePool 更新奖池
func (d *Dao) UpdatePool(ctx context.Context, data *v1pb.UpdatePoolReq) (status bool, err error) {
	id := data.Id
	coinId := data.CoinId
	title := data.Title
	startTime := data.StartTime
	endTime := data.EndTime
	rule := data.Rule
	res, err := d.db.Exec(ctx, _updatePool, coinId, title, rule, startTime, endTime, id)
	if err != nil {
		log.Error("[dao.pool | UpdatePool] query(%s) error(%v)", _updatePool, err)
		return
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.pool | UpdatePool] no pool is updated, err %v", err)
		return
	}
	status = rows > 0
	return
}
