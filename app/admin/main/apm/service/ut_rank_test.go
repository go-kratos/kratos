package service

import (
	"context"

	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_RankTen(t *testing.T) {
	convey.Convey("RankTen", t, func() {
		_, err := svr.RankTen(context.Background(), "desc")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestService_UserRank(t *testing.T) {
	username := "hedan"
	convey.Convey("UserRank", t, func() {
		_, err := svr.UserRank(context.Background(), username)
		convey.So(err, convey.ShouldBeNil)
	})
}
