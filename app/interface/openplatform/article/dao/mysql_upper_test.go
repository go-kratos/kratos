package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpperPassed(t *testing.T) {
	Convey("get data", t, WithMysql(func(d *Dao) {
		aids, err := d.UpperPassed(context.TODO(), dataMID)
		So(err, ShouldBeNil)
		So(aids, ShouldNotBeEmpty)
	}))
	Convey("no data", t, WithMysql(func(d *Dao) {
		aids, err := d.UpperPassed(context.TODO(), noDataMID)
		So(err, ShouldBeNil)
		So(aids, ShouldBeEmpty)
	}))
}

func Test_UppersPassed(t *testing.T) {
	Convey("get data", t, WithMysql(func(d *Dao) {
		arts, err := d.UppersPassed(context.TODO(), []int64{dataMID})
		So(err, ShouldBeNil)
		So(arts, ShouldNotBeEmpty)
	}))
	Convey("no data", t, WithMysql(func(d *Dao) {
		arts, err := d.UppersPassed(context.TODO(), []int64{noDataMID})
		So(err, ShouldBeNil)
		So(arts[noDataMID], ShouldBeEmpty)
	}))
}
