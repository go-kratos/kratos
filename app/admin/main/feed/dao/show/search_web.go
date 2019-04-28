package show

import (
	"encoding/json"
	"fmt"

	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/model/show"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// SearchWebAdd add search web
func (d *Dao) SearchWebAdd(param *show.SearchWebAP) (err error) {
	var (
		querys []*show.SearchWebQuery
	)
	if param.Query != "" {
		if err = json.Unmarshal([]byte(param.Query), &querys); err != nil {
			return
		}
	}
	tx := d.DB.Begin()
	if err = tx.Model(&show.SearchWeb{}).Create(param).Error; err != nil {
		log.Error("SearchWebAdd tx.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	if len(querys) > 0 {
		if err = tx.Model(&show.SearchWeb{}).Exec(show.BatchAddQuerySQL(param.ID, querys)).Error; err != nil {
			log.Error("SearchWebAdd tx.Model Exec(%+v) error(%v)", param, err)
			err = tx.Rollback().Error
			return
		}
	}
	if err = tx.Commit().Error; err != nil {
		return
	}
	return
}

// SearchWebUpdate update
func (d *Dao) SearchWebUpdate(param *show.SearchWebUP) (err error) {
	var (
		newQuerys []*show.SearchWebQuery
	)
	if param.Query != "" {
		if err = json.Unmarshal([]byte(param.Query), &newQuerys); err != nil {
			return
		}
	}
	tx := d.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("dao.SearchWebUpdate.DB.Begin error(%v)", err)
		return
	}
	if err = tx.Model(&show.SearchWeb{}).Update(param).Error; err != nil {
		log.Error("dao.SearchWebUpdate (%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	var (
		mapOldCData, mapNewQData    map[int64]*show.SearchWebQuery
		upQData, addQData, oldQData []*show.SearchWebQuery
		delQData                    []int64
	)
	if len(newQuerys) > 0 {
		if err = d.DB.Model(&show.SearchWeb{}).Where("sid=?", param.ID).Where("deleted=?", common.NotDeleted).Find(&oldQData).Error; err != nil {
			log.Error("dao.SearchWebUpdate Find Old data (%+v) error(%v)", param.ID, err)
			return
		}
		mapOldCData = make(map[int64]*show.SearchWebQuery, len(oldQData))
		for _, v := range oldQData {
			mapOldCData[v.ID] = v
		}
		//新数据在老数据中 更新老数据。新的数据不在老数据 添加新数据
		for _, qData := range newQuerys {
			if _, ok := mapOldCData[qData.ID]; ok {
				upQData = append(upQData, qData)
			} else {
				addQData = append(addQData, qData)
			}
		}
		mapNewQData = make(map[int64]*show.SearchWebQuery, len(newQuerys))
		for _, v := range newQuerys {
			mapNewQData[v.ID] = v
		}
		//老数据在新数据中 上面已经处理。老数据不在新数据中 删除老数据
		for _, qData := range oldQData {
			if _, ok := mapNewQData[qData.ID]; !ok {
				delQData = append(delQData, qData.ID)
			}
		}
		if len(upQData) > 0 {
			if err = tx.Model(&show.SearchWebQuery{}).Exec(show.BatchEditQuerySQL(upQData)).Error; err != nil {
				log.Error("dao.SearchWebUpdate tx.Model Exec(%+v) error(%v)", upQData, err)
				err = tx.Rollback().Error
				return
			}
		}
		if len(delQData) > 0 {
			if err = tx.Model(&show.SearchWebQuery{}).Where("id IN (?)", delQData).Updates(map[string]interface{}{"deleted": common.Deleted}).Error; err != nil {
				log.Error("dao.SearchWebUpdate Updates(%+v) error(%v)", delQData, err)
				err = tx.Rollback().Error
				return
			}
		}
		if len(addQData) > 0 {
			if err = tx.Model(&show.SearchWebQuery{}).Exec(show.BatchAddQuerySQL(param.ID, addQData)).Error; err != nil {
				log.Error("EditContest s.dao.DB.Model Create(%+v) error(%v)", addQData, err)
				err = tx.Rollback().Error
				return
			}
		}
	} else {
		if err = tx.Model(&show.SearchWebQuery{}).Where("sid IN (?)", param.ID).Updates(map[string]interface{}{"deleted": common.Deleted}).Error; err != nil {
			log.Error("dao.SearchWebUpdate Updates(%+v) error(%v)", param.ID, err)
			err = tx.Rollback().Error
			return
		}
	}
	err = tx.Commit().Error
	return
}

// SearchWebDelete delete search web
func (d *Dao) SearchWebDelete(id int64) (err error) {
	up := map[string]interface{}{
		"deleted": common.Deleted,
	}
	tx := d.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("dao.SearchWebDelete.DB.Begin error(%v)", err)
		return
	}
	if err = tx.Model(&show.SearchWeb{}).Where("id = ?", id).Update(up).Error; err != nil {
		log.Error("dao.show.SearchWebDelete(%+v) error(%v)", id, err)
		err = tx.Rollback().Error
		return
	}
	if err = tx.Model(&show.SearchWebQuery{}).Where("sid = (?)", id).Updates(map[string]interface{}{"deleted": common.Deleted}).Error; err != nil {
		log.Error("dao.SearchWebDelete Updates(%+v) error(%v)", id, err)
		err = tx.Rollback().Error
		return
	}
	err = tx.Commit().Error
	return
}

// SearchWebOption option search web
func (d *Dao) SearchWebOption(up *show.SearchWebOption) (err error) {
	if err = d.DB.Model(&show.SearchWebOption{}).Update(up).Error; err != nil {
		log.Error("dao.SearchWebOption Updates(%+v) error(%v)", up, err)
	}
	return
}

// SWTimeValid search web time validate
func (d *Dao) SWTimeValid(param *show.SWTimeValid) (count int, err error) {
	query := d.DB.Table("search_web_query").
		Select("search_web_query.id").
		Joins("left join search_web ON search_web.id = search_web_query.sid").
		Where("value = ?", param.Query).
		Where("priority = ?", param.Priority).
		Where("`check` in (?)", []int{common.Verify, common.Pass, common.Valid}).
		Where("stime <= ?", param.ETime).
		Where("etime >= ?", param.STime).
		Where("search_web_query.deleted = 0").
		Where("search_web.deleted = 0")
	if param.ID != 0 {
		query = query.Where("search_web.id != ?", param.ID)
	}
	if err = query.Count(&count).Error; err != nil {
		log.Error("dao.SWTimeValid Count error(%v)", err)
	}
	return
}

//SWFindByID search web table value find by id
func (d *Dao) SWFindByID(id int64) (value *show.SearchWeb, err error) {
	value = &show.SearchWeb{}
	w := map[string]interface{}{
		"deleted": common.NotDeleted,
		"id":      id,
	}
	if err = d.DB.Model(&show.SearchWeb{}).Where(w).Find(value).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = fmt.Errorf("ID为%d的数据不存在", id)
			return
		}
		return
	}
	return
}
