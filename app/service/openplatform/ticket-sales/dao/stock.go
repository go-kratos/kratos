package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_lockStockSQL        = "UPDATE sku_stock SET stock=stock-?, locked_stock=locked_stock+? WHERE sku_id=? AND stock>=?"
	_unlockStockSQL      = "UPDATE sku_stock SET stock=stock+?, locked_stock=locked_stock-? WHERE sku_id=? AND stock<=total_stock-? AND locked_stock>=?"
	_decrStockSQL        = "UPDATE sku_stock SET stock=stock-? WHERE sku_id=? AND stock>=?"
	_incrStockSQL        = "UPDATE sku_stock SET stock=stock+? WHERE sku_id=? AND stock<=total_stock-?"
	_decrStockLockedSQL  = "UPDATE sku_stock SET locked_stock=locked_stock-? WHERE sku_id=? AND total_stock>=?"
	_getStockBySkuIDSQL  = "SELECT sku_id, parent_sku_id, item_id, specs, total_stock, stock, locked_stock, sk_alert, ctime, mtime FROM sku_stock WHERE sku_id=? LIMIT 1"
	_getStocksBySkuIDSQL = "SELECT sku_id, parent_sku_id, item_id, specs, total_stock, stock, locked_stock, sk_alert, ctime, mtime FROM sku_stock WHERE sku_id IN (%s)"
	_getStockByItemIDSQL = "SELECT sku_id, parent_sku_id, item_id, specs, total_stock, stock, locked_stock, sk_alert, ctime, mtime FROM sku_stock WHERE item_id=?"
	_getStockBySpecsSQL  = "SELECT sku_id, parent_sku_id, item_id, specs, total_stock, stock, locked_stock, sk_alert, ctime, mtime FROM sku_stock WHERE item_id=? AND specs=? LIMIT 1"
	_insertStockSQL      = "INSERT INTO sku_stock (sku_id, parent_sku_id, item_id, specs, total_stock, stock, locked_stock, sk_alert) VALUES %s"
	_resetStockSQL       = "UPDATE sku_stock SET stock=stock+(?-total_stock),total_stock=?, sk_alert=? WHERE sku_id=? AND stock>=(total_stock-?)"
	_getStockBySkuID     = "SELECT sku_id, stock, locked_stock FROM sku_stock WHERE sku_id IN (%s)"

	_insertSKUStockLog   = "INSERT INTO sku_stock_log (sku_id, op_type, src_id, stock) VALUES %s"
	_selectSKUStockLog   = "SELECT id, sku_id, op_type, src_id, stock FROM sku_stock_log WHERE src_id=? AND op_type=? AND sku_id IN(%s)"
	_rollbackSKUStockLog = "UPDATE sku_stock_log set canceled_at=? where id=? AND canceled_at=0"
)

// StockLock lock stock
func (d *Dao) StockLock(c context.Context, skuID int64, cnt int64) (err error) {
	_, err = d.db.Exec(c, _lockStockSQL, cnt, cnt, skuID, cnt)
	if err != nil {
		log.Error("d.StockLock(%d, %d) error(%v)", skuID, cnt, err)
	}
	return
}

// TxStockLock lock stock with tx
func (d *Dao) TxStockLock(tx *xsql.Tx, skuID int64, cnt int64) (affected int64, err error) {
	res, err := tx.Exec(_lockStockSQL, cnt, cnt, skuID, cnt)
	if err != nil {
		log.Error("d.TxStockLock(%d, %d) error(%v)", skuID, cnt, err)
		return
	}
	affected, err = res.RowsAffected()
	if err != nil {
		log.Error("d.TxStockLock(%d, %d) res.RowsAffected() error(%v)", skuID, cnt, err)
		return
	}
	return
}

// StockDecr 减库存 DB
func (d *Dao) StockDecr(c context.Context, skuID int64, num int64) (affected int64, err error) {
	res, err := d.db.Exec(c, _decrStockSQL, num, skuID, num)
	if err != nil {
		log.Error("d.StockDecr(%d, %d) error(%v)", skuID, num, err)
		return
	}
	affected, err = res.RowsAffected()
	if err != nil {
		log.Error("d.StockDecr(%d, %d) res.RowsAffected() error(%v)", skuID, num, err)
		return
	}
	return
}

