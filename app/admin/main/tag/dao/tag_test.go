package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTag(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("Tag", t, func(ctx convey.C) {
		tag, err := d.Tag(c, tid)
		ctx.Convey("Then err should be nil.tag should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tag, convey.ShouldBeNil)
		})
	})
}

func TestDaoTags(t *testing.T) {
	var (
		c    = context.TODO()
		tids = []int64{1, 2, 3}
	)
	convey.Convey("Tags", t, func(ctx convey.C) {
		tags, tagMap, err := d.Tags(c, tids)
		ctx.Convey("Then err should be nil.tags,tagMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tagMap, convey.ShouldNotBeNil)
			ctx.So(tags, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagByName(t *testing.T) {
	var (
		c     = context.TODO()
		tname = ""
	)
	convey.Convey("TagByName", t, func(ctx convey.C) {
		tag, err := d.TagByName(c, tname)
		ctx.Convey("Then err should be nil.tag should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tag, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagByNames(t *testing.T) {
	var (
		c      = context.TODO()
		tnames = []string{"123", "456"}
	)
	convey.Convey("TagByNames", t, func(ctx convey.C) {
		tags, tagMap, tagNameMap, err := d.TagByNames(c, tnames)
		ctx.Convey("Then err should be nil.tags,tagMap,tagNameMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tagNameMap, convey.ShouldNotBeNil)
			ctx.So(tagMap, convey.ShouldNotBeNil)
			ctx.So(tags, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpVerifyState(t *testing.T) {
	var (
		c           = context.TODO()
		tid         = int64(0)
		verifyState = int32(0)
	)
	convey.Convey("UpVerifyState", t, func(ctx convey.C) {
		affect, err := d.UpVerifyState(c, tid, verifyState)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxInsertTag(t *testing.T) {
	var tag = &model.Tag{
		Name: "12345678",
	}
	convey.Convey("TxInsertTag", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxInsertTag(tx, tag)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoTxInsertTagCount(t *testing.T) {
	var (
		tid = int64(778899101010)
	)
	convey.Convey("TxInsertTagCount", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxInsertTagCount(tx, tid)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoUpdateTag(t *testing.T) {
	var (
		c   = context.TODO()
		tag = &model.Tag{}
	)
	convey.Convey("UpdateTag", t, func(ctx convey.C) {
		affect, err := d.UpdateTag(c, tag)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxUpdateTag(t *testing.T) {
	var (
		tag = &model.Tag{}
	)
	convey.Convey("TxUpdateTag", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		affect, err := d.TxUpdateTag(tx, tag)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}

func TestDaoUpTagState(t *testing.T) {
	var (
		c     = context.TODO()
		tid   = int64(0)
		state = int32(0)
	)
	convey.Convey("UpTagState", t, func(ctx convey.C) {
		affect, err := d.UpTagState(c, tid, state)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagCount(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("TagCount", t, func(ctx convey.C) {
		res, err := d.TagCount(c, tid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagCounts(t *testing.T) {
	var (
		c    = context.TODO()
		tids = []int64{1, 2, 3}
	)
	convey.Convey("TagCounts", t, func(ctx convey.C) {
		res, err := d.TagCounts(c, tids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
