package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/member/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestMemberList(t *testing.T) {
	convey.Convey("MemberList", t, func() {
		results, err := s.Members(context.Background(), &model.ArgList{
			Keyword: "123",
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(results, convey.ShouldNotBeNil)
	})
}

func TestMemberProfile(t *testing.T) {
	convey.Convey("MemberProfile", t, func() {
		result, err := s.MemberProfile(context.Background(), 123)
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeNil)
	})
}

func TestExpLog(t *testing.T) {
	convey.Convey("ExpLog", t, func() {
		result, err := s.ExpLog(context.Background(), 123)
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeNil)
	})
}
