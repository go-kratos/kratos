package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaosynonymLike(t *testing.T) {
	var (
		keyWord = ""
	)
	convey.Convey("synonymLike", t, func(ctx convey.C) {
		p1 := synonymLike(keyWord)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSynonymCount(t *testing.T) {
	var (
		c       = context.TODO()
		keyWord = ""
	)
	convey.Convey("SynonymCount", t, func(ctx convey.C) {
		count, err := d.SynonymCount(c, keyWord)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSynonyms(t *testing.T) {
	var (
		c       = context.TODO()
		keyWord = ""
		start   = int32(0)
		end     = int32(0)
	)
	convey.Convey("Synonyms", t, func(ctx convey.C) {
		stagMap, stag, ids, err := d.Synonyms(c, keyWord, start, end)
		ctx.Convey("Then err should be nil.stagMap,stag,ids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ids, convey.ShouldHaveLength, 0)
			ctx.So(stag, convey.ShouldHaveLength, 0)
			ctx.So(stagMap, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoSynonymIDs(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{1, 2, 3}
	)
	convey.Convey("SynonymIDs", t, func(ctx convey.C) {
		mapST, stag, sids, err := d.SynonymIDs(c, ids)
		ctx.Convey("Then err should be nil.mapST,stag,sids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(sids, convey.ShouldNotBeEmpty)
			ctx.So(stag, convey.ShouldNotBeEmpty)
			ctx.So(mapST, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaoSynonymByName(t *testing.T) {
	var (
		c     = context.TODO()
		tname = ""
	)
	convey.Convey("SynonymByName", t, func(ctx convey.C) {
		res, err := d.SynonymByName(c, tname)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoInsertSynonym(t *testing.T) {
	var (
		c     = context.TODO()
		uname = ""
		ptid  = int64(0)
		tid   = int64(0)
	)
	convey.Convey("InsertSynonym", t, func(ctx convey.C) {
		id, err := d.InsertSynonym(c, uname, ptid, tid)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelSynonym(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("DelSynonym", t, func(ctx convey.C) {
		affect, err := d.DelSynonym(c, tid)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelSynonymSon(t *testing.T) {
	var (
		c    = context.TODO()
		ptid = int64(0)
		tids = []int64{1, 2, 3}
	)
	convey.Convey("DelSynonymSon", t, func(ctx convey.C) {
		affect, err := d.DelSynonymSon(c, ptid, tids)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoInsertSynonyms(t *testing.T) {
	var (
		c     = context.TODO()
		uname = ""
		ptid  = int64(1)
		tids  = []int64{1, 2, 3}
	)
	convey.Convey("InsertSynonyms", t, func(ctx convey.C) {
		id, err := d.InsertSynonyms(c, uname, ptid, tids)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoSynonym(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("Synonym", t, func(ctx convey.C) {
		res, err := d.Synonym(c, tid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
