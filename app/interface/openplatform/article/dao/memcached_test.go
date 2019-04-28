package dao

import (
	"context"
	"go-common/app/interface/openplatform/article/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ArticlesCache(t *testing.T) {
	c := context.TODO()
	Convey("add cache", t, WithDao(func(d *Dao) {
		err := d.AddArticlesMetaCache(c, art.Meta)
		So(err, ShouldBeNil)
		err = d.AddArticleContentCache(c, art.ID, art.Content)
		So(err, ShouldBeNil)
		err = d.AddArticleStatsCache(c, art.ID, art.Stats)
		So(err, ShouldBeNil)

		Convey("get meta cache", func() {
			_, err = d.ArticleMetaCache(c, art.ID)
			So(err, ShouldBeNil)

			cached, missed, err1 := d.ArticlesMetaCache(c, []int64{art.ID})
			So(err1, ShouldBeNil)
			So(missed, ShouldBeEmpty)
			So(cached, ShouldResemble, map[int64]*model.Meta{art.ID: art.Meta})
		})

		Convey("get content cache", func() {
			content, err1 := d.ArticleContentCache(c, art.ID)
			So(err1, ShouldBeNil)
			So(content, ShouldEqual, art.Content)
		})
		Convey("get no filter content cache", func() {
			err = d.AddArticleContentCache(c, art.ID, art.Content)
			So(err, ShouldBeNil)
			content, err := d.ArticleContentCache(c, art.ID)
			So(err, ShouldBeNil)
			So(content, ShouldEqual, art.Content)
		})

		Convey("get stats cache", func() {
			cached, missed, err := d.ArticlesStatsCache(c, []int64{art.ID})
			So(err, ShouldBeNil)
			So(missed, ShouldBeEmpty)
			So(cached, ShouldResemble, map[int64]*model.Stats{art.ID: art.Stats})
		})

		Convey("get stat cache", func() {
			res, err := d.ArticleStatsCache(c, art.ID)
			So(err, ShouldBeNil)
			So(res, ShouldResemble, art.Stats)
		})

	}))
}

func Test_AudioCache(t *testing.T) {
	c := context.TODO()
	card := model.AudioCard{ID: 1, Title: "audio"}
	Convey("add cache", t, WithDao(func(d *Dao) {
		err := d.AddAudioCardsCache(c, map[int64]*model.AudioCard{1: &card})
		So(err, ShouldBeNil)
		x, err := d.AudioCardsCache(c, []int64{card.ID})
		So(err, ShouldBeNil)
		So(x[card.ID], ShouldResemble, &card)
	}))
}

func Test_Hotspots(t *testing.T) {
	c := context.TODO()
	hots := []*model.Hotspot{&model.Hotspot{ID: 1, Tag: "tag"}}
	Convey("add cache", t, func() {
		err := d.AddCacheHotspots(c, hots)
		So(err, ShouldBeNil)
		res, err := d.CacheHotspots(c)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, []*model.Hotspot{&model.Hotspot{ID: 1, Tag: "tag", TopArticles: []int64{}}})
		err = d.DelCacheHotspots(c)
		So(err, ShouldBeNil)
		// delete twice
		err = d.DelCacheHotspots(c)
		So(err, ShouldBeNil)
		res, err = d.CacheHotspots(c)
		So(err, ShouldBeNil)
		So(res, ShouldBeNil)
	})
}
