package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DelPKGCache(t *testing.T) {
	Convey("return someting", t, func() {
		err := d.DelPKGCache(context.Background(), []int64{22, 11})
		So(err, ShouldBeNil)
	})
}

func Test_DelEquipsCache(t *testing.T) {
	Convey("return someting", t, func() {
		err := d.DelEquipsCache(context.Background(), []int64{22, 11})
		So(err, ShouldBeNil)
	})
}
