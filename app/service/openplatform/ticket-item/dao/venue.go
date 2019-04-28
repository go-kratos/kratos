/*
	Package dao venue
	场馆（venue）=>场地（place）=>区域（area）的三级层次，均为一对多
	venue表冗余place_num表示场地数，
	place表通过ID对应place_polygon，存有area的地理位置信息
*/

package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// RawVenues 批量获取场馆
func (d *Dao) RawVenues(c context.Context, ids []int64) (vl map[int64]*model.Venue, err error) {
	vrows, err := d.db.Model(&model.Venue{}).Where("id in (?)", ids).Rows()
	vl = make(map[int64]*model.Venue)
	if err != nil {
		log.Error("RawVenues(%v) error(%v)", ids, err)
		return
	}
	defer vrows.Close()
	for vrows.Next() {
		var v model.Venue
		err = d.db.ScanRows(vrows, &v)
		vl[v.ID] = &v
	}
	return
}

// CacheVenues 缓存批量获取场馆
func (d *Dao) CacheVenues(c context.Context, ids []int64) (vl map[int64]*model.Venue, err error) {
	var keys []interface{}
	keyVidMap := make(map[string]int64, len(ids))
	for _, id := range ids {
		key := keyVenue(id)
		if _, ok := keyVidMap[key]; !ok {
			// duplicate pid
			keyVidMap[key] = id
			keys = append(keys, key)
		}
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	var data [][]byte
	log.Info("MGET %v", model.JSONEncode(keys))
	if data, err = redis.ByteSlices((conn.Do("MGET", keys))); err != nil {
		log.Error("VenueList MGET %v ERR: %v", model.JSONEncode(keys), err)
		return
	}
	vl = make(map[int64]*model.Venue)
	for _, d := range data {
		if d != nil {
			var v *model.Venue
			vl[v.ID] = v
			json.Unmarshal(d, &v)
		}
	}
	return
}

// AddCacheVenues 缓存场馆信息
func (d *Dao) AddCacheVenues(c context.Context, vl map[int64]*model.Venue) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	var data []interface{}
	var keys []string
	for k, v := range vl {
		b, _ := json.Marshal(v)
		key := keyVenue(k)
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

// VenueSearch 场馆搜索
func (d *Dao) VenueSearch(c context.Context, req *model.VenueSearchParam) (venues *model.VenueSearchList, err error) {
	r := d.es.NewRequest("ticket_venue").Index("ticket_venue")
	if req.ID > 0 {
		r.WhereEq("id", req.ID)
	} else if req.Name != "" {
		r.WhereLike([]string{"name"}, []string{req.Name}, false, elastic.LikeLevelLow)
	}
	if req.ProvinceID > 0 {
		r.WhereEq("province", req.ProvinceID)
	}
	if req.CityID > 0 {
		r.WhereEq("city", req.CityID)
	}
	r.Order("ctime", elastic.OrderDesc).Ps(req.Ps).Pn(req.Pn)
	log.Info(fmt.Sprintf("%s/x/admin/search/query?%s", d.c.URL.ElasticHost, r.Params()))

	venues = new(model.VenueSearchList)
	if err = r.Scan(c, venues); err != nil {
		log.Error("VenueSearch(%v) r.Query(%s) error(%s)", req, r.Params(), err)
	}
	return
}

// AddVenue 添加场馆信息
func (d *Dao) AddVenue(c context.Context, venue *model.Venue) (err error) {
	if err = d.db.Create(venue).Error; err != nil {
		log.Error("添加场馆信息失败:%s", err)
		err = ecode.NotModified
		return
	}
	return
}

// UpdateVenue 编辑场馆信息
func (d *Dao) UpdateVenue(c context.Context, venue *model.Venue) (err error) {
	// update venue with new info (using map can update the column with empty string)
	if err = d.db.Table("venue").Where("id = ?", venue.ID).Updates(
		map[string]interface{}{
			"name":           venue.Name,
			"city":           venue.City,
			"province":       venue.Province,
			"district":       venue.District,
			"address_detail": venue.AddressDetail,
			"status":         venue.Status,
			"traffic":        venue.Traffic,
		}).Error; err != nil {
		log.Error("更新场馆信息失败:%s", err)
		err = ecode.NotModified
		return
	}
	return
}

// TxIncPlaceNum 增加场馆的场地数（事务）
func (d *Dao) TxIncPlaceNum(c context.Context, tx *gorm.DB, venueID int64) (err error) {
	var oriVenue model.Venue
	if err = tx.First(&oriVenue, venueID).Error; err != nil {
		log.Error("查找对应的场馆信息（ID:%d）失败:%s", venueID, err)
		err = ecode.NotModified
		return
	}
	if err = tx.Model(&oriVenue).Updates(
		map[string]interface{}{
			"place_num": oriVenue.PlaceNum + 1,
		}).Error; err != nil {
		log.Error("更新场馆信息（ID:%d）失败:%s", venueID, err)
		err = ecode.NotModified
		return
	}
	return
}

// TxDecPlaceNum 减少场馆的场地数（事务）
func (d *Dao) TxDecPlaceNum(c context.Context, tx *gorm.DB, venueID int64) (err error) {
	var oriVenue model.Venue
	if err = tx.First(&oriVenue, venueID).Error; err != nil {
		log.Error("查找对应的场馆信息（ID:%d）失败:%s", venueID, err)
		err = ecode.NotModified
		return
	}
	if oriVenue.PlaceNum < 1 {
		log.Error("更新场馆信息（ID:%d）失败:场地数小于1", venueID)
		err = ecode.NotModified
		return
	}
	if err = tx.Model(&oriVenue).Updates(
		map[string]interface{}{
			"place_num": oriVenue.PlaceNum - 1,
		}).Error; err != nil {
		log.Error("更新场馆信息（ID:%d）失败:%s", venueID, err)
		err = ecode.NotModified
		return
	}
	return
}

// RawPlace 获取场地
func (d *Dao) RawPlace(c context.Context, id int64) (place *model.Place, err error) {
	place = new(model.Place)
	err = d.db.Model(&model.Place{}).First(&place, id).Scan(&place).Error

	if err != nil {
		log.Error("RawPlace(%v) error(%v)", id, err)
	}
	return
}

// RawPlacePolygon 获取场地坐标
func (d *Dao) RawPlacePolygon(c context.Context, id int64) (placePolygon *model.PlacePolygon, err error) {
	placePolygon = new(model.PlacePolygon)
	err = d.db.Model(&model.PlacePolygon{}).First(&placePolygon, id).Scan(&placePolygon).Error

	if err != nil {
		log.Error("RawPlacePolygon(%v) error(%v)", id, err)
	}
	return
}

// TxRawPlace 获取场地（事务）
func (d *Dao) TxRawPlace(c context.Context, tx *gorm.DB, id int64) (place *model.Place, err error) {
	place = new(model.Place)
	err = tx.Model(&model.Place{}).First(&place, id).Scan(&place).Error

	if err != nil {
		log.Error("TxRawPlace(%v) error(%v)", id, err)
	}
	return
}

// CachePlace 缓存获取场地
func (d *Dao) CachePlace(c context.Context, id int64) (place *model.Place, err error) {
	var data []byte
	key := keyPlace(id)
	conn := d.redis.Get(c)
	defer conn.Close()
	log.Info("GET %v", key)
	if data, err = redis.Bytes((conn.Do("GET", key))); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	json.Unmarshal(data, &place)
	return
}

// AddCachePlace 缓存场地信息
func (d *Dao) AddCachePlace(c context.Context, id int64, place *model.Place) (err error) {
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	var data []interface{}
	key := keyPlace(id)
	data = append(data, key, place)
	log.Info("SET %v", key)
	if err = conn.Send("SET", data...); err != nil {
		return
	}
	conn.Send("EXPIRE", data[0], CacheTimeout)
	return
}

// TxAddPlace 添加场地信息（事务）
func (d *Dao) TxAddPlace(c context.Context, tx *gorm.DB, place *model.Place) (err error) {
	if err = tx.Create(place).Error; err != nil {
		log.Error("添加场地信息失败:%s", err)
		err = ecode.NotModified
		return
	}
	return
}

// TxUpdatePlace 编辑场地信息（事务）
func (d *Dao) TxUpdatePlace(c context.Context, tx *gorm.DB, place *model.Place) (err error) {
	// find original place with id
	if err = tx.Table("place").Where("id = ?", place.ID).Update(
		map[string]interface{}{
			"name":     place.Name,
			"base_pic": place.BasePic,
			"status":   place.Status,
			"venue":    place.Venue,
			"d_width":  place.DWidth,
			"d_height": place.DHeight,
		}).Error; err != nil {
		log.Error("更新场地信息（ID:%d）失败:%s", place.ID, err)
		err = ecode.NotModified
		return
	}
	return
}

// TxAddOrUpdateAreaPolygon 添加或更新区域的场地坐标信息（事务）
func (d *Dao) TxAddOrUpdateAreaPolygon(c context.Context, tx *gorm.DB, place int64, area int64, coordinate *string) (err error) {
	var (
		oriPlacePolygon model.PlacePolygon
		coordinates     map[int64]string
		found           bool
	)
	if res := tx.First(&oriPlacePolygon, place); res.Error != nil {
		if !res.RecordNotFound() {
			log.Error("查找对应的场地坐标信息（ID:%d）失败:%s", place, res.Error)
			err = ecode.NotModified
			return
		}
		coordinates = make(map[int64]string)
		oriPlacePolygon.ID = place
		found = false
	} else {
		// 反序列化出coordinate字段
		if err = json.Unmarshal([]byte(oriPlacePolygon.Coordinate), &coordinates); err != nil {
			log.Error("数据库中的场地坐标信息（ID:%d）反序列化失败:%s", place, err)
			err = ecode.NotModified
			return
		}
		found = true
	}
	// 添加新的coordinate并序列化
	coordinates[area] = *coordinate
	b, _ := json.Marshal(coordinates)
	oriPlacePolygon.Coordinate = string(b)
	if found {
		if err = tx.Model(&oriPlacePolygon).Updates(
			map[string]interface{}{
				"coordinate": oriPlacePolygon.Coordinate,
			}).Error; err != nil {
			log.Error("更新场地坐标信息（ID:%d）失败:%s", place, err)
			err = ecode.NotModified
			return
		}
	} else {
		if err = tx.Create(oriPlacePolygon).Error; err != nil {
			log.Error("创建场地坐标信息（ID:%d）失败:%s", place, err)
			err = ecode.NotModified
			return
		}
	}
	*coordinate = oriPlacePolygon.Coordinate
	return
}

// TxDelAreaPolygon 删除区域的场地坐标信息（事务）
func (d *Dao) TxDelAreaPolygon(c context.Context, tx *gorm.DB, place int64, area int64) (err error) {
	var (
		oriPlacePolygon model.PlacePolygon
		coordinates     map[int64]string
	)
	if err = tx.First(&oriPlacePolygon, place).Error; err != nil {
		log.Error("查找对应的场地坐标信息（ID:%d）失败:%s", place, err)
		err = ecode.NotModified
		return
	}
	if oriPlacePolygon.Coordinate != "" {
		if err = json.Unmarshal([]byte(oriPlacePolygon.Coordinate), &coordinates); err != nil {
			log.Error("数据库中的场地坐标信息（ID:%d）反序列化失败:%s", place, err)
			err = ecode.NotModified
			return
		}
	} else {
		coordinates = make(map[int64]string)
	}

	// 添加新的coordinate并序列化
	if _, ok := coordinates[area]; !ok {
		log.Error("数据库中的场地坐标信息（ID:%d）并未包含该区域（ID:%d），删除失败", place, area)
		err = ecode.NotModified
		return
	}
	delete(coordinates, area)
	b, _ := json.Marshal(coordinates)
	oriPlacePolygon.Coordinate = string(b)
	if err = tx.Model(&oriPlacePolygon).Updates(
		map[string]interface{}{
			"coordinate": oriPlacePolygon.Coordinate,
		}).Error; err != nil {
		log.Error("创建场地坐标信息（ID:%d）失败:%s", place, err)
		err = ecode.NotModified
		return
	}
	return
}

// RawArea 获取区域
func (d *Dao) RawArea(c context.Context, id int64) (area *model.Area, err error) {
	area = new(model.Area)
	if err = d.db.Where("deleted_status = 0").First(&area, id).Error; err != nil {
		log.Error("RawArea(%v) error(%v)", id, err)
	}
	return
}

// TxRawArea 获取区域
func (d *Dao) TxRawArea(c context.Context, tx *gorm.DB, id int64) (area *model.Area, err error) {
	area = new(model.Area)
	if err = tx.Where("deleted_status = 0").First(&area, id).Error; err != nil {
		log.Error("TxRawArea(%v) error(%v)", id, err)
	}
	return
}

// TxRawDeletedAreaByAID 通过场地ID和自定义区域编号获取已删除的区域
func (d *Dao) TxRawDeletedAreaByAID(c context.Context, tx *gorm.DB, aid string, place int64) (area *model.Area, err error) {
	area = new(model.Area)
	if res := tx.Where("deleted_status = 1").Where("a_id = ? AND place = ?", aid, place).First(&area); res.Error != nil {
		if res.RecordNotFound() {
			return nil, nil
		}
		err = res.Error
		log.Error("TxRawAreaByAID(%v, %v) error(%v)", aid, place, err)
	}
	return
}

// TxAddArea 添加区域信息（事务）
func (d *Dao) TxAddArea(c context.Context, tx *gorm.DB, area *model.Area) (err error) {
	if res := tx.Create(area); res.Error != nil {
		log.Error("添加区域信息失败:%s", res.Error)
		err = ecode.NotModified
		return
	}
	return
}

// TxUpdateArea 编辑区域信息（事务）
func (d *Dao) TxUpdateArea(c context.Context, tx *gorm.DB, area *model.Area) (err error) {
	if err = tx.Table("area").Where("id = ?", area.ID).Updates(
		map[string]interface{}{
			"a_id":           area.AID,
			"name":           area.Name,
			"place":          area.Place,
			"deleted_status": area.DeletedStatus,
		}).Error; err != nil {
		log.Error("更新区域信息（ID:%d）失败:%s", area.ID, err)
		err = ecode.NotModified
		return
	}
	return
}

// TxDelArea 软删除区域（事务）
func (d *Dao) TxDelArea(c context.Context, tx *gorm.DB, id int64) (err error) {
	if err = tx.Table("area").Where("id = ?", id).Updates(
		map[string]interface{}{
			"deleted_status": 1,
		}).Error; err != nil {
		log.Error("删除区域信息失败:%s", err)
		err = ecode.NotModified
		return
	}
	return
}
