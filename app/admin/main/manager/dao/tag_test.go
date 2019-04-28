package dao

import (
	"context"
	"math/rand"
	"testing"

	"go-common/app/admin/main/manager/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddType(t *testing.T) {
	var (
		c  = context.TODO()
		tt = &model.TagType{}
	)
	convey.Convey("AddType", t, func(ctx convey.C) {
		err := d.AddType(c, tt)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateTypeName(t *testing.T) {
	var (
		c  = context.TODO()
		tt = &model.TagType{}
	)
	convey.Convey("UpdateTypeName", t, func(ctx convey.C) {
		err := d.UpdateTypeName(c, tt)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateType(t *testing.T) {
	var (
		c  = context.TODO()
		tt = &model.TagType{}
	)
	convey.Convey("UpdateType", t, func(ctx convey.C) {
		err := d.UpdateType(c, tt)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDeleteNonRole(t *testing.T) {
	var (
		c  = context.TODO()
		tt = &model.TagType{}
	)
	convey.Convey("DeleteNonRole", t, func(ctx convey.C) {
		err := d.DeleteNonRole(c, tt)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDeleteType(t *testing.T) {
	var (
		c  = context.TODO()
		td = &model.TagTypeDel{}
	)
	convey.Convey("DeleteType", t, func(ctx convey.C) {
		err := d.DeleteType(c, td)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddTag(t *testing.T) {
	var (
		c  = context.TODO()
		no = &model.Tag{}
	)
	convey.Convey("AddTag", t, func(ctx convey.C) {
		err := d.AddTag(c, no)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaomaxTagIDByBid(t *testing.T) {
	var (
		c   = context.TODO()
		bid = int64(0)
	)
	convey.Convey("maxTagIDByBid", t, func(ctx convey.C) {
		maxTagID, err := d.maxTagIDByBid(c, bid)
		ctx.Convey("Then err should be nil.maxTagID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(maxTagID, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateTag(t *testing.T) {
	var (
		c  = context.TODO()
		no = &model.Tag{}
	)
	convey.Convey("UpdateTag", t, func(ctx convey.C) {
		err := d.UpdateTag(c, no)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddControl(t *testing.T) {
	var (
		c  = context.TODO()
		tc = &model.TagControl{}
	)
	convey.Convey("AddControl", t, func(ctx convey.C) {
		err := d.AddControl(c, tc)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateControl(t *testing.T) {
	var (
		c  = context.TODO()
		tc = &model.TagControl{}
	)
	convey.Convey("UpdateControl", t, func(ctx convey.C) {
		err := d.UpdateControl(c, tc)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBatchUpdateState(t *testing.T) {
	var (
		c = context.TODO()
		b = &model.BatchUpdateState{}
	)
	convey.Convey("BatchUpdateState", t, func(ctx convey.C) {
		err := d.BatchUpdateState(c, b)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTagList(t *testing.T) {
	var (
		c  = context.TODO()
		no = &model.SearchTagParams{}
	)
	convey.Convey("TagList", t, func(ctx convey.C) {
		res, err := d.TagList(c, no)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagControl(t *testing.T) {
	var (
		c  = context.TODO()
		tc = &model.TagControlParam{
			BID: 1,
			TID: 2,
		}
	)
	convey.Convey("TagControl", t, func(ctx convey.C) {
		res, err := d.TagControl(c, tc)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAttrList(t *testing.T) {
	var (
		c   = context.TODO()
		bid = int64(0)
	)
	convey.Convey("AttrList", t, func(ctx convey.C) {
		res, err := d.AttrList(c, bid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInsertAttr(t *testing.T) {
	bid := rand.Intn(1000)
	tba := &model.TagBusinessAttr{
		Bid: int64(bid),
	}
	convey.Convey("InsertAttr", t, func(ctx convey.C) {
		err := d.InsertAttr(context.Background(), tba)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAttrUpdate(t *testing.T) {
	var (
		c   = context.TODO()
		tba = &model.TagBusinessAttr{}
	)
	convey.Convey("AttrUpdate", t, func(ctx convey.C) {
		err := d.AttrUpdate(c, tba)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTypeByIDs(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{}
	)
	convey.Convey("TypeByIDs", t, func(ctx convey.C) {
		res, err := d.TypeByIDs(c, ids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagTypeByBID(t *testing.T) {
	var (
		c   = context.TODO()
		bid = int64(0)
	)
	convey.Convey("TagTypeByBID", t, func(ctx convey.C) {
		res, err := d.TagTypeByBID(c, bid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagTypeRoleByTids(t *testing.T) {
	var (
		c    = context.TODO()
		tids = []int64{}
	)
	convey.Convey("TagTypeRoleByTids", t, func(ctx convey.C) {
		res, err := d.TagTypeRoleByTids(c, tids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagByType(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("TagByType", t, func(ctx convey.C) {
		res, err := d.TagByType(c, tid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
