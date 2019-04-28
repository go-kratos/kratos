package fav

import (
	"context"
	favmdl "go-common/app/service/main/favorite/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestFavfolderMcKey(t *testing.T) {
	convey.Convey("folderMcKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := folderMcKey(mid, fid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavrelationFidsKey(t *testing.T) {
	convey.Convey("relationFidsKey", t, func(ctx convey.C) {
		var (
			typ = int8(0)
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := relationFidsKey(typ, mid, oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavoidCountKey(t *testing.T) {
	convey.Convey("oidCountKey", t, func(ctx convey.C) {
		var (
			typ = int8(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := oidCountKey(typ, oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavbatchOidsKey(t *testing.T) {
	convey.Convey("batchOidsKey", t, func(ctx convey.C) {
		var (
			typ = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := batchOidsKey(typ, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavrecentOidsKey(t *testing.T) {
	convey.Convey("recentOidsKey", t, func(ctx convey.C) {
		var (
			typ = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := recentOidsKey(typ, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavPingMC(t *testing.T) {
	convey.Convey("PingMC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavSetFoldersMc(t *testing.T) {
	convey.Convey("SetFoldersMc", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			vs = &favmdl.Folder{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetFoldersMc(c, vs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavFolderMc(t *testing.T) {
	convey.Convey("FolderMc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			typ = int8(0)
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			f, err := d.FolderMc(c, typ, mid, fid)
			ctx.Convey("Then err should be nil.f should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(f, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavSetRelaitonFidsMc(t *testing.T) {
	convey.Convey("SetRelaitonFidsMc", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			typ  = int8(0)
			mid  = int64(0)
			oid  = int64(0)
			fids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRelaitonFidsMc(c, typ, mid, oid, fids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavRelaitonFidsMc(t *testing.T) {
	convey.Convey("RelaitonFidsMc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			typ = int8(0)
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			fids, err := d.RelaitonFidsMc(c, typ, mid, oid)
			ctx.Convey("Then err should be nil.fids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavDelRelationFidsMc(t *testing.T) {
	convey.Convey("DelRelationFidsMc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			typ = int8(0)
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelRelationFidsMc(c, typ, mid, oid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavSetOidCountMc(t *testing.T) {
	convey.Convey("SetOidCountMc", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			typ   = int8(0)
			oid   = int64(0)
			count = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetOidCountMc(c, typ, oid, count)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavDelBatchOidsMc(t *testing.T) {
	convey.Convey("DelBatchOidsMc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			typ = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelBatchOidsMc(c, typ, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavDelRecentOidsMc(t *testing.T) {
	convey.Convey("DelRecentOidsMc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			typ = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelRecentOidsMc(c, typ, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
