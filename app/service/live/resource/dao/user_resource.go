package dao

import (
	"context"
	"database/sql"
	"go-common/app/service/live/resource/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	addSql    = "INSERT INTO `user_resource`(`res_type`,`custom_id`,`title`,`url`,`weight`,`status`,`creator`) values (?,?,?,?,?,?,?);"
	rowFields = "`id`, `res_type`,`custom_id`,`title`,`url`,`weight`,`status`,`creator`,UNIX_TIMESTAMP(`ctime`), UNIX_TIMESTAMP(`mtime`)"
)

// AddUserResource 添加用户资源到DB
func (d *Dao) AddUserResource(c context.Context, res *model.UserResource) (newRes model.UserResource, err error) {
	if res == nil {
		return
	}

	var reply sql.Result
	if _, err = d.db.Begin(c); err != nil {
		log.Error("初始化DB错误(%v)", err)
		return
	}

	reply, err = d.db.Exec(c, addSql, res.ResType, res.CustomID, res.Title, res.URL, res.Weight, res.Status, res.Creator)
	if err != nil {
		log.Error("执行SQL语句 err: %v", err)
		return
	}

	lastID_, _ := reply.LastInsertId()
	newRes, err = d.GetUserResourceInfoByID(c, int32(lastID_))

	return
}

// EditUserResource 编辑已有资源
func (d *Dao) EditUserResource(c context.Context, resType int32, customID int32, update map[string]interface{}) (effectRow int32, newRes model.UserResource, err error) {
	if update == nil {
		return
	}

	var tx = d.rsDB
	tableInfo := &model.UserResource{}
	var reply = tx.Model(tableInfo).
		Where("`res_type` = ? AND `custom_id` = ?", resType, customID).
		Update(update)

	log.Info("effected rows: %d, res_type : %d custom_id : %d", reply.RowsAffected, resType, customID)
	if reply.Error != nil {
		log.Error("resource.editResource error: %v", err)
		return
	}

	effectRow = int32(reply.RowsAffected)
	newRes, err = d.GetUserResourceInfo(c, resType, customID)

	return
}

// SetUserResourceStatus 设置资源状态
func (d *Dao) SetUserResourceStatus(c context.Context, resType int32, customID int32, status int32) (effectRow int32, err error) {
	update := make(map[string]interface{})
	update["status"] = status

	effectRow, _, err = d.EditUserResource(c, resType, customID, update)
	if err != nil {
		log.Error("修改资源状态: %v", err)
	}

	return
}

// GetMaxCustomID 根据资源类型获取当前最大的资源ID
func (d *Dao) GetMaxCustomID(c context.Context, resType int32) (maxCustomID int32, err error) {
	tableInfo := &model.UserResource{}

	var ret sql.NullInt64

	err = d.rsDB.Model(tableInfo).Debug().
		Select("max(custom_id) as mcid").
		Where("res_type=?", resType).
		Row().Scan(&ret)

	if err != nil {
		log.Error("查找最大的资源ID res_type : %d : %v", resType, err)
		return
	}

	maxCustomID = int32(ret.Int64)
	log.Info("类型为 %d 最大的资源ID是 %d", resType, maxCustomID)

	return
}

// getRowResult Helper方法
func getRowResult(queryResult *gorm.DB) (res model.UserResource, err error) {
	var count int32

	err = queryResult.Count(&count).Error
	if err != nil {
		log.Error("user_resource.GetUserResourceInfoByID %v", err)
		err = ecode.SeltResErr
		return
	}

	if count == 0 {
		log.Info("user_resource.getRowResult 查询结果为空")
		err = ecode.SeltResErr
		return
	}

	var retID, retResType, retCustomID, retWeight, retStatus, retCtime, retMtime sql.NullInt64
	var retTitle, retURL, retCreator sql.NullString

	err = queryResult.Row().Scan(&retID, &retResType, &retCustomID, &retTitle, &retURL, &retWeight, &retStatus, &retCreator, &retCtime, &retMtime)

	if err != nil {
		log.Error("resource.GetUserResourceInfoByID error: %v", err)
		err = ecode.SeltResErr
		return
	}

	res.ID = int32(retID.Int64)
	res.ResType = int32(retResType.Int64)
	res.CustomID = int32(retCustomID.Int64)
	res.Title = retTitle.String
	res.URL = retURL.String
	res.Weight = int32(retWeight.Int64)
	res.Status = int32(retStatus.Int64)
	res.Creator = retCreator.String
	res.Ctime = time.Unix(retCtime.Int64, 0)
	res.Ctime = time.Unix(retMtime.Int64, 0)

	return
}

// GetUserResourceInfo 获取单个配置
func (d *Dao) GetUserResourceInfo(c context.Context, resType int32, customID int32) (res model.UserResource, err error) {
	tableInfo := &model.UserResource{}

	queryResult := d.rsDBReader.Model(tableInfo).Select(rowFields).
		Where("res_type=? AND custom_id=?", resType, customID)

	res, err = getRowResult(queryResult)

	return
}

// GetUserResourceInfoByID 根据ID获取单个配置
func (d *Dao) GetUserResourceInfoByID(c context.Context, id int32) (res model.UserResource, err error) {
	tableInfo := &model.UserResource{}

	queryResult := d.rsDBReader.Model(tableInfo).Select(rowFields).
		Where("id=?", id)

	res, err = getRowResult(queryResult)

	return
}

// ListUserResourceInfo 获取配置列表
func (d *Dao) ListUserResourceInfo(c context.Context, resType int32, page int32, pageSize int32) (list []model.UserResource, err error) {
	var tx = d.rsDBReader
	tableInfo := &model.UserResource{}

	err = tx.Model(tableInfo).
		Select("`id`, `res_type`,`custom_id`,`title`,`url`,`weight`,`status`,`creator`,`ctime`, `mtime`").
		Where("res_type=?", resType).
		Order("id ASC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&list).Error

	if err != nil {
		log.Error("resource.editResource error: %v", err)
		return
	}

	return
}
