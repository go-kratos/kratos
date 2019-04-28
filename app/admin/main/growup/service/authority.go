package service

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/growup/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// const for default groupid
const (
	// task group Admin id
	AdminID = 1
)

var (
	// task role list privilege id
	_privilegeTaskRoleList = "12,48"

	// admin
	_privilegePrivilege = "7,8,9,10,11,12,38,39,40,41,42,43,44,45,46,47,48"
)

func removeDuplicatesAndEmpty(a []string) (ret []string) {
	for i := 0; i < len(a); i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

// GetUserPri get user privilege
func (s *Service) GetUserPri(username string) (priMap map[int64]bool, err error) {
	priMap = make(map[int64]bool)

	query := fmt.Sprintf("username = '%s'", username)
	users, err := s.dao.GetAuthorityUsersInfo(query, "id, task_group, task_role")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo Error(%v)", err)
		return
	}
	if len(users) == 0 {
		return
	}
	if len(users) > 1 {
		err = errors.New("get user error")
		return
	}

	userPrivileges := []string{}
	var taskPrivileges string

	user := users[0]
	if user.TaskGroup != "" {
		groups := strings.Split(user.TaskGroup, ",")
		for _, groupID := range groups {
			id, _ := strconv.ParseInt(groupID, 10, 64)
			taskPrivileges, err = s.dao.GetAuthorityTaskGroupPrivileges(id)
			if err != nil {
				log.Error("s.dao.GetAuthorityTaskGroupPrivileges Error(%v)", err)
				return
			}
			userPrivileges = append(userPrivileges, strings.Split(taskPrivileges, ",")...)
		}
	}

	if user.TaskRole != "" {
		roles := strings.Split(user.TaskRole, ",")
		for _, roleID := range roles {
			id, _ := strconv.ParseInt(roleID, 10, 64)
			taskPrivileges, err = s.dao.GetAuthorityTaskRolePrivileges(id)
			if err != nil {
				log.Error("s.dao.GetAuthorityTaskRolePrivileges Error(%v)", err)
				return
			}
			userPrivileges = append(userPrivileges, strings.Split(taskPrivileges, ",")...)
		}
	}

	sort.Strings(userPrivileges)
	userPrivileges = removeDuplicatesAndEmpty(userPrivileges)

	for _, privilegesID := range userPrivileges {
		id, _ := strconv.ParseInt(privilegesID, 10, 64)
		priMap[id] = true
	}
	return
}

// GetAuthorityUserPrivileges get user all privileges
func (s *Service) GetAuthorityUserPrivileges(username string) (data interface{}, err error) {
	mPrivileges, err := s.GetUserPri(username)
	if err != nil {
		log.Error("s.GetUserPri Error(%v)", err)
		return
	}

	if len(mPrivileges) == 0 {
		data = map[string]interface{}{
			"router": []string{},
		}
		return
	}

	all, err := s.ListPrivilege()
	if err != nil {
		log.Error("s.ListPrivilege Error(%v)", err)
		return
	}

	router := []int64{}
	for _, level1 := range all {
		for _, level2 := range level1.Children {
			for _, level3 := range level2.Children {
				if _, ok := mPrivileges[level3.ID]; ok {
					level3.Selected = true
					level2.Selected = true
					level1.Selected = true
					if level3.IsRouter == 1 {
						router = append(router, level3.ID)
					}
				}
			}
		}
	}

	data = map[string]interface{}{
		"privileges": all,
		"router":     router,
	}
	return
}

// GetAuthorityUserGroup get authority user group
func (s *Service) GetAuthorityUserGroup(username string) (data []*model.Group, err error) {
	query := fmt.Sprintf("username = '%s'", username)
	users, err := s.dao.GetAuthorityUsersInfo(query, "id, task_group")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo Error(%v)", err)
		return
	}
	if len(users) == 0 {
		return
	}
	if len(users) > 1 {
		err = errors.New("get user error")
		return
	}
	if users[0].TaskGroup == "" {
		return
	}

	groups := strings.Split(users[0].TaskGroup, ",")
	allGroups, err := s.dao.GetAuthorityTaskGroups("")
	if err != nil {
		log.Error("s.dao.ListAuthorityTaskGroups Error(%v)", err)
		return
	}

	gMap := make(map[int64]*model.Group)
	for _, g := range allGroups {
		gMap[g.ID] = g
	}

	data = make([]*model.Group, 0)
	for _, id := range groups {
		gID, _ := strconv.ParseInt(id, 10, 64)
		if gID == 1 {
			data = allGroups
			return
		}
		data = append(data, gMap[gID])
	}

	return
}

