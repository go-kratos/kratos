package dao

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/live/resource/model"
	"go-common/library/log"

	"github.com/siddontang/go-mysql/mysql"
)

const (
	_addResourceSQL = "INSERT INTO resource (`platform`,`build`,`limit_type`,`start_time`,`end_time`,`type`,`title`,`image_info`) values(?,?,?,?,?,?,?,?)"
)

// AddResource -
func (d *Dao) AddResource(c context.Context, insert *model.Resource) (row int64, err error) {
	types := strings.Split(insert.Type, ",")
	for _, t := range types {
		insert.Type = t
		row, err = d.addDBResource(c, insert)
		if err != nil {
			return row, err
		}
	}
	return
}

// EditResource -
func (d *Dao) EditResource(c context.Context, id int64, update map[string]interface{}) (row int64, err error) {
	return d.editDBResource(c, id, update)
}

// GetResourceList -
func (d *Dao) GetResourceList(c context.Context, typ string, page int64, pageSize int64) (resp []model.Resource, err error) {
	return d.getDBResourceList(c, typ, page, pageSize)
}

// GetResourceListEx -
func (d *Dao) GetResourceListEx(c context.Context, typ []string, page int64, pageSize int64, devPlatform string, status string, startTime string, endTime string) (resp []model.Resource, count int64, err error) {
	return d.getDBResourceListEx(c, typ, page, pageSize, devPlatform, status, startTime, endTime)
}

// OfflineResource -
func (d *Dao) OfflineResource(c context.Context, id int64) (row int64, err error) {
	return d.offlineDBResource(c, id)
}

// SelectById -
func (d *Dao) SelectById(c context.Context, id int64) (resp *model.Resource, err error) {
	return d.selectDBById(c, id)
}

// GetInfo -
func (d *Dao) GetInfo(ctx context.Context, typ string, platform string, build int64) (resp *model.Resource, err error) {
	inst := rand.Intn(d.c.CacheInstCnt)
	res, ok := d.sCache[inst].Get(cacheResourceKey(typ, platform, build))
	if !ok {
		resp, err = d.getDBInfo(ctx, typ, platform, build)
		if err != nil {
			return
		}
		var resNew []model.Resource
		resNew = append(resNew, *resp)
		d.sCache[inst].Put(cacheResourceKey(typ, platform, build), resNew)
		return
	}
	r := res.([]model.Resource)
	if len(r) > 0 {
		return &r[0], nil
	}
	return nil, nil
}

// GetBanner -
func (d *Dao) GetBanner(ctx context.Context, platform string, build int64, t string) (resp []model.Resource, err error) {
	inst := rand.Intn(d.c.CacheInstCnt)
	res, ok := d.sCache[inst].Get(cacheResourceKey(t, platform, build))
	if !ok {
		resp, err = d.getDBBanner(ctx, platform, build, t)
		if err != nil {
			return
		}
		d.sCache[inst].Put(cacheResourceKey(t, platform, build), resp)
		return
	}
	resp = res.([]model.Resource)
	return
}

// SelectByTypeAndPlatform -
func (d *Dao) SelectByTypeAndPlatform(ctx context.Context, typ string, platform string) (resp *model.Resource, err error) {
	return d.selectDBByTypeAndPlatform(ctx, typ, platform)
}

// GetDBCount -
func (d *Dao) GetDBCount(ctx context.Context, typ string) (resp int64, err error) {
	return d.getDBCount(ctx, typ)
}

// get data from db source
// addSResource add resource to mysql
func (d *Dao) addDBResource(c context.Context, insert *model.Resource) (row int64, err error) {
	if insert == nil {
		return
	}
	var reply sql.Result
	if _, err = d.db.Begin(c); err != nil {
		log.Error("db.begin error(%v)", err)
		return
	}
	reply, err = d.db.Exec(c, _addResourceSQL, insert.Platform, insert.Build, insert.LimitType, insert.StartTime, insert.EndTime, insert.Type, insert.Title, insert.ImageInfo)
	if err != nil {
		log.Error("resource.addSResource d.db.Exec err: %v", err)
		return
	}
	row, err = reply.LastInsertId()
	return
}

func (d *Dao) editDBResource(c context.Context, id int64, update map[string]interface{}) (row int64, err error) {
	if update == nil || id < 0 {
		return
	}
	var tx = d.rsDB
	tableInfo := &model.Resource{}
	var reply = tx.Model(tableInfo).Where("id=?", id).Update(update)
	log.Info("effected rows: %d, id : %d", reply.RowsAffected, id)
	if reply.Error != nil {
		log.Error("resource.editResource error: %v", err)
		return
	}
	row = reply.RowsAffected
	return
}

