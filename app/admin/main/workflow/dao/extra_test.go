package dao

import (
	"context"
	"net/url"
	"testing"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBatchUperSpecial(t *testing.T) {
	convey.Convey("BatchUperSpecial", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{27515256}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			UperTagMap, err := d.BatchUperSpecial(c, mids)
			ctx.Convey("Then err should be nil.UperTagMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(UperTagMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArchiveRPC(t *testing.T) {
	convey.Convey("ArchiveRPC", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			archives, err := d.ArchiveRPC(c, oids)
			ctx.Convey("Then err should be nil.archives should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(archives, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagList(t *testing.T) {
	convey.Convey("TagList", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tags, err := d.TagList(c)
			ctx.Convey("Then err should be nil.tags should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tags, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCommonExtraInfo(t *testing.T) {
	convey.Convey("CommonExtraInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			bid  = int8(13)
			uri  = "http://uat-manager.bilibili.co/x/admin/reply/internal/reply"
			ids  = []int64{1}
			oids = []int64{1}
			eids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.CommonExtraInfo(c, bid, uri, ids, oids, eids)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAccountInfoRPC(t *testing.T) {
	convey.Convey("AccountInfoRPC", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			authors := d.AccountInfoRPC(c, mids)
			ctx.Convey("Then authors should not be nil.", func(ctx convey.C) {
				ctx.So(authors, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddMoral(t *testing.T) {
	convey.Convey("AddMoral", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1}
			gssp = &param.GroupStateSetParam{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddMoral(c, mids, gssp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err.Error(), convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddBlock(t *testing.T) {
	convey.Convey("AddBlock", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1}
			gssp = &param.GroupStateSetParam{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddBlock(c, mids, gssp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCreditBlockInfo(t *testing.T) {
	convey.Convey("AddCreditBlockInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			bus  map[int64]*model.Business
			gssp = &param.GroupStateSetParam{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCreditBlockInfo(c, bus, gssp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCreditCase(t *testing.T) {
	convey.Convey("AddCreditCase", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			uv url.Values
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCreditCase(c, uv)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, ecode.Int(56504))
			})
		})
	})
}

func TestDaoBlockNum(t *testing.T) {
	convey.Convey("BlockNum", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sum, err := d.BlockNum(c, mid)
			ctx.Convey("Then err should be nil.sum should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sum, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockInfo(t *testing.T) {
	convey.Convey("BlockInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			resp, err := d.BlockInfo(c, mid)
			ctx.Convey("Then err should be nil.resp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(resp, convey.ShouldNotBeNil)
			})
		})
	})
}
