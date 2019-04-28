package push

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_ZrangeList(t *testing.T) {
	Convey("TestDao_ZrangeList", t, WithDao(func(d *Dao) {
		data, err := d.ZrangeList(context.TODO())
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		Printf("%+v", data)
	}))
}

func TestDao_ZRem(t *testing.T) {
	var (
		conn = d.redis.Get(ctx)
		key  = _pushKey
		id   = "999"
	)
	Convey("everything is fine", t, WithDao(func(d *Dao) {
		conn.Do("ZADD", key, id)
		err := d.ZRem(ctx, id)
		So(err, ShouldBeNil)
	}))
}