// TxStockDecr 减库存 DB
func (d *Dao) TxStockDecr(tx *xsql.Tx, skuID int64, num int64) (affected int64, err error) {
	res, err := tx.Exec(_decrStockSQL, num, skuID, num)
	if err != nil {
		log.Error("d.TxStockDecr(%d, %d) error(%v)", skuID, num, err)
		return
	}
	affected, err = res.RowsAffected()
	if err != nil {
		log.Error("d.TxStockDecr(%d, %d) res.RowsAffected() error(%v)", skuID, num, err)
		return
	}
	return
}

// StockCacheDecr 减库存缓存
func (d *Dao) StockCacheDecr(c context.Context, skuID int64, total int64) (err error) {
	if err = d.RedisDecrExist(c, fmt.Sprintf(model.CacheKeyStock, skuID), total); err != nil {
		log.Error("d.StockCacheDecr(%d) error(%v)", skuID, err)
	}
	return
}

// StockLockedCacheDecr 减锁定库存缓存
func (d *Dao) StockLockedCacheDecr(c context.Context, skuID int64, total int64) (err error) {
	if err = d.RedisDecrExist(c, fmt.Sprintf(model.CacheKeyStockL, skuID), total); err != nil {
		log.Error("d.StockLockedCacheDecr(%d) error(%v)", skuID, err)
	}
	return
}

// StockCacheDel 删除库存缓存
func (d *Dao) StockCacheDel(c context.Context, skuID int64) (err error) {
	if err = d.RedisDel(c, fmt.Sprintf(model.CacheKeyStock, skuID)); err != nil {
		log.Error("d.StockCacheDel(%d) error(%v)", skuID, err)
	}
	return
}

// DelCacheSku 删除 skuId => sku 缓存
func (d *Dao) DelCacheSku(c context.Context, skuID int64) (err error) {
	if err = d.RedisDel(c, fmt.Sprintf(model.CacheKeySku, skuID)); err != nil {
		log.Error("d.StockCacheDel(%d) error(%v)", skuID, err)
	}
	return
}

// AddStockLog TxAddStockLog
func (d *Dao) AddStockLog() {

}

// TxAddStockLog 添加库存操作日志
func (d *Dao) TxAddStockLog(tx *xsql.Tx, stockLogs ...*model.SKUStockLog) (err error) {
	if len(stockLogs) == 0 {
		return
	}

	placeholder := strings.Trim(strings.Repeat("(?, ?, ?, ?),", len(stockLogs)), ",")
	var values []interface{}
	for _, stockLog := range stockLogs {
		values = append(values, stockLog.SKUID, stockLog.OpType, stockLog.SrcID, stockLog.Stock)
	}

	if _, err = tx.Exec(fmt.Sprintf(_insertSKUStockLog, placeholder), values...); err != nil {
		log.Error("d.TxAddStockLog() error(%v)", err)
		return
	}
	return
}

// TxStockUnlock 解锁库存（减去锁定库存增加库存）
func (d *Dao) TxStockUnlock(tx *xsql.Tx, skuID int64, count int64) (affected int64, err error) {
	res, err := tx.Exec(_unlockStockSQL, count, count, skuID, count, count)
	if err != nil {
		log.Error("d.TxStockUnlock(%d, %d) error(%v)", skuID, count, err)
		return
	}
	if affected, err = res.RowsAffected(); err != nil {
		log.Error("d.TxStockUnlock(%d, %d) res.RowsAffected() error(%v)", skuID, count, err)
	}
	return
}

// TxStockIncr 增加库存
func (d *Dao) TxStockIncr(tx *xsql.Tx, skuID int64, count int64) (affected int64, err error) {
	res, err := tx.Exec(_incrStockSQL, count, skuID, count)
	if err != nil {
		log.Error("d.TxStockIncr(%d, %d) error(%v)", skuID, count, err)
		return
	}
	if affected, err = res.RowsAffected(); err != nil {
		log.Error("d.TxStockIncr(%d, %d) res.RowsAffected() error(%v)", skuID, count, err)
	}
	return
}

// TxStockLockedDecr 减去锁定库存
func (d *Dao) TxStockLockedDecr(tx *xsql.Tx, skuID int64, count int64) (affected int64, err error) {
	res, err := tx.Exec(_decrStockLockedSQL, count, skuID, count)
	if err != nil {
		log.Error("d.TxStockLockedDecr(%d, %d) error(%v)", skuID, count, err)
		return
	}
	if affected, err = res.RowsAffected(); err != nil {
		log.Error("d.TxStockLockedDecr(%d, %d) res.RowsAffected() error(%v)", skuID, count, err)
	}
	return
}

