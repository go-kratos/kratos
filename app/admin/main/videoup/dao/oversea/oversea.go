package manager

import (
	"context"

	"go-common/app/admin/main/videoup/model/oversea"
	"go-common/library/log"
)

// UpPolicyRelation update or into archive_relation.
func (d *Dao) UpPolicyRelation(c context.Context, aid, gid int64) (relation *oversea.ArchiveRelation, err error) {
	var assign = map[string]interface{}{
		"policy_id": gid,
		"aid":       aid,
	}
	relation = &oversea.ArchiveRelation{}
	if err = d.OverseaDB.Where("aid=?", aid).Assign(assign).FirstOrCreate(&relation).Error; err != nil {
		log.Error("d.UpPolicyRelation.FirstOrCreate error(%v)", err)
		return
	}
	return
}

// PolicyRelation get archive policy group relation.
func (d *Dao) PolicyRelation(c context.Context, aid int64) (relation *oversea.ArchiveRelation, err error) {
	relation = &oversea.ArchiveRelation{}
	res := d.OverseaDB.Where("aid=?", aid).Find(&relation)
	if res.RecordNotFound() {
		relation = nil
		return
	}
	err = res.Error
	return
}

// PolicyGroups get policy group
func (d *Dao) PolicyGroups(c context.Context, uid, id int64, gType, state int8, count, page int64, order, sort string) (groups []*oversea.PolicyGroup, total int64, err error) {
	var (
		db     = d.OverseaDB.Model(&groups)
		orders = map[string]int{
			"mtime": 1,
		}
	)
	db = db.Where("is_global=?", 1)
	if uid > 0 {
		db = db.Where("uid=?", uid)
	}
	if id > 0 {
		db = db.Where("id=?", id)
	}
	if gType > 0 {
		db = db.Where("type=?", gType)
	}
	if state >= 0 {
		db = db.Where("state=?", state)
	}
	if order != "" && sort != "" {
		if _, ok := orders[order]; ok {
			db = db.Order(order + " " + sort)
		}
	}
	if count <= 0 {
		count = 20
	}
	if page <= 0 {
		page = 1
	}
	db.Count(&total)
	db = db.Offset((page - 1) * count)
	db = db.Limit(count)
	if err = db.Find(&groups).Error; err != nil {
		log.Error("d.PolicyGroups.Find error(%v)", err)
	}
	err = d.ItemsByGroup(groups)
	return
}

// PolicyGroupsByIds get policy groups by ids
func (d *Dao) PolicyGroupsByIds(c context.Context, ids []int64) (groups []*oversea.PolicyGroup, err error) {
	if err = d.OverseaDB.Where(ids).Find(&groups).Error; err != nil {
		log.Error("d.PolicyGroupsByIds.Find error(%v)", err)
	}
	return
}

// PolicyGroup get policy group by id
func (d *Dao) PolicyGroup(c context.Context, id int64) (group *oversea.PolicyGroup, err error) {
	var (
		groups []*oversea.PolicyGroup
	)
	group = &oversea.PolicyGroup{}
	res := d.OverseaDB.Where("id=?", id).Find(&group)
	if res.RecordNotFound() {
		group = nil
		return
	}
	err = res.Error
	if err != nil {
		log.Error("d.PolicyGroup.Find error(%v)", err)
		return
	}
	groups = append(groups, group)
	if err = d.ItemsByGroup(groups); err != nil {
		log.Error("d.ItemsByGroup.Find error(%v)", err)
		return
	}
	if len(groups) != 0 {
		group = groups[0]
	}
	return
}

// ArchiveGroups get archive's policy groups
func (d *Dao) ArchiveGroups(c context.Context, aid int64) (groups []*oversea.PolicyGroup, err error) {
	var (
		db        = d.OverseaDB
		relations []*oversea.ArchiveRelation
		gids      []int64
	)
	err = db.Where("aid=?", aid).Find(&relations).Error
	if err != nil {
		log.Error("d.ArchiveGroups.Find error(%v)", err)
		return
	}
	gids = make([]int64, len(relations))
	for i, v := range relations {
		gids[i] = v.GroupID
	}
	db = d.OverseaDB
	err = db.Where(gids).Find(&groups).Error
	if err != nil {
		log.Error("d.ArchiveGroups.Find error(%v)", err)
		return
	}
	err = d.ItemsByGroup(groups)
	return
}

