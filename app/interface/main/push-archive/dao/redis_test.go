package dao

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-common/app/interface/main/push-archive/model"

	"github.com/smartystreets/goconvey/convey"
)

func Test_upperlimit(t *testing.T) {
	upper := int64(998)
	d.UpperLimitExpire = 1 // 1s
	exist, err := d.ExistUpperLimitCache(context.TODO(), upper)
	convey.Convey("upper主推送频率限制，没存储过", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(exist, convey.ShouldEqual, false)
	})

	err = d.AddUpperLimitCache(context.TODO(), upper)
	convey.Convey("upper主推送频率限制,添加推送1次，再次获取已存在, 失效后不存在", t, func() {
		convey.So(err, convey.ShouldBeNil)
		exist, err = d.ExistUpperLimitCache(context.TODO(), upper)
		convey.So(err, convey.ShouldBeNil)
		convey.So(exist, convey.ShouldEqual, true)
		time.Sleep(2 * time.Second)
		exist, err = d.ExistUpperLimitCache(context.TODO(), upper)
		convey.So(err, convey.ShouldBeNil)
		convey.So(exist, convey.ShouldEqual, false)
	})
}

func Test_statisticscache(t *testing.T) {
	ps, err := d.GetStatisticsCache(context.TODO())
	convey.Convey("从redis获取统计数据, 没有数据", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(ps, convey.ShouldBeNil)
	})

	per := int64(1000)
	start := int64(1000000)
	var mids []int64
	for i := start; i < start+per; i++ {
		mids = append(mids, i)
	}
	midscount := len(mids)
	midsstr, _ := json.Marshal(mids)
	ps = &model.PushStatistic{
		Aid:         int64(121321),
		Group:       "ai:pushlist_offline_up",
		Type:        1,
		Mids:        string(midsstr),
		MidsCounter: midscount,
		CTime:       time.Now(),
	}

	err = d.AddStatisticsCache(context.TODO(), ps)
	convey.Convey("添加统计数据到redis", t, func() {
		convey.So(err, convey.ShouldBeNil)
	})
}

func Test_perupperlimit(t *testing.T) {
	upper := int64(10)
	fan := int64(20)
	err := d.AddPerUpperLimitCache(context.TODO(), fan, upper, 1, 1)
	convey.Convey("添加推送次数限制", t, func() {
		convey.So(err, convey.ShouldEqual, nil)
	})

	total, err := d.GetPerUpperLimitCache(context.TODO(), fan, upper)
	convey.Convey("获取推送次数限制, 失效后不存在", t, func() {
		convey.So(err, convey.ShouldEqual, nil)
		convey.So(total, convey.ShouldEqual, 1)
		time.Sleep(time.Second * 2)
		total, err = d.GetPerUpperLimitCache(context.TODO(), fan, upper)
		convey.So(err, convey.ShouldEqual, nil)
		convey.So(total, convey.ShouldEqual, 0)

	})
}