// Stock 查询库存信息
func (d *Dao) Stock(c context.Context, skuID int64) (stock *model.SKUStock, err error) {
	stock = new(model.SKUStock)
	if err = d.db.QueryRow(c, _getStockBySkuIDSQL, skuID).Scan(&stock.SKUID, &stock.ParentSKUID, &stock.ItemID, &stock.Specs, &stock.TotalStock, &stock.Stock, &stock.LockedStock, &stock.SkAlert, &stock.Ctime, &stock.Mtime); err != nil {
		log.Error("d.Stock(%d), error(%v)", skuID, err)
	}
	return
}

// StockLogs 查询库存操作记录
func (d *Dao) StockLogs(c context.Context, opType int16, srcID int64, skuIDs ...int64) (stockLogs []*model.SKUStockLog, err error) {
	if len(skuIDs) == 0 {
		return
	}

	rows, err := d.db.Query(c, fmt.Sprintf(_selectSKUStockLog, xstr.JoinInts(skuIDs)), srcID, opType)
	if err != nil {
		log.Error("d.StockLogs() error(%v)", err)
		return
	}
	defer rows.Close()
	stockLogs = make([]*model.SKUStockLog, 0)
	for rows.Next() {
		stockLog := &model.SKUStockLog{}
		if err = rows.Scan(&stockLog.ID, &stockLog.SKUID, &stockLog.OpType, &stockLog.SrcID, &stockLog.Stock); err != nil {
			log.Error("d.StockLogs() rows.Scan() error(%v)", err)
			return
		}
		stockLogs = append(stockLogs)
	}
	return
}

// TxAddStockInsert 插入 stock 数据
func (d *Dao) TxAddStockInsert(tx *xsql.Tx, stocks ...*model.SKUStock) (affected int64, err error) {
	placeholder := strings.Trim(strings.Repeat("(?, ?, ?, ?, ?, ?, ?, ?),", len(stocks)), ",")
	var values []interface{}
	for _, stock := range stocks {
		values = append(values, stock.SKUID, stock.ParentSKUID, stock.ItemID, stock.Specs, stock.TotalStock, stock.Stock, stock.LockedStock, stock.SkAlert)
	}

	res, err := tx.Exec(fmt.Sprintf(_insertStockSQL, placeholder), values...)
	if err != nil {
		log.Error("d.TxStockInsert() error(%v)", err)
		return
	}

	if affected, err = res.RowsAffected(); err != nil {
		log.Error("d.TxStockInsert() res.RowsAffected() error(%v)", err)
		return
	}
	return
}

// TxStockReset 重置库存
func (d *Dao) TxStockReset(tx *xsql.Tx, stock *model.SKUStock) (affected int64, err error) {
	res, err := tx.Exec(_resetStockSQL, stock.TotalStock, stock.TotalStock, stock.SkAlert, stock.SKUID, stock.TotalStock)
	fmt.Println(_resetStockSQL, stock.TotalStock, stock.TotalStock, stock.SkAlert, stock.SKUID, stock.TotalStock)
	if err != nil {
		log.Error("d.TxStockReset() error(%v)", err)
		return
	}
	if affected, err = res.RowsAffected(); err != nil {
		log.Error("d.TxStockReset() res.RowsAffected() error(%v)", err)
		return
	}
	return
}

// TxStockLogRollBack 回滚操作日志
func (d *Dao) TxStockLogRollBack(tx *xsql.Tx, stockLogID int64) (affected int64, err error) {
	res, err := tx.Exec(_rollbackSKUStockLog, time.Now().Unix(), stockLogID)
	if err != nil {
		log.Error("d.TxStockLogRollBack(%d) error(%v)", stockLogID, err)
		return
	}
	if affected, err = res.RowsAffected(); err != nil {
		log.Error("d.TxStockLogRollBack(%d) res.RowsAffected() error(%v)", stockLogID, err)
	}
	return
}

// SkuItemCacheDel 删除 itemId => sku 缓存
func (d *Dao) SkuItemCacheDel(c context.Context, itemID int64) (err error) {
	if err = d.RedisDel(c, fmt.Sprintf(model.CacheKeyItemSku, itemID)); err != nil {
		log.Error("d.SkuItemCacheDel(%d) error(%v)", itemID, err)
	}
	return
}

