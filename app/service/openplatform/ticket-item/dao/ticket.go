package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// TkTypeSingle 单场票
	TkTypeSingle = 1
	// TkTypePass 通票
	TkTypePass = 2
	// TkTypeAllPass 联票
	TkTypeAllPass = 3
	// TimeNull 空时间：0000-00-00 00:00:00
	TimeNull = -62135596800
)

// RawTkListByItem 批量取项目票价
func (d *Dao) RawTkListByItem(c context.Context, ids []int64) (info map[int64][]*model.TicketInfo, err error) {
	info = make(map[int64][]*model.TicketInfo)
	tkExt := make(map[int64]map[string]*model.TicketPriceExtra)
	rows, err := d.db.Model(&model.TicketPrice{}).Where("project_id in (?) and deleted_at = 0", ids).Rows()
	extRows, err := d.db.Model(&model.TicketPriceExtra{}).Where("project_id in (?) and is_deleted = 0", ids).Rows()
	if err != nil {
		log.Error("RawListByItem(%v) error(%v)", model.JSONEncode(ids), err)
		return
	}
	defer rows.Close()
	defer extRows.Close()
	for extRows.Next() {
		var ext model.TicketPriceExtra
		err = d.db.ScanRows(extRows, &ext)
		if err != nil {
			log.Error("RawListByItem(%v) error(%v)", model.JSONEncode(ids), err)
			return
		}
		if _, ok := tkExt[ext.SkuID]; !ok {
			tkExt[ext.SkuID] = make(map[string]*model.TicketPriceExtra)
		}
		if _, ok := tkExt[ext.SkuID][ext.Attrib]; !ok {
			tkExt[ext.SkuID][ext.Attrib] = new(model.TicketPriceExtra)
		}
		tkExt[ext.SkuID][ext.Attrib] = &ext
	}
	for rows.Next() {
		var tk model.TicketInfo
		err = d.db.ScanRows(rows, &tk)
		if err != nil {
			log.Error("RawListByItem(%v) error(%v)", model.JSONEncode(ids), err)
			return
		}
		if _, ok := tkExt[tk.ID]; ok {
			tk.BuyNumLimit = tkExt[tk.ID]
		}
		info[tk.ProjectID] = append(info[tk.ProjectID], &tk)
	}
	return
}

// CacheTkListByItem 缓存取项目票价
func (d *Dao) CacheTkListByItem(c context.Context, ids []int64) (info map[int64][]*model.TicketInfo, err error) {
	var keys []interface{}
	keyPidMap := make(map[string]int64, len(ids))
	for _, id := range ids {
		key := keyItemTicket(id)
		if _, ok := keyPidMap[key]; !ok {
			// duplicate mid
			keyPidMap[key] = id
			keys = append(keys, key)
		}
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	var data [][]byte
	log.Info("MGET %v", model.JSONEncode(keys))
	if data, err = redis.ByteSlices(conn.Do("mget", keys...)); err != nil {
		log.Error("TkList MGET %v ERR: %v", model.JSONEncode(keys), err)
		return
	}
	info = make(map[int64][]*model.TicketInfo)
	for _, d := range data {
		if d != nil {
			var tks []*model.TicketInfo
			json.Unmarshal(d, &tks)
			info[tks[0].ProjectID] = tks
		}
	}
	return
}

// AddCacheTkListByItem 取项目票价添加缓存
func (d *Dao) AddCacheTkListByItem(c context.Context, info map[int64][]*model.TicketInfo) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	var data []interface{}
	var keys []string
	for k, v := range info {
		b, _ := json.Marshal(v)
		key := keyItemTicket(k)
		keys = append(keys, key)
		data = append(data, key, b)
	}
	log.Info("MSET %v", keys)
	if err = conn.Send("MSET", data...); err != nil {
		return
	}
	for i := 0; i < len(data); i += 2 {
		conn.Send("EXPIRE", data[i], CacheTimeout)
	}
	return

}

