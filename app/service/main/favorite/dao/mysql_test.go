package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/favorite/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestFavfolderHit(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("folderHit", t, func(ctx convey.C) {
		p1 := folderHit(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavrelationHit(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("relationHit", t, func(ctx convey.C) {
		p1 := relationHit(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavusersHit(t *testing.T) {
	var (
		oid = int64(0)
	)
	convey.Convey("usersHit", t, func(ctx convey.C) {
		p1 := usersHit(oid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavcountHit(t *testing.T) {
	var (
		oid = int64(0)
	)
	convey.Convey("countHit", t, func(ctx convey.C) {
		p1 := countHit(oid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestFavpingMySQL(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("pingMySQL", t, func(ctx convey.C) {
		err := d.pingMySQL(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavFolder(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
		fid = int64(1)
	)
	convey.Convey("Folder", t, func(ctx convey.C) {
		f, err := d.Folder(c, tp, mid, fid)
		ctx.Convey("Then err should be nil.f should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(f, convey.ShouldNotBeNil)
		})
	})
}

func TestFavFolderByName(t *testing.T) {
	var (
		c    = context.TODO()
		tp   = int8(0)
		mid  = int64(0)
		name = ""
	)
	convey.Convey("FolderByName", t, func(ctx convey.C) {
		f, err := d.FolderByName(c, tp, mid, name)
		ctx.Convey("Then err should be nil.f should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(f, convey.ShouldNotBeNil)
		})
	})
}

func TestFavDefaultFolder(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(0)
		mid = int64(0)
	)
	convey.Convey("DefaultFolder", t, func(ctx convey.C) {
		f, err := d.DefaultFolder(c, tp, mid)
		ctx.Convey("Then err should be nil.f should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(f, convey.ShouldNotBeNil)
		})
	})
}

func TestFavAddFolder(t *testing.T) {
	var (
		c = context.TODO()
		f = &model.Folder{}
	)
	convey.Convey("AddFolder", t, func(ctx convey.C) {
		fid, err := d.AddFolder(c, f)
		ctx.Convey("Then err should be nil.fid should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fid, convey.ShouldNotBeNil)
		})
	})
}

func TestFavUpdateFolder(t *testing.T) {
	var (
		c = context.TODO()
		f = &model.Folder{}
	)
	convey.Convey("UpdateFolder", t, func(ctx convey.C) {
		fid, err := d.UpdateFolder(c, f)
		ctx.Convey("Then err should be nil.fid should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fid, convey.ShouldNotBeNil)
		})
	})
}

func TestFavUpFolderName(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		mid  = int64(0)
		fid  = int64(0)
		name = ""
	)
	convey.Convey("UpFolderName", t, func(ctx convey.C) {
		rows, err := d.UpFolderName(c, typ, mid, fid, name)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavUpFolderAttr(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		mid  = int64(0)
		fid  = int64(0)
		attr = int32(0)
	)
	convey.Convey("UpFolderAttr", t, func(ctx convey.C) {
		rows, err := d.UpFolderAttr(c, typ, mid, fid, attr)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavFolderRelations(t *testing.T) {
	var (
		c     = context.TODO()
		typ   = int8(0)
		mid   = int64(0)
		fid   = int64(0)
		start = int(0)
		end   = int(0)
	)
	convey.Convey("FolderRelations", t, func(ctx convey.C) {
		fr, err := d.FolderRelations(c, typ, mid, fid, start, end)
		ctx.Convey("Then err should be nil.fr should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fr, convey.ShouldNotBeNil)
		})
	})
}

func TestFavFolders(t *testing.T) {
	var (
		c      = context.TODO()
		fvmids = []*model.ArgFVmid{}
	)
	convey.Convey("Folders", t, func(ctx convey.C) {
		fs, err := d.Folders(c, fvmids)
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fs, convey.ShouldNotBeNil)
		})
	})
}

func TestFavRelationFidsByOid(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(1)
		mid = int64(88888894)
		oid = int64(1)
	)
	convey.Convey("RelationFidsByOid", t, func(ctx convey.C) {
		_, err := d.RelationFidsByOid(c, tp, mid, oid)
		ctx.Convey("Then err should be nil.fids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavRelationFidsByOids(t *testing.T) {
	var (
		c    = context.TODO()
		tp   = int8(1)
		mid  = int64(8888894)
		oids = []int64{1, 2, 3}
	)
	convey.Convey("RelationFidsByOids", t, func(ctx convey.C) {
		fidsMap, err := d.RelationFidsByOids(c, tp, mid, oids)
		ctx.Convey("Then err should be nil.fidsMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fidsMap, convey.ShouldNotBeNil)
		})
	})
}

func TestFavCntRelations(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		fid = int64(0)
	)
	convey.Convey("CntRelations", t, func(ctx convey.C) {
		count, err := d.CntRelations(c, mid, fid, 2)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestFavFolderCnt(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(0)
		mid = int64(0)
	)
	convey.Convey("FolderCnt", t, func(ctx convey.C) {
		count, err := d.FolderCnt(c, tp, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestFavAddFav(t *testing.T) {
	var (
		c  = context.TODO()
		fr = &model.Favorite{}
	)
	convey.Convey("AddFav", t, func(ctx convey.C) {
		rows, err := d.AddFav(c, fr)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavDelFav(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(0)
		mid = int64(0)
		fid = int64(0)
		oid = int64(0)
	)
	convey.Convey("DelFav", t, func(ctx convey.C) {
		rows, err := d.DelFav(c, tp, mid, fid, oid)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavAddRelation(t *testing.T) {
	var (
		c  = context.TODO()
		fr = &model.Favorite{}
	)
	convey.Convey("AddRelation", t, func(ctx convey.C) {
		rows, err := d.AddRelation(c, fr)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavRelation(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(0)
		mid = int64(0)
		fid = int64(0)
		oid = int64(0)
	)
	convey.Convey("Relation", t, func(ctx convey.C) {
		_, err := d.Relation(c, tp, mid, fid, oid)
		ctx.Convey("Then err should be nil.m should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavDelRelation(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(0)
		mid = int64(0)
		fid = int64(0)
		oid = int64(0)
	)
	convey.Convey("DelRelation", t, func(ctx convey.C) {
		rows, err := d.DelRelation(c, tp, mid, fid, oid)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavMultiDelRelations(t *testing.T) {
	var (
		c    = context.TODO()
		tp   = int8(1)
		mid  = int64(88888894)
		fid  = int64(0)
		oids = []int64{1, 2, 3}
	)
	convey.Convey("MultiDelRelations", t, func(ctx convey.C) {
		rows, err := d.MultiDelRelations(c, tp, mid, fid, oids)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavTxMultiDelRelations(t *testing.T) {
	var (
		tx, _ = d.BeginTran(context.TODO())
		tp    = int8(1)
		mid   = int64(88888894)
		fid   = int64(0)
		oids  = []int64{1, 2, 3}
	)
	convey.Convey("TxMultiDelRelations", t, func(ctx convey.C) {
		rows, err := d.TxMultiDelRelations(tx, tp, mid, fid, oids)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavMultiAddRelations(t *testing.T) {
	var (
		c    = context.TODO()
		tp   = int8(1)
		mid  = int64(88888894)
		fid  = int64(0)
		oids = []int64{1}
	)
	convey.Convey("MultiAddRelations", t, func(ctx convey.C) {
		rows, err := d.MultiAddRelations(c, tp, mid, fid, oids)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavDelFolder(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(0)
		mid = int64(0)
		fid = int64(0)
	)
	convey.Convey("DelFolder", t, func(ctx convey.C) {
		rows, err := d.DelFolder(c, tp, mid, fid)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavUserFolders(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		mid = int64(0)
	)
	convey.Convey("UserFolders", t, func(ctx convey.C) {
		fs, err := d.UserFolders(c, typ, mid)
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fs, convey.ShouldNotBeNil)
		})
	})
}

func TestFavFolderSort(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		mid = int64(0)
	)
	convey.Convey("FolderSort", t, func(ctx convey.C) {
		fst, err := d.FolderSort(c, typ, mid)
		ctx.Convey("Then err should be nil.fst should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fst, convey.ShouldNotBeNil)
		})
	})
}

func TestFavUpFolderSort(t *testing.T) {
	var (
		c   = context.TODO()
		fst = &model.FolderSort{}
	)
	convey.Convey("UpFolderSort", t, func(ctx convey.C) {
		rows, err := d.UpFolderSort(c, fst)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavRecentOids(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(88888894)
		fid = int64(1)
	)
	convey.Convey("RecentOids", t, func(ctx convey.C) {
		oids, err := d.RecentOids(c, mid, fid, 1)
		ctx.Convey("Then err should be nil.oids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(oids, convey.ShouldNotBeNil)
		})
	})
}

func TestFavTxCopyRelations(t *testing.T) {
	var (
		tx, _  = d.BeginTran(context.TODO())
		typ    = int8(1)
		oldmid = int64(88888894)
		mid    = int64(88888894)
		oldfid = int64(0)
		newfid = int64(0)
		oids   = []int64{1}
	)
	convey.Convey("TxCopyRelations", t, func(ctx convey.C) {
		rows, err := d.TxCopyRelations(tx, typ, oldmid, mid, oldfid, newfid, oids)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavCopyRelations(t *testing.T) {
	var (
		c      = context.TODO()
		typ    = int8(1)
		oldmid = int64(88888894)
		mid    = int64(88888894)
		oldfid = int64(0)
		newfid = int64(0)
		oids   = []int64{1}
	)
	convey.Convey("CopyRelations", t, func(ctx convey.C) {
		rows, err := d.CopyRelations(c, typ, oldmid, mid, oldfid, newfid, oids)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestFavCntUsers(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		oid = int64(0)
	)
	convey.Convey("CntUsers", t, func(ctx convey.C) {
		count, err := d.CntUsers(c, typ, oid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestFavUsers(t *testing.T) {
	var (
		c     = context.TODO()
		typ   = int8(0)
		oid   = int64(0)
		start = int(0)
		end   = int(0)
	)
	convey.Convey("Users", t, func(ctx convey.C) {
		us, err := d.Users(c, typ, oid, start, end)
		ctx.Convey("Then err should be nil.us should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(us, convey.ShouldNotBeNil)
		})
	})
}

func TestFavOidCount(t *testing.T) {
	var (
		c   = context.TODO()
		typ = int8(0)
		oid = int64(0)
	)
	convey.Convey("OidCount", t, func(ctx convey.C) {
		count, err := d.OidCount(c, typ, oid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestFavOidsCount(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(1)
		oids = []int64{1, 2, 3}
	)
	convey.Convey("OidsCount", t, func(ctx convey.C) {
		counts, err := d.OidsCount(c, typ, oids)
		ctx.Convey("Then err should be nil.counts should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(counts, convey.ShouldNotBeNil)
		})
	})
}

func TestFavBatchOids(t *testing.T) {
	var (
		c     = context.TODO()
		typ   = int8(1)
		mid   = int64(88888894)
		limit = int(10)
	)
	convey.Convey("BatchOids", t, func(ctx convey.C) {
		oids, err := d.BatchOids(c, typ, mid, limit)
		ctx.Convey("Then err should be nil.oids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(oids, convey.ShouldNotBeNil)
		})
	})
}
