package dao

import (
	"go-common/app/admin/ep/melloi/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTreesQuery(t *testing.T) {
	convey.Convey("TreesQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.TreesQuery()
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTreeNumQuery(t *testing.T) {
	convey.Convey("TreeNumQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.TreeNumQuery()
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTopHttpQuery(t *testing.T) {
	convey.Convey("TopHttpQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.TopHttpQuery()
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTopGrpcQuery(t *testing.T) {
	convey.Convey("TopGrpcQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.TopGrpcQuery()
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTopSceneQuery(t *testing.T) {
	convey.Convey("TopSceneQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.TopSceneQuery()
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTopDeptQuery(t *testing.T) {
	convey.Convey("TopDeptQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.TopDeptQuery()
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBuildLineQuery(t *testing.T) {
	convey.Convey("BuildLineQuery", t, func(convCtx convey.C) {
		var (
			rank    = &model.Rank{}
			summary = &model.ReportSummary{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BuildLineQuery(rank, summary)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStateLineQuery(t *testing.T) {
	convey.Convey("StateLineQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.StateLineQuery()
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
