package gorm

import (
	"context"

	"go-common/app/admin/main/aegis/model/business"

	"github.com/jinzhu/gorm"
)

// AddBusiness .
func (d *Dao) AddBusiness(c context.Context, e *business.Business) (id int64, err error) {
	err = d.orm.Table("business").Create(&e).Error
	id = e.ID
	return
}

// UpdateBusiness .
func (d *Dao) UpdateBusiness(c context.Context, e *business.Business) (err error) {
	return d.orm.Table("business").Where("id = ?", e.ID).Update(map[string]interface{}{
		"name":      e.Name,
		"desc":      e.Desc,
		"developer": e.Developer,
	}).Error
}

// EnableBusiness .
func (d *Dao) EnableBusiness(c context.Context, id int64) (err error) {
	return d.orm.Table("business").Where("id = ?", id).Update(map[string]interface{}{
		"state": business.StateEnable,
	}).Error
}

// DisableBusiness .
func (d *Dao) DisableBusiness(c context.Context, id int64) (err error) {
	return d.orm.Table("business").Where("id = ?", id).Update(map[string]interface{}{
		"state": business.StateDisable,
	}).Error
}

// Business .
func (d *Dao) Business(c context.Context, id int64) (res *business.Business, err error) {
	res = &business.Business{}
	if err = d.orm.Where("id = ?", id).First(&res).Error; err == gorm.ErrRecordNotFound {
		res = nil
		err = nil
	}
	return
}

// BusinessList .
func (d *Dao) BusinessList(c context.Context, tp int8, ids []int64, onlyEnable bool) (res []*business.Business, err error) {
	res = []*business.Business{}
	db := d.orm
	if len(ids) > 0 {
		db = db.Where("id in (?)", ids)
	}
	if onlyEnable {
		db = db.Where("state=?", business.StateEnable)
	}
	if tp > 0 {
		db = db.Where("type = ?", tp)
	}
	err = db.Find(&res).Error
	return
}
