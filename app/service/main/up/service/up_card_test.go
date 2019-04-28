package service

import (
	"testing"

	"go-common/library/net/http/blademaster"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceListCardBase(t *testing.T) {
	convey.Convey("ListCardBase", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			var (
				c = &blademaster.Context{}
			)
			mids, err := s.ListCardBase(c)
			ctx.Convey("Then err should be nil.mids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceGetCardInfo(t *testing.T) {
	convey.Convey("GetCardInfo", t, func(ctx convey.C) {
		var (
			c   = &blademaster.Context{}
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			card, err := s.GetCardInfo(c, mid)
			ctx.Convey("Then err should be nil.card should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(card, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceListCardDetail(t *testing.T) {
	convey.Convey("ListCardDetail", t, func(ctx convey.C) {
		var (
			c      = &blademaster.Context{}
			offset = uint(0)
			size   = uint(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			cards, total, err := s.ListCardDetail(c, offset, size)
			ctx.Convey("Then err should be nil.cards,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(cards, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceGetCardInfoByMids(t *testing.T) {
	convey.Convey("GetCardInfoByMids", t, func(ctx convey.C) {
		var (
			c    = &blademaster.Context{}
			mids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			cards, err := s.GetCardInfoByMids(c, mids)
			ctx.Convey("Then err should be nil.cards should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cards, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicegetUpAccountsMap(t *testing.T) {
	convey.Convey("getUpAccountsMap", t, func(ctx convey.C) {
		var (
			c    = &blademaster.Context{}
			mids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			upAccountsMap, err := s.getUpAccountsMap(c, mids)
			ctx.Convey("Then err should be nil.upAccountsMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upAccountsMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicegetUpVideosMap(t *testing.T) {
	convey.Convey("getUpVideosMap", t, func(ctx convey.C) {
		var (
			c    = &blademaster.Context{}
			mids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			upVideosMap, err := s.getUpVideosMap(c, mids)
			ctx.Convey("Then err should be nil.upVideosMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upVideosMap, convey.ShouldNotBeNil)
			})
		})
	})
}
