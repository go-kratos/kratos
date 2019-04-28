package fav

import (
	"context"
	favmdl "go-common/app/service/main/favorite/model"
	xtime "go-common/library/time"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestFavcntHit(t *testing.T) {
	convey.Convey("cntHit", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := cntHit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavfolderHit(t *testing.T) {
	convey.Convey("folderHit", t, func(ctx convey.C) {
		var (
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := folderHit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavrelationHit(t *testing.T) {
	convey.Convey("relationHit", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := relationHit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavusersHit(t *testing.T) {
	convey.Convey("usersHit", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := usersHit(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavPingMySQL(t *testing.T) {
	convey.Convey("PingMySQL", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMySQL(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestRecentFolder(t *testing.T) {
	convey.Convey("Folder Recents", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.RecentRes(c, mid, fid)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUpdateFavSequence(t *testing.T) {
	convey.Convey("Test Update Sequence", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			t, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			_, err = d.TxUpdateFavSequence(t, mid, fid, 1, 2, 123, xtime.Time(0))
			ctx.So(err, convey.ShouldBeNil)

		})
	})
}

func TestFavFolder(t *testing.T) {
	convey.Convey("Folder", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Folder(c, tp, mid, fid)
			ctx.Convey("Then err should be nil.f should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavUpFolderCnt(t *testing.T) {
	convey.Convey("UpFolderCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
			cnt = int(0)
			now xtime.Time
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpFolderCnt(c, mid, fid, cnt, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavRelation(t *testing.T) {
	convey.Convey("Relation", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
			fid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Relation(c, tp, mid, fid, oid)
			ctx.Convey("Then err should be nil.m should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavRelations(t *testing.T) {
	convey.Convey("Relations", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			typ   = int8(0)
			mid   = int64(0)
			fid   = int64(0)
			mtime xtime.Time
			limit = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			fr, err := d.Relations(c, typ, mid, fid, mtime, limit)
			ctx.Convey("Then err should be nil.fr should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fr, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavRelationFidsByOid(t *testing.T) {
	convey.Convey("RelationFidsByOid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.RelationFidsByOid(c, tp, mid, oid)
			ctx.Convey("Then err should be nil.fids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestFavRelationCnt(t *testing.T) {
	convey.Convey("RelationCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			fid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cnt, err := d.RelationCnt(c, mid, fid)
			ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cnt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavAddRelation(t *testing.T) {
	convey.Convey("AddRelation", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			fr = &favmdl.Favorite{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.AddRelation(c, fr)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavDelRelation(t *testing.T) {
	convey.Convey("DelRelation", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
			fid = int64(0)
			oid = int64(0)
			now xtime.Time
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelRelation(c, tp, mid, fid, oid, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavUpStatCnt(t *testing.T) {
	convey.Convey("UpStatCnt", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int8(0)
			oid  = int64(0)
			incr = int(0)
			now  xtime.Time
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpStatCnt(c, tp, oid, incr, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavStatCnt(t *testing.T) {
	convey.Convey("StatCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cnt, err := d.StatCnt(c, tp, oid)
			ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cnt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavRelationFids(t *testing.T) {
	convey.Convey("RelationFids", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int8(0)
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rfids, err := d.RelationFids(c, tp, mid)
			ctx.Convey("Then err should be nil.rfids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rfids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavOidsByFid(t *testing.T) {
	convey.Convey("OidsByFid", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			typ    = int8(0)
			mid    = int64(0)
			fid    = int64(0)
			offset = int(0)
			limit  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			oids, err := d.OidsByFid(c, typ, mid, fid, offset, limit)
			ctx.Convey("Then err should be nil.oids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(oids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavDelRelationsByOids(t *testing.T) {
	convey.Convey("DelRelationsByOids", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			typ  = int8(0)
			mid  = int64(0)
			fid  = int64(0)
			oids = []int64{}
			now  xtime.Time
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelRelationsByOids(c, typ, mid, fid, oids, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavAddUser(t *testing.T) {
	convey.Convey("AddUser", t, func(ctx convey.C) {
		var (
			c = context.Background()
			u = &favmdl.User{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.AddUser(c, u)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestFavDelUser(t *testing.T) {
	convey.Convey("DelUser", t, func(ctx convey.C) {
		var (
			c = context.Background()
			u = &favmdl.User{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelUser(c, u)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRelationsByOids(t *testing.T) {
	convey.Convey("RelationsByOids", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			favs, err := d.RelationsByOids(c, 2, 1501, 168, []int64{10108138})
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(favs, convey.ShouldNotBeNil)
			_, err = d.BatchUpdateSeq(c, 76, favs)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