// ItemsByGroup get policy items into group
func (d *Dao) ItemsByGroup(groups []*oversea.PolicyGroup) (err error) {
	var (
		items   []*oversea.PolicyItem
		itemMap = make(map[int64][]*oversea.PolicyItem)
	)
	gids := make([]int64, len(groups))
	for i, v := range groups {
		gids[i] = v.ID
	}
	db := d.OverseaDB
	err = db.Where("group_id in (?) and state=?", gids, oversea.StateOK).Find(&items).Error
	if err != nil {
		log.Error("d.ArchiveGroups.Find error(%v)", err)
		return
	}
	for _, v := range items {
		itemMap[v.GroupID] = append(itemMap[v.GroupID], v)
	}
	for i, g := range groups {
		if _, ok := itemMap[g.ID]; !ok {
			groups[i].Items = make([]*oversea.PolicyItem, 0)
			continue
		}
		groups[i].Items = itemMap[g.ID]
	}
	return
}

// AddPolicyGroup add policy group
func (d *Dao) AddPolicyGroup(c context.Context, group *oversea.PolicyGroup) (err error) {
	var (
		db = d.OverseaDB
	)
	group.IsGlobal = 1
	group.Aid = 0
	group.State = 1
	err = db.Create(&group).Error
	if err != nil {
		group = &oversea.PolicyGroup{}
		log.Error("d.AddPolicyGroup.Create error(%v)", err)
	}
	return
}

// UpdatePolicyGroup update policy group
func (d *Dao) UpdatePolicyGroup(c context.Context, id int64, attrs map[string]interface{}) (err error) {
	var (
		db = d.OverseaDB
	)
	err = db.Model(&oversea.PolicyGroup{}).Where("id=?", id).Update(attrs).Error
	if err != nil {
		log.Error("d.UpdatePolicyGroup.Update error(%v)", err)
	}
	return
}

// UpdatePolicyGroups multi update policy groups
func (d *Dao) UpdatePolicyGroups(c context.Context, ids []int64, attrs map[string]interface{}) (err error) {
	var (
		db = d.OverseaDB
	)
	err = db.Model(&oversea.PolicyGroup{}).Where(ids).Update(attrs).Error
	if err != nil {
		log.Error("d.UpdatePolicyGroup.Update error(%v)", err)
	}
	return
}

// PolicyItems get policy items
func (d *Dao) PolicyItems(c context.Context, gid int64) (items []*oversea.PolicyItem, err error) {
	err = d.OverseaDB.Where("group_id=? AND state=?", gid, oversea.StateOK).Find(&items).Error
	return
}

// ZoneIDs get zone ids by area ids
func (d *Dao) ZoneIDs(c context.Context, aids []int64) (ids []int64, err error) {
	var items []*oversea.Zone
	if err = d.OverseaDB.Where(aids).Find(&items).Pluck("zone_id", &ids).Error; err != nil {
		log.Error("d.ZoneIDs.Find error(%v)", err)
	}
	return
}

// AddPolicies add policy items
func (d *Dao) AddPolicies(c context.Context, policies []oversea.PolicyItem) (err error) {
	var assign = map[string]interface{}{
		"group_id":  0,
		"play_auth": 0,
		"down_auth": 0,
		"area_id":   "",
		"zone_id":   "",
	}
	for _, v := range policies {
		if v.ID > 0 {
			assign["group_id"] = v.GroupID
			assign["play_auth"] = v.PlayAuth
			assign["down_auth"] = v.DownAuth
			assign["area_id"] = v.AreaID
			assign["zone_id"] = v.ZoneID
			err = d.OverseaDB.Model(&v).Where("id=?", v.ID).Update(assign).Error
		} else {
			err = d.OverseaDB.Create(&v).Error
		}
		if err != nil {
			log.Error("d.AddPolicies.FirstOrCreate error(%v)", err)
			return
		}
	}
	return
}

// DelPolices soft delete policy items
func (d *Dao) DelPolices(c context.Context, gid int64, ids []int64) (err error) {
	err = d.OverseaDB.Debug().Model(&oversea.PolicyItem{}).Where(ids).Where("group_id=?", gid).Update("state", oversea.StateDeleted).Error
	if err != nil {
		log.Error("d.DelPolices.Update error(%v)", err)
	}
	return
}
