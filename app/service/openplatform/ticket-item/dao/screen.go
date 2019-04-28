package dao

import (
	"context"
	"encoding/json"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// TypeWithoutSeat 场次类型 站票
	TypeWithoutSeat = 2
	// TicketTypeElec 场次票类型 电子票
	TicketTypeElec = 2
	// scTypeNormal 普通场次
	scTypeNormal = 1
	// scTypePass 通票场次
	scTypePass = 2
	// scTypeAllPass 联票场次
	scTypeAllPass = 3
)

var scTypeName = map[int32]string{
	scTypeNormal:  "普通场次",
	scTypePass:    "通票",
	scTypeAllPass: "联票",
}

// RawScListByItem 批量取项目场次
func (d *Dao) RawScListByItem(c context.Context, ids []int64) (info map[int64][]*model.Screen, err error) {
	info = make(map[int64][]*model.Screen)
	rows, err := d.db.Model(&model.Screen{}).Where("project_id in (?)", ids).Rows()
	if err != nil {
		log.Error("RawScListByItem(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sc model.Screen
		err = d.db.ScanRows(rows, &sc)
		if err != nil {
			log.Error("RawScListByItem(%v) scan error(%v)", ids, err)
			return
		}
		info[sc.ProjectID] = append(info[sc.ProjectID], &sc)
	}
	return
}

// CacheScListByItem 缓存取项目票价
func (d *Dao) CacheScListByItem(c context.Context, ids []int64) (info map[int64][]*model.Screen, err error) {
	var keys []interface{}
	keyPidMap := make(map[string]int64, len(ids))
	for _, id := range ids {
		key := keyItemScreen(id)
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
		log.Error("ItemScList MGET %v ERR: %v", model.JSONEncode(keys), err)
		return
	}
	info = make(map[int64][]*model.Screen)
	for _, d := range data {
		if d != nil {
			var scs []*model.Screen
			json.Unmarshal(d, &scs)
			info[scs[0].ProjectID] = scs
		}
	}
	return
}

// AddCacheScListByItem 为项目场次添加缓存
func (d *Dao) AddCacheScListByItem(c context.Context, info map[int64][]*model.Screen) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	var data []interface{}
	var keys []string
	for k, v := range info {
		b, _ := json.Marshal(v)
		key := keyItemScreen(k)
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

// RawScList 批量取场次
func (d *Dao) RawScList(c context.Context, ids []int64) (info map[int64]*model.Screen, err error) {
	info = make(map[int64]*model.Screen)
	rows, err := d.db.Model(&model.Screen{}).Where("id in (?)", ids).Rows()
	if err != nil {
		log.Error("RawScList(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sc model.Screen
		err = d.db.ScanRows(rows, &sc)
		if err != nil {
			log.Error("RawScList(%v) scan error(%v)", ids, err)
			return
		}
		info[sc.ID] = &sc
	}
	return
}

// CacheScList 缓存批量取场次
func (d *Dao) CacheScList(c context.Context, ids []int64) (info map[int64]*model.Screen, err error) {
	var keys []interface{}
	keyPidMap := make(map[string]int64, len(ids))
	for _, id := range ids {
		key := keyScreen(id)
		if _, ok := keyPidMap[key]; !ok {
			// duplicate id
			keyPidMap[key] = id
			keys = append(keys, key)
		}
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	var data [][]byte
	log.Info("MGET %v", model.JSONEncode(keys))
	if data, err = redis.ByteSlices(conn.Do("mget", keys...)); err != nil {
		log.Error("ScList MGET %v ERR: %v", model.JSONEncode(keys), err)
		return
	}
	info = make(map[int64]*model.Screen)
	for _, d := range data {
		if d != nil {
			var sc model.Screen
			json.Unmarshal(d, &sc)
			info[sc.ID] = &sc
		}
	}
	return
}

// AddCacheScList 为场次添加缓存
func (d *Dao) AddCacheScList(c context.Context, info map[int64]*model.Screen) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	var data []interface{}
	var keys []string
	for k, v := range info {
		b, _ := json.Marshal(v)
		key := keyScreen(k)
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

// CreateOrUpdateScreen 创建或更新场次
func (d *Dao) CreateOrUpdateScreen(c context.Context, tx *gorm.DB, screenInfo model.Screen) (model.Screen, error) {

	if screenInfo.ID == 0 {
		// create
		if err := tx.Create(&screenInfo).Error; err != nil {
			log.Error("创建场次失败:%s", err)
			tx.Rollback()
			return model.Screen{}, err
		}
	} else {
		// update
		if err := tx.Model(&model.Screen{}).Where("id = ?", screenInfo.ID).Updates(map[string]interface{}{
			"name":          screenInfo.Name,
			"status":        screenInfo.Status,
			"type":          screenInfo.Type,
			"ticket_type":   screenInfo.TicketType,
			"screen_type":   screenInfo.ScreenType,
			"delivery_type": screenInfo.DeliveryType,
			"pick_seat":     screenInfo.PickSeat,
			"start_time":    screenInfo.StartTime,
			"end_time":      screenInfo.EndTime,
			"project_id":    screenInfo.ProjectID,
		}).Error; err != nil {
			log.Error("更新场次失败:%s", err)
			tx.Rollback()
			return model.Screen{}, err
		}
	}

	return screenInfo, nil

}

// GetOrUpdatePassSc 更新或新建通联票场次 返回id
func (d *Dao) GetOrUpdatePassSc(c context.Context, tx *gorm.DB, pid int64, tksPass []TicketPass, scStartTimes map[int32]int32,
	scEndTimes map[int32]int32, scType int32, opType int32) (int64, error) {
	var passScID int64

	scTime, err := d.CalPassScTime(scStartTimes, scEndTimes, tksPass)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var screen model.Screen
	if d.db.Where("project_id = ? and screen_type = ? and deleted_at = 0", pid, scType).First(&screen).RecordNotFound() {
		// 不存在场次并且有票种存在则新建一个
		if len(tksPass) > 0 {
			passScreen := model.Screen{
				ProjectID:    pid,
				Name:         scTypeName[scType],
				StartTime:    scTime[0],
				EndTime:      scTime[1],
				Type:         TypeWithoutSeat, // 站票
				TicketType:   TicketTypeElec,  // 电子票
				DeliveryType: 1,               // 默认不配送
				ScreenType:   scType,
				Status:       opType,
			}
			if err := tx.Create(&passScreen).Error; err != nil {
				log.Error("创建通票或联票场次失败:%s", err)
				tx.Rollback()
				return 0, err
			}
			passScID = passScreen.ID
		}

	} else {
		if len(tksPass) > 0 {
			// update
			if err := tx.Model(&model.Screen{}).Where("id = ?", screen.ID).Updates(map[string]interface{}{
				"status":     opType,
				"start_time": scTime[0],
				"end_time":   scTime[1],
			}).Error; err != nil {
				log.Error("更新通票联票场次失败:%s", err)
				tx.Rollback()
				return 0, err
			}
			passScID = screen.ID
		} else {
			// 如果没有通票/联票 删除通票/联票场次
			if err := d.DelScreen(c, tx, []int64{screen.ID}, nil, pid); err != nil {
				return 0, err
			}
			passScID = screen.ID
		}
	}

	return passScID, nil
}

// DelScreen 删除场次
func (d *Dao) DelScreen(c context.Context, tx *gorm.DB, oldIDs []int64, newIDs []int64, pid int64) error {
	delIDs, _ := model.ClassifyIDs(oldIDs, newIDs)
	for _, delID := range delIDs {
		if !d.CanDelScreen(delID) {
			tx.Rollback()
			return ecode.TicketCannotDelSc
		}
		// 删除场次前把改场次下的票价一起删除
		priceIDs, err := d.GetPriceIDs(delID, 1)
		if err != nil {
			return err
		}
		if err := d.DelTicket(c, tx, priceIDs, nil, pid, true); err != nil {
			return err
		}
		// 删除场次
		if err := tx.Exec("UPDATE screen SET deleted_at=? WHERE id = ? AND project_id = ?", time.Now().Format("2006-01-02 15:04:05"), delID, pid).Error; err != nil {
			log.Error("删除场次失败:%s", err)
			tx.Rollback()
			return err
		}
	}
	return nil
}

// CalPassScTime 获取通联场次的开始结束时间, 取关联场次时间的极值。返回[0]=starttime [1]=endtime
func (d *Dao) CalPassScTime(scStartTimes map[int32]int32, scEndTimes map[int32]int32, tksPass []TicketPass) ([]int32, error) {
	var startTimes, endTimes, currStartTimes, currEndTimes []int32
	for _, v := range tksPass {
		currStartTimes = []int32{}
		currEndTimes = []int32{}
		for _, v2 := range v.LinkScreens {
			if _, ok := scStartTimes[v2]; !ok {
				log.Error("关联的场次开始时间不存在")
				return []int32{}, ecode.TicketLkScTimeNotFound
			}
			if _, ok := scEndTimes[v2]; !ok {
				log.Error("关联的场次结束时间不存在")
				return []int32{}, ecode.TicketLkScTimeNotFound
			}
			currStartTimes = append(currStartTimes, scStartTimes[v2])
			currEndTimes = append(currEndTimes, scEndTimes[v2])
		}
		startTimes = append(startTimes, model.Min(currStartTimes))
		endTimes = append(endTimes, model.Max(currEndTimes))
	}
	return []int32{model.Min(startTimes), model.Max(endTimes)}, nil

}

// CanDelScreen 检查是否可以删除场次
func (d *Dao) CanDelScreen(id int64) bool {
	var priceIDs []int64

	// 获取场次下 所有票价id
	var prices []model.TicketPrice
	if err := d.db.Select("id").Where("screen_id = ? and deleted_at = 0", id).Find(&prices).Error; err != nil {
		log.Error("获取场次下所有票价id失败:%s", err)
		return false
	}
	for _, v := range prices {
		priceIDs = append(priceIDs, v.ID)
	}

	if d.HasPromotion(priceIDs, 2) || d.StockChanged(priceIDs) {
		log.Error("场次的票价下存在拼团或者库存有变动")
		return false
	}
	return true
}
