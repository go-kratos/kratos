package dao

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewArtsCache(t *testing.T) {
	var (
		err  error
		arts = [][2]int64{
			[2]int64{1, 4},
			[2]int64{2, 5},
			[2]int64{3, 6},
		}
		revids = []int64{4, 3, 2, 1}
		cid    = int64(0)
		field  = 1
	)
	// 1:4,2:5,3:6,4:8
	Convey("add cache", t, func() {
		for _, a := range arts {
			err = d.AddSortCache(context.TODO(), cid, field, a[0], a[1])
			So(err, ShouldBeNil)
		}
		err = d.AddSortCache(context.TODO(), cid, field, 4, 8)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.SortCache(context.TODO(), cid, field, 0, -1)
			So(err, ShouldBeNil)
			So(res, ShouldResemble, revids)
		})
		Convey("expire cache", func() {
			res, err := d.ExpireSortCache(context.TODO(), cid, field)
			So(err, ShouldBeNil)
			So(res, ShouldBeTrue)
		})
		Convey("count cache", func() {
			res, err := d.NewArticleCount(context.TODO(), 5)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 2)
		})
		Convey("delete cache", func() {
			err := d.DelSortCache(context.TODO(), cid, field, 1)
			So(err, ShouldBeNil)
			res, err := d.SortCache(context.TODO(), cid, field, 0, -1)
			So(err, ShouldBeNil)
			So(res, ShouldResemble, []int64{4, 3, 2})
		})
	})
}
func Test_randomSortTTL(t *testing.T) {
	d.redisSortTTL = 100
	Convey("random ttl should >= 95 && <= 105", t, func() {
		for i := 0; i < 20; i++ {
			ttl := d.randomSortTTL()
			So(ttl, ShouldBeBetween, 95, 105)
			fmt.Printf("%d ", ttl)
		}
	})
}
