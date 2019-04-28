package recommend

import (
	"context"
	"testing"

	"go-common/app/interface/main/app-card/model/card/ai"

	"github.com/smartystreets/goconvey/convey"
)

func TestAddRcmdAidsCache(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{1}
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.AddRcmdAidsCache(c, aids)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestRcmdAidsCache(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		_, err := d.RcmdAidsCache(c)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestAddRcmdCache(t *testing.T) {
	var (
		c  = context.TODO()
		is = []*ai.Item{}
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.AddRcmdCache(c, is)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestRcmdCache(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		_, err := d.RcmdCache(c)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
