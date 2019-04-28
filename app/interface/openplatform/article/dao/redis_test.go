package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/openplatform/article/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_pingRedis(t *testing.T) {
	Convey("ping redis", t, WithDao(func(d *Dao) {
		So(d.pingRedis(context.TODO()), ShouldBeNil)
	}))
}

func Test_UppersCache(t *testing.T) {
	var (
		mid  = int64(1)
		mid2 = int64(2)
		now  = time.Now().Unix()
		err  error
		a1   = model.Meta{ID: 1, PublishTime: xtime.Time(now), Author: &model.Author{Mid: mid}}
		a2   = model.Meta{ID: 2, PublishTime: xtime.Time(now - 1), Author: &model.Author{Mid: mid}}
		a3   = model.Meta{ID: 3, PublishTime: xtime.Time(now - 2), Author: &model.Author{Mid: mid2}}
		idsm = map[int64][][2]int64{
			mid:  [][2]int64{[2]int64{a1.ID, int64(a1.PublishTime)}, [2]int64{a2.ID, int64(a2.PublishTime)}},
			mid2: [][2]int64{[2]int64{a3.ID, int64(a3.PublishTime)}},
		}
	)
	Convey("add cache", t, WithDao(func(d *Dao) {
		err = d.AddUpperCaches(context.TODO(), idsm)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.UppersCaches(context.TODO(), []int64{mid, mid2}, 0, 2)
			So(res, ShouldResemble, map[int64][]int64{mid: []int64{1, 2}, mid2: []int64{3}})
			So(err, ShouldBeNil)
		})
		Convey("purge cache", func() {
			err := d.DelUpperCache(context.TODO(), a1.Author.Mid, a1.ID)
			So(err, ShouldBeNil)
		})
		Convey("count cache", func() {
			res, err := d.UpperArtsCountCache(context.TODO(), a1.Author.Mid)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 2)
		})
	}))
}

func Test_RankCache(t *testing.T) {
	var (
		cid  = int64(2)
		list = []*model.Rank{
			&model.Rank{
				Aid:   3,
				Score: 3,
			},
			&model.Rank{
				Aid:   2,
				Score: 2,
			},
			&model.Rank{
				Aid:   1,
				Score: 1,
			},
		}
		rank = model.RankResp{Note: "note", List: list}
		err  error
	)
	Convey("add cache", t, WithDao(func(d *Dao) {
		err = d.AddRankCache(context.TODO(), cid, rank)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.RankCache(context.TODO(), cid)
			So(err, ShouldBeNil)
			So(res, ShouldResemble, rank)
			res, err = d.RankCache(context.TODO(), 1000)
			So(err, ShouldBeNil)
			So(res.List, ShouldBeEmpty)
		})
		Convey("expire cache", func() {
			res, err := d.ExpireRankCache(context.TODO(), cid)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, true)
		})
	}))
}
func Test_HotspotCache(t *testing.T) {
	var (
		id   = int64(1)
		c    = context.TODO()
		arts = [][2]int64{[2]int64{0, -1}, [2]int64{1, 1}, [2]int64{2, 2}, [2]int64{3, 3}, [2]int64{4, 4}, [2]int64{5, 5}}
	)
	Convey("work", t, WithCleanCache(func() {
		ok, err := d.ExpireHotspotArtsCache(c, model.HotspotTypePtime, id)
		So(err, ShouldBeNil)
		So(ok, ShouldBeFalse)

		err = d.AddCacheHotspotArts(context.TODO(), model.HotspotTypePtime, id, arts, true)
		So(err, ShouldBeNil)

		ok, err = d.ExpireHotspotArtsCache(c, model.HotspotTypePtime, id)
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)

		var num int64
		num, err = d.HotspotArtsCacheCount(c, model.HotspotTypePtime, id)
		So(err, ShouldBeNil)
		So(num, ShouldEqual, len(arts))

		Convey("get cache", func() {
			res, err := d.HotspotArtsCache(context.TODO(), model.HotspotTypePtime, id, 0, -1)
			So(err, ShouldBeNil)
			So(res, ShouldResemble, []int64{5, 4, 3, 2, 1, 0})
		})
	}))
}