// RawTkList 批量取票价
func (d *Dao) RawTkList(c context.Context, ids []int64) (list map[int64]*model.TicketInfo, err error) {
	list = make(map[int64]*model.TicketInfo)
	tkExt := make(map[int64]map[string]*model.TicketPriceExtra)
	rows, err := d.db.Model(&model.TicketPrice{}).Where("id in (?) and deleted_at = 0", ids).Rows()
	extRows, err := d.db.Model(&model.TicketPriceExtra{}).Where("sku_id in (?) and is_deleted = 0", ids).Rows()
	if err != nil {
		log.Error("RawTkList(%v) error(%v)", model.JSONEncode(ids), err)
		return
	}
	defer rows.Close()
	defer extRows.Close()

	for extRows.Next() {
		var ext model.TicketPriceExtra
		err = d.db.ScanRows(extRows, &ext)
		if err != nil {
			log.Error("RawListByItem(%v) error(%v)", model.JSONEncode(ids), err)
			return
		}
		if _, ok := tkExt[ext.SkuID]; !ok {
			tkExt[ext.SkuID] = make(map[string]*model.TicketPriceExtra)
		}
		if _, ok := tkExt[ext.SkuID][ext.Attrib]; !ok {
			tkExt[ext.SkuID][ext.Attrib] = new(model.TicketPriceExtra)
		}
		tkExt[ext.SkuID][ext.Attrib] = &ext
	}

	for rows.Next() {
		var tk model.TicketInfo
		err = d.db.ScanRows(rows, &tk)
		if err != nil {
			log.Error("RawTkList(%v) error(%v)", model.JSONEncode(ids), err)
			return
		}
		if _, ok := tkExt[tk.ID]; ok {
			tk.BuyNumLimit = tkExt[tk.ID]
		}
		list[tk.ID] = &tk
	}
	return
}

// CacheTkList 缓存取项目票价
func (d *Dao) CacheTkList(c context.Context, ids []int64) (list map[int64]*model.TicketInfo, err error) {
	var keys []interface{}
	keyPidMap := make(map[string]int64, len(ids))
	for _, id := range ids {
		key := keyTicket(id)
		if _, ok := keyPidMap[key]; !ok {
			// duplicate mid
			keyPidMap[key] = id
			keys = append(keys, key)
		}
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	var data [][]byte
	log.Info("MGET %v", model.JSONEncode(keys))
	if data, err = redis.ByteSlices(conn.Do("mget", keys...)); err != nil {
		log.Error("TkList MGET %v ERR: %v", model.JSONEncode(keys), err)
		return
	}
	list = make(map[int64]*model.TicketInfo)
	for _, d := range data {
		if d != nil {
			var tk *model.TicketInfo
			json.Unmarshal(d, &tk)
			list[tk.ID] = tk
		}
	}
	return
}

// AddCacheTkList 取项目票价添加缓存
func (d *Dao) AddCacheTkList(c context.Context, list map[int64]*model.TicketInfo) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	var data []interface{}
	var keys []string
	for k, v := range list {
		b, _ := json.Marshal(v)
		key := keyTicket(k)
		keys = append(keys, key)
		data = append(data, key, b)
	}
	log.Info("MSET %v", keys)
	if err = conn.Send("MSET", data...); err != nil {
		return
	}
	for i := 0; i < len(data); i += 2 {
		conn.Send("EXPIRE", data[i], CacheTimeout)
	}
	return

}

// CreateOrUpdateTkPrice 创建或更新票种
func (d *Dao) CreateOrUpdateTkPrice(c context.Context, tx *gorm.DB, priceInfo model.TicketPrice, opType int32) (model.TicketPrice, error) {

	if opType == 0 {
		// create
		if err := tx.Create(&priceInfo).Error; err != nil {
			log.Error("创建票种失败:%s", err)
			tx.Rollback()
			return model.TicketPrice{}, err
		}
	} else {
		// update
		if err := tx.Model(&model.TicketPrice{}).Where("id = ?", priceInfo.ID).Updates(
			map[string]interface{}{
				"project_id":     priceInfo.ProjectID,
				"desc":           priceInfo.Desc,
				"type":           priceInfo.Type,
				"sale_type":      priceInfo.SaleType,
				"color":          priceInfo.Color,
				"buy_limit":      priceInfo.BuyLimit,
				"payment_method": priceInfo.PaymentMethod,
				"payment_value":  priceInfo.PaymentValue,
				"desc_detail":    priceInfo.DescDetail,
			}).Error; err != nil {
			log.Error("更新票种失败:%s", err)
			tx.Rollback()
			return model.TicketPrice{}, err
		}

	}
	return priceInfo, nil

}

