package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	xtime "time"
)

const (
	_getGroupByID                   = "select promo_id,group_id,uid,order_count,status,expire_at,ctime,mtime from promotion_group where group_id = ?"
	_getUserNotExpiredGroup         = "select promo_id,group_id,uid,order_count,status,expire_at,ctime,mtime from promotion_group where promo_id = ? and uid = ? and expire_at >= ? and status = ?"
	_updateGroupOrderCount          = "update promotion_group set order_count = order_count + ? where group_id = ? and status = ? and order_count < ? and expire_at > ?"
	_updateGroupStatus              = "update promotion_group set status = ? where group_id = ? and status = ?"
	_updateGroupStatusAndOrderCount = "update promotion_group set status = ?,order_count = order_count + ? where group_id = ? and status = ?"
	_insertGroupOrder               = "insert into promotion_group (promo_id,group_id,uid,order_count,status,expire_at) values (?,?,?,?,?,?)"
)

//keyPromoGroup 获取拼团缓存key
func keyPromoGroup(groupID int64) string {
	return fmt.Sprintf(model.CacheKeyPromoGroup, groupID)
}

//RawPromoGroup 根据id获取拼团信息
func (d *Dao) RawPromoGroup(c context.Context, groupID int64) (res *model.PromotionGroup, err error) {
	res = new(model.PromotionGroup)
	row := d.db.QueryRow(c, _getGroupByID, groupID)
	if err = row.Scan(&res.PromoID, &res.GroupID, &res.UID, &res.OrderCount, &res.Status, &res.ExpireAt, &res.Ctime, &res.Mtime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		}
		return
	}
	return
}

//AddPromoGroup add promo group into db
func (d *Dao) AddPromoGroup(c context.Context, promoID int64, groupID int64, uid int64, orderCount int64, status int16, expireAt int64) (id int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _insertGroupOrder, promoID, groupID, uid, orderCount, status, expireAt); err != nil {
		return
	}
	return res.LastInsertId()
}

//TxAddPromoGroup add promo group into db
func (d *Dao) TxAddPromoGroup(c context.Context, tx *xsql.Tx, promoID int64, groupID int64, uid int64, orderCount int64, status int16, expireAt int64) (id int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_insertGroupOrder, promoID, groupID, uid, orderCount, status, expireAt); err != nil {
		return
	}
	return res.LastInsertId()
}

//CachePromoGroup get promo group info from cache
func (d *Dao) CachePromoGroup(c context.Context, groupID int64) (res *model.PromotionGroup, err error) {
	var (
		data []byte
		key  = keyPromoGroup(groupID)
	)

	conn := d.redis.Get(c)
	defer conn.Close()

	if data, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	json.Unmarshal(data, &res)
	return
}

//AddCachePromoGroup add promo group info into cache
func (d *Dao) AddCachePromoGroup(c context.Context, groupID int64, group *model.PromotionGroup) (err error) {
	var (
		data []byte
		key  = keyPromoGroup(groupID)
	)

	conn := d.redis.Get(c)
	defer conn.Close()

	if data, err = json.Marshal(group); err != nil {
		return
	}
	conn.Do("SET", key, data, "EX", model.RedisExpirePromoGroup)
	return
}

//DelCachePromoGroup delete promo group cache
func (d *Dao) DelCachePromoGroup(c context.Context, groupID int64) {
	var key = keyPromo(groupID)
	conn := d.redis.Get(c)
	defer conn.Close()
	conn.Do("DEL", key)
}

//GetUserGroupDoing 获取用户正在进行中的拼团信息
func (d *Dao) GetUserGroupDoing(c context.Context, promoID int64, uid int64, status int16) (res *model.PromotionGroup, err error) {
	var (
		currentTime = xtime.Now().Unix()
	)

	res = new(model.PromotionGroup)
	row := d.db.QueryRow(c, _getUserNotExpiredGroup, promoID, uid, currentTime, status)
	if err = row.Scan(&res.PromoID, &res.GroupID, &res.UID, &res.OrderCount, &res.Status, &res.ExpireAt, &res.Ctime, &res.Mtime); err != nil {
		return
	}
	return
}

//TxUpdateGroupOrderCount 更新拼团的人数
func (d *Dao) TxUpdateGroupOrderCount(c context.Context, tx *xsql.Tx, step int64, groupID int64, skuCount int64) (number int64, err error) {
	var (
		currentTime = xtime.Now().Unix()
		res         sql.Result
	)
	if res, err = tx.Exec(_updateGroupOrderCount, step, groupID, consts.GroupDoing, skuCount, currentTime); err != nil {
		return
	}
	return res.RowsAffected()
}

//TxUpdateGroupStatus 更新拼团状态
func (d *Dao) TxUpdateGroupStatus(c context.Context, tx *xsql.Tx, groupID int64, oldStatus int16, newStatus int16) (number int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updateGroupStatus, newStatus, groupID, oldStatus); err != nil {
		return
	}
	return res.RowsAffected()
}

//UpdateGroupStatusAndOrderCount 更新拼团状态和人数
func (d *Dao) UpdateGroupStatusAndOrderCount(c context.Context, groupID int64, step int64, oldStatus int16, newStatus int16) (number int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _updateGroupStatusAndOrderCount, newStatus, step, groupID, oldStatus); err != nil {
		return
	}
	return res.RowsAffected()
}
