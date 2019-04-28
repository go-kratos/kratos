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
	_getPoolPrize    = "SELECT id, pool_id, type, num, object_id,expire, web_url, mobile_url, description, jump_url, pro_type, chance, loop_num, limit_num, weight FROM capsule_pool_prize WHERE pool_id = ? and status = 1 order by ctime"
	_delPoolPrize    = "UPDATE capsule_pool_prize set status = 0 WHERE id = ?"
	_addPoolPrize    = "INSERT INTO capsule_pool_prize(pool_id,type,num,object_id,expire,web_url,mobile_url,description,jump_url,pro_type,chance,loop_num,limit_num,weight,status) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_updatePoolPrize = "UPDATE capsule_pool_prize SET pool_id=?, type=?, num=?, object_id=?,expire=?, web_url=?, mobile_url=?, description=?, jump_url=?, pro_type=?, chance=?, loop_num=?, limit_num=?,weight=? WHERE id = ?"
	_getPoolPrizeMap = "SELECT id, pool_id, type, num, object_id,expire, web_url, mobile_url, description, jump_url, pro_type, chance, loop_num, limit_num, weight FROM capsule_pool_prize WHERE pool_id in (%s) and status = 1 order by ctime"
)

// GetPoolPrize 获取奖池奖品
func (d *Dao) GetPoolPrize(ctx context.Context, poolID int64) (poolPrize []*model.PoolPrize, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getPoolPrize, poolID); err != nil {
		log.Error("[dao.pool_prize | GetPoolPrize]query(%s) error(%v)", _getPoolPrize, err)
		return
	}
	defer rows.Close()

	poolPrize = make([]*model.PoolPrize, 0)
	for rows.Next() {
		p := &model.PoolPrize{}
		if err = rows.Scan(&p.Id, &p.PoolId, &p.Type, &p.Num, &p.ObjectId, &p.Expire, &p.WebUrl, &p.MobileUrl, &p.Description, &p.JumpUrl, &p.ProType, &p.Chance, &p.LoopNum, &p.LimitNum, &p.Weight); err != nil {
			log.Error("[dao.pool_prize | GetPoolPrize] scan error, err %v", err)
			return
		}
		poolPrize = append(poolPrize, p)
	}
	return
}

// GetPoolPrizeMap 批量奖池奖品
func (d *Dao) GetPoolPrizeMap(ctx context.Context, poolIds []int64) (poolPrizeMap map[int64][]*model.PoolPrize, err error) {
	var rows *sql.Rows
	stringPoolIds := make([]string, 0)
	for _, poolId := range poolIds {
		stringPoolIds = append(stringPoolIds, strconv.FormatInt(poolId, 10))
	}
	poolString := strings.Join(stringPoolIds, ",")
	if rows, err = d.db.Query(ctx, fmt.Sprintf(_getPoolPrizeMap, poolString)); err != nil {
		log.Error("[dao.pool_prize | GetPoolPrizeMap] query(%s) error(%v)", _getPoolPrizeMap, err)
		return
	}
	defer rows.Close()

	poolPrizeMap = make(map[int64][]*model.PoolPrize)
	for rows.Next() {
		d := &model.PoolPrize{}
		if err = rows.Scan(&d.Id, &d.PoolId, &d.Type, &d.Num, &d.ObjectId, &d.Expire, &d.WebUrl, &d.MobileUrl, &d.Description, &d.JumpUrl, &d.ProType, &d.Chance, &d.LoopNum, &d.LimitNum, &d.Weight); err != nil {
			log.Error("[dao.pool_prize | GetPoolPrizeMap] scan error, err %v", err)
			return
		}
		if _, ok := PrizeNameMap[d.Type]; !ok {
			continue
		}
		if _, ok := poolPrizeMap[d.PoolId]; !ok {
			poolPrizeMap[d.PoolId] = make([]*model.PoolPrize, 0)
		}
		poolPrizeMap[d.PoolId] = append(poolPrizeMap[d.PoolId], d)
	}
	return
}

// DeletePoolPrize 删除奖池奖品
func (d *Dao) DeletePoolPrize(ctx context.Context, id int64) (status bool, err error) {
	res, err := d.db.Exec(ctx, _delPoolPrize, id)
	if err != nil {
		log.Error("[dao.pool_prize | DelPoolDetail] query(%s) error (%v)", _delPoolPrize, err)
		return
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Error("[dao.pool_prize | DelPoolDetail] no poolDetail is deleted, err %v", err)
		return
	}
	status = rows > 0
	return
}

// AddPoolPrize 创建奖池奖品
func (d *Dao) AddPoolPrize(ctx context.Context, data *v1pb.UpdatePoolPrizeReq) (prizeId int64, err error) {
	res, err := d.db.Exec(ctx, _addPoolPrize, data.PoolId, data.Type, data.Num, data.ObjectId, data.Expire, data.WebUrl, data.MobileUrl, data.Description, data.JumpUrl, data.ProType, data.Chance, data.Loop, data.Limit, data.Weight, 1)
	if err != nil {
		log.Error("[dao.pool_prize | CreatePoolPrize] query(%s) error (%v)", _addPoolPrize, err)
		return
	}
	prizeId, err = res.LastInsertId()
	if err != nil {
		log.Error("[dao.pool_prize | CreatePoolPrize] insertLastId error(%v)", err)
		return
	}
	return
}

// UpdatePoolPrize 更新奖池奖品
func (d *Dao) UpdatePoolPrize(ctx context.Context, data *v1pb.UpdatePoolPrizeReq) (status bool, err error) {
	_, err = d.db.Exec(ctx, _updatePoolPrize, data.PoolId, data.Type, data.Num, data.ObjectId, data.Expire, data.WebUrl, data.MobileUrl, data.Description, data.JumpUrl, data.ProType, data.Chance, data.Loop, data.Limit, data.Weight, data.Id)
	if err != nil {
		log.Error("[dao.pool_prize | UpdatePoolPrize] query(%s) error(%v)", _updatePoolPrize, err)
		return
	}
	return true, nil
}