// InsertOrUpdateTkPass 新建或更新通票联票票价
func (d *Dao) InsertOrUpdateTkPass(c context.Context, tx *gorm.DB, pid int64, scID int64, tksPass []TicketPass, tkType int32,
	scIDList map[int32]int64, TkSingleIDList map[int32]int64, TkSingleTypeList map[int32]int32) ([]TicketPass, error) {
	alphabetTable := model.AlphabetTable()

	for k, v := range tksPass {
		if _, ok := TkSingleIDList[v.LinkTicket]; !ok {
			tx.Rollback()
			log.Error("关联票种不存在")
			return nil, ecode.TicketLkTkNotFound
		}
		if _, ok := TkSingleTypeList[v.LinkTicket]; !ok {
			tx.Rollback()
			log.Error("关联票种类型不存在")
			return nil, ecode.TicketLkTkTypeNotFound
		}
		var linkScIDs []int64
		for _, linkScID := range v.LinkScreens {
			if _, ok := scIDList[linkScID]; !ok {
				tx.Rollback()
				log.Error("关联场次不存在")
				return nil, ecode.TicketLkScNotFound
			}
			linkScIDs = append(linkScIDs, scIDList[linkScID])
		}
		tkID, _ := strconv.ParseInt(v.TicketID, 10, 64)

		symbol := alphabetTable[k]
		if tkID == 0 {
			// create
			newTkID, err := model.GetTicketIDFromBase()
			if err != nil {
				tx.Rollback()
				log.Error("basecenter获取通票票价id失败:%s", err)
				return nil, err
			}
			buyLimit, _ := strconv.ParseInt(v.BuyLimit, 10, 64)
			payMethod, _ := strconv.ParseInt(v.PayMethod, 10, 64)
			if err = tx.Create(&model.TicketPrice{
				ID:            newTkID,
				ProjectID:     pid,
				ScreenID:      scID,
				Desc:          v.Name,
				BuyLimit:      int32(buyLimit),
				ParentID:      TkSingleIDList[v.LinkTicket],
				Color:         v.Color,
				DescDetail:    v.Desc,
				PaymentMethod: int32(payMethod),
				PaymentValue:  v.PayValue,
				Type:          tkType,
				SaleType:      TkSingleTypeList[v.LinkTicket],
				Symbol:        symbol,
				LinkSc:        model.Implode(",", linkScIDs),
				IsSale:        0,   // 不可售
				IsRefund:      -10, // 不可退
				OriginPrice:   -1,  // 未設置
				MarketPrice:   -1,
				SaleStart:     TimeNull, // 0000-00-00 00:00:00
				SaleEnd:       TimeNull,
			}).Error; err != nil {
				log.Error("通票或联票创建失败:%s", err)
				tx.Rollback()
				return nil, err
			}
			//票价限购
			limitData := d.FormatByPrefix(v.BuyLimitNum, "buy_limit_")
			if err = d.CreateOrUpdateTkPriceExtra(c, tx, limitData, newTkID, pid); err != nil {
				return nil, err
			}
			tksPass[k].TicketID = strconv.FormatInt(newTkID, 10)
		} else {
			// update
			if err := tx.Model(&model.TicketPrice{}).Where("id = ?", tkID).Updates(map[string]interface{}{
				"screen_id":      scID,
				"desc":           v.Name,
				"buy_limit":      v.BuyLimit,
				"parent_id":      TkSingleIDList[v.LinkTicket],
				"color":          v.Color,
				"desc_detail":    v.Desc,
				"payment_method": v.PayMethod,
				"payment_value":  v.PayValue,
				"type":           tkType,
				"sale_type":      TkSingleTypeList[v.LinkTicket],
				"symbol":         symbol,
				"link_sc":        model.Implode(",", linkScIDs),
			}).Error; err != nil {
				log.Error("通票或联票票价信息更新失败:%s", err)
				tx.Rollback()
				return nil, err
			}
			//票价限购
			limitData := d.FormatByPrefix(v.BuyLimitNum, "buy_limit_")
			if err := d.CreateOrUpdateTkPriceExtra(c, tx, limitData, tkID, pid); err != nil {
				return nil, err
			}

		}
	}
	return tksPass, nil
}