// ListAuthorityUsers list all users in authority-manage
func (s *Service) ListAuthorityUsers(username string, from, limit int, sort string) (users []*model.User, total int, err error) {
	query := ""
	if len(username) != 0 {
		query = fmt.Sprintf("username = '%s'", username)
	}
	users, total, err = s.dao.ListAuthorityUsers(query, from, limit, sort)
	for i := 0; i < len(users); i++ {
		taskGroup := users[i].TaskGroup
		groups := []*model.Group{}
		if len(taskGroup) > 0 {
			query = fmt.Sprintf("id in (%s)", taskGroup)
			groups, err = s.dao.GetAuthorityTaskGroups(query)
			if err != nil {
				log.Error("s.dao.GetAuthorityTaskGroups Error(%v)", err)
				return
			}
		}
		users[i].Groups = groups

		taskRole := users[i].TaskRole
		roles := []*model.Role{}
		if len(taskRole) > 0 {
			query = fmt.Sprintf("id in (%s)", taskRole)
			roles, err = s.dao.GetAuthorityTaskRoles(query)
			if err != nil {
				log.Error("s.dao.GetAuthorityTaskRoles Error(%v)", err)
				return
			}
		}
		users[i].Roles = roles
	}
	return
}

// AddAuthorityUser add user to authority-manage
func (s *Service) AddAuthorityUser(username, nickname string) (err error) {
	user := model.User{
		Username: username,
		Nickname: nickname,
		ATime:    xtime.Time(time.Now().Unix()),
	}
	return s.dao.AddAuthorityUser(&user)
}

// UpdateAuthorityUserInfo update user's nickname
func (s *Service) UpdateAuthorityUserInfo(id int64, nickname string) (err error) {
	update := map[string]interface{}{
		"nickname": nickname,
	}
	return s.dao.UpdateAuthorityUser(id, update)
}

// UpdateAuthorityUserAuth update user's task group and task role
func (s *Service) UpdateAuthorityUserAuth(id int64, groupID, roleID string) (err error) {
	update := map[string]interface{}{
		"task_group": groupID,
		"task_role":  roleID,
	}
	return s.dao.UpdateAuthorityUser(id, update)
}

// DeleteAuthorityUser delete user from authority-manage from id
func (s *Service) DeleteAuthorityUser(id int64) (err error) {
	return s.dao.DeleteAuthorityUser(id)
}

func getAuthorityGroupUsers(groupID string, users []*model.User) (result []*model.SUser) {
	result = []*model.SUser{}
	for _, user := range users {
		groupIDs := strings.Split(user.TaskGroup, ",")
		for _, id := range groupIDs {
			if id == groupID {
				result = append(result, &model.SUser{
					ID:   user.ID,
					Name: user.Username,
				})
				break
			}
		}
	}
	return
}

// ListAuthorityTaskGroups list all task groups in authority-manage
func (s *Service) ListAuthorityTaskGroups(from, limit int, sort string) (groups []*model.TaskGroup, total int, err error) {
	groups, total, err = s.dao.ListAuthorityTaskGroups("", from, limit, sort)
	if err != nil {
		log.Error("s.dao.ListAuthorityTaskGroups Error(%v)", err)
		return
	}

	users, err := s.dao.GetAuthorityUsersInfo("", "id, username, task_group")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo Error(%v)", err)
		return
	}
	for i := 0; i < len(groups); i++ {
		groups[i].Users = getAuthorityGroupUsers(strconv.FormatInt(groups[i].ID, 10), users)
	}
	return
}

// AddAuthorityTaskGroup add new task group to authority-manage and add privilege 48
func (s *Service) AddAuthorityTaskGroup(name, desc string) (err error) {
	if len(name) == 0 {
		return errors.New("get group name error")
	}
	group := model.TaskGroup{
		Name:  name,
		Desc:  desc,
		ATime: xtime.Time(time.Now().Unix()),
	}
	err = s.dao.AddAuthorityTaskGroup(&group)
	if err != nil {
		log.Error("s.dao.AddAuthorityTaskGroup Error(%v)", err)
		return
	}

	newGroup, err := s.dao.GetAuthorityTaskGroup(fmt.Sprintf("name = '%s'", name))
	if err != nil {
		log.Error("s.dao.AddAuthorityTaskGroup Error(%v)", err)
		return
	}

	return s.UpdateAuthorityGroupPrivilege(newGroup.ID, "", "", 0)
}

