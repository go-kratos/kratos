package service

import (
	"go-common/app/admin/ep/melloi/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceTreesQuery(t *testing.T) {
	convey.Convey("TreesQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.TreesQuery()
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceTreeNumQuery(t *testing.T) {
	convey.Convey("TreeNumQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.TreeNumQuery()
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceTopHttpQuery(t *testing.T) {
	convey.Convey("TopHttpQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.TopHttpQuery()
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceTopGrpcQuery(t *testing.T) {
	convey.Convey("TopGrpcQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.TopGrpcQuery()
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceTopSceneQuery(t *testing.T) {
	convey.Convey("TopSceneQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.TopSceneQuery()
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceTopDeptQuery(t *testing.T) {
	convey.Convey("TopDeptQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.TopDeptQuery()
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceBuildLineQuery(t *testing.T) {
	convey.Convey("BuildLineQuery", t, func(convCtx convey.C) {
		var (
			rank    = &model.Rank{}
			summary = &model.ReportSummary{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.BuildLineQuery(rank, summary)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceStateLineQuery(t *testing.T) {
	convey.Convey("StateLineQuery", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.StateLineQuery()
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
