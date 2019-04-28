package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAll(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("All", t, func(ctx convey.C) {
		bs, err := d.All(c)
		ctx.Convey("Then err should be nil.bs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(bs, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoMaxSeq(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(2)
	)
	convey.Convey("MaxSeq", t, func(ctx convey.C) {
		maxSeq, err := d.MaxSeq(c, businessID)
		ctx.Convey("Then err should be nil.maxSeq should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(maxSeq, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpMaxSeq(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		maxSeq     = int64(10)
		lastSeq    = int64(2)
	)
	convey.Convey("UpMaxSeq", t, func(ctx convey.C) {
		rows, err := d.UpMaxSeq(c, businessID, maxSeq, lastSeq)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpMaxSeqToken(t *testing.T) {
	var (
		c          = context.TODO()
		businessID = int64(1)
		maxSeq     = int64(10)
		step       = int64(2)
		token      = ""
	)
	convey.Convey("UpMaxSeqToken", t, func(ctx convey.C) {
		rows, err := d.UpMaxSeqToken(c, businessID, maxSeq, step, token)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}
