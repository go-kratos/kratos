package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	dataID   = int64(175)
	noDataID = int64(100000000)
)

func Test_ArticleContent(t *testing.T) {
	Convey("get data", t, WithMysql(func(d *Dao) {
		res, err := d.ArticleContent(context.TODO(), dataID)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("no data", t, WithDao(func(d *Dao) {
		res, err := d.ArticleContent(context.TODO(), noDataID)
		So(err, ShouldBeNil)
		So(res, ShouldBeEmpty)
	}))
}

func Test_ArticleMeta(t *testing.T) {
	Convey("get data", t, WithMysql(func(d *Dao) {
		res, err := d.ArticleMeta(context.TODO(), dataID)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.PublishTime, ShouldNotEqual, 0)
	}))
	Convey("no data", t, WithDao(func(d *Dao) {
		res, err := d.ArticleMeta(context.TODO(), noDataID)
		So(err, ShouldBeNil)
		So(res, ShouldBeNil)
	}))
}

func Test_ArticleMetas(t *testing.T) {
	Convey("get data", t, WithMysql(func(d *Dao) {
		res, err := d.ArticleMetas(context.TODO(), []int64{dataID})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("no data", t, WithDao(func(d *Dao) {
		res, err := d.ArticleMetas(context.TODO(), []int64{noDataID})
		So(err, ShouldBeNil)
		So(res, ShouldBeEmpty)
	}))
}

func Test_UpperArticleCount(t *testing.T) {
	Convey("get data", t, WithMysql(func(d *Dao) {
		res, err := d.UpperArticleCount(context.TODO(), dataID)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
	}))
	Convey("no data", t, WithDao(func(d *Dao) {
		res, err := d.UpperArticleCount(context.TODO(), _noData)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 0)
	}))
}
