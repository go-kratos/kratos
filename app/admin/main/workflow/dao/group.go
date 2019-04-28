package dao

import (
	"context"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GroupByID will select a group by group id
func (d *Dao) GroupByID(c context.Context, gid int64) (g *model.Group, err error) {
	g = &model.Group{}
	if db := d.ReadORM.Table("workflow_group").Where("id=?", gid).First(g); db.Error != nil {
		err = db.Error
		if db.RecordNotFound() {
			g = nil
			err = nil
		} else {
			err = errors.Wrapf(err, "group(%d)", gid)
		}
	}
	return
}

// GroupByOid will select a group by oid and business
func (d *Dao) GroupByOid(c context.Context, oid int64, business int8) (g *model.Group, err error) {
	g = &model.Group{}
	err = d.ReadORM.Table("workflow_group").Where("oid=? AND business=?", oid, business).Find(g).Error
	return
}

// TxGroupsByOidsStates will select a set of groups by oids, business and states
func (d *Dao) TxGroupsByOidsStates(tx *gorm.DB, oids []int64, business, state int8) (groups map[int64]*model.Group, err error) {
	var groupSlice []*model.Group
	groups = make(map[int64]*model.Group, len(oids))
	if len(oids) <= 0 {
		return
	}

	if err = tx.Table("workflow_group").Where("oid IN (?) AND business=?", oids, business).Find(&groupSlice).Error; err != nil {
		return
	}

	for _, g := range groupSlice {
		groups[g.ID] = g
	}

	return
}

// Groups will select a set of groups by group ids
func (d *Dao) Groups(c context.Context, gids []int64) (groups map[int64]*model.Group, err error) {
	var groupSlice []*model.Group
	groups = make(map[int64]*model.Group, len(gids))
	if len(gids) <= 0 {
		return
	}
	if err = d.ReadORM.Table("workflow_group").Where("id IN (?)", gids).Find(&groupSlice).Error; err != nil {
		return
	}
	for _, g := range groupSlice {
		groups[g.ID] = g
	}
	return
}

// TxGroups .
func (d *Dao) TxGroups(tx *gorm.DB, gids []int64) (groups map[int64]*model.Group, err error) {
	var groupSlice []*model.Group
	groups = make(map[int64]*model.Group, len(gids))
	if len(gids) <= 0 {
		return
	}

	if err = tx.Set("gorm:query_option", "FOR UPDATE").Table("workflow_group").Where("id IN (?)", gids).Find(&groupSlice).Error; err != nil {
		return
	}

	for _, g := range groupSlice {
		groups[g.ID] = g
	}
	return
}

// TxUpGroup will update a group
func (d *Dao) TxUpGroup(tx *gorm.DB, oid int64, business int8, tid int64, note string, rid int8) (err error) {
	err = tx.Table("workflow_group").Where("oid=? AND business=?", oid, business).UpdateColumn(map[string]interface{}{
		"tid":  tid,
		"note": note,
		"rid":  rid,
	}).Error
	return
}

// UpGroupRole update tid note rid of groups
func (d *Dao) UpGroupRole(c context.Context, grsp *param.GroupRoleSetParam) (err error) {
	err = d.ORM.Table("workflow_group").Where("id IN (?)", grsp.GID).UpdateColumn(map[string]interface{}{
		"tid":  grsp.TID,
		"note": grsp.Note,
		"rid":  grsp.RID,
	}).Error
	return
}

// TxUpGroupState will update a group state
func (d *Dao) TxUpGroupState(tx *gorm.DB, gid int64, state int8) (err error) {
	err = tx.Table("workflow_group").Where("id=? AND state=0", gid).Update("state", state).Error

	return
}

// TxUpGroupHandling will use gorm update a group handling stat field
func (d *Dao) TxUpGroupHandling(tx *gorm.DB, gid int64, handling int32) (err error) {
	err = tx.Table("workflow_group").Where("id=?", gid).Update("handling", handling).Error
	return
}

// TxBatchUpGroupHandling will update a set of groups handling stat field
func (d *Dao) TxBatchUpGroupHandling(tx *gorm.DB, gids []int64, handling int32) (err error) {
	if len(gids) <= 0 {
		return
	}
	return tx.Table("workflow_group").Where("id IN (?)", gids).Update("handling", handling).Error
}

// TxBatchUpGroupState will update a group state
func (d *Dao) TxBatchUpGroupState(tx *gorm.DB, gids []int64, state int8) (err error) {
	if len(gids) <= 0 {
		return
	}
	return tx.Table("workflow_group").Where("id IN (?) AND state=0", gids).Update("state", state).Error
}

// TxSetGroupStateTid set group state,tid,rid
func (d *Dao) TxSetGroupStateTid(tx *gorm.DB, gids []int64, state, rid int8, tid int64) (err error) {
	if len(gids) <= 0 {
		return
	}
	return tx.Table("workflow_group").Where("id IN (?)", gids).UpdateColumn(map[string]interface{}{
		"state": state,
		"tid":   tid,
		"rid":   rid,
	}).Error
}

// TxSimpleSetGroupState only set group state
func (d *Dao) TxSimpleSetGroupState(tx *gorm.DB, gids []int64, state int8) (err error) {
	if len(gids) <= 0 {
		return
	}
	return tx.Table("workflow_group").Where("id IN (?)", gids).Update("state", state).Error
}
