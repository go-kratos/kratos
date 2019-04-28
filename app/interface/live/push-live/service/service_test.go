package service

import (
	"context"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/live/push-live/conf"
	"go-common/library/cache/redis"
	"path/filepath"
	"testing"
)

var (
	s        *Service
	targetID int64
)

func initd() {
	dir, _ := filepath.Abs("../cmd/push-live-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func TestService_ConvertStrToInt64(t *testing.T) {
	initd()
	Convey("test convert", t, func() {
		mStr := "1,2,3"
		mInt64 := []int64{
			int64(1), int64(2), int64(3),
		}
		mRes, err := s.convertStrToInt64(mStr)

		So(err, ShouldBeNil)
		So(mRes, ShouldResemble, mInt64)
	})
}

func TestService_limitDecreaseUnique(t *testing.T) {
	initd()
	Convey("test limit decrease request unique", t, func() {
		var (
			err  error
			conn redis.Conn
			key  string
		)
		Convey("test success request", func() {
			key = "test_request_unique"
			conn, err = redis.Dial(s.c.Redis.PushInterval.Proto, s.c.Redis.PushInterval.Addr, s.dao.RedisOption()...)
			So(err, ShouldBeNil)
			err = s.limitDecreaseUnique(key)
			So(err, ShouldBeNil)
			// clean
			conn.Do("DEL", key)
			conn.Close()
		})
	})
}

func TestService_LimitDecrease(t *testing.T) {
	initd()
	Convey("test LimitDecrease service", t, func() {
		var (
			ctx                              = context.Background()
			business, targetID, uuid, midStr string
			err                              error
			conn                             redis.Conn
		)
		Convey("test success", func() {
			business = "111"
			targetID = "123"
			uuid = "test"
			midStr = "1,2,3"
			conn, err = redis.Dial(s.c.Redis.PushInterval.Proto, s.c.Redis.PushInterval.Addr, s.dao.RedisOption()...)
			So(err, ShouldBeNil)
			err = s.LimitDecrease(ctx, business, targetID, uuid, midStr)
			So(err, ShouldBeNil)
			// clean
			key := getUniqueKey(business, targetID, uuid)
			conn.Do("DEL", key)
			conn.Close()
		})
	})
}
