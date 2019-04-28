package v2

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV2BuildsByAppID(t *testing.T) {
	var (
		appID = int64(24)
	)
	convey.Convey("BuildsByAppID", t, func(ctx convey.C) {
		builds, err := d.BuildsByAppID(appID)
		ctx.Convey("Then err should be nil.builds should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(builds, convey.ShouldNotBeNil)
		})
	})
}

func TestV2BuildsByAppIDs(t *testing.T) {
	var (
		appIDs = []int64{24}
	)
	convey.Convey("BuildsByAppIDs", t, func(ctx convey.C) {
		builds, err := d.BuildsByAppIDs(appIDs)
		ctx.Convey("Then err should be nil.builds should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(builds, convey.ShouldNotBeNil)
		})
	})
}

func TestV2TagID(t *testing.T) {
	var (
		appID = int64(24)
		build = "docker-1"
	)
	convey.Convey("TagID", t, func(ctx convey.C) {
		tagID, err := d.TagID(appID, build)
		ctx.Convey("Then err should be nil.tagID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tagID, convey.ShouldNotBeNil)
		})
	})
}
