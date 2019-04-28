package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	//insert
	_createPromo = "insert into promotion (promo_id,type,item_id,sku_id,extra,expire_sec,sku_count,amount,buyer_count,begin_time,end_time,status,priv_sku_id,usable_coupons) values(?,?,?,?,?,?,?,?,?,?,?,1,?,?)"
	//update
	_operatePromo          = "update promotion set status = ? where promo_id = ? and status= ?"
	_updatePromoBuyerCount = "update promotion set buyer_count = buyer_count + ? where promo_id ="
	_updatePromo           = "update promotion set amount = %v,expire_sec = %v,sku_count = %v,begin_time = %v,end_time = %v,priv_sku_id = %v,usable_coupons = '%v' where promo_id = %v"
	//select
	_promo           = "select promo_id,type,item_id,sku_id,extra,expire_sec,sku_count,amount,buyer_count,begin_time,end_time,status,ctime,mtime,priv_sku_id,usable_coupons from promotion where promo_id = ?"
	_selctPromoBySKU = "select count(*) as num from promotion where end_time >= ? and begin_time <= ? and sku_id = ? and status = 1"
	_getPromoList    = "select  promo_id,type,item_id,sku_id,extra,expire_sec,sku_count,amount,buyer_count,begin_time,end_time,status,ctime,mtime,priv_sku_id,usable_coupons from promotion ? order by mtime desc limit ?,?"
)

//keyPromo 获取活动缓存key
func keyPromo(promoID int64) string {
	return fmt.Sprintf(model.CacheKeyPromo, promoID)
}

//RawPromo get promo info from db
func (d *Dao) RawPromo(c context.Context, promoID int64) (res *model.Promotion, err error) {
	res = new(model.Promotion)
	row := d.db.QueryRow(c, _promo, promoID)
	if err = row.Scan(&res.PromoID, &res.Type, &res.ItemID, &res.SKUID, &res.Extra, &res.ExpireSec, &res.SKUCount, &res.Amount, &res.BuyerCount, &res.BeginTime, &res.EndTime, &res.Status, &res.Ctime, &res.Mtime, &res.PrivSKUID, &res.UsableCoupons); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		}
		return
	}
	return
}

//CachePromo get promo info from cache
func (d *Dao) CachePromo(c context.Context, promoID int64) (res *model.Promotion, err error) {
	var (
		data []byte
		key  = keyPromo(promoID)
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

//AddCachePromo add promo info into cache
func (d *Dao) AddCachePromo(c context.Context, promoID int64, promo *model.Promotion) (err error) {
	var (
		data []byte
		key  = keyPromo(promoID)
	)

	conn := d.redis.Get(c)
	defer conn.Close()

	if data, err = json.Marshal(promo); err != nil {
		return
	}
	conn.Do("SET", key, data, "EX", model.RedisExpirePromo)
	return
}

//DelCachePromo delete promo cache
func (d *Dao) DelCachePromo(c context.Context, promoID int64) {
	var key = keyPromo(promoID)
	conn := d.redis.Get(c)
	defer conn.Close()
	conn.Do("DEL", key)
}

//CreatePromo create user promo order
func (d *Dao) CreatePromo(c context.Context, promo *model.Promotion) (id int64, err error) {
	if _, err = d.db.Exec(c, _createPromo, promo.PromoID, promo.Type, promo.ItemID, promo.SKUID, promo.Extra, promo.ExpireSec, promo.SKUCount, promo.Amount, promo.BuyerCount, promo.BeginTime, promo.EndTime, promo.PrivSKUID, promo.UsableCoupons); err != nil {
		log.Warn("创建拼团活动失败:%d", promo.PromoID)
		return
	}
	id = promo.PromoID
	return
}

//OperatePromo 修改拼团活动状态
func (d *Dao) OperatePromo(c context.Context, promoID int64, fromStatus int16, toStatus int16) (num int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, _operatePromo, toStatus, promoID, fromStatus); err != nil {
		log.Warn("更新拼团状态失败:%d", promoID)
		return
	}
	num, err = res.RowsAffected()
	return
}

//HasPromoOfSKU 判断skuID是否有正在上架的拼团活动
func (d *Dao) HasPromoOfSKU(c context.Context, skuID int64, beginTime int64, endTime int64) (num int64, err error) {
	var res *xsql.Row
	if res = d.db.QueryRow(c, _selctPromoBySKU, skuID, beginTime, endTime); err != nil {
		return
	}
	err = res.Scan(&num)
	return
}

//GetPromoList 获取拼团列表
func (d *Dao) GetPromoList(c context.Context, where string, index int64, size int64) (res []*model.Promotion, err error) {
	var rows *xsql.Rows
	if rows, err = d.db.Query(c, _getPromoList, where, index, size); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Promotion)
		if err = rows.Scan(&r.PromoID, &r.Type, &r.ItemID, &r.SKUID, &r.Extra, &r.ExpireSec, &r.SKUCount, &r.Amount, &r.BuyerCount, &r.BeginTime, &r.EndTime, &r.Status, &r.Ctime, &r.Mtime, &r.PrivSKUID, &r.UsableCoupons); err != nil {
			return
		}
		res = append(res, r)
	}
	return
}

//TxUpdatePromoBuyerCount 更新活动的购买人数
func (d *Dao) TxUpdatePromoBuyerCount(c context.Context, tx *xsql.Tx, promoID int64, count int64) (number int64, err error) {
	var res sql.Result
	if res, err = tx.Exec(_updatePromoBuyerCount, count, promoID); err != nil {
		return
	}
	return res.RowsAffected()
}

//UpdatePromo 更新拼团
func (d *Dao) UpdatePromo(c context.Context, arg *model.Promotion) (num int64, err error) {
	var res sql.Result
	if res, err = d.db.Exec(c, fmt.Sprintf(_updatePromo, arg.Amount, arg.ExpireSec, arg.SKUCount, arg.BeginTime, arg.EndTime, arg.PrivSKUID, arg.UsableCoupons, arg.PromoID)); err != nil {
		log.Warn("更新拼团活动%d失败", arg.PromoID)
		return
	}
	return res.RowsAffected()
}
