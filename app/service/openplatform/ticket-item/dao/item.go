package dao

import (
	"context"
	"encoding/json"
	"go-common/library/cache/redis"

	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/log"
)

// RawItems 取项目信息
func (d *Dao) RawItems(c context.Context, ids []int64) (info map[int64]*model.Item, err error) {
	rows, err := d.db.Model(&model.Item{}).Where("id in (?)", ids).Rows()
	info = make(map[int64]*model.Item)
	if err != nil {
		log.Error("QueryItem(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item model.Item
		err = d.db.ScanRows(rows, &item)
		json.Unmarshal([]byte(item.PerformanceImage), &item.Img)
		info[item.ID] = &item
	}

	return
}

// CacheItems 缓存取项目信息
func (d *Dao) CacheItems(c context.Context, ids []int64) (info map[int64]*model.Item, err error) {
	var data [][]byte
	keys := make([]interface{}, len(ids))
	conn := d.redis.Get(c)
	defer conn.Close()
	keyPidMap := make(map[string]int64, len(ids))
	for k, id := range ids {
		key := keyItem(id)
		if _, ok := keyPidMap[key]; !ok {
			// duplicate id
			keyPidMap[key] = id
			keys[k] = key
		}
	}
	log.Info("MGET %v", model.JSONEncode(keys))
	if data, err = redis.ByteSlices(conn.Do("MGET", keys...)); err != nil {
		log.Error("MGET ERR: %v", err)
		return
	}
	info = make(map[int64]*model.Item)
	for _, d := range data {
		if d != nil {
			item := &model.Item{}
			json.Unmarshal(d, item)
			info[item.ID] = item
		}
	}
	return
}

// AddCacheItems 缓存取项目信息
func (d *Dao) AddCacheItems(c context.Context, info map[int64]*model.Item) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	var data []interface{}
	var keys []string
	for k, v := range info {
		b, _ := json.Marshal(v)
		key := keyItem(k)
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

// RawItemDetails 批量获取项目详情
func (d *Dao) RawItemDetails(c context.Context, ids []int64) (detail map[int64]*model.ItemDetail, err error) {
	rows, err := d.db.Model(&model.ItemDetail{}).Where("project_id in (?)", ids).Rows()
	detail = make(map[int64]*model.ItemDetail)
	if err != nil {
		log.Error("RawItemDetail(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item model.ItemDetail
		err = d.db.ScanRows(rows, &item)
		detail[item.ProjectID] = &item
	}
	return
}

// CacheItemDetails 缓存批量获取项目详情
func (d *Dao) CacheItemDetails(c context.Context, ids []int64) (detail map[int64]*model.ItemDetail, err error) {
	var data [][]byte
	var keys []interface{}
	keyPidMap := make(map[string]int64, len(ids))
	for _, id := range ids {
		key := keyItemDetail(id)
		if _, ok := keyPidMap[key]; !ok {
			// duplicate pid
			keyPidMap[key] = id
			keys = append(keys, key)
		}
	}
	conn := d.redis.Get(c)
	defer conn.Close()

	log.Info("MGET %v", model.JSONEncode(keys))
	if data, err = redis.ByteSlices((conn.Do("MGET", keys...))); err != nil {
		return
	}
	detail = make(map[int64]*model.ItemDetail)
	for _, d := range data {
		if d != nil {
			var v *model.ItemDetail
			json.Unmarshal(d, &v)
			detail[v.ProjectID] = v
		}
	}

	return

}

// AddCacheItemDetails 缓存取项目详情
func (d *Dao) AddCacheItemDetails(c context.Context, detail map[int64]*model.ItemDetail) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	var keys []string
	var data []interface{}
	for k, v := range detail {
		b, _ := json.Marshal(v)
		key := keyItemDetail(k)
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

// DelItemCache 删除项目相关缓存
func (d *Dao) DelItemCache(c context.Context, ids []int64) (res bool, err error) {
	var (
		keys []interface{}
	)

	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()

	for _, id := range ids {
		keys = append(keys, keyItem(id))
		keys = append(keys, keyItemDetail(id))
	}

	log.Info("DEL %v", keys)
	if err = conn.Send("DEL", keys...); err != nil {
		return
	}

	return
}
