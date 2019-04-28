package dao

import (
	"strings"
	"time"

	"go-common/app/admin/main/growup/model"
	"go-common/library/ecode"
	xtime "go-common/library/time"

	"github.com/jinzhu/gorm"
)

// const for db table name
const (
	// table user
	TableUser = "authority_user"
	// table taskgroup
	TableGroup = "authority_task_group"
	// table taskrole
	TableRole = "authority_task_role"
	// table privilege
	TablePrivilege = "authority_privilege"
)

func queryAddIsDeleted(query string) string {
	if len(query) != 0 {
		query += " AND is_deleted = 0"
	} else {
		query = "is_deleted = 0"
	}
	return query
}

// GetAuthorityUsersInfo get users info by query
func (d *Dao) GetAuthorityUsersInfo(query, members string) (users []*model.User, err error) {
	query = queryAddIsDeleted(query)
	err = d.db.Table(TableUser).Select(members).Where(query).Find(&users).Error
	return
}

// ListAuthorityUsers find all authority users from db
func (d *Dao) ListAuthorityUsers(query string, from, limit int, sort string) (users []*model.User, total int, err error) {
	query = queryAddIsDeleted(query)
	err = d.db.Table(TableUser).Where(query).Count(&total).Error
	if err != nil {
		return
	}
	if strings.HasPrefix(sort, "-") {
		sort = strings.TrimPrefix(sort, "-")
		sort = sort + " " + "desc"
	}
	err = d.db.Table(TableUser).Order(sort).Offset(from).Where(query).Limit(limit).Find(&users).Error
	return
}

// AddAuthorityUser add one user to authority-manage
func (d *Dao) AddAuthorityUser(user *model.User) (err error) {
	var u model.User
	err = d.db.Table(TableUser).Where("username = ?", user.Username).Find(&u).Error
	if err == gorm.ErrRecordNotFound {
		err = d.db.Table(TableUser).Create(user).Error
		return
	}
	if err != nil {
		return
	}
	if u.IsDeleted == 1 {
		update := map[string]interface{}{
			"nickname":   user.Nickname,
			"atime":      xtime.Time(time.Now().Unix()),
			"is_deleted": 0,
		}
		err = d.db.Table(TableUser).Where("username = ?", user.Username).Updates(update).Error
	} else {
		err = ecode.GrowupAuthorityExist
	}
	return
}

// UpdateAuthorityUser update user to db
func (d *Dao) UpdateAuthorityUser(id int64, update map[string]interface{}) (err error) {
	return d.db.Table(TableUser).Where("id = ? AND is_deleted = 0", id).Updates(update).Error
}

// DeleteAuthorityUser modify user's is_delete = 1
func (d *Dao) DeleteAuthorityUser(id int64) (err error) {
	update := map[string]interface{}{
		"task_group": "",
		"task_role":  "",
		"is_deleted": 1,
	}
	return d.db.Table(TableUser).Where("id = ? AND is_deleted = 0", id).Updates(update).Error
}

// ListAuthorityTaskGroups find all authority task groups from db
func (d *Dao) ListAuthorityTaskGroups(query string, from, limit int, sort string) (groups []*model.TaskGroup, total int, err error) {
	query = queryAddIsDeleted(query)
	err = d.db.Table(TableGroup).Where(query).Count(&total).Error
	if err != nil {
		return
	}
	if strings.HasPrefix(sort, "-") {
		sort = strings.TrimPrefix(sort, "-")
		sort = sort + " " + "desc"
	}
	err = d.db.Table(TableGroup).Order(sort).Offset(from).Where(query).Limit(limit).Find(&groups).Error
	return
}

// GetAuthorityTaskGroup get authority task group
func (d *Dao) GetAuthorityTaskGroup(query string) (group model.TaskGroup, err error) {
	query = queryAddIsDeleted(query)
	err = d.db.Table(TableGroup).Where(query).Find(&group).Error
	return
}

// AddAuthorityTaskGroup add new task group
func (d *Dao) AddAuthorityTaskGroup(taskGroup *model.TaskGroup) (err error) {
	var group model.TaskGroup
	err = d.db.Table(TableGroup).Where("name = ?", taskGroup.Name).Find(&group).Error
	if err == gorm.ErrRecordNotFound {
		err = d.db.Table(TableGroup).Create(taskGroup).Error
		return
	}
	if err != nil {
		return
	}
	if group.IsDeleted == 1 {
		update := map[string]interface{}{
			"desc":       taskGroup.Desc,
			"atime":      xtime.Time(time.Now().Unix()),
			"is_deleted": 0,
		}
		err = d.db.Table(TableGroup).Where("name = ?", taskGroup.Name).Updates(update).Error
	} else {
		return ecode.GrowupAuthorityExist
	}
	return
}

// UpdateAuthorityTaskGroup update task group info to db
func (d *Dao) UpdateAuthorityTaskGroup(id int64, update map[string]interface{}) (err error) {
	return d.db.Table(TableGroup).Where("id = ? AND is_deleted = 0", id).Updates(update).Error
}

// DeleteAuthorityTaskGroup modify task group's is_deleted = 1
func (d *Dao) DeleteAuthorityTaskGroup(id int64) (err error) {
	return d.db.Table(TableGroup).Where("id = ? AND is_deleted = 0", id).Updates(map[string]interface{}{"is_deleted": 1}).Error
}

