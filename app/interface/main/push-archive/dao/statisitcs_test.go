package dao

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_GetStatisticsIDRange(t *testing.T) {
	Convey("GetStatisticsIDRange", t, func() {
		deadline, _ := time.Parse("2006-01-02 15:04:05", "2018-05-01 00:00:00")
		min, max, err := d.GetStatisticsIDRange(context.TODO(), deadline)
		So(err, ShouldBeNil)
		So(min, ShouldBeLessThanOrEqualTo, max)
	})
}

func TestDao_DelStatisticsByID(t *testing.T) {
	Convey("DelStatisticsByID", t, func() {
		_, err := d.DelStatisticsByID(context.TODO(), 1, 10)
		So(err, ShouldBeNil)
	})
}
