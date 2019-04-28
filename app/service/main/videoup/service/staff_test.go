package service

import (
	"context"
	"testing"

	"fmt"
	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/videoup/model/archive"
)

func TestService_AddApply(t *testing.T) {
	var (
		c     = context.TODO()
		staff = &archive.ApplyParam{ApplyAID: 17191032, ApplyStaffMID: 17515232, ApplyTitle: "作词"}
	)
	Convey("AddApply", t, WithService(func(s *Service) {
		data, err := svr.AddApply(c, staff, "添加")
		spew.Dump(data)
		fmt.Println(data)
		So(err, ShouldBeNil)
	}))
}

func TestService_AddApply1(t *testing.T) {
	var (
		c     = context.TODO()
		staff = &archive.ApplyParam{State: 1, ID: 3}
	)
	Convey("AddApply staff 同意", t, WithService(func(s *Service) {
		data, err := svr.AddApply(c, staff, "申请单")
		spew.Dump(data)
		fmt.Println(data)
		So(err, ShouldBeNil)
	}))
}

func TestService_AddApply2(t *testing.T) {
	var (
		c     = context.TODO()
		staff = &archive.ApplyParam{State: 2, ID: 3, FlagAddBlack: true, FlagRefuse: true}
	)
	Convey("AddApply staff 拒绝", t, WithService(func(s *Service) {
		data, err := svr.AddApply(c, staff, "申请单")
		spew.Dump(data)
		fmt.Println(data)
		So(err, ShouldBeNil)
	}))
}
