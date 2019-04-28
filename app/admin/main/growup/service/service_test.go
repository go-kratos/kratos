package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/growup/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/growup-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestService_AddUpInfo(t *testing.T) {
	Convey("AddUpInfo", t, WithService(func(s *Service) {
		err := s.AddUp(context.Background(), 11, 2)
		So(err, ShouldBeNil)
	}))
}

func TestService_Block(t *testing.T) {
	Convey("Block", t, WithService(func(s *Service) {
		err := s.Block(context.Background(), 11)
		So(err, ShouldBeNil)
	}))
}

func TestService_DeleteFromBlocked(t *testing.T) {
	Convey("DeleteFromBlocked", t, WithService(func(s *Service) {
		err := s.DeleteFromBlocked(context.Background(), 11)
		So(err, ShouldBeNil)
	}))
}

func TestService_DeleteUp(t *testing.T) {
	Convey("DeleteUp", t, WithService(func(s *Service) {
		err := s.DeleteUp(context.Background(), 11)
		So(err, ShouldBeNil)
	}))
}

func TestService_Pass(t *testing.T) {
	Convey("Pass", t, WithService(func(s *Service) {
		err := s.Pass(context.Background(), []int64{11}, 0)
		So(err, ShouldBeNil)
	}))
}

func TestService_QueryFromBlocked(t *testing.T) {
	Convey("QueryFromBlocked", t, WithService(func(s *Service) {
		_, _, err := s.QueryFromBlocked(context.Background(), 0, 0, "", 0, 0, 0, 0, "mid")
		So(err, ShouldBeNil)
	}))
}

func TestService_QueryFromUpInfo(t *testing.T) {
	Convey("", t, WithService(func(s *Service) {
		_, _, err := s.QueryFromUpInfo(context.Background(), 0, 1, nil, 0, 1, 0, "", 0, 0, 1, 1, "-mid")
		So(err, ShouldBeNil)
	}))
}

func TestService_Recovery(t *testing.T) {
	Convey("Recovery", t, WithService(func(s *Service) {
		err := s.Recovery(context.Background(), 11)
		So(err, ShouldBeNil)
	}))
}

func TestService_Reject(t *testing.T) {
	Convey("Reject", t, WithService(func(s *Service) {
		mids := []int64{11}
		err := s.Reject(context.Background(), 0, mids, "违规", 1)
		So(err, ShouldBeNil)
	}))
}

func TestService_DelUpAccount(t *testing.T) {
	Convey("DelUpAccount", t, WithService(func(s *Service) {
		err := s.DelUpAccount(context.Background(), 11)
		So(err, ShouldBeNil)
	}))
}
