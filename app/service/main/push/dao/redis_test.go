package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PingRedis(t *testing.T) {
	Convey("ping redis", t, WithDao(func(d *Dao) {
		err := d.pingRedis(context.Background())
		So(err, ShouldBeNil)
	}))
}
