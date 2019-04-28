package dao

import (
	"context"
	"testing"

	"go-common/app/interface/openplatform/article/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ListCache(t *testing.T) {
	c := context.TODO()
	list := &model.List{ID: 1, Name: "name", UpdateTime: xtime.Time(100)}
	Convey("set cache", t, func() {
		err := d.AddCacheList(c, list.ID, list)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.CacheList(c, 1)
			So(err, ShouldBeNil)
			So(res, ShouldResemble, list)
		})
		Convey("get cache not exist", func() {
			res, err := d.CacheList(c, 2000000)
			So(err, ShouldBeNil)
			So(res, ShouldBeNil)
		})
	})
}

func Test_ListArtsCache(t *testing.T) {
	c := context.TODO()
	arts := []*model.ListArtMeta{&model.ListArtMeta{ID: 1, Title: "title", State: 1, PublishTime: xtime.Time(100)}}
	Convey("set cache", t, func() {
		err := d.AddCacheListArts(c, 1, arts)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.CacheListArts(c, 1)
			So(err, ShouldBeNil)
			So(res, ShouldResemble, arts)
		})
		Convey("get cache not exist", func() {
			res, err := d.CacheListArts(c, 20000000)
			So(err, ShouldBeNil)
			So(res, ShouldBeNil)
		})
	})
}

func Test_ListsCache(t *testing.T) {
	c := context.TODO()
	list := &model.List{ID: 1, Name: "name", UpdateTime: xtime.Time(100)}
	list2 := &model.List{ID: 2, Name: "name", UpdateTime: xtime.Time(100)}
	m := map[int64]*model.List{1: list, 2: list2}
	Convey("set cache", t, func() {
		err := d.AddCacheLists(c, m)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.CacheLists(c, []int64{1, 2})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, m)
		})
		Convey("get cache not exist", func() {
			res, err := d.CacheLists(c, []int64{300000})
			So(err, ShouldBeNil)
			So(res, ShouldBeNil)
		})
		Convey("get blank cache", func() {
			res, err := d.CacheLists(c, []int64{})
			So(err, ShouldBeNil)
			So(res, ShouldBeNil)
		})
	})
}

func Test_ListsArtsCache(t *testing.T) {
	c := context.TODO()
	arts := []*model.ListArtMeta{&model.ListArtMeta{ID: 1, Title: "title", State: 1, PublishTime: xtime.Time(100)}}
	arts2 := []*model.ListArtMeta{&model.ListArtMeta{ID: 2, Title: "title", State: 1, PublishTime: xtime.Time(100)}}
	Convey("set cache", t, func() {
		err := d.AddCacheListsArts(c, map[int64][]*model.ListArtMeta{1: arts, 2: arts2})
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.CacheListsArts(c, []int64{1, 2})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, map[int64][]*model.ListArtMeta{1: arts, 2: arts2})
		})
		Convey("get cache not exist", func() {
			res, err := d.CacheListsArts(c, []int64{200000})
			So(err, ShouldBeNil)
			So(res, ShouldBeNil)
		})
	})
}

func Test_SetArticleListCache(t *testing.T) {
	c := context.TODO()
	arts := []*model.ListArtMeta{&model.ListArtMeta{ID: 1, Title: "title", State: 1, PublishTime: xtime.Time(100)}}
	Convey("set cache", t, func() {
		err := d.SetArticleListCache(c, 100, arts)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.ArticleListCache(c, 1)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 100)
		})
		Convey("get cache not exist", func() {
			res, err := d.ArticleListCache(c, 1000000)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 0)
		})
		Convey("multi get", func() {
			res, err := d.CacheArtsListID(c, []int64{1, 100000})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, map[int64]int64{1: 100})
		})
	})
}

func Test_UpListsCache(t *testing.T) {
	c := context.TODO()
	lists := []int64{1}
	mid := int64(1)
	Convey("set cache", t, func() {
		err := d.AddCacheUpLists(c, mid, lists)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.CacheUpLists(c, mid)
			So(err, ShouldBeNil)
			So(res, ShouldResemble, lists)
		})
		Convey("get cache not exist", func() {
			res, err := d.CacheUpLists(c, 2000000)
			So(err, ShouldBeNil)
			So(res, ShouldBeNil)
		})
	})
}

func Test_ListReadCountCache(t *testing.T) {
	c := context.TODO()
	id := int64(1)
	count := int64(100)
	Convey("set cache", t, func() {
		err := d.AddCacheListReadCount(c, id, count)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.CacheListReadCount(c, id)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, count)
		})
		Convey("gets cache", func() {
			res, err := d.CacheListsReadCount(c, []int64{id})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, map[int64]int64{id: count})
		})
		Convey("get cache not exist", func() {
			res, err := d.CacheListReadCount(c, 2000000)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 0)
		})
	})
}
