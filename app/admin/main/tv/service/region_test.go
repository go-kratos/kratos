package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/tv/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_RegList(t *testing.T) {
	Convey("region list", t, WithService(func(s *Service) {
		var (
			err   error
			param = &model.Param{}
			res   []*model.RegList
			c     = context.Background()
		)
		res, err = s.RegList(c, param)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_AddReg(t *testing.T) {
	Convey("add region", t, WithService(func(s *Service) {
		var (
			err error
			c   = context.Background()
		)
		err = s.AddReg(c, "0", "0", "0", "1")
		So(err, ShouldBeNil)
	}))
}

func TestService_EditReg(t *testing.T) {
	Convey("edit region", t, WithService(func(s *Service) {
		var (
			err error
			c   = context.Background()
		)
		err = s.EditReg(c, "0", "0", "0", "0")
		So(err, ShouldBeNil)
	}))
}

func TestService_UpState(t *testing.T) {
	Convey("update state", t, WithService(func(s *Service) {
		var (
			err   error
			c     = context.Background()
			pids  = []int{1}
			state = "0"
		)
		err = s.UpState(c, pids, state)
		So(err, ShouldBeNil)
	}))
}

func TestService_RegSort(t *testing.T) {
	Convey("update region sort", t, WithService(func(s *Service) {
		var (
			err error
			c   = context.Background()
			ids = []int{1}
		)
		err = s.RegSort(c, ids)
		So(err, ShouldBeNil)
	}))
}
