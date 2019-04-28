package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/web/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyCard(t *testing.T) {
	convey.Convey("keyCard", t, func(ctx convey.C) {
		var (
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyCard(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetCardBakCache(t *testing.T) {
	convey.Convey("SetCardBakCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
			rs  = &model.Card{Card: &model.AccountCard{Mid: "2222"}, Space: &model.Space{}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetCardBakCache(c, mid, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCardBakCache(t *testing.T) {
	convey.Convey("CardBakCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rs, err := d.CardBakCache(c, mid)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", rs)
			})
		})
	})
}

func TestDaocommonSetBakCache(t *testing.T) {
	convey.Convey("commonSetBakCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
			bs  = []byte("")
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.commonSetBakCache(c, key, bs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
