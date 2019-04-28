package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DelMiInvalid(t *testing.T) {
	Convey("Test_DelMiInvalid", t, WithDao(func(d *Dao) {
		err := d.DelMiInvalid(context.Background())
		So(err, ShouldBeNil)
	}))
}

func Test_DelMiUninstalled(t *testing.T) {
	Convey("Test_DelMiUninstalled", t, WithDao(func(d *Dao) {
		err := d.DelMiUninstalled(context.Background())
		So(err, ShouldNotBeNil)
	}))
}

func Test_delInvalidMiReports(t *testing.T) {
	Convey("Test_delInvalidMiReports", t, WithDao(func(d *Dao) {
		err := d.delInvalidMiReports(context.Background(), 1, []string{"test"})
		So(err, ShouldBeNil)
	}))
}
