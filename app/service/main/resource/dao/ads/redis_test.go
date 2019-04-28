package ads

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAdsBuvidCache(t *testing.T) {
	var (
		buvidCache = map[string]map[int64]int64{
			"123456": {
				10097289: 1,
				10097290: 5,
			},
			"234567": {
				10097288: 3,
				10098523: 1,
			},
		}
		res int64
		err error
	)
	convey.Convey("AddBuvidCount - add cache", t, WithDao(func(d *Dao) {
		err = d.AddBuvidCount(context.Background(), buvidCache)
		convey.Convey("Error should be nil", func() {
			convey.So(err, convey.ShouldBeNil)
		})
		convey.Convey("BuvidCount - get cache", func() {
			convey.Convey("Case 1: faid = 10097290, buvid = 123456", func() {
				res, err = d.BuvidCount(context.Background(), 10097290, "123456")
				convey.Convey("Error should be nil, res should be 5", func() {
					convey.So(err, convey.ShouldBeNil)
					convey.So(res, convey.ShouldEqual, 5)
				})
			})
			convey.Convey("Case 2: faid = 10097288, buvid = 234567", func() {
				res, err := d.BuvidCount(context.Background(), 10097288, "234567")
				convey.Convey("Error should be nil, res should be 3", func() {
					convey.So(err, convey.ShouldBeNil)
					convey.So(res, convey.ShouldEqual, 3)
				})
			})
		})
		convey.Convey("keyBuvid", func() {
			key := d.keyBuvid("234567")
			convey.Convey("key should not be nil", func() {
				convey.So(key, convey.ShouldNotBeNil)
			})
			convey.Convey("ExistsAuth - expire cache", func() {
				res, err := d.ExistsAuth(context.Background(), key)
				convey.Convey("Error should be nil, res should be true", func() {
					convey.So(err, convey.ShouldBeNil)
					convey.So(res, convey.ShouldEqual, true)
				})
			})
		})
	}))
}
