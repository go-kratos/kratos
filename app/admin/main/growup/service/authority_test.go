package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_AddAuthorityUser(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "DELETE FROM authority_user WHERE username = 'maru'")
		err := s.AddAuthorityUser("maru", "maru")
		So(err, ShouldBeNil)
	}))
}

func TestService_AddAuthorityTaskRoleUser(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.AddAuthorityTaskRoleUser("maru", "1")
		So(err, ShouldBeNil)
	}))
}

func TestService_AddAuthorityTaskGroupUser(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.AddAuthorityTaskGroupUser("maru", "1")
		So(err, ShouldBeNil)
	}))
}

func TestService_AddAuthorityTaskRole(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "DELETE FROM authority_task_role WHERE name = 'test'")
		err := s.AddAuthorityTaskRole(1, "test", "test")
		So(err, ShouldBeNil)
	}))
}

func TestService_AddAuthorityTaskGroup(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "DELETE FROM authority_task_group WHERE name = 'test'")
		err := s.AddAuthorityTaskGroup("test", "test")
		So(err, ShouldBeNil)
	}))
}

func TestService_AddPrivilege(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "DELETE FROM authority_privilege WHERE name = 'test'")
		err := s.AddPrivilege("test", 1, 1, 1)
		So(err, ShouldBeNil)
	}))
}

func TestService_DeleteAuthorityTaskRoleUser(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.DeleteAuthorityTaskRoleUser(11, 1)
		So(err, ShouldBeNil)
	}))
}

func TestService_DeleteAuthorityUser(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.DeleteAuthorityUser(11)
		So(err, ShouldBeNil)
	}))
}

func TestService_DeleteAuthorityTaskRole(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.DeleteAuthorityTaskRole(1)
		So(err, ShouldBeNil)
	}))
}

func TestService_DeleteAuthorityTaskGroupUser(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.DeleteAuthorityTaskGroupUser(11, 1)
		So(err, ShouldBeNil)
	}))
}

func TestService_DeleteAuthorityTaskGroup(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.DeleteAuthorityTaskGroup(1)
		So(err, ShouldBeNil)
	}))
}

func TestService_GetAuthorityUserPrivileges(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		_, err := s.GetAuthorityUserPrivileges("maru")
		So(err, ShouldBeNil)
	}))
}

func TestService_ListPrivilege(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		_, err := s.ListPrivilege()
		So(err, ShouldBeNil)
	}))
}

func TestService_ListGroupAndRole(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		_, _, err := s.ListGroupAndRole()
		So(err, ShouldBeNil)
	}))
}

func TestService_ListAuthorityTaskRoles(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "UPDATE authority_user SET task_group = 1 WHERE username = 'test'")
		_, _, err := s.ListAuthorityTaskRoles("test", 0, 10, "id")
		So(err, ShouldBeNil)
	}))
}

func TestService_ListAuthorityTaskGroups(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		_, _, err := s.ListAuthorityTaskGroups(0, 10, "id")
		So(err, ShouldBeNil)
	}))
}

func TestService_ListAuthorityRolePrivilege(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "UPDATE authority_task_group SET is_deleted = 0 WHERE id = 1")
		s.dao.Exec(context.TODO(), "UPDATE authority_task_role SET is_deleted = 0 WHERE id = 1")
		_, err := s.ListAuthorityRolePrivilege(1, 1, 2)
		So(err, ShouldBeNil)
	}))
}

func TestService_ListAuthorityGroupPrivilege(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "UPDATE authority_task_group SET is_deleted = 0 WHERE id = 1")
		_, err := s.ListAuthorityGroupPrivilege(1, 1)
		So(err, ShouldBeNil)
	}))
}

func TestService_ListAuthorityUsers(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		_, _, err := s.ListAuthorityUsers("", 0, 0, "id")
		So(err, ShouldBeNil)
	}))
}

func TestService_UpdatePrivilege(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "DELETE FROM authority_privilege WHERE name = 'test'")
		err := s.UpdatePrivilege(1, "test", 1, 1, 1)
		So(err, ShouldBeNil)
	}))
}

func TestService_UpdateAuthorityUserInfo(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "UPDATE authority_user SET nickname = 'abc' WHERE id = 1")
		err := s.UpdateAuthorityUserInfo(1, "test")
		So(err, ShouldBeNil)
	}))
}

func TestService_UpdateAuthorityUserAuth(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.UpdateAuthorityUserAuth(1, "1", "1")
		So(err, ShouldBeNil)
	}))
}

func TestService_UpdateAuthorityTaskRoleInfo(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.UpdateAuthorityTaskRoleInfo(1, "maru", "test")
		So(err, ShouldBeNil)
	}))
}

func TestService_UpdateAuthorityRolePrivilege(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.UpdateAuthorityRolePrivilege(1, "", "")
		So(err, ShouldBeNil)
	}))
}

func TestService_UpdateAuthorityTaskGroupInfo(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		s.dao.Exec(context.TODO(), "UPDATE authority_task_group SET desc = 'tt' WHERE id = 1")
		err := s.UpdateAuthorityTaskGroupInfo(1, "", "test")
		So(err, ShouldBeNil)
	}))
}

func TestService_UpdateAuthorityGroupPrivilege(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		err := s.UpdateAuthorityGroupPrivilege(1, "", "", 0)
		So(err, ShouldBeNil)
	}))
}
