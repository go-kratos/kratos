package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRacks(t *testing.T) {
	convey.Convey("Racks", t, func(ctx convey.C) {
		var (
			cluster = "uat"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			racks, err := d.Racks(cluster)
			ctx.Convey("Then err should be nil.racks should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(racks, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoVolumes(t *testing.T) {
	convey.Convey("Volumes", t, func(ctx convey.C) {
		var (
			cluster = "uat"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			volumes, err := d.Volumes(cluster)
			ctx.Convey("Then err should be nil.volumes should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(volumes, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGroups(t *testing.T) {
	convey.Convey("Groups", t, func(ctx convey.C) {
		var (
			cluster = "uat"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			groups, err := d.Groups(cluster)
			ctx.Convey("Then err should be nil.groups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(groups, convey.ShouldNotBeNil)
			})
		})
	})
}
