package medal

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestMedalhit(t *testing.T) {
	convey.Convey("hit", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.hit(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalAddMedalOwner(t *testing.T) {
	convey.Convey("AddMedalOwner", t, func(ctx convey.C) {
		var (
			mid = time.Now().Unix()
			nid = int64(5)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddMedalOwner(c, mid, nid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalInstallMedalOwner(t *testing.T) {
	convey.Convey("InstallMedalOwner", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			nid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.InstallMedalOwner(c, mid, nid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalUninstallMedalOwner(t *testing.T) {
	convey.Convey("UninstallMedalOwner", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			nid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UninstallMedalOwner(c, mid, nid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalUninstallAllMedalOwner(t *testing.T) {
	convey.Convey("UninstallAllMedalOwner", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			nid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UninstallAllMedalOwner(c, mid, nid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalMedalInfoAll(t *testing.T) {
	convey.Convey("MedalInfoAll", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MedalInfoAll(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalMedalOwnerByMid(t *testing.T) {
	convey.Convey("MedalOwnerByMid", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MedalOwnerByMid(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalMedalInfoByNid(t *testing.T) {
	convey.Convey("MedalInfoByNid", t, func(ctx convey.C) {
		var (
			nid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MedalInfoByNid(c, nid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalActivatedOwnerByMid(t *testing.T) {
	convey.Convey("ActivatedOwnerByMid", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nid, err := d.ActivatedOwnerByMid(c, mid)
			ctx.Convey("Then err should be nil.nid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(nid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalCountOwnerBYNidMid(t *testing.T) {
	convey.Convey("CountOwnerBYNidMid", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			nid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.CountOwnerBYNidMid(c, mid, nid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalOwnerBYNidMid(t *testing.T) {
	convey.Convey("OwnerBYNidMid", t, func(ctx convey.C) {
		var (
			mid = int64(32141)
			nid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.OwnerBYNidMid(c, mid, nid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalMedalGroupAll(t *testing.T) {
	convey.Convey("MedalGroupAll", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MedalGroupAll(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
