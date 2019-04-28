package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoaddVolumeURI(t *testing.T) {
	convey.Convey("addVolumeURI", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.addVolumeURI()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoaddFreeVolumeURI(t *testing.T) {
	convey.Convey("addFreeVolumeURI", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.addFreeVolumeURI()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocompactURI(t *testing.T) {
	convey.Convey("compactURI", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.compactURI()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogroupStatusURI(t *testing.T) {
	convey.Convey("groupStatusURI", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.groupStatusURI()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

// func TestDaoAddVolume(t *testing.T) {
// 	convey.Convey("AddVolume", t, func(ctx convey.C) {
// 		var (
// 			c     = context.Background()
// 			group = "1"
// 			num   = int64(1)
// 		)
// 		ctx.Convey("add one volume to group 1", func(ctx convey.C) {
// 			err := d.AddVolume(c, group, num)
// 			if strings.Contains(err.Error(), "store response status code:7001") {
// 				t.Log("store have no free volume")
// 				err = nil // NOTE ignore error
// 			}
// 			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 			})
// 		})
// 	})
// }

// func TestDaoAddFreeVolume(t *testing.T) {
// 	convey.Convey("AddFreeVolume", t, func(ctx convey.C) {
// 		var (
// 			c     = context.Background()
// 			group = "1"
// 			dir   = "/mnt/storage00/bfsdata"
// 			num   = int64(1)
// 		)
// 		ctx.Convey("When everything goes positive", func(ctx convey.C) {
// 			err := d.AddFreeVolume(c, group, dir, num)
// 			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 			})
// 		})
// 	})
// }

func TestDaoCompact(t *testing.T) {
	convey.Convey("Compact", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			group = "1"
			vid   = int64(1)
		)
		ctx.Convey("compact group:1 vid:1", func(ctx convey.C) {
			err := d.Compact(c, group, vid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetGroupStatus(t *testing.T) {
	convey.Convey("SetGroupStatus", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			group  = "1"
			status = "health"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetGroupStatus(c, group, status)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
