package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestQusByType(t *testing.T) {
	convey.Convey("qusByType", t, func() {
		p1 := qusByType(0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestExtraQidByType(t *testing.T) {
	convey.Convey("extraQidByType", t, func() {
		p1 := extraQidByType(0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestDaopingRedis(t *testing.T) {
	convey.Convey("pingRedis", t, func() {
		err := d.pingRedis(context.Background())
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoSetQids(t *testing.T) {
	convey.Convey("SetQids", t, func() {
		err := d.SetQids(context.Background(), []int64{1, 2, 3}, 0)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("QidByType", t, func() {
		ids, err := d.QidByType(context.Background(), 0, 3)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ids, convey.ShouldNotBeNil)
	})
}

func TestDaoSetExtraQids(t *testing.T) {
	convey.Convey("SetExtraQids", t, func() {
		err := d.SetExtraQids(context.Background(), []int64{1, 2, 3}, 0)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("ExtraQidByType", t, func() {
		ids, err := d.ExtraQidByType(context.Background(), 0, 2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ids, convey.ShouldNotBeNil)
	})
	convey.Convey("DelQidsCache", t, func() {
		err := d.DelQidsCache(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoDelExtraQidsCache(t *testing.T) {
	convey.Convey("DelExtraQidsCache", t, func() {
		err := d.DelExtraQidsCache(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}
