package service

import (
	"context"
	"strings"

	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// StockLock 锁定库存 单个订单锁定多个 sku
func (s *Service) StockLock(ctx context.Context, opType int16, srcID int64, items []*model.SkuCnt, isLock bool) (err error) {
	if srcID == 0 || len(items) == 0 {
		err = ecode.ParamInvalid
		return
	}

	tx, err := s.dao.BeginTx(ctx)
	if err != nil {
		return
	}

	for _, skuCnt := range items {
		stockLog := &model.SKUStockLog{
			SKUID:  skuCnt.SkuID,
			OpType: opType,
			SrcID:  srcID,
			Stock:  skuCnt.Count,
		}
		if err = s.dao.TxAddStockLog(tx, stockLog); err != nil {
			// 如果是主键重复错误 可能是由于 slb 重试引起 所以跳过这条数据
			if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
				log.Error("s.StockLock() error(%v)", err)
				err = nil
				continue
			} else {
				tx.Rollback()
				return
			}
		}

		var affected int64
		if isLock {
			affected, err = s.dao.TxStockLock(tx, skuCnt.SkuID, skuCnt.Count) // 如果是锁定库存的话 就执行锁定
		} else {
			affected, err = s.dao.TxStockDecr(tx, skuCnt.SkuID, skuCnt.Count) // 不是锁定的话就直接减
		}
		if err != nil {
			tx.Rollback()
			return
		}
		// 影响了 0 条代表库存不足
		if affected == 0 {
			tx.Rollback()
			err = ecode.TicketStockLack
			log.Warn("s.StockLock() s.dao.TxStockLock(%d, %d) 库存不足 error(%v) ", skuCnt.SkuID, skuCnt.Count, err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("s.StockLock() tx.Commit() error(%v)", err)
		return
	}

	// 减少 redis 缓存计数
	for _, skuCnt := range items {
		s.dao.StockCacheDecr(ctx, skuCnt.SkuID, skuCnt.Count)
	}
	return
}

// StocksLock 锁定库存 多个订单锁定同一个 sku
func (s *Service) StocksLock(ctx context.Context, skuID int64, items []*model.OrderStockCnt, isLock bool) (err error) {
	skuLogs := make([]*model.SKUStockLog, 0)
	var total int64
	for _, item := range items {
		skuLogs = append(skuLogs, &model.SKUStockLog{
			SKUID:  skuID,
			OpType: consts.OpTypeOrder,
			SrcID:  item.OrderID,
			Stock:  item.Count,
		})

		total += item.Count
	}

	if len(skuLogs) == 0 {
		return
	}

	// 开启事务
	tx, err := s.dao.BeginTx(ctx)
	if err != nil {
		log.Error("s.LockStocks() d.BeginTx() error(%v)", err)
		return
	}

	var affected int64
	if isLock {
		affected, err = s.dao.TxStockLock(tx, skuID, total) // 如果是锁定库存的话 就执行锁定
	} else {
		affected, err = s.dao.TxStockDecr(tx, skuID, total) // 不是锁定的话就直接减
	}
	if err != nil {
		tx.Rollback()
		return
	}
	// 影响了 0 条代表库存不足
	if affected == 0 {
		tx.Rollback()
		err = ecode.TicketStockLack
		log.Warn("s.StocksLock() s.dao.TxStockLock(%d, %d) 库存不足 error(%v) ", skuID, total, err)
		return
	}

	if err = s.dao.TxAddStockLog(tx, skuLogs...); err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("s.LockStocks() tx.Commit() error(%v)", err)
		return
	}

	// 减库存缓存
	if err = s.dao.StockCacheDecr(ctx, skuID, total); err != nil {
		log.Error("s.LockStocks() error(%v)", err)
		return
	}
	return
}

// StocksUnlock 对之前做过的扣减库存操作进行回滚库存（取消预留和取消支付时调用）
// 这里的 opType 传锁定库存时候的 opType
func (s *Service) StocksUnlock(ctx context.Context, opType int16, srcID int64, skuIds []int64) (err error) {
	if srcID == 0 || len(skuIds) == 0 {
		err = ecode.ParamInvalid
		return
	}

	// 查询锁库存记录
	stockLogs, err := s.dao.StockLogs(ctx, opType, srcID, skuIds...)
	if err != nil {
		log.Error("s.StocksUnlock() error(%v)", err)
		return
	}

	tx, err := s.dao.BeginTx(ctx)
	if err != nil {
		return
	}
	for _, skuID := range skuIds {
		stockLog := s.getStockLog(stockLogs, skuID, opType, srcID)
		if stockLog == nil {
			log.Warn("s.StocksUnlock() 没有库存操作记录")
			err = ecode.TicketStockLogNotFound
			tx.Rollback()
			return
		}
		// 回滚锁定记录
		var affected int64
		affected, err = s.dao.TxStockLogRollBack(tx, stockLog.ID)
		if err != nil {
			tx.Rollback()
			return
		}
		if affected == 0 {
			continue
		}
		// 解锁库存
		affected, err = s.dao.TxStockUnlock(tx, skuID, stockLog.Stock)
		if err != nil || affected == 0 {
			log.Error("s.StocksUnlock() s.dao.TxStockUnlock() affected(%d) error(%v)", affected, err)
			tx.Rollback()
			return
		}
		if err = s.dao.StockCacheDel(ctx, skuID); err != nil {
			log.Error("s.StocksUnlock() s.dao.StockCacheDel(%d) error(%v)", skuID, err)
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("s.StocksUnlock() tx.Commit() error(%v)", err)
		return
	}
	return
}

// StockLockedDecr 扣除锁定的库存（支付减库存）对应 php 的 unlock
func (s *Service) StockLockedDecr(ctx context.Context, opType int16, srcID int64, items []*model.SkuCnt) (err error) {
	if srcID == 0 || len(items) == 0 {
		err = ecode.ParamInvalid
		return
	}

	tx, err := s.dao.BeginTx(ctx)
	if err != nil {
		log.Error("s.StockLockedDecr() s.dao.BeginTx() error(%v)", err)
		return
	}

	// 需要减缓存的列表
	decrCacheList := make([]*model.SkuCnt, 0)
	for _, item := range items {
		// 写缓存操作日志
		if err = s.dao.TxAddStockLog(tx, &model.SKUStockLog{
			SKUID:  item.SkuID,
			OpType: opType,
			SrcID:  srcID,
			Stock:  item.Count,
		}); err != nil {
			// 重复写入跳过，可能因为 slb 重试导致
			if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
				log.Error("s.StockLockedDecr() error(%v)", err)
				err = nil
				continue
			} else {
				tx.Rollback()
				return
			}
		}

		// 扣除锁定库存
		var affected int64
		if affected, err = s.dao.TxStockLockedDecr(tx, item.SkuID, item.Count); err != nil || affected == 0 {
			tx.Rollback()
			return
		}

		decrCacheList = append(decrCacheList, item)
	}

	if err = tx.Commit(); err != nil {
		log.Error("s.StockLockedDecr() tx.Commit() error(%v)", err)
		return
	}

	// 减缓存
	for _, item := range decrCacheList {
		s.dao.StockLockedCacheDecr(ctx, item.SkuID, item.Count)
	}
	return
}

// StockIncr 增加库存(发生退票时调用) 对应 php incrStock
func (s *Service) StockIncr(ctx context.Context, opType int16, srcID int64, items []*model.SkuCnt) (err error) {
	if srcID == 0 || len(items) == 0 {
		err = ecode.ParamInvalid
		return
	}

	tx, err := s.dao.BeginTx(ctx)
	if err != nil {
		log.Error("s.StockIncr() s.dao.BeginTx() error(%v)", err)
		return
	}

	for _, item := range items {
		// 写缓存操作日志
		if err = s.dao.TxAddStockLog(tx, &model.SKUStockLog{
			SKUID:  item.SkuID,
			OpType: opType,
			SrcID:  srcID,
			Stock:  item.Count,
		}); err != nil {
			// 重复写入跳过，可能因为 slb 重试导致
			if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
				log.Error("s.StockIncr() error(%v)", err)
				err = nil
				continue
			} else {
				tx.Rollback()
				return
			}
		}

		var affected int64
		if affected, err = s.dao.TxStockIncr(tx, item.SkuID, item.Count); err != nil {
			tx.Rollback()
			return
		}
		if affected == 0 {
			// 如果是拼团分出来的子库存，拼团过期之后，拼团 sku 的 total_stock 就会被置为 0
			// 这时候子 sku 总库存为 0，加父 sku 库存，父 sku 库存不会为 0
			var stock *model.SKUStock
			if stock, err = s.dao.Stock(ctx, item.SkuID); err != nil {
				log.Error("s.StockIncr() s.dao.Stock(%d) error(%v)", item.SkuID, err)
				tx.Rollback()
				return
			}
			if stock.TotalStock == 0 {
				if _, err = s.dao.TxStockIncr(tx, stock.ParentSKUID, item.Count); err != nil {
					log.Error("s.StockIncr() s.dao.TxStockIncr(%d, %d) error(%v)", stock.ParentSKUID, item.Count, err)
					tx.Rollback()
					return
				}
			}
		}

		// redis 删除缓存
		if err = s.dao.StockCacheDel(ctx, item.SkuID); err != nil {
			log.Error("s.StockIncr() s.dao.StockCacheDel(%d), error(%v)", item.SkuID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error("s.StockLockedDecr() tx.Commit() error(%v)", err)
		return
	}
	return
}

// StockInit 设置库存，初始化场次库存时调用
func (s *Service) StockInit(ctx context.Context, batchID int64, itemID int64, batches []*model.Batch) (skuIds []int64, err error) {
	if batchID == 0 || len(batches) == 0 {
		err = ecode.ParamInvalid
		return
	}

	skuIds = make([]int64, 0)

	tx, err := s.dao.BeginTx(ctx)
	if err != nil {
		log.Error("s.StockInit() s.dao.BeginTx() error(%v)", err)
		return
	}

	for _, batch := range batches {
		var stock *model.SKUStock
		stockExist := false
		if stock, err = s.dao.RawSkuByItemSpecs(ctx, itemID, batch.Specs()); err != nil {
			if err != sql.ErrNoRows {
				log.Error("s.StockInit() s.dao.SkuByItemSpecs error()")
				tx.Rollback()
				return
			}
			stock = &model.SKUStock{
				SKUID:      batch.TicketPriceID,
				ItemID:     itemID,
				Specs:      batch.Specs(),
				TotalStock: batch.TotalStock,
				Stock:      batch.TotalStock,
				SkAlert:    batch.SkAlert,
			}
			if _, err = s.dao.TxAddStockInsert(tx, stock); err != nil {
				log.Error("s.StockInit() s.dao.TxAddStockInsert(%+v) error(%v)", stock, err)
				tx.Rollback()
				return
			}
		} else {
			stockExist = true
		}

		skuIds = append(skuIds, stock.SKUID)

		stockLog := &model.SKUStockLog{
			SKUID:  stock.SKUID,
			OpType: consts.OpTypeActive,
			SrcID:  batchID,
			Stock:  batch.TotalStock,
		}
		if err = s.dao.TxAddStockLog(tx, stockLog); err != nil {
			if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
				continue
			}
			log.Error("s.StockInit() s.dao.TxAddStockLog(%+v) error(%v)", stockLog, err)
			tx.Rollback()
			return
		}

		if stockExist {
			affected, err1 := s.dao.TxStockReset(tx, stock)
			if err1 != nil {
				err = err1
				log.Error("s.StockInit() s.dao.TxStockReset() error(%v)", err)
				tx.Rollback()
				return
			}
			if affected == 0 {
				err = ecode.TicketStockUpdateFail
				log.Error("s.StockInit() s.dao.TxStockReset() affected(%d) error(%v)", affected, err)
				tx.Rollback()
				return
			}
		}

		// 删除 skuId => stock 缓存
		s.dao.StockCacheDel(ctx, stock.SKUID)
		// 删除 skuId => sku 缓存
		s.dao.DelCacheSku(ctx, stock.SKUID)
	}

	if err = tx.Commit(); err != nil {
		log.Error("s.StockInit() tx.Commit() error(%v)", err)
		return
	}

	// 删除 sku 详情缓存
	s.dao.SkuItemCacheDel(ctx, itemID)
	return
}

func (s *Service) getStockLog(stockLogs []*model.SKUStockLog, skuID int64, opType int16, srcID int64) (stockLog *model.SKUStockLog) {
	for _, tmp := range stockLogs {
		if tmp.SKUID == skuID && tmp.OpType == opType && tmp.SrcID == srcID {
			stockLog = tmp
			return
		}
	}
	return
}
