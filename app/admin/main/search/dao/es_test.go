package dao

import (
	"context"
	"go-common/app/admin/main/search/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/olivere/elastic.v5"
)

func TestDaoUpdateMapBulk(t *testing.T) {
	convey.Convey("UpdateMapBulk", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			esName   = ""
			bulkData = []BulkMapItem{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//err :=
			d.UpdateMapBulk(c, esName, bulkData)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateBulk(t *testing.T) {
	convey.Convey("UpdateBulk", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			esName   = ""
			bulkData = []BulkItem{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//err :=
			d.UpdateBulk(c, esName, bulkData)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpsertBulk(t *testing.T) {
	convey.Convey("UpsertBulk", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			esCluster = ""
			up        = &model.UpsertParams{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpsertBulk(c, esCluster, up)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaosearchResult(t *testing.T) {
	convey.Convey("searchResult", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			esClusterName = ""
			indexName     = ""
			query         elastic.Query
			bsp           = &model.BasicSearchParams{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.searchResult(c, esClusterName, indexName, query, bsp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryResult(t *testing.T) {
	convey.Convey("QueryResult", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query elastic.Query
			sp    = &model.QueryParams{
				AppIDConf: &model.QueryConfDetail{
					ESCluster: "",
				},
			}
			qbDebug = &model.QueryDebugResult{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, qrDebug, err := d.QueryResult(c, query, sp, qbDebug)
			ctx.Convey("Then err should be nil.res,qrDebug should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(qrDebug, convey.ShouldNotBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBulkIndex(t *testing.T) {
	convey.Convey("BulkIndex", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			esName   = ""
			bulkData = []BulkItem{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//err :=
			d.BulkIndex(c, esName, bulkData)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoExistIndex(t *testing.T) {
	convey.Convey("ExistIndex", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			esClusterName = ""
			indexName     = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			exist, _ := d.ExistIndex(c, esClusterName, indexName)
			ctx.Convey("Then err should be nil.exist should not be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
				ctx.So(exist, convey.ShouldNotBeNil)
			})
		})
	})
}
