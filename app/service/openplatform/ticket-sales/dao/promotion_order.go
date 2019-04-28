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
	"go-common/library/log"
	xtime "time"
)

const (
	_addPromoOrder                    = "insert into promotion_order (promo_id,group_id,order_id,is_master,uid,status,ctime,sku_id) values(?,?,?,?,?,?,?,?)"
	_updatePromoOrderStatus           = "update promotion_order set status = ? where order_id = ?"
	_updatePromoOrderGroupIDAndStatus = "update promotion_order set group_id = ?,status = ? where order_id = ?"
	_promoOrderByStatus               = "select promo_id,group_id,order_id,is_master,uid,status,ctime,mtime,sku_id from promotion_order where promo_id = ? and group_id = ? and uid = ? and status = ?"
	_promoOrderDoing                  = "select promo_id,group_id,order_id,is_master,uid,status,ctime,mtime,sku_id from promotion_order where promo_id = ? and group_id = ? and uid = ? and status in (?,?)"
	_groupOrdersByStatus              = "select promo_id,group_id,order_id,is_master,uid,status,ctime,mtime,sku_id from promotion_order where group_id = ? and status = ?"
	_promoOrdersByGroupID             = "select promo_id,group_id,order_id,is_master,uid,status,ctime,mtime,sku_id from promotion_order where group_id = ? and status in (?,?,?)"
	_promoOrder                       = "select promo_id,group_id,order_id,is_master,uid,status,ctime,mtime,sku_id from promotion_order where order_id = ?"
)

//keyPromoOrder 获取拼团订单缓存key
func keyPromoOrder(orderID int64) string {
	return fmt.Sprintf(model.CacheKeyPromoOrder, orderID)
}

//RawPromoOrder get promo order info from db
func (d *Dao) RawPromoOrder(c context.Context, orderID int64) (res *model.PromotionOrder, err error) {
	res = new(model.PromotionOrder)
	row := d.db.QueryRow(c, _promoOrder, orderID)
	if err = row.Scan(&res.PromoID, &res.GroupID, &res.OrderID, &res.IsMaster, &res.UID, &res.Status, &res.Ctime, &res.Mtime, &res.SKUID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		}
		return
	}
	return
}

//CachePromoOrder get promo order info from cache
func (d *Dao) CachePromoOrder(c context.Context, orderID int64) (res *model.PromotionOrder, err error) {
	var (
		data []byte
		key  = keyPromoOrder(orderID)
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

//AddCachePromoOrder add promo order info into cache
func (d *Dao) AddCachePromoOrder(c context.Context, orderID int64, promoOrder *model.PromotionOrder) (err error) {
	var (
		data []byte
		key  = keyPromoOrder(orderID)
	)

	conn := d.redis.Get(c)
	defer conn.Close()

	if data, err = json.Marshal(promoOrder); err != nil {
		return
	}

	conn.Do("SET", key, data, "EX", model.RedisExpirePromoOrder)
	return
}

//DelCachePromoOrder delete promo order cache
func (d *Dao) DelCachePromoOrder(c context.Context, orderID int64) {
	var key = keyPromoOrder(orderID)
	conn := d.redis.Get(c)
	defer conn.Close()
	conn.Do("DEL", key)
}

//keyPromoOrders 获取拼团团订单缓存key
func keyPromoOrders(groupID int64) string {
	return fmt.Sprintf(model.CacheKeyPromoOrders, groupID)
}

//RawPromoOrders get promo orders info from db
func (d *Dao) RawPromoOrders(c context.Context, groupID int64) (res []*model.PromotionOrder, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _promoOrdersByGroupID, groupID, consts.PromoOrderUnpaid, consts.PromoOrderPaid, consts.PromoOrderRefund); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.PromotionOrder)
		if err = rows.Scan(&r.PromoID, &r.GroupID, &r.OrderID, &r.IsMaster, &r.UID, &r.Status, &r.Ctime, &r.Mtime, &r.SKUID); err != nil {
			return
		}
		res = append(res, r)
	}
	return
}