// AddAuthorityTaskGroupUser add user to task group
func (s *Service) AddAuthorityTaskGroupUser(username, groupID string) (err error) {
	if len(username) == 0 {
		return errors.New("get username error")
	}
	var users []*model.User
	query := fmt.Sprintf("username = '%s'", username)
	users, err = s.dao.GetAuthorityUsersInfo(query, "id, task_group")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo error(%v)", err)
		return
	}
	if len(users) != 1 {
		return ecode.GrowupAuthorityUserNotFound
	}

	if len(users[0].TaskGroup) != 0 {
		groupID = users[0].TaskGroup + "," + groupID
	}

	update := map[string]interface{}{"task_group": groupID}
	err = s.dao.UpdateAuthorityUser(users[0].ID, update)
	if err != nil {
		log.Error("s.dao.UpdateAuthorityUser error(%v)", err)
	}
	return
}

// UpdateAuthorityTaskGroupInfo update task group info
func (s *Service) UpdateAuthorityTaskGroupInfo(groupID int64, name, desc string) (err error) {
	update := make(map[string]interface{})
	if len(desc) != 0 {
		update["desc"] = desc
	}
	if len(name) != 0 {
		update["name"] = name
	}
	return s.dao.UpdateAuthorityTaskGroup(groupID, update)
}

func spliceStrs(strs []string) (ret string) {
	for _, str := range strs {
		ret += str + ","
	}
	if len(ret) == 0 {
		return
	}
	return ret[:len(ret)-1]
}

func (s *Service) updateAuthorityUsersGroup(groupID string, users []*model.User) (err error) {
	for _, user := range users {
		groupIDs := strings.Split(user.TaskGroup, ",")
		for i := 0; i < len(groupIDs); i++ {
			if groupIDs[i] == groupID {
				groupIDs = append(groupIDs[:i], groupIDs[i+1:]...)
				err = s.dao.UpdateAuthorityUser(user.ID, map[string]interface{}{"task_group": spliceStrs(groupIDs)})
				if err != nil {
					log.Error("s.dao.UpdateAuthorityUser error(%v)", err)
					return
				}
				break
			}
		}
	}
	return
}

// DeleteAuthorityTaskGroup delete task group
func (s *Service) DeleteAuthorityTaskGroup(groupID int64) (err error) {
	err = s.dao.DeleteAuthorityTaskGroup(groupID)
	if err != nil {
		log.Error("s.dao.DeleteAuthorityTaskGroup error(%v)", err)
		return
	}

	// update users task group
	users, err := s.dao.GetAuthorityUsersInfo("", "id, task_group")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo Error(%v)", err)
		return
	}
	err = s.updateAuthorityUsersGroup(strconv.FormatInt(groupID, 10), users)
	if err != nil {
		log.Error("s.updateAuthorityUsersGroup error(%v)", err)
		return
	}

	// delete task role which belong this group
	query := fmt.Sprintf("group_id = %d", groupID)
	roles, err := s.dao.GetAuthorityTaskRoles(query)
	for _, role := range roles {
		err = s.DeleteAuthorityTaskRole(role.ID)
		if err != nil {
			log.Error("s.DeleteAuthorityTaskRole error(%v)", err)
			return
		}
	}
	return
}

// DeleteAuthorityTaskGroupUser delete user from task group
func (s *Service) DeleteAuthorityTaskGroupUser(id, groupID int64) (err error) {
	query := fmt.Sprintf("id = %d", id)
	users, err := s.dao.GetAuthorityUsersInfo(query, "id, task_group")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo Error(%v)", err)
		return
	}
	err = s.updateAuthorityUsersGroup(strconv.FormatInt(groupID, 10), users)
	if err != nil {
		log.Error("s.updateAuthorityUsersGroup error(%v)", err)
	}
	return
}

