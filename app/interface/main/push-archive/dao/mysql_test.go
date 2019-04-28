package dao

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-common/app/interface/main/push-archive/model"

	"github.com/smartystreets/goconvey/convey"
)

func Test_mxID(t *testing.T) {
	_, err := d.SettingsMaxID(context.TODO())
	convey.Convey("获取最大的设置id", t, func() {
		convey.So(err, convey.ShouldBeNil)
	})
}

func Test_settingsall(t *testing.T) {
	res := make(map[int64]*model.Setting)
	start, end := int64(2), int64(3)
	err := d.SettingsAll(context.TODO(), start, end, &res)
	convey.Convey("batch search settings", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(res), convey.ShouldEqual, 0)
	})

	start, end = int64(0), int64(10)
	err = d.SettingsAll(context.TODO(), start, end, &res)
	convey.Convey("batch search settings", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(res), convey.ShouldBeGreaterThan, 0)
	})
}

func Test_statistics(t *testing.T) {
	fans := []int64{1, 2, 3, 4, 5}
	b, err1 := json.Marshal(fans)
	ps := model.PushStatistic{
		Aid:         int64(101),
		Group:       "ai:pushlist_follow_recent",
		Type:        model.StatisticsUnpush,
		Mids:        string(b),
		MidsCounter: len(fans),
		CTime:       time.Now(),
	}
	rows, err := d.SetStatistics(context.TODO(), &ps)
	convey.Convey("添加统计数据", t, func() {
		convey.So(err1, convey.ShouldBeNil)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rows, convey.ShouldEqual, 1)
	})
}
