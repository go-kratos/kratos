package dao

import (
	"context"
	xsql "database/sql"

	"go-common/app/admin/main/vip/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	_vipPrivileges          = "vip_privileges"
	_vipPrivilegesResources = "vip_privileges_resources"
	updateOrderSQL          = "UPDATE vip_privileges a, vip_privileges b SET a.order_num = b.order_num, b.order_num = a.order_num WHERE a.id = ? AND b.id = ?;"
)

// PrivilegeList query .
func (d *Dao) PrivilegeList(c context.Context, langType int8) (res []*model.Privilege, err error) {
	db := d.vip.Table(_vipPrivileges).Where("deleted=0 AND lang_type=?", langType).Order("order_num ASC")
	if err := db.Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// PrivilegeResourcesList query privilege resources .
func (d *Dao) PrivilegeResourcesList(c context.Context) (res []*model.PrivilegeResources, err error) {
	db := d.vip.Table(_vipPrivilegesResources)
	if err := db.Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// UpdateStatePrivilege update state privilege.
func (d *Dao) UpdateStatePrivilege(c context.Context, p *model.Privilege) (a int64, err error) {
	stage := map[string]interface{}{
		"state": p.State,
	}
	db := d.vip.Table(_vipPrivileges).Where("id = ?", p.ID).Updates(stage)
	if err = db.Error; err != nil {
		return
	}
	a = db.RowsAffected
	return
}

// DeletePrivilege dekete privilege.
func (d *Dao) DeletePrivilege(c context.Context, id int64) (a int64, err error) {
	stage := map[string]interface{}{
		"deleted": 1,
	}
	db := d.vip.Table(_vipPrivileges).Where("id = ?", id).Updates(stage)
	if err = db.Error; err != nil {
		return
	}
	a = db.RowsAffected
	return
}

// AddPrivilege add privilege.
func (d *Dao) AddPrivilege(tx *gorm.DB, ps *model.Privilege) (id int64, err error) {
	db := tx.Table(_vipPrivileges).Save(ps)
	if err = db.Error; err != nil {
		return
	}
	id = ps.ID
	return
}

// MaxOrder max priivilege order.
func (d *Dao) MaxOrder(c context.Context) (order int64, err error) {
	p := new(model.Privilege)
	db := d.vip.Table(_vipPrivileges).Order("order_num DESC").First(&p)
	if err = db.Error; err != nil {
		return
	}
	return p.Order, err
}

// AddPrivilegeResources add privilege resources.
func (d *Dao) AddPrivilegeResources(tx *gorm.DB, p *model.PrivilegeResources) (a int64, err error) {
	db := tx.Table(_vipPrivilegesResources).Save(p)
	if err = db.Error; err != nil {
		return
	}
	a = db.RowsAffected
	return
}

// UpdatePrivilege update privilege .
func (d *Dao) UpdatePrivilege(tx *gorm.DB, ps *model.Privilege) (a int64, err error) {
	val := map[string]interface{}{
		"privileges_name": ps.Name,
		"title":           ps.Title,
		"explains":        ps.Explain,
		"privileges_type": ps.Type,
		"operator":        ps.Operator,
	}
	if ps.IconURL != "" {
		val["icon_url"] = ps.IconURL
	}
	if ps.IconGrayURL != "" {
		val["icon_gray_url"] = ps.IconGrayURL
	}
	db := tx.Table(_vipPrivileges).Where("id = ?", ps.ID).Updates(val)
	if err = db.Error; err != nil {
		return
	}
	a = db.RowsAffected
	return
}

// UpdatePrivilegeResources update privilege resources .
func (d *Dao) UpdatePrivilegeResources(tx *gorm.DB, ps *model.PrivilegeResources) (aff int64, err error) {
	stage := map[string]interface{}{
		"link": ps.Link,
	}
	if ps.ImageURL != "" {
		stage["image_url"] = ps.ImageURL
	}
	db := tx.Table(_vipPrivilegesResources).Where("pid = ? AND resources_type = ?", ps.PID, ps.Type).Updates(stage)
	if err = db.Error; err != nil {
		return
	}
	aff = db.RowsAffected
	return
}

// UpdateOrder update privilege order.
func (d *Dao) UpdateOrder(c context.Context, aid, bid int64) (a int64, err error) {
	var res xsql.Result
	if res, err = d.db.Exec(c, updateOrderSQL, aid, bid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if a, err = res.RowsAffected(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
