package api

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var client MemberClient

func init() {
	var err error
	client, err = NewClient(nil)
	if err != nil {
		panic(err)
	}
}

func TestUserDetails(t *testing.T) {
	convey.Convey("Block User Details", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When get block user details", func(ctx convey.C) {
			details, err := client.BlockBatchDetail(c, &MemberMidsReq{Mids: []int64{2, 3}})
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(details, convey.ShouldNotBeNil)
		})
	})
}

func TestUserInfo(t *testing.T) {
	convey.Convey("Block User Infos", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When get block user infos", func(ctx convey.C) {
			infos, err := client.BlockBatchInfo(c, &MemberMidsReq{Mids: []int64{2}})
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(infos, convey.ShouldNotBeNil)
		})
	})
}
