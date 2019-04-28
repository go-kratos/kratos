package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	r1 = &model.Recommend{ArticleID: 1, RecImageURL: "xx", RecImageStartTime: 0, RecImageEndTime: 1998603966, Rec: true, RecFlag: true}
	r2 = &model.Recommend{ArticleID: 2, RecImageURL: "xx", RecImageStartTime: 0, RecImageEndTime: 1398603966, Rec: true}
	r3 = &model.Recommend{ArticleID: 3, Rec: true}
	r4 = &model.Recommend{ArticleID: 4, Rec: true}
	rs = [][]*model.Recommend{
		[]*model.Recommend{r1},
		[]*model.Recommend{r2},
		[]*model.Recommend{r3},
		[]*model.Recommend{r4},
	}
	cid           = int64(4)
	recommendAids = map[int64][]int64{
		0:   []int64{r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID},
		cid: []int64{r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID},
	}
)

func Test_Recommends_1(t *testing.T) {
	Convey("get data from page 1", t, WithCleanCache(func() {
		s.setting.ShowRecommendNewArticles = true
		//s.updateNewArts(context.TODO(), cid)
		s.RecommendsMap = map[int64][][]*model.Recommend{cid: rs}
		res, err := s.Recommends(context.TODO(), cid, 1, 3, []int64{}, model.FieldDefault)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 3)
		// 不改变原始值
		So(r1, ShouldResemble, &model.Recommend{ArticleID: 1, RecImageURL: "xx", RecImageStartTime: 0, RecImageEndTime: 1998603966, Rec: true, RecFlag: true})
		So(len(s.RecommendsMap[cid]), ShouldEqual, 4)
		So(res[0].Recommend, ShouldResemble, model.Recommend{ArticleID: 0, RecImageURL: "xx", RecImageStartTime: 0, RecImageEndTime: 1998603966, Rec: true, RecFlag: true, RecText: "编辑推荐"})
		So(res[0].ID, ShouldEqual, 1)
		So(res[1].Recommend, ShouldResemble, model.Recommend{ArticleID: 0, RecImageURL: "", RecImageStartTime: 0, RecImageEndTime: 1398603966, Rec: true, RecText: ""})
		So(res[1].ID, ShouldEqual, 2)
		So(res[2].ID, ShouldEqual, 3)
	}))

	Convey("get data from page 1 with aids", t, WithCleanCache(func() {
		s.setting.ShowRecommendNewArticles = true
		//s.updateNewArts(context.TODO(), cid)
		s.RecommendsMap = map[int64][][]*model.Recommend{cid: rs}
		res, err := s.Recommends(context.TODO(), cid, 1, 2, []int64{1, 2}, model.FieldDefault)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 2)
		So(res[0].ID, ShouldEqual, 3)
		So(res[1].ID, ShouldEqual, 4)
	}))
}
func Test_Recommends_2(t *testing.T) {
	Convey("get data from page 2", t, WithCleanCache(func() {
		//s.updateNewArts(context.TODO(), cid)
		s.RecommendsMap = map[int64][][]*model.Recommend{cid: rs}
		s.recommendAids = recommendAids
		Convey("show new art", func() {
			s.setting.ShowRecommendNewArticles = true
			res, err := s.Recommends(context.TODO(), cid, 2, 3, []int64{}, model.FieldDefault)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 3)
			So(res[0].ID, ShouldEqual, 4)
			So(res[1].ID, ShouldNotBeIn, r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID)
			So(res[2].ID, ShouldNotBeIn, r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID)
		})
	}))
}

func Test_Recommends_recomend_category(t *testing.T) {
	Convey("get data from page 1", t, WithCleanCache(func() {
		//s.updateNewArts(context.TODO(), 0)
		rss := [][]*model.Recommend{[]*model.Recommend{r1, r2, r3, r4}}
		s.RecommendsMap = map[int64][][]*model.Recommend{0: rss}
		s.recommendAids = recommendAids
		Convey("show new art", func() {
			s.setting.ShowRecommendNewArticles = true
			res, err := s.Recommends(context.TODO(), 0, 2, 3, []int64{}, model.FieldDefault)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 3)
			So(res[0].ID, ShouldBeIn, r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID)
			So(res[1].ID, ShouldNotBeIn, r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID)
			So(res[2].ID, ShouldNotBeIn, r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID)
		})
		Convey("hide new art", func() {
			s.setting.ShowRecommendNewArticles = false
			res, err := s.Recommends(context.TODO(), 0, 2, 3, []int64{}, model.FieldDefault)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 1)
			So(res[0].ID, ShouldBeIn, r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID)
		})
	}))
}
func Test_Recommends_3(t *testing.T) {
	Convey("get data from page 3", t, WithCleanCache(func() {
		//s.updateNewArts(context.TODO(), cid)
		s.RecommendsMap = map[int64][][]*model.Recommend{cid: rs}
		s.recommendAids = recommendAids
		Convey("show new art", func() {
			s.setting.ShowRecommendNewArticles = true
			res, err := s.Recommends(context.TODO(), cid, 3, 3, []int64{}, model.FieldDefault)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, 3)
			So(res[0].ID, ShouldNotBeIn, r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID)
			So(res[1].ID, ShouldNotBeIn, r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID)
			So(res[2].ID, ShouldNotBeIn, r1.ArticleID, r2.ArticleID, r3.ArticleID, r4.ArticleID)
		})
	}))
	Convey("other category no data", t, WithCleanCache(func() {
		//s.updateNewArts(context.TODO(), cid)
		s.RecommendsMap = map[int64][][]*model.Recommend{cid: rs}
		res, err := s.Recommends(context.TODO(), 100, 1, 10, []int64{}, model.FieldDefault)
		So(err, ShouldNotBeNil)
		So(res, ShouldBeNil)
	}))
}

