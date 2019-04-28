package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/push-archive/model"

	"github.com/smartystreets/goconvey/convey"
)

func Test_onekey(t *testing.T) {
	var included bool
	var err error
	included, err = d.filterFanByUpper(context.TODO(), int64(12312313), int64(275152561), "ai:pushlist_follow_recent", []string{"m"})
	convey.Convey("hbase过滤up主, 不存在", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(included, convey.ShouldEqual, false)
	})

	included, err = d.filterFanByUpper(context.TODO(), int64(27515303), int64(27515256), "ai:pushlist_follow_recent", []string{"m", "m1"})
	convey.Convey("hbase过滤up主，增加1个", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(included, convey.ShouldEqual, true)
	})
	included, err = d.filterFanByUpper(context.TODO(), int64(27515401), int64(27515256), "ai:pushlist_follow_recent", []string{"m"})
	convey.Convey("hbase过滤up主，增加1个", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(included, convey.ShouldEqual, true)
	})
	included, err = d.filterFanByUpper(context.TODO(), int64(27515300), int64(27515256), "ai:pushlist_follow_recent", []string{"m"})
	convey.Convey("hbase过滤up主，增加1个", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(included, convey.ShouldEqual, true)
	})

}

func Test_keys(t *testing.T) {
	var result, excluded []int64
	params := map[string]interface{}{
		"base":     int64(27515256),
		"table":    "ai:pushlist_follow_recent",
		"family":   []string{"m"},
		"result":   &result,
		"excluded": &excluded,
		"handler":  d.filterFanByUpper,
	}
	err := d.FilterFans(&[]int64{27515303, 27515401, 27515300, 12312313}, params)
	convey.Convey("多协程过滤up主,3个符合，1个排除", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(result), convey.ShouldEqual, 3)
		convey.So(len(excluded), convey.ShouldEqual, 1)
	})
}

func Test_batchfilter(t *testing.T) {
	var result, excluded []int64
	params := model.NewBatchParam(map[string]interface{}{
		"base":     int64(27515256),
		"table":    "ai:pushlist_follow_recent",
		"family":   []string{"m"},
		"result":   &result,
		"excluded": &excluded,
		"handler":  d.filterFanByUpper,
	}, nil)
	Batch(&[]int64{27515303, 27515401, 27515300, 12312313}, 1, 2, params, d.FilterFans)
	convey.Convey("批量过滤up主, ,3个符合，1个排除", t, func() {
		convey.So(len(result), convey.ShouldEqual, 3)
		convey.So(len(excluded), convey.ShouldEqual, 1)
	})
	t.Logf("the result(%v), excluded(%v)", result, excluded)
}

func Test_addfans(t *testing.T) {
	err := d.AddFans(context.TODO(), int64(275152561), int64(121212), model.RelationAttention)
	convey.Convey("添加粉丝到up主", t, func() {
		convey.So(err, convey.ShouldBeNil)
	})
}

func Test_delfans(t *testing.T) {
	err := d.DelFans(context.TODO(), int64(275152561), int64(121212))
	convey.Convey("删除up主的粉丝", t, func() {
		convey.So(err, convey.ShouldBeNil)
	})
}

func Test_fansbyupper(t *testing.T) {
	Test_addfans(t)
	fans, err := d.Fans(context.TODO(), int64(275152561), false)
	convey.Convey("up主增加一个粉丝后", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(fans), convey.ShouldEqual, 1)
	})

	fans, err = d.Fans(context.TODO(), int64(275152561), true)
	convey.Convey("up主增加一个普通关注粉丝后, pgc稿件只有特殊关注粉丝", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(fans), convey.ShouldEqual, 0)
	})

	Test_delfans(t)
	fans, err = d.Fans(context.TODO(), int64(275152561), false)
	convey.Convey("up主删除一个粉丝后", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(fans), convey.ShouldEqual, 0)
	})
}

func Test_fansbyactive(t *testing.T) {
	// 18507659 + 37118721 + 88889069
	fan := int64(88889069)
	hour := 21
	table := "dm_member_push_active_hour"
	family := []string{"p"}
	included, err := d.filterFanByActive(context.TODO(), fan, hour, table, family)
	t.Logf("the included(%v) err(%v)", included, err)
}
