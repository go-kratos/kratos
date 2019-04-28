package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/manager/model"
	"go-common/library/ecode"
)

// AddType .
func (d *Dao) AddType(c context.Context, tt *model.TagType) (err error) {
	tx := d.db.Begin()
	if tx.Error != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Create(tt).Error; err != nil {
		return
	}
	if len(tt.Rids) > 0 {
		values := []string{}
		valueArgs := []interface{}{}
		for _, rid := range tt.Rids {
			values = append(values, "(?,?)")
			valueArgs = append(valueArgs, tt.ID, rid)
		}
		stmt := fmt.Sprintf("INSERT INTO manager_tag_type_role(tid,rid) VALUES %s ON DUPLICATE KEY UPDATE tid=values(tid),rid=values(rid)", strings.Join(values, ","))
		if err = tx.Exec(stmt, valueArgs...).Error; err != nil {
			return
		}
	}
	err = tx.Commit().Error
	return
}

// UpdateTypeName .
func (d *Dao) UpdateTypeName(c context.Context, tt *model.TagType) error {
	return d.db.Table("manager_tag_type").Where("id = ?", tt.ID).Update("name", tt.Name).Error
}

// UpdateType .
func (d *Dao) UpdateType(c context.Context, tt *model.TagType) (err error) {
	if len(tt.Rids) <= 0 {
		return
	}
	values := []string{}
	valueArgs := []interface{}{}
	for _, rid := range tt.Rids {
		values = append(values, "(?,?)")
		valueArgs = append(valueArgs, tt.ID, rid)
	}
	stmt := fmt.Sprintf("INSERT INTO manager_tag_type_role(tid,rid) VALUES %s ON DUPLICATE KEY UPDATE tid=values(tid),rid=values(rid)", strings.Join(values, ","))
	err = d.db.Exec(stmt, valueArgs...).Error
	return
}

// DeleteNonRole .
func (d *Dao) DeleteNonRole(c context.Context, tt *model.TagType) error {
	if len(tt.Rids) <= 0 {
		stmt := "DELETE FROM manager_tag_type_role WHERE tid=?"
		return d.db.Exec(stmt, tt.ID).Error
	}
	stmt := "DELETE FROM manager_tag_type_role WHERE id IN (SELECT id FROM (select id from manager_tag_type_role WHERE tid=? AND rid NOT IN (?)) as temp)"
	return d.db.Exec(stmt, tt.ID, tt.Rids).Error
}

// DeleteType .
func (d *Dao) DeleteType(c context.Context, td *model.TagTypeDel) error {
	return d.db.Where("id = ?", td.ID).Delete(&model.TagType{}).Error
}

// AddTag .
func (d *Dao) AddTag(c context.Context, t *model.Tag) (err error) {
	maxTagID := int64(0)
	if maxTagID, err = d.maxTagIDByBid(c, t.Bid); err != nil {
		return
	}
	maxTagID++
	t.TagID = maxTagID
	err = d.db.Create(t).Error
	return
}

// maxTagIDByBid .
func (d *Dao) maxTagIDByBid(c context.Context, bid int64) (maxTagID int64, err error) {
	var tagID interface{}
	err = d.db.Table("manager_tag").Select("max(tag_id) as max_tag_id").Where("bid = ?", bid).Row().Scan(&tagID)
	if tagID == nil {
		maxTagID = int64(0)
		return
	}
	maxTagID = tagID.(int64)
	return
}

// UpdateTag .
func (d *Dao) UpdateTag(c context.Context, t *model.Tag) error {
	return d.db.Table("manager_tag").Where("id = ?", t.ID).
		Updates(map[string]interface{}{"bid": t.Bid, "tid": t.Tid, "rid": t.Rid, "name": t.Name, "weight": t.Weight, "description": t.Description}).Error
}

// AddControl .
func (d *Dao) AddControl(c context.Context, tc *model.TagControl) error {
	return d.db.Create(tc).Error
}