// ListAuthorityGroupPrivilege list task group's privileges
func (s *Service) ListAuthorityGroupPrivilege(groupID int64, fatherID int64) (ret *model.SPrivilege, err error) {
	var privilege string
	privilege, err = s.dao.GetAuthorityTaskGroupPrivileges(groupID)
	if err != nil {
		log.Error("s.dao.GetAuthorityTaskGroupPrivileges Error(%v)", err)
		return
	}
	privileges := strings.Split(privilege, ",")

	var data []*model.SPrivilege
	data, err = s.ListPrivilege()
	if err != nil {
		log.Error("s.ListPrivilege Error(%v)", err)
		return
	}
	for i := 0; i < len(data); i++ {
		if data[i].ID == fatherID {
			ret = data[i]
			for _, idStr := range privileges {
				id, _ := strconv.ParseInt(idStr, 10, 64)
				for _, level2 := range ret.Children {
					for _, level3 := range level2.Children {
						if level3.ID == id {
							level3.Selected = true
							level2.Selected = true
						}
					}
				}
			}
			ret.Selected = true
			break
		}
	}
	return
}

// UpdateAuthorityGroupPrivilege update group task privileges
func (s *Service) UpdateAuthorityGroupPrivilege(groupID int64, add, minus string, authType int) (err error) {
	// get old privilege by group id
	privilege, err := s.dao.GetAuthorityTaskGroupPrivileges(groupID)
	if err != nil {
		log.Error("s.dao.GetAuthorityTaskGroupPrivileges Error(%v)", err)
		return
	}
	newP := make(map[string]struct{})
	privilegeSli := strings.Split(privilege, ",")
	for _, p := range privilegeSli {
		if p == "" {
			continue
		}
		newP[p] = struct{}{}
	}
	// default add task role list privilege
	// 数据源权限不需要添加
	if authType == 0 {
		add += "," + _privilegeTaskRoleList
		if groupID == 1 {
			add += "," + _privilegePrivilege
		}
		add = strings.TrimPrefix(add, ",")
	}
	for _, p := range strings.Split(add, ",") {
		if p == "" {
			continue
		}
		newP[p] = struct{}{}
	}
	// minus
	for _, p := range strings.Split(minus, ",") {
		if p == "" {
			continue
		}
		delete(newP, p)
	}
	privileges := ""
	for p := range newP {
		privileges += p + ","
	}
	update := map[string]interface{}{
		"privileges": strings.TrimSuffix(privileges, ","),
	}
	return s.dao.UpdateAuthorityTaskGroup(groupID, update)
}

func getAuthorityRoleUsers(roleID string, users []*model.User) (result []*model.SUser) {
	result = []*model.SUser{}
	for _, user := range users {
		roleIDs := strings.Split(user.TaskRole, ",")
		for _, id := range roleIDs {
			if id == roleID {
				result = append(result, &model.SUser{
					ID:   user.ID,
					Name: user.Username,
				})
				break
			}
		}
	}
	return
}

// ListAuthorityTaskRoles list user's task roles
func (s *Service) ListAuthorityTaskRoles(username string, from, limit int, sort string) (roles []*model.TaskRole, total int, err error) {
	query := fmt.Sprintf("username = '%s'", username)
	users, err := s.dao.GetAuthorityUsersInfo(query, "task_group")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo Error(%v)", err)
		return
	}
	if len(users) == 0 || len(users) > 1 {
		err = ecode.GrowupAuthorityUserNotFound
		return
	}
	user := users[0]
	query = fmt.Sprintf("group_id in (%s)", user.TaskGroup)
	groupIDs := strings.Split(user.TaskGroup, ",")
	for _, groupID := range groupIDs {
		if groupID == strconv.Itoa(AdminID) { // Admin
			query = ""
			break
		}
	}

	roles, total, err = s.dao.ListAuthorityTaskRoles(query, from, limit, sort)
	if err != nil {
		log.Error("s.dao.ListAuthorityTaskRoles Error(%v)", err)
		return
	}

	users, err = s.dao.GetAuthorityUsersInfo("", "id, username, task_role")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo Error(%v)", err)
		return
	}
	for i := 0; i < len(roles); i++ {
		roles[i].Users = getAuthorityRoleUsers(strconv.FormatInt(roles[i].ID, 10), users)
		roles[i].GroupName, err = s.dao.GetAuthorityTaskGroupName(roles[i].GroupID)
		if err != nil {
			log.Error("s.dao.GetAuthorityTaskGroupName Error(%v)", err)
			return
		}
	}
	return
}

