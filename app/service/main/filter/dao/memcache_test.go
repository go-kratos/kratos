package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/filter/model"
)

func TestDaoMCFilterKey(t *testing.T) {
	convey.Convey("MCFilterKey", t, func(ctx convey.C) {
		var (
			area    = "danmu"
			tpid    = int64(0)
			keys    = []string{"cid:3", "hello"}
			content = "23333"
		)
		ctx.Convey("When everything looks good.", func(ctx convey.C) {
			key := mcFilterKey(area, tpid, keys, content)
			ctx.Convey("Then key should not be equal .", func(ctx convey.C) {
				ctx.So(key, convey.ShouldNotBeNil)
				ctx.So(key, convey.ShouldEqual, "f_danmu_0_Y2lkOjN8aGVsbG8=_MjMzMzM=")
			})
		})
	})
}

func TestDaoFilterCache(t *testing.T) {
	convey.Convey("SetFilterCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			area    = "common"
			tpid    = int64(1)
			keys    = []string{"cid:3", "hello"}
			content = "ut_test"
			res     = &model.FilterCacheRes{
				Fmsg:     "filter_*",
				Level:    20,
				TpIDs:    []int64{0, 1},
				HitRules: []string{"hit_test"},
				Limit:    233,
				AI: &model.AiScore{
					Scores:    []float64{0.23},
					Threshold: 0.9888,
					Note:      "test_note",
				},
			}
		)
		ctx.Convey("When SetFilterCache", func(ctx convey.C) {
			err := d.SetFilterCache(c, area, tpid, keys, content, res)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			ctx.Convey("When get FilterCache.", func(ctx convey.C) {
				res2, err := d.FilterCache(c, area, tpid, keys, content)
				ctx.Convey("Then err should be nil.res2 should equal res.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(res2, convey.ShouldResemble, res)
				})
			})
		})
	})
}

func TestDaoMCKey(t *testing.T) {
	convey.Convey("mcKey", t, func(ctx convey.C) {
		var (
			key  = "test_key"
			area = "danmu"
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			res := mcKey(key, area)
			ctx.Convey("Then res should not be empty.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldEqual, "test_key_danmu")
			})
		})
	})
}

func TestDaoKeyAreaCache(t *testing.T) {
	convey.Convey("KeyAreaCache", t, func(ctx convey.C) {
		var (
			c    = context.TODO()
			key  = "test_key"
			area = "common"
			res  = &model.KeyAreaInfo{
				ID:      2333,
				FKID:    3222,
				Key:     key,
				Mode:    1,
				Filter:  "test_filter",
				Level:   20,
				Area:    area,
				State:   1,
				Comment: "test_comment",
			}
		)
		ctx.Convey("When SetKeyAreaCache", func(ctx convey.C) {
			err := d.SetKeyAreaCache(c, key, area, []*model.KeyAreaInfo{res})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			ctx.Convey("When get KeyAreaCache.", func(ctx convey.C) {
				ress, miss, err := d.KeyAreaCache(c, key, []string{area})
				ctx.Convey("Then err should be nil.ress2[0] should equal res.miss should have length 0", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(miss, convey.ShouldHaveLength, 1)
					ctx.So(ress, convey.ShouldHaveLength, 1)
					ctx.So(ress[0], convey.ShouldResemble, res)
				})
			})
			ctx.Convey("When DelKeyAreaCache.", func(ctx convey.C) {
				err := d.DelKeyAreaCache(c, key, area)
				ctx.Convey("Then err should be nil", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
				})
			})
		})
	})
}