// GetAuthorityTaskGroups get task groups id and name
func (d *Dao) GetAuthorityTaskGroups(query string) (groups []*model.Group, err error) {
	query = queryAddIsDeleted(query)
	err = d.db.Table(TableGroup).Where(query).Scan(&groups).Error
	return
}

// GetAuthorityTaskGroupName get task group name
func (d *Dao) GetAuthorityTaskGroupName(groupID int64) (name string, err error) {
	var group model.Group
	err = d.db.Table(TableGroup).Where("id = ? AND is_deleted = 0", groupID).Scan(&group).Error
	if err != nil {
		return
	}
	name = group.Name
	return
}

// GetAuthorityTaskGroupPrivileges get task group privileges
func (d *Dao) GetAuthorityTaskGroupPrivileges(groupID int64) (privileges string, err error) {
	var group model.TaskGroup
	err = d.db.Table(TableGroup).Where("id = ? AND is_deleted = 0", groupID).Scan(&group).Error
	if err != nil {
		return
	}
	privileges = group.Privileges
	return
}

// GetAuthorityTaskGroupNames get task group names from ids(not use)
func (d *Dao) GetAuthorityTaskGroupNames(groupIDs []string) (names []string, err error) {
	err = d.db.Table(TableGroup).Where("id in (?) AND is_deleted = 0", groupIDs).Pluck("name", &names).Error
	return
}

// ListAuthorityTaskRoles find all authority task roles from db
func (d *Dao) ListAuthorityTaskRoles(query string, from, limit int, sort string) (roles []*model.TaskRole, total int, err error) {
	query = queryAddIsDeleted(query)
	err = d.db.Table(TableRole).Where(query).Count(&total).Error
	if err != nil {
		return
	}
	if strings.HasPrefix(sort, "-") {
		sort = strings.TrimPrefix(sort, "-")
		sort = sort + " " + "desc"
	}
	err = d.db.Table(TableRole).Order(sort).Offset(from).Where(query).Limit(limit).Find(&roles).Error
	return
}

// AddAuthorityTaskRole add task role to db
func (d *Dao) AddAuthorityTaskRole(taskRole *model.TaskRole) (err error) {
	var role model.TaskRole
	err = d.db.Table(TableRole).Where("name = ?", taskRole.Name).Find(&role).Error
	if err == gorm.ErrRecordNotFound {
		err = d.db.Table(TableRole).Create(taskRole).Error
		return
	}
	if err != nil {
		return
	}
	if role.IsDeleted == 1 {
		update := map[string]interface{}{
			"desc":       taskRole.Desc,
			"group_id":   taskRole.GroupID,
			"atime":      xtime.Time(time.Now().Unix()),
			"is_deleted": 0,
		}
		err = d.db.Table(TableRole).Where("name = ?", taskRole.Name).Updates(update).Error
	} else {
		err = ecode.GrowupAuthorityExist
	}
	return
}

// UpdateAuthorityTaskRole update task role  to db
func (d *Dao) UpdateAuthorityTaskRole(id int64, update map[string]interface{}) (err error) {
	return d.db.Table(TableRole).Where("id = ? AND is_deleted = 0", id).Updates(update).Error
}

// DeleteAuthorityTaskRole modify task role's is_deleted = 1
func (d *Dao) DeleteAuthorityTaskRole(id int64) (err error) {
	return d.db.Table(TableRole).Where("id = ? AND is_deleted = 0", id).Updates(map[string]interface{}{"is_deleted": 1}).Error
}

// GetAuthorityTaskRoles get task roles id and name
func (d *Dao) GetAuthorityTaskRoles(query string) (roles []*model.Role, err error) {
	query = queryAddIsDeleted(query)
	err = d.db.Table(TableRole).Where(query).Scan(&roles).Error
	return
}

// GetAuthorityTaskRolePrivileges get task role's privileges
func (d *Dao) GetAuthorityTaskRolePrivileges(roleID int64) (privileges string, err error) {
	var role model.TaskRole
	err = d.db.Table(TableRole).Where("id = ? AND is_deleted = 0", roleID).Scan(&role).Error
	if err != nil {
		return
	}
	privileges = role.Privileges
	return
}

// AddPrivilege add privilege to db
func (d *Dao) AddPrivilege(privilege *model.Privilege) (err error) {
	var p model.Privilege
	err = d.db.Table(TablePrivilege).Where("name = ?", privilege.Name).Find(&p).Error
	if err == gorm.ErrRecordNotFound {
		err = d.db.Table(TablePrivilege).Create(privilege).Error
		return
	}
	if err != nil {
		return
	}
	if p.IsDeleted == 1 {
		update := map[string]interface{}{
			"name":       privilege.Name,
			"level":      privilege.Level,
			"father_id":  privilege.FatherID,
			"is_deleted": 0,
		}
		err = d.db.Table(TablePrivilege).Where("name = ?", privilege.Name).Updates(update).Error
	} else {
		err = ecode.GrowupAuthorityExist
	}
	return
}

// UpdatePrivilege update privilege to db
func (d *Dao) UpdatePrivilege(id int64, update map[string]interface{}) (err error) {
	return d.db.Table(TablePrivilege).Where("id = ? AND is_deleted = 0", id).Updates(update).Error
}

// GetLevelPrivileges get sprivileges by query
func (d *Dao) GetLevelPrivileges(query string) (p []*model.SPrivilege, err error) {
	query = queryAddIsDeleted(query)
	err = d.db.Table(TablePrivilege).Select("id, name, level, is_router").Where(query).Scan(&p).Error
	return
}
