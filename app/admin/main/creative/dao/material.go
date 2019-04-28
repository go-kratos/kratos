package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/creative/model/material"
	"go-common/library/log"
)

// CategoryByID .
func (d *Dao) CategoryByID(c context.Context, id int64) (cate *material.Category, err error) {
	cate = &material.Category{}
	if err = d.DB.Where("id=?", id).First(&cate).Error; err != nil {
		log.Error("d.CategoryByID.Find error(%v)", err)
		return
	}
	return
}

// BindWithCategory .
func (d *Dao) BindWithCategory(c context.Context, MaterialID, CategoryID, index int64) (id int64, err error) {
	var state int
	cate := &material.WithCategory{}
	if err = d.DB.Where("material_id=?", MaterialID).First(&cate).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("d.BindWithCategory.Find error(%v)", err)
		return
	}
	cate.CategoryID = CategoryID
	cate.MaterialID = MaterialID
	cate.Index = index
	if err != nil && err == gorm.ErrRecordNotFound {
		//添加关联
		if CategoryID == 0 {
			return
		}
		if err = d.DB.Create(cate).Error; err != nil {
			log.Error("BindWithCategory  Create error(%+v)", err)
			return
		}
	} else {
		if CategoryID == 0 {
			//删除关联
			state = material.StateOff
		} else {
			state = material.StateOn
		}
		if err = d.DB.Model(&material.WithCategory{}).Where("id=?", cate.ID).Update(cate).Update(map[string]int{"state": state}).Error; err != nil {
			log.Error("dao BindWithCategory error(%v)", err)
			return
		}
	}
	id = cate.ID
	err = nil
	return
}
