package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/push/dao/oppo"
	"go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_roundIndex(t *testing.T) {
	Convey("ping redis", t, WithDao(func(d *Dao) {
		err := d.pingRedis(context.Background())
		PromChanLen("a", 1)
		So(err, ShouldBeNil)
		_, err = d.roundIndex(0, 0)
		So(err, ShouldNotBeNil)
		logPushError("asd", 1, []string{"asd"})
	}))

}

func TestPushIPhone(t *testing.T) {
	Convey("push ipad", t, WithDao(func(d *Dao) {
		_, err := d.PushIPhone(context.Background(), &model.PushInfo{}, &model.PushItem{})
		So(err, ShouldNotBeNil)
	}))
}

func TestPushIPad(t *testing.T) {
	Convey("push ipad", t, WithDao(func(d *Dao) {
		_, err := d.PushIPad(context.Background(), &model.PushInfo{}, &model.PushItem{})
		So(err, ShouldNotBeNil)
	}))
}

func TestPushMi(t *testing.T) {
	Convey("push mi", t, WithDao(func(d *Dao) {
		_, err := d.PushMi(context.Background(), &model.PushInfo{}, "", "123")
		So(err, ShouldNotBeNil)
	}))
}

func TestPushMiByMids(t *testing.T) {
	Convey("push mi by mids", t, WithDao(func(d *Dao) {
		_, err := d.PushMiByMids(context.Background(), &model.PushInfo{}, "", "123")
		So(err, ShouldBeNil)
	}))
}

func TestPushHuawei(t *testing.T) {
	Convey("push huawei", t, WithDao(func(d *Dao) {
		_, err := d.PushHuawei(context.Background(), &model.PushInfo{}, "", []string{"123"})
		So(err, ShouldNotBeNil)
	}))
}

func TestPushOppo(t *testing.T) {
	Convey("push oppo", t, WithDao(func(d *Dao) {
		_, err := d.PushOppo(context.Background(), &model.PushInfo{}, "", []string{"123"})
		So(err, ShouldNotBeNil)
	}))
}

func TestPushOppoOne(t *testing.T) {
	Convey("push oppo one", t, WithDao(func(d *Dao) {
		_, err := d.PushOppoOne(context.Background(), &model.PushInfo{}, &oppo.Message{}, "123")
		So(err, ShouldNotBeNil)
	}))
}

func TestPushJpush(t *testing.T) {
	Convey("push jpush", t, WithDao(func(d *Dao) {
		_, err := d.PushJpush(context.Background(), &model.PushInfo{}, "", []string{"123"})
		So(err, ShouldNotBeNil)
	}))
}

func TestPushFCM(t *testing.T) {
	Convey("push fcm", t, WithDao(func(d *Dao) {
		_, err := d.PushFCM(context.Background(), &model.PushInfo{}, "", []string{"123"})
		So(err, ShouldNotBeNil)
	}))
}
