package like

import (
	"context"

	"testing"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeMatch(t *testing.T) {
	convey.Convey("Match", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Match(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeActMatch(t *testing.T) {
	convey.Convey("ActMatch", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActMatch(c, sid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeObject(t *testing.T) {
	convey.Convey("Object", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Object(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeRawMatchSubjects(t *testing.T) {
	convey.Convey("RawMatchSubjects", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{10256}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawMatchSubjects(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeObjectsUnStart(t *testing.T) {
	convey.Convey("ObjectsUnStart", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ObjectsUnStart(c, sid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeAddGuess(t *testing.T) {
	convey.Convey("AddGuess", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(77)
			matID  = int64(7)
			objID  = int64(7)
			sid    = int64(10256)
			result = int64(7)
			stake  = int64(7)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			lastID, err := d.AddGuess(c, mid, matID, objID, sid, result, stake)
			ctx.Convey("Then err should be nil.lastID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(lastID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeListGuess(t *testing.T) {
	convey.Convey("ListGuess", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ListGuess(c, sid, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
