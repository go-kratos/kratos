package bws

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-common/app/interface/main/activity/model/bws"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_RawUserPoints(t *testing.T) {
	Convey("test cache user points", t, WithDao(func(d *Dao) {
		bid := int64(3)
		key := "9875fa517967622b"
		data, err := d.RawUserPoints(context.Background(), bid, key)
		So(err, ShouldBeNil)
		bs, _ := json.Marshal(data)
		Println(string(bs))
	}))
}

func TestDao_CacheUserPoints(t *testing.T) {
	Convey("test cache user points", t, WithDao(func(d *Dao) {
		bid := int64(3)
		key := "9875fa517967622b"
		data, err := d.CacheUserPoints(context.Background(), bid, key)
		So(err, ShouldBeNil)
		bs, _ := json.Marshal(data)
		Println(string(bs))
	}))
}

func TestDao_AddCacheUserPoints(t *testing.T) {
	Convey("add user points cache", t, WithDao(func(d *Dao) {
		bid := int64(3)
		key := "9875fa517967622b"
		data := []*bws.UserPoint{
			{ID: 1, Pid: 2, Points: 3, Ctime: xtime.Time(time.Now().Unix())},
			{ID: 2, Pid: 3, Points: 4, Ctime: xtime.Time(time.Now().Unix())},
		}
		err := d.AddCacheUserPoints(context.Background(), bid, data, key)
		So(err, ShouldBeNil)
	}))
}
