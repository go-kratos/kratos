package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLimitUsers(t *testing.T) {
	var (
		c     = context.TODO()
		start = int32(0)
		end   = int32(0)
	)
	convey.Convey("LimitUsers", t, func(ctx convey.C) {
		res, err := d.LimitUsers(c, start, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoLimitUserCount(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("LimitUserCount", t, func(ctx convey.C) {
		count, err := d.LimitUserCount(c)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoLimitUser(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("LimitUser", t, func(ctx convey.C) {
		res, err := d.LimitUser(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoInsertLimitUser(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		name  = ""
		cname = ""
	)
	convey.Convey("InsertLimitUser", t, func(ctx convey.C) {
		id, err := d.InsertLimitUser(c, mid, name, cname)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelLimitUser(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("DelLimitUser", t, func(ctx convey.C) {
		affect, err := d.DelLimitUser(c, mid)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoResLimitByOid(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(0)
		typ = int32(0)
	)
	convey.Convey("ResLimitByOid", t, func(ctx convey.C) {
		res, err := d.ResLimitByOid(c, oid, typ)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoResLimitCount(t *testing.T) {
	var (
		c     = context.TODO()
		state = int32(0)
	)
	convey.Convey("ResLimitCount", t, func(ctx convey.C) {
		count, err := d.ResLimitCount(c, state)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpResLimitState(t *testing.T) {
	var (
		c     = context.TODO()
		oid   = int64(0)
		tp    = int32(0)
		opera = int32(0)
	)
	convey.Convey("UpResLimitState", t, func(ctx convey.C) {
		affect, err := d.UpResLimitState(c, oid, tp, opera)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoResLimitAdd(t *testing.T) {
	var (
		c         = context.TODO()
		oid       = int64(0)
		tp        = int32(0)
		operation = int32(0)
	)
	convey.Convey("ResLimitAdd", t, func(ctx convey.C) {
		id, err := d.ResLimitAdd(c, oid, tp, operation)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoResLimitByOpState(t *testing.T) {
	var (
		c     = context.TODO()
		state = int32(0)
		start = int32(0)
		end   = int32(0)
	)
	convey.Convey("ResLimitByOpState", t, func(ctx convey.C) {
		res, oids, err := d.ResLimitByOpState(c, state, start, end)
		ctx.Convey("Then err should be nil.res,oids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(oids, convey.ShouldHaveLength, 0)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}
