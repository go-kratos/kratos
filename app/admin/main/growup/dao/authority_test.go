package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/growup/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_AddAuthorityTaskGroup(t *testing.T) {
	Convey("", t, func() {
		now := xtime.Time(time.Now().Unix())
		user := &model.SUser{
			Name: "admin",
		}
		sUser := []*model.SUser{user}
		group := &model.TaskGroup{
			Name:       "test-dao",
			Desc:       "test",
			Privileges: "test",
			ATime:      now,
			Users:      sUser,
			IsDeleted:  0,
		}
		Exec(context.Background(), "DELETE FROM authority_task_group WHERE name = 'test-dao'")
		err := d.AddAuthorityTaskGroup(group)
		So(err, ShouldBeNil)
		group.Desc = "test1"
		Exec(context.Background(), "UPDATE authority_task_group SET is_deleted = 1 WHERE name = 'test-dao'")
		err = d.AddAuthorityTaskGroup(group)
		So(err, ShouldBeNil)
	})
}

func TestDao_AddAuthorityTaskRole(t *testing.T) {
	Convey("", t, func() {
		now := xtime.Time(time.Now().Unix())
		user := &model.SUser{
			Name: "admin",
		}
		users := []*model.SUser{user}
		role := &model.TaskRole{
			Name:       "test-dao",
			Desc:       "test",
			GroupID:    1,
			Privileges: "test",
			ATime:      now,
			GroupName:  "test",
			Users:      users,
		}
		Exec(context.Background(), "DELETE FROM authority_task_role WHERE name = 'test-dao'")
		err := d.AddAuthorityTaskRole(role)
		So(err, ShouldBeNil)
		Exec(context.Background(), "UPDATE authority_task_role SET is_deleted = 1 WHERE name = 'test-dao'")
		role.Desc = "test1"
		err = d.AddAuthorityTaskRole(role)
		So(err, ShouldBeNil)

	})
}

func TestDao_AddAuthorityUser(t *testing.T) {
	Convey("", t, func() {
		now := xtime.Time(time.Now().Unix())
		group := &model.Group{
			Name: "TEST_GROUP",
		}
		role := &model.Role{
			Name: "test_role",
		}
		user := &model.User{
			Username:  "test-dao",
			Nickname:  "test",
			TaskGroup: "test",
			TaskRole:  "test",
			ATime:     now,
			IsDeleted: 0,
			Groups:    []*model.Group{group},
			Roles:     []*model.Role{role},
		}
		Exec(context.Background(), "DELETE FROM authority_user WHERE username = 'test-dao'")
		err := d.AddAuthorityUser(user)
		So(err, ShouldBeNil)
		Exec(context.Background(), "UPDATE authority_user SET is_deleted = 1 WHERE username = 'test-dao'")
		user.Nickname = "test1"
		err = d.AddAuthorityUser(user)
		So(err, ShouldBeNil)
	})
}

func TestDao_AddPrivilege(t *testing.T) {
	Convey("", t, func() {
		privilege := &model.Privilege{
			Name:      "test-dao",
			Level:     1,
			FatherID:  2,
			IsRouter:  1,
			IsDeleted: 0,
		}
		Exec(context.Background(), "DELETE FROM authority_privilege WHERE name = 'test-dao'")
		err := d.AddPrivilege(privilege)
		So(err, ShouldBeNil)
		Exec(context.Background(), "UPDATE authority_privilege SET is_deleted = 1 WHERE name = 'test-dao'")
		privilege.IsRouter = 0
		err = d.AddPrivilege(privilege)
		So(err, ShouldBeNil)
	})
}

func TestDao_DeleteAuthorityTaskGroup(t *testing.T) {
	Convey("", t, func() {
		Exec(context.Background(), "INSERT INTO authority_task_group(id) VALUES(9999)")
		err := d.DeleteAuthorityTaskGroup(9999)
		So(err, ShouldBeNil)
	})
}

func TestDao_DeleteAuthorityTaskRole(t *testing.T) {
	Convey("", t, func() {
		Exec(context.Background(), "INSERT INTO authority_task_role(id) VALUES(9999)")
		err := d.DeleteAuthorityTaskRole(9999)
		So(err, ShouldBeNil)
	})
}

func TestDao_DeleteAuthorityUser(t *testing.T) {
	Convey("", t, func() {
		Exec(context.Background(), "INSERT INTO authority_user(id) VALUES(9999)")
		err := d.DeleteAuthorityUser(1)
		So(err, ShouldBeNil)
	})
}