// DelTicket 根据id删除票种或票价
func (d *Dao) DelTicket(c context.Context, tx *gorm.DB, oldIDs []int64, newIDs []int64, pid int64, isPrice bool) error {
	delIDs, _ := model.ClassifyIDs(oldIDs, newIDs)
	for _, delID := range delIDs {
		if !d.CanDelTicket(delID, isPrice) {
			tx.Rollback()
			return ecode.TicketCannotDelTk
		}
		if isPrice {
			// TODO 存在需要删除的票价时 检查票价是否在坐票可选座场次下 是的话需要删除对应的座位图

			// 删除票价
			if err := tx.Exec("UPDATE ticket_price SET deleted_at=? WHERE id = ? AND project_id = ?", time.Now().Format("2006-01-02 15:04:05"), delID, pid).Error; err != nil {
				log.Error("删除票价失败:%s", err)
				tx.Rollback()
				return ecode.TicketDelTkFailed
			}
			// 删除票价额外信息表
			if err := tx.Model(&model.TicketPriceExtra{}).Where("sku_id = ? project_id = ?", delID, pid).Update("is_deleted", 1).Error; err != nil {
				log.Error("删除票种额外信息记录失败:%s", err)
				tx.Rollback()
				return ecode.TicketDelTkExFailed
			}

		} else {
			// 票种 需要获取 所有票价id
			priceIDs, err := d.GetPriceIDs(delID, 2)
			if err != nil {
				tx.Rollback()
				return err
			}
			// TODO 存在需要删除的票价时 检查票价是否在坐票可选座场次下 是的话需要删除对应的座位图

			// 将票种id加到需要删除的票价array里
			priceIDs = append(priceIDs, delID)
			// 删除票种票价
			if err := tx.Exec("UPDATE ticket_price SET deleted_at=? WHERE id IN (?) AND project_id = ?", time.Now().Format("2006-01-02 15:04:05"), priceIDs, pid).Error; err != nil {
				log.Error("删除票种及其票价失败:%s", err)
				tx.Rollback()
				return ecode.TicketDelTkFailed
			}
			//删除票价额外信息表
			if err := tx.Model(&model.TicketPriceExtra{}).Where("sku_id IN (?) AND project_id = ?", priceIDs, pid).Update("is_deleted", 1).Error; err != nil {
				log.Error("删除票种及票价额外信息记录失败:%s", err)
				tx.Rollback()
				return ecode.TicketDelTkExFailed
			}
		}

	}
	return nil
}

// CanDelTicket 检查是否可以删除票价或票种
func (d *Dao) CanDelTicket(id int64, isPrice bool) bool {
	var priceIDs []int64
	if isPrice {
		// 票价
		priceIDs = append(priceIDs, id)
	} else {
		// 票种 需要获取 所有票价id
		ids, err := d.GetPriceIDs(id, 2)
		if err != nil {
			return false
		}
		priceIDs = ids
	}

	if d.HasPromotion(priceIDs, 2) || d.StockChanged(priceIDs) {
		log.Error("票价下存在拼团或者库存有变动:%d", id)
		return false
	}
	return true
}

// GetPriceIDs 获取场次或票种下所有票价id inputType 1-场次id 2-票种id
func (d *Dao) GetPriceIDs(id int64, inputType int32) ([]int64, error) {
	var priceIDs []int64
	var prices []model.TicketPrice
	var whereStr string
	if inputType == 1 {
		// id = screenID
		whereStr = "screen_id = ? and deleted_at = 0"
	} else {
		// id = skuID
		whereStr = "parent_id = ? and deleted_at = 0"
	}

	if err := d.db.Select("id").Where(whereStr, id).Find(&prices).Error; err != nil {
		log.Error("获取场次或票种下所有票价id失败:%s", err)
		return nil, err
	}

	for _, v := range prices {
		priceIDs = append(priceIDs, v.ID)
	}

	return priceIDs, nil
}

// ticket_price表同时包含票价和票种，票种的parent_id为0，票价的parent_id为票种ID
// 票种不直接使用于场次，票价继承自票种，指定某一场次
// 以此实现“票种可以在一个项目下多个场次通用”
// 坐票只存在单场票

// TxGetTicketPrice 获取票价的价格标志和场次（事务）
func (d *Dao) TxGetTicketPrice(c context.Context, tx *gorm.DB, id int64) (ticketPrice *model.TicketPrice, err error) {
	ticketPrice = new(model.TicketPrice)
	if err = tx.Select("symbol, screen_id").Where("id = ? AND parent_id <> 0 AND deleted_at = 0", id).First(ticketPrice).Error; err != nil {
		log.Error("TxGetTicketPrice error(%v)", err)
	}
	return
}

// TxGetPriceSymbols 获取场次下的所有票价的父票种ID、价格和标志（事务）
func (d *Dao) TxGetPriceSymbols(c context.Context, tx *gorm.DB, screen int64) (ticketPrices []*model.TicketPrice, err error) {
	if err = tx.Select("parent_id, price, symbol").Where("screen_id = ? AND type = ? AND parent_id <> 0 AND deleted_at = 0", screen, TkTypeSingle).Find(&ticketPrices).Error; err != nil {
		log.Error("TxGetPriceSymbols error(%v)", err)
	}
	return
}

