package data

import (
	"context"
	"go-common/app/admin/main/up/model/datamodel"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDataGetUpArchiveInfo(t *testing.T) {
	convey.Convey("GetUpArchiveInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &datamodel.GetUpArchiveInfoArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.GetUpArchiveInfo(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataGetUpArchiveTagInfo(t *testing.T) {
	convey.Convey("GetUpArchiveTagInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &datamodel.GetUpArchiveTagInfoArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.GetUpArchiveTagInfo(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataGetUpArchiveTypeInfo(t *testing.T) {
	convey.Convey("GetUpArchiveTypeInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &datamodel.GetUpArchiveTypeInfoArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.GetUpArchiveTypeInfo(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
