package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/favorite/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestMultiExpireRelations(t *testing.T) {
	var (
		c    = context.TODO()
		fids = []int64{1, 2}
		mid  = int64(0)
	)
	convey.Convey("RemFidsRedis", t, func(ctx convey.C) {
		_, err := d.MultiExpireRelations(c, mid, fids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavRemFidsRedis(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		mid = int64(0)
		fs  = &model.Folder{}
	)
	convey.Convey("RemFidsRedis", t, func(ctx convey.C) {
		err := d.RemFidsRedis(c, typ, mid, fs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavfavedBitKey(t *testing.T) {
	var (
		tp  = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("favedBitKey", t, func(ctx convey.C) {
		p1 := favedBitKey(tp, mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavAddRelationCache(t *testing.T) {
	var (
		c = context.TODO()
		m = &model.Favorite{
			Oid:  1,
			Fid:  1,
			Mid:  88888894,
			Type: 1,
		}
	)
	convey.Convey("AddRelationCache", t, func(ctx convey.C) {
		err := d.AddRelationCache(c, m)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavfolderKey(t *testing.T) {
	var (
		tp  = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("folderKey", t, func(ctx convey.C) {
		p1 := folderKey(tp, mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavrelationKey(t *testing.T) {
	var (
		mid = int64(88888894)
		fid = int64(1)
	)
	convey.Convey("relationKey", t, func(ctx convey.C) {
		p1 := relationKey(mid, fid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavrelationOidsKey(t *testing.T) {
	var (
		tp  = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("relationOidsKey", t, func(ctx convey.C) {
		p1 := relationOidsKey(tp, mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavpingRedis(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		err := d.pingRedis(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavExpireRelations(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(88888894)
		fid = int64(1)
	)
	convey.Convey("ExpireRelations", t, func(ctx convey.C) {
		ok, err := d.ExpireRelations(c, mid, fid)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestFavExpireFolder(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("ExpireFolder", t, func(ctx convey.C) {
		ok, err := d.ExpireFolder(c, tp, mid)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestFavAddFidsRedis(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		mid = int64(0)
		fs  = &model.Folder{}
	)
	convey.Convey("AddFidsRedis", t, func(ctx convey.C) {
		err := d.AddFidsRedis(c, typ, mid, fs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavFidsRedis(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(0)
		mid = int64(0)
	)
	convey.Convey("FidsRedis", t, func(ctx convey.C) {
		_, err := d.FidsRedis(c, tp, mid)
		ctx.Convey("Then err should be nil.fids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavDelFidsRedis(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("DelFidsRedis", t, func(ctx convey.C) {
		err := d.DelFidsRedis(c, typ, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavAddFoldersCache(t *testing.T) {
	var (
		c       = context.TODO()
		tp      = int8(1)
		mid     = int64(88888894)
		folders = []*model.Folder{
			&model.Folder{
				ID: 1,
			},
		}
	)
	convey.Convey("AddFoldersCache", t, func(ctx convey.C) {
		err := d.AddFoldersCache(c, tp, mid, folders)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavFolderRelationsCache(t *testing.T) {
	var (
		c     = context.TODO()
		typ   = int8(1)
		mid   = int64(88888894)
		fid   = int64(1)
		start = int(1)
		end   = int(2)
	)
	convey.Convey("FolderRelationsCache", t, func(ctx convey.C) {
		_, err := d.FolderRelationsCache(c, typ, mid, fid, start, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavAllFolderRelationsCache(t *testing.T) {
	var (
		c     = context.TODO()
		typ   = int8(1)
		mid   = int64(88888894)
		fid   = int64(1)
		start = int(1)
		end   = int(2)
	)
	convey.Convey("FolderAllRelationsCache", t, func(ctx convey.C) {
		_, err := d.FolderAllRelationsCache(c, typ, mid, fid, start, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavCntRelationsCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(88888894)
		fid = int64(1)
	)
	convey.Convey("CntRelationsCache", t, func(ctx convey.C) {
		cnt, err := d.CntRelationsCache(c, mid, fid)
		ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(cnt, convey.ShouldNotBeNil)
		})
	})
}

func TestFavCntAllRelationsCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(88888894)
		fid = int64(1)
	)
	convey.Convey("CntAllRelationsCache", t, func(ctx convey.C) {
		cnt, err := d.CntAllRelationsCache(c, mid, fid)
		ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(cnt, convey.ShouldNotBeNil)
		})
	})
}

func TestFavExpireRelationOids(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("ExpireRelationOids", t, func(ctx convey.C) {
		ok, err := d.ExpireRelationOids(c, tp, mid)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestFavAddRelationOidCache(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
		oid = int64(1)
	)
	convey.Convey("AddRelationOidCache", t, func(ctx convey.C) {
		err := d.AddRelationOidCache(c, tp, mid, oid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavRemRelationOidCache(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
		oid = int64(123)
	)
	convey.Convey("RemRelationOidCache", t, func(ctx convey.C) {
		err := d.RemRelationOidCache(c, tp, mid, oid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavIsFavedCache(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
		oid = int64(123)
	)
	convey.Convey("IsFavedCache", t, func(ctx convey.C) {
		isFaved, err := d.IsFavedCache(c, tp, mid, oid)
		ctx.Convey("Then err should be nil.isFaved should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(isFaved, convey.ShouldNotBeNil)
		})
	})
}

func TestFavIsFavedsCache(t *testing.T) {
	var (
		c    = context.TODO()
		tp   = int8(1)
		mid  = int64(88888894)
		oids = []int64{1, 2, 3}
	)
	convey.Convey("IsFavedsCache", t, func(ctx convey.C) {
		favoreds, err := d.IsFavedsCache(c, tp, mid, oids)
		ctx.Convey("Then err should be nil.favoreds should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(favoreds, convey.ShouldNotBeNil)
		})
	})
}

func TestFavSetFavedBit(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("SetFavedBit", t, func(ctx convey.C) {
		err := d.SetFavedBit(c, tp, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavSetUnFavedBit(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("SetUnFavedBit", t, func(ctx convey.C) {
		err := d.SetUnFavedBit(c, tp, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavFavedBit(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
	)
	convey.Convey("FavedBit", t, func(ctx convey.C) {
		unfaved, err := d.FavedBit(c, tp, mid)
		ctx.Convey("Then err should be nil.unfaved should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(unfaved, convey.ShouldNotBeNil)
		})
	})
}

func TestFavDelRelationsCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(88888894)
		fid = int64(1)
	)
	convey.Convey("DelRelationsCache", t, func(ctx convey.C) {
		err := d.DelRelationsCache(c, mid, fid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavDelAllRelationsCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(88888894)
		fid = int64(1)
	)
	convey.Convey("DelRelationsCache", t, func(ctx convey.C) {
		err := d.DelAllRelationsCache(c, mid, fid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavRecentOidsCache(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(1)
		mid  = int64(88888894)
		fids = []int64{1}
	)
	convey.Convey("RecentOidsCache", t, func(ctx convey.C) {
		rctFidsMap, missFids, err := d.RecentOidsCache(c, typ, mid, fids)
		ctx.Convey("Then err should be nil.rctFidsMap,missFids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missFids, convey.ShouldNotBeNil)
			ctx.So(rctFidsMap, convey.ShouldNotBeNil)
		})
	})
}

func TestFavBatchOidsRedis(t *testing.T) {
	var (
		c     = context.TODO()
		tp    = int8(1)
		mid   = int64(88888894)
		limit = int(10)
	)
	convey.Convey("BatchOidsRedis", t, func(ctx convey.C) {
		oids, err := d.BatchOidsRedis(c, tp, mid, limit)
		ctx.Convey("Then err should be nil.oids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(oids, convey.ShouldNotBeNil)
		})
	})
}

func TestIsCleaned(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(2)
		mid = int64(88888894)
		fid = int64(59)
	)
	convey.Convey("IsCleaned", t, func(ctx convey.C) {
		_, err := d.IsCleaned(c, typ, mid, fid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestSetCleanedCache(t *testing.T) {
	var (
		c      = context.TODO()
		typ    = int8(2)
		mid    = int64(88888894)
		fid    = int64(59)
		ftime  = int64(88888894)
		expire = int64(86400)
	)
	convey.Convey("SetCleanedCache", t, func(ctx convey.C) {
		err := d.SetCleanedCache(c, typ, mid, fid, ftime, expire)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