// AddAuthorityTaskRole add task role to authority-manage
func (s *Service) AddAuthorityTaskRole(groupID int64, name, desc string) (err error) {
	if len(name) == 0 {
		return errors.New("get role name error")
	}

	role := model.TaskRole{
		Name:    name,
		Desc:    desc,
		GroupID: groupID,
		ATime:   xtime.Time(time.Now().Unix()),
	}
	return s.dao.AddAuthorityTaskRole(&role)
}

// AddAuthorityTaskRoleUser add user to task group
func (s *Service) AddAuthorityTaskRoleUser(username, roleID string) (err error) {
	if len(username) == 0 {
		return errors.New("get username error")
	}
	var users []*model.User
	query := fmt.Sprintf("username = '%s'", username)
	users, err = s.dao.GetAuthorityUsersInfo(query, "id, task_role")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo error(%v)", err)
		return
	}
	if len(users) != 1 {
		return ecode.GrowupAuthorityUserNotFound
	}

	if len(users[0].TaskRole) != 0 {
		roleID = users[0].TaskRole + "," + roleID
	}

	update := map[string]interface{}{"task_role": roleID}
	err = s.dao.UpdateAuthorityUser(users[0].ID, update)
	if err != nil {
		log.Error("s.dao.UpdateAuthorityUser error(%v)", err)
	}
	return
}

// UpdateAuthorityTaskRoleInfo update task role info
func (s *Service) UpdateAuthorityTaskRoleInfo(roleID int64, name, desc string) (err error) {
	update := make(map[string]interface{})
	if len(desc) != 0 {
		update["desc"] = desc
	}
	if len(name) != 0 {
		update["name"] = name
	}
	return s.dao.UpdateAuthorityTaskRole(roleID, update)
}

func (s *Service) updateAuthorityUsersRole(roleID string, users []*model.User) (err error) {
	for _, user := range users {
		roleIDs := strings.Split(user.TaskRole, ",")
		for i := 0; i < len(roleIDs); i++ {
			if roleIDs[i] == roleID {
				roleIDs = append(roleIDs[:i], roleIDs[i+1:]...)
				err = s.dao.UpdateAuthorityUser(user.ID, map[string]interface{}{"task_role": spliceStrs(roleIDs)})
				if err != nil {
					log.Error("s.dao.UpdateAuthorityUser error(%v)", err)
					return
				}
				break
			}
		}
	}
	return
}

// DeleteAuthorityTaskRole delete task role
func (s *Service) DeleteAuthorityTaskRole(roleID int64) (err error) {
	err = s.dao.DeleteAuthorityTaskRole(roleID)
	if err != nil {
		log.Error("s.dao.DeleteAuthorityTaskRole error(%v)", err)
		return
	}

	// update users task role
	users, err := s.dao.GetAuthorityUsersInfo("", "id, task_role")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo Error(%v)", err)
		return
	}
	err = s.updateAuthorityUsersRole(strconv.FormatInt(roleID, 10), users)
	if err != nil {
		log.Error("s.updateAuthorityUsersRole error(%v)", err)
	}
	return
}

// DeleteAuthorityTaskRoleUser delete user from task role
func (s *Service) DeleteAuthorityTaskRoleUser(id, roleID int64) (err error) {
	query := fmt.Sprintf("id = %d", id)
	users, err := s.dao.GetAuthorityUsersInfo(query, "id, task_role")
	if err != nil {
		log.Error("s.dao.GetAuthorityUsersInfo Error(%v)", err)
		return
	}
	err = s.updateAuthorityUsersRole(strconv.FormatInt(roleID, 10), users)
	if err != nil {
		log.Error("s.updateAuthorityUsersRole error(%v)", err)
	}
	return
}

