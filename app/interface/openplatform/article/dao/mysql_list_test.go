package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/openplatform/article/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_List(t *testing.T) {
	c := context.TODO()
	list := &model.List{Name: "name", Mid: 100}
	Convey("add list", t, func() {
		id, err := d.CreativeListAdd(c, list.Mid, list.Name, "", "summary", xtime.Time(200), 200)
		So(err, ShouldBeNil)
		So(id, ShouldBeGreaterThan, 0)
		Convey("get list", func() {
			res, err := d.RawList(c, id)
			So(err, ShouldBeNil)
			res.Ctime = 0
			So(res, ShouldResemble, &model.List{Name: "name", Mid: 100, ID: id, Summary: "summary", PublishTime: xtime.Time(200), Words: 200})
		})
		Convey("update time", func() {
			t := time.Now()
			err := d.CreativeListUpdateTime(c, id, t)
			So(err, ShouldBeNil)
			Convey("get list", func() {
				res, err := d.RawList(c, id)
				So(err, ShouldBeNil)
				res.Ctime = 0
				So(res, ShouldResemble, &model.List{Name: "name", Mid: 100, ID: id, UpdateTime: xtime.Time(t.Unix()), PublishTime: xtime.Time(200), Words: 200, Summary: "summary"})
			})
		})
		Convey("update name", func() {
			err := d.CreativeListUpdate(c, id, "new name", "", "summary", xtime.Time(300), 300)
			So(err, ShouldBeNil)
			Convey("get list", func() {
				res, err := d.RawList(c, id)
				So(err, ShouldBeNil)
				res.Ctime = 0
				So(res, ShouldResemble, &model.List{Name: "new name", Mid: 100, ID: id, Summary: "summary", Words: 300, PublishTime: xtime.Time(300)})
			})
		})
		Convey("up list", func() {
			res, err := d.CreativeUpLists(c, list.Mid)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		})
		Convey("del", func() {
			err := d.CreativeListDel(c, id)
			So(err, ShouldBeNil)
			Convey("get list", func() {
				res, err := d.RawList(c, id)
				So(err, ShouldBeNil)
				So(res, ShouldBeNil)
			})
		})
		Convey("del all articles", func() {
			err := d.CreativeListDelAllArticles(c, id)
			So(err, ShouldBeNil)
		})
	})
}
func Test_CreativeListArticlesCount(t *testing.T) {
	c := context.TODO()
	Convey("get count", t, func() {
		res, err := d.CreativeCountArticles(c, 88888929, []int64{25, 38})
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
	})
}

func Test_CreativeListArticles(t *testing.T) {
	c := context.TODO()
	Convey("get data", t, func() {
		res, err := d.CreativeListArticles(c, 8)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	})
}
func Test_CreativeListsArticles(t *testing.T) {
	c := context.TODO()
	Convey("get data", t, func() {
		res, err := d.CreativeListsArticles(c, []int64{8})
		So(err, ShouldBeNil)
		So(len(res[8]), ShouldBeGreaterThan, 0)
	})
}

func Test_CreativeCategoryArticles(t *testing.T) {
	c := context.TODO()
	Convey("get count", t, func() {
		res, err := d.CreativeCategoryArticles(c, 88888929)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_passedListArts(t *testing.T) {
	Convey("get data", t, func() {
		res, err := d.RawListArts(context.TODO(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_CpListArts(t *testing.T) {
	Convey("get data", t, func() {
		res, err := d.ListArts(context.TODO(), 8)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("null data", t, func() {
		res, err := d.ListArts(context.TODO(), 999999999)
		So(err, ShouldBeNil)
		So(res, ShouldBeEmpty)
	})
}

func Test_ArtsList(t *testing.T) {
	Convey("get data", t, func() {
		res, err := d.ArtsList(context.TODO(), []int64{821})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("list not exist", t, func() {
		res, err := d.ArtsList(context.TODO(), []int64{99999})
		So(err, ShouldBeNil)
		So(res, ShouldBeEmpty)
	})
	Convey("list blank", t, func() {
		res, err := d.ArtsList(context.TODO(), []int64{})
		So(err, ShouldBeNil)
		So(res, ShouldBeNil)
	})
}

func Test_AllArtsList(t *testing.T) {
	Convey("get data", t, func() {
		res, err := d.RawAllLists(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