// TxGetParentTicketPrice 获取票种-单场票（事务）
func (d *Dao) TxGetParentTicketPrice(c context.Context, tx *gorm.DB, id int64) (ticketPrice *model.TicketPrice, err error) {
	ticketPrice = new(model.TicketPrice)
	if err = tx.Where("id = ? AND type = ? AND deleted_at = 0", id, TkTypeSingle).First(ticketPrice).Error; err != nil {
		log.Error("TxGetParentTicketPrice error(%v)", err)
	}
	return
}

// TxBatchAddTicketPrice 批量添加票价（事务）
func (d *Dao) TxBatchAddTicketPrice(c context.Context, tx *gorm.DB, ticketPrices []*model.TicketPrice) (err error) {
	if len(ticketPrices) == 0 {
		return
	}
	var values = make([]string, len(ticketPrices))
	for i, tp := range ticketPrices {
		values[i] = fmt.Sprintf("(%d,%d,%d,%d,%d,'%s','%s',%d,'%s',%d,%d,%d,%d,%d,%d,%d,'%s',%d,'%s',%d,'%s',%d,%d,%d)", tp.ID, tp.ProjectID, tp.ScreenID, tp.Price, tp.BuyLimit, tp.Desc, tp.Color, tp.ParentID, tp.Symbol, tp.IsSale, tp.OriginPrice, tp.PaymentMethod, tp.PaymentValue, tp.Type, tp.IsRefund, tp.IsVisible, tp.DescDetail, tp.SaleType, tp.SaleTime, tp.LinkTicketID, tp.LinkSc, tp.SaleStart, tp.SaleEnd, tp.MarketPrice)
	}
	var sql = fmt.Sprintf("INSERT INTO `ticket_price` (`id`, `project_id`, `screen_id`, `price`, `buy_limit`, `desc`, `color`, `parent_id`, `symbol`, `is_sale`, `origin_price`, `payment_method`, `payment_value`, `type`, `is_refund`, `is_visible`, `desc_detail`, `sale_type`, `sale_time`, `link_ticket_id`, `link_sc`, `sale_start`, `sale_end`, `market_price`) VALUES %s;", strings.Join(values, ","))
	if err = tx.Exec(sql).Error; err != nil {
		log.Error("批量添加票种（%s）失败:%s", sql, err)
		err = ecode.NotModified
		return
	}
	return
}

// CreateOrUpdateTkPriceExtra 创建票价额外信息记录
func (d *Dao) CreateOrUpdateTkPriceExtra(c context.Context, tx *gorm.DB, input map[string]string, skuID int64, pid int64) (err error) {

	data := d.FormatInputData(input, skuID, pid)

	var tmpTkExtra model.TicketPriceExtra
	for attrib, val := range data {
		if err = tx.Where("sku_id = ? and project_id = ? and attrib = ? and is_deleted=0", skuID, pid, attrib).First(&tmpTkExtra).Error; err != nil {
			//除去没查找到记录的报错 其他直接抛错
			if err != ecode.NothingFound {
				log.Error("获取票价%s额外信息失败:%s", skuID, err)
				tx.Rollback()
				return
			}
		}

		if tmpTkExtra.ID == 0 {
			// create
			if err = tx.Create(&val).Error; err != nil {
				log.Error("创建票价额外信息记录失败:%s", err)
				tx.Rollback()
				return
			}
		} else {
			// update
			if err = tx.Model(&model.TicketPriceExtra{}).Where("sku_id = ? and project_id = ? and attrib = ?", skuID, pid, attrib).Update("value", val).Error; err != nil {
				log.Error("更新票价额外信息记录失败:%s", err)
				tx.Rollback()
				return
			}
		}

	}
	return
}

// FormatInputData 格式化input数据
func (d *Dao) FormatInputData(input map[string]string, skuID int64, pid int64) (res []model.TicketPriceExtra) {
	for key, value := range input {
		res = append(res, model.TicketPriceExtra{
			Attrib:    key,
			Value:     value,
			SkuID:     skuID,
			ProjectID: pid,
		})
	}
	return
}

// FormatByPrefix 给键值加前缀
func (d *Dao) FormatByPrefix(input []string, prefix string) map[string]string {
	result := make(map[string]string)
	for k, v := range input {
		result[prefix+strconv.Itoa(k)] = v
	}
	return result
}
