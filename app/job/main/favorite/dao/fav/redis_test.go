package fav

import (
	"context"
	favmdl "go-common/app/service/main/favorite/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPing(t *testing.T) {
	convey.Convey("ping test", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Ping(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.pingMC(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.pingMySQL(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.pingRedis(c)
			ctx.So(err, convey.ShouldBeNil)
			t, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(t, convey.ShouldNotBeNil)
			t.Rollback()
		})
	})
}

func TestDelNewCovers(t *testing.T) {
	convey.Convey("del covers", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelNewCoverCache(c, mid, fid)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
func TestExpireAll(t *testing.T) {
	convey.Convey("ExpireAll", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ExpireAllRelations(c, mid, fid)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
func TestDelAll(t *testing.T) {
	convey.Convey("DelAll", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelAllRelationCache(c, mid, fid, 1, 2)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
func TestFavfavedBitKey(t *testing.T) {
	convey.Convey("favedBitKey", t, func(ctx convey.C) {
		var (
			tp  = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := favedBitKey(tp, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavfolderKey(t *testing.T) {
	convey.Convey("folderKey", t, func(ctx convey.C) {
		var (
			tp  = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := folderKey(tp, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavrelationKey(t *testing.T) {
	convey.Convey("relationKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := relationKey(mid, fid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavoldRelationKey(t *testing.T) {
	convey.Convey("oldRelationKey", t, func(ctx convey.C) {
		var (
			typ = int8(0)
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := oldRelationKey(typ, mid, fid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavrelationOidsKey(t *testing.T) {
	convey.Convey("relationOidsKey", t, func(ctx convey.C) {
		var (
			tp  = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := relationOidsKey(tp, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavcleanedKey(t *testing.T) {
	convey.Convey("cleanedKey", t, func(ctx convey.C) {
		var (
			tp  = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := cleanedKey(tp, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavSetUnFavedBit(t *testing.T) {
	convey.Convey("SetUnFavedBit", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUnFavedBit(c, tp, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavSetFavedBit(t *testing.T) {
	convey.Convey("SetFavedBit", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetFavedBit(c, tp, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavExpireRelations(t *testing.T) {
	convey.Convey("ExpireRelations", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.ExpireRelations(c, mid, fid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavFolderCache(t *testing.T) {
	convey.Convey("FolderCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(1)
			mid = int64(1)
			fid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.FolderCache(c, tp, mid, fid)
			ctx.Convey("Then err should be nil.folder should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavDefaultFolderCache(t *testing.T) {
	convey.Convey("DefaultFolderCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(1)
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.DefaultFolderCache(c, tp, mid)
			ctx.Convey("Then err should be nil.folder should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavfoldersCache(t *testing.T) {
	convey.Convey("foldersCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(1)
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.foldersCache(c, tp, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavRelationCntCache(t *testing.T) {
	convey.Convey("RelationCntCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cnt, err := d.RelationCntCache(c, mid, fid)
			ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cnt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavAddRelationCache(t *testing.T) {
	convey.Convey("AddRelationCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
			m = &favmdl.Favorite{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddRelationCache(c, m)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavAddRelationsCache(t *testing.T) {
	convey.Convey("AddRelationsCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
			fid = int64(0)
			fs  = []*favmdl.Favorite{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddRelationsCache(c, tp, mid, fid, fs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavDelRelationsCache(t *testing.T) {
	convey.Convey("DelRelationsCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelRelationsCache(c, mid, fid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavDelRelationCache(t *testing.T) {
	convey.Convey("DelRelationCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelRelationCache(c, mid, fid, oid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavDelOldRelationsCache(t *testing.T) {
	convey.Convey("DelOldRelationsCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			typ = int8(0)
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelOldRelationsCache(c, typ, mid, fid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavExpireRelationOids(t *testing.T) {
	convey.Convey("ExpireRelationOids", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.ExpireRelationOids(c, tp, mid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavAddRelationOidCache(t *testing.T) {
	convey.Convey("AddRelationOidCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddRelationOidCache(c, tp, mid, oid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavRemRelationOidCache(t *testing.T) {
	convey.Convey("RemRelationOidCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemRelationOidCache(c, tp, mid, oid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavSetRelationOidsCache(t *testing.T) {
	convey.Convey("SetRelationOidsCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int8(0)
			mid  = int64(0)
			oids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRelationOidsCache(c, tp, mid, oids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavSetCleanedCache(t *testing.T) {
	convey.Convey("SetCleanedCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			typ    = int8(0)
			mid    = int64(0)
			fid    = int64(0)
			ftime  = int64(0)
			expire = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetCleanedCache(c, typ, mid, fid, ftime, expire)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavDelRelationOidsCache(t *testing.T) {
	convey.Convey("DelRelationOidsCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			typ = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelRelationOidsCache(c, typ, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
