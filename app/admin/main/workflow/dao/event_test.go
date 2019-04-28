package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoEventsByCid(t *testing.T) {
	var c = context.TODO()
	convey.Convey("events by cid", t, func() {
		events, err := d.EventsByCid(c, 1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(events, convey.ShouldNotBeNil)
	})
}

func TestDaoEventsByIDs(t *testing.T) {
	var c = context.TODO()
	convey.Convey("events by multi event ids", t, func() {
		events, err := d.EventsByIDs(c, []int64{1})
		convey.So(err, convey.ShouldBeNil)
		convey.So(events, convey.ShouldNotBeNil)
	})
}

func TestDaoLastEventByCid(t *testing.T) {
	var c = context.TODO()
	convey.Convey("last event by eids", t, func() {
		events, err := d.EventsByIDs(c, []int64{1})
		convey.So(err, convey.ShouldBeNil)
		convey.So(events, convey.ShouldNotBeNil)
	})
}

func TestDaoBatchLastEventIDs(t *testing.T) {
	var c = context.TODO()
	convey.Convey("batch last event by multi cids", t, func() {
		events, err := d.BatchLastEventIDs(c, []int64{1})
		convey.So(err, convey.ShouldBeNil)
		convey.So(events, convey.ShouldNotBeNil)
	})
}
