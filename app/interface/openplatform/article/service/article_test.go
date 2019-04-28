package service

import (
	"context"
	"testing"

	artmdl "go-common/app/interface/openplatform/article/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Article(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		res, err := s.Article(context.TODO(), dataID)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("no data", t, WithService(func(s *Service) {
		res, err := s.Article(context.TODO(), noDataID)
		So(err, ShouldBeNil)
		So(res, ShouldBeNil)
	}))
	Convey("ArticleRemainCount", t, WithService(func(s *Service) {
		_, err := s.ArticleRemainCount(context.TODO(), art.Author.Mid)
		So(err, ShouldBeNil)
	}))
}

func Test_ArticleMetas(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		res, err := s.ArticleMetas(context.TODO(), []int64{dataID})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("no data", t, WithService(func(s *Service) {
		res, err := s.ArticleMetas(context.TODO(), []int64{noDataID})
		So(err, ShouldBeNil)
		So(res, ShouldBeEmpty)
	}))
}

func Test_AddArticleCache(t *testing.T) {
	Convey("add data", t, WithService(func(s *Service) {
		var c = context.TODO()
		err := s.AddArticleCache(context.TODO(), dataID)
		So(err, ShouldBeNil)

		Convey("del cache", func() {
			err := s.DelArticleCache(c, 175, dataID)
			So(err, ShouldBeNil)
			Convey("delete twice return null", func() {
				err := s.DelArticleCache(c, 175, dataID)
				So(err, ShouldBeNil)
			})
		})
	}))
}

func Test_FilterNoDistributeArts(t *testing.T) {
	a1 := artmdl.Meta{ID: 1}
	a2 := artmdl.Meta{ID: 2}
	a3 := artmdl.Meta{ID: 3}
	a2.AttrSet(int32(1), artmdl.AttrBitNoDistribute)

	Convey("array work", t, WithService(func(s *Service) {
		res := filterNoDistributeArts([]*artmdl.Meta{&a1, &a2, &a3})
		So(res, ShouldResemble, []*artmdl.Meta{&a1, &a3})
	}))

	Convey("map work", t, WithService(func(s *Service) {
		arg := map[int64]*artmdl.Meta{1: &a1, 2: &a2, 3: &a3}
		res := map[int64]*artmdl.Meta{1: &a1, 3: &a3}
		filterNoDistributeArtsMap(arg)
		So(res, ShouldResemble, res)
	}))
}

func Test_fmtMoreArts(t *testing.T) {
	a1 := &artmdl.Meta{ID: 1, PublishTime: xtime.Time(1)}
	a2 := &artmdl.Meta{ID: 2, PublishTime: xtime.Time(2)}
	a3 := &artmdl.Meta{ID: 3, PublishTime: xtime.Time(3)}
	a4 := &artmdl.Meta{ID: 4, PublishTime: xtime.Time(4)}
	a5 := &artmdl.Meta{ID: 5, PublishTime: xtime.Time(5)}
	m := map[int64]*artmdl.Meta{1: a1, 2: a2, 3: a3, 4: a4, 5: a5}
	Convey("position: x5432", t, func() {
		res := fmtMoreArts([]int64{2, 3, 4, 5}, []int64{}, m)
		So(res, ShouldResemble, []*artmdl.Meta{a5, a4, a3, a2})
	})
	Convey("position: x32", t, func() {
		res := fmtMoreArts([]int64{2, 3}, []int64{}, m)
		So(res, ShouldResemble, []*artmdl.Meta{a3, a2})
	})
	Convey("position: 54x321", t, func() {
		res := fmtMoreArts([]int64{3, 2, 1}, []int64{5, 4}, m)
		So(res, ShouldResemble, []*artmdl.Meta{a4, a3, a2, a1})
	})
	Convey("position: 5432x1", t, func() {
		res := fmtMoreArts([]int64{1}, []int64{5, 4, 3, 2}, m)
		So(res, ShouldResemble, []*artmdl.Meta{a4, a3, a2, a1})
	})
	Convey("position: 4321x", t, func() {
		res := fmtMoreArts([]int64{}, []int64{4, 3, 2, 1}, m)
		So(res, ShouldResemble, []*artmdl.Meta{a4, a3, a2, a1})
	})
	Convey("position: 2x1", t, func() {
		res := fmtMoreArts([]int64{1}, []int64{2}, m)
		So(res, ShouldResemble, []*artmdl.Meta{a2, a1})
	})
}

func Test_splitAids(t *testing.T) {
	aids := []int64{4, 3, 2, 1}
	Convey("position: 4", t, func() {
		before, after := splitAids(aids, 4)
		So(after, ShouldResemble, []int64{})
		So(before, ShouldResemble, []int64{3, 2, 1})
	})
	Convey("position: 3", t, func() {
		before, after := splitAids(aids, 3)
		So(after, ShouldResemble, []int64{4})
		So(before, ShouldResemble, []int64{2, 1})
	})
	Convey("position: 2", t, func() {
		before, after := splitAids(aids, 2)
		So(after, ShouldResemble, []int64{4, 3})
		So(before, ShouldResemble, []int64{1})
	})
	Convey("position: 1", t, func() {
		before, after := splitAids(aids, 1)
		So(after, ShouldResemble, []int64{4, 3, 2})
		So(before, ShouldResemble, []int64{})
	})
}