func TestDao_GetAuthorityTaskGroupName(t *testing.T) {
	Convey("", t, func() {
		Exec(context.Background(), "UPDATE authority_task_group SET is_deleted = 0 WHERE id = 1")
		_, err := d.GetAuthorityTaskGroupName(1)
		So(err, ShouldBeNil)
	})
}

func TestDao_GetAuthorityTaskGroupNames(t *testing.T) {
	Convey("", t, func() {
		_, err := d.GetAuthorityTaskGroupNames([]string{"1"})
		So(err, ShouldBeNil)
	})
}

func TestDao_GetAuthorityTaskGroup(t *testing.T) {
	Convey("", t, func() {
		_, err := d.GetAuthorityTaskGroup("id > 0")
		So(err, ShouldBeNil)
	})
}

func TestDao_GetAuthorityTaskGroupPrivileges(t *testing.T) {
	Convey("", t, func() {
		Exec(context.Background(), "insert into authority_task_group(id,name,privileges) values(1001, '12131','1,2,3,4,5,6,7')")
		_, err := d.GetAuthorityTaskGroupPrivileges(1001)
		So(err, ShouldBeNil)
	})
}

func TestDao_GetAuthorityTaskGroups(t *testing.T) {
	Convey("", t, func() {
		_, err := d.GetAuthorityTaskGroups("")
		So(err, ShouldBeNil)
	})
}

func TestDao_GetAuthorityTaskRolePrivileges(t *testing.T) {
	Convey("", t, func() {
		Exec(context.Background(), "insert into authority_task_role(id,name,group_id,privileges) values(1000, '12312',1001, '1,2,3,4,5,6,7')")
		_, err := d.GetAuthorityTaskRolePrivileges(1000)
		So(err, ShouldBeNil)
	})
}

func TestDao_GetAuthorityTaskRoles(t *testing.T) {
	Convey("", t, func() {
		_, err := d.GetAuthorityTaskRoles("")
		So(err, ShouldBeNil)
	})
}

func TestDao_GetAuthorityUsersInfo(t *testing.T) {
	Convey("", t, func() {
		_, err := d.GetAuthorityUsersInfo("", "id")
		So(err, ShouldBeNil)
	})
}

func TestDao_GetLevelPrivileges(t *testing.T) {
	Convey("", t, func() {
		_, err := d.GetLevelPrivileges("")
		So(err, ShouldBeNil)
	})
}

func TestDao_ListAuthorityTaskGroups(t *testing.T) {
	Convey("", t, func() {
		_, _, err := d.ListAuthorityTaskGroups("", 0, 0, "-id")
		So(err, ShouldBeNil)
	})
}

func TestDao_ListAuthorityTaskRoles(t *testing.T) {
	Convey("", t, func() {
		_, _, err := d.ListAuthorityTaskRoles("", 0, 0, "-id")
		So(err, ShouldBeNil)
	})
}

func TestDao_ListAuthorityUsers(t *testing.T) {
	Convey("", t, func() {
		_, _, err := d.ListAuthorityUsers("", 0, 0, "-id")
		So(err, ShouldBeNil)
	})
}

func TestDao_UpdateAuthorityTaskGroup(t *testing.T) {
	Convey("", t, func() {
		updates := map[string]interface{}{
			"desc": "test11",
		}
		Exec(context.Background(), "UPDATE authority_task_group SET desc = 'aaa' WHERE id = 1")
		err := d.UpdateAuthorityTaskGroup(1, updates)
		So(err, ShouldBeNil)
	})
}

func TestDao_UpdateAuthorityTaskRole(t *testing.T) {
	Convey("", t, func() {
		updates := map[string]interface{}{
			"desc": "test",
		}
		Exec(context.Background(), "UPDATE authority_task_role SET desc = 'aaa' WHERE id = 1")
		err := d.UpdateAuthorityTaskRole(1, updates)
		So(err, ShouldBeNil)
	})
}

func TestDao_UpdateAuthorityUser(t *testing.T) {
	Convey("", t, func() {
		updates := map[string]interface{}{
			"username": "test",
		}
		err := d.UpdateAuthorityUser(1, updates)
		So(err, ShouldBeNil)
	})
}

func TestDao_UpdatePrivilege(t *testing.T) {
	Convey("", t, func() {
		updates := map[string]interface{}{
			"name": "test",
		}
		err := d.UpdatePrivilege(1, updates)
		So(err, ShouldBeNil)
	})
}
