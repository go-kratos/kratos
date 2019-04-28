package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/space/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_Channel(t *testing.T) {
	convey.Convey("test channel", t, func(ctx convey.C) {
		mid := int64(27515260)
		cid := int64(38)
		res, err := d.Channel(context.Background(), mid, cid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldHaveSameTypeAs, &model.Channel{})
	})
}

func TestDao_ChannelCnt(t *testing.T) {
	convey.Convey("test channel cnt", t, func(ctx convey.C) {
		mid := int64(27515260)
		res, err := d.ChannelCnt(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDao_ChannelExtra(t *testing.T) {
	convey.Convey("test channel extra", t, func(ctx convey.C) {
		mid := int64(27515260)
		cid := int64(38)
		res, err := d.ChannelExtra(context.Background(), mid, cid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldHaveSameTypeAs, &model.ChannelExtra{})
	})
}

func TestDao_EditChannelArc(t *testing.T) {
	convey.Convey("test edit channel arc", t, func(ctx convey.C) {
		mid := int64(27515260)
		cid := int64(37)
		sort := []*model.ChannelArcSort{{Aid: 5462035, OrderNum: 4}, {Aid: 5461964, OrderNum: 3}}
		res, err := d.EditChannelArc(context.Background(), mid, cid, time.Now(), sort)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", res)
	})
}

func TestDao_ChannelList(t *testing.T) {
	convey.Convey("test channel list", t, func(ctx convey.C) {
		mid := int64(27515260)
		res, err := d.ChannelList(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", res)
	})
}

func TestDao_ChannelVideos(t *testing.T) {
	convey.Convey("test channel video", t, func(ctx convey.C) {
		mid := int64(27515260)
		cid := int64(37)
		res1, err1 := d.ChannelVideos(context.Background(), mid, cid, true)
		convey.So(err1, convey.ShouldBeNil)
		convey.Printf("%+v", res1)
		res2, err2 := d.ChannelVideos(context.Background(), mid, cid, false)
		convey.So(err2, convey.ShouldBeNil)
		convey.Printf("%+v", res2)
	})
}
