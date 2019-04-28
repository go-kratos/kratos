package dao

import (
	"context"
	"go-common/app/job/main/figure/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testRedisMid int64 = 15555180
)

func Test_PingRedis(t *testing.T) {
	Convey("ping redis", t, WithDao(func(d *Dao) {
		So(d.PingRedis(context.TODO()), ShouldBeNil)
	}))
}

func Test_SetWaiteUserCache(t *testing.T) {
	Convey("set waite user cache", t, WithDao(func(d *Dao) {
		So(d.SetWaiteUserCache(context.TODO(), testRedisMid, 111), ShouldBeNil)
	}))
}

func Test_AddFigureInfoCache(t *testing.T) {
	Convey("add figure info cache", t, WithDao(func(d *Dao) {
		f := &model.Figure{Mid: testRedisMid}
		So(d.AddFigureInfoCache(context.TODO(), f), ShouldBeNil)
	}))
}
