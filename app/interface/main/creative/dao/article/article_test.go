package article

import (
	"context"
	artMdl "go-common/app/interface/main/creative/model/article"
	"testing"

	"go-common/app/interface/openplatform/article/model"
	"go-common/app/interface/openplatform/article/rpc/client"
	"go-common/library/ecode"
	"reflect"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestArticleArticles(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(1)
		pn       = int(0)
		ps       = int(0)
		sort     = int(0)
		group    = int(0)
		category = int(0)
		ip       = ""
	)

	convey.Convey("Articles", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "CreationUpperArticles",
			func(_ *client.Service, _ context.Context, _ *model.ArgCreationArts) (res *model.CreationArts, err error) {
				return nil, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		res, err := d.Articles(c, mid, pn, ps, sort, group, category, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestArticleCategories(t *testing.T) {
	var (
		c  = context.TODO()
		ip = ""
	)
	convey.Convey("Categories", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "Categories",
			func(_ *client.Service, _ context.Context, _ *model.ArgIP) (res *model.Categories, err error) {
				return nil, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		res, err := d.Categories(c, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestArticleCategoriesMap(t *testing.T) {
	var (
		c  = context.TODO()
		ip = ""
	)
	convey.Convey("CategoriesMap", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "CategoriesMap",
			func(_ *client.Service, _ context.Context, _ *model.ArgIP) (res map[int64]*model.Category, err error) {
				return nil, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		res, err := d.CategoriesMap(c, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, 20017)
			ctx.So(len(res), convey.ShouldEqual, 0)
		})
	})
}

func TestArticleArticle(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1198)
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("Article", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "CreationArticle",
			func(_ *client.Service, _ context.Context, _ *model.ArgAidMid) (res *model.Article, err error) {
				return nil, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		res, err := d.Article(c, aid, mid, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestArticleAddArticle(t *testing.T) {
	var (
		c   = context.TODO()
		art = &artMdl.ArtParam{}
	)
	convey.Convey("AddArticle", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "AddArticle",
			func(_ *client.Service, _ context.Context, _ *model.ArgArticle) (id int64, err error) {
				return 0, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		id, err := d.AddArticle(c, art)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(id, convey.ShouldEqual, 0)
		})
	})
}

func TestArticleUpdateArticle(t *testing.T) {
	var (
		c   = context.TODO()
		art = &artMdl.ArtParam{}
	)
	convey.Convey("UpdateArticle", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "UpdateArticle",
			func(_ *client.Service, _ context.Context, _ *model.ArgArticle) (err error) {
				return ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		err := d.UpdateArticle(c, art)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestArticleDelArticle(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(0)
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("DelArticle", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "DelArticle",
			func(_ *client.Service, _ context.Context, _ *model.ArgAidMid) (err error) {
				return ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		err := d.DelArticle(c, aid, mid, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestArticleWithDrawArticle(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(0)
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("WithDrawArticle", t, func(ctx convey.C) {
		err := d.WithDrawArticle(c, aid, mid, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestArticleIsAuthor(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("IsAuthor", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "IsAuthor",
			func(_ *client.Service, _ context.Context, _ *model.ArgMid) (res bool, err error) {
				return false, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		res, err := d.IsAuthor(c, mid, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldEqual, false)
		})
	})
}

func TestArticleRemainCount(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("ArticleRemainCount", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "ArticleRemainCount",
			func(_ *client.Service, _ context.Context, _ *model.ArgMid) (res int, err error) {
				return 0, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		res, err := d.RemainCount(c, mid, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldEqual, 0)
		})
	})
}

func TestArticleArticleStat(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2333)
		ip  = ""
	)
	convey.Convey("ArticleStat", t, func(ctx convey.C) {
		// mock
		//mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "CreationUpStat",
		//	func(_ *client.Service, _ context.Context, _ *model.ArgMid) (res model.UpStat, err error) {
		//		return new(model.UpStat), ecode.CreativeArticleRPCErr
		//	})
		//defer mock.Unpatch()
		res, err := d.ArticleStat(c, mid, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestArticleThirtyDayArticle(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("ThirtyDayArticle", t, func(ctx convey.C) {
		res, err := d.ThirtyDayArticle(c, mid, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestArticleArticleMetas(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{233}
		ip   = ""
	)
	convey.Convey("ArticleMetas", t, func(ctx convey.C) {
		//mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.art), "ArticleMetas",
			func(_ *client.Service, _ context.Context, _ *model.ArgAids) (res map[int64]*model.Meta, err error) {
				return nil, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		res, err := d.ArticleMetas(c, aids, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}
