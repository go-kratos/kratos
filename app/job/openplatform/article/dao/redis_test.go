package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Redis(t *testing.T) {
	Convey("test redis", t, WithDao(func(d *Dao) {
		c := context.TODO()
		conn := d.redis.Get(c)
		defer conn.Close()

		conn.Do("SET", "name", "echo")

		if t, err := redis.Bool(conn.Do("EXISTS", "name")); err != nil {
			fmt.Println(t)
			fmt.Println(err)
			_ = t
		} else {
			fmt.Println(t)
			fmt.Println(err)
			_ = t
		}
		fmt.Println("done")
		err := d.PushStat(c, &StatRetry{})
		So(err, ShouldBeNil)
		res, err := d.PopStat(c)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(d.Intercept(c, 1, 2, ""), ShouldBeFalse)
		So(d.DupViewIntercept(c, 1, 2), ShouldBeTrue)
		So(d.PushStat(c, nil), ShouldBeNil)
		_, err = d.PopStat(c)
		So(err, ShouldBeNil)
		err = d.PushReply(c, 1, 2)
		So(err, ShouldBeNil)
		_, _, err = d.PopReply(c)
		So(err, ShouldBeNil)
		err = d.PushCDN(c, "")
		So(err, ShouldBeNil)
		_, err = d.PopCDN(c)
		So(err, ShouldBeNil)
		err = d.PushArtCache(c, nil)
		So(err, ShouldBeNil)
		_, err = d.PopArtCache(c)
		So(err, ShouldBeNil)

		err = d.PushGameCache(c, nil)
		So(err, ShouldBeNil)
		_, err = d.PopGameCache(c)
		So(err, ShouldBeNil)

		err = d.PushFlowCache(c, nil)
		So(err, ShouldBeNil)
		_, err = d.PopFlowCache(c)
		So(err, ShouldBeNil)

		err = d.PushDynamicCache(c, nil)
		So(err, ShouldBeNil)
		_, err = d.PopDynamicCache(c)
		So(err, ShouldBeNil)
	}))
}