// SkuByItemSpecs 通过 itemID specs 获取单个 sku
func (d *Dao) SkuByItemSpecs(c context.Context, itemID int64, specs string) (stock *model.SKUStock, err error) {
	res, err := d.SkuByItemID(c, itemID)
	if err != nil {
		log.Error("d.SkuByItemSpecs(%d, %s) d.SkuByItemID() error(%v)", itemID, specs, err)
		return
	}

	if item, ok := res[specs]; ok {
		stock = item
	}
	return
}

// RawSkuByItemSpecs 根据 itemID 和规格获取单个 sku
func (d *Dao) RawSkuByItemSpecs(c context.Context, itemID int64, specs string) (stock *model.SKUStock, err error) {
	stock = new(model.SKUStock)
	if err = d.db.QueryRow(c, _getStockBySpecsSQL, itemID, specs).Scan(&stock.SKUID, &stock.ParentSKUID, &stock.ItemID, &stock.Specs, &stock.TotalStock, &stock.Stock, &stock.LockedStock, &stock.SkAlert, &stock.Ctime, &stock.Mtime); err != nil {
		if err != sql.ErrNoRows {
			log.Error("d.SkuByItemSpecs(%d, %s) error(%v)", itemID, specs, err)
		}
	}
	return
}

// RawSkuByItemID 根据规格获取 sku
func (d *Dao) RawSkuByItemID(c context.Context, itemID int64) (stocks map[string]*model.SKUStock, err error) {
	rows, err := d.db.Query(c, _getStockByItemIDSQL, itemID)
	if err != nil {
		log.Error("d.RawSkuByItemID(%d) error(%v)", itemID, err)
		return
	}
	defer rows.Close()

	stocks = make(map[string]*model.SKUStock)
	for rows.Next() {
		stock := new(model.SKUStock)
		if err = rows.Scan(&stock.SKUID, &stock.ParentSKUID, &stock.ItemID, &stock.Specs, &stock.TotalStock, &stock.Stock, &stock.LockedStock, &stock.SkAlert, &stock.Ctime, &stock.Mtime); err != nil {
			log.Error("d.RawSkuByItemID(%d) rows.Scan() error(%v)", itemID, err)
			return
		}
		stocks[stock.Specs] = stock
	}
	return
}

// CacheSkuByItemID 根据 itemID 获取 sku 缓存
func (d *Dao) CacheSkuByItemID(c context.Context, itemID int64) (stocks map[string]*model.SKUStock, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", fmt.Sprintf(model.CacheKeyItemSku, itemID)))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("d.CacheSkuByItemID(%d) redis.Bytes() error(%v)", itemID, err)
		return
	}

	stocks = make(map[string]*model.SKUStock)
	if err = json.Unmarshal(reply, &stocks); err != nil {
		log.Error("d.CacheSkuByItemID(%d) json.Unmarshal() error(%v)", itemID, err)
		return
	}
	return
}

// AddCacheSkuByItemID 添加 itemId => sku 缓存
func (d *Dao) AddCacheSkuByItemID(c context.Context, itemID int64, stocks map[string]*model.SKUStock) (err error) {
	if stocks == nil {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()

	s, err := json.Marshal(stocks)
	if err != nil {
		log.Error("d.AddCacheSkuByItemID(%d, %+v) json.Marshal() error(%v)", itemID, stocks, err)
		return
	}
	if _, err = conn.Do("SETEX", fmt.Sprintf(model.CacheKeyItemSku, itemID), model.RedisExpireSku, s); err != nil {
		log.Error("d.AddCacheSkuByItemID(%d, %+v) conn.Do() error(%v)", itemID, stocks, err)
		return
	}
	return
}

// CacheStocks 获取 skuID => stock 库存缓存
func (d *Dao) CacheStocks(c context.Context, keys []int64, isLocked bool) (res map[int64]int64, err error) {
	if len(keys) == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()

	cacheKey := model.CacheKeyStock
	if isLocked {
		cacheKey = model.CacheKeyStockL
	}

	args := make([]interface{}, 0)
	for _, key := range keys {
		args = append(args, fmt.Sprintf(cacheKey, key))
	}
	int64s, err := redis.Int64s(conn.Do("MGET", args...))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("d.CacheStocks(%v, %t) error(%v)", keys, isLocked, err)
		return
	}

	res = make(map[int64]int64)
	for index, val := range int64s {
		if val < 0 {
			val = 0
		}
		res[keys[index]] = val
	}
	return
}

