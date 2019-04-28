package dao

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_RecommendByCategory(t *testing.T) {
	Convey("get data", t, WithMysql(func(d *Dao) {
		res, err := d.RecommendByCategory(context.TODO(), 0)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("no data", t, WithMysql(func(d *Dao) {
		res, err := d.RecommendByCategory(context.TODO(), 1000)
		So(err, ShouldBeNil)
		So(res, ShouldBeEmpty)
	}))
}

func Test_AllRecommendCount(t *testing.T) {
	Convey("get data", t, func() {
		res, err := d.AllRecommendCount(context.TODO(), time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
	})
}

func Test_AllRecommends(t *testing.T) {
	Convey("get data", t, func() {
		res, err := d.AllRecommends(context.TODO(), time.Now(), 1, 5)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
