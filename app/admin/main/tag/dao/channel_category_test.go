package dao

import (
	"context"
	"go-common/app/admin/main/tag/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaosetChannelCategory(t *testing.T) {
	var (
		categories = []*model.ChannelCategory{}
	)
	convey.Convey("setChannelCategory", t, func(ctx convey.C) {
		p1 := setChannelCategory(categories)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCountChannelCategory(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("CountChannelCategory", t, func(ctx convey.C) {
		count, err := d.CountChannelCategory(c)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoChannelCategory(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("ChannelCategory", t, func(ctx convey.C) {
		res, err := d.ChannelCategory(c, id)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoChannelCategoryByName(t *testing.T) {
	var (
		c    = context.TODO()
		name = "bilibili"
	)
	convey.Convey("ChannelCategoryByName", t, func(ctx convey.C) {
		res, err := d.ChannelCategoryByName(c, name)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoChannelCategories(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("ChannelCategories", t, func(ctx convey.C) {
		res, err := d.ChannelCategories(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInsertChannelCategory(t *testing.T) {
	var (
		c        = context.TODO()
		category = &model.ChannelCategory{
			Name:      "ut",
			Order:     0,
			State:     0,
			INTShield: 1,
		}
	)
	convey.Convey("InsertChannelCategory", t, func(ctx convey.C) {
		id, err := d.InsertChannelCategory(c, category)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoStateChannelCategory(t *testing.T) {
	var (
		c     = context.TODO()
		id    = int64(0)
		state = int32(0)
	)
	convey.Convey("StateChannelCategory", t, func(ctx convey.C) {
		affect, err := d.StateChannelCategory(c, id, state)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateChannelCategories(t *testing.T) {
	var (
		c          = context.TODO()
		categories = []*model.ChannelCategory{
			{
				ID:   123,
				Name: "unit test",
			},
		}
	)
	convey.Convey("UpdateChannelCategories", t, func(ctx convey.C) {
		affect, err := d.UpdateChannelCategories(c, categories)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoTxUpChannelCategoryAttr(t *testing.T) {
	var category = &model.ChannelCategory{
		ID:   1,
		Attr: 1,
	}
	convey.Convey("TxUpChannelCategoryAttr", t, func(ctx convey.C) {
		tx, err := d.BeginTran(context.TODO())
		if err != nil {
			return
		}
		id, err := d.TxUpChannelCategoryAttr(tx, category)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		tx.Rollback()
	})
}
