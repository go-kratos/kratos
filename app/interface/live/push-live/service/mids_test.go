package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/live/push-live/dao"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/cache/redis"
	"math/rand"
)

func TestService_MidFilter(t *testing.T) {
	initd()
	Convey("test mid filter", t, func() {
		var (
			total    int
			midList  map[int64]bool
			business int
			resMids  []int64
			task     *model.ApPushTask
			err      error
			conn     redis.Conn
		)
		business = rand.Intn(99) + 1 // should through all filters
		// init mids input
		total = 10
		midList = make(map[int64]bool, total)
		for i := 0; i < total; i++ {
			mid := rand.Int63()
			midList[mid] = true
		}
		// init task
		task = &model.ApPushTask{
			LinkValue: "test",
		}
		// do mid filter
		resMids = s.midFilter(midList, business, task)
		So(len(resMids), ShouldEqual, total)
		// clean
		conn, err = redis.Dial(s.c.Redis.PushInterval.Proto, s.c.Redis.PushInterval.Addr, s.dao.RedisOption()...)
		So(err, ShouldBeNil)
		for mid := range midList {
			keys := []string{
				dao.GetDailyLimitKey(mid),
				dao.GetIntervalKey(mid),
			}
			for _, key := range keys {
				conn.Do("DEL", key)
			}
		}
	})
}
