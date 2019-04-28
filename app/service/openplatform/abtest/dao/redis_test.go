package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPingRedis(t *testing.T) {
	Convey("TestPingRedis: ", t, func() {
		err := d.PingRedis(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestGetRedisVersionID(t *testing.T) {
	var (
		ver int64
		err error
	)
	ctx := context.TODO()
	conn := d.redis.Get(ctx)
	defer conn.Close()
	Convey("TestGetRedisVersionID: ", t, func() {
		if ver, _ = redis.Int64(conn.Do("GET", fmt.Sprintf(_keyVersionID, 1))); ver > 0 {
			_, err = conn.Do("DEL", fmt.Sprintf(_keyVersionID, 1))
		}
		conn.Do("SET", fmt.Sprintf(_keyVersionID, 1), 123)
		ver, err = d.RedisVersionID(ctx, 1)
		So(ver, ShouldEqual, 123)
		So(err, ShouldBeNil)
		_, err = conn.Do("DEL", fmt.Sprintf(_keyVersionID, 1))
	})
}

func TestSetnxRedisVersionID(t *testing.T) {
	var (
		ver int64
		err error
	)
	ctx := context.TODO()
	conn := d.redis.Get(ctx)
	defer conn.Close()
	Convey("TestSetnxRedisVersionID: ", t, func() {
		if ver, _ = redis.Int64(conn.Do("GET", fmt.Sprintf(_keyVersionID, 1))); ver > 0 {
			_, err = conn.Do("DEL", fmt.Sprintf(_keyVersionID, 1))
		}
		err = d.SetnxRedisVersionID(ctx, 1, 146)
		So(err, ShouldBeNil)
		ver, _ = redis.Int64(conn.Do("GET", fmt.Sprintf(_keyVersionID, 1)))
		So(ver, ShouldEqual, 146)
		_, err = conn.Do("DEL", fmt.Sprintf(_keyVersionID, 1))
	})
}

func TestUpdateRedisVersionID(t *testing.T) {
	var (
		ver int64
		err error
	)
	ctx := context.TODO()
	conn := d.redis.Get(ctx)
	defer conn.Close()
	Convey("TestUpdateRedisVersionID: ", t, func() {
		if ver, err = redis.Int64(conn.Do("GET", fmt.Sprintf(_keyVersionID, 1))); err == nil {
			_, err = conn.Do("DEL", fmt.Sprintf(_keyVersionID, 1))
			So(err, ShouldBeNil)
		} else {
			So(err, ShouldEqual, redis.ErrNil)
		}

		ver, err = redis.Int64(conn.Do("GET", fmt.Sprintf(_keyVersionID, 1)))
		fmt.Println(ver)
		So(err, ShouldEqual, redis.ErrNil)

		err = d.UpdateRedisVersionID(ctx, 1, 123)
		So(err, ShouldBeNil)
		ver, err = redis.Int64(conn.Do("GET", fmt.Sprintf(_keyVersionID, 1)))
		So(err, ShouldBeNil)
		So(ver, ShouldEqual, 123)

		err = d.UpdateRedisVersionID(ctx, 1, 132)
		So(err, ShouldBeNil)
		ver, err = redis.Int64(conn.Do("GET", fmt.Sprintf(_keyVersionID, 1)))
		So(ver, ShouldEqual, 132)
		So(err, ShouldBeNil)

		_, err = conn.Do("DEL", fmt.Sprintf(_keyVersionID, 1))
		So(err, ShouldBeNil)
	})
}

func TestEnd(t *testing.T) {
	d.Close()
}
