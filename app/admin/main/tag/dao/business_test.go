package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInBusiness(t *testing.T) {
	var (
		c      = context.TODO()
		tp     = int32(3)
		name   = "稿件"
		appkey = "3c4e41f926e51656"
		remark = "稿件"
		alias  = "archive"
	)
	convey.Convey("InBusiness", t, func(ctx convey.C) {
		id, err := d.InBusiness(c, tp, name, appkey, remark, alias)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoUpBusiness(t *testing.T) {
	var (
		c      = context.TODO()
		tp     = int32(0)
		name   = "23"
		appkey = "3c4e41f926e51656"
		remark = "23"
		alias  = "23"
	)
	convey.Convey("UpBusiness", t, func(ctx convey.C) {
		id, err := d.UpBusiness(c, name, appkey, remark, alias, tp)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoUpBusinessState(t *testing.T) {
	var (
		c     = context.TODO()
		state = int32(1)
		tp    = int32(0)
	)
	convey.Convey("UpBusinessState", t, func(ctx convey.C) {
		id, err := d.UpBusinessState(c, state, tp)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBusiness(t *testing.T) {
	var (
		c  = context.TODO()
		tp = int32(1)
	)
	convey.Convey("Business", t, func(ctx convey.C) {
		_, err := d.Business(c, tp)
		ctx.Convey("Then err should be nil.business should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoListBusiness(t *testing.T) {
	var (
		c     = context.TODO()
		state = int32(0)
	)
	convey.Convey("ListBusiness", t, func(ctx convey.C) {
		business, err := d.ListBusiness(c, state)
		ctx.Convey("Then err should be nil.business should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(business, convey.ShouldNotBeNil)
		})
	})
}
