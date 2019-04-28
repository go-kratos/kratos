package assist

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAssistAddAssist(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
	)
	convey.Convey("AddAssist", t, func(ctx convey.C) {
		id, err := d.AddAssist(c, mid, assistMid)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistDelAssist(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
	)
	convey.Convey("DelAssist", t, func(ctx convey.C) {
		row, err := d.DelAssist(c, mid, assistMid)
		ctx.Convey("Then err should be nil.row should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(row, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistAssist(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(27515409)
		assistMid = int64(27515235)
		err       error
	)
	convey.Convey("Assist", t, func(ctx convey.C) {
		_, err = d.Assist(c, mid, assistMid)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistAssists(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("Assists", t, func(ctx convey.C) {
		as, err := d.Assists(c, mid)
		ctx.Convey("Then err should be nil.as should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(as, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistAssistCnt(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("AssistCnt", t, func(ctx convey.C) {
		count, err := d.AssistCnt(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistUps(t *testing.T) {
	var (
		c         = context.Background()
		assistMid = int64(0)
		pn        = int64(0)
		ps        = int64(0)
	)
	convey.Convey("Ups", t, func(ctx convey.C) {
		mids, ups, total, err := d.Ups(c, assistMid, pn, ps)
		ctx.Convey("Then err should be nil.mids,ups,total should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(total, convey.ShouldNotBeNil)
			ctx.So(ups, convey.ShouldNotBeNil)
			ctx.So(mids, convey.ShouldNotBeNil)
		})
	})
}
