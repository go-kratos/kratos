package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPolicies(t *testing.T) {
	convey.Convey("Policies", t, func(ctx convey.C) {
		res, err := d.Policies(context.Background())
		ctx.Convey("Error should be nil, res should not be empty", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaoGroupPolicies(t *testing.T) {
	convey.Convey("GroupPolicies", t, func(ctx convey.C) {
		res, err := d.GroupPolicies(context.Background())
		ctx.Convey("Error should be nil, res should not be empty", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaoGroupid(t *testing.T) {
	convey.Convey("should get group_id", t, func(ctx convey.C) {
		_, err := d.Groupid(context.Background(), 11424224)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGroupAuthZone(t *testing.T) {
	convey.Convey("GroupAuthZone", t, func(ctx convey.C) {
		res, err := d.GroupAuthZone(context.Background())
		ctx.Convey("Error should be nil, res should not be empty", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}
