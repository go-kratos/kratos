package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/passport-auth/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddToken(t *testing.T) {
	convey.Convey("AddToken", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			no    = &model.Token{}
			token = []byte("9df38fe4b94a47baad001ad823b84110")
			ct    = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.AddToken(c, no, token, ct)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelToken(t *testing.T) {
	convey.Convey("DelToken", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			token = []byte("9df38fe4b94a47baad001ad823b84110")
			ct    = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.DelToken(c, token, ct)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddTokenDeleted(t *testing.T) {
	convey.Convey("AddTokenDeleted", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			no    = &model.Token{}
			token = []byte("9df38fe4b94a47baad001ad823b84110")
			ct    = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affected, err := d.AddTokenDeleted(c, no, token, ct)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoformatSuffix(t *testing.T) {
	convey.Convey("formatSuffix", t, func(ctx convey.C) {
		var (
			no = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := formatSuffix(no)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
