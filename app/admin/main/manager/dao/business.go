package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/admin/main/manager/model"

	"github.com/jinzhu/gorm"
)

// AddBusiness .
func (d *Dao) AddBusiness(c context.Context, b *model.Business) (err error) {
	maxBid := int64(-1)
	if b.PID != 0 {
		if maxBid, err = d.maxBid(c); err != nil {
			return
		}
	}
	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, name := range b.Names {
		maxBid++
		valueStrings = append(valueStrings, "(?,?,?,?,?,?)")
		valueArgs = append(valueArgs, b.PID, maxBid, name, b.Flow, b.FlowState, b.State)
	}
	stmt := fmt.Sprintf("INSERT INTO manager_business(pid,bid,name,flow,flow_state,state) VALUES %s ON DUPLICATE KEY UPDATE pid=values(pid),name=values(name),flow=values(flow),flow_state=values(flow_state),state=values(state)", strings.Join(valueStrings, ","))
	err = d.db.Exec(stmt, valueArgs...).Error
	return
}

// maxBidByPid .
func (d *Dao) maxBid(c context.Context) (maxBid int64, err error) {
	err = d.db.Table("manager_business").Select("max(bid)").Row().Scan(&maxBid)
	return
}

// UpdateBusiness .
func (d *Dao) UpdateBusiness(c context.Context, b *model.Business) error {
	return d.db.Table("manager_business").Where("id = ?", b.ID).
		Update(map[string]interface{}{
			"name": b.Name,
			"flow": b.Flow,
		}).Error
}

// BusinessChilds .
func (d *Dao) BusinessChilds(c context.Context, pid int64) (res []*model.Business, err error) {
	err = d.db.Table("manager_business").Where("pid = ?", pid).Find(&res).Error
	return
}

// AddRole .
func (d *Dao) AddRole(c context.Context, br *model.BusinessRole) (err error) {
	stmt := fmt.Sprintf("INSERT INTO manager_business_role(bid,rid,name,type,state) VALUES %s ON DUPLICATE KEY UPDATE bid=values(bid),name=values(name),type=values(type),state=values(state)", "(?,?,?,?,?)")
	err = d.db.Exec(stmt, br.BID, br.RID, br.Name, br.Type, br.State).Error
	return
}

// MaxRidByBid .
func (d *Dao) MaxRidByBid(c context.Context, bid int64) (maxRid interface{}, err error) {
	err = d.db.Table("manager_business_role").Select("max(rid) as mrid").Where("bid = ?", bid).Row().Scan(&maxRid)
	return
}

// UpdateRole .
func (d *Dao) UpdateRole(c context.Context, br *model.BusinessRole) error {
	return d.db.Table("manager_business_role").Where("id = ?", br.ID).
		Update("name", br.Name).Error
}

// RoleByRIDs .
func (d *Dao) RoleByRIDs(c context.Context, bid int64, rids []int64) (res map[int64]*model.BusinessRole, err error) {
	t := []*model.BusinessRole{}
	res = make(map[int64]*model.BusinessRole)
	err = d.db.Where("rid IN (?)", rids).Where("bid = ?", bid).Find(&t).Error
	for _, r := range t {
		res[r.RID] = r
	}
	return
}

// AddUser add users with associated roles .
func (d *Dao) AddUser(c context.Context, bur *model.BusinessUserRole) error {
	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, uid := range bur.UIDs {
		valueStrings = append(valueStrings, "(?,?,?,?)")
		valueArgs = append(valueArgs, uid, bur.CUID, bur.BID, bur.Role)
	}
	stmt := fmt.Sprintf("INSERT INTO manager_business_user_role(uid, cuid, bid, role) VALUES %s ON DUPLICATE KEY UPDATE uid=values(uid),cuid=values(cuid),bid=values(bid),role=values(role)", strings.Join(valueStrings, ","))
	return d.db.Exec(stmt, valueArgs...).Error
}

// UpdateUser .
func (d *Dao) UpdateUser(c context.Context, bur *model.BusinessUserRole) error {
	return d.db.Table("manager_business_user_role").Where("id = ?", bur.ID).
		Update("role", bur.Role).Error
}

// UpdateBusinessState .
func (d *Dao) UpdateBusinessState(c context.Context, su *model.StateUpdate) error {
	return d.db.Table("manager_business").Where("id = ?", su.ID).
		Update("state", su.State).Error
}

// BatchUpdateChildState .
func (d *Dao) BatchUpdateChildState(c context.Context, pid int64, flowStates []int64) error {
	return d.db.Table("manager_business").Where("flow_state NOT IN (?)", flowStates).Where("pid = ?", pid).
		Update("state", 0).Error
}