// ListAuthorityRolePrivilege list task role's privileges
func (s *Service) ListAuthorityRolePrivilege(groupID, roleID int64, fatherID int64) (data *model.SPrivilege, err error) {
	var privilege string
	privilege, err = s.dao.GetAuthorityTaskRolePrivileges(roleID)
	if err != nil {
		log.Error("s.dao.GetAuthorityTaskRolePrivileges Error(%v)", err)
		return
	}
	privileges := strings.Split(privilege, ",")

	data, err = s.ListAuthorityGroupPrivilege(groupID, fatherID)
	if err != nil {
		log.Error("s.ListAuthorityGroupPrivilege Error(%v)", err)
		return
	}

	level2 := data.Children
	for i := 0; i < len(level2); i++ {
		if !level2[i].Selected {
			level2 = append(level2[:i], level2[i+1:]...)
			i--
		} else {
			level3 := level2[i].Children
			for j := 0; j < len(level3); j++ {
				if !level3[j].Selected {
					level3 = append(level3[:j], level3[j+1:]...)
					j--
				} else {
					level3[j].Selected = false
				}
			}
			level2[i].Children = level3
			level2[i].Selected = false
		}
	}
	data.Children = level2
	data.Selected = false

	for _, idStr := range privileges {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		for _, level2 := range data.Children {
			for _, level3 := range level2.Children {
				if level3.ID == id {
					level3.Selected = true
					level2.Selected = true
				}
			}
		}
	}
	data.Selected = true
	return
}

// UpdateAuthorityRolePrivilege update role task privileges
func (s *Service) UpdateAuthorityRolePrivilege(roleID int64, add, minus string) (err error) {
	privilege, err := s.dao.GetAuthorityTaskRolePrivileges(roleID)
	if err != nil {
		log.Error("s.dao.GetAuthorityTaskRolePrivileges Error(%v)", err)
		return
	}
	newP := make(map[string]struct{})
	privilegeSli := strings.Split(privilege, ",")
	for _, p := range privilegeSli {
		if p == "" {
			continue
		}
		newP[p] = struct{}{}
	}
	// add
	for _, p := range strings.Split(add, ",") {
		if p == "" {
			continue
		}
		newP[p] = struct{}{}
	}
	// minus
	for _, p := range strings.Split(minus, ",") {
		if p == "" {
			continue
		}
		delete(newP, p)
	}
	privileges := ""
	for p := range newP {
		privileges += p + ","
	}
	update := map[string]interface{}{
		"privileges": strings.TrimSuffix(privileges, ","),
	}
	return s.dao.UpdateAuthorityTaskRole(roleID, update)
}

// ListGroupAndRole list all task groups and task roles to admin
func (s *Service) ListGroupAndRole() (groups []*model.Group, roles []*model.Role, err error) {
	groups, err = s.dao.GetAuthorityTaskGroups("")
	if err != nil {
		log.Error("s.dao.GetAuthorityTaskGroups Error(%v)", err)
		return
	}
	roles, err = s.dao.GetAuthorityTaskRoles("")
	if err != nil {
		log.Error("s.dao.GetAuthorityTaskRoles Error(%v)", err)
		return
	}
	return
}

// AddPrivilege all privilege
func (s *Service) AddPrivilege(name string, level, fatherID int64, isRouter uint8) (err error) {
	privilege := model.Privilege{
		Name:     name,
		Level:    level,
		FatherID: fatherID,
		IsRouter: isRouter,
	}
	return s.dao.AddPrivilege(&privilege)
}

// UpdatePrivilege update privilege info
func (s *Service) UpdatePrivilege(id int64, name string, level, fatherID int64, isRouter uint8) (err error) {
	update := map[string]interface{}{
		"name":      name,
		"level":     level,
		"father_id": fatherID,
		"is_router": isRouter,
	}
	return s.dao.UpdatePrivilege(id, update)
}

// ListPrivilege list privilege by level
func (s *Service) ListPrivilege() (data []*model.SPrivilege, err error) {
	var level1, level2, level3 []*model.SPrivilege

	query := fmt.Sprintf("level = 1")
	level1, err = s.dao.GetLevelPrivileges(query)
	if err != nil {
		log.Error("s.dao.GetLevelPrivileges Error(%v)", err)
		return
	}
	for _, p1 := range level1 {
		query = fmt.Sprintf("level = 2 AND father_id = %d", p1.ID)
		level2, err = s.dao.GetLevelPrivileges(query)
		if err != nil {
			log.Error("s.dao.GetLevelPrivileges Error(%v)", err)
			return
		}
		for _, p2 := range level2 {
			query = fmt.Sprintf("level = 3 AND father_id = %d", p2.ID)
			level3, err = s.dao.GetLevelPrivileges(query)
			if err != nil {
				log.Error("s.dao.GetLevelPrivileges Error(%v)", err)
				return
			}
			p2.Children = level3
			p2.Level = 2
		}
		p1.Children = level2
		p1.Level = 1
	}
	data = level1
	return
}
