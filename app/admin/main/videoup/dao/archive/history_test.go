package archive

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_historyByAID(t *testing.T) {
	convey.Convey("根据aid获取稿件编辑历史", t, WithDao(func(d *Dao) {
		aid := int64(10107879)
		h, err := d.HistoryByAID(context.TODO(), aid, time.Now().Add(720*-1*time.Hour))
		for _, o := range h {
			t.Logf("%+v", o)
		}
		convey.So(err, convey.ShouldBeNil)
	}))
}

func Test_historyByID(t *testing.T) {
	convey.Convey("根据id获取稿件编辑历史", t, WithDao(func(d *Dao) {
		hid := int64(1)
		_, err := d.HistoryByID(context.TODO(), hid)
		convey.So(err, convey.ShouldBeNil)
	}))
}

func Test_videoHistoryByHID(t *testing.T) {
	convey.Convey("根据id获取分p编辑历史", t, WithDao(func(d *Dao) {
		hid := int64(1)
		_, err := d.VideoHistoryByHID(context.TODO(), hid)
		convey.So(err, convey.ShouldBeNil)
	}))
}