func Test_CalculateRecommends(t *testing.T) {
	r10 := &model.Recommend{ArticleID: 1, Position: 2, EndTime: 1}
	r12 := &model.Recommend{ArticleID: 2, Position: 2, EndTime: 2}
	r13 := &model.Recommend{ArticleID: 2, Position: 2, EndTime: 0}
	r14 := &model.Recommend{ArticleID: 2, Position: 2, EndTime: 0}
	r20 := &model.Recommend{ArticleID: 1, Position: 1, EndTime: 1}
	Convey("diffrent position", t, func() {
		res := calculateRecommends([]*model.Recommend{r20, r10})
		exp := [][]*model.Recommend{[]*model.Recommend{r10}, []*model.Recommend{r20}}
		So(res, ShouldResemble, exp)
	})
	Convey("same position", t, func() {
		res := calculateRecommends([]*model.Recommend{r12, r10})
		exp1 := [][]*model.Recommend{[]*model.Recommend{r10, r12}}
		exp2 := [][]*model.Recommend{[]*model.Recommend{r12, r10}}
		So(res, ShouldBeIn, exp1, exp2)
	})
	Convey("one no endtime", t, func() {
		res := calculateRecommends([]*model.Recommend{r13, r10, r20})
		exp := [][]*model.Recommend{[]*model.Recommend{r10}, []*model.Recommend{r20}}
		So(res, ShouldResemble, exp)
	})
	Convey("all no endtime", t, func() {
		res := calculateRecommends([]*model.Recommend{r13, r14})
		exp1 := [][]*model.Recommend{[]*model.Recommend{r13, r14}}
		exp2 := [][]*model.Recommend{[]*model.Recommend{r14, r13}}
		So(res, ShouldBeIn, exp1, exp2)
	})

	Convey("no endtime and have endtime ", t, func() {
		res := calculateRecommends([]*model.Recommend{r13, r14, r12})
		exp := [][]*model.Recommend{[]*model.Recommend{r12}}
		So(res, ShouldResemble, exp)
	})
}

func Test_DelRecommendArt(t *testing.T) {
	Convey("del recommend", t, WithService(func(s *Service) {
		s.RecommendsMap = map[int64][][]*model.Recommend{0: rs}
		So(s.RecommendsMap, ShouldNotBeNil)
		So(len(s.RecommendsMap[0]), ShouldEqual, 4)
		s.DelRecommendArt(0, 1)
		time.Sleep(50 * time.Millisecond)
		So(s.RecommendsMap[0][0][0], ShouldResemble, r2)
	}))
}

func Test_genRecommendArtFromPool(t *testing.T) {
	Convey("should generate arts", t, WithService(func(s *Service) {
		res := s.genRecommendArtFromPool([][]*model.Recommend{[]*model.Recommend{r1, r2, r3, r4}}, s.c.Article.RecommendRegionLen)
		So(len(res), ShouldEqual, 4)
	}))
}

func Test_sortRecs(t *testing.T) {
	Convey("should sort recommends by ptime", t, WithCleanCache(func() {
		a1 := &model.RecommendArt{Meta: model.Meta{ID: 1, PublishTime: 1}}
		a1.Rec = true
		a2 := &model.RecommendArt{Meta: model.Meta{ID: 2, PublishTime: 2}}
		a2.Rec = true
		a3 := &model.RecommendArt{Meta: model.Meta{ID: 3, PublishTime: 3}}
		a3.Rec = true
		res := []*model.RecommendArt{a1, a3, a2}
		sortRecs(res)
		So(res, ShouldResemble, []*model.RecommendArt{a3, a2, a1})
	}))
}

func Test_skyHorseGray(t *testing.T) {
	Convey("mid", t, func() {
		s.c.Article.SkyHorseGray = []int64{}
		s.c.Article.SkyHorseGrayUsers = []int64{123}
		So(s.skyHorseGray("1", 123), ShouldBeTrue)
		So(s.skyHorseGray("", 12), ShouldBeFalse)
		So(s.skyHorseGray("", 0), ShouldBeFalse)
		So(s.skyHorseGray("1", 0), ShouldBeFalse)
	})
	Convey("gray", t, func() {
		s.c.Article.SkyHorseGray = []int64{3}
		s.c.Article.SkyHorseGrayUsers = []int64{}
		So(s.skyHorseGray("1", 123), ShouldBeTrue)
		So(s.skyHorseGray("", 3), ShouldBeTrue)
		So(s.skyHorseGray("", 5), ShouldBeFalse)
		So(s.skyHorseGray("", 0), ShouldBeFalse)
		So(s.skyHorseGray("1", 0), ShouldBeFalse)
	})
}
