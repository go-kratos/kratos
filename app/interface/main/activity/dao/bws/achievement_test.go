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

func TestDao_CacheAchieveCounts(t *testing.T) {
	Convey("test cache achieve count", t, WithDao(func(d *Dao) {
		bid := int64(3)
		day := "20180712"
		data, err := d.CacheAchieveCounts(context.Background(), bid, day)
		So(err, ShouldBeNil)
		bs, _ := json.Marshal(data)
		Printf("%v", string(bs))
	}))
}

func TestDao_AddCacheAchieveCounts(t *testing.T) {
	Convey("test add cache achieve count", t, WithDao(func(d *Dao) {
		bid := int64(3)
		day := "20180712"
		list := []*bws.CountAchieves{
			{Aid: 111, Count: 222},
			{Aid: 222, Count: 333},
		}
		err := d.AddCacheAchieveCounts(context.Background(), bid, list, day)
		So(err, ShouldBeNil)
	}))
}

func TestDao_AddCacheUserAchieves(t *testing.T) {
	Convey("test add cache", t, WithDao(func(d *Dao) {
		bid := int64(3)
		list := []*bws.UserAchieve{
			{ID: 2, Aid: 3, Award: 0, Ctime: xtime.Time(time.Now().Unix())},
			{ID: 3, Aid: 4, Award: 0, Ctime: xtime.Time(time.Now().Unix())},
		}
		key := "9abf1997abe851e6"
		err := d.AddCacheUserAchieves(context.Background(), bid, list, key)
		So(err, ShouldBeNil)
	}))
}

func TestDao_CacheUserAchieves(t *testing.T) {
	Convey("test cache user achieves", t, WithDao(func(d *Dao) {
		bid := int64(3)
		key := "9abf1997abe851e6"
		data, err := d.CacheUserAchieves(context.Background(), bid, key)
		So(err, ShouldBeNil)
		bs, _ := json.Marshal(data)
		Printf("%v", string(bs))
	}))
}

func TestDao_AddLotteryMidCache(t *testing.T) {
	Convey("test add lottery mid cache", t, WithDao(func(d *Dao) {
		aid := int64(3)
		mid := int64(908085)
		for i := 0; i < 10; i++ {
			err := d.AddLotteryMidCache(context.Background(), aid, mid+int64(i))
			So(err, ShouldBeNil)
		}
	}))
}

func TestDao_LotteryMidCache(t *testing.T) {
	Convey("test get lottery mid cache", t, WithDao(func(d *Dao) {
		aid := int64(3)
		mid, err := d.CacheLotteryMid(context.Background(), aid, "")
		So(err, ShouldBeNil)
		Println(mid)
	}))
}

func TestDao_RawAchieveCounts(t *testing.T) {
	Convey("test achieve count", t, WithDao(func(d *Dao) {
		bid := int64(1)
		day := "20180712"
		data, err := d.RawAchieveCounts(context.Background(), bid, day)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}

func TestDao_RawAchievements(t *testing.T) {
	Convey("test raw achievements", t, WithDao(func(d *Dao) {
		bid := int64(1)
		data, err := d.RawAchievements(context.Background(), bid)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}
