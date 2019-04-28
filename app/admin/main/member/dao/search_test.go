package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/member/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSearchMember(t *testing.T) {
	convey.Convey("SearchMember", t, func() {
		result, err := d.SearchMember(context.Background(), &model.ArgList{Mid: 627272, PN: 1, PS: 10})
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeNil)
	})
}

func TestSearchLog(t *testing.T) {
	convey.Convey("SearchLog", t, func() {
		result, err := d.SearchLog(context.Background(), 1164, 0, "", "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeNil)
	})
}

func TestDao_SearchFaceCheckRes(t *testing.T) {
	convey.Convey("SearchFaceCheckRes", t, func() {
		result, err := d.SearchFaceCheckRes(context.Background(), "2879cd5fb8518f7c6da75887994c1b2a7fe670bd.png")
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeNil)
	})
}

func TestDao_SearchUserAuditLog(t *testing.T) {
	convey.Convey("SearchUserAuditLog", t, func() {
		result, err := d.SearchUserAuditLog(context.Background(), 1164)
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeNil)
	})
}

func TestDao_SearchUserPropertyReview(t *testing.T) {
	convey.Convey("SearchFaceCheckRes", t, func() {
		etime := time.Now().Format("2006-01-02 15:04:05")
		stime := time.Now().AddDate(0, 0, -26).Format("2006-01-02 15:04:05")
		result, err := d.SearchUserPropertyReview(context.Background(), 0, []int{1}, []int{1, 2}, false, false, "", stime, etime, 1, 10)
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeNil)
	})
}
