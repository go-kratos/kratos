package dao

import (
	"context"
	"go-common/app/service/main/favorite/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestFavsend(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		msg = &model.Message{}
	)
	convey.Convey("send", t, func(ctx convey.C) {
		err := d.send(c, mid, msg)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavPubAddFav(t *testing.T) {
	var (
		c    = context.TODO()
		tp   = int8(0)
		mid  = int64(0)
		fid  = int64(0)
		oid  = int64(0)
		attr = int32(0)
		ts   = int64(0)
	)
	convey.Convey("PubAddFav", t, func(ctx convey.C) {
		d.PubAddFav(c, tp, mid, fid, oid, attr, ts, tp)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestFavPubDelFav(t *testing.T) {
	var (
		c    = context.TODO()
		tp   = int8(0)
		mid  = int64(0)
		fid  = int64(0)
		oid  = int64(0)
		attr = int32(0)
		ts   = int64(0)
	)
	convey.Convey("PubDelFav", t, func(ctx convey.C) {
		d.PubDelFav(c, tp, mid, fid, oid, attr, ts, tp)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestFavPubInitRelationFids(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(0)
		mid = int64(0)
	)
	convey.Convey("PubInitRelationFids", t, func(ctx convey.C) {
		d.PubInitRelationFids(c, tp, mid)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestFavPubInitFolderRelations(t *testing.T) {
	var (
		c   = context.TODO()
		tp  = int8(0)
		mid = int64(0)
		fid = int64(0)
	)
	convey.Convey("PubInitFolderRelations", t, func(ctx convey.C) {
		d.PubInitFolderRelations(c, tp, mid, fid)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestFavPubAddFolder(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		mid  = int64(0)
		fid  = int64(0)
		attr = int32(0)
	)
	convey.Convey("PubAddFolder", t, func(ctx convey.C) {
		d.PubAddFolder(c, typ, mid, fid, attr)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestFavPubDelFolder(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		mid  = int64(0)
		fid  = int64(0)
		attr = int32(0)
		ts   = int64(0)
	)
	convey.Convey("PubDelFolder", t, func(ctx convey.C) {
		d.PubDelFolder(c, typ, mid, fid, attr, ts)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestFavPubMultiDelFavs(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		mid  = int64(0)
		fid  = int64(0)
		rows = int64(0)
		attr = int32(0)
		oids = []int64{}
		ts   = int64(0)
	)
	convey.Convey("PubMultiDelFavs", t, func(ctx convey.C) {
		d.PubMultiDelFavs(c, typ, mid, fid, rows, attr, oids, ts)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestFavPubMultiAddFavs(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		mid  = int64(0)
		fid  = int64(0)
		rows = int64(0)
		attr = int32(0)
		oids = []int64{}
		ts   = int64(0)
	)
	convey.Convey("PubMultiAddFavs", t, func(ctx convey.C) {
		d.PubMultiAddFavs(c, typ, mid, fid, rows, attr, oids, ts)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestFavPubMoveFavs(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		mid  = int64(0)
		ofid = int64(0)
		nfid = int64(0)
		rows = int64(0)
		oids = []int64{}
		ts   = int64(0)
	)
	convey.Convey("PubMoveFavs", t, func(ctx convey.C) {
		d.PubMoveFavs(c, typ, mid, ofid, nfid, rows, oids, ts)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestFavPubCopyFavs(t *testing.T) {
	var (
		c    = context.TODO()
		typ  = int8(0)
		mid  = int64(0)
		ofid = int64(0)
		nfid = int64(0)
		rows = int64(0)
		oids = []int64{}
		ts   = int64(0)
	)
	convey.Convey("PubCopyFavs", t, func(ctx convey.C) {
		d.PubCopyFavs(c, typ, mid, ofid, nfid, rows, oids, ts)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}
