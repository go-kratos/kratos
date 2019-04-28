package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/manager/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
	ctx = context.Background()
)

func init() {
	dir, _ := filepath.Abs("../cmd/manager-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		// Reset(func() { CleanCache() })
		f(srv)
	}
}

func Test_Admins(t *testing.T) {
	Convey("admins", t, WithService(func(s *Service) {
		res, err := s.adms()
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		t.Logf("admins len(%d)", len(res))
	}))
}

func Test_Pointers(t *testing.T) {
	Convey("pointers", t, WithService(func(s *Service) {
		res, _, err := s.ptrs()
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		t.Logf("points len(%d)", len(res))
	}))
}

func Test_RoleAuths(t *testing.T) {
	Convey("role auths", t, WithService(func(s *Service) {
		res, err := s.roleAuths()
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		t.Logf("role auths len(%d)", len(res))
	}))
}

func Test_GroupAuths(t *testing.T) {
	Convey("group auths", t, WithService(func(s *Service) {
		res, err := s.groupAuths()
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		t.Logf("group auths len(%d)", len(res))
	}))
}

func TestService_Unames(t *testing.T) {
	Convey("unames check", t, WithService(func(s *Service) {
		var uids []int64
		uids = append(uids, 1, 2, 3)
		res := s.Unames(ctx, uids)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_UsersTotal(t *testing.T) {
	Convey("TestService_UsersTotal", t, WithService(func(s *Service) {
		res, err := s.UsersTotal(ctx)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
	}))
}

func TestService_Users(t *testing.T) {
	Convey("TestService_Users", t, WithService(func(s *Service) {
		res, err := s.Users(ctx, 1, 20)
		So(err, ShouldBeNil)
		So(len(res.Items), ShouldBeGreaterThan, 0)
	}))
}

func TestService_RankUsers(t *testing.T) {
	Convey("TestService_RankUsers", t, WithService(func(s *Service) {
		res, count, err := s.RankUsers(ctx, 1, 20, "zhaoshichen")
		So(err, ShouldBeNil)
		fmt.Println(res)
		fmt.Println(count)
	}))
}

func TestService_Ping(t *testing.T) {
	Convey("TestService_RankUsers", t, WithService(func(s *Service) {
		err := s.Ping(ctx)
		So(err, ShouldBeNil)
	}))
}

func TestService_Heartbeat(t *testing.T) {
	Convey("TestService_RankUsers", t, WithService(func(s *Service) {
		err := s.Heartbeat(ctx, "zhaoshichen")
		So(err, ShouldBeNil)
	}))
}

func TestService_Close(t *testing.T) {
	Convey("TestService_Close", t, WithService(func(s *Service) {
		s.Close()
	}))
}

func TestService_Wait(t *testing.T) {
	Convey("TestService_Wait", t, WithService(func(s *Service) {
		s.Wait()
	}))
}