// UpdateBusinessRoleState .
func (d *Dao) UpdateBusinessRoleState(c context.Context, su *model.StateUpdate) error {
	return d.db.Table("manager_business_role").Where("id = ?", su.ID).
		Update("state", su.State).Error
}

// ParentBusiness .
func (d *Dao) ParentBusiness(c context.Context, state int64) (res map[int64]*model.BusinessList, err error) {
	temp := []*model.BusinessList{}
	res = make(map[int64]*model.BusinessList)
	db := d.db.Where("pid = ?", 0)
	if state != -1 {
		db = db.Where("state = ?", state)
	}
	if err = db.Order("id asc", true).Find(&temp).Error; err != nil {
		return
	}
	for _, t := range temp {
		res[t.ID] = t
	}
	return
}

// ChildBusiness .
func (d *Dao) ChildBusiness(c context.Context, bp *model.BusinessListParams) (res map[int64]*model.BusinessList, err error) {
	temp := []*model.BusinessList{}
	res = make(map[int64]*model.BusinessList)
	db := d.db.Where("pid != ?", 0)
	if bp.State != -1 {
		db = db.Where("state = ?", bp.State)
	}
	if bp.Flow != 0 {
		db = db.Where("flow_state = ?", bp.Flow)
	}
	err = db.Order("ctime desc").Find(&temp).Error
	for _, t := range temp {
		res[t.ID] = t
	}
	return
}

// ChildBusinessByPIDs .
func (d *Dao) ChildBusinessByPIDs(c context.Context, pids []int64) (res map[int64][]*model.BusinessList, err error) {
	res = make(map[int64][]*model.BusinessList)
	temp := []*model.BusinessList{}
	if err = d.db.Where("pid IN (?)", pids).Find(&temp).Error; err != nil {
		return
	}
	for _, t := range temp {
		res[t.PID] = append(res[t.PID], t)
	}
	return
}

// FlowList .
func (d *Dao) FlowList(c context.Context, bp *model.BusinessListParams) (res []*model.BusinessList, err error) {
	db := d.db.Where("pid != ?", 0)
	if bp.Flow > 0 {
		db = db.Where("flow_state = ?", bp.Flow)
	}
	if bp.State != -1 {
		db = db.Where("state = ?", bp.State)
	}
	err = db.Find(&res).Error
	return
}

// RoleListByBID .
func (d *Dao) RoleListByBID(c context.Context, br *model.BusinessRole) (res []*model.BusinessRole, err error) {
	db := d.db
	if br.BID > 0 {
		db = db.Where("bid = ?", br.BID)
	}
	if br.Type != -1 {
		db = db.Where("type = ?", br.Type)
	}
	if br.State != -1 {
		db = db.Where("state = ?", br.State)
	}
	err = db.Order("rid desc").Find(&res).Error
	return
}

// RoleListByRIDs .
func (d *Dao) RoleListByRIDs(c context.Context, bid int64, rids []int64) (res []*model.BusinessRole, err error) {
	err = d.DB().Table("manager_business_role").Where("bid = ? AND rid IN (?)", bid, rids).Where("state = ?", 1).Find(&res).Error
	return
}

// UserList .
func (d *Dao) UserList(c context.Context, u *model.UserListParams) (res []*model.BusinessUserRoleList, err error) {
	db := d.db.Table("manager_business_user_role")
	if u.BID > 0 {
		db = db.Where("bid = ?", u.BID)
	}
	if u.UID > 0 {
		db = db.Where("uid = ?", u.UID)
	}
	if u.Role != -1 {
		db = db.Where("find_in_set(?, role) > 0", strconv.FormatInt(u.Role, 10))
	}
	err = db.Where("state = ?", model.UserOnState).Find(&res).Error
	return
}

// DeleteUser .
func (d *Dao) DeleteUser(c context.Context, bur *model.BusinessUserRole) error {
	return d.db.Table("manager_business_user_role").Where("id = ?", bur.ID).Update("state", model.UserOffState).Error
}

// BusinessByID .
func (d *Dao) BusinessByID(c context.Context, id int64) (res *model.BusinessList, err error) {
	res = &model.BusinessList{}
	if err = d.db.Where("id = ?", id).First(res).Error; err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

// UserRoles .
func (d *Dao) UserRoles(c context.Context, uid int64) (res []*model.BusinessUserRoleList, err error) {
	err = d.db.Table("manager_business_user_role").Where("uid = ?", uid).Find(&res).Error
	return
}

// UserRoleByBIDs .
func (d *Dao) UserRoleByBIDs(c context.Context, uid int64, bids []int64) (res []*model.BusinessUserRoleList, err error) {
	err = d.DB().Table("manager_business_user_role").Where("bid IN (?)", bids).Where("uid = ? AND state = ?", uid, 1).Find(&res).Error
	return
}
