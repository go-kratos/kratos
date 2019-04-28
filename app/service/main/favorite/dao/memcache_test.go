package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/favorite/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestGETANDSetUserRecentResourcesMc(t *testing.T) {
	var (
		mid = int64(76)
		typ = int8(2)
		c   = context.Background()
	)
	convey.Convey("folderMcKey", t, func(ctx convey.C) {
		recents, err := d.UserRecentResourcesMc(c, typ, mid)
		ctx.So(err, convey.ShouldBeNil)
		err = d.SetUserRecentResourcesMc(c, typ, mid, recents)
		ctx.So(err, convey.ShouldBeNil)
		err = d.DelRecentResMc(c, typ, mid)
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestFavfolderMcKey(t *testing.T) {
	var (
		mid = int64(0)
		fid = int64(0)
	)
	convey.Convey("folderMcKey", t, func(ctx convey.C) {
		p1 := folderMcKey(mid, fid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavfsortMcKey(t *testing.T) {
	var (
		typ = int8(0)
		mid = int64(0)
	)
	convey.Convey("fsortMcKey", t, func(ctx convey.C) {
		p1 := fsortMcKey(typ, mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavrelationFidsKey(t *testing.T) {
	var (
		typ = int8(0)
		mid = int64(0)
		oid = int64(0)
	)
	convey.Convey("relationFidsKey", t, func(ctx convey.C) {
		p1 := relationFidsKey(typ, mid, oid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavoidCountKey(t *testing.T) {
	var (
		typ = int8(0)
		oid = int64(0)
	)
	convey.Convey("oidCountKey", t, func(ctx convey.C) {
		p1 := oidCountKey(typ, oid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavbatchOidsKey(t *testing.T) {
	var (
		typ = int8(0)
		mid = int64(0)
	)
	convey.Convey("batchOidsKey", t, func(ctx convey.C) {
		p1 := batchOidsKey(typ, mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavpingMC(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("pingMC", t, func(ctx convey.C) {
		err := d.pingMC(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavSetFoldersMc(t *testing.T) {
	var (
		c  = context.TODO()
		vs = &model.Folder{}
	)
	convey.Convey("SetFoldersMc", t, func(ctx convey.C) {
		err := d.SetFoldersMc(c, vs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavFoldersMc(t *testing.T) {
	var (
		c      = context.TODO()
		fvmids = []*model.ArgFVmid{
			&model.ArgFVmid{
				Fid:  1,
				Vmid: 88888894,
			},
		}
	)
	convey.Convey("FoldersMc", t, func(ctx convey.C) {
		fs, missFvmids, err := d.FoldersMc(c, fvmids)
		ctx.Convey("Then err should be nil.fs,missFvmids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missFvmids, convey.ShouldNotBeNil)
			ctx.So(fs, convey.ShouldNotBeNil)
		})
	})
}

func TestFavFolderMc(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		mid = int64(0)
		fid = int64(0)
	)
	convey.Convey("FolderMc", t, func(ctx convey.C) {
		f, err := d.FolderMc(c, typ, mid, fid)
		ctx.Convey("Then err should be nil.f should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(f, convey.ShouldNotBeNil)
		})
	})
}

func TestFavDelFolderMc(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		mid = int64(0)
		fid = int64(0)
	)
	convey.Convey("DelFolderMc", t, func(ctx convey.C) {
		err := d.DelFolderMc(c, typ, mid, fid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavSetFolderSortMc(t *testing.T) {
	var (
		c   = context.TODO()
		fst = &model.FolderSort{}
	)
	convey.Convey("SetFolderSortMc", t, func(ctx convey.C) {
		err := d.SetFolderSortMc(c, fst)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavFolderSortMc(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		mid = int64(0)
	)
	convey.Convey("FolderSortMc", t, func(ctx convey.C) {
		fst, err := d.FolderSortMc(c, typ, mid)
		ctx.Convey("Then err should be nil.fst should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fst, convey.ShouldNotBeNil)
		})
	})
}

func TestFavDelFolderSortMc(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		mid = int64(0)
	)
	convey.Convey("DelFolderSortMc", t, func(ctx convey.C) {
		err := d.DelFolderSortMc(c, typ, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavSetRelaitonFidsMc(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(1)
		mid  = int64(88888894)
		oid  = int64(1)
		fids = []int64{1}
	)
	convey.Convey("SetRelaitonFidsMc", t, func(ctx convey.C) {
		err := d.SetRelaitonFidsMc(c, typ, mid, oid, fids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavRelaitonFidsMc(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(1)
		mid = int64(88888894)
		oid = int64(1)
	)
	convey.Convey("RelaitonFidsMc", t, func(ctx convey.C) {
		fids, err := d.RelaitonFidsMc(c, typ, mid, oid)
		ctx.Convey("Then err should be nil.fids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fids, convey.ShouldNotBeNil)
		})
	})
}

func TestFavDelRelationFidsMc(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		mid  = int64(0)
		oids = int64(0)
	)
	convey.Convey("DelRelationFidsMc", t, func(ctx convey.C) {
		err := d.DelRelationFidsMc(c, typ, mid, oids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavSetOidCountMc(t *testing.T) {
	var (
		c     = context.TODO()
		typ   = int8(0)
		oid   = int64(0)
		count = int64(0)
	)
	convey.Convey("SetOidCountMc", t, func(ctx convey.C) {
		err := d.SetOidCountMc(c, typ, oid, count)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavOidCountMc(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		oid = int64(0)
	)
	convey.Convey("OidCountMc", t, func(ctx convey.C) {
		count, err := d.OidCountMc(c, typ, oid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestFavOidsCountMc(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(1)
		oids = []int64{1, 2, 3}
	)
	convey.Convey("OidsCountMc", t, func(ctx convey.C) {
		counts, misOids, err := d.OidsCountMc(c, typ, oids)
		ctx.Convey("Then err should be nil.counts,misOids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(misOids, convey.ShouldNotBeNil)
			ctx.So(counts, convey.ShouldNotBeNil)
		})
	})
}

func TestFavSetOidsCountMc(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		cnts map[int64]int64
	)
	convey.Convey("SetOidsCountMc", t, func(ctx convey.C) {
		err := d.SetOidsCountMc(c, typ, cnts)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavBatchOidsMc(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("BatchOidsMc", t, func(ctx convey.C) {
		_, err := d.BatchOidsMc(c, typ, mid)
		ctx.Convey("Then err should be nil.oids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavSetBatchOidsMc(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(1)
		mid  = int64(88888894)
		oids = []int64{1, 2, 3}
	)
	convey.Convey("SetBatchOidsMc", t, func(ctx convey.C) {
		err := d.SetBatchOidsMc(c, typ, mid, oids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
