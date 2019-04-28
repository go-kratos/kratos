package medal

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMedalloadMedal(t *testing.T) {
	convey.Convey("loadMedal", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, p2, err := d.loadMedal(c, mid)
			ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p2, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalstoreMedal(t *testing.T) {
	convey.Convey("storeMedal", t, func(ctx convey.C) {
		var (
			mid     = int64(88889017)
			nid     = int64(1)
			nofound bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.storeMedal(mid, nid, nofound)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestMedallocalMedal(t *testing.T) {
	convey.Convey("localMedal", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.localMedal(mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalMedalActivatedCache(t *testing.T) {
	convey.Convey("MedalActivatedCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, p2, err := d.MedalActivatedCache(c, mid)
			ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p2, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalMedalsActivatedCache(t *testing.T) {
	convey.Convey("MedalsActivatedCache", t, func(ctx convey.C) {
		var (
			mids = []int64{88889017, 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, p2, err := d.MedalsActivatedCache(c, mids)
			ctx.Convey("Then err should be nil.p1,p2 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
				ctx.So(p2, convey.ShouldNotBeNil)
			})
		})
	})
}
