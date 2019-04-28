package dao

import (
	"context"

	"go-common/app/admin/main/manager/model"
)

// AddCateSecExt .
func (d *Dao) AddCateSecExt(c context.Context, e *model.CateSecExt) (err error) {
	str := "INSERT INTO manager_reason_catesecext (bid,name,type) VALUES (?,?,?) ON DUPLICATE KEY UPDATE bid=values(bid),name=values(name),type=values(type)"
	return d.db.Exec(str, e.BusinessID, e.Name, e.Type).Error
}

// UpdateCateSecExt .
func (d *Dao) UpdateCateSecExt(c context.Context, e *model.CateSecExt) (err error) {
	return d.db.Table("manager_reason_catesecext").Where("id = ?", e.ID).Update("name", e.Name).Error
}

// BanCateSecExt .
func (d *Dao) BanCateSecExt(c context.Context, e *model.CateSecExt) (err error) {
	return d.db.Table("manager_reason_catesecext").Where("id = ?", e.ID).Update("state", e.State).Error
}

// AddAssociation .
func (d *Dao) AddAssociation(c context.Context, e *model.Association) (err error) {
	return d.db.Table("manager_reason_association").Create(&e).Error
}

// UpdateAssociation .
func (d *Dao) UpdateAssociation(c context.Context, e *model.Association) (err error) {
	return d.db.Table("manager_reason_association").Where("id = ?", e.ID).Updates(
		map[string]interface{}{
			"rid":  e.RoleID,
			"cid":  e.CategoryID,
			"sids": e.SecondIDs,
		}).Error
}

// BanAssociation .
func (d *Dao) BanAssociation(c context.Context, e *model.Association) (err error) {
	return d.db.Table("manager_reason_association").Where("id = ?", e.ID).Update("state", e.State).Error
}

// AddReason .
func (d *Dao) AddReason(c context.Context, e *model.Reason) (err error) {
	return d.db.Table("manager_reason").Create(e).Error
}

// UpdateReason .
func (d *Dao) UpdateReason(c context.Context, e *model.Reason) (err error) {
	return d.db.Table("manager_reason").Where("id = ?", e.ID).
		Updates(map[string]interface{}{
			"rid":         e.RoleID,
			"cid":         e.CategoryID,
			"sid":         e.SecondID,
			"state":       e.State,
			"common":      e.Common,
			"uid":         e.UID,
			"description": e.Description,
			"weight":      e.Weight,
			"flag":        e.Flag,
			"lid":         e.LinkID,
			"type_id":     e.TypeID,
			"tid":         e.TagID,
		}).Error
}

// ReasonList .
func (d *Dao) ReasonList(c context.Context, e *model.SearchReasonParams) (res []*model.Reason, err error) {
	db := d.db.Table("manager_reason")
	if e.BusinessID != -1 {
		db = db.Where("bid = ?", e.BusinessID)
	}
	if e.KeyWord != "" {
		db = db.Where("description LIKE ?", "%"+e.KeyWord+"%")
	}
	if e.RoleID != 0 {
		db = db.Where("rid = ?", e.RoleID)
	}
	if e.CategoryID != 0 {
		db = db.Where("cid = ?", e.CategoryID)
	}
	if e.SecondID != 0 {
		db = db.Where("sid = ?", e.SecondID)
	}
	if e.State != -1 {
		db = db.Where("state = ?", e.State)
	}
	if e.UName != "" {
		db = db.Where("uid = ?", e.UID)
	}
	if e.Order != "" {
		db = db.Order(e.Order+" "+e.Sort, true)
	}
	err = db.Find(&res).Error
	return
}

// CateSecByIDs .
func (d *Dao) CateSecByIDs(c context.Context, ids []int64) (res map[int64]string, err error) {
	r := []*model.CateSecExt{}
	res = make(map[int64]string)
	if err = d.db.Table("manager_reason_catesecext").Where("id IN (?)", ids).Find(&r).Error; err != nil {
		return
	}
	for _, cn := range r {
		res[cn.ID] = cn.Name
	}
	return
}

// BatchUpdateReasonState .
func (d *Dao) BatchUpdateReasonState(c context.Context, b *model.BatchUpdateReasonState) (err error) {
	return d.db.Table("manager_reason").Where("id IN (?)", b.IDs).Update("state", b.State).Error
}

// CateSecExtList .
func (d *Dao) CateSecExtList(c context.Context, e *model.CateSecExt) (res []*model.CateSecExt, err error) {
	// Display all record
	db := d.db.Table("manager_reason_catesecext").Where("bid = ? and type = ?", e.BusinessID, e.Type)
	if e.State != -1 {
		db = db.Where("state = ?", e.State)
	}
	err = db.Find(&res).Error
	return
}

// CateSecList .
func (d *Dao) CateSecList(c context.Context, bid int64) (res []*model.CateSecExt, err error) {
	err = d.db.Table("manager_reason_catesecext").Where("bid = ?", bid).Find(&res).Error
	return
}

// AssociationList .
func (d *Dao) AssociationList(c context.Context, state int64, bid int64) (res []*model.Association, err error) {
	db := d.db.Table("manager_reason_association").Where("bid = ?", bid)
	if state != -1 {
		db = db.Where("state = ?", state)
	}
	err = db.Find(&res).Error
	return
}