//CachePromoOrders get promo orders info from cache
func (d *Dao) CachePromoOrders(c context.Context, groupID int64) (res []*model.PromotionOrder, err error) {
	var (
		data []byte
		key  = keyPromoOrders(groupID)
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

//AddCachePromoOrders add promo orders info into cache
func (d *Dao) AddCachePromoOrders(c context.Context, groupID int64, promoOrders []*model.PromotionOrder) (err error) {
	var (
		data []byte
		key  = keyPromoOrders(groupID)
	)

	conn := d.redis.Get(c)
	defer conn.Close()

	if data, err = json.Marshal(promoOrders); err != nil {
		return
	}

	conn.Do("SET", key, data, "EX", model.RedisExpirePromoOrders)
	return
}

//DelCachePromoOrders delete promo orders cache
func (d *Dao) DelCachePromoOrders(c context.Context, groupID int64) {
	var key = keyPromoOrders(groupID)
	conn := d.redis.Get(c)
	defer conn.Close()
	conn.Do("DEL", key)
}

//PromoOrderByStatus get user promo order by status
func (d *Dao) PromoOrderByStatus(c context.Context, promoID int64, groupID int64, uid int64, status int16) (res *model.PromotionOrder, err error) {
	res = new(model.PromotionOrder)
	row := d.db.QueryRow(c, _promoOrderByStatus, promoID, groupID, uid, status)
	if err = row.Scan(&res.PromoID, &res.GroupID, &res.OrderID, &res.IsMaster, &res.UID, &res.Status, &res.Ctime, &res.Mtime, &res.SKUID); err != nil {
		return
	}
	return
}

//PromoOrderDoing get user promo order where status in unpaid and paid
func (d *Dao) PromoOrderDoing(c context.Context, promoID int64, groupID int64, uid int64) (res *model.PromotionOrder, err error) {
	res = new(model.PromotionOrder)
	row := d.db.QueryRow(c, _promoOrderDoing, promoID, groupID, uid, consts.PromoOrderUnpaid, consts.PromoOrderPaid)
	if err = row.Scan(&res.PromoID, &res.GroupID, &res.OrderID, &res.IsMaster, &res.UID, &res.Status, &res.Ctime, &res.Mtime, &res.SKUID); err != nil {
		return
	}
	return
}

//AddPromoOrder create user promo order
func (d *Dao) AddPromoOrder(c context.Context, promoID int64, groupID int64, orderID int64, isMaster int16, uid int64, status int16, skuID int64, ctime int64) (id int64, err error) {
	var (
		res       sql.Result
		ctimeDate = xtime.Unix(ctime, 0)
	)
	if res, err = d.db.Exec(c, _addPromoOrder, promoID, groupID, orderID, isMaster, uid, status, ctimeDate, skuID); err != nil {
		log.Warn("创建活动订单%d失败", orderID)
		return
	}
	return res.LastInsertId()
}

//TxAddPromoOrder create user promo order
func (d *Dao) TxAddPromoOrder(c context.Context, tx *xsql.Tx, promoID int64, groupID int64, orderID int64, isMaster int16, uid int64, status int16, skuID int64, ctime int64) (id int64, err error) {
	var (
		res       sql.Result
		ctimeDate = xtime.Unix(ctime, 0)
	)
	if res, err = tx.Exec(_addPromoOrder, promoID, groupID, orderID, isMaster, uid, status, ctimeDate, skuID); err != nil {
		log.Warn("创建活动订单%d失败", orderID)
		return
	}
	return res.LastInsertId()
}

//UpdatePromoOrderStatus update promo order status
func (d *Dao) UpdatePromoOrderStatus(c context.Context, orderID int64, status int16) (number int64, err error) {
	var (
		res sql.Result
	)
	if res, err = d.db.Exec(c, _updatePromoOrderStatus, status, orderID); err != nil {
		log.Warn("更新活动订单%d失败", orderID)
		return
	}
	return res.RowsAffected()
}

//TxUpdatePromoOrderStatus update promo order status
func (d *Dao) TxUpdatePromoOrderStatus(c context.Context, tx *xsql.Tx, orderID int64, status int16) (number int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updatePromoOrderStatus, status, orderID); err != nil {
		log.Warn("更新活动订单%d失败", orderID)
		return
	}
	return res.RowsAffected()
}

//TxUpdatePromoOrderGroupIDAndStatus update promo order groupid and status
func (d *Dao) TxUpdatePromoOrderGroupIDAndStatus(c context.Context, tx *xsql.Tx, orderID int64, groupID int64, status int16) (number int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updatePromoOrderGroupIDAndStatus, groupID, status, orderID); err != nil {
		log.Warn("更新活动订单%d失败", orderID)
		return
	}
	return res.RowsAffected()
}

//GroupOrdersByStatus 根据groupid和status获取活动订单
func (d *Dao) GroupOrdersByStatus(c context.Context, groupID int64, status int16) (res []*model.PromotionOrder, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _groupOrdersByStatus, groupID, status); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		temp := new(model.PromotionOrder)
		if err = rows.Scan(&temp.PromoID, &temp.GroupID, &temp.OrderID, &temp.IsMaster, &temp.UID, &temp.Status, &temp.Ctime, &temp.Mtime, &temp.SKUID); err != nil {
			return
		}
		res = append(res, temp)
	}
	return
}