// UpdateControl .
func (d *Dao) UpdateControl(c context.Context, tc *model.TagControl) error {
	return d.db.Table("manager_tag_control").Where("id = ?", tc.ID).
		Update(map[string]interface{}{
			"tid":         tc.Tid,
			"name":        tc.Name,
			"title":       tc.Title,
			"weight":      tc.Weight,
			"component":   tc.Component,
			"placeholder": tc.Placeholder,
			"required":    tc.Required,
		}).Error
}

// BatchUpdateState .
func (d *Dao) BatchUpdateState(c context.Context, b *model.BatchUpdateState) (err error) {
	return d.db.Table("manager_tag").Where("id IN (?)", b.IDs).Update("state", b.State).Error
}

// TagList .
func (d *Dao) TagList(c context.Context, t *model.SearchTagParams) (res []*model.Tag, err error) {
	db := d.db.Table("manager_tag")
	if t.Bid != -1 {
		db = db.Where("bid = ?", t.Bid)
	}
	if t.KeyWord != "" {
		db = db.Where("name LIKE ?", "%"+t.KeyWord+"%")
	}
	if t.Tid != -1 {
		db = db.Where("tid = ?", t.Tid)
	}
	if t.Rid != -1 {
		db = db.Where("rid = ?", t.Rid)
	}
	if t.State != -1 {
		db = db.Where("state = ?", t.State)
	}
	if t.UName != "" {
		db = db.Where("uid = ?", t.UID)
	}
	if t.Order != "" {
		db = db.Order(t.Order+" "+t.Sort, true)
	}
	err = db.Find(&res).Error
	return
}

// TagControl .
func (d *Dao) TagControl(c context.Context, tc *model.TagControlParam) (res []*model.TagControl, err error) {
	db := d.db.Table("manager_tag_control")
	if tc.BID != 0 {
		db = db.Where("bid = ?", tc.BID)
	}
	if tc.TID != 0 {
		db = db.Where("tid = ?", tc.TID)
	}
	err = db.Find(&res).Error
	return
}

// AttrList .
func (d *Dao) AttrList(c context.Context, bid int64) (res *model.TagBusinessAttr, err error) {
	res = &model.TagBusinessAttr{}
	err = d.db.Where("bid = ?", bid).First(res).Error
	if err == ecode.NothingFound {
		err = nil
	}
	return
}

// InsertAttr .
func (d *Dao) InsertAttr(c context.Context, tba *model.TagBusinessAttr) error {
	return d.db.Create(tba).Error
}

// AttrUpdate .
func (d *Dao) AttrUpdate(c context.Context, tba *model.TagBusinessAttr) error {
	return d.db.Table("manager_tag_business_attr").Where("bid = ?", tba.Bid).Update("button", tba.Button).Error
}

// TypeByIDs .
func (d *Dao) TypeByIDs(c context.Context, ids []int64) (res map[int64]*model.TagType, err error) {
	t := []*model.TagType{}
	res = make(map[int64]*model.TagType)
	err = d.db.Where("id IN (?)", ids).Find(&t).Error
	for _, r := range t {
		res[r.ID] = r
	}
	return
}

// TagTypeByBID .
func (d *Dao) TagTypeByBID(c context.Context, bid int64) (res []*model.TagType, err error) {
	err = d.db.Where("bid = ?", bid).Find(&res).Error
	return
}

// TagTypeRoleByTids .
func (d *Dao) TagTypeRoleByTids(c context.Context, tids []int64) (res []*model.TagTypeRole, err error) {
	err = d.db.Where("tid IN (?)", tids).Find(&res).Error
	return
}

// TagByType accross type get tag list .
func (d *Dao) TagByType(c context.Context, tid int64) (res []*model.Tag, err error) {
	err = d.db.Where("tid = ?", tid).Find(&res).Error
	return
}
