package pendant

import (
	"go-common/app/service/main/usersuit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPendantloadEquip(t *testing.T) {
	convey.Convey("loadEquip", t, func(ctx convey.C) {
		var (
			mid  = int64(650454)
			info = &model.PendantEquip{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.AddEquipCache(c, mid, info)
			p1, err := d.loadEquip(c, mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantstoreEquip(t *testing.T) {
	convey.Convey("storeEquip", t, func(ctx convey.C) {
		var (
			mid   = int64(650454)
			equip = &model.PendantEquip{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.storeEquip(mid, equip)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestPendantlocalEquip(t *testing.T) {
	convey.Convey("localEquip", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.localEquip(mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantEquipCache(t *testing.T) {
	convey.Convey("EquipCache", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.EquipCache(c, mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantEquipsCache(t *testing.T) {
	convey.Convey("EquipsCache", t, func(ctx convey.C) {
		var (
			mids = []int64{650454, 1}
			info = &model.PendantEquip{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.AddEquipCache(c, int64(650454), info)
			d.AddEquipCache(c, int64(1), info)
			p1, p2, err := d.EquipsCache(c, mids)
			ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p2, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
