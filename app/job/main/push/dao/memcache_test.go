package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Ping(t *testing.T) {
	Convey("ping mc ", t, func() {
		err := d.Ping(context.Background())
		So(err, ShouldBeNil)
	})
}

func Test_ReportsCacheByMids(t *testing.T) {
	Convey("ReportsCacheByMids", t, func() {
		_, _, err := d.ReportsCacheByMids(context.Background(), []int64{0, 1})
		So(err, ShouldBeNil)
	})
}
