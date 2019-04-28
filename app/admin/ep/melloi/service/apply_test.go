package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/ep/melloi/model"
)

var (
	apply = model.Apply{ID: 1, Path: "bilibili.test.ep.melloi", From: "hujianping",
		To: "hujianping", Status: 1, Active: 1, StartTime: "1541007000", EndTime: "1541235600"}

	qar = model.QueryApplyRequest{
		Apply:      apply,
		Pagination: model.Pagination{PageNum: 1, PageSize: 1, TotalSize: 1},
	}
	userName = "hujianping"
	cookie   = "_AJSESSIONID=e2df43ed324d20811e8d1be1a9fb36d5"
)

func Test_Apply(t *testing.T) {
	Convey("query apply info", t, func() {
		_, err := s.QueryApply(&qar)
		So(err, ShouldBeNil)
	})

	Convey("query user applyList", t, func() {
		_, err := s.QueryUserApplyList(userName)
		So(err, ShouldBeNil)
	})

	Convey("add apply", t, func() {
		err := s.AddApply(c, cookie, &apply)
		So(err, ShouldBeNil)
	})

	Convey("update apply", t, func() {
		err := s.UpdateApply(cookie, &apply)
		So(err, ShouldBeNil)
	})

	Convey("delete apply", t, func() {
		err := s.DeleteApply(apply.ID)
		So(err, ShouldBeNil)
	})
}
