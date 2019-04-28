package dao

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	model "go-common/app/interface/main/shorturl/model"
	xtime "go-common/library/time"
	"testing"
	"time"
)

func TestDao_Short(t *testing.T) {
	Convey("Short", t, WithDao(func(d *Dao) {
		_, err := d.Short(context.TODO(), "http://b23.tv/EbUzmu")
		So(err, ShouldBeNil)
	}))
}

func TestDao_ShortbyID(t *testing.T) {
	Convey("ShortbyID", t, WithDao(func(d *Dao) {
		_, err := d.ShortbyID(context.TODO(), 1)
		So(err, ShouldBeNil)
	}))
}

func TestDao_AllShorts(t *testing.T) {
	Convey("AllShorts", t, WithDao(func(d *Dao) {
		_, err := d.AllShorts(context.TODO())
		So(err, ShouldBeNil)
	}))
}

func TestDao_ShortCount(t *testing.T) {
	Convey("ShortCount", t, WithDao(func(d *Dao) {
		_, err := d.ShortCount(context.TODO(), 1, "http://www.baidu.com")
		So(err, ShouldBeNil)
	}))
}

func TestDao_InShort(t *testing.T) {
	Convey("InShort", t, WithDao(func(d *Dao) {
		su := &model.ShortUrl{
			Long:  "http://www.baidu.com",
			Mid:   279,
			State: model.StateNormal,
			CTime: xtime.Time(time.Now().Unix()),
		}
		_, err := d.InShort(context.TODO(), su)
		So(err, ShouldBeNil)
	}))
}

func TestDao_ShortUp(t *testing.T) {
	Convey("ShortUp", t, WithDao(func(d *Dao) {
		_, err := d.ShortUp(context.TODO(), 1, 20, "http://www.baidu.com")
		So(err, ShouldBeNil)
	}))
}

func TestDao_UpdateState(t *testing.T) {
	Convey("UpdateState", t, WithDao(func(d *Dao) {
		_, err := d.UpdateState(context.TODO(), 1, 279, 0)
		So(err, ShouldBeNil)
	}))
}

func TestDao_ShortLimit(t *testing.T) {
	Convey("ShortLimit", t, WithDao(func(d *Dao) {
		_, err := d.ShortLimit(context.TODO(), 1, 20, 279, "http://www.baidu.com")
		So(err, ShouldBeNil)
	}))
}