// RawStocks skuID => stock 缓存回源
func (d *Dao) RawStocks(c context.Context, keys []int64, isLocked bool) (res map[int64]int64, err error) {
	if len(keys) == 0 {
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_getStockBySkuID, xstr.JoinInts(keys)))
	if err != nil {
		log.Error("d.RawStocks() error(%v)", err)
		return
	}
	defer rows.Close()

	res = make(map[int64]int64)
	for rows.Next() {
		var skuID, stock, lockedStock int64
		if err = rows.Scan(&skuID, &stock, &lockedStock); err != nil {
			log.Error("d.RawStocks() rows.Scan() error(%v)", err)
			return
		}
		if isLocked {
			res[skuID] = lockedStock
		} else {
			res[skuID] = stock
		}
	}
	return
}

// AddCacheStocks skuID => stock 加入缓存
func (d *Dao) AddCacheStocks(c context.Context, stocks map[int64]int64, isLocked bool) (err error) {
	cacheKey := model.CacheKeyStock
	if isLocked {
		cacheKey = model.CacheKeyStockL
	}

	for skuID, stock := range stocks {
		if err1 := d.RedisSetnx(c, fmt.Sprintf(cacheKey, skuID), stock, model.RedisExpireStock); err1 != nil {
			log.Warn("d.AddCacheStocks() d.RedisSetnx(%s, %d, %d) error(%v)", fmt.Sprintf(cacheKey, skuID), stock, model.RedisExpireStock, err)
		}
	}
	return
}

// CacheGetSKUs 根据 skuID 获取 sku
// withNewStock 是否获取最新库存信息
func (d *Dao) CacheGetSKUs(c context.Context, skuIds []int64, withNewStock bool) (skuMap map[int64]*model.SKUStock, err error) {
	if len(skuIds) == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()

	args := make([]interface{}, 0)
	for _, skuID := range skuIds {
		args = append(args, fmt.Sprintf(model.CacheKeySku, skuID))
	}

	res, err := redis.ByteSlices(conn.Do("MGET", args...))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("d.CacheGetSKUs(%v, %t) error(%v)", skuIds, withNewStock, err)
		return
	}

	skuMap = make(map[int64]*model.SKUStock, len(res))
	for _, v := range res {
		if len(v) == 0 {
			continue
		}
		sku := &model.SKUStock{}
		if err = json.Unmarshal(v, sku); err != nil {
			log.Error("d.CacheGetSKUs() json.Unmarshal(%s) error(%v)", v, err)
			return
		}
		skuMap[sku.SKUID] = sku
	}

	if withNewStock {
		var stockMap map[int64]int64
		if stockMap, err = d.Stocks(c, skuIds, false); err != nil {
			log.Error("d.CacheGetSKUs() d.Stocks(%v) error(%v)", skuIds, err)
			return
		}
		for _, sku := range skuMap {
			sku.Stock = stockMap[sku.SKUID]
		}
	}
	return
}

// RawGetSKUs .
func (d *Dao) RawGetSKUs(c context.Context, skuIds []int64, withNewStock bool) (skuMap map[int64]*model.SKUStock, err error) {
	if len(skuIds) == 0 {
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_getStocksBySkuIDSQL, xstr.JoinInts(skuIds)))
	if err != nil {
		log.Error("d.RawGetSKUs() error(%v)", err)
	}
	defer rows.Close()

	skuMap = make(map[int64]*model.SKUStock)
	for rows.Next() {
		sku := &model.SKUStock{}
		if err = rows.Scan(&sku.SKUID, &sku.ParentSKUID, &sku.ItemID, &sku.Specs, &sku.TotalStock, &sku.Stock, &sku.LockedStock, &sku.SkAlert, &sku.Ctime, &sku.Mtime); err != nil {
			log.Error("d.RawGetSKUs() rows.Scan error(%v)", err)
			return
		}
		skuMap[sku.SKUID] = sku
	}
	return
}

// AddCacheGetSKUs .
func (d *Dao) AddCacheGetSKUs(c context.Context, skuMap map[int64]*model.SKUStock, withNewStock bool) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	for skuID, sku := range skuMap {
		var v []byte
		if v, err = json.Marshal(sku); err != nil {
			log.Warn("d.AddCacheGetSKUs() json.Marshal(%v) error(%v)", sku, err)
			err = nil
			continue
		}
		conn.Send("SETEX", fmt.Sprintf(model.CacheKeySku, skuID), model.RedisExpireSkuTmp, v)
	}
	return
}