func (d *Dao) getDBResourceList(c context.Context, typ string, page int64, pageSize int64) (resp []model.Resource, err error) {
	if typ == "" {
		return
	}
	var tx = d.rsDBReader
	tableInfo := &model.Resource{}
	err = tx.Model(tableInfo).
		Select("`id`,`platform`,`build`,`limit_type`,`start_time`,`end_time`,`title`,`image_info`").
		Where("type=?", typ).
		Order("id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&resp).Error
	if err != nil {
		log.Error("resource.editResource error: %v", err)
		return
	}
	return
}

func (d *Dao) getDBResourceListEx(c context.Context, typ []string, page int64, pageSize int64, devPlatform string, status string, startTime string, endTime string) (resp []model.Resource, count int64, err error) {
	if len(typ) <= 0 {
		return
	}
	whereStr := "1=1"
	for i, t := range typ {
		if i != 0 {
			whereStr += " or "
		} else {
			whereStr += " and ("
		}
		whereStr += fmt.Sprintf("`type` like \"%s%%\"", mysql.Escape(t))
	}
	whereStr += ")"
	if devPlatform != "" {
		whereStr += fmt.Sprintf(" and `platform`=\"%s\"", mysql.Escape(devPlatform))
	}
	if status != "" {
		var i int
		if i, err = strconv.Atoi(status); err != nil {
			return
		}
		now := time.Now().Format("2006-01-02 15:04:05")
		switch i {
		case 0:
			whereStr += fmt.Sprintf(" and `start_time`>=\"%s\"", now)
		case 1:
			whereStr += fmt.Sprintf(" and `start_time`<=\"%s\" and `end_time`>=\"%s\"", now, now)
		case -1:
			whereStr += fmt.Sprintf(" and `end_time`<=\"%s\"", now)
		}
	}
	if startTime != "" {
		whereStr += fmt.Sprintf(" and `start_time`>=\"%s\"", mysql.Escape(startTime))
	}
	if endTime != "" {
		whereStr += fmt.Sprintf(" and `end_time`<=\"%s\"", mysql.Escape(endTime))
	}
	var tx = d.rsDBReader
	tableInfo := &model.Resource{}
	err = tx.Model(tableInfo).
		Select("`id`,`platform`,`build`,`limit_type`,`start_time`,`end_time`,`title`,`image_info`,`type`").
		Where(whereStr).
		Order("id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&resp).Error
	if err != nil {
		log.Error("resource.editResource error: %v", err)
		return
	}
	err = tx.Model(tableInfo).
		Select("`id`").
		Where(whereStr).
		Count(&count).Error
	if err != nil {
		log.Error("resource.editResource error: %v", err)
		return
	}
	return
}
func (d *Dao) getDBCount(c context.Context, typ string) (resp int64, err error) {
	if typ == "" {
		return
	}
	var tx = d.rsDBReader
	tableInfo := &model.Resource{}
	err = tx.Model(tableInfo).
		Select("`id`").
		Where("type=?", typ).
		Count(&resp).Error
	if err != nil {
		log.Error("resource.editResource error: %v", err)
		return
	}
	return
}
func (d *Dao) offlineDBResource(c context.Context, id int64) (row int64, err error) {
	if id == 0 {
		return
	}
	reply, err := d.SelectById(c, id)
	if err != nil || reply == nil {
		log.Error("resource.OfflineResource select error: %v reply %v", err, reply)
		return
	}
	var tx = d.rsDB
	tableInfo := &model.Resource{}
	var updateReply = tx.Model(tableInfo).Where("id=?", id).Update("end_time", time.Now())
	if updateReply.Error != nil {
		log.Error("resource.OfflineResource update error: %v", err)
		return
	}
	return
}

func (d *Dao) selectDBById(c context.Context, id int64) (resp *model.Resource, err error) {
	if id == 0 {
		return
	}
	resp = &model.Resource{}
	var tx = d.rsDBReader
	var reply = tx.Model(&model.Resource{}).Where("id=?", id).Find(resp)
	if reply.Error != nil {
		log.Error("resource.SelectById error: %v", err)
		return
	}
	return
}

func (d *Dao) getDBInfo(ctx context.Context, typ string, platform string, build int64) (resp *model.Resource, err error) {
	if platform == "" || build == 0 {
		return
	}
	resp = &model.Resource{}
	var tx = d.rsDBReader
	now := time.Now()
	var reply = tx.Model(&model.Resource{}).Where("start_time<? and end_time>? and type=? and `platform`=? and ((`limit_type`=0 and `build`<=?) or (`limit_type`=1 and `build`=?) or (`limit_type`=2 and `build`>=?))", now, now, typ, platform, build, build, build).Limit(1).Find(resp)
	if reply.Error != nil {
		log.Error("resource.GetInfo error: %v", err)
		return
	}
	return
}

func (d *Dao) getDBBanner(ctx context.Context, platform string, build int64, t string) (resp []model.Resource, err error) {
	if platform == "" || build == 0 || t == "" {
		return
	}
	var tx = d.rsDBReader
	var reply = tx.Model(&model.Resource{}).Where("`start_time`<? and `end_time`>? and `type`=? and (`platform`='' or `platform`=?) and ((`limit_type`=0 and `build`<=?) or (`limit_type`=1 and `build`=?) or (`limit_type`=2 and `build`>=?))", time.Now(), time.Now(), t, platform, build, build, build).Order("mtime DESC").Find(&resp)
	if reply.Error != nil {
		log.Error("resource.GetBanner error: %v", err)
		return
	}
	return
}

func (d *Dao) selectDBByTypeAndPlatform(ctx context.Context, typ string, platform string) (resp *model.Resource, err error) {
	if typ == "" || platform == "" {
		return
	}
	resp = &model.Resource{}
	var tx = d.rsDBReader
	now := time.Now()
	var reply = tx.Model(&model.Resource{}).Where("`type`=? and `end_time`>? and `platform`=? ", typ, now, platform).Limit(1).Find(resp)
	if reply.Error != nil {
		resp = nil
		log.Error("resource.SelectByTypeAndPlatform error: %v", err)
		return
	}
	return
}
