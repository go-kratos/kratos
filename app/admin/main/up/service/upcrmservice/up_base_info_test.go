package upcrmservice

import (
	"context"
	"testing"

	"go-common/app/admin/main/up/model/upcrmmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmserviceUpBaseInfoQuery(t *testing.T) {
	convey.Convey("UpBaseInfoQuery", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &upcrmmodel.InfoQueryArgs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.UpBaseInfoQuery(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceupBaseInfoQueryBatch(t *testing.T) {
	convey.Convey("upBaseInfoQueryBatch", t, func(ctx convey.C) {
		var (
			queryfunc QueryDbFunc
			ids       = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.upBaseInfoQueryBatch(queryfunc, ids)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceUpAccountInfo(t *testing.T) {
	convey.Convey("UpAccountInfo", t, func(ctx convey.C) {
		var (
			con = context.Background()
			arg = &upcrmmodel.InfoAccountInfoArgs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.UpAccountInfo(con, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicegetAttrFormat(t *testing.T) {
	convey.Convey("getAttrFormat", t, func(ctx convey.C) {
		var (
			attrs upcrmmodel.UpAttr
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := getAttrFormat(attrs)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicegetEsCombo(t *testing.T) {
	convey.Convey("getEsCombo", t, func(ctx convey.C) {
		var (
			attrs upcrmmodel.UpAttr
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			combos := getEsCombo(attrs)
			ctx.Convey("Then combos should not be nil.", func(ctx convey.C) {
				ctx.So(combos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceUpInfoSearch(t *testing.T) {
	convey.Convey("UpInfoSearch", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upcrmmodel.InfoSearchArgs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.UpInfoSearch(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceToPageInfo(t *testing.T) {
	convey.Convey("ToPageInfo", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			var p = esPage{}
			pageInfo := p.ToPageInfo()
			ctx.Convey("Then pageInfo should not be nil.", func(ctx convey.C) {
				ctx.So(pageInfo, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicesearchFromEs(t *testing.T) {
	convey.Convey("searchFromEs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upcrmmodel.InfoSearchArgs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			searchData, err := s.searchFromEs(c, arg)
			ctx.Convey("Then err should be nil.searchData should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(searchData, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicequeryUpBaseInfo(t *testing.T) {
	convey.Convey("queryUpBaseInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.queryUpBaseInfo(c, ids)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceQueryUpInfoWithViewerData(t *testing.T) {
	convey.Convey("QueryUpInfoWithViewerData", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upcrmmodel.UpInfoWithViewerArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.QueryUpInfoWithViewerData(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicegetOrCreateUpInfo(t *testing.T) {
	convey.Convey("getOrCreateUpInfo", t, func(ctx convey.C) {
		var (
			dataMap map[int64]*upcrmmodel.UpInfoWithViewerData
			mid     = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := getOrCreateUpInfo(dataMap, mid)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
